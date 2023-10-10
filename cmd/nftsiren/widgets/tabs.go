package widgets

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Tab struct {
	Button widget.Clickable
	Title  string
	Widget layout.Widget
}

type TabsState struct {
	List     layout.List
	Tabs     []Tab
	Slider   Slider
	Selected int
}

type TabsStyle struct {
	Theme     *material.Theme
	State     *TabsState
	TextSize  unit.Sp
	TextInset layout.Inset
}

func Tabs(theme *material.Theme, state *TabsState) TabsStyle {
	return TabsStyle{
		Theme:     theme,
		State:     state,
		TextSize:  theme.TextSize,
		TextInset: layout.UniformInset(8),
	}
}

func (state *TabsState) AddTab(text string, w layout.Widget) {
	state.Tabs = append(state.Tabs, Tab{
		Title:  text,
		Widget: w,
	})
}

func (style TabsStyle) Layout(gtx layout.Context) layout.Dimensions {
	state := style.State
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return state.List.Layout(gtx, len(state.Tabs), func(gtx layout.Context, tabIdx int) layout.Dimensions {
				t := &state.Tabs[tabIdx]
				if t.Button.Clicked() {
					if state.Selected < tabIdx {
						state.Slider.PushLeft()
					} else if state.Selected > tabIdx {
						state.Slider.PushRight()
					}
					state.Selected = tabIdx
				}
				var tabWidth int
				return layout.Stack{Alignment: layout.S}.Layout(gtx,
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						dims := material.Clickable(gtx, &t.Button, func(gtx layout.Context) layout.Dimensions {
							return style.TextInset.Layout(gtx, material.Label(style.Theme, style.TextSize, t.Title).Layout)
						})
						tabWidth = dims.Size.X
						return dims
					}),
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						if state.Selected != tabIdx {
							return layout.Dimensions{}
						}
						tabHeight := gtx.Dp(unit.Dp(4))
						tabRect := image.Rect(0, 0, tabWidth, tabHeight)
						paint.FillShape(gtx.Ops, style.Theme.ContrastBg, clip.Rect(tabRect).Op())
						return layout.Dimensions{
							Size: image.Point{X: tabWidth, Y: tabHeight},
						}
					}),
				)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return state.Slider.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// fill(gtx, dynamicColor(tabs.selected), dynamicColor(tabs.selected+1))
				// return layout.Center.Layout(gtx,
				// 	material.H1(th, fmt.Sprintf("Tab content #%d", tabs.selected+1)).Layout,
				// )
				// fmt.Println("selected:", state.Selected)
				// if state.Selected >= 0 && state.Selected < len(state.Tabs) {
				return state.Tabs[state.Selected].Widget(gtx)
				// }
				// return layout.Dimensions{}
			})
		}),
	)
}

/*
func fill(gtx layout.Context, col1, col2 color.NRGBA) {
	dr := image.Rectangle{Max: gtx.Constraints.Min}
	paint.FillShape(gtx.Ops,
		color.NRGBA{R: 0, G: 0, B: 0, A: 0xFF},
		clip.Rect(dr).Op(),
	)

	col2.R = byte(float32(col2.R))
	col2.G = byte(float32(col2.G))
	col2.B = byte(float32(col2.B))
	paint.LinearGradientOp{
		Stop1:  f32.Pt(float32(dr.Min.X), 0),
		Stop2:  f32.Pt(float32(dr.Max.X), 0),
		Color1: col1,
		Color2: col2,
	}.Add(gtx.Ops)
	defer clip.Rect(dr).Push(gtx.Ops).Pop()
	paint.PaintOp{}.Add(gtx.Ops)
}

func dynamicColor(i int) color.NRGBA {
	sn, cs := math.Sincos(float64(i) * math.Phi)
	return color.NRGBA{
		R: 0xA0 + byte(0x30*sn),
		G: 0xA0 + byte(0x30*cs),
		B: 0xD0,
		A: 0xFF,
	}
}
*/
