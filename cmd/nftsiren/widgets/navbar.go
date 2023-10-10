package widgets

import (
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type NavbarButton struct {
	Text   string
	Icon   *Icon
	Button widget.Clickable
}

func (button *NavbarButton) Layout(gtx layout.Context, state *NavbarState, style *NavbarStyle, index int) layout.Dimensions {
	return style.ButtonInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// return material.Clickable(gtx, &button.Button, func(gtx layout.Context) layout.Dimensions {
		return button.Button.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			// semantic.Button.Add(gtx.Ops)
			pointer.CursorPointer.Add(gtx.Ops)
			if button.Button.Clicked() {
				state.Active = index
				state.changed = true
			}
			fg := style.ForegroundInactive
			if state.Active == index {
				fg = style.ForegroundActive
			}
			return layout.Flex{
				Axis:      layout.Vertical,
				Spacing:   layout.SpaceSides,
				Alignment: layout.Middle,
			}.Layout(gtx,
				// icon
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return button.Icon.Layout(gtx, style.IconSize, fg)
				}),
				// padding
				layout.Rigid(style.ButtonSpacer.Layout),
				// text
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.LabelStyle{
						Color:     fg,
						Alignment: text.Middle,
						MaxLines:  1,
						Text:      button.Text,
						TextSize:  style.TextSize,
						Shaper:    style.Theme.Shaper,
					}.Layout(gtx)
				}),
			)
		})
	})
}

type NavbarState struct {
	Buttons []NavbarButton
	Active  int
	changed bool
}

type NavbarStyle struct {
	Theme              *material.Theme
	State              *NavbarState
	Axis               layout.Axis
	ForegroundInactive color.NRGBA
	ForegroundActive   color.NRGBA
	IconSize           unit.Dp
	IconColor          color.NRGBA
	TextSize           unit.Sp
	ButtonColorActive  color.NRGBA
	ButtonInset        layout.Inset
	ButtonSpacer       layout.Spacer
}

func Navbar(theme *material.Theme, state *NavbarState, axis layout.Axis) NavbarStyle {
	return NavbarStyle{
		Theme:              theme,
		State:              state,
		Axis:               axis,
		ForegroundInactive: MulAlpha(theme.Fg, 0.25),
		ForegroundActive:   theme.ContrastBg,
		IconSize:           unit.Dp(24),
		IconColor:          theme.Fg,
		TextSize:           theme.TextSize,
		ButtonInset:        layout.UniformInset(unit.Dp(6)),
		ButtonSpacer:       layout.Spacer{Height: unit.Dp(4)},
	}
}

func (state *NavbarState) AddButton(text string, icon *Icon) {
	state.Buttons = append(state.Buttons, NavbarButton{
		Text: text,
		Icon: icon,
	})
}

func (state *NavbarState) Changed() bool {
	if state.changed {
		state.changed = false
		return true
	}
	return false
}

func (style NavbarStyle) Layout(gtx layout.Context) layout.Dimensions {
	if len(style.State.Buttons) == 0 {
		return layout.Dimensions{}
	}
	buttons := make([]layout.FlexChild, len(style.State.Buttons))
	for i := range style.State.Buttons {
		index := i
		buttons[i] = layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return style.State.Buttons[index].Layout(gtx, style.State, &style, index)
		})
	}
	spacing := layout.SpaceEvenly
	alignment := layout.Middle
	// if style.Axis == layout.Vertical {
	// 	spacing = layout.SpaceEnd
	// 	alignment = layout.Start
	// }
	// return layout.UniformInset(4).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis:      style.Axis,
		Spacing:   spacing,
		Alignment: alignment,
	}.Layout(gtx, buttons...)
	// })
}
