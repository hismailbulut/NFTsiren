package magiceden

import (
	"errors"
	"net/http"
	"time"

	"nftsiren/pkg/httpclient"
	"nftsiren/pkg/nft"
	"nftsiren/pkg/number"
)

// This public API is free to use and the default limit is 120 QPM or 2 QPS
var client = httpclient.NewClientWithLimit("https://api-mainnet.magiceden.dev/v2", 120, time.Minute)

func SetApiKey(apiKey string) {
	client.SetBearerAuth(apiKey)
}

type errorFields struct {
	StatusCode *int    `json:"statusCode"`
	Error      *string `json:"error"`
	Message    *string `json:"message"`
}

func (resp errorFields) check() error {
	if resp.StatusCode != nil && *resp.StatusCode >= 400 {
		if resp.Message != nil {
			return errors.New(*resp.Message)
		}
		if resp.Error != nil {
			return errors.New(*resp.Error)
		}
		return errors.New(http.StatusText(*resp.StatusCode))
	}
	return nil
}

type hasErrorCheck interface {
	check() error
}

func get[T hasErrorCheck](path []string) (T, error) {
	var resp T
	status, err := client.GetJson(path, nil, &resp)
	if err != nil {
		// Error is either json or network error
		// Return status code if it
		if status >= 400 {
			return *new(T), errors.New(http.StatusText(status))
		}
		return *new(T), err
	}
	if err := resp.check(); err != nil {
		return *new(T), err
	}
	return resp, nil
}

type collection struct {
	errorFields
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	Twitter     string   `json:"twitter"`
	Discord     string   `json:"discord"`
	Website     string   `json:"website"`
	IsFlagged   bool     `json:"isFlagged"`
	FlagMessage string   `json:"flagMessage"`
	Categories  []string `json:"categories"`
	IsBadged    bool     `json:"isBadged"`
	collectionStats
}

func FetchCollection(symbol string) (nft.Collection, error) {
	resp, err := get[collection]([]string{"collections", symbol})
	if err != nil {
		return nft.Collection{}, err
	}
	stats := convertStats(resp.collectionStats)
	return nft.Collection{
		Time:        time.Now(),
		Currency:    nft.SOL,
		Marketplace: nft.Magiceden,
		Symbol:      symbol,
		Address:     "", // We don't know
		Name:        resp.Name,
		Description: resp.Description,
		ImageURL:    resp.Image,
		Marketpage:  nft.Magiceden.MakeCollectionURL(symbol),
		Website:     resp.Website,
		Twitter:     resp.Twitter,
		Discord:     resp.Discord,
		Stats:       &stats,
	}, nil
}

type collectionStats struct {
	errorFields
	Symbol       string        `json:"symbol"`
	FloorPrice   number.Number `json:"floorPrice"`
	ListedCount  number.Number `json:"listedCount"`
	AvgPrice24Hr number.Number `json:"avgPrice24hr"`
	VolumeAll    number.Number `json:"volumeAll"`
}

func FetchCollectionStats(symbol string) (nft.CollectionStats, error) {
	resp, err := get[collectionStats]([]string{"collections", symbol, "stats"})
	if err != nil {
		return nft.CollectionStats{}, err
	}
	return convertStats(resp), nil
}

func convertStats(stats collectionStats) nft.CollectionStats {
	return nft.CollectionStats{
		Time:        time.Now(),
		Floor:       nft.LamportsToSol(stats.FloorPrice),
		TotalVolume: nft.LamportsToSol(stats.VolumeAll),
		Listed:      stats.ListedCount,
	}
}
