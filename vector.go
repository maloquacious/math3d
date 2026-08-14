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

// Reflect returns v reflected around normal. Normal need not be unit length.
// It fails when v or normal is non-finite, normal has zero or overflowed
// squared magnitude, or the result is non-finite.
func (v Vec2) Reflect(normal Vec2) (Vec2, bool) {
	lengthSquared := normal.Dot(normal)
	if !v.IsFinite() || !normal.IsFinite() || lengthSquared == 0 || !finite32(lengthSquared) {
		return Vec2{}, false
	}
	result := v.Sub(normal.Scale(2 * v.Dot(normal) / lengthSquared))
	if !result.IsFinite() {
		return Vec2{}, false
	}
	return result, true
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

// CatmullRom interpolates between value2 and value3 using v and value4 as the
// neighboring control values. Amount is not clamped.
func (v Vec3) CatmullRom(value2, value3, value4 Vec3, amount float32) Vec3 {
	return V3(
		CatmullRom(v.X, value2.X, value3.X, value4.X, amount),
		CatmullRom(v.Y, value2.Y, value3.Y, value4.Y, amount),
		CatmullRom(v.Z, value2.Z, value3.Z, value4.Z, amount),
	)
}

// Hermite interpolates from v with tangent1 to value2 with tangent2. Amount is
// not clamped.
func (v Vec3) Hermite(tangent1, value2, tangent2 Vec3, amount float32) Vec3 {
	return V3(
		Hermite(v.X, tangent1.X, value2.X, tangent2.X, amount),
		Hermite(v.Y, tangent1.Y, value2.Y, tangent2.Y, amount),
		Hermite(v.Z, tangent1.Z, value2.Z, tangent2.Z, amount),
	)
}

// SmoothStep smoothly interpolates from v to value2, clamping amount to [0, 1].
func (v Vec3) SmoothStep(value2 Vec3, amount float32) Vec3 {
	return V3(
		SmoothStep(v.X, value2.X, amount),
		SmoothStep(v.Y, value2.Y, amount),
		SmoothStep(v.Z, value2.Z, amount),
	)
}

// Projection returns the projection of v onto other. It fails when either
// vector is non-finite, other has zero or overflowed squared magnitude, or the
// result is non-finite.
func (v Vec3) Projection(other Vec3) (Vec3, bool) {
	lengthSquared := other.Dot(other)
	if !v.IsFinite() || !other.IsFinite() || lengthSquared == 0 || !finite32(lengthSquared) {
		return Vec3{}, false
	}
	result := other.Scale(v.Dot(other) / lengthSquared)
	if !result.IsFinite() {
		return Vec3{}, false
	}
	return result, true
}

// Rejection returns the component of v perpendicular to other. It has the
// same failure conditions as Projection.
func (v Vec3) Rejection(other Vec3) (Vec3, bool) {
	projection, ok := v.Projection(other)
	if !ok {
		return Vec3{}, false
	}
	result := v.Sub(projection)
	if !result.IsFinite() {
		return Vec3{}, false
	}
	return result, true
}

// Reflect returns v reflected around normal. Normal need not be unit length.
// It has the same failure conditions as Projection.
func (v Vec3) Reflect(normal Vec3) (Vec3, bool) {
	projection, ok := v.Projection(normal)
	if !ok {
		return Vec3{}, false
	}
	result := v.Sub(projection.Scale(2))
	if !result.IsFinite() {
		return Vec3{}, false
	}
	return result, true
}

// IsPerpendicular reports whether both vectors are finite and nonzero and the
// absolute value of their dot product is strictly less than tolerance. Like
// upstream, tolerance is absolute and therefore depends on vector scale.
func (v Vec3) IsPerpendicular(other Vec3, tolerance float32) bool {
	if !v.IsFinite() || !other.IsFinite() || !finite32(tolerance) || tolerance <= 0 {
		return false
	}
	vLengthSquared, otherLengthSquared := v.Dot(v), other.Dot(other)
	return vLengthSquared > 0 && finite32(vLengthSquared) &&
		otherLengthSquared > 0 && finite32(otherLengthSquared) &&
		abs32(v.Dot(other)) < tolerance
}

// Angle returns the smaller angle between v and other in radians. It fails
// when either vector is zero or non-finite, an intermediate overflows, or no
// finite angle can be calculated.
func (v Vec3) Angle(other Vec3) (float32, bool) {
	if !v.IsFinite() || !other.IsFinite() {
		return 0, false
	}
	lengthProduct := v.Magnitude() * other.Magnitude()
	if lengthProduct == 0 || math.IsNaN(lengthProduct) || math.IsInf(lengthProduct, 0) {
		return 0, false
	}
	cosine := float64(v.Dot(other)) / lengthProduct
	if math.IsNaN(cosine) || math.IsInf(cosine, 0) {
		return 0, false
	}
	angle := float32(math.Acos(max(-1, min(cosine, 1))))
	if !finite32(angle) {
		return 0, false
	}
	return angle, true
}

// SignedAngle returns the smaller angle from v to other in radians, with its
// sign determined by axis·(v×other). It fails when Angle fails, axis is zero or
// non-finite, or a nonzero angle has no orientation around axis.
func (v Vec3) SignedAngle(other, axis Vec3) (float32, bool) {
	axisLengthSquared := axis.Dot(axis)
	if !axis.IsFinite() || axisLengthSquared == 0 || !finite32(axisLengthSquared) {
		return 0, false
	}
	angle, ok := v.Angle(other)
	if !ok || angle == 0 {
		return angle, ok
	}
	orientation := axis.Dot(v.Cross(other))
	if !finite32(orientation) || orientation == 0 {
		return 0, false
	}
	if orientation < 0 {
		angle = -angle
	}
	return angle, true
}

// IsCollinear reports whether the smaller angle between the vectors, ignoring
// direction, is at most tolerance radians. Unlike upstream's defective
// signed-angle test, it works in three dimensions and treats opposite vectors
// as collinear. Zero, non-finite, and overflowed vectors return false.
func (v Vec3) IsCollinear(other Vec3, tolerance float32) bool {
	if !finite32(tolerance) || tolerance < 0 {
		return false
	}
	if !v.IsFinite() || !other.IsFinite() {
		return false
	}
	lengthProduct := v.Magnitude() * other.Magnitude()
	if lengthProduct == 0 || math.IsNaN(lengthProduct) || math.IsInf(lengthProduct, 0) {
		return false
	}
	cosine := math.Abs(float64(v.Dot(other)) / lengthProduct)
	if math.IsNaN(cosine) || math.IsInf(cosine, 0) {
		return false
	}
	return float32(math.Acos(min(cosine, 1))) <= tolerance
}

// IsBackFace reports whether normal faces away from lineOfSight, represented
// by a strict negative dot product. The vectors need not be normalized; zero
// and non-finite inputs return false.
func (v Vec3) IsBackFace(lineOfSight Vec3) bool {
	vLengthSquared, sightLengthSquared := v.Dot(v), lineOfSight.Dot(lineOfSight)
	dot := v.Dot(lineOfSight)
	if !v.IsFinite() || !lineOfSight.IsFinite() ||
		vLengthSquared == 0 || !finite32(vLengthSquared) ||
		sightLengthSquared == 0 || !finite32(sightLengthSquared) || !finite32(dot) {
		return false
	}
	return dot < 0
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

// Coplanar reports whether four finite points have an absolute scalar triple
// product strictly less than tolerance. As in upstream, tolerance is absolute
// and scales with the cube of the coordinates. A degenerate point set is
// coplanar; a non-positive or non-finite tolerance returns false.
func Coplanar(a, b, c, d Vec3, tolerance float32) bool {
	if !a.IsFinite() || !b.IsFinite() || !c.IsFinite() || !d.IsFinite() ||
		!finite32(tolerance) || tolerance <= 0 {
		return false
	}
	volume := c.Sub(a).Dot(b.Sub(a).Cross(d.Sub(a)))
	return finite32(volume) && abs32(volume) < tolerance
}
