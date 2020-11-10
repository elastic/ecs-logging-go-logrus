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
	"errors"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
)

func ExampleFormatter() {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&ecslogrus.Formatter{})

	epoch := time.Unix(0, 0).UTC()
	log.WithTime(epoch).WithField("custom", "foo").Info("hello")

	// Output:
	// {"@timestamp":"1970-01-01T00:00:00.000Z","custom":"foo","ecs.version":"1.6.0","log.level":"info","message":"hello"}
}

func ExampleFormatterError() {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&ecslogrus.Formatter{})

	epoch := time.Unix(0, 0).UTC()
	log.WithTime(epoch).WithError(errors.New("boom!")).Error("an error occurred")

	// Output:
	// {"@timestamp":"1970-01-01T00:00:00.000Z","ecs.version":"1.6.0","error":{"message":"boom!"},"log.level":"error","message":"an error occurred"}
}
