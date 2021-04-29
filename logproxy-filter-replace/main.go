package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/philips-software/logproxy/shared"
	"regexp"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/philips-software/go-hsdp-api/logging"
)

var (
	log = hclog.Default()
)

type FilterReplace struct {
	pattern *regexp.Regexp
	replace string
}

type Filter struct {
	filterList []FilterReplace
}

func parse(config []Config) (ret []FilterReplace) {
	for _, filterConfig := range config {
		filter := &FilterReplace{}
		compiled, err := regexp.Compile(filterConfig.Pattern)
		if err != nil {
			log.Error("Failed to compile regexp", "regexp", filterConfig.Pattern)
			return
		}
		filter.pattern = compiled
		filter.replace = filterConfig.Replace
		ret = append(ret, *filter)
	}
	return ret
}

func (f Filter) Filter(msg logging.Resource) (logging.Resource, bool, bool, error) {
	modified := false
	for _, filter := range f.filterList {
		if req := filter.pattern.FindAllStringSubmatch(msg.LogData.Message, -1); req != nil {
			for i := range req {
				for j := range req[i] {
					msg.LogData.Message = strings.ReplaceAll(msg.LogData.Message, req[i][j], filter.replace)
				}
			}
			modified = true
		}
	}

	return msg, false, modified, nil
}

func main() {
	filter := &Filter{}
	filter.filterList = parse(get("FILTER_CONFIG"))

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"filter": &shared.FilterGRPCPlugin{Impl: filter},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
