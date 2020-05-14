# datadog-agent-stub

A lightweight, no-op substitute for the datadog-agent in development
environments.

## Usage

### Example

Here's an example Docker Compose file that demonstrates how to use the
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
