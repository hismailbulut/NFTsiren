package widgets

/*
import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type DropdownState struct {
	Group  widget.Enum
	Button widget.Clickable
	List   widget.List
	Open   bool
	// animation
	anim    Animation
	reverse bool // will be true when need to animate closing
}

type DropdownStyle struct {
	Theme *material.Theme
	State *DropdownState
	Hint  string
	Keys  []string

	TextSize            unit.Sp
	BorderColorInactive color.NRGBA
	BorderColorHovered  color.NRGBA
	BorderColorActive   color.NRGBA
	BorderWidthInactive unit.Dp
	BorderWidthActive   unit.Dp
	CornerRadius        unit.Dp
	ActiveInset         layout.Inset
	KeyInset            layout.Inset
	ModalBackground     color.NRGBA
	KeyColorHovered     color.NRGBA
}

func Dropdown(theme *material.Theme, state *DropdownState, hint string, keys ...string) DropdownStyle {
	state.List.Axis = layout.Vertical
	return DropdownStyle{
		Theme: theme,
		State: state,
		Hint:  hint,
		Keys:  keys,
		// These matches with component.TextField
		TextSize:            theme.TextSize,
		BorderColorInactive: WithAlpha(theme.Fg, 128),
		BorderColorHovered:  WithAlpha(theme.Fg, 221),
		BorderColorActive:   theme.ContrastBg,
		BorderWidthInactive: unit.Dp(0.5),
		BorderWidthActive:   unit.Dp(2),
		CornerRadius:        unit.Dp(4),
		ActiveInset:         layout.UniformInset(unit.Dp(12)),
		KeyInset:            layout.UniformInset(unit.Dp(6)),
		ModalBackground:     Darker(theme.Bg),
		KeyColorHovered:     WithAlpha(theme.ContrastBg, 20),
	}
}

func (style DropdownStyle) Layout(gtx layout.Context) layout.Dimensions {
	openPreviously := style.State.Open
	buttonClicked := style.State.Button.Clicked()
	buttonHasFocus := style.State.Button.Focused()
	_, groupHasFocus := style.State.Group.Focused()
	_, groupHasHover := style.State.Group.Hovered()
	groupChanged := style.State.Group.Changed()
	switch {
	case buttonClicked:
		// Toggle modal on clicks
		style.State.Open = !style.State.Open
	case groupChanged:
		// Close if it changed
		style.State.Open = false
	case !(buttonHasFocus || groupHasFocus || groupHasHover):
		// Close if modal has no focus
		style.State.Open = false
	}
	// Open/Close animation
	style.State.anim.Update(gtx.Now)
	if openPreviously && !style.State.Open {
		style.State.reverse = true
		style.State.anim.Start(gtx.Now, ANIM_DURATION)
	} else if !openPreviously && style.State.Open {
		style.State.reverse = false
		style.State.anim.Start(gtx.Now, ANIM_DURATION)
	}
	// colors
	borderColor := style.BorderColorInactive
	borderWidth := style.BorderWidthInactive
	if style.State.Button.Focused() {
		borderColor = style.BorderColorActive
		borderWidth = style.BorderWidthActive
	} else if style.State.Button.Hovered() {
		borderColor = style.BorderColorHovered
	}
	// Layout entry
	dims := Border(gtx, borderColor, style.CornerRadius, borderWidth, func(gtx layout.Context) layout.Dimensions {
		return style.State.Button.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Spacing:   layout.SpaceBetween,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					text := style.Hint
					if style.State.Group.Value != "" {
						text = style.State.Group.Value
					}
					return style.ActiveInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return material.Label(style.Theme, style.TextSize, text).Layout(gtx)
						})
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// Draw triangle
					width := float32(gtx.Sp(style.Theme.TextSize))
					p := clip.Path{}
					p.Begin(gtx.Ops)
					var size f32.Point
					if style.State.Open {
						size = f32.Pt(width, width/2)
						p.Move(f32.Pt(-size.X/4, 0))
						p.Line(f32.Pt(size.X, 0))
						p.Line(f32.Pt(-size.Y, size.Y))
						p.Line(f32.Pt(-size.Y, -size.Y))
					} else {
						size = f32.Pt(width/2, width)
						p.Line(f32.Pt(0, size.Y))
						p.Line(f32.Pt(-size.X, -size.X))
						p.Line(f32.Pt(size.X, -size.X))
					}
					p.Close()
					paint.FillShape(gtx.Ops, borderColor, clip.Outline{Path: p.End()}.Op())
					return layout.Dimensions{Size: size.Round()}
				}),
			)
		})
	})
	// Draw selection list modal
	progress := style.State.anim.Progress()
	if style.State.Open || progress < 1 {
		modalSize := image.Pt(dims.Size.X, dims.Size.Y*5)
		macro, _ := Macro(gtx, func(gtx layout.Context) layout.Dimensions {
			// Calculate animation if we are animating
			if progress < 1 {
				if style.State.reverse {
					progress = 1 - progress
				}
				op.InvalidateOp{At: gtx.Now.Add(ANIM_NEXTTIME)}.Add(gtx.Ops)
				offset := int(float32(modalSize.Y) * progress)
				clipRect := image.Rect(0, dims.Size.Y, modalSize.X, dims.Size.Y+offset)
				defer clip.Rect(clipRect).Push(gtx.Ops).Pop()
				defer op.Offset(image.Pt(0, offset-modalSize.Y)).Push(gtx.Ops).Pop()
			}
			// offset modal by field size
			defer op.Offset(image.Pt(0, dims.Size.Y)).Push(gtx.Ops).Pop()
			// layout modal
			gtx.Constraints = layout.Exact(modalSize)
			return BackgroundRect(gtx, style.ModalBackground, style.CornerRadius, func(gtx layout.Context) layout.Dimensions {
				return material.List(style.Theme, &style.State.List).Layout(gtx, len(style.Keys),
					func(gtx layout.Context, index int) layout.Dimensions {
						macro, dims := Macro(gtx, func(gtx layout.Context) layout.Dimensions {
							return style.State.Group.Layout(gtx, style.Keys[index], func(gtx layout.Context) layout.Dimensions {
								return style.KeyInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return material.Label(style.Theme, style.Theme.TextSize, style.Keys[index]).Layout(gtx)
								})
							})
						})
						hoveredKey, hovered := style.State.Group.Hovered()
						if hovered && hoveredKey == style.Keys[index] {
							FillRect(gtx, style.KeyColorHovered, image.Rect(0, 0, dims.Size.X, dims.Size.Y), style.CornerRadius)
						}
						macro.Add(gtx.Ops)
						return dims
					},
				)
			})
		})
		op.Defer(gtx.Ops, macro)
	}
	return dims
}

func (state *DropdownState) Selected() string {
	return state.Group.Value
}

type TypedDropdown[T fmt.Stringer] struct {
	State DropdownState
	Types []T
	Keys  []string
}

func (td *TypedDropdown[T]) SetKeys(types ...T) {
	td.Types = types
	td.Keys = make([]string, len(td.Types))
	for i, t := range td.Types {
		td.Keys[i] = t.String()
	}
}

func (td *TypedDropdown[T]) Layout(gtx layout.Context, theme *material.Theme, hint string) layout.Dimensions {
	return Dropdown(theme, &td.State, hint, td.Keys...).Layout(gtx)
}

func (td *TypedDropdown[T]) SelectedType() (T, bool) {
	selected := td.State.Selected()
	for i, k := range td.Keys {
		if selected == k {
			return td.Types[i], true
		}
	}
	return *new(T), false
}
*/