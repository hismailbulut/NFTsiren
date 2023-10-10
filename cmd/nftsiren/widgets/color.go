package widgets

import (
	"image/color"
)

func AddRGB(c color.NRGBA, a uint8) color.NRGBA {
	c.R = clamp(c.R+a, 0, 0xff)
	c.G = clamp(c.G+a, 0, 0xff)
	c.B = clamp(c.B+a, 0, 0xff)
	return c
}

func SubRGB(c color.NRGBA, a uint8) color.NRGBA {
	c.R = clamp(c.R-a, 0, 0xff)
	c.G = clamp(c.G-a, 0, 0xff)
	c.B = clamp(c.B-a, 0, 0xff)
	return c
}

func MulRGB(c color.NRGBA, a float32) color.NRGBA {
	c.R -= uint8(clamp(float32(c.R)*a, 0, 0xff))
	c.G -= uint8(clamp(float32(c.G)*a, 0, 0xff))
	c.B -= uint8(clamp(float32(c.B)*a, 0, 0xff))
	return c
}

func MulAlpha(c color.NRGBA, a float32) color.NRGBA {
	c.A = uint8(clamp(float32(c.A)*a, 0, 0xff))
	return c
}

func WithAlpha(c color.NRGBA, a uint8) color.NRGBA {
	c.A = a
	return c
}

func Lighter(c color.NRGBA) color.NRGBA {
	return AddRGB(c, 0x08)
}

func Darker(c color.NRGBA) color.NRGBA {
	// c.R = clamp(c.R-c.R/0x10, 0, 0xff)
	// c.G = clamp(c.G-c.G/0x10, 0, 0xff)
	// c.B = clamp(c.B-c.B/0x10, 0, 0xff)
	// return c
	return SubRGB(c, 0x08)
}