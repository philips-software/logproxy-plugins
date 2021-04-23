package main

import (
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/philips-software/go-hsdp-api/logging"
	"github.com/philips-software/logproxy/shared"
)

var (
	log = hclog.Default()
)

type TriggerReplace struct {
	pattern *regexp.Regexp
	replace string
}

func (f TriggerReplace) Filter(msg logging.Resource) (logging.Resource, bool, bool, error) {
	modified := false
	dropped := false
	if req := f.pattern.FindAllStringSubmatch(msg.LogData.Message, -1); req != nil {
		for i := range req {
			for j := range req[i] {
				msg.LogData.Message = strings.ReplaceAll(msg.LogData.Message, req[i][j], f.replace)
			}
		}
		modified = true
	}
	return msg, dropped, modified, nil
}

func main() {
	filter := &TriggerReplace{}
	reg := os.Getenv("FILTER_REGEXP")
	pattern, err := regexp.Compile(reg)
	if err != nil {
		log.Error("failed to compile FILTER_REGEXP", "regexp", reg)
		return
	}
	filter.pattern = pattern
	filter.replace = os.Getenv("REPLACE_STRING")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"filter": &shared.FilterGRPCPlugin{Impl: filter},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
