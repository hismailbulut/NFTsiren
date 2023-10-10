package opensea

import (
	"errors"
	"net/http"
	"time"

	"nftsiren/pkg/httpclient"
	"nftsiren/pkg/nft"
)

var errSomethingWentWrong = errors.New("something went wrong")

// GET requests are limited to 4/sec per API key. POST requests are limited to 2/sec per API key.
var client = httpclient.NewClientWithLimit("https://api.opensea.io/api/v1", 4, time.Second)

func SetApiKey(apiKey string) {
	client.SetDefaultHeader("X-API-KEY", apiKey)
}

func get(path []string, resp hasErrorCheck) error {
	status, err := client.GetJson(path, nil, resp)
	if err != nil {
		if status >= 400 {
			return errors.New(http.StatusText(status))
		}
		return err
	}
	if err := resp.check(); err != nil {
		return err
	}
	return nil
}

func FetchCollection(slug string) (nft.Collection, error) {
	var resp struct {
		Collection *collection `json:"collection"`
		errorFields
	}
	err := get([]string{"collection", slug}, &resp)
	if err != nil {
		return nft.Collection{}, err
	}
	c := resp.Collection
	if c == nil {
		return nft.Collection{}, errors.New("collection information not available")
	}
	ret := nft.Collection{
		Time:        time.Now(),
		Currency:    nft.ETH,
		Marketplace: nft.Opensea,
		Symbol:      slug,
		Address:     "", // TODO: find
		Name:        c.Name,
		Description: c.Description,
		ImageURL:    c.ImageURL,
		Marketpage:  nft.Opensea.MakeCollectionURL(slug),
		Website:     c.ExternalURL,
		Twitter:     "https://twitter.com/" + c.TwitterUsername,
		Discord:     c.DiscordURL,
	}
	if c.Stats != nil {
		stats := c.Stats.convert()
		ret.Stats = &stats
	}
	return ret, nil
}

func FetchCollectionStats(slug string) (nft.CollectionStats, error) {
	var resp struct {
		Stats *collectionStats `json:"stats"`
		errorFields
	}
	err := get([]string{"collection", slug, "stats"}, &resp)
	if err != nil {
		return nft.CollectionStats{}, err
	}
	if resp.Stats == nil {
		return nft.CollectionStats{}, errors.New("collections statistics not available")
	}
	return resp.Stats.convert(), nil
}
