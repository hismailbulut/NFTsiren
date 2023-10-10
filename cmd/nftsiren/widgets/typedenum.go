package widgets

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type TypedEnum[T fmt.Stringer] struct {
	State  widget.Enum
	Types  []T
	Keys   []string
	Childs []layout.FlexChild
}

func (td *TypedEnum[T]) SetKeys(types ...T) {
	td.Types = types
	td.Keys = make([]string, len(td.Types))
	for i, t := range td.Types {
		td.Keys[i] = t.String()
	}
}

func (td *TypedEnum[T]) SelectedType() (T, bool) {
	selected := td.State.Value
	for i, k := range td.Keys {
		if selected == k {
			return td.Types[i], true
		}
	}
	return *new(T), false
}

func (td *TypedEnum[T]) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	if len(td.Keys) == 0 {
		return layout.Dimensions{}
	}
	if len(td.Childs) != len(td.Keys) {
		td.Childs = make([]layout.FlexChild, len(td.Keys))
		for i, k := range td.Keys {
			k := k
			td.Childs[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.RadioButton(theme, &td.State, k, k).Layout(gtx)
			})
		}
	}
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx, td.Childs...)
}
