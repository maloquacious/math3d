package math3d

import (
	"math"
	"testing"
)

const algorithmTolerance = float32(1e-5)

func sameRotation(a, b Quat, tolerance float32) bool {
	return a.AlmostEqual(b, tolerance) || a.AlmostEqual(b.Negated(), tolerance)
}

func TestQuaternionAlgebraAndFailureResults(t *testing.T) {
	a, b := NewQuat(1, 2, 3, 4), NewQuat(5, 6, 7, 8)
	if got, want := a.Mul(b), NewQuat(24, 48, 48, -6); got != want {
		t.Fatalf("Mul = %#v, want %#v", got, want)
	}
	if got := Concatenate(a, b); got != b.Mul(a) {
		t.Fatalf("Concatenate = %#v, want %#v", got, b.Mul(a))
	}
	q, ok := NewQuat(1, 2, 3, 4).Normalized()
	if !ok || abs32(q.LengthSquared()-1) > algorithmTolerance {
		t.Fatalf("Normalized = %#v, %v", q, ok)
	}
	inverse, ok := a.Inverse()
	if !ok || !a.Mul(inverse).AlmostEqual(IdentityQuat(), algorithmTolerance) {
		t.Fatalf("Inverse = %#v, %v", inverse, ok)
	}
	if _, ok := (Quat{}).Normalized(); ok {
		t.Fatal("zero quaternion normalization succeeded")
	}
	if _, ok := (Quat{}).Inverse(); ok {
		t.Fatal("zero quaternion inversion succeeded")
	}
}

func TestQuaternionCreationInterpolationAndRotation(t *testing.T) {
	q := QuatRotationZ(HalfPi)
	if got := q.Rotate(V3(1, 0, 0)); !got.AlmostEqual(V3(0, 1, 0), algorithmTolerance) {
		t.Fatalf("Rotate = %#v", got)
	}
	angles := V3(Pi/5, 2*Pi/7, Pi/3)
	if got := QuatFromEulerAngles(angles).EulerAngles(); !got.AlmostEqual(angles, 1e-3) {
		t.Fatalf("Euler round trip = %#v, want %#v", got, angles)
	}
	ypr := QuatFromYawPitchRoll(.3, .4, .5)
	want := QuatRotationY(.3).Mul(QuatRotationX(.4)).Mul(QuatRotationZ(.5))
	if !sameRotation(ypr, want, algorithmTolerance) {
		t.Fatalf("yaw/pitch/roll = %#v, want %#v", ypr, want)
	}
	start, end := QuatRotationZ(Pi/18), QuatRotationZ(Pi/6)
	if got := start.Slerp(end, .5); !sameRotation(got, QuatRotationZ(Pi/9), algorithmTolerance) {
		t.Fatalf("Slerp midpoint = %#v", got)
	}
	if got, ok := start.Lerp(end, .5); !ok || !sameRotation(got, QuatRotationZ(Pi/9), 1e-4) {
		t.Fatalf("Lerp midpoint = %#v, %v", got, ok)
	}
	between, ok := RotationBetween(V3(1, 0, 0), V3(0, 1, 0), V3(0, 0, 1))
	if !ok || !between.Rotate(V3(1, 0, 0)).AlmostEqual(V3(0, 1, 0), algorithmTolerance) {
		t.Fatalf("RotationBetween = %#v, %v", between, ok)
	}
	if _, ok := RotationBetween(V3(2, 0, 0), V3(0, 1, 0), V3(0, 0, 1)); ok {
		t.Fatal("RotationBetween accepted a non-unit input")
	}
}

func TestMatrixRowVectorCompositionAndQuaternionRoundTrip(t *testing.T) {
	translation := TranslationMat4(V3(10, 20, 30))
	if got := translation.TransformPoint(V3(1, 2, 3)); got != V3(11, 22, 33) {
		t.Fatalf("translated point = %#v", got)
	}
	if got := translation.TransformDirection(V3(1, 2, 3)); got != V3(1, 2, 3) {
		t.Fatalf("translated direction = %#v", got)
	}
	composed := ScaleMat4(V3(2, 2, 2)).Mul(translation)
	if got := composed.TransformPoint(V3(1, 0, 0)); got != V3(12, 20, 30) {
		t.Fatalf("composition order = %#v", got)
	}
	q := QuatRotationZ(.7).Mul(QuatRotationY(-.3)).Mul(QuatRotationX(.2))
	m := Mat4FromQuat(q)
	if !sameRotation(m.Quat(), q, algorithmTolerance) {
		t.Fatalf("matrix/quaternion round trip = %#v, want %#v", m.Quat(), q)
	}
	want := RotationXMat4(.2).Mul(RotationYMat4(-.3)).Mul(RotationZMat4(.7))
	if !m.AlmostEqual(want, algorithmTolerance) {
		t.Fatalf("quaternion matrix order mismatch\n got %#v\nwant %#v", m, want)
	}
}

func TestMatrixInverseDeterminantTransposeAndReflection(t *testing.T) {
	m := ComposeMat4(V3(2, 3, 4), QuatFromYawPitchRoll(.2, .3, .4), V3(5, 6, 7))
	inverse, ok := m.Inverse()
	if !ok || !m.Mul(inverse).AlmostEqual(IdentityMat4(), 2e-5) {
		t.Fatalf("inverse = %#v, %v", inverse, ok)
	}
	if _, ok := (Mat4{}).Inverse(); ok {
		t.Fatal("singular matrix inversion succeeded")
	}
	if got := m.Transposed().Transposed(); got != m {
		t.Fatal("double transpose did not round trip")
	}
	if abs32(m.Determinant()-24) > 1e-4 {
		t.Fatalf("determinant = %v", m.Determinant())
	}
	if !ScaleMat4(V3(-1, 1, 1)).IsReflection() || ScaleMat4(V3(1, 1, 1)).IsReflection() {
		t.Fatal("reflection detection failed")
	}
}

func TestMatrixProjectionViewAndDecomposition(t *testing.T) {
	p, ok := PerspectiveFOVMat4(HalfPi, 1, 1, 10)
	if !ok {
		t.Fatal("valid perspective rejected")
	}
	project := func(v Vec3) Vec3 {
		r := p.TransformVec4(V4(v.X, v.Y, v.Z, 1))
		return V3(r.X/r.W, r.Y/r.W, r.Z/r.W)
	}
	if got := project(V3(0, 0, -1)); abs32(got.Z) > algorithmTolerance {
		t.Fatalf("near depth = %v", got.Z)
	}
	if got := project(V3(0, 0, -10)); abs32(got.Z-1) > algorithmTolerance {
		t.Fatalf("far depth = %v", got.Z)
	}
	if _, ok := PerspectiveFOVMat4(0, 1, 1, 10); ok {
		t.Fatal("invalid FOV accepted")
	}
	view, ok := LookAtMat4(V3(0, 0, 10), Vec3{}, V3(0, 1, 0))
	if !ok || !view.TransformPoint(V3(0, 0, 10)).AlmostEqual(Vec3{}, algorithmTolerance) {
		t.Fatalf("look-at = %#v, %v", view, ok)
	}
	scale, rotation, translation := V3(2, 3, 4), QuatFromYawPitchRoll(.2, .3, .4), V3(5, 6, 7)
	m := ComposeMat4(scale, rotation, translation)
	gotScale, gotRotation, gotTranslation, ok := m.Decompose()
	if !ok || !gotScale.AlmostEqual(scale, algorithmTolerance) || gotTranslation != translation || !sameRotation(gotRotation, rotation, algorithmTolerance) {
		t.Fatalf("decompose = %#v %#v %#v %v", gotScale, gotRotation, gotTranslation, ok)
	}
	if _, _, _, ok := NewMat4(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16).Decompose(); ok {
		t.Fatal("non-SRT decomposition succeeded")
	}
	if _, _, _, ok := ScaleMat4(Vec3{}).Decompose(); !ok {
		t.Fatal("zero scale decomposition should retain upstream behavior")
	}
	if math.IsNaN(float64(p.M11)) {
		t.Fatal("projection unexpectedly contains NaN")
	}
}

func TestMatrixBillboardReflectionAndShadowFactories(t *testing.T) {
	billboard, ok := BillboardMat4(V3(1, 2, 3), V3(1, 2, 10), V3(0, 1, 0), V3(0, 0, -1))
	if !ok || billboard.Translation() != V3(1, 2, 3) || !V3(billboard.M31, billboard.M32, billboard.M33).AlmostEqual(V3(0, 0, -1), algorithmTolerance) {
		t.Fatalf("billboard = %#v, %v", billboard, ok)
	}
	reflection, ok := ReflectionMat4(NewPlane(V3(1, 0, 0), -2))
	if !ok || !reflection.TransformPoint(V3(3, 4, 5)).AlmostEqual(V3(1, 4, 5), algorithmTolerance) {
		t.Fatalf("reflection = %#v, %v", reflection, ok)
	}
	shadow, ok := ShadowMat4(V3(0, -1, 0), NewPlane(V3(0, 1, 0), 0))
	if !ok {
		t.Fatal("valid shadow projection rejected")
	}
	r := shadow.TransformVec4(V4(2, 3, 4, 1))
	projected := V3(r.X/r.W, r.Y/r.W, r.Z/r.W)
	if !projected.AlmostEqual(V3(2, 0, 4), algorithmTolerance) {
		t.Fatalf("shadow point = %#v", projected)
	}
	if _, ok := ReflectionMat4(Plane{}); ok {
		t.Fatal("zero plane reflection succeeded")
	}
}
