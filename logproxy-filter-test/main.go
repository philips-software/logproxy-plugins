package main

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/philips-software/go-hsdp-api/logging"
	"github.com/philips-software/logproxy/shared"
)

type Filter struct{}

func (Filter) Filter(msg logging.Resource) (logging.Resource, bool, bool, error) {
	hclog.Default().Info(fmt.Sprintf("Processed: %s", msg.ID))
	return msg, false, false, nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"filter": &shared.FilterGRPCPlugin{Impl: &Filter{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
