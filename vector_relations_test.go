package math3d

import (
	"math"
	"testing"
)

func TestVec3ProjectionAndRejection(t *testing.T) {
	v := V3(3, 4, 5)
	onto := V3(2, 0, 0)
	projection, ok := v.Projection(onto)
	if !ok || projection != V3(3, 0, 0) {
		t.Fatalf("Projection() = %#v, %v; want (3,0,0), true", projection, ok)
	}
	rejection, ok := v.Rejection(onto)
	if !ok || rejection != V3(0, 4, 5) {
		t.Fatalf("Rejection() = %#v, %v; want (0,4,5), true", rejection, ok)
	}
	if got := projection.Add(rejection); got != v {
		t.Fatalf("projection + rejection = %#v, want %#v", got, v)
	}
	if projection.Dot(rejection) != 0 {
		t.Fatalf("projection dot rejection = %v, want 0", projection.Dot(rejection))
	}
}

func TestVectorReflection(t *testing.T) {
	for _, normal := range []Vec3{V3(0, 1, 0), V3(0, 3, 0)} {
		got, ok := V3(2, -1, 4).Reflect(normal)
		if !ok || got != V3(2, 1, 4) {
			t.Fatalf("Reflect(%#v) = %#v, %v; want (2,1,4), true", normal, got, ok)
		}
	}
	for _, normal := range []Vec2{V2(0, 1), V2(0, 4)} {
		got, ok := V2(2, -1).Reflect(normal)
		if !ok || got != V2(2, 1) {
			t.Fatalf("Vec2.Reflect(%#v) = %#v, %v; want (2,1), true", normal, got, ok)
		}
	}
}

func TestVectorRelationInvalidInputs(t *testing.T) {
	nan := float32(math.NaN())
	inf := float32(math.Inf(1))
	for name, operation := range map[string]func() bool{
		"projection onto zero":      func() bool { _, ok := V3(1, 2, 3).Projection(Vec3{}); return ok },
		"rejection from zero":       func() bool { _, ok := V3(1, 2, 3).Rejection(Vec3{}); return ok },
		"reflection around zero":    func() bool { _, ok := V3(1, 2, 3).Reflect(Vec3{}); return ok },
		"2D reflection around zero": func() bool { _, ok := V2(1, 2).Reflect(Vec2{}); return ok },
		"non-finite projection":     func() bool { _, ok := V3(nan, 0, 0).Projection(V3(1, 0, 0)); return ok },
		"non-finite reflection":     func() bool { _, ok := V3(1, 0, 0).Reflect(V3(inf, 0, 0)); return ok },
		"zero angle":                func() bool { _, ok := Vec3{}.Angle(V3(1, 0, 0)); return ok },
		"zero signed axis":          func() bool { _, ok := V3(1, 0, 0).SignedAngle(V3(0, 1, 0), Vec3{}); return ok },
	} {
		t.Run(name, func(t *testing.T) {
			if operation() {
				t.Fatal("operation unexpectedly succeeded")
			}
		})
	}

	if (Vec3{}).IsPerpendicular(V3(1, 0, 0), Tolerance) ||
		V3(nan, 0, 0).IsCollinear(V3(1, 0, 0), Tolerance) ||
		V3(1, 0, 0).IsBackFace(Vec3{}) ||
		Coplanar(Vec3{}, V3(1, 0, 0), V3(0, 1, 0), V3(1, 1, 0), nan) {
		t.Fatal("invalid predicate input must return false")
	}
}

func TestVec3PerpendicularTolerance(t *testing.T) {
	x := V3(1, 0, 0)
	if !x.IsPerpendicular(V3(Tolerance/2, 1, 0), Tolerance) {
		t.Fatal("dot product inside tolerance was not perpendicular")
	}
	if x.IsPerpendicular(V3(Tolerance, 1, 0), Tolerance) {
		t.Fatal("dot product at strict tolerance boundary was perpendicular")
	}
	if x.IsPerpendicular(V3(Tolerance*2, 1, 0), Tolerance) {
		t.Fatal("dot product outside tolerance was perpendicular")
	}
}

func TestVec3Collinearity(t *testing.T) {
	x := V3(1, 0, 0)
	if !x.IsCollinear(V3(3, 0, 0), 0) || !x.IsCollinear(V3(-3, 0, 0), 0) {
		t.Fatal("parallel and antiparallel vectors must be collinear")
	}
	if x.IsCollinear(V3(0, 1, 0), HalfPi/2) {
		t.Fatal("orthogonal vectors must not be collinear")
	}
	if !x.IsCollinear(V3(0, 1, 0), HalfPi) {
		t.Fatal("angle at the inclusive tolerance boundary must be collinear")
	}
	if x.IsCollinear(V3(1, 0, 1), Tolerance) {
		t.Fatal("non-collinear 3D vectors were reported collinear")
	}
}

func TestVec3Angles(t *testing.T) {
	x, y, z := V3(1, 0, 0), V3(0, 1, 0), V3(0, 0, 1)
	if got, ok := x.Angle(y); !ok || got != HalfPi {
		t.Fatalf("Angle() = %v, %v; want %v, true", got, ok, HalfPi)
	}
	if got, ok := x.SignedAngle(y, z); !ok || got != HalfPi {
		t.Fatalf("positive SignedAngle() = %v, %v; want %v, true", got, ok, HalfPi)
	}
	if got, ok := y.SignedAngle(x, z); !ok || got != -HalfPi {
		t.Fatalf("negative SignedAngle() = %v, %v; want %v, true", got, ok, -HalfPi)
	}
	if _, ok := x.SignedAngle(x.Negated(), z); ok {
		t.Fatal("antiparallel signed angle has ambiguous orientation")
	}
}

func TestCoplanar(t *testing.T) {
	a, b, c := Vec3{}, V3(1, 0, 0), V3(0, 1, 0)
	if !Coplanar(a, b, c, V3(2, 3, 0), Tolerance) {
		t.Fatal("points in the XY plane were not coplanar")
	}
	if Coplanar(a, b, c, V3(0, 0, 1), 1) {
		t.Fatal("scalar triple product at strict tolerance boundary was coplanar")
	}
	if !Coplanar(a, b, c, V3(0, 0, 1), math.Nextafter32(1, 2)) {
		t.Fatal("scalar triple product inside tolerance was not coplanar")
	}
	if !Coplanar(a, a, c, V3(0, 0, 1), Tolerance) {
		t.Fatal("degenerate point set must be coplanar")
	}
}

func TestBackFace(t *testing.T) {
	normal := V3(0, 0, 2)
	if !normal.IsBackFace(V3(0, 0, -3)) {
		t.Fatal("opposing line of sight was not a back face")
	}
	if normal.IsBackFace(V3(0, 0, 3)) || normal.IsBackFace(V3(1, 0, 0)) {
		t.Fatal("front-facing or edge-on normal was a back face")
	}
}

func TestVectorRelationsDoNotMutateInputs(t *testing.T) {
	v, other := V3(1, 2, 3), V3(4, 5, 6)
	wantV, wantOther := v, other
	_, _ = v.Projection(other)
	_, _ = v.Rejection(other)
	_, _ = v.Reflect(other)
	_ = v.IsPerpendicular(other, Tolerance)
	_, _ = v.Angle(other)
	_, _ = v.SignedAngle(other, V3(0, 0, 1))
	_ = v.IsCollinear(other, Tolerance)
	_ = v.IsBackFace(other)
	_ = Coplanar(v, other, Vec3{}, V3(1, 0, 0), Tolerance)
	if v != wantV || other != wantOther {
		t.Fatalf("inputs changed: %#v, %#v", v, other)
	}
}
