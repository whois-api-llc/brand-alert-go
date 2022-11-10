package brandalert

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

const (
	pathBrandAlertResponseOK         = "/BrandAlert/ok"
	pathBrandAlertResponseError      = "/BrandAlert/error"
	pathBrandAlertResponse500        = "/BrandAlert/500"
	pathBrandAlertResponsePartial1   = "/BrandAlert/partial"
	pathBrandAlertResponsePartial2   = "/BrandAlert/partial2"
	pathBrandAlertResponseUnparsable = "/BrandAlert/unparsable"
)

const apiKey = "at_LoremIpsumDolorSitAmetConsect"

// dummyServer is the sample of the Brand Alert API server for testing.
func dummyServer(resp, respUnparsable string, respErr string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var response string

		response = resp

		switch req.URL.Path {
		case pathBrandAlertResponseOK:
		case pathBrandAlertResponseError:
			w.WriteHeader(499)
			response = respErr
		case pathBrandAlertResponse500:
			w.WriteHeader(500)
			response = respUnparsable
		case pathBrandAlertResponsePartial1:
			response = response[:len(response)-10]
		case pathBrandAlertResponsePartial2:
			w.Header().Set("Content-Length", strconv.Itoa(len(response)))
			response = response[:len(response)-10]
		case pathBrandAlertResponseUnparsable:
			response = respUnparsable
		default:
			panic(req.URL.Path)
		}
		_, err := w.Write([]byte(response))
		if err != nil {
			panic(err)
		}
	}))

	return server
}

// newAPI returns new Brand Alert API client for testing.
func newAPI(apiServer *httptest.Server, link string) *Client {
	apiURL, err := url.Parse(apiServer.URL)
	if err != nil {
		panic(err)
	}

	apiURL.Path = link

	params := ClientParams{
		HTTPClient:        apiServer.Client(),
		BrandAlertBaseURL: apiURL,
	}

	return NewClient(apiKey, params)
}

// TestBrandAlertPreview tests the Preview function.
func TestBrandAlertPreview(t *testing.T) {
	checkResultRec := func(res int) bool {
		return res != 0
	}

	ctx := context.Background()

	const resp = `{"domainsCount":4}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":["Test error message."]}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory1 *SearchTerms
		mandatory2 *SearchTerms
		option     Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request",
			path: pathBrandAlertResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{"whois"},
					nil,
					OptionResponseFormat("json"),
				},
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathBrandAlertResponse500,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{"whois"},
					nil,
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "partial response 1",
			path: pathBrandAlertResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{"whois"},
					nil,
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: unexpected EOF",
		},
		{
			name: "partial response 2",
			path: pathBrandAlertResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{"whois"},
					nil,
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathBrandAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{"whois"},
					nil,
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "API error: [499] [Test error message.]",
		},
		{
			name: "unparsable response",
			path: pathBrandAlertResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{"whois"},
					nil,
					OptionResponseFormat("xml"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "invalid argument1",
			path: pathBrandAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{},
					nil,
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "includeSearchTerms" must have between 1 and 4 items.`,
		},
		{
			name: "invalid argument2",
			path: pathBrandAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					nil,
					nil,
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "includeSearchTerms" must have between 1 and 4 items.`,
		},
		{
			name: "invalid argument3",
			path: pathBrandAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{"whois"},
					&SearchTerms{"1", "2", "3", "4", "5"},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "excludeSearchTerms" must have between 0 and 4 items.`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			gotRec, _, err := api.Preview(tt.args.ctx, tt.args.options.mandatory1, tt.args.options.mandatory2, tt.args.options.option)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("BrandAlert.Preview() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.want {
				if !checkResultRec(gotRec) {
					t.Errorf("BrandAlert.Preview() got = %v, expected something else", gotRec)
				}
			} else {
				if gotRec != 0 {
					t.Errorf("BrandAlert.Get() got = %v, expected nil", gotRec)
				}
			}
		})
	}
}

// TestBrandAlertPurchase tests the Purchase function.
func TestBrandAlertPurchase(t *testing.T) {
	checkResultRec := func(res *BrandAlertResponse) bool {
		return res != nil
	}

	ctx := context.Background()

	const resp = `{"domainsCount":4,"domainsList":[
{"domainName":"batchwhois.com","date":"2022-10-30","action":"discovered"},
{"domainName":"betterwhoislookup.com","date":"2022-10-30","action":"discovered"},
{"domainName":"whoisdomainlookup.info","date":"2022-10-30","action":"updated"},
{"domainName":"whoisdodster.com","date":"2022-10-30","action":"added"}]}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory1 *SearchTerms
		mandatory2 *SearchTerms
		option     Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request",
			path: pathBrandAlertResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{"whois"},
					nil,
					OptionResponseFormat("json"),
				},
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathBrandAlertResponse500,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{"whois"},
					nil,
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "partial response 1",
			path: pathBrandAlertResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{"whois"},
					nil,
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: unexpected EOF",
		},
		{
			name: "partial response 2",
			path: pathBrandAlertResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{"whois"},
					nil,
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathBrandAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{"whois"},
					nil,
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: "API error: [499] [Test error message.]",
		},
		{
			name: "unparsable response",
			path: pathBrandAlertResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{"whois"},
					nil,
					OptionResponseFormat("xml"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "invalid argument1",
			path: pathBrandAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{},
					nil,
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "includeSearchTerms" must have between 1 and 4 items.`,
		},
		{
			name: "invalid argument2",
			path: pathBrandAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					nil,
					nil,
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "includeSearchTerms" must have between 1 and 4 items.`,
		},
		{
			name: "invalid argument3",
			path: pathBrandAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{"whois"},
					&SearchTerms{"1", "2", "3", "4", "5"},
					OptionResponseFormat("json"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "excludeSearchTerms" must have between 0 and 4 items.`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			gotRec, _, err := api.Purchase(tt.args.ctx, tt.args.options.mandatory1, tt.args.options.mandatory2, tt.args.options.option)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("BrandAlert.Purchase() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.want {
				if !checkResultRec(gotRec) {
					t.Errorf("BrandAlert.Purchase() got = %v, expected something else", gotRec)
				}
			} else {
				if gotRec != nil {
					t.Errorf("BrandAlert.Purchase() got = %v, expected nil", gotRec)
				}
			}
		})
	}
}

// TestBrandAlertRawData tests the GetRaw function.
func TestBrandAlertRawData(t *testing.T) {
	checkResultRaw := func(res []byte) bool {
		return len(res) != 0
	}

	ctx := context.Background()

	const resp = `{"domainsCount":4,"domainsList":[
{"domainName":"batchwhois.com","date":"2022-10-30","action":"discovered"},
{"domainName":"betterwhoislookup.com","date":"2022-10-30","action":"discovered"},
{"domainName":"whoisdomainlookup.info","date":"2022-10-30","action":"updated"},
{"domainName":"whoisdodster.com","date":"2022-10-30","action":"added"}]}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory *SearchTerms
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		wantErr string
	}{
		{
			name: "successful request",
			path: pathBrandAlertResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{"whois"},
					OptionResponseFormat("json"),
				},
			},
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathBrandAlertResponse500,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{"whois"},
					OptionResponseFormat("json"),
				},
			},
			wantErr: "API failed with status code: 500",
		},
		{
			name: "partial response 1",
			path: pathBrandAlertResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{"whois"},
					OptionResponseFormat("json"),
				},
			},
			wantErr: "",
		},
		{
			name: "partial response 2",
			path: pathBrandAlertResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{"whois"},
					OptionResponseFormat("json"),
				},
			},
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "unparsable response",
			path: pathBrandAlertResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{"whois"},
					OptionResponseFormat("xml"),
				},
			},
			wantErr: "",
		},
		{
			name: "could not process request",
			path: pathBrandAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{"whois"},
					OptionResponseFormat("json"),
				},
			},
			wantErr: "API failed with status code: 499",
		},
		{
			name: "invalid argument1",
			path: pathBrandAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					&SearchTerms{},
					OptionResponseFormat("json"),
				},
			},
			wantErr: `invalid argument: "includeSearchTerms" must have between 1 and 4 items.`,
		},
		{
			name: "invalid argument2",
			path: pathBrandAlertResponseError,
			args: args{
				ctx: ctx,
				options: options{
					nil,
					OptionResponseFormat("json"),
				},
			},
			wantErr: `invalid argument: "includeSearchTerms" must have between 1 and 4 items.`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			resp, err := api.RawData(tt.args.ctx, tt.args.options.mandatory, nil)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("BrandAlert.GetRaw() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if resp != nil && !checkResultRaw(resp.Body) {
				t.Errorf("BrandAlert.GetRaw() got = %v, expected something else", string(resp.Body))
			}
		})
	}
}
