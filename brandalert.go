package brandalert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// BrandAlert is an interface for Brand Alert API.
type BrandAlert interface {
	// Purchase returns parsed Brand Alert API response.
	Purchase(ctx context.Context, includeSearchTerms *SearchTerms, excludeSearchTerms *SearchTerms, option ...Option) (*BrandAlertResponse, *Response, error)

	// Preview returns only the number of domains. No credits deducted.
	Preview(ctx context.Context, includeSearchTerms *SearchTerms, excludeSearchTerms *SearchTerms, option ...Option) (int, *Response, error)

	// RawData returns raw Brand Alert API response as the Response struct with Body saved as a byte slice.
	RawData(ctx context.Context, includeSearchTerms *SearchTerms, excludeSearchTerms *SearchTerms, option ...Option) (*Response, error)
}

// Response is the http.Response wrapper with Body saved as a byte slice.
type Response struct {
	*http.Response

	// Body is the byte slice representation of http.Response Body.
	Body []byte
}

// brandAlertServiceOp is the type implementing the BrandAlert interface.
type brandAlertServiceOp struct {
	client  *Client
	baseURL *url.URL
}

var _ BrandAlert = &brandAlertServiceOp{}

// newRequest creates the API request with default parameters and specified body.
func (service brandAlertServiceOp) newRequest(body []byte) (*http.Request, error) {
	req, err := service.client.NewRequest(http.MethodPost, service.baseURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	return req, nil
}

// apiResponse is used for parsing Brand Alert API response as a model instance.
type apiResponse struct {
	BrandAlertResponse
	ErrorMessage
}

// validateSearchTerms validates the terms of search.
func validateSearchTerms(includeSearchTerms *SearchTerms, excludeSearchTerms *SearchTerms) error {
	const limitOfSearchTerms = 4

	if includeSearchTerms == nil || len(*includeSearchTerms) == 0 || len(*includeSearchTerms) > limitOfSearchTerms {
		return &ArgError{"includeSearchTerms", "must have between 1 and 4 items."}
	}

	if excludeSearchTerms != nil && len(*excludeSearchTerms) > limitOfSearchTerms {
		return &ArgError{"excludeSearchTerms", "must have between 0 and 4 items."}
	}

	return nil
}

// validateOptions validates options.
func validateOptions(opts ...Option) error {
	for _, opt := range opts {
		if opt == nil {
			return &ArgError{"Option", "can not be nil"}
		}
	}
	return nil
}

// request returns intermediate API response for further actions.
func (service brandAlertServiceOp) request(
	ctx context.Context,
	includeSearchTerms *SearchTerms, excludeSearchTerms *SearchTerms,
	purchase bool,
	opts ...Option) (*Response, error) {
	err := validateSearchTerms(includeSearchTerms, excludeSearchTerms)
	if err != nil {
		return nil, err
	}

	var request = &brandAlertRequest{
		service.client.apiKey,
		includeSearchTerms,
		excludeSearchTerms,
		"",
		"preview",
		false,
		true,
		"json",
	}

	if purchase {
		request.Mode = "purchase"
	}

	if err := validateOptions(opts...); err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(request)
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := service.newRequest(requestBody)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer

	resp, err := service.client.Do(ctx, req, &b)
	if err != nil {
		return &Response{
			Response: resp,
			Body:     b.Bytes(),
		}, err
	}

	return &Response{
		Response: resp,
		Body:     b.Bytes(),
	}, nil
}

// parse parses raw Brand Alert API response.
func parse(raw []byte) (*apiResponse, error) {
	var response apiResponse

	err := json.NewDecoder(bytes.NewReader(raw)).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("cannot parse response: %w", err)
	}

	return &response, nil
}

// Purchase returns parsed Brand Alert API response.
func (service brandAlertServiceOp) Purchase(
	ctx context.Context,
	includeSearchTerms *SearchTerms, excludeSearchTerms *SearchTerms,
	opts ...Option,
) (brandAlertResponse *BrandAlertResponse, resp *Response, err error) {
	optsJSON := make([]Option, 0, len(opts)+1)
	optsJSON = append(optsJSON, opts...)
	optsJSON = append(optsJSON, OptionResponseFormat("json"))

	resp, err = service.request(ctx, includeSearchTerms, excludeSearchTerms, true, optsJSON...)
	if err != nil {
		return nil, resp, err
	}

	brandAlertResp, err := parse(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	if brandAlertResp.Message != nil || brandAlertResp.Code != 0 {
		return nil, nil, &ErrorMessage{
			Code:    brandAlertResp.Code,
			Message: brandAlertResp.Message,
		}
	}

	return &brandAlertResp.BrandAlertResponse, resp, nil
}

// Preview returns only the number of domains. No credits deducted.
func (service brandAlertServiceOp) Preview(
	ctx context.Context,
	includeSearchTerms *SearchTerms, excludeSearchTerms *SearchTerms,
	opts ...Option,
) (domainsCount int, resp *Response, err error) {
	optsJSON := make([]Option, 0, len(opts)+1)
	optsJSON = append(optsJSON, opts...)
	optsJSON = append(optsJSON, OptionResponseFormat("json"))

	resp, err = service.request(ctx, includeSearchTerms, excludeSearchTerms, false, optsJSON...)
	if err != nil {
		return 0, resp, err
	}

	brandAlertResp, err := parse(resp.Body)
	if err != nil {
		return 0, resp, err
	}

	if brandAlertResp.Message != nil || brandAlertResp.Code != 0 {
		return 0, nil, &ErrorMessage{
			Code:    brandAlertResp.Code,
			Message: brandAlertResp.Message,
		}
	}

	return brandAlertResp.DomainsCount, resp, nil
}

// RawData returns raw Brand Alert API response as the Response struct with Body saved as a byte slice.
func (service brandAlertServiceOp) RawData(
	ctx context.Context,
	includeSearchTerms *SearchTerms, excludeSearchTerms *SearchTerms,
	opts ...Option,
) (resp *Response, err error) {
	resp, err = service.request(ctx, includeSearchTerms, excludeSearchTerms, true, opts...)
	if err != nil {
		return resp, err
	}

	if respErr := checkResponse(resp.Response); respErr != nil {
		return resp, respErr
	}

	return resp, nil
}

// ArgError is the argument error.
type ArgError struct {
	Name    string
	Message string
}

// Error returns error message as a string.
func (a *ArgError) Error() string {
	return `invalid argument: "` + a.Name + `" ` + a.Message
}
