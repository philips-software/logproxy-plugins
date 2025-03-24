package main

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"

	"github.com/dip-software/go-dip-api/logging"
)

// SETUP
// Importantly you need to call Run() once you've done what you need
func TestMain(m *testing.M) {
	log = hclog.NewNullLogger()
	os.Exit(m.Run())
}

func TestFilter_Filter(t *testing.T) {
	type fields struct {
		filterList []FilterReplace
	}
	type args struct {
		msg string
	}
	config := []Config{
		{
			Pattern: "[a-z0-9._%+\\-]+(@|%40)[a-z0-9.\\-]+\\.[a-z]{2,4}",
			Replace: "<email obfuscated>",
		},
		{
			Pattern: "identifier(%3D|=)[^&\\s]*",
			Replace: "identifier=*****",
		},
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantLogMessage string
		wantDropped    bool
		wantModified   bool
		wantErr        error
	}{
		{
			"Modified",
			fields{filterList: parse(config)},
			args{msg: base64.StdEncoding.EncodeToString([]byte("dumbmessage lilpotato@potato.com somemessage russet.potatoes%40potato.com foobarmessage bigpotato@potato.com"))},
			"dumbmessage <email obfuscated> somemessage <email obfuscated> foobarmessage <email obfuscated>",
			false,
			true,
			nil,
		},
		{
			"NotModified",
			fields{filterList: parse(config)},
			args{msg: base64.StdEncoding.EncodeToString([]byte("dumbmessage foobarmessage"))},
			"dumbmessage foobarmessage",
			false,
			false,
			nil,
		},
		{
			"Modified",
			fields{filterList: parse(config)},
			args{msg: base64.StdEncoding.EncodeToString([]byte("GET /logging?procedure.identifier=https://www.philips.com/identifiers/ProcedureIdentifier|ACC_20241120110729422 HTTP/1.1"))},
			"GET /logging?procedure.identifier=***** HTTP/1.1",
			false,
			true,
			nil,
		},
		{
			"Modified",
			fields{filterList: parse(config)},
			args{msg: base64.StdEncoding.EncodeToString([]byte("GET /logging?procedure.identifier%3Dhttps://www.philips.com/identifiers/ProcedureIdentifier|ACC_20241120110729422 HTTP/1.1"))},
			"GET /logging?procedure.identifier=***** HTTP/1.1",
			false,
			true,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Filter{
				filterList: tt.fields.filterList,
			}
			log := &logging.Resource{}
			log.LogData.Message = tt.args.msg
			msg, dropped, modified, err := f.Filter(*log)
			decodedMsg, _ := base64.StdEncoding.DecodeString(msg.LogData.Message)
			if err != tt.wantErr {
				t.Errorf("Filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(decodedMsg) != tt.wantLogMessage {
				t.Errorf("FilterReplace.Filter() got = %v, want %v", string(decodedMsg), tt.wantLogMessage)
			}
			if dropped != tt.wantDropped {
				t.Errorf("Filter() got1 = %v, want %v", dropped, tt.wantDropped)
			}
			if modified != tt.wantModified {
				t.Errorf("Filter() got2 = %v, want %v", modified, tt.wantModified)
			}
		})
	}
}
