package example

import (
	"context"
	"errors"
	brandalert "github.com/whois-api-llc/brand-alert-go"
	"log"
	"net/http"
	"time"
)

func BrandAlertPreview(apikey string) {
	client := brandalert.NewBasicClient(apikey)

	// Get a number of domains matching the criteria.
	domainsCount, _, err := client.Preview(context.Background(),
		// specify the including search terms
		&brandalert.SearchTerms{"google"},
		// excluding search terms can be unspecified
		nil)

	if err != nil {
		// Handle error message returned by server
		var apiErr *brandalert.ErrorMessage
		if errors.As(err, &apiErr) {
			log.Println(apiErr.Code)
			log.Println(apiErr.Message)
		}
		log.Println(err)
		return
	}

	log.Println(domainsCount)
}

func BrandAlertPurchase(apikey string) {
	client := brandalert.NewClient(apikey, brandalert.ClientParams{
		HTTPClient: &http.Client{
			Transport: nil,
			Timeout:   40 * time.Second,
		},
	})

	// Get parsed Brand Alert API response as a model instance.
	brandAlertResp, resp, err := client.Purchase(context.Background(),
		&brandalert.SearchTerms{"whois"},
		nil,
		// this option is ignored, as the inner parser works with JSON only
		brandalert.OptionResponseFormat("XML"),
		// this option results in the search terms set will be enriched with their possible typos
		brandalert.OptionWithTypos(true),
		// this option results in domain names in the response will be encoded to Punycode
		brandalert.OptionPunycode(true))

	if err != nil {
		// Handle error message returned by server
		var apiErr *brandalert.ErrorMessage
		if errors.As(err, &apiErr) {
			log.Println(apiErr.Code)
			log.Println(apiErr.Message)
		}
		log.Println(err)
		return
	}

	// Then print all "added" domains.
	for _, obj := range brandAlertResp.DomainsList {
		if obj.Action == brandalert.Added {
			log.Println(obj.DomainName, obj.Action, time.Time(obj.Date).Format("2006-01-02"))
		}
	}

	log.Println("raw response is always in JSON format. Most likely you don't need it.")
	log.Printf("raw response: %s\n", string(resp.Body))
}

func BrandAlertRawData(apikey string) {
	client := brandalert.NewBasicClient(apikey)

	// Get raw API response.
	resp, err := client.RawData(context.Background(),
		// specify the including search terms
		&brandalert.SearchTerms{"google", "blog"},
		// specify the excluding search terms
		&brandalert.SearchTerms{"analytics"},
		// this option results in search through activities discovered since the given date
		brandalert.OptionSinceDate(time.Date(2022, 10, 14, 0, 0, 0, 0, time.UTC)))

	if err != nil {
		// Handle error message returned by server
		log.Println(err)
	}

	if resp != nil {
		log.Println(string(resp.Body))
	}
}
