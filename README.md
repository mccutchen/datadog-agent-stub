# datadog-agent-stub

A lightweight, no-op substitute for the datadog-agent in development
environments, useful when you want your applications to connect to real statsd
(UDP) and APM trace (HTTP) endpoints but don't want to set up a fully
functional datadog-agent process.

This stub agent logs nothing by default, but when run with `VERBOSE=true` it
logs the statsd metrics and HTTP requests it receives.

## Usage

### Configuration

Configure datadog-agent-stub using environment variables:

| Env var | Default | Description |
| ------- | ------- | ------------|
| `APM_ADDR` | :8126 | Address on which to listen for APM traces, in IP:PORT form |
| `STATSD_ADDR` | :8125 | Address on which to listen for statsd metrics, in IP:PORT form |
| `VERBOSE` | false | Whether to log incoming metrics and traces |

### Example

An example Docker Compose file that demonstrates how to use the
datadog-agent-stub container as a replacement for a production datadog-agent
deployment during local development:

```yaml
version: "3"
services:
  foo-api:
    image: foo-api:latest
    environment:
      DATADOG_TRACE_ENABLED: "true"
      DD_AGENT_HOST: "datadog"
      DD_DOGSTATSD_PORT: "8125"
      DD_SERVICE: "foo-api"
      DD_TAGS: "env:dev"
      DD_TRACE_AGENT_PORT: "8126"
  datadog:
    image: mccutchen/datadog-agent-stub:latest
    environment:
      VERBOSE: "true" # enable logging of statsd metrics and http requests
    ports:
    - "8125/udp"
    - "8126/tcp"
```
