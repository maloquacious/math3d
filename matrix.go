package math3d

import "math"

// Mat4 is the upstream Matrix4x4 value. Fields are in row-major order, but
// transformations use row vectors: translation occupies M41, M42, and M43.
// The zero value is a zero matrix, not an identity transform.
type Mat4 struct {
	M11, M12, M13, M14 float32
	M21, M22, M23, M24 float32
	M31, M32, M33, M34 float32
	M41, M42, M43, M44 float32
}

// NewMat4 constructs a Mat4 from rows in row-major order.
func NewMat4(
	m11, m12, m13, m14 float32,
	m21, m22, m23, m24 float32,
	m31, m32, m33, m34 float32,
	m41, m42, m43, m44 float32,
) Mat4 {
	return Mat4{
		M11: m11, M12: m12, M13: m13, M14: m14,
		M21: m21, M22: m22, M23: m23, M24: m24,
		M31: m31, M32: m32, M33: m33, M34: m34,
		M41: m41, M42: m42, M43: m43, M44: m44,
	}
}

// Mat4FromRows constructs a Mat4 from four rows.
func Mat4FromRows(row0, row1, row2, row3 Vec4) Mat4 {
	return NewMat4(
		row0.X, row0.Y, row0.Z, row0.W,
		row1.X, row1.Y, row1.Z, row1.W,
		row2.X, row2.Y, row2.Z, row2.W,
		row3.X, row3.Y, row3.Z, row3.W,
	)
}

// IdentityMat4 returns the identity matrix.
func IdentityMat4() Mat4 {
	return NewMat4(
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
}

// At returns the element at the zero-based row and column. It panics if either
// index is outside [0, 4).
func (m Mat4) At(row, column int) float32 {
	if row < 0 || row >= 4 || column < 0 || column >= 4 {
		panic("math3d: Mat4 index out of range")
	}
	switch row*4 + column {
	case 0:
		return m.M11
	case 1:
		return m.M12
	case 2:
		return m.M13
	case 3:
		return m.M14
	case 4:
		return m.M21
	case 5:
		return m.M22
	case 6:
		return m.M23
	case 7:
		return m.M24
	case 8:
		return m.M31
	case 9:
		return m.M32
	case 10:
		return m.M33
	case 11:
		return m.M34
	case 12:
		return m.M41
	case 13:
		return m.M42
	case 14:
		return m.M43
	case 15:
		return m.M44
	}
	panic("math3d: unreachable Mat4 index")
}

// Row returns the zero-based row. It panics if row is outside [0, 4).
func (m Mat4) Row(row int) Vec4 {
	switch row {
	case 0:
		return V4(m.M11, m.M12, m.M13, m.M14)
	case 1:
		return V4(m.M21, m.M22, m.M23, m.M24)
	case 2:
		return V4(m.M31, m.M32, m.M33, m.M34)
	case 3:
		return V4(m.M41, m.M42, m.M43, m.M44)
	default:
		panic("math3d: Mat4 row out of range")
	}
}

// Column returns the zero-based column. It panics if column is outside [0, 4).
func (m Mat4) Column(column int) Vec4 {
	switch column {
	case 0:
		return V4(m.M11, m.M21, m.M31, m.M41)
	case 1:
		return V4(m.M12, m.M22, m.M32, m.M42)
	case 2:
		return V4(m.M13, m.M23, m.M33, m.M43)
	case 3:
		return V4(m.M14, m.M24, m.M34, m.M44)
	default:
		panic("math3d: Mat4 column out of range")
	}
}

// Translation returns the row-vector translation stored in m.
func (m Mat4) Translation() Vec3 { return V3(m.M41, m.M42, m.M43) }

// WithTranslation returns m with its row-vector translation replaced.
func (m Mat4) WithTranslation(translation Vec3) Mat4 {
	m.M41, m.M42, m.M43 = translation.X, translation.Y, translation.Z
	return m
}

// IsIdentity reports exact representation equality with the identity matrix.
func (m Mat4) IsIdentity() bool { return m == IdentityMat4() }

// AlmostEqual compares corresponding elements using a strict absolute tolerance.
func (m Mat4) AlmostEqual(other Mat4, tolerance float32) bool {
	for row := 0; row < 4; row++ {
		for column := 0; column < 4; column++ {
			if !(abs32(m.At(row, column)-other.At(row, column)) < tolerance) {
				return false
			}
		}
	}
	return true
}

func (m Mat4) Add(o Mat4) Mat4 {
	return NewMat4(m.M11+o.M11, m.M12+o.M12, m.M13+o.M13, m.M14+o.M14,
		m.M21+o.M21, m.M22+o.M22, m.M23+o.M23, m.M24+o.M24,
		m.M31+o.M31, m.M32+o.M32, m.M33+o.M33, m.M34+o.M34,
		m.M41+o.M41, m.M42+o.M42, m.M43+o.M43, m.M44+o.M44)
}

func (m Mat4) Sub(o Mat4) Mat4 { return m.Add(o.Scale(-1)) }

func (m Mat4) Scale(s float32) Mat4 {
	return NewMat4(m.M11*s, m.M12*s, m.M13*s, m.M14*s, m.M21*s, m.M22*s, m.M23*s, m.M24*s,
		m.M31*s, m.M32*s, m.M33*s, m.M34*s, m.M41*s, m.M42*s, m.M43*s, m.M44*s)
}

// Mul multiplies two matrices. With row vectors, m is applied before o.
func (m Mat4) Mul(o Mat4) Mat4 {
	return NewMat4(
		m.Row(0).Dot(o.Column(0)), m.Row(0).Dot(o.Column(1)), m.Row(0).Dot(o.Column(2)), m.Row(0).Dot(o.Column(3)),
		m.Row(1).Dot(o.Column(0)), m.Row(1).Dot(o.Column(1)), m.Row(1).Dot(o.Column(2)), m.Row(1).Dot(o.Column(3)),
		m.Row(2).Dot(o.Column(0)), m.Row(2).Dot(o.Column(1)), m.Row(2).Dot(o.Column(2)), m.Row(2).Dot(o.Column(3)),
		m.Row(3).Dot(o.Column(0)), m.Row(3).Dot(o.Column(1)), m.Row(3).Dot(o.Column(2)), m.Row(3).Dot(o.Column(3)),
	)
}

func (m Mat4) Transposed() Mat4 {
	return Mat4FromRows(m.Column(0), m.Column(1), m.Column(2), m.Column(3))
}

func (m Mat4) Lerp(o Mat4, amount float32) Mat4 { return m.Add(o.Sub(m).Scale(amount)) }

// TransformVec4 applies m to the row vector v.
func (m Mat4) TransformVec4(v Vec4) Vec4 {
	return V4(v.Dot(m.Column(0)), v.Dot(m.Column(1)), v.Dot(m.Column(2)), v.Dot(m.Column(3)))
}

// TransformPoint transforms a position with homogeneous W=1.
func (m Mat4) TransformPoint(v Vec3) Vec3 {
	r := m.TransformVec4(V4(v.X, v.Y, v.Z, 1))
	return V3(r.X, r.Y, r.Z)
}

// TransformDirection transforms a direction with homogeneous W=0.
func (m Mat4) TransformDirection(v Vec3) Vec3 {
	r := m.TransformVec4(V4(v.X, v.Y, v.Z, 0))
	return V3(r.X, r.Y, r.Z)
}

// TransformNormal transforms a surface normal using the inverse transpose of
// m, as required by the row-vector convention. Translation has no effect, and
// the result is not normalized. It fails for a non-finite normal, a singular or
// non-finite matrix, or a non-finite result.
func (m Mat4) TransformNormal(normal Vec3) (Vec3, bool) {
	if !normal.IsFinite() {
		return Vec3{}, false
	}
	inverse, ok := m.Inverse()
	if !ok {
		return Vec3{}, false
	}
	result := inverse.Transposed().TransformDirection(normal)
	if !result.IsFinite() {
		return Vec3{}, false
	}
	return result, true
}

func TranslationMat4(v Vec3) Mat4 { return IdentityMat4().WithTranslation(v) }

func ScaleMat4(v Vec3) Mat4 {
	return NewMat4(v.X, 0, 0, 0, 0, v.Y, 0, 0, 0, 0, v.Z, 0, 0, 0, 0, 1)
}

// ScaleAroundMat4 creates a scale about center.
func ScaleAroundMat4(v, center Vec3) Mat4 {
	return TranslationMat4(center.Negated()).Mul(ScaleMat4(v)).Mul(TranslationMat4(center))
}

func RotationXMat4(angle float32) Mat4 {
	c, s := float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle)))
	return NewMat4(1, 0, 0, 0, 0, c, s, 0, 0, -s, c, 0, 0, 0, 0, 1)
}
func RotationYMat4(angle float32) Mat4 {
	c, s := float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle)))
	return NewMat4(c, 0, -s, 0, 0, 1, 0, 0, s, 0, c, 0, 0, 0, 0, 1)
}
func RotationZMat4(angle float32) Mat4 {
	c, s := float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle)))
	return NewMat4(c, s, 0, 0, -s, c, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1)
}

func RotationAroundMat4(rotation Mat4, center Vec3) Mat4 {
	return TranslationMat4(center.Negated()).Mul(rotation).Mul(TranslationMat4(center))
}

// BillboardMat4 creates a spherical billboard at objectPosition.
func BillboardMat4(objectPosition, cameraPosition, cameraUp, cameraForward Vec3) (Mat4, bool) {
	z := objectPosition.Sub(cameraPosition)
	if z.MagnitudeSquared() < 1e-4 {
		z = cameraForward.Negated()
	} else {
		var ok bool
		z, ok = z.Normalized()
		if !ok {
			return Mat4{}, false
		}
	}
	x, ok := cameraUp.Cross(z).Normalized()
	if !ok {
		return Mat4{}, false
	}
	y := z.Cross(x)
	return NewMat4(x.X, x.Y, x.Z, 0, y.X, y.Y, y.Z, 0, z.X, z.Y, z.Z, 0,
		objectPosition.X, objectPosition.Y, objectPosition.Z, 1), true
}

// ConstrainedBillboardMat4 creates a cylindrical billboard around a unit axis.
func ConstrainedBillboardMat4(objectPosition, cameraPosition, rotateAxis, cameraForward, objectForward Vec3) (Mat4, bool) {
	face := objectPosition.Sub(cameraPosition)
	if face.MagnitudeSquared() < 1e-4 {
		face = cameraForward.Negated()
	} else {
		var ok bool
		face, ok = face.Normalized()
		if !ok {
			return Mat4{}, false
		}
	}
	const minAngle = 1 - float32(0.1)*(Pi/180)
	y := rotateAxis
	var x, z Vec3
	if abs32(rotateAxis.Dot(face)) > minAngle {
		z = objectForward
		if abs32(rotateAxis.Dot(z)) > minAngle {
			if abs32(rotateAxis.Z) > minAngle {
				z = V3(1, 0, 0)
			} else {
				z = V3(0, 0, -1)
			}
		}
		x, _ = rotateAxis.Cross(z).Normalized()
		z, _ = x.Cross(rotateAxis).Normalized()
	} else {
		x, _ = rotateAxis.Cross(face).Normalized()
		z, _ = x.Cross(y).Normalized()
	}
	if !x.IsFinite() || !z.IsFinite() || x == (Vec3{}) || z == (Vec3{}) {
		return Mat4{}, false
	}
	return NewMat4(x.X, x.Y, x.Z, 0, y.X, y.Y, y.Z, 0, z.X, z.Y, z.Z, 0,
		objectPosition.X, objectPosition.Y, objectPosition.Z, 1), true
}

// Mat4FromAxisAngle creates a rotation around a unit axis.
func Mat4FromAxisAngle(axis Vec3, angle float32) Mat4 {
	return Mat4FromQuat(QuatFromAxisAngle(axis, angle))
}

// Mat4FromYawPitchRoll creates roll (Z), then pitch (X), then yaw (Y).
func Mat4FromYawPitchRoll(yaw, pitch, roll float32) Mat4 {
	return Mat4FromQuat(QuatFromYawPitchRoll(yaw, pitch, roll))
}

// Mat4FromQuat creates the row-vector matrix represented by q. The quaternion
// must have unit length for the matrix to be a pure rotation.
func Mat4FromQuat(q Quat) Mat4 {
	xx, yy, zz := q.X*q.X, q.Y*q.Y, q.Z*q.Z
	xy, xz, yz := q.X*q.Y, q.X*q.Z, q.Y*q.Z
	wx, wy, wz := q.W*q.X, q.W*q.Y, q.W*q.Z
	return NewMat4(
		1-2*(yy+zz), 2*(xy+wz), 2*(xz-wy), 0,
		2*(xy-wz), 1-2*(zz+xx), 2*(yz+wx), 0,
		2*(xz+wy), 2*(yz-wx), 1-2*(yy+xx), 0,
		0, 0, 0, 1,
	)
}

// Quat returns the quaternion represented by the upper-left rotation matrix.
func (m Mat4) Quat() Quat {
	trace := m.M11 + m.M22 + m.M33
	if trace > 0 {
		s := float32(math.Sqrt(float64(trace + 1)))
		return NewQuat((m.M23-m.M32)*0.5/s, (m.M31-m.M13)*0.5/s, (m.M12-m.M21)*0.5/s, s*0.5)
	}
	if m.M11 >= m.M22 && m.M11 >= m.M33 {
		s := float32(math.Sqrt(float64(1 + m.M11 - m.M22 - m.M33)))
		return NewQuat(s*.5, (m.M12+m.M21)*.5/s, (m.M13+m.M31)*.5/s, (m.M23-m.M32)*.5/s)
	}
	if m.M22 > m.M33 {
		s := float32(math.Sqrt(float64(1 + m.M22 - m.M11 - m.M33)))
		return NewQuat((m.M21+m.M12)*.5/s, s*.5, (m.M32+m.M23)*.5/s, (m.M31-m.M13)*.5/s)
	}
	s := float32(math.Sqrt(float64(1 + m.M33 - m.M11 - m.M22)))
	return NewQuat((m.M31+m.M13)*.5/s, (m.M32+m.M23)*.5/s, s*.5, (m.M12-m.M21)*.5/s)
}

func (m Mat4) RotationDeterminant() float32 {
	return m.M11*(m.M22*m.M33-m.M23*m.M32) - m.M12*(m.M21*m.M33-m.M23*m.M31) + m.M13*(m.M21*m.M32-m.M22*m.M31)
}
func (m Mat4) IsReflection() bool { return m.RotationDeterminant() < 0 }

func (m Mat4) Determinant() float32 {
	a, b, c, d := m.M11, m.M12, m.M13, m.M14
	e, f, g, h := m.M21, m.M22, m.M23, m.M24
	i, j, k, l := m.M31, m.M32, m.M33, m.M34
	m0, n, o, p := m.M41, m.M42, m.M43, m.M44
	kpLo, jpLn, joKn := k*p-l*o, j*p-l*n, j*o-k*n
	ipLm, ioKm, inJm := i*p-l*m0, i*o-k*m0, i*n-j*m0
	return a*(f*kpLo-g*jpLn+h*joKn) - b*(e*kpLo-g*ipLm+h*ioKm) +
		c*(e*jpLn-f*ipLm+h*inJm) - d*(e*joKn-f*ioKm+g*inJm)
}

// Inverse returns the inverse matrix, or false when m is singular or non-finite.
func (m Mat4) Inverse() (Mat4, bool) {
	// Gauss-Jordan elimination in float64 avoids duplicating an opaque cofactor expansion.
	var a [4][8]float64
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			a[r][c] = float64(m.At(r, c))
		}
		a[r][r+4] = 1
	}
	for c := 0; c < 4; c++ {
		pivot := c
		for r := c + 1; r < 4; r++ {
			if math.Abs(a[r][c]) > math.Abs(a[pivot][c]) {
				pivot = r
			}
		}
		if a[pivot][c] == 0 || math.IsNaN(a[pivot][c]) || math.IsInf(a[pivot][c], 0) {
			return Mat4{}, false
		}
		a[c], a[pivot] = a[pivot], a[c]
		div := a[c][c]
		for j := 0; j < 8; j++ {
			a[c][j] /= div
		}
		for r := 0; r < 4; r++ {
			if r == c {
				continue
			}
			factor := a[r][c]
			for j := 0; j < 8; j++ {
				a[r][j] -= factor * a[c][j]
			}
		}
	}
	return NewMat4(float32(a[0][4]), float32(a[0][5]), float32(a[0][6]), float32(a[0][7]),
		float32(a[1][4]), float32(a[1][5]), float32(a[1][6]), float32(a[1][7]),
		float32(a[2][4]), float32(a[2][5]), float32(a[2][6]), float32(a[2][7]),
		float32(a[3][4]), float32(a[3][5]), float32(a[3][6]), float32(a[3][7])), true
}

// PerspectiveFOVMat4 creates a right-handed, negative-Z-forward perspective
// projection with depth [0, 1]. It returns false when fieldOfView is at most
// zero or at least Pi, near or far is at most zero, or near is not less than
// far. An infinite positive far plane is supported. All other inputs must be
// finite. Aspect is the viewport width divided by its height and must be
// non-zero.
func PerspectiveFOVMat4(fieldOfView, aspect, near, far float32) (Mat4, bool) {
	if fieldOfView <= 0 || fieldOfView >= Pi || near <= 0 || far <= 0 || near >= far {
		return Mat4{}, false
	}
	y := 1 / float32(math.Tan(float64(fieldOfView)*.5))
	rangeScale := far / (near - far)
	if math.IsInf(float64(far), 1) {
		rangeScale = -1
	}
	return NewMat4(y/aspect, 0, 0, 0, 0, y, 0, 0, 0, 0, rangeScale, -1, 0, 0, near*rangeScale, 0), true
}

// PerspectiveMat4 creates a perspective projection from near-plane dimensions.
func PerspectiveMat4(width, height, near, far float32) (Mat4, bool) {
	if near <= 0 || far <= 0 || near >= far {
		return Mat4{}, false
	}
	rangeScale := far / (near - far)
	if math.IsInf(float64(far), 1) {
		rangeScale = -1
	}
	return NewMat4(2*near/width, 0, 0, 0, 0, 2*near/height, 0, 0,
		0, 0, rangeScale, -1, 0, 0, near*rangeScale, 0), true
}

// PerspectiveOffCenterMat4 creates an asymmetric perspective projection.
func PerspectiveOffCenterMat4(left, right, bottom, top, near, far float32) (Mat4, bool) {
	if near <= 0 || far <= 0 || near >= far {
		return Mat4{}, false
	}
	rangeScale := far / (near - far)
	if math.IsInf(float64(far), 1) {
		rangeScale = -1
	}
	return NewMat4(2*near/(right-left), 0, 0, 0, 0, 2*near/(top-bottom), 0, 0,
		(left+right)/(right-left), (top+bottom)/(top-bottom), rangeScale, -1,
		0, 0, near*rangeScale, 0), true
}

func OrthographicMat4(width, height, near, far float32) Mat4 {
	return NewMat4(2/width, 0, 0, 0, 0, 2/height, 0, 0, 0, 0, 1/(near-far), 0, 0, 0, near/(near-far), 1)
}

func OrthographicOffCenterMat4(left, right, bottom, top, near, far float32) Mat4 {
	return NewMat4(2/(right-left), 0, 0, 0, 0, 2/(top-bottom), 0, 0, 0, 0, 1/(near-far), 0,
		(left+right)/(left-right), (top+bottom)/(bottom-top), near/(near-far), 1)
}

// LookAtMat4 creates a right-handed view matrix with negative Z forward. It
// fails when position and target do not define a finite direction or when up
// and that direction cannot define a finite perpendicular axis.
func LookAtMat4(position, target, up Vec3) (Mat4, bool) {
	z, ok := position.Sub(target).Normalized()
	if !ok {
		return Mat4{}, false
	}
	x, ok := up.Cross(z).Normalized()
	if !ok {
		return Mat4{}, false
	}
	y := z.Cross(x)
	return NewMat4(x.X, y.X, z.X, 0, x.Y, y.Y, z.Y, 0, x.Z, y.Z, z.Z, 0,
		-x.Dot(position), -y.Dot(position), -z.Dot(position), 1), true
}

// WorldMat4 creates a world matrix whose local forward direction is negative Z.
func WorldMat4(position, forward, up Vec3) (Mat4, bool) {
	z, ok := forward.Negated().Normalized()
	if !ok {
		return Mat4{}, false
	}
	x, ok := up.Cross(z).Normalized()
	if !ok {
		return Mat4{}, false
	}
	y := z.Cross(x)
	return NewMat4(x.X, x.Y, x.Z, 0, y.X, y.Y, y.Z, 0, z.X, z.Y, z.Z, 0, position.X, position.Y, position.Z, 1), true
}

// ComposeMat4 creates scale, then rotation, then translation for row vectors.
// Rotation must be a unit quaternion for the linear part to contain only the
// supplied scale and rotation.
func ComposeMat4(scale Vec3, rotation Quat, translation Vec3) Mat4 {
	return ScaleMat4(scale).Mul(Mat4FromQuat(rotation)).Mul(TranslationMat4(translation))
}

// Decompose extracts scale, rotation, and translation from an affine SRT
// matrix. It returns false when the linear part is non-finite or cannot be
// represented as scale followed by an orthonormal rotation within its fixed
// tolerance, as with shear.
func (m Mat4) Decompose() (scale Vec3, rotation Quat, translation Vec3, ok bool) {
	rows := [3]Vec3{V3(m.M11, m.M12, m.M13), V3(m.M21, m.M22, m.M23), V3(m.M31, m.M32, m.M33)}
	scales := [3]float32{float32(rows[0].Magnitude()), float32(rows[1].Magnitude()), float32(rows[2].Magnitude())}
	order := [3]int{0, 1, 2}
	for i := 0; i < 2; i++ {
		for j := i + 1; j < 3; j++ {
			if scales[order[j]] > scales[order[i]] {
				order[i], order[j] = order[j], order[i]
			}
		}
	}
	canonical := [3]Vec3{V3(1, 0, 0), V3(0, 1, 0), V3(0, 0, 1)}
	a, b, c := order[0], order[1], order[2]
	if scales[a] < 1e-4 {
		rows[a] = canonical[a]
	}
	rows[a], _ = rows[a].Normalized()
	if scales[b] < 1e-4 {
		least := 0
		leastMagnitude := abs32(rows[a].X)
		if abs32(rows[a].Y) < leastMagnitude {
			least, leastMagnitude = 1, abs32(rows[a].Y)
		}
		if abs32(rows[a].Z) < leastMagnitude {
			least = 2
		}
		rows[b] = rows[a].Cross(canonical[least])
	}
	rows[b], _ = rows[b].Normalized()
	if scales[c] < 1e-4 {
		rows[c] = rows[a].Cross(rows[b])
	}
	rows[c], _ = rows[c].Normalized()
	scale = V3(scales[0], scales[1], scales[2])
	rot := NewMat4(rows[0].X, rows[0].Y, rows[0].Z, 0, rows[1].X, rows[1].Y, rows[1].Z, 0,
		rows[2].X, rows[2].Y, rows[2].Z, 0, 0, 0, 0, 1)
	det := rot.RotationDeterminant()
	if det < 0 {
		switch a {
		case 0:
			scale.X = -scale.X
		case 1:
			scale.Y = -scale.Y
		case 2:
			scale.Z = -scale.Z
		}
		rows[a] = rows[a].Negated()
		rot = Mat4FromRows(V4(rows[0].X, rows[0].Y, rows[0].Z, 0), V4(rows[1].X, rows[1].Y, rows[1].Z, 0), V4(rows[2].X, rows[2].Y, rows[2].Z, 0), V4(0, 0, 0, 1))
		det = -det
	}
	if abs32(det-1) >= 1e-2 {
		return scale, IdentityQuat(), m.Translation(), false
	}
	return scale, rot.Quat(), m.Translation(), true
}

// ReflectionMat4 creates a reflection about dot(Normal, point)+D=0.
func ReflectionMat4(plane Plane) (Mat4, bool) {
	p, ok := plane.Normalized()
	if !ok {
		return Mat4{}, false
	}
	a, b, c := p.Normal.X, p.Normal.Y, p.Normal.Z
	fa, fb, fc := -2*a, -2*b, -2*c
	return NewMat4(fa*a+1, fb*a, fc*a, 0, fa*b, fb*b+1, fc*b, 0,
		fa*c, fb*c, fc*c+1, 0, fa*p.D, fb*p.D, fc*p.D, 1), true
}

// ShadowMat4 projects onto plane along lightDirection.
func ShadowMat4(lightDirection Vec3, plane Plane) (Mat4, bool) {
	p, ok := plane.Normalized()
	if !ok {
		return Mat4{}, false
	}
	dot := p.Normal.Dot(lightDirection)
	a, b, c, d := -p.Normal.X, -p.Normal.Y, -p.Normal.Z, -p.D
	return NewMat4(
		a*lightDirection.X+dot, a*lightDirection.Y, a*lightDirection.Z, 0,
		b*lightDirection.X, b*lightDirection.Y+dot, b*lightDirection.Z, 0,
		c*lightDirection.X, c*lightDirection.Y, c*lightDirection.Z+dot, 0,
		d*lightDirection.X, d*lightDirection.Y, d*lightDirection.Z, dot,
	), true
}
