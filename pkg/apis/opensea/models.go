package opensea

import (
	"errors"
	"strings"
	"time"

	"nftsiren/pkg/nft"
	"nftsiren/pkg/number"
)

type errorFields struct {
	Success *bool    `json:"success"`
	Errors  []string `json:"errors"`
}

type hasErrorCheck interface {
	check() error
}

func (resp errorFields) check() error {
	if resp.Success != nil && !*resp.Success {
		if len(resp.Errors) > 0 {
			return errors.New(strings.Join(resp.Errors, ", "))
		}
		return errSomethingWentWrong
	}
	return nil
}

type collection struct {
	Editors                     []string         `json:"editors"`
	PaymentTokens               []any            `json:"payment_tokens"`
	PrimaryAssetContracts       []any            `json:"primary_asset_contracts"`
	Traits                      any              `json:"traits"`
	Stats                       *collectionStats `json:"stats"`
	BannerImageURL              string           `json:"banner_image_url"`
	ChatURL                     string           `json:"chat_url"`
	CreatedDate                 string           `json:"created_date"`
	DefaultToFiat               bool             `json:"default_to_fiat"`
	Description                 string           `json:"description"`
	DevBuyerFeeBasisPoints      number.Number    `json:"dev_buyer_fee_basis_points"`
	DevSellerFeeBasisPoints     number.Number    `json:"dev_seller_fee_basis_points"`
	DiscordURL                  string           `json:"discord_url"`
	DisplayData                 any              `json:"display_data"`
	ExternalURL                 string           `json:"external_url"`
	Featured                    bool             `json:"featured"`
	FeaturedImageURL            string           `json:"featured_image_url"`
	Hidden                      bool             `json:"hidden"`
	SafelistRequestStatus       string           `json:"safelist_request_status"`
	ImageURL                    string           `json:"image_url"`
	IsSubjectToWhitelist        bool             `json:"is_subject_to_whitelist"`
	LargeImageURL               string           `json:"large_image_url"`
	MediumUsername              string           `json:"medium_username"`
	Name                        string           `json:"name"`
	OnlyProxiedTransfers        bool             `json:"only_proxied_transfers"`
	OpenseaBuyerFeeBasisPoints  number.Number    `json:"opensea_buyer_fee_basis_points"`
	OpenseaSellerFeeBasisPoints number.Number    `json:"opensea_seller_fee_basis_points"`
	PayoutAddress               string           `json:"payout_address"`
	RequireEmail                bool             `json:"require_email"`
	ShortDescription            string           `json:"short_description"`
	Slug                        string           `json:"slug"`
	TelegramURL                 string           `json:"telegram_url"`
	TwitterUsername             string           `json:"twitter_username"`
	InstagramUsername           string           `json:"instagram_username"`
	WikiURL                     string           `json:"wiki_url"`
	IsNsfw                      bool             `json:"is_nsfw"`
}

type collectionStats struct {
	OneDayVolume          number.Number `json:"one_day_volume"`
	OneDayChange          number.Number `json:"one_day_change"`
	OneDaySales           number.Number `json:"one_day_sales"`
	OneDayAveragePrice    number.Number `json:"one_day_average_price"`
	SevenDayVolume        number.Number `json:"seven_day_volume"`
	SevenDayChange        number.Number `json:"seven_day_change"`
	SevenDaySales         number.Number `json:"seven_day_sales"`
	SevenDayAveragePrice  number.Number `json:"seven_day_average_price"`
	ThirtyDayVolume       number.Number `json:"thirty_day_volume"`
	ThirtyDayChange       number.Number `json:"thirty_day_change"`
	ThirtyDaySales        number.Number `json:"thirty_day_sales"`
	ThirtyDayAveragePrice number.Number `json:"thirty_day_average_price"`
	TotalVolume           number.Number `json:"total_volume"`
	TotalSales            number.Number `json:"total_sales"`
	TotalSupply           number.Number `json:"total_supply"`
	Count                 number.Number `json:"count"`
	NumOwners             number.Number `json:"num_owners"`
	AveragePrice          number.Number `json:"average_price"`
	NumReports            number.Number `json:"num_reports"`
	MarketCap             number.Number `json:"market_cap"`
	FloorPrice            number.Number `json:"floor_price"`
}

func (stats collectionStats) convert() nft.CollectionStats {
	return nft.CollectionStats{
		Time:        time.Now(),
		Floor:       stats.FloorPrice,
		DaySales:    stats.OneDaySales,
		DayVolume:   stats.OneDayVolume,
		TotalSales:  stats.TotalSales,
		TotalVolume: stats.TotalVolume,
		NumOwners:   stats.NumOwners,
		TotalSupply: stats.Count,
	}
}
