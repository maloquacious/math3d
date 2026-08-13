package math3d

import "testing"

func TestSegment3MappingAndGeometry(t *testing.T) {
	segment := NewSegment3(V3(1, 2, 3), V3(4, 6, 3))
	if segment.Vector() != V3(3, 4, 0) || segment.Length() != 5 || segment.LengthSquared() != 25 {
		t.Fatalf("unexpected segment geometry: %#v", segment)
	}
	if segment.Midpoint() != V3(2.5, 4, 3) || segment.PointAt(2) != V3(7, 10, 3) {
		t.Fatal("segment mapping or extrapolation is incorrect")
	}
	if got, ok := segment.WithLength(10); !ok || got.B != V3(7, 10, 3) {
		t.Fatalf("WithLength = %#v, %v", got, ok)
	}
	if _, ok := NewSegment3(Vec3{}, Vec3{}).WithLength(1); ok {
		t.Fatal("degenerate segment accepted a new length")
	}
	translation := TranslationMat4(V3(10, -2, 1))
	if got := segment.Transformed(translation); got != NewSegment3(V3(11, 0, 4), V3(14, 4, 4)) {
		t.Fatalf("Transformed = %#v", got)
	}
}

func TestSegment2Intersection(t *testing.T) {
	tests := []struct {
		name       string
		a, b       Segment2
		intersects bool
	}{
		{"crossing", NewSegment2(V2(0, 0), V2(2, 2)), NewSegment2(V2(0, 2), V2(2, 0)), true},
		{"shared endpoint", NewSegment2(V2(0, 0), V2(1, 0)), NewSegment2(V2(1, 0), V2(1, 1)), true},
		{"collinear overlap", NewSegment2(V2(0, 0), V2(3, 0)), NewSegment2(V2(2, 0), V2(4, 0)), true},
		{"collinear disjoint", NewSegment2(V2(0, 0), V2(1, 0)), NewSegment2(V2(2, 0), V2(3, 0)), false},
		{"separated", NewSegment2(V2(0, 0), V2(1, 0)), NewSegment2(V2(0, 1), V2(1, 1)), false},
		{"point on segment", NewSegment2(V2(1, 0), V2(1, 0)), NewSegment2(V2(0, 0), V2(2, 0)), true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.a.Intersects(test.b, Tolerance); got != test.intersects {
				t.Fatalf("Intersects = %v, want %v", got, test.intersects)
			}
			if got := test.b.Intersects(test.a, Tolerance); got != test.intersects {
				t.Fatalf("reverse Intersects = %v, want %v", got, test.intersects)
			}
		})
	}
}

func TestTriangleGeometryAndTransform(t *testing.T) {
	triangle := NewTriangle3(Vec3{}, V3(4, 0, 0), V3(0, 3, 0))
	if triangle.Area() != 6 || triangle.Perimeter() != 12 || triangle.Centroid() != V3(4.0/3.0, 1, 0) {
		t.Fatalf("unexpected triangle measurements: area=%v perimeter=%v centroid=%#v", triangle.Area(), triangle.Perimeter(), triangle.Centroid())
	}
	if normal, ok := triangle.Normal(); !ok || normal != V3(0, 0, 1) {
		t.Fatalf("Normal = %#v, %v", normal, ok)
	}
	if triangle.Degenerate() || triangle.IsSliver(1) || triangle.Bounds() != NewBox3(Vec3{}, V3(4, 3, 0)) {
		t.Fatal("triangle validity or bounds are incorrect")
	}
	collinear := NewTriangle3(Vec3{}, V3(1, 0, 0), V3(2, 0, 0))
	if !collinear.Degenerate() {
		t.Fatal("distinct collinear points were not degenerate")
	}
	if _, ok := collinear.Normal(); ok {
		t.Fatal("degenerate triangle produced a normal")
	}
	transformed := triangle.Transformed(TranslationMat4(V3(10, 20, 30)))
	if transformed != NewTriangle3(V3(10, 20, 30), V3(14, 20, 30), V3(10, 23, 30)) {
		t.Fatalf("Transformed = %#v", transformed)
	}
}

func TestTriangle2AreaAndContainment(t *testing.T) {
	triangle := NewTriangle2(Vec2{}, V2(2, 0), V2(0, 2))
	if triangle.SignedArea() != 2 || triangle.Area() != 2 {
		t.Fatalf("areas = %v, %v", triangle.SignedArea(), triangle.Area())
	}
	if !triangle.Contains(V2(0.5, 0.5)) || triangle.Contains(V2(1, 1)) || triangle.Contains(V2(2, 2)) {
		t.Fatal("strict triangle containment is incorrect")
	}
}

func TestRayTriangleIntersection(t *testing.T) {
	triangle := NewTriangle3(V3(-1, -1, 5), V3(1, -1, 5), V3(0, 1, 5))
	for _, direction := range []Vec3{V3(0, 0, 2), V3(0, 0, -2)} {
		origin := Vec3{}
		if direction.Z < 0 {
			origin = V3(0, 0, 10)
		}
		if got, ok := NewRay(origin, direction).IntersectTriangle(triangle, Tolerance); !ok || got != 2.5 {
			t.Fatalf("IntersectTriangle = %v, %v", got, ok)
		}
	}
	if _, ok := NewRay(V3(2, 0, 0), V3(0, 0, 1)).IntersectTriangle(triangle, Tolerance); ok {
		t.Fatal("ray outside triangle hit")
	}
	if _, ok := NewRay(V3(0, 0, 5), V3(0, 0, 1)).IntersectTriangle(triangle, Tolerance); ok {
		t.Fatal("hit at ray origin was accepted")
	}
	if _, ok := NewRay(Vec3{}, V3(1, 0, 0)).IntersectTriangle(triangle, Tolerance); ok {
		t.Fatal("parallel ray hit triangle")
	}
	if got, ok := NewRay(Vec3{}, V3(0, 0, 1)).IntersectTriangle(NewTriangle3(triangle.C, triangle.B, triangle.A), Tolerance); !ok || got != 5 {
		t.Fatalf("reversed-winding IntersectTriangle = %v, %v", got, ok)
	}
}

func TestQuadMappingAndTransform(t *testing.T) {
	quad := NewQuad3(Vec3{}, V3(1, 0, 0), V3(1, 1, 0), V3(0, 1, 0))
	got := quad.Map(func(v Vec3) Vec3 { return v.Scale(2) }).Transformed(TranslationMat4(V3(3, 4, 5)))
	want := NewQuad3(V3(3, 4, 5), V3(5, 4, 5), V3(5, 6, 5), V3(3, 6, 5))
	if got != want {
		t.Fatalf("mapped and transformed quad = %#v", got)
	}
}
