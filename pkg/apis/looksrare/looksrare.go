package looksrare

import (
	"errors"
	"net/http"
	"time"

	"nftsiren/pkg/httpclient"
	"nftsiren/pkg/nft"
	"nftsiren/pkg/number"
)

var errSomethingWentWrong = errors.New("something went wrong")

var client = httpclient.NewClientWithLimit("https://api.looksrare.org/api/v1", 120, time.Minute)

func SetApiKey(apiKey string) {
	client.SetDefaultHeader("X-Looks-Api-Key", apiKey)
}

type genericResponse[T any] struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    T             `json:"data"`
	Errors  []interface{} `json:"errors"`
}

func get[T any](path []string, params map[string]string) (T, error) {
	var resp genericResponse[T]
	status, err := client.GetJson(path, params, &resp)
	if err != nil {
		if status >= 400 {
			return *new(T), errors.New(http.StatusText(status))
		}
		return *new(T), err
	}
	if !resp.Success {
		// TODO: find the actual error message
		if resp.Message != "" {
			return *new(T), errors.New(resp.Message)
		}
		return *new(T), errSomethingWentWrong
	}
	return resp.Data, nil
}

type collection struct {
	Address       string `json:"address"`
	Owner         string `json:"owner"`
	Setter        string `json:"setter"`
	Admin         string `json:"admin"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Symbol        string `json:"symbol"`
	Type          string `json:"type"`
	WebsiteLink   string `json:"websiteLink"`
	FacebookLink  string `json:"facebookLink"`
	TwitterLink   string `json:"twitterLink"`
	InstagramLink string `json:"instagramLink"`
	TelegramLink  string `json:"telegramLink"`
	MediumLink    string `json:"mediumLink"`
	DiscordLink   string `json:"discordLink"`
	IsVerified    bool   `json:"isVerified"`
	IsExplicit    bool   `json:"isExplicit"`
	LogoURI       string `json:"logoURI"`
	BannerURI     string `json:"bannerURI"`
}

func FetchCollection(address string) (nft.Collection, error) {
	c, err := get[collection]([]string{"collections"}, map[string]string{"address": address})
	if err != nil {
		return nft.Collection{}, err
	}
	return nft.Collection{
		Time:        time.Now(),
		Currency:    nft.ETH,
		Marketplace: nft.Looksrare,
		Symbol:      address,
		Address:     c.Address,
		Name:        c.Name,
		Description: c.Description,
		ImageURL:    c.LogoURI,
		Marketpage:  nft.Looksrare.MakeCollectionURL(address),
		Website:     c.WebsiteLink,
		Twitter:     c.TwitterLink,
		Discord:     c.DiscordLink,
		Stats:       nil,
	}, nil
}

type collectionStats struct {
	Address        string        `json:"address"`
	CountOwners    number.Number `json:"countOwners"`
	TotalSupply    number.Number `json:"totalSupply"`
	FloorPrice     number.Number `json:"floorPrice"`
	FloorChange24H number.Number `json:"floorChange24h"`
	FloorChange7D  number.Number `json:"floorChange7d"`
	FloorChange30D number.Number `json:"floorChange30d"`
	MarketCap      number.Number `json:"marketCap"`
	Volume24H      number.Number `json:"volume24h"`
	Average24H     number.Number `json:"average24h"`
	Count24H       number.Number `json:"count24h"`
	Change24H      number.Number `json:"change24h"`
	Volume7D       number.Number `json:"volume7d"`
	Average7D      number.Number `json:"average7d"`
	Count7D        number.Number `json:"count7d"`
	Change7D       number.Number `json:"change7d"`
	Volume1M       number.Number `json:"volume1m"`
	Average1M      number.Number `json:"average1m"`
	Count1M        number.Number `json:"count1m"`
	Change1M       number.Number `json:"change1m"`
	Volume3M       number.Number `json:"volume3m"`
	Average3M      number.Number `json:"average3m"`
	Count3M        number.Number `json:"count3m"`
	Change3M       number.Number `json:"change3m"`
	Volume6M       number.Number `json:"volume6m"`
	Average6M      number.Number `json:"average6m"`
	Count6M        number.Number `json:"count6m"`
	Change6M       number.Number `json:"change6m"`
	Volume1Y       number.Number `json:"volume1y"`
	Average1Y      number.Number `json:"average1y"`
	Count1Y        number.Number `json:"count1y"`
	Change1Y       number.Number `json:"change1y"`
	VolumeAll      number.Number `json:"volumeAll"`
	AverageAll     number.Number `json:"averageAll"`
	CountAll       number.Number `json:"countAll"`
}

func FetchCollectionStats(address string) (nft.CollectionStats, error) {
	stats, err := get[collectionStats]([]string{"collections", "stats"}, map[string]string{"address": address})
	if err != nil {
		return nft.CollectionStats{}, err
	}
	// Looksrare returns eth in wei format, make sure they converted correctly
	return nft.CollectionStats{
		Time:        time.Now(),
		Floor:       nft.WeiToEth(stats.FloorPrice),
		DaySales:    stats.Count24H,
		DayVolume:   nft.WeiToEth(stats.Volume24H),
		TotalSales:  stats.CountAll,
		TotalVolume: nft.WeiToEth(stats.VolumeAll),
		NumOwners:   stats.CountOwners,
		TotalSupply: stats.TotalSupply,
	}, nil
}
