package main

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"os"
	"testing"

	"github.com/philips-software/go-hsdp-api/logging"
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
			Pattern: "[a-z0-9._%+\\-]+@[a-z0-9.\\-]+\\.[a-z]{2,4}",
			Replace: "<email obfuscated>",
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
			args{msg: "dumbmessage lilpotato@potato.com foobarmessage bigpotato@potato.com"},
			fmt.Sprintf("dumbmessage %s foobarmessage %s", "<email obfuscated>", "<email obfuscated>"),
			false,
			true,
			nil,
		},
		{
			"NotModified",
			fields{filterList: parse(config)},
			args{msg: "dumbmessage foobarmessage"},
			"dumbmessage foobarmessage",
			false,
			false,
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
			if err != tt.wantErr {
				t.Errorf("Filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if msg.LogData.Message != tt.wantLogMessage {
				t.Errorf("FilterReplace.Filter() got = %v, want %v", msg.LogData.Message, tt.wantLogMessage)
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
