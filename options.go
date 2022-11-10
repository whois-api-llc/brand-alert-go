package brandalert

import (
	"time"
)

// Option adds parameters to the query.
type Option func(v *brandAlertRequest)

var _ = []Option{
	OptionResponseFormat("JSON"),
	OptionSinceDate(time.Time{}),
	OptionWithTypos(true),
	OptionPunycode(true),
}

// OptionResponseFormat sets Response output format json | xml. Default: json.
func OptionResponseFormat(outputFormat string) Option {
	return func(v *brandAlertRequest) {
		v.ResponseFormat = outputFormat
	}
}

// OptionSinceDate results in search through activities discovered since the given date.
func OptionSinceDate(date time.Time) Option {
	return func(v *brandAlertRequest) {
		v.SinceDate = date.Format(dateFormat)
	}
}

// OptionWithTypos sets the withTypos option.
// If true, the search terms set will be enriched with their possible typos. Default: false.
func OptionWithTypos(withTypos bool) Option {
	return func(v *brandAlertRequest) {
		v.WithTypos = withTypos
	}
}

// OptionPunycode sets the punicode option.
// If true, domain names in the response will be encoded to Punycode. Default: true.
func OptionPunycode(punycode bool) Option {
	return func(v *brandAlertRequest) {
		v.Punycode = punycode
	}
}
