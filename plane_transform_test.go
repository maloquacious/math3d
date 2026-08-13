package math3d

import (
	"math"
	"testing"
)

func TestPlaneFactoriesUseConsistentEquation(t *testing.T) {
	a, b, c := V3(2, 0, 0), V3(2, 1, 0), V3(2, 0, 1)
	fromPoints, ok := PlaneFromPoints(a, b, c)
	if !ok {
		t.Fatal("PlaneFromPoints rejected non-collinear points")
	}
	fromNormal, ok := PlaneFromNormalPoint(V3(4, 0, 0), a)
	if !ok {
		t.Fatal("PlaneFromNormalPoint rejected finite normal")
	}
	for _, plane := range []Plane{fromPoints, fromNormal} {
		if plane != NewPlane(V3(1, 0, 0), -2) {
			t.Fatalf("plane = %#v, want x-2=0", plane)
		}
		if plane.DotCoordinate(a) != 0 || plane.ClassifyPoint(V3(3, 0, 0)) <= 0 || plane.ClassifyPoint(V3(1, 0, 0)) >= 0 {
			t.Fatalf("inconsistent classification for %#v", plane)
		}
	}
	if _, ok := PlaneFromPoints(a, a, c); ok {
		t.Fatal("degenerate points produced a plane")
	}
	if _, ok := PlaneFromNormalPoint(Vec3{}, a); ok {
		t.Fatal("zero normal produced a plane")
	}
	if XYPlane() != NewPlane(V3(0, 0, 1), 0) || XZPlane() != NewPlane(V3(0, 1, 0), 0) || YZPlane() != NewPlane(V3(1, 0, 0), 0) {
		t.Fatal("coordinate planes are incorrect")
	}
}

func TestPlaneNormalizationProjectionAndDots(t *testing.T) {
	p := NewPlane(V3(2, 0, 0), -4)
	normalized, ok := p.Normalized()
	if !ok || normalized != NewPlane(V3(1, 0, 0), -2) {
		t.Fatalf("Normalized = %#v, %v", normalized, ok)
	}
	projected, ok := p.Project(V3(7, 3, 4))
	if !ok || projected != V3(2, 3, 4) {
		t.Fatalf("Project = %#v, %v", projected, ok)
	}
	if p.Dot(V4(1, 2, 3, 4)) != -14 || p.DotNormal(V3(3, 4, 5)) != 6 {
		t.Fatal("plane dot operations are incorrect")
	}
	if _, ok := (Plane{}).Normalized(); ok {
		t.Fatal("zero plane normalization succeeded")
	}
	if _, ok := (Plane{}).Project(Vec3{}); ok {
		t.Fatal("zero plane projection succeeded")
	}
	if _, ok := NewPlane(V3(1, 0, 0), float32(math.Inf(1))).Normalized(); ok {
		t.Fatal("non-finite plane normalization succeeded")
	}
}

func TestPlaneMatrixAndQuaternionTransforms(t *testing.T) {
	p := NewPlane(V3(1, 0, 0), -2)
	transformed, ok := p.Transformed(TranslationMat4(V3(3, 4, 5)))
	if !ok || transformed != NewPlane(V3(1, 0, 0), -5) {
		t.Fatalf("translated plane = %#v, %v", transformed, ok)
	}
	if transformed.DotCoordinate(V3(5, 9, -1)) != 0 {
		t.Fatal("translated point is not on transformed plane")
	}
	rotated := p.Rotated(QuatRotationZ(HalfPi))
	if !rotated.Normal.AlmostEqual(V3(0, 1, 0), algorithmTolerance) || rotated.D != -2 {
		t.Fatalf("rotated plane = %#v", rotated)
	}
	if _, ok := p.Transformed(Mat4{}); ok {
		t.Fatal("plane transformation by singular matrix succeeded")
	}
}

func TestRigidTransformBehavior(t *testing.T) {
	transform := NewTransform(V3(10, 20, 30), QuatRotationZ(HalfPi))
	point := V3(2, 0, 1)
	direction := V3(2, 0, 1)
	if got := transform.TransformPoint(point); !got.AlmostEqual(V3(10, 22, 31), algorithmTolerance) {
		t.Fatalf("TransformPoint = %#v", got)
	}
	if got := transform.TransformDirection(direction); !got.AlmostEqual(V3(0, 2, 1), algorithmTolerance) {
		t.Fatalf("TransformDirection = %#v", got)
	}
	if got, want := transform.Mat4().TransformPoint(point), transform.TransformPoint(point); !got.AlmostEqual(want, algorithmTolerance) {
		t.Fatalf("Mat4 point = %#v, direct = %#v", got, want)
	}
	if got, want := transform.Mat4().TransformDirection(direction), transform.TransformDirection(direction); !got.AlmostEqual(want, algorithmTolerance) {
		t.Fatalf("Mat4 direction = %#v, direct = %#v", got, want)
	}
	if IdentityTransform().Mat4() != IdentityMat4() {
		t.Fatal("identity transform did not produce identity matrix")
	}
}
