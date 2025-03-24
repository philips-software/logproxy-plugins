FROM golang:1.24.1-alpine3.20 as build_base
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

FROM ghcr.io/philips-software/logproxy:v1.11.0
RUN apk --no-cache add ca-certificates
COPY --from=builder /plugin/logproxy-filter-replace /app

