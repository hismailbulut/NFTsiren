package widgets

/*
import (
	"image"
	"nftsiren/pkg/log"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

type DynamicGridState component.GridState

type DynamicGridStyle struct {
	Theme       *material.Theme
	State       *DynamicGridState
	MinWidth    unit.Dp
	MinHeight   unit.Dp
	SpaceEvenly bool
}

func DynamicGrid(theme *material.Theme, state *DynamicGridState, minWidth, minHeight unit.Dp) *DynamicGridStyle {
	return &DynamicGridStyle{
		Theme:       theme,
		State:       state,
		MinWidth:    minWidth,
		MinHeight:   minHeight,
		SpaceEvenly: true,
	}
}

func (dyn *DynamicGridStyle) Layout(gtx layout.Context, count int, layouter func(layout.Context, int) layout.Dimensions) layout.Dimensions {
	maxWidth := gtx.Constraints.Max.X
	width := gtx.Dp(dyn.MinWidth)
	height := gtx.Dp(dyn.MinHeight)
	// calculate column and row count
	cols := maxWidth / width
	if cols <= 0 {
		cols = 1
	}
	rows := count / cols
	if count%cols > 0 {
		rows++
	}
	log.Debug().Field("constraints", gtx.Constraints).Println("items:", count, "cols:", cols, "rows:", rows)
	// calculate remaining width of row and share it between widgets
	if dyn.SpaceEvenly {
		width += (maxWidth - (width * cols)) / cols
	}
	size := image.Pt(width, height)
	return component.Grid(dyn.Theme, (*component.GridState)(dyn.State)).Layout(gtx, rows, cols,
		func(axis layout.Axis, index, constraint int) int {
			return height
		},
		func(gtx layout.Context, row, col int) layout.Dimensions {
			index := row*cols + col
			if index >= count {
				return layout.Dimensions{Size: size}
			}
			gtx.Constraints = layout.Exact(size)
			return layouter(gtx, index)
		},
	)
}
*/