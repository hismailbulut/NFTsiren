package widgets

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type ModalStyle struct {
	Background      color.NRGBA
	ModalBackground color.NRGBA
	CornerRadius    unit.Dp
	ShadowElevation unit.Dp
	Inset           layout.Inset
}

func Modal(theme *material.Theme) ModalStyle {
	return ModalStyle{
		Background:      MulAlpha(theme.Bg, 0.75),
		ModalBackground: Darker(theme.Bg),
		CornerRadius:    unit.Dp(4),
		ShadowElevation: unit.Dp(4),
		Inset:           layout.UniformInset(unit.Dp(12)),
	}
}

// This should be called in the root of the page
func (style ModalStyle) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	macro, dims := Macro(gtx, func(gtx layout.Context) layout.Dimensions {
		// Fill background with background color smaller alpha
		paint.Fill(gtx.Ops, style.Background)
		// Modal should be smaller than the area, divide it to golden ratio
		size := gtx.Constraints.Max
		gtx.Constraints.Max.X = int(float32(size.X) / 1.618)
		gtx.Constraints.Max.Y = int(float32(size.Y) / 1.618)
		// Center our modal
		defer op.Offset(size.Sub(gtx.Constraints.Max).Div(2)).Push(gtx.Ops).Pop()
		// Draw shadow around modal
		return Shadow(gtx, style.CornerRadius, style.ShadowElevation, func(gtx layout.Context) layout.Dimensions {
			// Draw modal
			return BackgroundRect(gtx, style.ModalBackground, style.CornerRadius, func(gtx layout.Context) layout.Dimensions {
				// Draw widget with inset
				return style.Inset.Layout(gtx, w)
			})
		})
	})
	op.Defer(gtx.Ops, macro)
	return dims
}
