package math3d

import (
	"math"
	"testing"
)

func TestColorConstructors(t *testing.T) {
	if got, want := NewRGB(1, 2, 3), (RGB{R: 1, G: 2, B: 3}); got != want {
		t.Fatalf("NewRGB() = %#v, want %#v", got, want)
	}
	if got, want := NewRGBA(1, 2, 3, 4), (RGBA{R: 1, G: 2, B: 3, A: 4}); got != want {
		t.Fatalf("NewRGBA() = %#v, want %#v", got, want)
	}
	if got, want := NewHDRColor(1, 2, 3, 4), (HDRColor{R: 1, G: 2, B: 3, A: 4}); got != want {
		t.Fatalf("NewHDRColor() = %#v, want %#v", got, want)
	}
}

func TestNamedColors(t *testing.T) {
	tests := []struct {
		name string
		got  RGBA
		want RGBA
	}{
		{"LightRed", LightRed(), NewRGBA(255, 128, 128, 255)},
		{"DarkRed", DarkRed(), NewRGBA(255, 0, 0, 255)},
		{"LightGreen", LightGreen(), NewRGBA(128, 255, 128, 255)},
		{"DarkGreen", DarkGreen(), NewRGBA(0, 255, 0, 255)},
		{"LightBlue", LightBlue(), NewRGBA(128, 128, 255, 255)},
		{"DarkBlue", DarkBlue(), NewRGBA(0, 0, 255, 255)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.got != test.want {
				t.Fatalf("got %#v, want %#v", test.got, test.want)
			}
		})
	}
}

func TestComplexInterop(t *testing.T) {
	want := NewComplex(2.5, -3.25)
	if got := ComplexFrom128(complex(2.5, -3.25)); got != want {
		t.Fatalf("ComplexFrom128() = %#v, want %#v", got, want)
	}
	if got := want.Complex128(); got != complex(2.5, -3.25) {
		t.Fatalf("Complex128() = %v", got)
	}
}

func TestFloatingRecordAlmostEqual(t *testing.T) {
	if !NewHDRColor(1, 2, 3, 4).AlmostEqual(NewHDRColor(1.01, 2, 3, 4), 0.02) {
		t.Fatal("near HDR colors should compare approximately equal")
	}
	difference := abs32(float32(1) - float32(1.02))
	if NewHDRColor(1, 2, 3, 4).AlmostEqual(NewHDRColor(1.02, 2, 3, 4), difference) {
		t.Fatal("tolerance boundary must be strict")
	}
	if !NewComplex(1, 2).AlmostEqual(NewComplex(1.01, 2), 0.02) {
		t.Fatal("near complex values should compare approximately equal")
	}
	if NewComplex(math.NaN(), 0).AlmostEqual(NewComplex(math.NaN(), 0), 1) {
		t.Fatal("NaN must not compare approximately equal")
	}
	if NewComplex(math.Inf(1), 0).AlmostEqual(NewComplex(math.Inf(1), 0), 1) {
		t.Fatal("infinity follows upstream subtraction behavior")
	}
}
