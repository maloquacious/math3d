package math3d

import "math"

// RGB is the upstream ColorRGB value.
type RGB struct {
	R, G, B byte
}

// NewRGB constructs an RGB color.
func NewRGB(r, g, b byte) RGB { return RGB{R: r, G: g, B: b} }

// RGBA is the upstream ColorRGBA value.
type RGBA struct {
	R, G, B, A byte
}

// NewRGBA constructs an RGBA color.
func NewRGBA(r, g, b, a byte) RGBA { return RGBA{R: r, G: g, B: b, A: a} }

// LightRed returns the upstream light-red color.
func LightRed() RGBA { return NewRGBA(255, 128, 128, 255) }

// DarkRed returns the upstream dark-red color.
func DarkRed() RGBA { return NewRGBA(255, 0, 0, 255) }

// LightGreen returns the upstream light-green color.
func LightGreen() RGBA { return NewRGBA(128, 255, 128, 255) }

// DarkGreen returns the upstream dark-green color.
func DarkGreen() RGBA { return NewRGBA(0, 255, 0, 255) }

// LightBlue returns the upstream light-blue color.
func LightBlue() RGBA { return NewRGBA(128, 128, 255, 255) }

// DarkBlue returns the upstream dark-blue color.
func DarkBlue() RGBA { return NewRGBA(0, 0, 255, 255) }

// HDRColor is the upstream ColorHDR value.
type HDRColor struct {
	R, G, B, A float32
}

// NewHDRColor constructs an HDRColor.
func NewHDRColor(r, g, b, a float32) HDRColor {
	return HDRColor{R: r, G: g, B: b, A: a}
}

// AlmostEqual reports whether corresponding components differ by less than
// tolerance. It follows normal IEEE behavior: NaN never compares equal.
func (c HDRColor) AlmostEqual(other HDRColor, tolerance float32) bool {
	return abs32(c.R-other.R) < tolerance &&
		abs32(c.G-other.G) < tolerance &&
		abs32(c.B-other.B) < tolerance &&
		abs32(c.A-other.A) < tolerance
}

// Complex is the upstream Complex value. It interoperates explicitly with
// Go's complex128 rather than shadowing the built-in conversion syntax.
type Complex struct {
	Real, Imaginary float64
}

// NewComplex constructs a Complex value.
func NewComplex(real, imaginary float64) Complex {
	return Complex{Real: real, Imaginary: imaginary}
}

// ComplexFrom128 converts a built-in complex value.
func ComplexFrom128(value complex128) Complex {
	return NewComplex(real(value), imag(value))
}

// Complex128 converts c to a built-in complex value.
func (c Complex) Complex128() complex128 { return complex(c.Real, c.Imaginary) }

// AlmostEqual reports whether corresponding components differ by less than
// tolerance. It follows normal IEEE behavior: NaN never compares equal.
func (c Complex) AlmostEqual(other Complex, tolerance float64) bool {
	return math.Abs(c.Real-other.Real) < tolerance &&
		math.Abs(c.Imaginary-other.Imaginary) < tolerance
}

func abs32(value float32) float32 {
	if value < 0 {
		return -value
	}
	return value
}
