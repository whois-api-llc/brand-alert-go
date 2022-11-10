package brandalert

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

// TestOptions tests the Options functions.
func TestOptions(t *testing.T) {
	tests := []struct {
		name   string
		values *brandAlertRequest
		option Option
		want   string
	}{
		{
			name:   "responseFormat",
			values: &brandAlertRequest{},
			option: OptionResponseFormat("json"),
			want:   "json",
		},
		{
			name:   "sinceDate",
			values: &brandAlertRequest{},
			option: OptionSinceDate(time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC)),
			want:   "2021-01-01",
		},
		{
			name:   "withTypos",
			values: &brandAlertRequest{},
			option: OptionWithTypos(true),
			want:   "true",
		},
		{
			name:   "punycode",
			values: &brandAlertRequest{},
			option: OptionPunycode(false),
			want:   "false",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			tt.option(tt.values)

			switch tt.name {
			case "responseFormat":
				got = tt.values.ResponseFormat
			case "sinceDate":
				got = tt.values.SinceDate
			case "withTypos":
				got = strconv.FormatBool(tt.values.WithTypos)
			case "punycode":
				got = strconv.FormatBool(tt.values.WithTypos)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Option() = %v, want %v", got, tt.want)
			}
		})
	}
}
