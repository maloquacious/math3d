package math3d

import "math"

// DVec2 is a float64 2D vector. Its zero value is the zero vector.
type DVec2 struct{ X, Y float64 }

func DV2(x, y float64) DVec2                  { return DVec2{X: x, Y: y} }
func SplatDV2(v float64) DVec2                { return DV2(v, v) }
func (v DVec2) Add(o DVec2) DVec2             { return DV2(v.X+o.X, v.Y+o.Y) }
func (v DVec2) Sub(o DVec2) DVec2             { return DV2(v.X-o.X, v.Y-o.Y) }
func (v DVec2) Mul(o DVec2) DVec2             { return DV2(v.X*o.X, v.Y*o.Y) }
func (v DVec2) Div(o DVec2) DVec2             { return DV2(v.X/o.X, v.Y/o.Y) }
func (v DVec2) Scale(s float64) DVec2         { return DV2(v.X*s, v.Y*s) }
func (v DVec2) Negated() DVec2                { return DV2(-v.X, -v.Y) }
func (v DVec2) Dot(o DVec2) float64           { return v.X*o.X + v.Y*o.Y }
func (v DVec2) Cross(o DVec2) float64         { return v.X*o.Y - v.Y*o.X }
func (v DVec2) SumComponents() float64        { return v.X + v.Y }
func (v DVec2) SumSquaredComponents() float64 { return v.Dot(v) }
func (v DVec2) ProductComponents() float64    { return v.X * v.Y }
func (v DVec2) MinComponent() float64         { return min(v.X, v.Y) }
func (v DVec2) MaxComponent() float64         { return max(v.X, v.Y) }
func (v DVec2) MagnitudeSquared() float64     { return v.SumSquaredComponents() }
func (v DVec2) Magnitude() float64            { return math.Sqrt(v.MagnitudeSquared()) }
func (v DVec2) Min(o DVec2) DVec2             { return DV2(min(v.X, o.X), min(v.Y, o.Y)) }
func (v DVec2) Max(o DVec2) DVec2             { return DV2(max(v.X, o.X), max(v.Y, o.Y)) }
func (v DVec2) Normalized() (DVec2, bool) {
	m := v.Magnitude()
	if m == 0 || math.IsNaN(m) || math.IsInf(m, 0) {
		return DVec2{}, false
	}
	return DV2(v.X/m, v.Y/m), true
}
func (v DVec2) IsNaN() bool    { return math.IsNaN(v.X) || math.IsNaN(v.Y) }
func (v DVec2) IsInf() bool    { return math.IsInf(v.X, 0) || math.IsInf(v.Y, 0) }
func (v DVec2) IsFinite() bool { return !v.IsNaN() && !v.IsInf() }

// AlmostEqual compares corresponding components using a strict absolute
// tolerance.
func (v DVec2) AlmostEqual(o DVec2, tolerance float64) bool {
	return math.Abs(v.X-o.X) < tolerance && math.Abs(v.Y-o.Y) < tolerance
}

// Vec2 narrows v to float32 precision. Components can be rounded or overflow
// to an infinity.
func (v DVec2) Vec2() Vec2 { return V2(float32(v.X), float32(v.Y)) }

// DVec3 is a float64 3D vector. Its zero value is the zero vector.
type DVec3 struct{ X, Y, Z float64 }

func DV3(x, y, z float64) DVec3               { return DVec3{X: x, Y: y, Z: z} }
func SplatDV3(v float64) DVec3                { return DV3(v, v, v) }
func (v DVec3) Add(o DVec3) DVec3             { return DV3(v.X+o.X, v.Y+o.Y, v.Z+o.Z) }
func (v DVec3) Sub(o DVec3) DVec3             { return DV3(v.X-o.X, v.Y-o.Y, v.Z-o.Z) }
func (v DVec3) Mul(o DVec3) DVec3             { return DV3(v.X*o.X, v.Y*o.Y, v.Z*o.Z) }
func (v DVec3) Div(o DVec3) DVec3             { return DV3(v.X/o.X, v.Y/o.Y, v.Z/o.Z) }
func (v DVec3) Scale(s float64) DVec3         { return DV3(v.X*s, v.Y*s, v.Z*s) }
func (v DVec3) Negated() DVec3                { return DV3(-v.X, -v.Y, -v.Z) }
func (v DVec3) Dot(o DVec3) float64           { return v.X*o.X + v.Y*o.Y + v.Z*o.Z }
func (v DVec3) Cross(o DVec3) DVec3           { return DV3(v.Y*o.Z-v.Z*o.Y, v.Z*o.X-v.X*o.Z, v.X*o.Y-v.Y*o.X) }
func (v DVec3) SumComponents() float64        { return v.X + v.Y + v.Z }
func (v DVec3) SumSquaredComponents() float64 { return v.Dot(v) }
func (v DVec3) ProductComponents() float64    { return v.X * v.Y * v.Z }
func (v DVec3) MinComponent() float64         { return min(v.X, v.Y, v.Z) }
func (v DVec3) MaxComponent() float64         { return max(v.X, v.Y, v.Z) }
func (v DVec3) MagnitudeSquared() float64     { return v.SumSquaredComponents() }
func (v DVec3) Magnitude() float64            { return math.Sqrt(v.MagnitudeSquared()) }
func (v DVec3) Min(o DVec3) DVec3             { return DV3(min(v.X, o.X), min(v.Y, o.Y), min(v.Z, o.Z)) }
func (v DVec3) Max(o DVec3) DVec3             { return DV3(max(v.X, o.X), max(v.Y, o.Y), max(v.Z, o.Z)) }
func (v DVec3) Normalized() (DVec3, bool) {
	m := v.Magnitude()
	if m == 0 || math.IsNaN(m) || math.IsInf(m, 0) {
		return DVec3{}, false
	}
	return DV3(v.X/m, v.Y/m, v.Z/m), true
}
func (v DVec3) IsNaN() bool    { return math.IsNaN(v.X) || math.IsNaN(v.Y) || math.IsNaN(v.Z) }
func (v DVec3) IsInf() bool    { return math.IsInf(v.X, 0) || math.IsInf(v.Y, 0) || math.IsInf(v.Z, 0) }
func (v DVec3) IsFinite() bool { return !v.IsNaN() && !v.IsInf() }

// AlmostEqual compares corresponding components using a strict absolute
// tolerance.
func (v DVec3) AlmostEqual(o DVec3, tolerance float64) bool {
	return math.Abs(v.X-o.X) < tolerance && math.Abs(v.Y-o.Y) < tolerance && math.Abs(v.Z-o.Z) < tolerance
}

// Vec3 narrows v to float32 precision. Components can be rounded or overflow
// to an infinity.
func (v DVec3) Vec3() Vec3 { return V3(float32(v.X), float32(v.Y), float32(v.Z)) }

// DVec4 is a float64 4D vector. Its zero value is the zero vector.
type DVec4 struct{ X, Y, Z, W float64 }

func DV4(x, y, z, w float64) DVec4            { return DVec4{X: x, Y: y, Z: z, W: w} }
func SplatDV4(v float64) DVec4                { return DV4(v, v, v, v) }
func (v DVec4) Add(o DVec4) DVec4             { return DV4(v.X+o.X, v.Y+o.Y, v.Z+o.Z, v.W+o.W) }
func (v DVec4) Sub(o DVec4) DVec4             { return DV4(v.X-o.X, v.Y-o.Y, v.Z-o.Z, v.W-o.W) }
func (v DVec4) Mul(o DVec4) DVec4             { return DV4(v.X*o.X, v.Y*o.Y, v.Z*o.Z, v.W*o.W) }
func (v DVec4) Div(o DVec4) DVec4             { return DV4(v.X/o.X, v.Y/o.Y, v.Z/o.Z, v.W/o.W) }
func (v DVec4) Scale(s float64) DVec4         { return DV4(v.X*s, v.Y*s, v.Z*s, v.W*s) }
func (v DVec4) Negated() DVec4                { return DV4(-v.X, -v.Y, -v.Z, -v.W) }
func (v DVec4) Dot(o DVec4) float64           { return v.X*o.X + v.Y*o.Y + v.Z*o.Z + v.W*o.W }
func (v DVec4) SumComponents() float64        { return v.X + v.Y + v.Z + v.W }
func (v DVec4) SumSquaredComponents() float64 { return v.Dot(v) }
func (v DVec4) ProductComponents() float64    { return v.X * v.Y * v.Z * v.W }
func (v DVec4) MinComponent() float64         { return min(v.X, v.Y, v.Z, v.W) }
func (v DVec4) MaxComponent() float64         { return max(v.X, v.Y, v.Z, v.W) }
func (v DVec4) MagnitudeSquared() float64     { return v.SumSquaredComponents() }
func (v DVec4) Magnitude() float64            { return math.Sqrt(v.MagnitudeSquared()) }
func (v DVec4) Min(o DVec4) DVec4 {
	return DV4(min(v.X, o.X), min(v.Y, o.Y), min(v.Z, o.Z), min(v.W, o.W))
}
func (v DVec4) Max(o DVec4) DVec4 {
	return DV4(max(v.X, o.X), max(v.Y, o.Y), max(v.Z, o.Z), max(v.W, o.W))
}
func (v DVec4) Normalized() (DVec4, bool) {
	m := v.Magnitude()
	if m == 0 || math.IsNaN(m) || math.IsInf(m, 0) {
		return DVec4{}, false
	}
	return DV4(v.X/m, v.Y/m, v.Z/m, v.W/m), true
}
func (v DVec4) IsNaN() bool {
	return math.IsNaN(v.X) || math.IsNaN(v.Y) || math.IsNaN(v.Z) || math.IsNaN(v.W)
}
func (v DVec4) IsInf() bool {
	return math.IsInf(v.X, 0) || math.IsInf(v.Y, 0) || math.IsInf(v.Z, 0) || math.IsInf(v.W, 0)
}
func (v DVec4) IsFinite() bool { return !v.IsNaN() && !v.IsInf() }

// AlmostEqual compares corresponding components using a strict absolute
// tolerance.
func (v DVec4) AlmostEqual(o DVec4, tolerance float64) bool {
	return math.Abs(v.X-o.X) < tolerance && math.Abs(v.Y-o.Y) < tolerance && math.Abs(v.Z-o.Z) < tolerance && math.Abs(v.W-o.W) < tolerance
}

// Vec4 narrows v to float32 precision. Components can be rounded or overflow
// to an infinity.
func (v DVec4) Vec4() Vec4 { return V4(float32(v.X), float32(v.Y), float32(v.Z), float32(v.W)) }
