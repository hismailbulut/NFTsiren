package main

import (
	"errors"
	"fmt"
	"image"
	"net/url"
	"time"

	"nftsiren/cmd/nftsiren/alerts"
	"nftsiren/cmd/nftsiren/cache"
	"nftsiren/cmd/nftsiren/widgets"
	"nftsiren/pkg/apis"
	"nftsiren/pkg/images"
	"nftsiren/pkg/log"
	"nftsiren/pkg/mutex"
	"nftsiren/pkg/nft"
	"nftsiren/pkg/number"
	"nftsiren/pkg/worker"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Collection struct {
	Daemon *Daemon
	Market mutex.Value[nft.Marketplace] // Which marketplace this collection will be fetched from
	Symbol mutex.Value[string]          // Unique id of this collection, maybe contract address
	// worker runs the given func constantly in given period
	worker *worker.Worker
	// alerts of this collection
	alerts *AlertList
	// updated at runtime
	err   mutex.Value[error]                // Whether an error happened while fetching this collection
	info  mutex.Value[*nft.Collection]      // We only need to fetch this first time
	stats mutex.Value[*nft.CollectionStats] // This will be updated on every check
	// gui stuff
	img    mutex.Value[*widgets.Icon] // May be nil on error or while loading
	imgErr mutex.Value[error]         // Image fetching or parsing error
	// gui state variables dont require mutex because they will only be used in rendering thread
	// alertsButton   widget.Clickable
	detailsList    widget.List
	detailsButton  widget.Clickable // This must be controlled by parent widget
	deleteButton   widget.Clickable
	alertListState AlertListState
	statsState     CollectionStatsState
}

func NewCollection(daemon *Daemon, market nft.Marketplace, slug string) *Collection {
	collection := &Collection{
		Daemon: daemon,
		alerts: new(AlertList),
	}
	collection.Market.Store(market)
	collection.Symbol.Store(slug)
	collection.worker = worker.New(worker.Settings{
		Name:        collection.String(),
		Interval:    time.Minute,
		Work:        collection.FetchAndCheck,
		InitialRun:  true,
		PanicHanler: ReportPanic,
	})
	// gui initilization
	collection.detailsList.Axis = layout.Vertical
	// alertList initialization
	collection.alertListState.Title = "Alerts"
	collection.alertListState.AlertCreationPage = NewAlertCreationPage("New Collection Alert",
		func(t alerts.Condition, n number.Number, loop bool) {
			params := alerts.CollectionAlert{
				Type: t.(alerts.CollectionAlertType),
				Base: n,
				Loop: loop,
			}
			alert := NewCollectionAlert(params, collection)
			go collection.AddAlert(alert) // TODO: we may not need to go
		},
		[]alerts.Condition{
			alerts.CollectionAlertTypeFloorLessThan,
			alerts.CollectionAlertTypeFloorGreaterThan,
			alerts.CollectionAlertTypeSalesGreaterThan,
		}...,
	)
	return collection
}

func (collection *Collection) Start() {
	collection.worker.Start()
}

func (collection *Collection) Stop() {
	collection.worker.Stop()
}

func (collection *Collection) String() string {
	return fmt.Sprintf("Collection(%s:%s)", collection.Market.Load(), collection.Symbol.Load())
}

// This will return the slug if collection name is not known yet
func (collection *Collection) Name() (string, bool) {
	info := collection.info.Load()
	if info != nil && info.Name != "" {
		return info.Name, true
	}
	return collection.Symbol.Load(), false
}

func (collection *Collection) Floor() (number.Number, bool) {
	stats := collection.stats.Load()
	if stats == nil {
		return number.Number{}, false
	}
	if stats.Floor.IsNil() {
		return number.Number{}, false
	}
	return stats.Floor, true
}

func (collection *Collection) HasValidInfo() bool {
	info := collection.info.Load()
	if info == nil {
		return false
	}
	// We have to figure out whether every collection has an image
	return info.Name != "" && info.ImageURL != ""
}

func (collection *Collection) HasValidStats() bool {
	stats := collection.stats.Load()
	if stats == nil {
		return false
	}
	return stats.IsValid()
}

func (collection *Collection) NumAlerts() int {
	return collection.alerts.Len()
}

func (collection *Collection) AddAlert(alert *Alert[alerts.CollectionAlert]) bool {
	return collection.alerts.Add(alert)
}

func (collection *Collection) RemoveAlert(alert *Alert[alerts.CollectionAlert]) bool {
	return collection.alerts.Remove(alert)
}

func (collection *Collection) FetchAndCheck() {
	collection.Fetch()
	collection.Check()
}

// Fetches full collection if it is not fetched already,
// Otherwise only fetches the collection stats
func (collection *Collection) Fetch() {
	// Only fetch collection once
	if !collection.HasValidInfo() {
		collection.FetchCollection()
	}
	// Always fetch stats if info is fetched
	if collection.info.Load() != nil {
		collection.FetchStats()
	}
	RefreshWindowChan <- struct{}{}
}

func (collection *Collection) FetchCollection() {
	info, err := apis.FetchCollection(collection.Market.Load(), collection.Symbol.Load())
	collection.err.Store(err)
	// Check error
	if err != nil {
		log.Warn().Printf("Failed to fetch %s: %s", collection, err)
		return
	}
	// Debug
	/*
		var desc string
		if len(info.Description) > 40 {
			desc = info.Description[:40]
		} else {
			desc = info.Description
		}
		desc = strings.ReplaceAll(desc, "\n", " ")
		desc = strings.TrimSpace(desc)
		log.Debug().
			Field("time", info.Time).
			Field("currency", info.Currency).
			Field("marketplace", info.Marketplace).
			Field("symbol", info.Symbol).
			Field("address", info.Address).
			Field("name", info.Name).
			Field("description", desc).
			Field("imageurl", info.ImageURL).
			Field("marketpage", info.Marketpage).
			Field("website", info.Website).
			Field("twitter", info.Twitter).
			Field("discord", info.Discord).
			Println("Collection fetched")
	*/
	// Set collection info
	collection.info.Store(&info)
	// Download collection image
	go collection.FetchImage(info.ImageURL)
}

func (collection *Collection) FetchImage(imgURL string) {
	if imgURL == "" {
		collection.setImage(nil, errors.New("no image url"))
		return
	}
	if _, err := url.Parse(imgURL); err != nil {
		collection.setImage(nil, fmt.Errorf("invalid image url: %s", imgURL))
		return
	}
	// Try to load from cache
	img, err := cache.LoadImage(imgURL)
	if err != nil {
		log.Warn().Println("Couldn't load cached image:", err)
	} else if img != nil {
		collection.setImage(img, nil)
		return
	}
	// Download image and shrink to reduce ram usage
	// Because this will be done once in a while, we can use catmull-rom to create high quality images
	log.Debug().Println("Downloading collection image:", imgURL)
	const maxImageSize = 256
	img, err = images.DownloadAndShrink(imgURL, maxImageSize)
	collection.setImage(img, err)
	if err != nil {
		log.Warn().Println("Couldn't download image:", err)
		return
	}
	// Cache this image
	err = cache.SaveImage(imgURL, img)
	if err != nil {
		log.Warn().Println("Failed to cache image:", err)
	}
}

func (collection *Collection) reFetchImage() {
	if collection.hasImage() {
		return
	}
	info := collection.info.Load()
	if info == nil {
		return
	}
	go collection.FetchImage(info.ImageURL)
}

func (collection *Collection) setImage(img *image.RGBA, err error) {
	var icon *widgets.Icon
	if img != nil {
		icon = widgets.NewIconFromImage(img)
	}
	collection.img.Store(icon)
	collection.imgErr.Store(err)
	RefreshWindowChan <- struct{}{}
}

func (collection *Collection) hasImage() bool {
	return collection.img.Load() != nil
}

func (collection *Collection) FetchStats() {
	info := collection.info.Load()
	assert(info != nil, "collection info must not nil here")
	if info.Stats != nil && info.Stats.IsValid() && info.Stats.IsRecent(time.Minute) {
		collection.stats.Store(info.Stats)
		// Free info.Stats otherwise we have to check this everytime
		info.Stats = nil
		return
	}
	// Fetch additionally
	stats, err := apis.FetchCollectionStats(collection.Market.Load(), collection.Symbol.Load())
	if err != nil {
		log.Warn().Println("Failed to fetch", collection, "stats:", err)
		return
	}
	if !stats.IsValid() {
		log.Warn().Printf("%v stats is not valid %+v", collection, stats)
		return
	}
	collection.stats.Store(&stats)
}

func (collection *Collection) Check() {
	collection.alerts.ForEach(func(index int, alert alerts.Alert) {
		alert.Check(number.Number{})
	})
}

// This is for implementing ChildPage
func (collection *Collection) Title() string {
	name, _ := collection.Name()
	return name
}

func (collection *Collection) Entering() {}

func (collection *Collection) Leaving() {}

// Layouts much more detailed child page
func (collection *Collection) Layout(gtx layout.Context, theme *Theme, pages *PageStack) layout.Dimensions {
	items := make([]layout.Widget, 0)
	// Header
	items = append(items, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(gtx,
			// Image
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return collection.layoutImage(gtx, theme, unit.Dp(theme.TextSize*6))
				})
			}),
			layout.Rigid(theme.MediumHSpacer.Layout),
			// Header
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				// return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return collection.layoutDetailedHeader(gtx, theme, pages)
				// })
			}),
		)
	})
	// Update time
	/*
		items = append(items, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(gtx,
					// Marketplace logo
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return theme.MarketplaceLogo(collection.Market.Load()).Layout(gtx, theme.IconSize*0.75, theme.Fg)
					}),
					// Space
					layout.Rigid(theme.SmallHSpacer.Layout),
					// Marketplace name
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := material.Subtitle2(theme.Material(), collection.Market.Load().String())
						return label.Layout(gtx)
					}),
					// Space
					layout.Rigid(theme.SmallHSpacer.Layout),
					// Last update
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						stats := collection.stats.Load()
						text := "Outdated"
						color := theme.Error
						if stats != nil && stats.IsRecent(time.Minute*2) {
							text = "Up to date"
							color = theme.Link
						}
						label := material.Overline(theme.Material(), text)
						label.Color = color
						return label.Layout(gtx)
					}),
				)
			})
		})
	*/
	// Error message
	if collection.err.Load() != nil {
		items = append(items, func(gtx layout.Context) layout.Dimensions {
			errLabel := material.Body2(theme.Material(), collection.err.Load().Error())
			errLabel.Alignment = text.Middle
			errLabel.Color = theme.Error
			return layout.Center.Layout(gtx, errLabel.Layout)
		})
	}
	// Image error message
	if collection.imgErr.Load() != nil {
		items = append(items, func(gtx layout.Context) layout.Dimensions {
			errLabel := material.Body2(theme.Material(), collection.imgErr.Load().Error())
			errLabel.Alignment = text.Middle
			errLabel.Color = theme.Error
			return layout.Center.Layout(gtx, errLabel.Layout)
		})
	}
	// Stats
	items = append(items, func(gtx layout.Context) layout.Dimensions {
		return collection.statsState.Layout(gtx, theme,
			collection.Market.Load(),
			collection.info.Load(),
			collection.stats.Load(),
		)
	})
	/*
		if collection.HasValidStats() {
			items = append(items, func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return collection.layoutStats(gtx, theme)
				})
			})
		}
	*/
	// Alerts
	items = append(items, func(gtx layout.Context) layout.Dimensions {
		return collection.alertListState.Layout(gtx, theme, pages, collection.alerts)
	})
	return theme.LayoutListSpaced(gtx, &collection.detailsList, theme.LargeVSpacer, items...)
}

// Minimal stats with only image, name and floor
func (collection *Collection) LayoutMinimal(gtx layout.Context, theme *Theme, pages *PageStack) layout.Dimensions {
	return theme.Background(gtx, theme.DarkerBg, func(gtx layout.Context) layout.Dimensions {
		return theme.SmallInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			if collection.detailsButton.Clicked() {
				// collection implements ChildPage
				pages.Push(collection)
			}
			return collection.detailsButton.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Horizontal,
					// Spacing: layout.SpaceBetween,
					Alignment: layout.Middle,
				}.Layout(gtx,
					// Image
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return collection.layoutImage(gtx, theme, unit.Dp(theme.TextSize*3))
					}),
					// Header
					layout.Flexed(0.7, func(gtx layout.Context) layout.Dimensions {
						return collection.layoutMinimalHeader(gtx, theme)
					}),
					// Floor price
					layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
						if !collection.HasValidStats() {
							return layout.Dimensions{}
						}
						// Layout floor price
						return collection.layoutFloorPrice(gtx, theme)
					}),
				)
			})
		})
	})
}

func (collection *Collection) layoutImage(gtx layout.Context, theme *Theme, size unit.Dp) layout.Dimensions {
	if collection.img.Load() == nil {
		if collection.err.Load() != nil || collection.imgErr.Load() != nil {
			// Error while fetching image or collection, or when parsing image
			return theme.BrokenIcon.Layout(gtx, size, theme.LowImpFg)
		} else {
			// Most probably loading the image
			p := gtx.Dp(size)
			gtx.Constraints = layout.Exact(image.Pt(p, p))
			return layout.Center.Layout(gtx, material.Loader(theme.Material()).Layout)
		}
	} else {
		return collection.img.Load().Layout(gtx, size, theme.Fg)
	}
}

func (collection *Collection) layoutDetailedHeader(gtx layout.Context, theme *Theme, pages *PageStack) layout.Dimensions {
	return theme.Background(gtx, theme.DarkerBg, func(gtx layout.Context) layout.Dimensions {
		return theme.SmallInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				// Name and delete button
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(gtx,
						// Name
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							name, _ := collection.Name()
							title := material.Subtitle1(theme.Material(), name)
							title.Alignment = text.Middle
							title.MaxLines = 1
							return title.Layout(gtx)
						}),
						// Space
						layout.Rigid(theme.MediumHSpacer.Layout),
						// Delete button
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if collection.deleteButton.Clicked() {
								TODO("show are you sure dialog")
								// delete collection and close the page
								go collection.Daemon.RemoveCollection(collection)
								pages.Pop()
							}
							return theme.IconButton(theme.DeleteIcon, &collection.deleteButton).Layout(gtx)
						}),
					)
				}),
				// Space
				layout.Rigid(theme.MediumVSpacer.Layout),
				// Description
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					info := collection.info.Load()
					if info == nil {
						return layout.Dimensions{}
					}
					desc := material.Body2(theme.Material(), info.Description)
					desc.MaxLines = 3
					desc.Color = theme.MediumImpFg
					return desc.Layout(gtx)
				}),
			)
		})
	})
}

func (collection *Collection) layoutMinimalHeader(gtx layout.Context, theme *Theme) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		// Space
		layout.Rigid(layout.Spacer{Width: theme.LargeSpace}.Layout),
		// Markeplace logo
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return theme.MarketplaceLogo(collection.Market.Load()).Layout(gtx, theme.IconSize*0.75, theme.Fg)
		}),
		// Space
		layout.Rigid(layout.Spacer{Width: theme.MediumSpace}.Layout),
		// Name
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			name, _ := collection.Name()
			title := material.Subtitle1(theme.Material(), name)
			// title.Alignment = text.Middle
			// title.TextSize = theme.TextSize
			title.MaxLines = 1
			return title.Layout(gtx)
		}),
	)
}

func (collection *Collection) layoutFloorPrice(gtx layout.Context, theme *Theme) layout.Dimensions {
	info := collection.info.Load()
	assert(info != nil, "info is nil")
	stats := collection.stats.Load()
	assert(stats != nil, "stats is nil")
	floor := stats.Floor
	assert(!floor.IsNil(), "%v floor is nil, other stats: %+v", collection, stats)
	floorText := floor.StringPretty() + " " + info.Currency.String()
	floorLabel := material.Body1(theme.Material(), floorText)
	floorLabel.Alignment = text.End
	floorLabel.Color = theme.ContrastBg
	// usd price for ethereum only
	// TODO: solana
	if info.Currency == nft.ETH && collection.Daemon.GasTracker.EthStillValid() {
		floorUsd := floor.Mul(collection.Daemon.GasTracker.GetEth())
		floorUsdLabel := material.Body2(theme.Material(), floorUsd.StringFixed(0)+"$")
		floorUsdLabel.Alignment = text.End
		floorUsdLabel.Color = theme.MediumImpFg
		return layout.Flex{
			Axis:    layout.Vertical,
			Spacing: layout.SpaceEvenly,
			// Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Rigid(floorLabel.Layout),
			layout.Rigid(floorUsdLabel.Layout),
		)
	}
	return floorLabel.Layout(gtx)
}

/*
func (collection *Collection) layoutStats(gtx layout.Context, theme *Theme) layout.Dimensions {
	info := collection.info.Load()
	assert(info != nil, "info is nil")
	stats := collection.stats.Load()
	assert(stats != nil, "stats is nil")
	availableStats := stats.All()
	const cols = 4
	rows := ((len(availableStats) - 1) / cols) + 1
	widgets := make([]layout.Widget, len(availableStats))
	for i := range availableStats {
		info := availableStats[i]
		widgets[i] = func(gtx layout.Context) layout.Dimensions {
			return theme.SmallInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Vertical,
					Spacing:   layout.SpaceEvenly,
					Alignment: layout.Middle,
				}.Layout(gtx,
					// Label
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := material.Body2(theme.Material(), info.Label)
						label.Alignment = text.Middle
						label.Color = theme.MediumImpFg
						return label.Layout(gtx)
					}),
					// Value
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := material.Body1(theme.Material(), info.Value.StringPretty())
						label.Alignment = text.Middle
						return label.Layout(gtx)
					}),
				)
			})
		}
	}
	return theme.Background(gtx, theme.DarkerBg, func(gtx layout.Context) layout.Dimensions {
		return theme.SmallInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return theme.LayoutFlexGrid(gtx, rows, cols, widgets...)
		})
	})
}
*/
