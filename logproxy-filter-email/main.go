package main

import (
	"fmt"
	_ "github.com/cloudfoundry-community/gautocloud/connectors/smtp"
	"github.com/cloudfoundry-community/gautocloud"
	"github.com/cloudfoundry-community/gautocloud/connectors/smtp/smtptype"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/philips-software/go-hsdp-api/logging"
	"github.com/philips-software/logproxy/shared"
	_ "github.com/philips-software/gautocloud-connectors/hsdp"
	"net/smtp"
	"os"
	"regexp"
)

var (
	log = hclog.Default()
)

type TriggerEmail struct{
	emailTo string
	emailFrom string
	emailSubject string
	pattern *regexp.Regexp
	smtpData smtptype.Smtp
}


func (f TriggerEmail) Filter(msg logging.Resource) (logging.Resource, bool, bool, error) {
	if req := f.pattern.FindStringSubmatch(msg.LogData.Message); req != nil {
		// Send email
		go func() {
			log.Info("Triggering email!")
			auth := smtp.PlainAuth("", f.smtpData.User, f.smtpData.Password, f.smtpData.Host)
			to := []string{f.emailTo}
			msg := []byte("To: "+f.emailTo+"\r\n" +
				"Subject: "+f.emailSubject+"\r\n" +
				"\r\n" +
				"Here's the email triggered by our plugin\r\n" +
				"Timestamp: " + msg.LogTime + "\r\n" +
				"Message: " + msg.LogData.Message + "\r\n")
			err := smtp.SendMail(fmt.Sprintf("%s:%d", f.smtpData.Host, f.smtpData.Port), auth, f.emailFrom, to, msg)
			if err != nil {
				log.Error("error sending email", "error", err.Error())
			}
		}()
	}
	return msg, false, false, nil
}

func main() {
	log.Info("Starting plugin")
	filter := &TriggerEmail{}
	err := gautocloud.Inject(&filter.smtpData)
	if err != nil {
		log.Error("no SMTP service available", "err", err.Error())
		return
	}
	filter.emailTo = os.Getenv("EMAIL_TO")
	filter.emailSubject = os.Getenv("EMAIL_SUBJECT")
	filter.emailFrom = os.Getenv("EMAIL_FROM")
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