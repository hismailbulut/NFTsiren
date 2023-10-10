package widgets

import (
	"image"
	"image/color"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/x/component"
)

const (
	ANIM_NEXTTIME = time.Second / 144
	ANIM_DURATION = time.Second / 10
)

func Macro(gtx layout.Context, w layout.Widget) (op.CallOp, layout.Dimensions) {
	macro := op.Record(gtx.Ops)
	dims := w(gtx)
	return macro.Stop(), dims
}

// Draws a shadow around the widget
func Shadow(gtx layout.Context, radius unit.Dp, elevation unit.Dp, w layout.Widget) layout.Dimensions {
	macro, dims := Macro(gtx, w)
	gtx.Constraints = layout.Exact(dims.Size)
	component.Shadow(radius, elevation).Layout(gtx)
	macro.Add(gtx.Ops)
	return dims
}

func DividerLine(gtx layout.Context, axis layout.Axis, width int, fg color.NRGBA) layout.Dimensions {
	// TODO: make shadowed
	end := f32.Point{}
	var size image.Point
	if axis == layout.Horizontal {
		end.Y = float32(gtx.Constraints.Min.Y)
		size = image.Pt(width, gtx.Constraints.Min.Y)
	} else {
		end.X = float32(gtx.Constraints.Min.X)
		size = image.Pt(gtx.Constraints.Min.X, width)
		// size = image.Pt(width, gtx.Constraints.Max.Y)
	}
	p := clip.Path{}
	p.Begin(gtx.Ops)
	p.Line(end)
	p.Close()
	paint.FillShape(gtx.Ops, fg, clip.Stroke{Path: p.End(), Width: float32(width)}.Op())
	return layout.Dimensions{Size: size}
}

func DrawBorder(gtx layout.Context, c color.NRGBA, rect image.Rectangle, radius unit.Dp, width unit.Dp) {
	var path clip.PathSpec
	if radius <= 0 {
		path = clip.Rect(rect).Path()
	} else {
		path = clip.UniformRRect(rect, gtx.Dp(radius)).Path(gtx.Ops)
	}
	paint.FillShape(gtx.Ops, c, clip.Stroke{
		Path:  path,
		Width: float32(gtx.Dp(width)),
	}.Op())
}

func Border(gtx layout.Context, c color.NRGBA, radius unit.Dp, width unit.Dp, w layout.Widget) layout.Dimensions {
	macro, dims := Macro(gtx, w)
	DrawBorder(gtx, c, image.Rect(0, 0, dims.Size.X, dims.Size.Y), radius, width)
	macro.Add(gtx.Ops)
	wdp2 := 2 * gtx.Dp(width)
	return layout.Dimensions{Size: image.Pt(dims.Size.X+wdp2, dims.Size.Y+wdp2)}
}

func FillRect(gtx layout.Context, c color.NRGBA, rect image.Rectangle, radius unit.Dp) {
	// rect := image.Rect(pos.X, pos.Y, size.X, size.Y)
	if radius <= 0 {
		paint.FillShape(gtx.Ops, c, clip.Rect(rect).Op())
		return
	}
	paint.FillShape(gtx.Ops, c, clip.UniformRRect(rect, gtx.Dp(radius)).Op(gtx.Ops))
}

func BackgroundRect(gtx layout.Context, c color.NRGBA, radius unit.Dp, w layout.Widget) layout.Dimensions {
	macro, dims := Macro(gtx, w)
	FillRect(gtx, c, image.Rect(0, 0, dims.Size.X, dims.Size.Y), radius)
	macro.Add(gtx.Ops)
	return dims
}
