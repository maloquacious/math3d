package math3d

import "math"

// Quat is the upstream Quaternion value. X, Y, and Z are its vector part and
// W is its scalar part. The zero value is not a valid rotation.
//
// Quat values need not be normalized. Operations that require a unit
// quaternion are introduced with the quaternion algorithms.
type Quat struct {
	X, Y, Z, W float32
}

// NewQuat constructs a Quat without normalizing it.
func NewQuat(x, y, z, w float32) Quat { return Quat{X: x, Y: y, Z: z, W: w} }

// QuatFromVector constructs a Quat from its vector and scalar parts.
func QuatFromVector(vector Vec3, scalar float32) Quat {
	return NewQuat(vector.X, vector.Y, vector.Z, scalar)
}

// IdentityQuat returns the identity rotation.
func IdentityQuat() Quat { return NewQuat(0, 0, 0, 1) }

// Vector returns q's vector part.
func (q Quat) Vector() Vec3 { return V3(q.X, q.Y, q.Z) }

// AlmostEqual compares corresponding components using a strict absolute tolerance.
// It compares representations, so q and -q are not treated as equivalent rotations.
func (q Quat) AlmostEqual(other Quat, tolerance float32) bool {
	return abs32(q.X-other.X) < tolerance && abs32(q.Y-other.Y) < tolerance &&
		abs32(q.Z-other.Z) < tolerance && abs32(q.W-other.W) < tolerance
}

// DQuat converts q to float64 precision.
func (q Quat) DQuat() DQuat {
	return NewDQuat(float64(q.X), float64(q.Y), float64(q.Z), float64(q.W))
}

func (q Quat) Add(other Quat) Quat {
	return NewQuat(q.X+other.X, q.Y+other.Y, q.Z+other.Z, q.W+other.W)
}

func (q Quat) Sub(other Quat) Quat {
	return NewQuat(q.X-other.X, q.Y-other.Y, q.Z-other.Z, q.W-other.W)
}

func (q Quat) Negated() Quat { return NewQuat(-q.X, -q.Y, -q.Z, -q.W) }

func (q Quat) Scale(scalar float32) Quat {
	return NewQuat(q.X*scalar, q.Y*scalar, q.Z*scalar, q.W*scalar)
}

// Dot returns the four-component dot product.
func (q Quat) Dot(other Quat) float32 {
	return q.X*other.X + q.Y*other.Y + q.Z*other.Z + q.W*other.W
}

func (q Quat) LengthSquared() float32 { return q.Dot(q) }
func (q Quat) Length() float64        { return math.Sqrt(float64(q.LengthSquared())) }

// Normalized returns a unit quaternion. It fails for zero and non-finite lengths.
func (q Quat) Normalized() (Quat, bool) {
	length := q.Length()
	if length == 0 || math.IsNaN(length) || math.IsInf(length, 0) {
		return Quat{}, false
	}
	return q.Scale(1 / float32(length)), true
}

func (q Quat) Conjugated() Quat { return NewQuat(-q.X, -q.Y, -q.Z, q.W) }

// Inverse returns the multiplicative inverse. It fails for zero and non-finite norms.
func (q Quat) Inverse() (Quat, bool) {
	norm := q.LengthSquared()
	if norm == 0 || math.IsNaN(float64(norm)) || math.IsInf(float64(norm), 0) {
		return Quat{}, false
	}
	return q.Conjugated().Scale(1 / norm), true
}

// Mul returns the Hamilton product q·other. Quaternion multiplication applies
// other first and q second; use Concatenate for explicitly ordered rotations.
func (q Quat) Mul(other Quat) Quat {
	return NewQuat(
		q.W*other.X+q.X*other.W+q.Y*other.Z-q.Z*other.Y,
		q.W*other.Y-q.X*other.Z+q.Y*other.W+q.Z*other.X,
		q.W*other.Z+q.X*other.Y-q.Y*other.X+q.Z*other.W,
		q.W*other.W-q.X*other.X-q.Y*other.Y-q.Z*other.Z,
	)
}

// Div multiplies q by the inverse of other.
func (q Quat) Div(other Quat) (Quat, bool) {
	inverse, ok := other.Inverse()
	if !ok {
		return Quat{}, false
	}
	return q.Mul(inverse), true
}

// Concatenate returns the rotation first followed by second.
func Concatenate(first, second Quat) Quat { return second.Mul(first) }

// QuatFromAxisAngle creates a quaternion for an angle in radians. Axis must be unit length.
func QuatFromAxisAngle(axis Vec3, angle float32) Quat {
	half := float64(angle) * 0.5
	s := float32(math.Sin(half))
	return QuatFromVector(axis.Scale(s), float32(math.Cos(half)))
}

func QuatRotationX(angle float32) Quat { return QuatFromAxisAngle(V3(1, 0, 0), angle) }
func QuatRotationY(angle float32) Quat { return QuatFromAxisAngle(V3(0, 1, 0), angle) }
func QuatRotationZ(angle float32) Quat { return QuatFromAxisAngle(V3(0, 0, 1), angle) }

// QuatFromEulerAngles creates an intrinsic X, then Y, then Z rotation.
func QuatFromEulerAngles(angles Vec3) Quat {
	c1, s1 := math.Cos(float64(angles.X)/2), math.Sin(float64(angles.X)/2)
	c2, s2 := math.Cos(float64(angles.Y)/2), math.Sin(float64(angles.Y)/2)
	c3, s3 := math.Cos(float64(angles.Z)/2), math.Sin(float64(angles.Z)/2)
	return NewQuat(
		float32(s1*c2*c3+c1*s2*s3), float32(c1*s2*c3-s1*c2*s3),
		float32(c1*c2*s3+s1*s2*c3), float32(c1*c2*c3-s1*s2*s3),
	)
}

// EulerAngles returns intrinsic X, Y, Z angles in radians.
func (q Quat) EulerAngles() Vec3 {
	return V3(
		float32(math.Atan2(float64(-2*(q.Y*q.Z-q.W*q.X)), float64(q.W*q.W-q.X*q.X-q.Y*q.Y+q.Z*q.Z))),
		float32(math.Asin(float64(2*(q.X*q.Z+q.W*q.Y)))),
		float32(math.Atan2(float64(-2*(q.X*q.Y-q.W*q.Z)), float64(q.W*q.W+q.X*q.X-q.Y*q.Y-q.Z*q.Z))),
	)
}

// QuatFromYawPitchRoll creates roll (Z), then pitch (X), then yaw (Y).
func QuatFromYawPitchRoll(yaw, pitch, roll float32) Quat {
	sy, cy := float32(math.Sin(float64(yaw)/2)), float32(math.Cos(float64(yaw)/2))
	sp, cp := float32(math.Sin(float64(pitch)/2)), float32(math.Cos(float64(pitch)/2))
	sr, cr := float32(math.Sin(float64(roll)/2)), float32(math.Cos(float64(roll)/2))
	return NewQuat(cy*sp*cr+sy*cp*sr, sy*cp*cr-cy*sp*sr,
		cy*cp*sr-sy*sp*cr, cy*cp*cr+sy*sp*sr)
}

// Lerp linearly interpolates along the shortest rotational arc and normalizes the result.
func (q Quat) Lerp(other Quat, amount float32) (Quat, bool) {
	if q.Dot(other) < 0 {
		other = other.Negated()
	}
	return q.Scale(1 - amount).Add(other.Scale(amount)).Normalized()
}

// Slerp spherically interpolates along the shortest rotational arc.
func (q Quat) Slerp(other Quat, amount float32) Quat {
	cosOmega := q.Dot(other)
	if cosOmega < 0 {
		cosOmega, other = -cosOmega, other.Negated()
	}
	if cosOmega > 1-1e-6 {
		return q.Scale(1 - amount).Add(other.Scale(amount))
	}
	omega := math.Acos(float64(cosOmega))
	invSin := 1 / math.Sin(omega)
	a := float32(math.Sin(float64(1-amount)*omega) * invSin)
	b := float32(math.Sin(float64(amount)*omega) * invSin)
	return q.Scale(a).Add(other.Scale(b))
}

// Rotate rotates v. q must be a unit quaternion.
func (q Quat) Rotate(v Vec3) Vec3 {
	qv := q.Vector()
	t := qv.Cross(v).Scale(2)
	return v.Add(t.Scale(q.W)).Add(qv.Cross(t))
}

// RotationBetween returns a rotation from one unit vector to another. For
// opposite vectors, oppositeAxis must be unit and perpendicular to from.
func RotationBetween(from, to, oppositeAxis Vec3) (Quat, bool) {
	unit := func(v Vec3) bool {
		return v.IsFinite() && math.Abs(v.MagnitudeSquared()-1) < float64(Tolerance)
	}
	if !unit(from) || !unit(to) {
		return Quat{}, false
	}
	axis := from.Cross(to)
	if axis.MagnitudeSquared() > 0 {
		axis, ok := axis.Normalized()
		if !ok {
			return Quat{}, false
		}
		dot := max(float32(-1), min(float32(1), from.Dot(to)))
		return QuatFromAxisAngle(axis, float32(math.Acos(float64(dot)))), true
	}
	if from.Add(to).MagnitudeSquared() < float64(Tolerance*Tolerance) {
		if !unit(oppositeAxis) || abs32(from.Dot(oppositeAxis)) >= Tolerance {
			return Quat{}, false
		}
		return QuatFromAxisAngle(oppositeAxis, float32(math.Pi)), true
	}
	return IdentityQuat(), true
}

// DQuat is the upstream DQuaternion value. The zero value is not a valid
// rotation, and direct construction does not normalize the value.
type DQuat struct {
	X, Y, Z, W float64
}

// NewDQuat constructs a DQuat without normalizing it.
func NewDQuat(x, y, z, w float64) DQuat { return DQuat{X: x, Y: y, Z: z, W: w} }

// DQuatFromVector constructs a DQuat from its vector and scalar parts.
func DQuatFromVector(vector DVec3, scalar float64) DQuat {
	return NewDQuat(vector.X, vector.Y, vector.Z, scalar)
}

// IdentityDQuat returns the float64 identity rotation.
func IdentityDQuat() DQuat { return NewDQuat(0, 0, 0, 1) }

// Vector returns q's vector part.
func (q DQuat) Vector() DVec3 { return DV3(q.X, q.Y, q.Z) }

// AlmostEqual compares corresponding components using a strict absolute tolerance.
// It compares representations, so q and -q are not treated as equivalent rotations.
func (q DQuat) AlmostEqual(other DQuat, tolerance float64) bool {
	return almostEqual64(tolerance, q.X-other.X, q.Y-other.Y, q.Z-other.Z, q.W-other.W)
}

// Quat converts q to float32 precision.
func (q DQuat) Quat() Quat {
	return NewQuat(float32(q.X), float32(q.Y), float32(q.Z), float32(q.W))
}

// AxisAngle is the upstream AxisAngle value. Its float64 precision matches the
// source despite the absence of a D prefix. Angle is in radians. A zero axis is
// safe to represent but does not define a rotation.
type AxisAngle struct {
	Axis  DVec3
	Angle float64
}

// NewAxisAngle constructs an AxisAngle without normalizing axis.
func NewAxisAngle(axis DVec3, angle float64) AxisAngle {
	return AxisAngle{Axis: axis, Angle: angle}
}

// AlmostEqual compares the stored axis and angle, not rotational equivalence.
func (a AxisAngle) AlmostEqual(other AxisAngle, tolerance float64) bool {
	return a.Axis.AlmostEqual(other.Axis, tolerance) &&
		almostEqual64(tolerance, a.Angle-other.Angle)
}

// Euler is the upstream Euler value. Its yaw, pitch, and roll angles are
// radians. The zero value represents no rotation; conversion conventions are
// defined with the quaternion algorithms.
type Euler struct {
	Yaw, Pitch, Roll float32
}

// NewEuler constructs an Euler value in yaw, pitch, roll order.
func NewEuler(yaw, pitch, roll float32) Euler {
	return Euler{Yaw: yaw, Pitch: pitch, Roll: roll}
}

// AlmostEqual compares corresponding angles using a strict absolute tolerance.
func (e Euler) AlmostEqual(other Euler, tolerance float32) bool {
	return abs32(e.Yaw-other.Yaw) < tolerance &&
		abs32(e.Pitch-other.Pitch) < tolerance &&
		abs32(e.Roll-other.Roll) < tolerance
}
