package main

import (
	"encoding/base64"
	"regexp"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/philips-software/logproxy/shared"

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
		decodedMsg, _ := base64.StdEncoding.DecodeString(msg.LogData.Message)
		if req := filter.pattern.FindAllStringSubmatch(string(decodedMsg), -1); req != nil {
			modifiedMsg := string(decodedMsg)
			for i := range req {
				for j := range req[i] {
					modifiedMsg = strings.ReplaceAll(modifiedMsg, req[i][j], filter.replace)
				}
			}
			msg.LogData.Message = base64.StdEncoding.EncodeToString([]byte(modifiedMsg))
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
