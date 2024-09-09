FROM golang:1.23.1-alpine3.19 as build_base
RUN apk add --no-cache git openssh gcc musl-dev
WORKDIR /plugin
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
LABEL builder=true

# Build
FROM build_base AS builder
WORKDIR /plugin
COPY . .
RUN go build .

FROM ghcr.io/philips-software/logproxy:v1.9.0
RUN apk --no-cache add ca-certificates
COPY --from=builder /plugin/logproxy-filter-replace /app

