package math3d

import "testing"

func TestIntervalValidityMergeAndIntersection(t *testing.T) {
	a := NewInterval(-2, 4)
	b := NewInterval(3, 8)
	if !a.Valid() || a.Empty() || a.Extent() != 6 || a.Center() != 1 {
		t.Fatalf("unexpected interval properties: %#v", a)
	}
	if got := a.Merge(b); got != NewInterval(-2, 8) {
		t.Fatalf("Merge = %#v", got)
	}
	if got, ok := a.Intersection(b); !ok || got != NewInterval(3, 4) {
		t.Fatalf("Intersection = %#v, %v", got, ok)
	}
	if got, ok := a.Intersection(NewInterval(4, 9)); !ok || got != NewInterval(4, 4) {
		t.Fatalf("touching intersection = %#v, %v", got, ok)
	}
	if got, ok := a.Intersection(NewInterval(5, 9)); ok || !got.Empty() {
		t.Fatalf("disjoint intersection = %#v, %v", got, ok)
	}
	if !EmptyInterval().Empty() || EmptyInterval().Merge(a) != a {
		t.Fatal("canonical empty interval is not a merge identity")
	}
	if got, ok := NewDInterval(0, 2).Intersection(NewDInterval(1, 3)); !ok || got != NewDInterval(1, 2) {
		t.Fatalf("DInterval intersection = %#v, %v", got, ok)
	}
}

func TestBoxValidityContainmentAndIntersection(t *testing.T) {
	outer := NewBox2(V2(0, 0), V2(4, 4))
	inner := NewBox2(V2(1, 1), V2(2, 2))
	touching := NewBox2(V2(4, 4), V2(5, 5))
	if !outer.Valid() || outer.Empty() || outer.Center() != V2(2, 2) || outer.Extent() != V2(4, 4) {
		t.Fatalf("unexpected Box2 properties: %#v", outer)
	}
	if outer.ContainsBox(inner) != ContainmentContains || outer.ContainsBox(touching) != ContainmentIntersects {
		t.Fatal("Box2 containment did not include boundaries")
	}
	if !outer.Intersects(touching) || outer.ContainsPoint(V2(4, 2)) != true {
		t.Fatal("Box2 boundary contact did not intersect")
	}
	if _, ok := outer.Intersection(NewBox2(V2(5, 5), V2(6, 6))); ok {
		t.Fatal("disjoint boxes intersected")
	}
	if !EmptyBox4().Empty() || !EmptyDBox2().Empty() || !EmptyDBox4().Empty() {
		t.Fatal("canonical empty boxes are valid")
	}
}

func TestBox3CornersAndTransform(t *testing.T) {
	box := NewBox3(V3(0, 0, 0), V3(2, 4, 6))
	corners := box.Corners()
	for _, corner := range corners {
		if !box.ContainsPoint(corner) {
			t.Fatalf("corner %#v is outside box", corner)
		}
	}
	matrix := RotationZMat4(HalfPi).Mul(TranslationMat4(V3(10, 0, 0)))
	got, ok := box.Transformed(matrix)
	want := NewBox3(V3(6, 0, 0), V3(10, 2, 6))
	if !ok || !got.AlmostEqual(want, algorithmTolerance) {
		t.Fatalf("Transformed = %#v, %v; want %#v", got, ok, want)
	}
	fromPoints, ok := Box3FromPoints([]Vec3{V3(2, 3, 4), V3(-1, 5, 0)})
	if !ok || fromPoints != NewBox3(V3(-1, 3, 0), V3(2, 5, 4)) {
		t.Fatalf("Box3FromPoints = %#v, %v", fromPoints, ok)
	}
	if _, ok := Box3FromPoints(nil); ok {
		t.Fatal("empty point set produced bounds")
	}
}

func TestRayBoxAndPlaneIntersections(t *testing.T) {
	box := NewBox3(V3(-1, -1, -1), V3(1, 1, 1))
	ray := NewRay(V3(-5, 0, 0), V3(2, 0, 0))
	if got, ok := ray.IntersectBox(box); !ok || got != 2 || ray.PointAt(got) != V3(-1, 0, 0) {
		t.Fatalf("IntersectBox = %v, %v", got, ok)
	}
	if got, ok := NewRay(Vec3{}, V3(0, 1, 0)).IntersectBox(box); !ok || got != 0 {
		t.Fatalf("inside IntersectBox = %v, %v", got, ok)
	}
	if _, ok := NewRay(V3(2, 0, 0), V3(0, 1, 0)).IntersectBox(box); ok {
		t.Fatal("parallel outside ray hit box")
	}
	if _, ok := NewRay(Vec3{}, Vec3{}).IntersectBox(box); ok {
		t.Fatal("zero-direction ray hit box")
	}
	plane := NewPlane(V3(2, 0, 0), -2)
	if got, ok := NewRay(Vec3{}, V3(4, 0, 0)).IntersectPlane(plane, Tolerance); !ok || got != 0.25 {
		t.Fatalf("IntersectPlane = %v, %v", got, ok)
	}
	if _, ok := NewRay(Vec3{}, V3(0, 1, 0)).IntersectPlane(plane, Tolerance); ok {
		t.Fatal("parallel ray hit plane")
	}
}

func TestRaySphereSupportsNonUnitDirections(t *testing.T) {
	sphere := NewSphere(Vec3{}, 1)
	ray := NewRay(V3(-5, 0, 0), V3(2, 0, 0))
	if got, ok := ray.IntersectSphere(sphere); !ok || got != 2 || ray.PointAt(got) != V3(-1, 0, 0) {
		t.Fatalf("IntersectSphere = %v, %v", got, ok)
	}
	if got, ok := NewRay(Vec3{}, V3(10, 0, 0)).IntersectSphere(sphere); !ok || got != 0 {
		t.Fatalf("inside sphere = %v, %v", got, ok)
	}
	if _, ok := NewRay(V3(-5, 0, 0), V3(-2, 0, 0)).IntersectSphere(sphere); ok {
		t.Fatal("ray pointing away hit sphere")
	}
	dray := NewDRay(DV3(-5, 0, 0), DV3(2, 0, 0))
	if got, ok := dray.IntersectSphere(NewDSphere(DVec3{}, 1)); !ok || got != 2 {
		t.Fatalf("DRay.IntersectSphere = %v, %v", got, ok)
	}
}

func TestSphereContainmentMergeAndTransform(t *testing.T) {
	sphere := NewSphere(Vec3{}, 2)
	if sphere.ContainsPoint(Vec3{}) != ContainmentContains || sphere.ContainsPoint(V3(2, 0, 0)) != ContainmentIntersects {
		t.Fatal("sphere point containment is incorrect")
	}
	if sphere.ContainsSphere(NewSphere(Vec3{}, 3)) != ContainmentIntersects {
		t.Fatal("smaller sphere incorrectly contains larger sphere")
	}
	box := NewBox3(V3(-0.5, -0.5, -0.5), V3(0.5, 0.5, 0.5))
	if sphere.ContainsBox(box) != ContainmentContains || box.ContainsSphere(NewSphere(Vec3{}, 0.25)) != ContainmentContains {
		t.Fatal("sphere/box containment is incorrect")
	}
	if !NewSphere(V3(0.9, 0.9, 0), 0.2).IntersectsBox(NewBox3(V3(-1, -1, -1), V3(1, 1, 1))) {
		t.Fatal("sphere centered inside a box did not intersect it")
	}
	merged, ok := NewSphere(V3(0, 0, 0), 1).Merge(NewSphere(V3(4, 0, 0), 1))
	if !ok || merged != NewSphere(V3(2, 0, 0), 3) {
		t.Fatalf("Merge = %#v, %v", merged, ok)
	}
	transformed, ok := sphere.Transformed(ScaleMat4(V3(1, 3, 2)).Mul(TranslationMat4(V3(10, 0, 0))))
	if !ok || transformed != NewSphere(V3(10, 0, 0), 6) {
		t.Fatalf("Transformed = %#v, %v", transformed, ok)
	}
	rotated, ok := sphere.Transformed(RotationZMat4(HalfPi))
	if !ok || !rotated.AlmostEqual(sphere, algorithmTolerance) {
		t.Fatalf("rotation changed sphere: %#v, %v", rotated, ok)
	}
	shear := IdentityMat4()
	shear.M12 = 1
	sheared, ok := NewSphere(Vec3{}, 1).Transformed(shear)
	if !ok || sheared.Radius < 1.618 || sheared.Radius > 1.619 {
		t.Fatalf("shear enclosure radius = %v, %v", sheared.Radius, ok)
	}
	projective := IdentityMat4()
	projective.M14 = 1
	if _, ok := sphere.Transformed(projective); ok {
		t.Fatal("projective matrix produced a transformed sphere")
	}
}

func TestBoxAndSpherePlaneClassification(t *testing.T) {
	plane := NewPlane(V3(2, 0, 0), -4)
	if got := NewBox3(V3(3, -1, -1), V3(4, 1, 1)).IntersectsPlane(plane); got != PlaneIntersectionFront {
		t.Fatalf("front box = %v", got)
	}
	if got := NewBox3(V3(1, -1, -1), V3(3, 1, 1)).IntersectsPlane(plane); got != PlaneIntersectionIntersecting {
		t.Fatalf("crossing box = %v", got)
	}
	if got := NewSphere(V3(4, 0, 0), 1).IntersectsPlane(plane); got != PlaneIntersectionFront {
		t.Fatalf("front sphere = %v", got)
	}
	if got := NewSphere(V3(2.5, 0, 0), 1).IntersectsPlane(plane); got != PlaneIntersectionIntersecting {
		t.Fatalf("crossing sphere = %v", got)
	}
}
