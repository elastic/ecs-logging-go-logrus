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

package ecslogrus_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.elastic.co/ecslogrus"
	"go.elastic.co/ecslogrus/internal/spec"
)

func TestFormatter(t *testing.T) {
	var buf bytes.Buffer
	log := newLogger(&buf, &ecslogrus.Formatter{
		DataKey:          "labels",
		CallerPrettyfier: func(*runtime.Frame) (function string, file string) { return "function_name", "file_name" },
	})

	epoch := time.Unix(0, 0).UTC()
	err := errors.New("oy vey")
	log.WithTime(epoch).WithError(err).WithField("custom", "field").Error("oh noes")
	assert.Equal(t,
		`{"@timestamp":"1970-01-01T00:00:00.000Z","ecs.version":"1.6.0","error":{"message":"oy vey"},"labels":{"custom":"field"},"log.level":"error","log.origin.file.name":"file_name","log.origin.function":"function_name","message":"oh noes"}`+"\n",
		buf.String(),
	)
}

func TestSpecValidation(t *testing.T) {
	t.Run("minimal", func(t *testing.T) {
		var buf bytes.Buffer
		log := newLogger(&buf, &ecslogrus.Formatter{})
		log.ReportCaller = false
		log.Info()

		var decoded map[string]interface{}
		require.NoError(t, json.Unmarshal(buf.Bytes(), &decoded))
		validateSpec(t, decoded, spec.V1)
	})
	t.Run("caller", func(t *testing.T) {
		var buf bytes.Buffer
		log := newLogger(&buf, &ecslogrus.Formatter{
			CallerPrettyfier: func(*runtime.Frame) (function string, file string) { return "function_name", "file_name" },
		})
		log.Info()

		var decoded map[string]interface{}
		require.NoError(t, json.Unmarshal(buf.Bytes(), &decoded))
		validateSpec(t, decoded, spec.V1)
	})
	t.Run("error", func(t *testing.T) {
		var buf bytes.Buffer
		log := newLogger(&buf, &ecslogrus.Formatter{})
		log.WithError(errors.New("oy vey")).Error("oh noes")

		var decoded map[string]interface{}
		require.NoError(t, json.Unmarshal(buf.Bytes(), &decoded))
		validateSpec(t, decoded, spec.V1)
	})
}

func validateSpec(t *testing.T, fields map[string]interface{}, spec *spec.Spec, additionalFields ...string) {
	flattened := make(map[string]interface{})
	flattenMap(fields, flattened, "")
	for name, field := range spec.Fields {
		validateField(t, flattened, name, field)
	}
	for name := range flattened {
		if _, ok := spec.Fields[name]; !ok {
			assert.Contains(t, additionalFields, name, "unexpected field %q", name)
		}
	}
}

func validateField(t *testing.T, fields map[string]interface{}, name string, field spec.Field) {
	val, ok := fields[name]
	if field.Required { // all required fields must be present in the log line
		require.True(t, ok)
		require.NotNil(t, val)
	} else if !ok {
		return
	}
	if field.Type != "" { // the defined type must be met
		var ok bool
		switch field.Type {
		case "string":
			_, ok = val.(string)
		case "datetime":
			var s string
			s, ok = val.(string)
			if _, err := time.Parse("2006-01-02T15:04:05.000Z0700", s); err == nil {
				ok = true
			}
		case "integer":
			// json.Unmarshal unmarshals into float64 instead of int
			if _, floatOK := val.(float64); floatOK {
				_, err := strconv.ParseInt(fmt.Sprintf("%v", val), 10, 64)
				if err == nil {
					ok = true
				}
			}
		default:
			panic(fmt.Errorf("unhandled type %s from specification for field %s", field.Type, name))
		}
		require.True(t, ok, fmt.Sprintf("%s: %v", name, val))
	}
}

func flattenMap(m, flattened map[string]interface{}, prefix string) {
	for k, v := range m {
		switch v := v.(type) {
		case map[string]interface{}:
			flattenMap(v, flattened, prefix+k+".")
		default:
			flattened[prefix+k] = v
		}
	}
}

func BenchmarkFormatter(b *testing.B) {
	test := func(b *testing.B, f func(b *testing.B, logger *logrus.Logger)) {
		b.Run("json", func(b *testing.B) {
			f(b, newLogger(ioutil.Discard, &logrus.JSONFormatter{
				DataKey:          "labels",
				CallerPrettyfier: func(*runtime.Frame) (function string, file string) { return "function_name", "file_name" },
			}))
		})
		b.Run("ecs", func(b *testing.B) {
			f(b, newLogger(ioutil.Discard, &ecslogrus.Formatter{
				DataKey:          "labels",
				CallerPrettyfier: func(*runtime.Frame) (function string, file string) { return "function_name", "file_name" },
			}))
		})
	}

	err := errors.New("oy vey")
	test(b, func(b *testing.B, logger *logrus.Logger) {
		for i := 0; i < b.N; i++ {
			logger.WithError(err).WithField("custom", "field").Error("oh noes")
		}
	})
}

func newLogger(w io.Writer, formatter logrus.Formatter) *logrus.Logger {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.SetOutput(w)
	log.SetFormatter(formatter)
	log.ReportCaller = true
	return log
}
