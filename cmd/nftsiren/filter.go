package main

import (
	"nftsiren/cmd/nftsiren/config"
	"nftsiren/pkg/nft"
	"sort"
	"strings"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type CollectionListSortingStrategy string

const (
	SortByFloorAscending  CollectionListSortingStrategy = "Sort by floor price (low to high)"
	SortByFloorDescending CollectionListSortingStrategy = "Sort by floor price (high to low)"
	SortByNameAZ          CollectionListSortingStrategy = "Sort by collection name (A-Z)"
	SortByNameZA          CollectionListSortingStrategy = "Sort by collection name (Z-A)"
	SortByMarket          CollectionListSortingStrategy = "Sort by marketplace (group them)"
)

// TODO: we need to save this to user config
// TODO: add currency filters (ETH, SOL, MATIC etc.)
type CollectionFilterPage struct {
	List                  widget.List
	FilterEnableOpensea   widget.Bool
	FilterEnableLooksrare widget.Bool
	FilterEnableMagiceden widget.Bool
	SortGroup             widget.Enum
}

func NewCollectionFilterPage() *CollectionFilterPage {
	page := &CollectionFilterPage{}
	page.List.Axis = layout.Vertical
	page.FilterEnableOpensea.Value = true
	page.FilterEnableLooksrare.Value = true
	page.FilterEnableMagiceden.Value = true
	return page
}

func (page *CollectionFilterPage) Title() string {
	return "Filter Collections"
}

func (page *CollectionFilterPage) Entering() {
	// NOTE: We also need to load this in daemon
	sort, err := config.Load[CollectionListSortingStrategy]("collectionFilter")
	if err == nil {
		page.SortGroup.Value = string(sort)
	}
}

func (page *CollectionFilterPage) Leaving() {
	config.Store("collectionFilter", page.SortGroup.Value)
}

func (page *CollectionFilterPage) Layout(gtx layout.Context, theme *Theme, pages *PageStack) layout.Dimensions {
	return theme.LayoutListSpaced(gtx, &page.List, theme.SmallVSpacer,
		material.Subtitle1(theme.Material(), "Filter").Layout,
		material.CheckBox(theme.Material(), &page.FilterEnableOpensea, "Opensea collections").Layout,
		material.CheckBox(theme.Material(), &page.FilterEnableLooksrare, "Looksrare collections").Layout,
		material.CheckBox(theme.Material(), &page.FilterEnableMagiceden, "Magiceden collections").Layout,
		material.Subtitle1(theme.Material(), "Sort").Layout,
		material.RadioButton(theme.Material(), &page.SortGroup, string(SortByFloorAscending), string(SortByFloorAscending)).Layout,
		material.RadioButton(theme.Material(), &page.SortGroup, string(SortByFloorDescending), string(SortByFloorDescending)).Layout,
		material.RadioButton(theme.Material(), &page.SortGroup, string(SortByNameAZ), string(SortByNameAZ)).Layout,
		material.RadioButton(theme.Material(), &page.SortGroup, string(SortByNameZA), string(SortByNameZA)).Layout,
		material.RadioButton(theme.Material(), &page.SortGroup, string(SortByMarket), string(SortByMarket)).Layout,
	)
}

func (page *CollectionFilterPage) Filter(c *Collection) bool {
	market := c.Market.Load()
	switch market {
	case nft.Opensea:
		return page.FilterEnableOpensea.Value
	case nft.Looksrare:
		return page.FilterEnableLooksrare.Value
	case nft.Magiceden:
		return page.FilterEnableMagiceden.Value
	}
	return true
}

func (page *CollectionFilterPage) Sort(collections []*Collection) {
	if page.SortGroup.Value == "" {
		return
	}
	less := func(i, j int) bool {
		return page.Less(collections[i], collections[j])
	}
	if sort.SliceIsSorted(collections, less) {
		return
	}
	sort.Slice(collections, less)
}

func (page *CollectionFilterPage) Less(i, j *Collection) bool {
	v := CollectionListSortingStrategy(page.SortGroup.Value)
	switch v {
	case SortByFloorAscending, SortByFloorDescending:
		i_floor, i_ok := i.Floor()
		j_floor, j_ok := j.Floor()
		switch {
		case !i_ok:
			return false
		case !j_ok:
			return true
		case v == SortByFloorAscending:
			return i_floor.LessThan(j_floor)
		case v == SortByFloorDescending:
			return i_floor.GreaterThan(j_floor)
		}
	case SortByNameAZ, SortByNameZA:
		i_name, i_ok := i.Name()
		j_name, j_ok := j.Name()
		switch {
		case !i_ok:
			return false
		case !j_ok:
			return true
		case v == SortByNameAZ:
			return strings.ToLower(i_name) < strings.ToLower(j_name)
		case v == SortByNameZA:
			return strings.ToLower(i_name) > strings.ToLower(j_name)
		}
	case SortByMarket:
		i_market := i.Market.Load()
		j_market := j.Market.Load()
		return i_market < j_market
	}
	return true
}
