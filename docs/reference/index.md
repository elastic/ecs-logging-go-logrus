---
mapped_pages:
  - https://www.elastic.co/guide/en/ecs-logging/go-logrus/current/intro.html
  - https://www.elastic.co/guide/en/ecs-logging/go-logrus/current/index.html
---

# ECS Logging Go (Logrus) [intro]

ECS loggers are formatter/encoder plugins for your favorite logging libraries. They make it easy to format your logs into ECS-compatible JSON.

The encoder logs in JSON format, relying on the default [logrus.JSONFormatter](https://pkg.go.dev/github.com/sirupsen/logrus#JSONFormatter) internally. It also handles the logging of error fields in [ECS error format](ecs://reference/ecs-error.md).

By default, the following fields are added:

```json
{
  "log.level": "info",
  "@timestamp": "2020-09-13T10:48:03.000Z",
  "message":" some logging info",
  "ecs.version": "1.6.0"
}
```

::::{tip}
Want to learn more about ECS, ECS logging, and other available language plugins? See the [ECS logging guide](ecs-logging://reference/intro.md).
::::


Ready to jump into `ecslogrus`? [Get started](/reference/setup.md).

