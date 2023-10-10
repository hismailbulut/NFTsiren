package main

import (
	"errors"
	"fmt"
	"strings"

	"nftsiren/cmd/nftsiren/widgets"
	"nftsiren/pkg/bench"
	"nftsiren/pkg/nft"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

type CollectionCreationPage struct {
	Daemon *Daemon
	List   widget.List
	Market widgets.TypedEnum[nft.Marketplace]
	Url    component.TextField
	Error  error
	Ok     widget.Clickable
}

func NewCollectionCreationPage(daemon *Daemon) *CollectionCreationPage {
	page := &CollectionCreationPage{
		Daemon: daemon,
	}
	page.List.Axis = layout.Vertical
	page.Market.SetKeys([]nft.Marketplace{
		nft.Opensea,
		nft.Looksrare,
		nft.Magiceden,
	}...)
	page.Url.SingleLine = true
	page.Url.Submit = true
	return page
}

func (page *CollectionCreationPage) Title() string {
	return "Add New Collection"
}

func (page *CollectionCreationPage) Entering() {}

func (page *CollectionCreationPage) Leaving() {
	// reset page
	page.Market.State.Value = ""
	page.Url.SetText("")
	page.Error = nil
}

func (page *CollectionCreationPage) Layout(gtx layout.Context, theme *Theme, pages *PageStack) layout.Dimensions {
	if page.Ok.Clicked() {
		err := page.AddCollection()
		if err != nil {
			page.Error = err
		} else {
			pages.Pop()
		}
	}
	return theme.LayoutForm(gtx, &page.List, &page.Ok,
		// Market label
		func(gtx layout.Context) layout.Dimensions {
			return material.Body1(theme.Material(), "Market").Layout(gtx)
		},
		// Type select
		func(gtx layout.Context) layout.Dimensions {
			return page.Market.Layout(gtx, theme.Material())
		},
		// URL entry
		func(gtx layout.Context) layout.Dimensions {
			return page.Url.Layout(gtx, theme.Material(), "URL or slug")
			// return material.Editor(page.Theme, &newAlertPage.Value, hint).Layout(gtx)
		},
		// Error
		func(gtx layout.Context) layout.Dimensions {
			if page.Error == nil {
				return layout.Dimensions{}
			}
			label := material.Caption(theme.Material(), page.Error.Error())
			label.Color = theme.Error
			label.Alignment = text.Middle
			return label.Layout(gtx)
		},
	)
}

func (page *CollectionCreationPage) AddCollection() error {
	market, ok := page.Market.SelectedType()
	if !ok {
		return fmt.Errorf("select marketplace")
	}
	var symbol string
	urlstr := page.Url.Text()
	if urlstr == "" {
		return fmt.Errorf("enter collection url")
	} else if strings.Contains(urlstr, "/") {
		// This is an URL, parse it and get symbol
		var err error
		symbol, err = market.ParseCollectionURL(urlstr)
		if err != nil {
			return fmt.Errorf("collection url is not valid")
		}
	} else {
		// Maybe user directly entered collection slug or address and not URL
		symbol = urlstr
	}
	collection := NewCollection(page.Daemon, market, symbol)
	ok = page.Daemon.AddCollection(collection)
	if !ok {
		return errors.New("this collection is already in the list")
	}
	return nil
}

type CollectionsPage struct {
	Daemon *Daemon
	// State
	List widget.List
	// filter
	FilterButton widget.Clickable
	FilterPage   *CollectionFilterPage
	// add new collection
	AddCollectionButton    widget.Clickable
	CollectionCreationPage *CollectionCreationPage
}

func NewCollectionsPage(daemon *Daemon) *CollectionsPage {
	page := &CollectionsPage{
		Daemon: daemon,
	}
	page.List.Axis = layout.Vertical
	page.FilterPage = NewCollectionFilterPage()
	page.CollectionCreationPage = NewCollectionCreationPage(daemon)
	return page
}

func (page *CollectionsPage) Title() string {
	return "Collections"
}

func (page *CollectionsPage) Entering() {}

func (page *CollectionsPage) Leaving() {}

func (page *CollectionsPage) Layout(gtx layout.Context, theme *Theme, pages *PageStack) layout.Dimensions {
	defer bench.Begin()()
	return layout.Flex{
		Axis: layout.Vertical,
		// Spacing:   layout.SpaceEvenly,
		// Alignment: layout.Middle,
	}.Layout(gtx,
		// Top actions
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// return page.Theme.SmallInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Spacing:   layout.SpaceBetween,
				Alignment: layout.Middle,
			}.Layout(gtx,
				// Add button
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if page.AddCollectionButton.Clicked() {
						pages.Push(page.CollectionCreationPage)
					}
					return theme.Button("Add collection", &page.AddCollectionButton, PrimaryButton).Layout(gtx)
				}),
				// Total collection count
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					count := page.Daemon.CollectionCount()
					countstr := fmt.Sprintf("Total of %d collections", count)
					label := material.Caption(theme.Material(), countstr)
					label.Alignment = text.End
					return label.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: theme.SmallSpace}.Layout),
				// Filter button
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if page.FilterButton.Clicked() {
						pages.Push(page.FilterPage)
					}
					return theme.IconButton(theme.FilterIcon, &page.FilterButton).Layout(gtx)
				}),
			)
			// })
		}),
		layout.Rigid(theme.MediumVSpacer.Layout),
		// Collection list
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// return page.Theme.SmallInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return page.layoutList(gtx, theme, pages)
			// })
		}),
	)
}

func (page *CollectionsPage) layoutList(gtx layout.Context, theme *Theme, pages *PageStack) layout.Dimensions {
	defer bench.Begin()()
	collectionCount := page.Daemon.CollectionCount()
	if collectionCount <= 0 {
		// TODO: empty label
		return layout.Dimensions{}
	}
	// sort collections
	page.Daemon.collectionsMutex.Lock()
	page.FilterPage.Sort(page.Daemon.collections)
	page.Daemon.collectionsMutex.Unlock()
	// layout
	childs := make([]layout.Widget, 0)
	for i := 0; i < collectionCount; i++ {
		c := page.Daemon.CollectionAtIndex(i)
		if page.FilterPage.Filter(c) {
			childs = append(childs, func(gtx layout.Context) layout.Dimensions {
				return c.LayoutMinimal(gtx, theme, pages)
			})
		}
	}
	return theme.LayoutListSpaced(gtx, &page.List, theme.SmallVSpacer, childs...)
}
