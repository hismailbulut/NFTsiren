package main

import (
	"strconv"

	"nftsiren/cmd/nftsiren/config"
	"nftsiren/cmd/nftsiren/widgets"
	"nftsiren/pkg/bench"
	"nftsiren/pkg/log"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

type SettingsPage struct {
	Daemon *Daemon
	// State
	List widget.List
	// Settings
	NotificationCooldown component.TextField
	StartInBackground    widget.Bool // Silent start
	Autostart            widget.Bool
	// Api key entries
	ApiKeys widget.Clickable
	// Buttons
	// Save     widget.Clickable
	// SaveText string
	// Quit widget.Clickable
	// Debug
	OpenIconGallery widget.Clickable
}

func NewSettingsPage(daemon *Daemon) *SettingsPage {
	settings := &SettingsPage{
		Daemon: daemon,
	}

	settings.List.Axis = layout.Vertical

	settings.NotificationCooldown.SingleLine = true
	settings.NotificationCooldown.Filter = "0123456789"
	settings.NotificationCooldown.InputHint = key.HintNumeric

	return settings
}

func (page *SettingsPage) Title() string {
	return "Settings"
}

func (page *SettingsPage) Entering() {
	log.Debug().Println("Loading preferences")
	// Load user setting
	cooldown := config.LoadFallback[int]("notificationCooldown", 60)
	page.NotificationCooldown.SetText(strconv.Itoa(cooldown))
	page.StartInBackground.Value = config.LoadFallback[bool]("startInBackground", false)
	page.Autostart.Value = config.GetAutostart()
}

func (page *SettingsPage) Leaving() {
	log.Debug().Println("Saving preferences")
	// Save settings
	config.Store("notificationCooldown", mustParseNumber(page.NotificationCooldown.Text()))
	config.Store("startInBackground", page.StartInBackground.Value)
	config.SetAutostart(page.Autostart.Value)
}

func (page *SettingsPage) Layout(gtx layout.Context, theme *Theme, pages *PageStack) layout.Dimensions {
	defer bench.Begin()()

	items := []layout.Widget{}

	items = append(items, func(gtx layout.Context) layout.Dimensions {
		return page.NotificationCooldown.Layout(gtx, theme.Material(), "Notification cooldown in seconds")
	})

	if isDesktop() {
		// Start in background
		items = append(items, material.CheckBox(theme.Material(), &page.StartInBackground,
			"Start in background. The program will be started in background without disturbing you.").Layout)
		// Autostart at system startup
		items = append(items, material.CheckBox(theme.Material(), &page.Autostart,
			"Autostart with system. The program will be automatically started just after system startup. Changing executable path may broke this feature.").Layout)
	}

	// Api keys button
	if page.ApiKeys.Clicked() {
		pages.Push(&ApiKeysPage{Daemon: page.Daemon})
	}
	items = append(items, theme.Hyperlink("Api Keys", &page.ApiKeys))
	// items = append(items, func(gtx layout.Context) layout.Dimensions {
	// 	return page.ApiKeys.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
	// 		pointer.CursorPointer.Add(gtx.Ops)
	// 		label := material.Body1(theme.Material(), "Api Keys")
	// 		label.Alignment = text.End
	// 		label.Color = theme.Link
	// 		return label.Layout(gtx)
	// 	})
	// })

	// Save button
	/*
		items = append(items, func(gtx layout.Context) layout.Dimensions {
			if page.Save.Clicked() {
				prefs.SetSilentStart(page.StartInBackground.Value)
			}
			return theme.Button("Save", &page.Save).Layout(gtx)
		})
		if page.SaveText != "" {
			items = append(items, func(gtx layout.Context) layout.Dimensions {
				label := material.Body1(theme.Material(), page.SaveText)
				label.Color = theme.Link
				return label.Layout(gtx)
			})
		}
	*/
	// Quit
	/*
		items = append(items, func(gtx layout.Context) layout.Dimensions {
			if page.Quit.Clicked() {
				TODO("quit")
			}
			return theme.IconButton(theme.QuitIcon, &page.Quit).Layout(gtx)
		})
	*/

	// DEBUG MODE
	if isDebug() {
		items = append(items, layout.Spacer{Height: theme.MediumSpace}.Layout)
		items = append(items, material.Subtitle1(theme.Material(), "DEBUG MODE").Layout)
		if page.OpenIconGallery.Clicked() {
			pages.Push(&IconGalleryPage{})
		}
		items = append(items, theme.Button("Show icon gallery", &page.OpenIconGallery, RegularButton).Layout)
	}

	return theme.LayoutListSpaced(gtx, &page.List, theme.MediumVSpacer, items...)
}

func mustParseNumber(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

type ApiKeys struct {
	Etherscan  string `json:"etherscan,omitempty"`
	Opensea    string `json:"opensea,omitempty"`
	Blur       string `json:"blur,omitempty"`
	Looksrare  string `json:"looksrare,omitempty"`
	Magiceden  string `json:"magiceden,omitempty"`
	Solanart   string `json:"solanart,omitempty"`
	Reservoir  string `json:"reservoir,omitempty"`
	Simplehash string `json:"simplehash,omitempty"`
}

type ApiKeysPage struct {
	Daemon     *Daemon
	List       widget.List
	Etherscan  component.TextField
	Opensea    component.TextField
	Blur       component.TextField
	Looksrare  component.TextField
	Magiceden  component.TextField
	Solanart   component.TextField
	Reservoir  component.TextField
	Simplehash component.TextField
}

func (page *ApiKeysPage) Title() string {
	return "Api Keys"
}

func (page *ApiKeysPage) Entering() {
	page.List.Axis = layout.Vertical
	page.Etherscan.SingleLine = true
	page.Opensea.SingleLine = true
	page.Blur.SingleLine = true
	page.Looksrare.SingleLine = true
	page.Magiceden.SingleLine = true
	page.Solanart.SingleLine = true
	page.Reservoir.SingleLine = true
	page.Simplehash.SingleLine = true
	// Load current keys
	keys, err := config.Load[ApiKeys]("apiKeys")
	if err == nil {
		page.Etherscan.SetText(keys.Etherscan)
		page.Opensea.SetText(keys.Opensea)
		page.Blur.SetText(keys.Blur)
		page.Looksrare.SetText(keys.Looksrare)
		page.Magiceden.SetText(keys.Magiceden)
		page.Solanart.SetText(keys.Solanart)
		page.Reservoir.SetText(keys.Reservoir)
		page.Simplehash.SetText(keys.Simplehash)
	}
}

func (apiKeys *ApiKeysPage) Leaving() {
	config.Store("apiKeys", ApiKeys{
		Etherscan:  apiKeys.Etherscan.Text(),
		Opensea:    apiKeys.Opensea.Text(),
		Blur:       apiKeys.Blur.Text(),
		Looksrare:  apiKeys.Looksrare.Text(),
		Magiceden:  apiKeys.Magiceden.Text(),
		Solanart:   apiKeys.Solanart.Text(),
		Reservoir:  apiKeys.Reservoir.Text(),
		Simplehash: apiKeys.Simplehash.Text(),
	})
	apiKeys.Daemon.ResetApiKeys()
}

func (apiKeys *ApiKeysPage) Layout(gtx layout.Context, theme *Theme, pages *PageStack) layout.Dimensions {
	items := make([]layout.Widget, 0)
	// Credentials
	items = append(items, func(gtx layout.Context) layout.Dimensions {
		return apiKeys.Etherscan.Layout(gtx, theme.Material(), "Etherscan Api Key (Required)")
	})
	items = append(items, func(gtx layout.Context) layout.Dimensions {
		return apiKeys.Opensea.Layout(gtx, theme.Material(), "Opensea Api Key (Required)")
	})
	items = append(items, func(gtx layout.Context) layout.Dimensions {
		return apiKeys.Looksrare.Layout(gtx, theme.Material(), "Looksrare Api Key (Required)")
	})
	items = append(items, func(gtx layout.Context) layout.Dimensions {
		return apiKeys.Magiceden.Layout(gtx, theme.Material(), "Magiceden Api Key (Required)")
	})
	return theme.LayoutListSpaced(gtx, &apiKeys.List, theme.MediumVSpacer, items...)
}

// DEBUG

type NamedIcon struct {
	Name string
	Icon *widgets.Icon
}

type IconGalleryPage struct {
	List  widget.List
	Icons []NamedIcon
}

func (page *IconGalleryPage) Title() string {
	return "Icon Gallery"
}

func (page *IconGalleryPage) Entering() {
	page.List.Axis = layout.Vertical
	// Load icons
	for _, icon := range iconList {
		page.Icons = append(page.Icons, NamedIcon{
			Name: icon.name,
			Icon: widgets.NewIconFromIconVG(icon.data),
		})
	}
}

func (page *IconGalleryPage) Leaving() {
	// Free icons
	page.Icons = nil
}

func (page *IconGalleryPage) Layout(gtx layout.Context, theme *Theme, pages *PageStack) layout.Dimensions {
	widgets := make([]layout.Widget, len(page.Icons))
	for i, icon := range page.Icons {
		icon := icon
		widgets[i] = func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Horizontal,
				// Spacing:   layout.SpaceEvenly,
				Alignment: layout.Middle,
			}.Layout(gtx,
				// Icon
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return icon.Icon.Layout(gtx, theme.IconSize*3, theme.Fg)
				}),
				layout.Rigid(layout.Spacer{Width: theme.LargeSpace}.Layout),
				// Label
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					label := material.Body2(theme.Material(), icon.Name)
					label.Alignment = text.Middle
					label.Color = theme.MediumImpFg
					return label.Layout(gtx)
				}),
			)
		}
	}
	return theme.LayoutListSpaced(gtx, &page.List, theme.SmallVSpacer, widgets...)
}
