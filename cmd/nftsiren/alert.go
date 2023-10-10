package main

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"nftsiren/cmd/nftsiren/alerts"
	"nftsiren/cmd/nftsiren/config"
	"nftsiren/cmd/nftsiren/widgets"
	"nftsiren/pkg/bench"
	"nftsiren/pkg/log"
	"nftsiren/pkg/mutex"
	"nftsiren/pkg/number"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"gioui.org/x/notify"
)

// Alert encapsulates alerts.Alert and also implements it
type Alert[T alerts.Alert] struct {
	handle mutex.Value[T]
	// latest notification time
	last mutex.Value[time.Time]
	// checker is mutex locked, do not use any of the alerts methods, use handle directly
	checker func() bool
	// will be called when removing this alert
	onRemove func()
	// will be called constantly while layouting
	layouter func(layout.Context, *Theme) layout.Dimensions
	// for gui
	remove widget.Clickable
}

func (alert *Alert[T]) String() string {
	return alert.Handle().String()
}

func (alert *Alert[T]) Description() string {
	return alert.Handle().Description()
}

func (alert *Alert[T]) NeedsInterval() bool {
	return alert.Handle().NeedsInterval()
}

func (alert *Alert[T]) Interval() int {
	return alert.Handle().Interval()
}

func (alert *Alert[T]) Looping() bool {
	return alert.Handle().Looping()
}

func (alert *Alert[T]) Handle() T {
	return alert.handle.Load()
}

// Check will send notification if required and returns ok on notification
// number argument is here because of implementing alerts.Alert, it is noop
func (alert *Alert[T]) Check(number.Number) bool {
	// Check cooldown
	notifCooldown := config.LoadFallback[int]("notificationCooldown", 60)
	if time.Since(alert.last.Load()) < time.Duration(notifCooldown)*time.Second {
		// Do not notify continuously
		log.Info().Printf("Passed check for %s because cooldown", alert.handle.String())
		return false
	}
	// Check
	ok := alert.checker()
	// Set last notification time
	if ok {
		alert.last.Store(time.Now())
		// Remove this alert if it's not looping
		if !alert.Looping() {
			alert.onRemove()
		}
	}
	return ok
}

func (alert *Alert[T]) NotificationText() string {
	return alert.Handle().NotificationText()
}

func (alert *Alert[T]) Layout(gtx layout.Context, theme *Theme) layout.Dimensions {
	return alert.layouter(gtx, theme)
}

func (alert *Alert[T]) defaultLayouter(gtx layout.Context, theme *Theme) layout.Dimensions {
	return theme.Background(gtx, theme.DarkerBg, func(gtx layout.Context) layout.Dimensions {
		return theme.SmallInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			spacer := layout.Rigid(theme.MediumHSpacer.Layout)
			return layout.Flex{
				Axis: layout.Horizontal,
				// Spacing:   layout.SpaceBetween,
				Alignment: layout.Middle,
			}.Layout(gtx,
				spacer,
				// alert icon
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return theme.AlertIcon(alert.Handle()).Layout(gtx, theme.IconSize, theme.Fg)
				}),
				// center
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Horizontal,
						Spacing:   layout.SpaceSides,
						Alignment: layout.Middle,
					}.Layout(gtx,
						// description
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Body2(theme.Material(), alert.Description())
							label.Alignment = text.Middle
							label.Color = theme.MediumImpFg
							return label.Layout(gtx)
						}),
						// space
						layout.Rigid(layout.Spacer{Width: theme.MediumSpace}.Layout),
						// loop
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if !alert.Looping() {
								return layout.Dimensions{}
							}
							return theme.LoopIcon.Layout(gtx, theme.IconSize, theme.Fg)
						}),
					)
				}),
				// remove button
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if alert.remove.Clicked() {
						TODO("Show are you sure modal")
						alert.onRemove()
					}
					return theme.IconButton(theme.DeleteIcon, &alert.remove).Layout(gtx)
				}),
				spacer,
			)
		})
	})
}

func NewEthAlert(handle alerts.EthereumAlert, daemon *Daemon) *Alert[alerts.EthereumAlert] {
	alert := &Alert[alerts.EthereumAlert]{}
	alert.handle.Store(handle)
	alert.checker = func() bool {
		if !daemon.GasTracker.EthStillValid() {
			return false
		}
		if alert.Handle().Check(daemon.GasTracker.GetEth()) {
			notify.Push(NAME, alert.NotificationText())
			return true
		}
		return false
	}
	alert.onRemove = func() {
		go daemon.RemoveEthAlert(alert)
	}
	alert.layouter = alert.defaultLayouter
	return alert
}

func NewGasAlert(handle alerts.GasAlert, daemon *Daemon) *Alert[alerts.GasAlert] {
	alert := &Alert[alerts.GasAlert]{}
	alert.handle.Store(handle)
	alert.checker = func() bool {
		if !daemon.GasTracker.GasStillValid() {
			return false
		}
		if alert.Handle().Check(daemon.GasTracker.GetGas()) {
			notify.Push(NAME, alert.NotificationText())
			return true
		}
		return false
	}
	alert.onRemove = func() {
		go daemon.RemoveGasAlert(alert)
	}
	alert.layouter = alert.defaultLayouter
	return alert
}

func NewCollectionAlert(handle alerts.CollectionAlert, collection *Collection) *Alert[alerts.CollectionAlert] {
	alert := &Alert[alerts.CollectionAlert]{}
	alert.handle.Store(handle)
	alert.checker = func() bool {
		checkresult := false
		switch alert.Handle().Type {
		case alerts.CollectionAlertTypeFloorLessThan, alerts.CollectionAlertTypeFloorGreaterThan:
			// Direct check against floor
			floor, ok := collection.Floor()
			if ok {
				checkresult = alert.Handle().Check(floor)
			}
		case alerts.CollectionAlertTypeSalesGreaterThan:
			// Compare current TotalSales with previous TotalSales
			// First check interval
			// TODO
		}
		if checkresult {
			name, _ := collection.Name()
			title := fmt.Sprintf("%s | %s", collection.Market.Load(), name)
			txt := alert.NotificationText()
			notify.Push(title, txt)
			return true
		}
		return false
	}
	alert.onRemove = func() {
		go collection.RemoveAlert(alert)
	}
	alert.layouter = alert.defaultLayouter
	return alert
}

type AlertCreationPage struct {
	title    string
	onCreate func(t alerts.Condition, n number.Number, loop bool)
	List     widget.List
	Type     widgets.TypedEnum[alerts.Condition]
	Value    component.TextField
	Loop     widget.Bool
	Error    error
	Ok       widget.Clickable
}

func NewAlertCreationPage(title string, onCreate func(t alerts.Condition, n number.Number, loop bool), types ...alerts.Condition) *AlertCreationPage {
	page := &AlertCreationPage{
		title:    title,
		onCreate: onCreate,
	}
	page.List.Axis = layout.Vertical
	page.Type.SetKeys(types...)
	page.Value.SingleLine = true
	page.Value.Submit = true // TODO: check submit
	page.Value.InputHint = key.HintNumeric
	page.Value.Filter = "0123456789."
	return page
}

func (page *AlertCreationPage) Title() string {
	return page.title
}

func (page *AlertCreationPage) Entering() {
	log.Debug().Println("Entering", page.title)
}

func (page *AlertCreationPage) Leaving() {
	// reset everything
	page.Type.State.Value = ""
	page.Value.SetText("")
	page.Loop.Value = false
	page.Error = nil
}

func (page *AlertCreationPage) Layout(gtx layout.Context, theme *Theme, pages *PageStack) layout.Dimensions {
	if page.Ok.Clicked() {
		err := page.AddAlert()
		if err != nil {
			page.Error = err
		} else {
			pages.Pop()
		}
	}
	return theme.LayoutForm(gtx, &page.List, &page.Ok,
		// Type label
		func(gtx layout.Context) layout.Dimensions {
			return material.Body1(theme.Material(), "Type").Layout(gtx)
		},
		// Type select
		func(gtx layout.Context) layout.Dimensions {
			return page.Type.Layout(gtx, theme.Material())
		},
		// Value entry
		func(gtx layout.Context) layout.Dimensions {
			hint := "Select type first"
			t, ok := page.Type.SelectedType()
			if ok {
				hint = t.Label()
			} else {
				// Disable widget
				gtx.Queue = nil
			}
			return page.Value.Layout(gtx, theme.Material(), hint)
		},
		// Loop
		func(gtx layout.Context) layout.Dimensions {
			return material.CheckBox(theme.Material(), &page.Loop, "Loop").Layout(gtx)
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

func (page *AlertCreationPage) AddAlert() error {
	// Validate type
	t, ok := page.Type.SelectedType()
	if !ok {
		return errors.New("select type")
	}
	// Validate number
	n, ok := number.NewFromString(page.Value.Text())
	if !ok {
		return errors.New("invalid number")
	}
	// Create alert
	page.onCreate(t, n, page.Loop.Value)
	return nil
}

// can be initialized by new(AlertList)
type AlertList struct {
	mutex sync.Mutex
	slice []alerts.Alert
}

func (alertList *AlertList) Len() int {
	alertList.mutex.Lock()
	defer alertList.mutex.Unlock()
	return len(alertList.slice)
}

// Has reports whether this alert is in the list
func (alertList *AlertList) Has(alert alerts.Alert) bool {
	alertList.mutex.Lock()
	defer alertList.mutex.Unlock()
	for _, a := range alertList.slice {
		if a.String() == alert.String() {
			return true
		}
	}
	return false
}

func (alertList *AlertList) Add(alert alerts.Alert) bool {
	if alertList.Has(alert) {
		return false
	}
	// Append and return true
	alertList.mutex.Lock()
	defer alertList.mutex.Unlock()
	alertList.slice = append(alertList.slice, alert)
	return true
}

func (alertList *AlertList) Remove(alert alerts.Alert) bool {
	alertList.mutex.Lock()
	defer alertList.mutex.Unlock()
	for i, a := range alertList.slice {
		if a.String() == alert.String() {
			alertList.slice = append(alertList.slice[:i], alertList.slice[i+1:]...)
			return true
		}
	}
	return false
}

func (alertList *AlertList) ForEach(fn func(index int, alert alerts.Alert)) {
	alertList.mutex.Lock()
	for i, alert := range alertList.slice {
		fn(i, alert)
	}
	alertList.mutex.Unlock()
}

type AlertListState struct {
	Title string
	// List              widget.List
	AddNewAlert       widget.Clickable
	AlertCreationPage *AlertCreationPage
}

func (state *AlertListState) Layout(gtx layout.Context, theme *Theme, pages *PageStack, alertList *AlertList) layout.Dimensions {
	defer bench.Begin()()

	if state.AddNewAlert.Clicked() {
		pages.Push(state.AlertCreationPage)
	}

	count := max(alertList.Len(), 1) + 1
	widgets := make([]layout.FlexChild, count)

	widgets[0] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Spacing:   layout.SpaceSides,
			Alignment: layout.Middle,
		}.Layout(gtx,
			// Title
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.H6(theme.Material(), state.Title)
				label.Alignment = text.Middle
				return label.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Width: theme.MediumSpace}.Layout),
			// Button
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return theme.IconButton(theme.NewAlarmIcon, &state.AddNewAlert).Layout(gtx)
			}),
		)
	})

	if alertList.Len() <= 0 {
		// TODO: we can store this label
		label := material.Body1(theme.Material(), "It's empty here")
		label.Color = theme.LowImpFg
		label.Alignment = text.Middle
		label.MaxLines = 1
		widgets[1] = layout.Rigid(label.Layout)
	}

	type hasLayout interface {
		Layout(layout.Context, *Theme) layout.Dimensions
	}

	alertList.ForEach(func(index int, alert alerts.Alert) {
		widgets[index+1] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return alert.(hasLayout).Layout(gtx, theme)
		})
	})

	spacer := layout.Spacer{Height: theme.SmallSpace}
	return theme.LayoutFlexListSpaced(gtx, spacer, widgets...)
}
