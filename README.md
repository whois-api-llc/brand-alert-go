[![brand-alert-go license](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![brand-alert-go made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://pkg.go.dev/github.com/whois-api-llc/brand-alert-go)
[![brand-alert-go test](https://github.com/whois-api-llc/brand-alert-go/workflows/Test/badge.svg)](https://github.com/whois-api-llc/brand-alert-go/actions/)

# Overview

The client library for
[Brand Alert API](https://brand-alert.whoisxmlapi.com/)
in Go language.

The minimum go version is 1.17.

# Installation

The library is distributed as a Go module

```bash
go get github.com/whois-api-llc/brand-alert-go
```

# Examples

Full API documentation available [here](https://brand-alert.whoisxmlapi.com/api/documentation/making-requests)

You can find all examples in `example` directory.

## Create a new client

To start making requests you need the API Key. 
You can find it on your profile page on [whoisxmlapi.com](https://whoisxmlapi.com/).
Using the API Key you can create Client.

Most users will be fine with `NewBasicClient` function. 
```go
client := brandalert.NewBasicClient(apiKey)
```

If you want to set custom `http.Client` to use proxy then you can use `NewClient` function.
```go
transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

client := brandalert.NewClient(apiKey, brandalert.ClientParams{
    HTTPClient: &http.Client{
        Transport: transport,
        Timeout:   20 * time.Second,
    },
})
```

## Make basic requests

Brand Alert API searches across all recently registered & deleted domain names and returns result sets consisting of domain names that contain term(s) that are specified by you.

```go

// Make request to get a list of all domains matching the criteria.
brandAlertResp, resp, err := client.Purchase(ctx,
    &brandalert.SearchTerms{"google"}
    nil)

for _, obj := range brandAlertResp.DomainsList {
    log.Println(obj.DomainName)
}


// Make request to get only domains count.
domainsCount, _, err := client.Preview(ctx,
    &brandalert.SearchTerms{"google"},
    nil)

log.Println(domainsCount)

// Make request to get raw data in XML.
resp, err := client.RawData(ctx,
    &brandalert.SearchTerms{"google", "blog"},
    &brandalert.SearchTerms{"analytics"},
    brandalert.OptionResponseFormat("XML"))

log.Println(string(resp.Body))

```
