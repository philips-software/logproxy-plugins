package main

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/philips-software/go-hsdp-api/logging"
)

func TestTriggerReplace_Filter(t *testing.T) {
	tests := []struct {
		name           string
		regex          string
		replaceStr     string
		logMessage     string
		wantLogMessage string
		wantDropped    bool
		wantModified   bool
		wantErr        error
	}{
		{
			"Modified",
			"[a-z0-9._%+\\-]+@[a-z0-9.\\-]+\\.[a-z]{2,4}",
			"<email obfuscated>",
			"dumbmessage lilpotato@potato.com foobarmessage bigpotato@potato.com",
			fmt.Sprintf("dumbmessage %s foobarmessage %s", "<email obfuscated>", "<email obfuscated>"),
			false,
			true,
			nil,
		},
		{
			"NotModified",
			"[a-z0-9._%+\\-]+@[a-z0-9.\\-]+\\.[a-z]{2,4}",
			"",
			"dumbmessage foobarmessage",
			"dumbmessage foobarmessage",
			false,
			false,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := &TriggerReplace{}
			reg := tt.regex
			pattern, err := regexp.Compile(reg)
			if err != nil {
				log.Error("failed to compile FILTER_REGEXP", "regexp", reg)
				return
			}
			filter.pattern = pattern
			filter.replace = tt.replaceStr
			logging := &logging.Resource{}
			logging.LogData.Message = tt.logMessage
			got, got1, got2, err := filter.Filter(*logging)
			if got.LogData.Message != tt.wantLogMessage {
				t.Errorf("TriggerReplace.Filter() got = %v, want %v", got.LogData.Message, tt.wantLogMessage)
			}
			if got1 != tt.wantDropped {
				t.Errorf("TriggerReplace.Filter() got1 = %v, want %v", got1, tt.wantDropped)
			}
			if got2 != tt.wantModified {
				t.Errorf("TriggerReplace.Filter() got2 = %v, want %v", got2, tt.wantModified)
			}
			if err != tt.wantErr {
				t.Errorf("TriggerReplace.Filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
