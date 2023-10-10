package main

import (
	"nftsiren/cmd/nftsiren/alerts"
	"nftsiren/pkg/bench"
	"nftsiren/pkg/log"
	"nftsiren/pkg/number"

	"gioui.org/layout"
	"gioui.org/widget"
)

// TODO: rename this page and extend it's usage
type HomePage struct {
	Daemon *Daemon
	// States
	List     widget.List
	Ethereum AlertListState
	Gas      AlertListState
}

func NewHomePage(daemon *Daemon) *HomePage {
	page := &HomePage{
		Daemon: daemon,
	}
	page.List.Axis = layout.Vertical
	// init ethereum alerts
	page.Ethereum.Title = "Ethereum Alerts"
	page.Ethereum.AlertCreationPage = NewAlertCreationPage("New Ethereum Alert",
		func(t alerts.Condition, n number.Number, loop bool) {
			params := alerts.EthereumAlert{
				Type: t.(alerts.EthereumAlertType),
				Base: n,
				Loop: loop,
			}
			alert := NewEthAlert(params, page.Daemon)
			go page.Daemon.AddEthAlert(alert)
		},
		[]alerts.Condition{
			alerts.EthereumAlertTypeLessThan,
			alerts.EthereumAlertTypeGreaterThan,
		}...,
	)
	// init gas alerts
	page.Gas.Title = "Gas Alerts"
	page.Gas.AlertCreationPage = NewAlertCreationPage("New Gas Alert",
		func(t alerts.Condition, n number.Number, loop bool) {
			params := alerts.GasAlert{
				Type: t.(alerts.GasAlertType),
				Base: n,
				Loop: loop,
			}
			alert := NewGasAlert(params, page.Daemon)
			go page.Daemon.AddGasAlert(alert)
		},
		[]alerts.Condition{
			alerts.GasAlertTypeLessThan,
			alerts.GasAlertTypeGreaterThan,
		}...,
	)
	return page
}

func (page *HomePage) Title() string {
	return "Home"
}

func (page *HomePage) Entering() {
	log.Debug().Println("Entering Home")
}

func (page *HomePage) Leaving() {
	log.Debug().Println("Leaving Home")
}

func (page *HomePage) Layout(gtx layout.Context, theme *Theme, pages *PageStack) layout.Dimensions {
	defer bench.Begin()()
	// TODO: overview
	return theme.LayoutListSpaced(gtx, &page.List, theme.MediumVSpacer,
		func(gtx layout.Context) layout.Dimensions {
			return page.layoutEthAlertList(gtx, theme, pages)
		},
		func(gtx layout.Context) layout.Dimensions {
			return page.layoutGasAlertList(gtx, theme, pages)
		},
	)
}

func (page *HomePage) layoutEthAlertList(gtx layout.Context, theme *Theme, pages *PageStack) layout.Dimensions {
	return page.Ethereum.Layout(gtx, theme, pages, page.Daemon.ethAlerts)
}

func (page *HomePage) layoutGasAlertList(gtx layout.Context, theme *Theme, pages *PageStack) layout.Dimensions {
	return page.Gas.Layout(gtx, theme, pages, page.Daemon.gasAlerts)
}
