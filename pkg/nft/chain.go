package nft

import "errors"

const UNKNOWN_CHAIN = "unknown chain"

type Chain uint32

const (
	ETH Chain = iota
	SOL
)

func (c Chain) String() string {
	switch c {
	case ETH:
		return "ETH"
	case SOL:
		return "SOL"
	}
	return UNKNOWN_CHAIN
}

func (chain Chain) MarshalText() ([]byte, error) {
	return []byte(chain.String()), nil
}

func (chain *Chain) UnmarshalText(text []byte) error {
	switch string(text) {
	case "ETH":
		*chain = ETH
	case "SOL":
		*chain = SOL
	default:
		return errors.New(UNKNOWN_CHAIN)
	}
	return nil
}
