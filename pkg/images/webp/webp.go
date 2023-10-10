// Pure go webp decoder
// This is a custom version of x/image/webp
// This version can decode first frame of an animation
// https://developers.google.com/speed/webp/docs/riff_container
package webp

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"io"

	"golang.org/x/image/riff"
	"golang.org/x/image/vp8"
	"golang.org/x/image/vp8l"
)

var errInvalidFormat = errors.New("webp: invalid format")

var (
	fccWEBP = riff.FourCC{'W', 'E', 'B', 'P'}
	fccVP8X = riff.FourCC{'V', 'P', '8', 'X'}
	fccALPH = riff.FourCC{'A', 'L', 'P', 'H'}
	fccVP8  = riff.FourCC{'V', 'P', '8', ' '}
	fccVP8L = riff.FourCC{'V', 'P', '8', 'L'}
	fccANIM = riff.FourCC{'A', 'N', 'I', 'M'}
	fccANMF = riff.FourCC{'A', 'N', 'M', 'F'}
)

func decode(r io.Reader, configOnly bool) (image.Image, *image.Config, error) {
	formType, riffReader, err := riff.NewReader(r)
	if err != nil {
		return nil, nil, err
	}
	if formType != fccWEBP {
		return nil, nil, errInvalidFormat
	}
	dec := decoder{}
	return dec.decode(riffReader, configOnly)
}

type decoder struct {
	alpha       []byte
	alphaStride int
	// wantAlpha      bool
	widthMinusOne  uint32
	heightMinusOne uint32
	buf            [16]byte
}

func (dec *decoder) decode(riffReader *riff.Reader, configOnly bool) (image.Image, *image.Config, error) {
	for {
		chunkID, chunkLen, chunkData, err := riffReader.Next()
		if err == io.EOF {
			err = errInvalidFormat
		}
		if err != nil {
			return nil, nil, err
		}
		img, config, err := dec.decodeChunk(chunkID, chunkLen, chunkData, configOnly)
		if err != nil {
			return nil, nil, err
		}
		if config != nil {
			return nil, config, nil
		}
		if img != nil {
			return img, nil, nil
		}
	}
}

func (dec *decoder) decodeChunk(chunkID riff.FourCC, chunkLen uint32, chunkData io.Reader, configOnly bool) (image.Image, *image.Config, error) {
	switch chunkID {
	case fccVP8X:
		if chunkLen != 10 {
			return nil, nil, errInvalidFormat
		}
		if _, err := io.ReadFull(chunkData, dec.buf[:10]); err != nil {
			return nil, nil, err
		}
		// const (
		// 	animationBit    = 1 << 1
		// 	xmpMetadataBit  = 1 << 2
		// 	exifMetadataBit = 1 << 3
		// 	alphaBit        = 1 << 4
		// 	iccProfileBit   = 1 << 5
		// )
		// dec.wantAlpha = (dec.buf[0] & alphaBit) != 0
		dec.widthMinusOne = uint32(dec.buf[4]) | uint32(dec.buf[5])<<8 | uint32(dec.buf[6])<<16
		dec.heightMinusOne = uint32(dec.buf[7]) | uint32(dec.buf[8])<<8 | uint32(dec.buf[9])<<16

	case fccALPH:
		// if !dec.wantAlpha {
		// 	return nil, nil, errInvalidFormat
		// }
		// dec.wantAlpha = false
		// Read the Pre-processing | Filter | Compression byte.
		if _, err := io.ReadFull(chunkData, dec.buf[:1]); err != nil {
			if err == io.EOF {
				err = errInvalidFormat
			}
			return nil, nil, err
		}
		var err error
		dec.alpha, dec.alphaStride, err = readAlpha(chunkData, dec.widthMinusOne, dec.heightMinusOne, dec.buf[0]&0x03)
		if err != nil {
			return nil, nil, err
		}
		unfilterAlpha(dec.alpha, dec.alphaStride, (dec.buf[0]>>2)&0x03)

	case fccVP8:
		// if dec.wantAlpha || int32(chunkLen) < 0 {
		// 	return nil, nil, errInvalidFormat
		// }
		d := vp8.NewDecoder()
		d.Init(chunkData, int(chunkLen))
		fh, err := d.DecodeFrameHeader()
		if err != nil {
			return nil, nil, err
		}
		if configOnly {
			if dec.alpha != nil {
				config := &image.Config{
					ColorModel: color.NYCbCrAModel,
					Width:      fh.Width,
					Height:     fh.Height,
				}
				return nil, config, nil
			}
			config := &image.Config{
				ColorModel: color.YCbCrModel,
				Width:      fh.Width,
				Height:     fh.Height,
			}
			return nil, config, nil
		}
		m, err := d.DecodeFrame()
		if err != nil {
			return nil, nil, err
		}
		if dec.alpha != nil {
			img := &image.NYCbCrA{
				YCbCr:   *m,
				A:       dec.alpha,
				AStride: dec.alphaStride,
			}
			return img, nil, nil
		}
		return m, nil, nil

	case fccVP8L:
		// if dec.wantAlpha || dec.alpha != nil {
		// 	return nil, nil, errInvalidFormat
		// }
		if configOnly {
			c, err := vp8l.DecodeConfig(chunkData)
			return nil, &c, err
		}
		m, err := vp8l.Decode(chunkData)
		return m, nil, err

	case fccANIM:
		// Currently we don't need this, it only contains
		// background color of all frames (can be ignored?)
		// and loop count

	case fccANMF:
		if _, err := io.ReadFull(chunkData, dec.buf[:16]); err != nil {
			return nil, nil, err
		}
		// frameX := uint32(dec.buf[0]) | uint32(dec.buf[1])<<8 | uint32(dec.buf[2])<<16
		// frameY := uint32(dec.buf[3]) | uint32(dec.buf[4])<<8 | uint32(dec.buf[5])<<16
		// frameWidthMinusOne := uint32(dec.buf[6]) | uint32(dec.buf[7])<<8 | uint32(dec.buf[8])<<16
		// frameHeightMinusOne := uint32(dec.buf[9]) | uint32(dec.buf[10])<<8 | uint32(dec.buf[11])<<16
		// frameDuration := uint32(dec.buf[12]) | uint32(dec.buf[13])<<8 | uint32(dec.buf[14])<<16
		// blendingMethod := (uint8(dec.buf[15]) & 1 << 1) >> 1
		// disposalMethod := uint8(dec.buf[15]) & 1
		dummyListType := bytes.NewReader([]byte{0, 0, 0, 0})
		_, riffReader, err := riff.NewListReader(chunkLen-16+4, io.MultiReader(dummyListType, chunkData))
		if err != nil {
			return nil, nil, err
		}
		return dec.decode(riffReader, configOnly)
	}

	return nil, nil, nil
}

func readAlpha(chunkData io.Reader, widthMinusOne, heightMinusOne uint32, compression byte) (
	alpha []byte, alphaStride int, err error) {

	switch compression {
	case 0:
		w := int(widthMinusOne) + 1
		h := int(heightMinusOne) + 1
		alpha = make([]byte, w*h)
		if _, err := io.ReadFull(chunkData, alpha); err != nil {
			return nil, 0, err
		}
		return alpha, w, nil

	case 1:
		// Read the VP8L-compressed alpha values. First, synthesize a 5-byte VP8L header:
		// a 1-byte magic number, a 14-bit widthMinusOne, a 14-bit heightMinusOne,
		// a 1-bit (ignored, zero) alphaIsUsed and a 3-bit (zero) version.
		// TODO(nigeltao): be more efficient than decoding an *image.NRGBA just to
		// extract the green values to a separately allocated []byte. Fixing this
		// will require changes to the vp8l package's API.
		if widthMinusOne > 0x3fff || heightMinusOne > 0x3fff {
			return nil, 0, errors.New("webp: invalid format")
		}
		alphaImage, err := vp8l.Decode(io.MultiReader(
			bytes.NewReader([]byte{
				0x2f, // VP8L magic number.
				uint8(widthMinusOne),
				uint8(widthMinusOne>>8) | uint8(heightMinusOne<<6),
				uint8(heightMinusOne >> 2),
				uint8(heightMinusOne >> 10),
			}),
			chunkData,
		))
		if err != nil {
			return nil, 0, err
		}
		// The green values of the inner NRGBA image are the alpha values of the
		// outer NYCbCrA image.
		pix := alphaImage.(*image.NRGBA).Pix
		alpha = make([]byte, len(pix)/4)
		for i := range alpha {
			alpha[i] = pix[4*i+1]
		}
		return alpha, int(widthMinusOne) + 1, nil
	}
	return nil, 0, errInvalidFormat
}

func unfilterAlpha(alpha []byte, alphaStride int, filter byte) {
	if len(alpha) == 0 || alphaStride == 0 {
		return
	}
	switch filter {
	case 1: // Horizontal filter.
		for i := 1; i < alphaStride; i++ {
			alpha[i] += alpha[i-1]
		}
		for i := alphaStride; i < len(alpha); i += alphaStride {
			// The first column is equivalent to the vertical filter.
			alpha[i] += alpha[i-alphaStride]

			for j := 1; j < alphaStride; j++ {
				alpha[i+j] += alpha[i+j-1]
			}
		}

	case 2: // Vertical filter.
		// The first row is equivalent to the horizontal filter.
		for i := 1; i < alphaStride; i++ {
			alpha[i] += alpha[i-1]
		}

		for i := alphaStride; i < len(alpha); i++ {
			alpha[i] += alpha[i-alphaStride]
		}

	case 3: // Gradient filter.
		// The first row is equivalent to the horizontal filter.
		for i := 1; i < alphaStride; i++ {
			alpha[i] += alpha[i-1]
		}

		for i := alphaStride; i < len(alpha); i += alphaStride {
			// The first column is equivalent to the vertical filter.
			alpha[i] += alpha[i-alphaStride]

			// The interior is predicted on the three top/left pixels.
			for j := 1; j < alphaStride; j++ {
				c := int(alpha[i+j-alphaStride-1])
				b := int(alpha[i+j-alphaStride])
				a := int(alpha[i+j-1])
				x := a + b - c
				if x < 0 {
					x = 0
				} else if x > 255 {
					x = 255
				}
				alpha[i+j] += uint8(x)
			}
		}
	}
}

// Decode reads a WEBP image from r and returns it as an image.Image.
func Decode(r io.Reader) (image.Image, error) {
	m, _, err := decode(r, false)
	if err != nil {
		return nil, err
	}
	return m, err
}

// DecodeConfig returns the color model and dimensions of a WEBP image without
// decoding the entire image.
func DecodeConfig(r io.Reader) (image.Config, error) {
	_, c, err := decode(r, true)
	if c != nil {
		return *c, err
	}
	return image.Config{}, err
}

func init() {
	image.RegisterFormat("webp", "RIFF????WEBPVP8", Decode, DecodeConfig)
}
