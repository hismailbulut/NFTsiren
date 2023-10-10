package widgets

/*
import (
	"image"
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type TabbarState struct {
	Group widget.Enum
	List  widget.List
	// animation
	prev int
	next int
	anim Animation
}

type TabbarStyle struct {
	Theme             *material.Theme
	State             *TabbarState
	Keys              []string
	BorderColor       color.NRGBA
	BorderWidth       unit.Dp
	TextSize          unit.Sp
	ButtonColorActive color.NRGBA
	ButtonInset       layout.Inset
	CornerRadius      unit.Dp
}

func Tabbar(theme *material.Theme, state *TabbarState, keys ...string) TabbarStyle {
	return TabbarStyle{
		Theme:             theme,
		State:             state,
		Keys:              keys,
		BorderColor:       Darker(theme.Bg),
		BorderWidth:       unit.Dp(0.5),
		TextSize:          theme.TextSize,
		ButtonColorActive: WithAlpha(theme.ContrastBg, 64),
		ButtonInset:       layout.UniformInset(unit.Dp(4)),
		CornerRadius:      unit.Dp(0),
	}
}

func (state *TabbarState) Active() int {
	return state.next
}

func (state *TabbarState) Previous() int {
	return state.prev
}

func (state *TabbarState) Progress() float32 {
	return state.anim.progress
}

func (style TabbarStyle) Layout(gtx layout.Context) layout.Dimensions {
	keyChanged := style.State.Group.Changed()
	macro, dims := Macro(gtx, func(gtx layout.Context) layout.Dimensions {
		return Border(gtx, style.BorderColor, style.CornerRadius, style.BorderWidth, func(gtx layout.Context) layout.Dimensions {
			return material.List(style.Theme, &style.State.List).Layout(gtx, len(style.Keys),
				func(gtx layout.Context, index int) layout.Dimensions {
					key := style.Keys[index]
					return style.State.Group.Layout(gtx, key, func(gtx layout.Context) layout.Dimensions {
						if keyChanged && style.State.Group.Value == key {
							// start animation if user presses this button and this is not active
							if style.State.next != index {
								style.State.prev = style.State.next
								style.State.next = index
								style.State.anim.Start(gtx.Now, ANIM_DURATION)
							}
						}
						pointer.CursorPointer.Add(gtx.Ops)
						return style.ButtonInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Center.Layout(gtx, material.Label(style.Theme, style.TextSize, key).Layout)
						})
					})
				},
			)
		})
	})
	// Draw animated color bar
	state := style.State
	buttonWidth := dims.Size.X / len(style.Keys)
	// fmt.Println("buttonWidth:", buttonWidth)
	var rect image.Rectangle
	state.anim.Update(gtx.Now)
	progress := state.anim.Progress()
	if progress < 1 {
		// Draw animated
		posNorm := float32(state.prev) + (float32(state.next-state.prev) * progress)
		rect.Min.X = int(posNorm * float32(buttonWidth))
		op.InvalidateOp{At: gtx.Now.Add(ANIM_NEXTTIME)}.Add(gtx.Ops)
	} else {
		// Draw exact location
		rect.Min.X = state.next * buttonWidth
	}
	rect.Max.X = rect.Min.X + buttonWidth
	rect.Min.Y = dims.Size.Y / 4 * 3
	rect.Max.Y = dims.Size.Y
	// Draw border around tabbar
	DrawBorder(gtx, style.BorderColor, image.Rect(0, 0, dims.Size.X, dims.Size.Y), style.CornerRadius, style.BorderWidth)
	// FillRect(gtx, style.Background, image.Rect(0, 0, dims.Size.X, dims.Size.Y), style.CornerRadius)
	// Draw a small rectangle under active tab
	FillRect(gtx, style.ButtonColorActive, rect, style.CornerRadius)
	macro.Add(gtx.Ops)
	return dims
}
*/
