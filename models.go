package brandalert

import (
	"encoding/json"
	"fmt"
	"time"
)

func unmarshalString(raw json.RawMessage) (string, error) {
	var val string
	err := json.Unmarshal(raw, &val)
	if err != nil {
		return "", err
	}
	return val, nil
}

// Time is a helper wrapper on time.Time
type Time time.Time

var emptyTime Time

const dateFormat = "2006-01-02"

// UnmarshalJSON decodes time as Brand Alert API does.
func (t *Time) UnmarshalJSON(b []byte) error {
	str, err := unmarshalString(b)
	if err != nil {
		return err
	}
	if str == "" {
		*t = emptyTime
		return nil
	}
	v, err := time.Parse(dateFormat, str)
	if err != nil {
		return err
	}
	*t = Time(v)
	return nil
}

// MarshalJSON encodes time as Brand Alert API does.
func (t Time) MarshalJSON() ([]byte, error) {
	if t == emptyTime {
		return []byte(`""`), nil
	}
	return []byte(`"` + time.Time(t).Format(dateFormat) + `"`), nil
}

// SearchTerms is a set of including or excluding search terms.
type SearchTerms []string

// brandAlertRequest is the request struct for Brand Alert API.
type brandAlertRequest struct {
	// APIKey is the user's API key.
	APIKey string `json:"apiKey"`

	// IncludeSearchTerms is an array of search terms. All of them should be present in the domain name.
	IncludeSearchTerms *SearchTerms `json:"includeSearchTerms,omitempty"`

	// ExcludeSearchTerms is an array of search terms. All of them should NOT be present in the domain name.
	ExcludeSearchTerms *SearchTerms `json:"excludeSearchTerms,omitempty"`

	// SinceDate If present, search through activities discovered since the given date.
	SinceDate string `json:"sinceDate,omitempty"`

	// Mode is the mode of the API call. Acceptable values: preview | purchase.
	Mode string `json:"mode,omitempty"`

	// WithTypos If true, the search terms set will be enriched with their possible typos.
	WithTypos bool `json:"withTypos,omitempty"`

	// Punycode If true, domain names in the response will be encoded to punycode.
	Punycode bool `json:"punycode,omitempty"`

	// ResponseFormat is the response output format JSON | XML.
	ResponseFormat string `json:"responseFormat,omitempty"`
}

// Action is a wrapper on string.
type Action string

// List of possible actions.
const (
	Added      Action = "added"
	Updated           = "updated"
	Dropped           = "dropped"
	Discovered        = "discovered"
)

var _ = []Action{
	Added,
	Updated,
	Dropped,
	Discovered,
}

// DomainItem is a part of the Brand Alert API response.
type DomainItem struct {
	// DomainName is the full domain name.
	DomainName string `json:"domainName"`

	// Action is the related action. Possible actions: added | updated | dropped | discovered.
	Action Action `json:"action"`

	// Date is the event date.
	Date Time `json:"date"`
}

// BrandAlertResponse is a response of Brand Alert API.
type BrandAlertResponse struct {
	// DomainsList is the list of domains matching the criteria.
	DomainsList []DomainItem `json:"domainsList"`

	// DomainsCount is the number of domains matching the criteria.
	DomainsCount int `json:"domainsCount"`
}

// Messages is a wrapper on []string.
type Messages []string

// UnmarshalJSON decodes the error messages returned by Registrant Alert API.
func (m *Messages) UnmarshalJSON(b []byte) error {
	var msgs []string

	if err := json.Unmarshal(b, &msgs); err == nil {
		*m = msgs
		return nil
	}

	var x interface{}
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}
	*m = append(*m, fmt.Sprintf("%s", x))
	return nil
}

// ErrorMessage is the error message.
type ErrorMessage struct {
	Code    int      `json:"code"`
	Message Messages `json:"messages"`
}

// Error returns error message as a string.
func (e *ErrorMessage) Error() string {
	return fmt.Sprintf("API error: [%d] %s", e.Code, e.Message)
}
