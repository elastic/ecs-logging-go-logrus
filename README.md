# Elastic Common Schema (ECS) support for Logrus

Use this library for automatically adding a minimal set of ECS fields to your logs, when using [logrus](https://github.com/sirupsen/logrus).
 
---

**Please note** that this library is in a **beta** version and backwards-incompatible changes might be introduced in future releases.
While we strive to comply to [semver](https://semver.org/), we can not guarantee to avoid breaking changes in minor releases.

---
 
The encoder logs in JSON format, relying on the default [logrus.JSONFormatter](https://pkg.go.dev/github.com/sirupsen/logrus#JSONFormatter) internally. 

The following fields will be added by default:
```
{
    "@timestamp":"1970-01-01T0000:00.000Z",
    "ecs.version":"1.6.0",
    "log.level":"info",
    "message":"some logging info"
}
```

The formatter takes care of logging error fields in the [ECS error format](https://www.elastic.co/guide/en/ecs/current/ecs-error.html),
and adding optional caller (file/function) information the [ECS log format](https://www.elastic.co/guide/en/ecs/current/ecs-log.html).
Additional fields will be recorded as ECS "labels".

## What is ECS?

Elastic Common Schema (ECS) defines a common set of fields for ingesting data into Elasticsearch.
For more information about ECS, visit the [ECS Reference Documentation](https://www.elastic.co/guide/en/ecs/current/ecs-reference.html).

## Example usage

### Set up a logger
```go
import (
	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
)

log := logrus.New()
log.SetFormatter(&ecslogrus.Formatter{})
```

### Use structured logging
```go
epoch := time.Unix(0, 0).UTC()
log.WithTime(epoch).WithField("custom", "foo").Info("hello")

// Output:
// {"@timestamp":"1970-01-01T00:00:00.000Z","custom":"foo","ecs.version":"1.6.0","log.level":"info","message":"hello"}
```

### Log errors
```go
epoch := time.Unix(0, 0).UTC()
log.WithTime(epoch).WithError(errors.New("boom!")).Error("an error occurred")

// Output:
// {"@timestamp":"1970-01-01T00:00:00.000Z","ecs.version":"1.6.0","error":{"message":"boom!"},"log.level":"error","message":"an error occurred"}
```

## References
* Introduction to ECS [blog post](https://www.elastic.co/blog/introducing-the-elastic-common-schema).
* Logs UI [blog post](https://www.elastic.co/blog/infrastructure-and-logs-ui-new-ways-for-ops-to-interact-with-elasticsearch).

## Test
```
go test ./...
```

## Contribute
Create a Pull Request from your own fork. 

Run `mage` to update and format you changes before submitting.

Add new dependencies to the NOTICE.txt.

## License
This software is licensed under the [Apache 2 license](./LICENSE). 
