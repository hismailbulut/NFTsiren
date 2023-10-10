package widgets

import (
	"image"
	"image/color"

	"nftsiren/pkg/bench"
	"nftsiren/pkg/images"
	"nftsiren/pkg/log"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/inkeliz/giosvg"
)

type Icon struct {
	// iconvg icon source
	// TODO: we should cache this icons size's
	icon *widget.Icon
	// svg icon source
	svg *giosvg.Icon
	// image icon source
	image *image.RGBA
	ops   map[int]paint.ImageOp
}

func NewIconFromImage(img *image.RGBA) *Icon {
	return &Icon{
		image: img,
		ops:   make(map[int]paint.ImageOp),
	}
}

func NewIconFromIconVG(data []byte) *Icon {
	return &Icon{
		icon: mustIcon(data),
	}
}

func NewIconFromSvgVector(vector giosvg.Vector) *Icon {
	return &Icon{
		svg: giosvg.NewIcon(vector),
	}
}

// clr is only required for iconvg icons
func (icon *Icon) Layout(gtx layout.Context, size unit.Dp, clr color.NRGBA) layout.Dimensions {
	if size <= 0 {
		size = 20
	}
	p := gtx.Dp(size)
	gtx.Constraints = layout.Exact(image.Pt(p, p))
	if icon.image != nil {
		op := icon.ImageOp(p)
		return widget.Image{
			Src:      op,
			Fit:      widget.Contain,
			Position: layout.Center,
			Scale:    1,
		}.Layout(gtx)
	} else if icon.icon != nil {
		return icon.icon.Layout(gtx, clr)
	} else if icon.svg != nil {
		return icon.svg.Layout(gtx)
	}
	// This should not happen
	panic("icon has no source")
}

func (icon *Icon) WidgetIcon() *widget.Icon {
	return icon.icon
}

func (icon *Icon) ImageOp(size int) paint.ImageOp {
	defer bench.Begin()()
	op, ok := icon.ops[size]
	if !ok {
		// Do not create new image if image is already in the size we want
		var resized *image.RGBA
		b := icon.image.Bounds()
		if (b.Dx() == size && b.Dy() <= size) || (b.Dx() <= size && b.Dy() == size) {
			// log.Debug().Println("image size is same with render size, using default image")
			resized = icon.image
		} else {
			resized = images.ResizeKeepRatio(icon.image, size)
		}
		op = paint.NewImageOp(resized)
		icon.ops[size] = op
		log.Debug().Println("New image created with size:", size)
	}
	return op
}

func mustIcon(data []byte) *widget.Icon {
	icon, err := widget.NewIcon(data)
	if err != nil {
		panic(err)
	}
	return icon
}
