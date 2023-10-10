package main

import (
	"nftsiren/pkg/nft"
	"time"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget/material"
)

type CollectionStatsState struct{}

func (state *CollectionStatsState) Layout(gtx layout.Context, theme *Theme, market nft.Marketplace, info *nft.Collection, stats *nft.CollectionStats) layout.Dimensions {
	return theme.Background(gtx, theme.DarkerBg, func(gtx layout.Context) layout.Dimensions {
		return theme.SmallInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(gtx,
						// Marketplace logo
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return theme.MarketplaceLogo(market).Layout(gtx, theme.IconSize*0.75, theme.Fg)
						}),
						// Space
						layout.Rigid(theme.MediumHSpacer.Layout),
						// Marketplace name
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							label := material.Subtitle2(theme.Material(), market.String())
							return label.Layout(gtx)
						}),
						// Space
						// layout.Rigid(theme.SmallHSpacer.Layout),
						// Last update
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
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
					// })
				}),
				layout.Rigid(theme.MediumVSpacer.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if info == nil || stats == nil || !stats.IsValid() {
						// not available
						// return layout.Dimensions{}
						return material.Body1(theme.Material(), "Not available").Layout(gtx)
					}
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
					return theme.LayoutFlexGrid(gtx, rows, cols, widgets...)
				}),
			)
		})
	})
}
