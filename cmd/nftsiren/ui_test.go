package main

import (
	"image"
	"testing"
	"time"

	"gioui.org/io/event"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func Benchmark_LayoutDynamicList(b *testing.B) {
	gtx, theme := mockLayoutNeeds()
	list := widget.List{}
	w := func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		widgets := make([]layout.Widget, 0)
		for j := 0; j < 10; j++ {
			widgets = append(widgets, w)
		}
		theme.LayoutList(gtx, &list, widgets...)
	}
}

func Benchmark_Flex(b *testing.B) {
	gtx, _ := mockLayoutNeeds()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		layout.Flex{}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{}
			}),
		)
	}
}

func Benchmark_LayoutDynamicFlex(b *testing.B) {
	gtx, _ := mockLayoutNeeds()
	w := func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		widgets := make([]layout.FlexChild, 0)
		for j := 0; j < 10; j++ {
			widgets = append(widgets, layout.Rigid(w))
		}
		layout.Flex{}.Layout(gtx, widgets...)
	}
}

func Benchmark_LayoutLabel(b *testing.B) {
	gtx, theme := mockLayoutNeeds()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		material.Label(theme.Material(), theme.TextSize, "Hello World!").Layout(gtx)
	}
}

type mockQueue struct{}

func (q *mockQueue) Events(event.Tag) []event.Event { return nil }

func mockFrameEvent() system.FrameEvent {
	return system.FrameEvent{
		Now:    time.Now(),
		Metric: unit.Metric{PxPerDp: 1, PxPerSp: 1},
		Size:   image.Pt(800, 600),
		Insets: system.Insets{Top: 4, Bottom: 4, Left: 4, Right: 4},
		Frame:  func(frame *op.Ops) {},
		Queue:  &mockQueue{},
	}
}

func mockLayoutNeeds() (gtx layout.Context, thm *Theme) {
	evt := mockFrameEvent()
	ops := op.Ops{}
	gtx = layout.NewContext(&ops, evt)
	thm = DefaultTheme()
	return
}
