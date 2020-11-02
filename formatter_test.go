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
	"errors"
	"io/ioutil"
	"runtime"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.elastic.co/ecslogrus"
)

func TestFormatter(t *testing.T) {
	var buf bytes.Buffer
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.SetOutput(&buf)
	log.SetFormatter(&ecslogrus.Formatter{
		DataKey:          "labels",
		CallerPrettyfier: func(*runtime.Frame) (function string, file string) { return "function_name", "file_name" },
	})
	log.ReportCaller = true

	epoch := time.Unix(0, 0).UTC()
	err := errors.New("oy vey")
	log.WithTime(epoch).WithError(err).WithField("custom", "field").Error("oh noes")
	assert.Equal(t,
		`{"@timestamp":"1970-01-01T00:00:00.000Z","ecs.version":"1.6.0","error":{"message":"oy vey"},"labels":{"custom":"field"},"log.level":"error","log.origin.file.name":"file_name","log.origin.function":"function_name","message":"oh noes"}`+"\n",
		buf.String(),
	)
}

func BenchmarkFormatter(b *testing.B) {
	newLogger := func(formatter logrus.Formatter) *logrus.Logger {
		log := logrus.New()
		log.SetLevel(logrus.DebugLevel)
		log.SetOutput(ioutil.Discard)
		log.SetFormatter(formatter)
		log.ReportCaller = true
		return log
	}

	test := func(b *testing.B, f func(b *testing.B, logger *logrus.Logger)) {
		b.Run("json", func(b *testing.B) {
			f(b, newLogger(&logrus.JSONFormatter{
				DataKey:          "labels",
				CallerPrettyfier: func(*runtime.Frame) (function string, file string) { return "function_name", "file_name" },
			}))
		})
		b.Run("ecs", func(b *testing.B) {
			f(b, newLogger(&ecslogrus.Formatter{
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
