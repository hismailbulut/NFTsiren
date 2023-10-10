package nft

import (
	"time"

	"nftsiren/pkg/number"
)

type StatDescription struct {
	Label string
	Value number.Number
}

type CollectionStats struct {
	Time        time.Time // Fetch time
	Floor       number.Number
	DaySales    number.Number
	DayVolume   number.Number
	TotalSales  number.Number
	TotalVolume number.Number
	NumOwners   number.Number
	TotalSupply number.Number
	Listed      number.Number
}

func (stats CollectionStats) IsValid() bool {
	return !stats.Time.IsZero() && !stats.Floor.IsNil()
}

func (stats CollectionStats) IsRecent(duration time.Duration) bool {
	return stats.Time.After(time.Now().Add(-duration))
}

func (stats *CollectionStats) All() []StatDescription {
	// Change this number depending on available statistics
	allArr := [8]StatDescription{}
	index := 0

	add := func(lbl string, val number.Number) {
		if !val.IsNil() {
			allArr[index].Label = lbl
			allArr[index].Value = val
			index++
		}
	}

	add("Floor", stats.Floor)
	add("Sales 24h", stats.DaySales)
	add("Volume 24h", stats.DayVolume)
	add("Sales", stats.TotalSales)
	add("Volume", stats.TotalVolume)
	add("Owners", stats.NumOwners)
	add("Supply", stats.TotalSupply)
	add("Listed", stats.Listed)

	return allArr[:index]
}
