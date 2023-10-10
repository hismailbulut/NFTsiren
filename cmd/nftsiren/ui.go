package main

import (
	"nftsiren/cmd/nftsiren/widgets"
	"nftsiren/pkg/bench"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

type UI struct {
	Daemon *Daemon
	Theme  *Theme

	Pages PageStack

	HomePage        *HomePage
	CollectionsPage *CollectionsPage
	SettingsPage    *SettingsPage

	Navbar         widgets.NavbarState
	AboutButton    widget.Clickable
	ShowAboutPopup bool

	LayoutPage layout.Widget
}

func NewUI(daemon *Daemon, theme *Theme) *UI {
	ui := &UI{
		Daemon: daemon,
		Theme:  theme,
	}

	ui.HomePage = NewHomePage(ui.Daemon)
	ui.Navbar.AddButton("Home", ui.Theme.HomeIcon)

	ui.CollectionsPage = NewCollectionsPage(ui.Daemon)
	ui.Navbar.AddButton("Collections", ui.Theme.CollectionsIcon)

	ui.SettingsPage = NewSettingsPage(ui.Daemon)
	ui.Navbar.AddButton("Settings", ui.Theme.SettingsIcon)

	if isMobile() {
		ui.LayoutPage = ui.LayoutPageMobile
	} else {
		ui.LayoutPage = ui.LayoutPageDesktop
	}

	ui.Navbar.Active = 1

	return ui
}

func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	defer bench.Begin()()
	// Fill background
	paint.Fill(gtx.Ops, ui.Theme.Bg)
	// Layout active page
	return ui.LayoutPage(gtx)
}

func (ui *UI) LayoutPageDesktop(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Spacing:   layout.SpaceEvenly,
		Alignment: layout.Middle,
	}.Layout(gtx,
		// Navbar and stats on the left
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Vertical,
				// Spacing:   layout.SpaceBetween,
				Alignment: layout.Middle,
			}.Layout(gtx,
				// Icon
				// TODO: our icon will open about page
				layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
					if ui.AboutButton.Clicked() {
						ui.ShowAboutPopup = true
					}
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return ui.AboutButton.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return ui.Theme.SirenIcon.Layout(gtx, 64, ui.Theme.Fg)
						})
					})
				}),
				// Gas area
				layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
					return ui.Daemon.GasTracker.Layout(gtx, ui.Theme, layout.Vertical)
				}),
				// Divider
				// layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				// 	return widgets.DividerLine(gtx, layout.Horizontal, 1, widgets.Lighter(ui.Theme.Bg))
				// }),
				// Navbar
				layout.Flexed(0.8, func(gtx layout.Context) layout.Dimensions {
					return widgets.Navbar(ui.Theme.Material(), &ui.Navbar, layout.Vertical).Layout(gtx)
				}),
			)
		}),
		// Divider
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return widgets.DividerLine(gtx, layout.Horizontal, 1, widgets.Lighter(ui.Theme.Bg))
		}),
		// Page on the right
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return ui.Theme.MediumInset.Layout(gtx, ui.LayoutActivePage)
		}),
	)
}

func (ui *UI) LayoutPageMobile(gtx layout.Context) layout.Dimensions {
	panic("mobile not implemented yet")
}

func (ui *UI) LayoutActivePage(gtx layout.Context) layout.Dimensions {
	if ui.Navbar.Changed() {
		ui.Pages.Clear()
	}
	if !ui.Pages.HasPage() {
		switch ui.Navbar.Active {
		case 0:
			ui.Pages.Push(ui.HomePage)
		case 1:
			ui.Pages.Push(ui.CollectionsPage)
		case 2:
			ui.Pages.Push(ui.SettingsPage)
		default:
			panic("unknown page number")
		}
	}
	return ui.Pages.Layout(gtx, ui.Theme)
}
