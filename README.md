# Logproxy plugins
This projects contains a number of [logproxy](https://github.com/philips-software/logproxy) plugins

Plugins are compiled binaries running alongside the main logproxy process which manages the complete lifecycle. 
It detects and loads any plugins found in the search paths. The 
mechanism is based on [Hashicorp's Go plugin system](https://github.com/hashicorp/go-plugin). Some use cases
for writing a plugin:

- Dropping verbose or high frequency logs which are not interesting (saving costs)
- Trigger events (e.g. send an email, call a webhook) when certain patterns in your logs are detected
- Controlled forwarding of logs to other systems

## Plugin interface
Plugins implement a single method  from the `Filter` interface

```go
package main

import (
  "github.com/philips-software/go-hsdp-api/logging"
)

type Filter interface {
    Filter(in logging.Resource) (out logging.Resource, drop bool, modified bool, err error)
}
```

The incoming Resource is a single log message. You can examine its content and decide to leave it as is, drop it, or make changes to various fields.
Setting the `drop` boolean `true` to true will instruct Logproxy to
discard the message without further processing. If the message
contains modifications set the `modified` boolean to `true` as well as 
returning the modified message.

## Note on aggregation in filters
When you want to aggregate data (e.g. a counter) over a number of messages, or trigger events, keep in mind there
might be multiple Logproxy instances running, each with their own copy of the plugin running. 
You will need to use a backing store service, e.g. Redis or PostgreSQL, to synchronize across instances. 

# Building and deploying your plugin
The quickest way to build and deploy your plugin is using a well crafted Dockerfile. We will build the plugin using Docker. 

# Naming
Your plugin binary should follow the `logproxy-filter-*` glob naming convention. In future we might support other types of plugins.

# Dockerfile
We use the official [Logproxy Docker image](https://hub.docker.com/r/philipssoftware/logproxy) as base and simply copy
the plugin binary to the app folder. When the image starts your plugin will be auto-detected. Example:

```Dockerfile
FROM golang:1.14.4-alpine3.11 as build_base
RUN apk add --no-cache git openssh gcc musl-dev
WORKDIR /plugin
COPY go.mod .
COPY go.sum .

# Get plugin dependancies
RUN go mod download
LABEL builder=true

# Build
FROM build_base AS builder
WORKDIR /plugin
COPY . .
RUN go build .

FROM philipssoftware/logproxy:latest
COPY --from=builder /plugin/logproxy-filter-myplugin /app
```

Build and push your Docker image to a registry and you are now
ready to deploy. Please see the [Logproxy project](https://github.com/philips-software/logproxy) for
details on the configuration, specifically the required ENV variables.

# Contact / Getting help

Andy Lo-A-Foe <andy.lo-a-foe@philips.com>

# License

License is MIT
