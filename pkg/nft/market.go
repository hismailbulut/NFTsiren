package nft

import (
	"errors"
	"net/url"
	"strings"
)

const UNKNOWN_MARKETPLACE = "unknown marketplace"

type Marketplace uint32

const (
	Opensea Marketplace = iota
	Looksrare
	Magiceden
)

var MarketplaceNames = map[Marketplace]string{
	Opensea:   "Opensea",
	Looksrare: "Looksrare",
	Magiceden: "Magiceden",
}

func (market Marketplace) String() string {
	s, ok := MarketplaceNames[market]
	if ok {
		return s
	}
	return UNKNOWN_MARKETPLACE
}

func (market Marketplace) Host() string {
	switch market {
	case Opensea:
		return "opensea.io"
	case Looksrare:
		return "looksrare.org"
	case Magiceden:
		return "magiceden.io"
	}
	return UNKNOWN_MARKETPLACE
}

func (market Marketplace) CollectionsPath() string {
	switch market {
	case Opensea:
		return "collection"
	case Looksrare:
		return "collections"
	case Magiceden:
		return "marketplace"
	}
	return UNKNOWN_MARKETPLACE
}

func (market Marketplace) MakeCollectionURL(symbol string) string {
	return "https://" + market.Host() + "/" + market.CollectionsPath() + "/" + symbol
}

// Parses given raw url for the given collection and returns the slug
// used in the fetch apis, or error if there is
func (market Marketplace) ParseCollectionURL(rawurl string) (string, error) {
	invalidErr := errors.New("invalid url")
	parsed, err := url.Parse(rawurl)
	if err != nil {
		return "", invalidErr
	}
	if parsed.Hostname() != market.Host() {
		return "", invalidErr
	}
	elements := strings.Split(strings.Trim(parsed.EscapedPath(), "/"), "/")
	for i, e := range elements {
		if e == market.CollectionsPath() && len(elements) > i+1 {
			return elements[i+1], nil
		}
	}
	return "", invalidErr
}

func (market Marketplace) MarshalText() ([]byte, error) {
	s, ok := MarketplaceNames[market]
	if ok {
		return []byte(s), nil
	}
	return nil, errors.New(UNKNOWN_MARKETPLACE)
}

func (market *Marketplace) UnmarshalText(text []byte) error {
	for k, v := range MarketplaceNames {
		if v == string(text) {
			*market = k
			return nil
		}
	}
	return errors.New(UNKNOWN_MARKETPLACE)
}
