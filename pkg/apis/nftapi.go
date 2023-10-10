package apis

import (
	"errors"

	"nftsiren/pkg/apis/looksrare"
	"nftsiren/pkg/apis/magiceden"
	"nftsiren/pkg/apis/opensea"
	"nftsiren/pkg/nft"
)

type ApiProvider int32

const (
	// TODO: we will use these instead of nft.Marketplace
	// to allow different apis returning information about various marketplaces
	// To use them we will ask to user which api should be used
	// instead of the marketplace, api results already contains
	// marketplace information either in the collection info or stats
	// We only allow api's has a public api (either with api key or keyless)
	// and has docs/sdk about this api
	OpenseaApi ApiProvider = iota
	LooksrareApi
	MagicedenApi
)

func FetchCollection(marketplace nft.Marketplace, symbol string) (nft.Collection, error) {
	switch marketplace {
	case nft.Opensea:
		return opensea.FetchCollection(symbol)
	case nft.Looksrare:
		return looksrare.FetchCollection(symbol)
	case nft.Magiceden:
		return magiceden.FetchCollection(symbol)
	}
	return nft.Collection{}, errors.New(nft.UNKNOWN_MARKETPLACE)
}

func FetchCollectionStats(marketplace nft.Marketplace, symbol string) (nft.CollectionStats, error) {
	switch marketplace {
	case nft.Opensea:
		return opensea.FetchCollectionStats(symbol)
	case nft.Looksrare:
		return looksrare.FetchCollectionStats(symbol)
	case nft.Magiceden:
		return magiceden.FetchCollectionStats(symbol)
	}
	return nft.CollectionStats{}, errors.New(nft.UNKNOWN_MARKETPLACE)
}
