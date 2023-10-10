package main

import (
	"nftsiren/cmd/nftsiren/widgets"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Page interface {
	// Title will be visible for child-pages, it is not visible on top pages
	Title() string
	// Will be called once before the first Layout call for a page
	Entering()
	// Will be called once after the last Layout call for a page
	Leaving()
	// Layout will be called every frame while this page is visible
	Layout(layout.Context, *Theme, *PageStack) layout.Dimensions
}

type PageStack struct {
	goBack   widget.Clickable
	stack    []Page
	prevPage Page
	nextPage Page
	slider   widgets.Slider
}

func (stack *PageStack) Len() int {
	return len(stack.stack)
}

func (stack *PageStack) HasPage() bool {
	return stack.Len() > 0
}

func (stack *PageStack) Push(page Page) {
	// TODO: do not allow same page twice in stack
	if stack.Len() > 0 {
		stack.slider.PushLeft()
	}
	stack.stack = append(stack.stack, page)
	stack.nextPage = page
	RefreshWindowChan <- struct{}{}
}

func (stack *PageStack) Pop() Page {
	if stack.Len() > 0 {
		stack.slider.PushRight()
		last := stack.Last()
		stack.stack = stack.stack[:stack.Len()-1]
		stack.prevPage = last
		RefreshWindowChan <- struct{}{}
		return last
	}
	return nil
}

func (stack *PageStack) Clear() {
	if stack.Len() > 0 {
		last := stack.Last()
		stack.stack = stack.stack[:0]
		stack.prevPage = last
		RefreshWindowChan <- struct{}{}
	}
}

func (stack *PageStack) Last() Page {
	if stack.Len() > 0 {
		return stack.stack[stack.Len()-1]
	}
	return nil
}

// Caller must make sure stack is not empty
func (stack *PageStack) Layout(gtx layout.Context, theme *Theme) layout.Dimensions {
	if stack.prevPage != nil {
		stack.prevPage.Leaving()
		stack.prevPage = nil
	}
	if stack.nextPage != nil {
		stack.nextPage.Entering()
		stack.nextPage = nil
	}
	if stack.Len() > 0 {
		return stack.slider.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			if stack.Len() == 1 {
				// Root-page
				return stack.Last().Layout(gtx, theme, stack)
			} else {
				// Child-page
				return stack.LayoutAsChild(gtx, theme)
			}
		})
	}
	return material.H2(theme.Material(), "ERROR: No pages in stack").Layout(gtx)
}

func (stack *PageStack) LayoutAsChild(gtx layout.Context, theme *Theme) layout.Dimensions {
	last := stack.Last()
	// Do not pop before Last because last may become nil
	if stack.goBack.Clicked() {
		stack.Pop()
	}
	// Layout with last child page, panics if last is nil
	return layout.Flex{
		Axis:      layout.Vertical,
		Spacing:   layout.SpaceEnd,
		Alignment: layout.Start,
	}.Layout(gtx,
		// Header
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Spacing:   layout.SpaceEnd,
				Alignment: layout.Start,
			}.Layout(gtx,
				// Back button
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return theme.IconButton(theme.GoBackIcon, &stack.goBack).Layout(gtx)
					})
				}),
				layout.Rigid(layout.Spacer{Width: theme.MediumSpace}.Layout),
				// Title
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					label := material.H5(theme.Material(), last.Title())
					label.MaxLines = 1
					return label.Layout(gtx)
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: theme.LargeSpace}.Layout),
		// Body
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return theme.MediumInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return last.Layout(gtx, theme, stack)
			})
		}),
	)
}
