package main

import (
	"encoding/base64"
	"os"
	"reflect"
	"testing"
)

func Test_get(t *testing.T) {
	var valid = []Config{
		{
			Pattern: "([^&\\s]*)",
			Replace: "<FooBar>",
		},
	}
	const envName = "FILTER_CONFIG"
	tests := []struct {
		name       string
		env        string
		base64json string
		wantRet    []Config
	}{
		{"valid", envName, base64.StdEncoding.EncodeToString([]byte(`[{"pattern": "([^&\\s]*)", "replace": "<FooBar>"}]`)), valid},
		{"invalidJson", envName, base64.StdEncoding.EncodeToString([]byte(`invalid json`)), nil},
		{"invalidBase64", envName, `not base 64`, nil},
	}
	for _, tt := range tests {
		_ = os.Setenv(envName, tt.base64json)
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := get(tt.env); !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("get() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}
