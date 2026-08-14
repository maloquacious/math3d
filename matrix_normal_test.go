package math3d

import (
	"math"
	"testing"
)

func TestMatrixTransformNormalBasicTransforms(t *testing.T) {
	normal := V3(1, 2, 3)
	tests := []struct {
		name   string
		matrix Mat4
		want   Vec3
	}{
		{"identity", IdentityMat4(), normal},
		{"translation", TranslationMat4(V3(10, 20, 30)), normal},
		{"rotation", RotationZMat4(HalfPi), V3(-2, 1, 3)},
		{"non-uniform scale", ScaleMat4(V3(2, 4, 8)), V3(0.5, 0.5, 0.375)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := test.matrix.TransformNormal(normal)
			if !ok || !got.AlmostEqual(test.want, algorithmTolerance) {
				t.Fatalf("TransformNormal = %#v, %v; want %#v", got, ok, test.want)
			}
		})
	}
	if normal != V3(1, 2, 3) {
		t.Fatalf("TransformNormal mutated its input: %#v", normal)
	}
}

func TestMatrixTransformNormalUsesRowVectorInverseTranspose(t *testing.T) {
	// In row-vector form this maps the X tangent from (1,0,0) to (1,2,0).
	shear := IdentityMat4()
	shear.M12 = 2
	normal, ok := shear.TransformNormal(V3(0, 1, 0))
	if !ok || !normal.AlmostEqual(V3(-2, 1, 0), algorithmTolerance) {
		t.Fatalf("sheared normal = %#v, %v", normal, ok)
	}

	tangentX := shear.TransformDirection(V3(1, 0, 0))
	tangentZ := shear.TransformDirection(V3(0, 0, 1))
	if abs32(normal.Dot(tangentX)) > algorithmTolerance || abs32(normal.Dot(tangentZ)) > algorithmTolerance {
		t.Fatalf("normal %#v is not perpendicular to transformed tangents %#v and %#v", normal, tangentX, tangentZ)
	}
	if shear.TransformDirection(V3(0, 1, 0)).AlmostEqual(normal, algorithmTolerance) {
		t.Fatal("normal transform unexpectedly matched direction transform under shear")
	}
}

func TestMatrixTransformNormalFailures(t *testing.T) {
	if got, ok := ScaleMat4(V3(1, 0, 1)).TransformNormal(V3(0, 1, 0)); ok || got != (Vec3{}) {
		t.Fatalf("singular TransformNormal = %#v, %v", got, ok)
	}
	if got, ok := IdentityMat4().TransformNormal(V3(float32(math.NaN()), 0, 0)); ok || got != (Vec3{}) {
		t.Fatalf("non-finite normal TransformNormal = %#v, %v", got, ok)
	}
	nonFinite := IdentityMat4()
	nonFinite.M11 = float32(math.Inf(1))
	if got, ok := nonFinite.TransformNormal(V3(1, 0, 0)); ok || got != (Vec3{}) {
		t.Fatalf("non-finite matrix TransformNormal = %#v, %v", got, ok)
	}
}
