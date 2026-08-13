package math3d

import (
	"math"
	"testing"
)

func TestVectorConstructorsAndConversions(t *testing.T) {
	if got, want := V2(1, 2), (Vec2{X: 1, Y: 2}); got != want {
		t.Fatalf("V2() = %#v, want %#v", got, want)
	}
	if got, want := V3(1, 2, 3), (Vec3{X: 1, Y: 2, Z: 3}); got != want {
		t.Fatalf("V3() = %#v, want %#v", got, want)
	}
	if got, want := V4(1, 2, 3, 4), (Vec4{X: 1, Y: 2, Z: 3, W: 4}); got != want {
		t.Fatalf("V4() = %#v, want %#v", got, want)
	}
	if got, want := SplatV3(5), V3(5, 5, 5); got != want {
		t.Fatalf("SplatV3() = %#v, want %#v", got, want)
	}
	if got, want := V3(1.25, -2.5, 3.75).DVec3(), DV3(1.25, -2.5, 3.75); got != want {
		t.Fatalf("DVec3() = %#v, want %#v", got, want)
	}
	if got, want := DV3(1.25, -2.5, 3.75).Vec3(), V3(1.25, -2.5, 3.75); got != want {
		t.Fatalf("Vec3() = %#v, want %#v", got, want)
	}
	if got := DV2(math.MaxFloat64, 0).Vec2(); !got.IsInf() {
		t.Fatalf("narrowing overflow = %#v, want infinity", got)
	}
}

func TestVec3Arithmetic(t *testing.T) {
	a, b := V3(1, 2, 3), V3(4, 5, 6)
	tests := []struct {
		name      string
		got, want Vec3
	}{
		{"Add", a.Add(b), V3(5, 7, 9)},
		{"Sub", a.Sub(b), V3(-3, -3, -3)},
		{"Mul", a.Mul(b), V3(4, 10, 18)},
		{"Div", b.Div(a), V3(4, 2.5, 2)},
		{"Scale", a.Scale(2), V3(2, 4, 6)},
		{"Negated", a.Negated(), V3(-1, -2, -3)},
		{"Min", a.Min(V3(0, 4, 2)), V3(0, 2, 2)},
		{"Max", a.Max(V3(0, 4, 2)), V3(1, 4, 3)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.got != test.want {
				t.Fatalf("got %#v, want %#v", test.got, test.want)
			}
		})
	}
	if got := a.Dot(b); got != 32 {
		t.Fatalf("Dot() = %v, want 32", got)
	}
	if got := V3(3, 4, 0).Magnitude(); got != 5 {
		t.Fatalf("Magnitude() = %v, want 5", got)
	}
	if got := a.SumComponents(); got != 6 {
		t.Fatalf("SumComponents() = %v, want 6", got)
	}
	if got := a.ProductComponents(); got != 6 {
		t.Fatalf("ProductComponents() = %v, want 6", got)
	}
	if a.MinComponent() != 1 || a.MaxComponent() != 3 {
		t.Fatal("component reductions returned incorrect values")
	}
}

func TestVectorCrossProducts(t *testing.T) {
	x, y, z := V3(1, 0, 0), V3(0, 1, 0), V3(0, 0, 1)
	if got := x.Cross(y); got != z {
		t.Fatalf("x.Cross(y) = %#v, want %#v", got, z)
	}
	if got := V2(1, 0).Cross(V2(0, 1)); got != 1 {
		t.Fatalf("2D Cross() = %v, want 1", got)
	}
	dx, dy := DV3(1, 0, 0), DV3(0, 1, 0)
	if got, want := dx.Cross(dy), DV3(0, 0, 1); got != want {
		t.Fatalf("DVec3.Cross() = %#v, want %#v", got, want)
	}
}

func TestVectorNormalization(t *testing.T) {
	if got, ok := V3(3, 4, 0).Normalized(); !ok || !got.AlmostEqual(V3(0.6, 0.8, 0), 1e-6) {
		t.Fatalf("Normalized() = %#v, %v", got, ok)
	}
	if got, ok := (Vec3{}).Normalized(); ok || got != (Vec3{}) {
		t.Fatalf("zero Normalized() = %#v, %v", got, ok)
	}
	if got, ok := V2(math.MaxFloat32, math.MaxFloat32).Normalized(); ok || got != (Vec2{}) {
		t.Fatalf("overflowed Normalized() = %#v, %v", got, ok)
	}
	if got, ok := DV4(1, 2, 2, 4).Normalized(); !ok || math.Abs(got.Magnitude()-1) > 1e-15 {
		t.Fatalf("DVec4.Normalized() magnitude = %.17g, ok %v", got.Magnitude(), ok)
	}
}

func TestVectorFloatingPointPredicatesAndEquality(t *testing.T) {
	finite := V4(1, 2, 3, 4)
	if !finite.IsFinite() || finite.IsNaN() || finite.IsInf() {
		t.Fatal("finite vector predicates are inconsistent")
	}
	if !V4(float32(math.NaN()), 0, 0, 0).IsNaN() {
		t.Fatal("NaN component was not detected")
	}
	if !V4(float32(math.Inf(1)), 0, 0, 0).IsInf() {
		t.Fatal("infinite component was not detected")
	}
	if V2(float32(math.NaN()), 0).AlmostEqual(V2(float32(math.NaN()), 0), 1) {
		t.Fatal("NaN must not compare approximately equal")
	}
	if V2(float32(math.Inf(1)), 0).AlmostEqual(V2(float32(math.Inf(1)), 0), 1) {
		t.Fatal("infinity follows upstream subtraction behavior")
	}
	if V2(1, 2).AlmostEqual(V2(1.25, 2), 0.25) {
		t.Fatal("tolerance boundary must be strict")
	}
}

func TestVectorOperationsDoNotMutateInputs(t *testing.T) {
	a, b := V3(1, 2, 3), V3(4, 5, 6)
	wantA, wantB := a, b
	_, _ = a.Add(b).Normalized()
	_ = a.Cross(b)
	if a != wantA || b != wantB {
		t.Fatalf("inputs changed: %#v, %#v", a, b)
	}
}

func FuzzVec3Properties(f *testing.F) {
	f.Add(float32(1), float32(2), float32(3), float32(4), float32(5), float32(6))
	f.Add(float32(-1), float32(0), float32(1), float32(2), float32(-3), float32(4))
	f.Fuzz(func(t *testing.T, ax, ay, az, bx, by, bz float32) {
		a, b := V3(ax, ay, az), V3(bx, by, bz)
		if !a.IsFinite() || !b.IsFinite() {
			t.Skip()
		}
		if got, want := a.Add(b), b.Add(a); got != want {
			t.Fatalf("addition is not commutative: %#v != %#v", got, want)
		}
		if got, want := a.Dot(b), b.Dot(a); got != want && !(math.IsNaN(float64(got)) && math.IsNaN(float64(want))) {
			t.Fatalf("dot is not symmetric: %v != %v", got, want)
		}
		cross := a.Cross(b)
		if cross.IsFinite() {
			scale := max(float64(1), a.Magnitude()*cross.Magnitude(), b.Magnitude()*cross.Magnitude())
			if float64(abs32(cross.Dot(a))) > 2e-5*scale || float64(abs32(cross.Dot(b))) > 2e-5*scale {
				t.Fatalf("cross product is not orthogonal: a=%#v b=%#v cross=%#v", a, b, cross)
			}
		}
		if normalized, ok := a.Normalized(); ok && normalized.IsFinite() && math.Abs(normalized.Magnitude()-1) > 2e-6 {
			t.Fatalf("normalized magnitude = %v", normalized.Magnitude())
		}
	})
}
