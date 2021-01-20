// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package ecslogrus

import (
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	// ecsVersion holds the version of ECS with which the formatter is compatible.
	ecsVersion = "1.6.0"
)

var (
	ecsFieldMap = logrus.FieldMap{
		logrus.FieldKeyTime:  "@timestamp",
		logrus.FieldKeyMsg:   "message",
		logrus.FieldKeyLevel: "log.level",
	}
)

// Formatter is a logrus.Formatter, formatting log entries as ECS-compliant JSON.
type Formatter struct {
	// DisableHTMLEscape allows disabling html escaping in output
	DisableHTMLEscape bool

	// DataKey allows users to put all the log entry parameters into a
	// nested dictionary at a given key.
	//
	// DataKey is ignored for well-defined fields, such as "error",
	// which will instead be stored under the appropriate ECS fields.
	DataKey string

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the json data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from json fields.
	CallerPrettyfier func(*runtime.Frame) (function string, file string)

	// PrettyPrint will indent all json logs
	PrettyPrint bool
}

// Format formats e as ECS-compliant JSON.
func (f *Formatter) Format(e *logrus.Entry) ([]byte, error) {
	datahint := len(e.Data)
	if f.DataKey != "" {
		datahint = 2
	}
	data := make(logrus.Fields, datahint)
	if len(e.Data) > 0 {
		extraData := data
		if f.DataKey != "" {
			extraData = make(logrus.Fields, len(e.Data))
		}
		for k, v := range e.Data {
			switch k {
			case logrus.ErrorKey:
				err, ok := v.(error)
				if ok {
					data["error"] = errorObject{
						Message: err.Error(),
					}
					break
				}
				fallthrough // error has unexpected type
			default:
				extraData[k] = v
			}
		}
		if f.DataKey != "" && len(extraData) > 0 {
			data[f.DataKey] = extraData
		}
	}
	if e.HasCaller() {
		// Logrus has a single configurable field (logrus.FieldKeyFile)
		// for storing a combined filename and line number, but we want
		// to split them apart into two fields. Remove the event's Caller
		// field, and encode the ECS fields explicitly.
		var funcVal, fileVal string
		var lineVal int
		if f.CallerPrettyfier != nil {
			var fileLineVal string
			funcVal, fileLineVal = f.CallerPrettyfier(e.Caller)
			if sep := strings.IndexRune(fileLineVal, ':'); sep != -1 {
				fileVal = fileLineVal[:sep]
				lineVal, _ = strconv.Atoi(fileLineVal[sep+1:])
			} else {
				fileVal = fileLineVal
				lineVal = 0
			}
		} else {
			funcVal = e.Caller.Function
			fileVal = e.Caller.File
			lineVal = e.Caller.Line
		}
		e.Caller = nil
		if funcVal != "" {
			data["log.origin.function"] = funcVal
		}
		if fileVal != "" {
			data["log.origin.file.name"] = fileVal
		}
		if lineVal > 0 {
			data["log.origin.file.line"] = lineVal
		}
	}
	data["ecs.version"] = ecsVersion
	ecopy := *e
	ecopy.Data = data
	e = &ecopy

	jf := logrus.JSONFormatter{
		TimestampFormat:   "2006-01-02T15:04:05.000Z0700",
		DisableHTMLEscape: f.DisableHTMLEscape,
		FieldMap:          ecsFieldMap,
		CallerPrettyfier:  f.CallerPrettyfier,
		PrettyPrint:       f.PrettyPrint,
	}
	return jf.Format(e)
}

type errorObject struct {
	Message string `json:"message,omitempty"`
}
