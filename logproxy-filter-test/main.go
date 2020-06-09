package main

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/philips-software/go-hsdp-api/logging"
	"github.com/philips-software/logproxy/shared"
)

type Filter struct{}

func (f Filter) Filter(msg logging.Resource) (logging.Resource, bool, bool, error) {
	modified := false
	dropped := false
	hclog.Default().Info(fmt.Sprintf("Processed: %s", msg.ID))
	msg.ID = "42"
	modified = true	
	return msg, dropped, modified, nil
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
