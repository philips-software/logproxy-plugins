package main

import (
	"encoding/base64"
	"os"
	"regexp"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/philips-software/go-hsdp-api/logging"
	"github.com/philips-software/logproxy/shared"
)

var (
	log = hclog.Default()
)

type TriggerDrop struct {
	pattern *regexp.Regexp
}

func (f TriggerDrop) Filter(msg logging.Resource) (logging.Resource, bool, bool, error) {
	drop := false
	decodedMsg, _ := base64.StdEncoding.DecodeString(msg.LogData.Message)
	if req := f.pattern.FindStringSubmatch(string(decodedMsg)); req != nil {
		drop = true
	}
	return msg, drop, false, nil
}

func main() {
	filter := &TriggerDrop{}
	reg := os.Getenv("FILTER_REGEXP")
	pattern, err := regexp.Compile(reg)
	if err != nil {
		log.Error("failed to compile FILTER_REGEXP", "regexp", reg)
		return
	}
	filter.pattern = pattern

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"filter": &shared.FilterGRPCPlugin{Impl: filter},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
