package images

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"net/http"

	"nftsiren/pkg/bench"
	"nftsiren/pkg/images/webp"

	// "github.com/tidbyt/go-libwebp/webp"
	"golang.org/x/image/bmp"
	"golang.org/x/image/draw"
)

const MAX_IMAGE_SIZE = 1024 * 1024 * 32 // 32 mb

func Download(url string) (*image.RGBA, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Security check, do not download images larger than MAX_IMAGE_SIZE
	data, err := io.ReadAll(io.LimitReader(resp.Body, MAX_IMAGE_SIZE))
	if err != nil {
		return nil, err
	}
	return Parse(data)
}

func DownloadAndShrink(url string, maxDim int) (*image.RGBA, error) {
	img, err := Download(url)
	if err != nil {
		return nil, err
	}
	return Shrink(img, maxDim), nil
}

// Parses given image data
func Parse(data []byte) (*image.RGBA, error) {
	defer bench.Begin()()
	buf := bytes.NewBuffer(data)
	// Find data format
	format := http.DetectContentType(data)
	var img image.Image
	var err error
	endBench := bench.Begin()
	switch format {
	case "image/png":
		img, err = png.Decode(buf)
		endBench("png.Decode")
	case "image/jpeg":
		img, err = jpeg.Decode(buf)
		endBench("jpeg.Decode")
	case "image/gif":
		img, err = gif.Decode(buf)
		endBench("gif.Decode")
	case "image/webp":
		// x/image/webp library doesn't support animated webp images so we are using our custom version can decode first frame of the animation
		img, err = webp.Decode(buf)
		endBench("webp.Decode")
	case "image/bmp":
		img, err = bmp.Decode(buf)
		endBench("bmp.Decode")
	default:
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}
	if err != nil {
		return nil, err
	}
	return ConvertToRgba(img), nil
}

func MustParse(data []byte) *image.RGBA {
	img, err := Parse(data)
	if err != nil {
		panic(err)
	}
	return img
}

func ParseAndShrink(data []byte, maxDim int) (*image.RGBA, error) {
	img, err := Parse(data)
	if err != nil {
		return nil, err
	}
	return Shrink(img, maxDim), nil
}

func MustParseAndShrink(data []byte, maxDim int) *image.RGBA {
	img, err := ParseAndShrink(data, maxDim)
	if err != nil {
		panic(err)
	}
	return img
}

// Shrinks image to 'max' while keeping aspect ratio
func Shrink(img image.Image, max int) *image.RGBA {
	defer bench.Begin()()
	if max <= 0 {
		panic("image can not be resized to zero size")
	}
	b := img.Bounds()
	if b.Dx() <= max && b.Dy() <= max { // No need to shrink
		return ConvertToRgba(img)
	}
	return ResizeKeepRatio(img, max)
}

// Resizes image while keeping aspect ratio
// Bigger dimension will be resized to size, and smaller one
// will be calculated depending on the aspect ratio
func ResizeKeepRatio(img image.Image, size int) *image.RGBA {
	defer bench.Begin()()
	b := img.Bounds()
	w := float64(b.Dx())
	h := float64(b.Dy())
	sizef := float64(size)
	var tw, th float64
	switch {
	case w > h:
		tw = sizef
		th = sizef * (h / w)
	case w < h:
		tw = sizef * (w / h)
		th = sizef
	case w == h:
		tw = sizef
		th = sizef
	}
	return Resize(img, image.Pt(int(math.Round(tw)), int(math.Round(th))))
}

func Resize(img image.Image, size image.Point) *image.RGBA {
	defer bench.Begin()()
	b := img.Bounds()
	if b.Dx() == size.X && b.Dy() == size.Y {
		return ConvertToRgba(img)
	}
	newImg := image.NewRGBA(image.Rect(0, 0, size.X, size.Y))
	draw.BiLinear.Scale(newImg, newImg.Bounds(), img, b, draw.Over, nil)
	return newImg
}

func ConvertToRgba(img image.Image) *image.RGBA {
	defer bench.Begin()()
	switch img := img.(type) {
	case *image.RGBA:
		return img
	default:
		b := img.Bounds()
		rgba := image.NewRGBA(b)
		draw.Draw(rgba, b, img, image.Point{}, draw.Over)
		// log.Debug().Println("Image converted to RGBA")
		return rgba
	}
}

func Encode(img image.Image, format string) ([]byte, error) {
	out := new(bytes.Buffer)
	err := EncodeTo(img, format, out)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func EncodeTo(img image.Image, format string, out io.Writer) error {
	var err error
	endBench := bench.Begin()
	switch format {
	case "png":
		err = png.Encode(out, img)
		endBench("png.Encode")
	case "jpg", "jpeg":
		err = jpeg.Encode(out, img, nil)
		endBench("jpeg.Encode")
	case "gif":
		err = gif.Encode(out, img, nil)
		endBench("gif.Encode")
	case "webp":
		err = fmt.Errorf("webp encoding is not supported")
	case "bmp":
		err = bmp.Encode(out, img)
		endBench("bmp.Encode")
	default:
		err = fmt.Errorf("unsupported image format: %s", format)
	}
	return err
}
