---
mapped_pages:
  - https://www.elastic.co/guide/en/ecs-logging/go-logrus/current/setup.html
navigation_title: Get started
---

# Get started with ECS Logging Go (Logrus) [setup]


## Step 1: Install [setup-step-1]

Add the package to your `go.mod` file:

```go
require go.elastic.co/ecslogrus master
```


## Step 2: Configure [setup-step-2]

Set up a default logger. For example:

```go
log := logrus.New()
log.SetFormatter(&ecslogrus.Formatter{})
```


## Examples [examples]


### Use structured logging [use-structured-logging]

```go
// Add custom fields.
log.WithError(errors.New("boom!")).WithField("custom", "foo").Info("hello")
```

The example above produces the following log output:

```json
{
  "@timestamp": "2021-01-20T11:12:43.061+0800",
  "custom":"foo",
  "ecs.version": "1.6.0",
  "error": {
    "message": "boom!"
  },
  "log.level": "info",
  "message":"hello"
}
```


### Nest custom fields under "labels" [nest-data-labels]

For complete ECS compliance, custom fields should be nested in a "labels" object.

```go
log := logrus.New()
log.SetFormatter(&ecslogrus.Formatter{
    DataKey: "labels",
})
log.WithError(errors.New("boom!")).WithField("custom", "foo").Info("hello")
```

The example above produces the following log output:

```json
{
  "@timestamp": "2021-01-20T11:12:43.061+0800",
  "ecs.version": "1.6.0",
  "error": {
    "message": "boom!"
  },
  "labels": {
    "custom": "foo"
  },
  "log.level": "info",
  "message":"hello"
}
```


### Report caller information [report-caller]

```go
log := logrus.New()
log.SetFormatter(&ecslogrus.Formatter{})
log.ReportCaller = true
log.Info("hello")
```

The example above produces the following log output:

```json
{
  "@timestamp": "2021-01-20T11:12:43.061+0800",
  "ecs.version": "1.6.0",
  "log.level": "info",
  "log.origin.file.line": 48,
  "log.origin.file.name": "/path/to/example_test.go",
  "log.origin.function": "go.elastic.co/ecslogrus_test.ExampleFoo",
  "message":"hello"
}
```


## Step 3: Configure Filebeat [setup-step-3]

:::::::{tab-set}

::::::{tab-item} Log file
1. Follow the [Filebeat quick start](beats://reference/filebeat/filebeat-installation-configuration.md)
2. Add the following configuration to your `filebeat.yaml` file.

For Filebeat 7.16+

```yaml
filebeat.inputs:
- type: filestream <1>
  paths: /path/to/logs.json
  parsers:
    - ndjson:
      overwrite_keys: true <2>
      add_error_key: true <3>
      expand_keys: true <4>

processors: <5>
  - add_host_metadata: ~
  - add_cloud_metadata: ~
  - add_docker_metadata: ~
  - add_kubernetes_metadata: ~
```

1. Use the filestream input to read lines from active log files.
2. Values from the decoded JSON object overwrite the fields that {{filebeat}} normally adds (type, source, offset, etc.) in case of conflicts.
3. {{filebeat}} adds an "error.message" and "error.type: json" key in case of JSON unmarshalling errors.
4. {{filebeat}} will recursively de-dot keys in the decoded JSON, and expand them into a hierarchical object structure.
5. Processors enhance your data. See [processors](beats://reference/filebeat/filtering-enhancing-data.md) to learn more.


For Filebeat < 7.16

```yaml
filebeat.inputs:
- type: log
  paths: /path/to/logs.json
  json.keys_under_root: true
  json.overwrite_keys: true
  json.add_error_key: true
  json.expand_keys: true

processors:
- add_host_metadata: ~
- add_cloud_metadata: ~
- add_docker_metadata: ~
- add_kubernetes_metadata: ~
```
::::::

::::::{tab-item} Kubernetes
1. Make sure your application logs to stdout/stderr.
2. Follow the [Run Filebeat on Kubernetes](beats://reference/filebeat/running-on-kubernetes.md) guide.
3. Enable [hints-based autodiscover](beats://reference/filebeat/configuration-autodiscover-hints.md) (uncomment the corresponding section in `filebeat-kubernetes.yaml`).
4. Add these annotations to your pods that log using ECS loggers. This will make sure the logs are parsed appropriately.

```yaml
annotations:
  co.elastic.logs/json.overwrite_keys: true <1>
  co.elastic.logs/json.add_error_key: true <2>
  co.elastic.logs/json.expand_keys: true <3>
```

1. Values from the decoded JSON object overwrite the fields that {{filebeat}} normally adds (type, source, offset, etc.) in case of conflicts.
2. {{filebeat}} adds an "error.message" and "error.type: json" key in case of JSON unmarshalling errors.
3. {{filebeat}} will recursively de-dot keys in the decoded JSON, and expand them into a hierarchical object structure.
::::::

::::::{tab-item} Docker
1. Make sure your application logs to stdout/stderr.
2. Follow the [Run Filebeat on Docker](beats://reference/filebeat/running-on-docker.md) guide.
3. Enable [hints-based autodiscover](beats://reference/filebeat/configuration-autodiscover-hints.md).
4. Add these labels to your containers that log using ECS loggers. This will make sure the logs are parsed appropriately.

```yaml
labels:
  co.elastic.logs/json.overwrite_keys: true <1>
  co.elastic.logs/json.add_error_key: true <2>
  co.elastic.logs/json.expand_keys: true <3>
```

1. Values from the decoded JSON object overwrite the fields that {{filebeat}} normally adds (type, source, offset, etc.) in case of conflicts.
2. {{filebeat}} adds an "error.message" and "error.type: json" key in case of JSON unmarshalling errors.
3. {{filebeat}} will recursively de-dot keys in the decoded JSON, and expand them into a hierarchical object structure.
::::::

:::::::
For more information, see the [Filebeat reference](beats://reference/filebeat/configuring-howto-filebeat.md).

