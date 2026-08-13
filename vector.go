package math3d

import "math"

// Vec2 is the upstream Vector2 value. Its components use source-compatible
// float32 precision.
type Vec2 struct {
	X, Y float32
}

// V2 constructs a Vec2.
func V2(x, y float32) Vec2 { return Vec2{X: x, Y: y} }

// SplatV2 constructs a Vec2 whose components are value.
func SplatV2(value float32) Vec2 { return V2(value, value) }

func (v Vec2) Add(other Vec2) Vec2 { return V2(v.X+other.X, v.Y+other.Y) }
func (v Vec2) Sub(other Vec2) Vec2 { return V2(v.X-other.X, v.Y-other.Y) }

// Mul multiplies corresponding components.
func (v Vec2) Mul(other Vec2) Vec2 { return V2(v.X*other.X, v.Y*other.Y) }

// Div divides corresponding components using normal IEEE 754 semantics.
func (v Vec2) Div(other Vec2) Vec2       { return V2(v.X/other.X, v.Y/other.Y) }
func (v Vec2) Scale(scalar float32) Vec2 { return V2(v.X*scalar, v.Y*scalar) }
func (v Vec2) Negated() Vec2             { return V2(-v.X, -v.Y) }
func (v Vec2) Dot(other Vec2) float32    { return v.X*other.X + v.Y*other.Y }

// Cross returns the signed scalar cross product of v and other.
func (v Vec2) Cross(other Vec2) float32      { return v.X*other.Y - v.Y*other.X }
func (v Vec2) SumComponents() float32        { return v.X + v.Y }
func (v Vec2) SumSquaredComponents() float32 { return v.Dot(v) }
func (v Vec2) ProductComponents() float32    { return v.X * v.Y }
func (v Vec2) MinComponent() float32         { return min(v.X, v.Y) }
func (v Vec2) MaxComponent() float32         { return max(v.X, v.Y) }
func (v Vec2) MagnitudeSquared() float64     { return float64(v.SumSquaredComponents()) }
func (v Vec2) Magnitude() float64            { return math.Sqrt(v.MagnitudeSquared()) }
func (v Vec2) Min(other Vec2) Vec2           { return V2(min(v.X, other.X), min(v.Y, other.Y)) }
func (v Vec2) Max(other Vec2) Vec2           { return V2(max(v.X, other.X), max(v.Y, other.Y)) }

// Normalized returns a unit vector. It fails when the magnitude is zero or
// non-finite, including when its float32 calculation overflows.
func (v Vec2) Normalized() (Vec2, bool) {
	magnitude := v.Magnitude()
	if magnitude == 0 || math.IsNaN(float64(magnitude)) || math.IsInf(float64(magnitude), 0) {
		return Vec2{}, false
	}
	return V2(v.X/float32(magnitude), v.Y/float32(magnitude)), true
}

func (v Vec2) IsNaN() bool    { return math.IsNaN(float64(v.X)) || math.IsNaN(float64(v.Y)) }
func (v Vec2) IsInf() bool    { return math.IsInf(float64(v.X), 0) || math.IsInf(float64(v.Y), 0) }
func (v Vec2) IsFinite() bool { return !v.IsNaN() && !v.IsInf() }
func (v Vec2) AlmostEqual(other Vec2, tolerance float32) bool {
	return abs32(v.X-other.X) < tolerance && abs32(v.Y-other.Y) < tolerance
}

// DVec2 converts v to float64 precision.
func (v Vec2) DVec2() DVec2 { return DV2(float64(v.X), float64(v.Y)) }

// Vec3 promotes v with a zero Z component.
func (v Vec2) Vec3() Vec3 { return V3(v.X, v.Y, 0) }

// Vec4 promotes v with zero Z and W components.
func (v Vec2) Vec4() Vec4 { return V4(v.X, v.Y, 0, 0) }

// Vec3 is the upstream Vector3 value. Its components use source-compatible
// float32 precision.
type Vec3 struct {
	X, Y, Z float32
}

// V3 constructs a Vec3.
func V3(x, y, z float32) Vec3 { return Vec3{X: x, Y: y, Z: z} }

// SplatV3 constructs a Vec3 whose components are value.
func SplatV3(value float32) Vec3 { return V3(value, value, value) }

func (v Vec3) Add(other Vec3) Vec3       { return V3(v.X+other.X, v.Y+other.Y, v.Z+other.Z) }
func (v Vec3) Sub(other Vec3) Vec3       { return V3(v.X-other.X, v.Y-other.Y, v.Z-other.Z) }
func (v Vec3) Mul(other Vec3) Vec3       { return V3(v.X*other.X, v.Y*other.Y, v.Z*other.Z) }
func (v Vec3) Div(other Vec3) Vec3       { return V3(v.X/other.X, v.Y/other.Y, v.Z/other.Z) }
func (v Vec3) Scale(scalar float32) Vec3 { return V3(v.X*scalar, v.Y*scalar, v.Z*scalar) }
func (v Vec3) Negated() Vec3             { return V3(-v.X, -v.Y, -v.Z) }
func (v Vec3) Dot(other Vec3) float32    { return v.X*other.X + v.Y*other.Y + v.Z*other.Z }
func (v Vec3) Cross(other Vec3) Vec3 {
	return V3(v.Y*other.Z-v.Z*other.Y, v.Z*other.X-v.X*other.Z, v.X*other.Y-v.Y*other.X)
}
func (v Vec3) SumComponents() float32        { return v.X + v.Y + v.Z }
func (v Vec3) SumSquaredComponents() float32 { return v.Dot(v) }
func (v Vec3) ProductComponents() float32    { return v.X * v.Y * v.Z }
func (v Vec3) MinComponent() float32         { return min(v.X, v.Y, v.Z) }
func (v Vec3) MaxComponent() float32         { return max(v.X, v.Y, v.Z) }
func (v Vec3) MagnitudeSquared() float64     { return float64(v.SumSquaredComponents()) }
func (v Vec3) Magnitude() float64            { return math.Sqrt(v.MagnitudeSquared()) }
func (v Vec3) Min(other Vec3) Vec3 {
	return V3(min(v.X, other.X), min(v.Y, other.Y), min(v.Z, other.Z))
}
func (v Vec3) Max(other Vec3) Vec3 {
	return V3(max(v.X, other.X), max(v.Y, other.Y), max(v.Z, other.Z))
}
func (v Vec3) Normalized() (Vec3, bool) {
	magnitude := v.Magnitude()
	if magnitude == 0 || math.IsNaN(float64(magnitude)) || math.IsInf(float64(magnitude), 0) {
		return Vec3{}, false
	}
	scale := float32(magnitude)
	return V3(v.X/scale, v.Y/scale, v.Z/scale), true
}
func (v Vec3) IsNaN() bool {
	return math.IsNaN(float64(v.X)) || math.IsNaN(float64(v.Y)) || math.IsNaN(float64(v.Z))
}
func (v Vec3) IsInf() bool {
	return math.IsInf(float64(v.X), 0) || math.IsInf(float64(v.Y), 0) || math.IsInf(float64(v.Z), 0)
}
func (v Vec3) IsFinite() bool { return !v.IsNaN() && !v.IsInf() }
func (v Vec3) AlmostEqual(other Vec3, tolerance float32) bool {
	return abs32(v.X-other.X) < tolerance && abs32(v.Y-other.Y) < tolerance && abs32(v.Z-other.Z) < tolerance
}
func (v Vec3) DVec3() DVec3 { return DV3(float64(v.X), float64(v.Y), float64(v.Z)) }

func (v Vec3) XY() Vec2  { return V2(v.X, v.Y) }
func (v Vec3) XZ() Vec2  { return V2(v.X, v.Z) }
func (v Vec3) YZ() Vec2  { return V2(v.Y, v.Z) }
func (v Vec3) XZY() Vec3 { return V3(v.X, v.Z, v.Y) }
func (v Vec3) ZXY() Vec3 { return V3(v.Z, v.X, v.Y) }
func (v Vec3) ZYX() Vec3 { return V3(v.Z, v.Y, v.X) }
func (v Vec3) YXZ() Vec3 { return V3(v.Y, v.X, v.Z) }
func (v Vec3) YZX() Vec3 { return V3(v.Y, v.Z, v.X) }

// Vec4 promotes v with a zero W component.
func (v Vec3) Vec4() Vec4 { return V4(v.X, v.Y, v.Z, 0) }

// Vec4 is the upstream Vector4 value. Its components use source-compatible
// float32 precision.
type Vec4 struct {
	X, Y, Z, W float32
}

// V4 constructs a Vec4.
func V4(x, y, z, w float32) Vec4 { return Vec4{X: x, Y: y, Z: z, W: w} }

// SplatV4 constructs a Vec4 whose components are value.
func SplatV4(value float32) Vec4 { return V4(value, value, value, value) }

func (v Vec4) Add(o Vec4) Vec4               { return V4(v.X+o.X, v.Y+o.Y, v.Z+o.Z, v.W+o.W) }
func (v Vec4) Sub(o Vec4) Vec4               { return V4(v.X-o.X, v.Y-o.Y, v.Z-o.Z, v.W-o.W) }
func (v Vec4) Mul(o Vec4) Vec4               { return V4(v.X*o.X, v.Y*o.Y, v.Z*o.Z, v.W*o.W) }
func (v Vec4) Div(o Vec4) Vec4               { return V4(v.X/o.X, v.Y/o.Y, v.Z/o.Z, v.W/o.W) }
func (v Vec4) Scale(s float32) Vec4          { return V4(v.X*s, v.Y*s, v.Z*s, v.W*s) }
func (v Vec4) Negated() Vec4                 { return V4(-v.X, -v.Y, -v.Z, -v.W) }
func (v Vec4) Dot(o Vec4) float32            { return v.X*o.X + v.Y*o.Y + v.Z*o.Z + v.W*o.W }
func (v Vec4) SumComponents() float32        { return v.X + v.Y + v.Z + v.W }
func (v Vec4) SumSquaredComponents() float32 { return v.Dot(v) }
func (v Vec4) ProductComponents() float32    { return v.X * v.Y * v.Z * v.W }
func (v Vec4) MinComponent() float32         { return min(v.X, v.Y, v.Z, v.W) }
func (v Vec4) MaxComponent() float32         { return max(v.X, v.Y, v.Z, v.W) }
func (v Vec4) MagnitudeSquared() float64     { return float64(v.SumSquaredComponents()) }
func (v Vec4) Magnitude() float64            { return math.Sqrt(v.MagnitudeSquared()) }
func (v Vec4) Min(o Vec4) Vec4               { return V4(min(v.X, o.X), min(v.Y, o.Y), min(v.Z, o.Z), min(v.W, o.W)) }
func (v Vec4) Max(o Vec4) Vec4               { return V4(max(v.X, o.X), max(v.Y, o.Y), max(v.Z, o.Z), max(v.W, o.W)) }
func (v Vec4) Normalized() (Vec4, bool) {
	magnitude := v.Magnitude()
	if magnitude == 0 || math.IsNaN(float64(magnitude)) || math.IsInf(float64(magnitude), 0) {
		return Vec4{}, false
	}
	scale := float32(magnitude)
	return V4(v.X/scale, v.Y/scale, v.Z/scale, v.W/scale), true
}
func (v Vec4) IsNaN() bool {
	return math.IsNaN(float64(v.X)) || math.IsNaN(float64(v.Y)) || math.IsNaN(float64(v.Z)) || math.IsNaN(float64(v.W))
}
func (v Vec4) IsInf() bool {
	return math.IsInf(float64(v.X), 0) || math.IsInf(float64(v.Y), 0) || math.IsInf(float64(v.Z), 0) || math.IsInf(float64(v.W), 0)
}
func (v Vec4) IsFinite() bool { return !v.IsNaN() && !v.IsInf() }
func (v Vec4) AlmostEqual(o Vec4, tolerance float32) bool {
	return abs32(v.X-o.X) < tolerance && abs32(v.Y-o.Y) < tolerance && abs32(v.Z-o.Z) < tolerance && abs32(v.W-o.W) < tolerance
}
func (v Vec4) DVec4() DVec4 { return DV4(float64(v.X), float64(v.Y), float64(v.Z), float64(v.W)) }

// Vec2 returns the X and Y components.
func (v Vec4) Vec2() Vec2 { return V2(v.X, v.Y) }

// Vec3 returns the X, Y, and Z components.
func (v Vec4) Vec3() Vec3 { return V3(v.X, v.Y, v.Z) }
