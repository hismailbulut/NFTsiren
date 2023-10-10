package widgets

/*
import (
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

type NavigationState struct {
	Navbar NavbarState
	Pages  []layout.Widget
}

type NavigationStyle struct {
	Theme *material.Theme
	State *NavigationState
	Axis  layout.Axis
}

func Navigation(theme *material.Theme, state *NavigationState, axis layout.Axis) NavigationStyle {
	return NavigationStyle{
		Theme: theme,
		State: state,
		Axis:  axis,
	}
}

func (state *NavigationState) AddPage(text string, icon paint.ImageOp, w layout.Widget) {
	state.Navbar.AddButton(text, icon)
	state.Pages = append(state.Pages, w)
}

func (style NavigationStyle) Layout(gtx layout.Context) layout.Dimensions {
	childs := []layout.FlexChild{
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// return layout.UniformInset(8).Layout(gtx,
			return style.State.Pages[style.State.Navbar.Active](gtx)
			// )
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Draw a 1px line between navbar and page
			return DividerLine(gtx, style.Axis, 1, Lighter(style.Theme.Bg))
			// return component.Divider(style.Theme).Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return Navbar(style.Theme, &style.State.Navbar, OppositeAxis(style.Axis)).Layout(gtx)
		}),
	}
	if style.Axis == layout.Horizontal {
		childs[0], childs[2] = childs[2], childs[0]
	}
	return layout.Flex{
		Axis:    style.Axis,
		Spacing: layout.SpaceEvenly,
		// Alignment: layout.Middle,
	}.Layout(gtx, childs...)
}
*/
