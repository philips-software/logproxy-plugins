package main

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/philips-software/go-hsdp-api/logging"
	"github.com/philips-software/logproxy/shared"
	"regexp"
)

var (
	dropPatterns = []*regexp.Regexp{
		regexp.MustCompile(`Consul Health Check`),
		regexp.MustCompile(`GET /api/version`),
		regexp.MustCompile(`POST /syslog/drain`),
	}
)

type DropFilter struct{}

func (DropFilter) Filter(msg logging.Resource) (logging.Resource, bool, bool, error) {
	// If any of the dropPatterns matches return drop
	for i, r := range dropPatterns {
		if r.MatchString(msg.LogData.Message) {
			hclog.Fmt("dropping message -- matched %d", i)
			return msg, true, false, nil
		}
	}
	return msg, false, false, nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"filter": &shared.FilterGRPCPlugin{Impl: &DropFilter{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}