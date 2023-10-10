package main

import (
	"time"

	"nftsiren/cmd/nftsiren/widgets"
	"nftsiren/pkg/apis/etherscan"
	"nftsiren/pkg/log"
	"nftsiren/pkg/mutex"
	"nftsiren/pkg/number"
	"nftsiren/pkg/worker"

	"gioui.org/layout"
	"gioui.org/widget/material"
)

type GasTracker struct {
	worker *worker.Worker
	// TODO: solana price in usd
	// TODO: eip-1159 gas (base fee and priority fee)
	// TODO: move ethereum price to some other price api
	eth           mutex.Value[etherscan.EthPrice]
	ethUpdateTime mutex.Value[time.Time]
	gas           mutex.Value[etherscan.GasPrice]
	gasUpdateTime mutex.Value[time.Time]
}

func NewGasTracker() *GasTracker {
	tracker := &GasTracker{}
	tracker.worker = worker.New(worker.Settings{
		Name:              "GasTracker",
		Interval:          time.Second * 5,
		Work:              tracker.FetchEthInfo,
		InitialRun:        true,
		PanicHanler:       ReportPanic,
		RestartAfterPanic: true,
		MaxPanics:         3,
	})
	return tracker
}

func (tracker *GasTracker) Start() {
	tracker.worker.Start()
}

func (tracker *GasTracker) Stop() {
	tracker.worker.Stop()
}

func (tracker *GasTracker) FetchEthInfo() {
	eth, err := etherscan.FetchEthPrice()
	if err != nil {
		log.Warn().Println("Failed to fetch eth price from etherscan:", err)
	} else {
		tracker.eth.Store(eth)
		tracker.ethUpdateTime.Store(time.Now())
	}
	gas, err := etherscan.FetchGasPrice()
	if err != nil {
		log.Warn().Println("Failed to fetch gas price from etherscan:", err)
	} else {
		tracker.gas.Store(gas)
		tracker.gasUpdateTime.Store(time.Now())
	}
	RefreshWindowChan <- struct{}{}
}

func (tracker *GasTracker) EthStillValid() bool {
	updateTime := tracker.ethUpdateTime.Load()
	const d = time.Second * 10
	return !updateTime.IsZero() && time.Since(updateTime) < d && tracker.GetEth().Int64() > 0
}

func (tracker *GasTracker) GasStillValid() bool {
	updateTime := tracker.gasUpdateTime.Load()
	const d = time.Second * 10
	return !updateTime.IsZero() && time.Since(updateTime) < d && tracker.GetGas().Int64() > 0
}

func (tracker *GasTracker) GetEth() number.Number {
	return tracker.eth.Load().Ethusd
}

func (tracker *GasTracker) GetGas() number.Number {
	return tracker.gas.Load().ProposeGasPrice // Average
}

// A generic layout for gas tracker, you don't have to use it
func (tracker *GasTracker) Layout(gtx layout.Context, theme *Theme, axis layout.Axis) layout.Dimensions {
	return layout.Flex{
		Axis:    axis,
		Spacing: layout.SpaceEvenly,
		// Alignment: layout.Middle,
	}.Layout(gtx,
		// Ethereum
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return tracker.layoutPrice(gtx, theme, theme.EthereumIcon, tracker.GetEth().StringPretty()+"$")
		}),
		// Gas
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return tracker.layoutPrice(gtx, theme, theme.GasIcon, tracker.GetGas().StringPretty())
		}),
	)
}

func (tracker *GasTracker) layoutPrice(gtx layout.Context, theme *Theme, icon *widgets.Icon, price string) layout.Dimensions {
	// return theme.SmallInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Horizontal,
		// Spacing: layout.SpaceSides,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(layout.Spacer{Width: theme.SmallSpace}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return icon.Layout(gtx, theme.IconSize, theme.Fg)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Body2(theme.Material(), price)
			// label.Alignment = text.Middle
			label.MaxLines = 1
			return label.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Width: theme.SmallSpace}.Layout),
	)
	// })
}
