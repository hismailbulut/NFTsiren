package nft

import (
	"time"
)

type Collection struct {
	Time        time.Time        // Fetch time of this collection
	Currency    Chain            // Either ETH or SOL
	Marketplace Marketplace      //
	Symbol      string           // Unique identifier for marketplace
	Address     string           // Collection's primary contract address
	Name        string           //
	Description string           //
	ImageURL    string           //
	Marketpage  string           // URL of the collection webpage
	Website     string           //
	Twitter     string           //
	Discord     string           //
	Stats       *CollectionStats // Collection may have stats information but not guaranteed, fetch it additionally if it is nil
}
