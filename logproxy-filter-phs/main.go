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
	requestUsersAPIPattern = regexp.MustCompile(`/api/user[s]?/(?P<userID>[^?/\s]+)`)
	log = hclog.Default()
)

type PHSFilter struct{}

func (PHSFilter) Filter(msg logging.Resource) (logging.Resource, bool, bool, error) {
	modified := false
	// If any of the dropPatterns matches return drop
	for i, r := range dropPatterns {
		if r.MatchString(msg.LogData.Message) {
			log.Info("message dropped", "regex", i)
			return msg, true, false, nil
		}
	}
	// Enhance the log if a userUUID is detected in a request path
	if req := requestUsersAPIPattern.FindStringSubmatch(msg.LogData.Message); req != nil {
		log.Info("user api pattern", "uuid", req[1])
		msg.OriginatingUser = req[1]
		modified = true
	}
	return msg, false, modified, nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"filter": &shared.FilterGRPCPlugin{Impl: &PHSFilter{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}