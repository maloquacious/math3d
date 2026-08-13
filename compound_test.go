package math3d

import (
	"math"
	"reflect"
	"testing"
)

func TestCompoundConstructorsPreserveRepresentation(t *testing.T) {
	v2a, v2b := V2(1, 2), V2(3, 4)
	v3a, v3b := V3(1, 2, 3), V3(4, 5, 6)
	v4a, v4b := V4(1, 2, 3, 4), V4(5, 6, 7, 8)
	dv2a, dv2b := DV2(1, 2), DV2(3, 4)
	dv3a, dv3b := DV3(1, 2, 3), DV3(4, 5, 6)
	dv4a, dv4b := DV4(1, 2, 3, 4), DV4(5, 6, 7, 8)

	tests := []struct {
		name      string
		got, want any
	}{
		{"Quat", NewQuat(1, 2, 3, 4), Quat{X: 1, Y: 2, Z: 3, W: 4}},
		{"DQuat", NewDQuat(1, 2, 3, 4), DQuat{X: 1, Y: 2, Z: 3, W: 4}},
		{"AxisAngle", NewAxisAngle(dv3a, 4), AxisAngle{Axis: dv3a, Angle: 4}},
		{"Euler", NewEuler(1, 2, 3), Euler{Yaw: 1, Pitch: 2, Roll: 3}},
		{"Transform", NewTransform(v3a, NewQuat(4, 5, 6, 7)), Transform{Position: v3a, Orientation: NewQuat(4, 5, 6, 7)}},
		{"Plane", NewPlane(v3a, 4), Plane{Normal: v3a, D: 4}},
		{"DPlane", NewDPlane(dv3a, 4), DPlane{Normal: dv3a, D: 4}},
		{"Interval", NewInterval(2, -1), Interval{Min: 2, Max: -1}},
		{"DInterval", NewDInterval(2, -1), DInterval{Min: 2, Max: -1}},
		{"Box2", NewBox2(v2b, v2a), Box2{Min: v2b, Max: v2a}},
		{"Box3", NewBox3(v3b, v3a), Box3{Min: v3b, Max: v3a}},
		{"Box4", NewBox4(v4b, v4a), Box4{Min: v4b, Max: v4a}},
		{"DBox2", NewDBox2(dv2b, dv2a), DBox2{Min: dv2b, Max: dv2a}},
		{"DBox3", NewDBox3(dv3b, dv3a), DBox3{Min: dv3b, Max: dv3a}},
		{"DBox4", NewDBox4(dv4b, dv4a), DBox4{Min: dv4b, Max: dv4a}},
		{"Ray", NewRay(v3a, v3b), Ray{Origin: v3a, Direction: v3b}},
		{"DRay", NewDRay(dv3a, dv3b), DRay{Origin: dv3a, Direction: dv3b}},
		{"Sphere", NewSphere(v3a, -2), Sphere{Center: v3a, Radius: -2}},
		{"DSphere", NewDSphere(dv3a, -2), DSphere{Center: dv3a, Radius: -2}},
		{"Segment3", NewSegment3(v3a, v3b), Segment3{A: v3a, B: v3b}},
		{"Segment2", NewSegment2(v2a, v2b), Segment2{A: v2a, B: v2b}},
		{"Triangle3", NewTriangle3(v3a, v3b, Vec3{}), Triangle3{A: v3a, B: v3b, C: Vec3{}}},
		{"Triangle2", NewTriangle2(v2a, v2b, Vec2{}), Triangle2{A: v2a, B: v2b, C: Vec2{}}},
		{"Quad3", NewQuad3(v3a, v3b, Vec3{}, V3(7, 8, 9)), Quad3{A: v3a, B: v3b, C: Vec3{}, D: V3(7, 8, 9)}},
		{"Quad2", NewQuad2(v2a, v2b, Vec2{}, V2(5, 6)), Quad2{A: v2a, B: v2b, C: Vec2{}, D: V2(5, 6)}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !reflect.DeepEqual(test.got, test.want) {
				t.Fatalf("constructor = %#v, want %#v", test.got, test.want)
			}
		})
	}
}

func TestMat4Representation(t *testing.T) {
	m := NewMat4(
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	)
	if got, want := m.Row(2), V4(9, 10, 11, 12); got != want {
		t.Fatalf("Row(2) = %#v, want %#v", got, want)
	}
	if got, want := m.Column(1), V4(2, 6, 10, 14); got != want {
		t.Fatalf("Column(1) = %#v, want %#v", got, want)
	}
	for row := 0; row < 4; row++ {
		for column := 0; column < 4; column++ {
			want := float32(row*4 + column + 1)
			if got := m.At(row, column); got != want {
				t.Fatalf("At(%d, %d) = %v, want %v", row, column, got, want)
			}
		}
	}
	if got := m.Translation(); got != V3(13, 14, 15) {
		t.Fatalf("Translation() = %#v", got)
	}
	updated := m.WithTranslation(V3(20, 21, 22))
	if updated.Translation() != V3(20, 21, 22) || m.Translation() != V3(13, 14, 15) {
		t.Fatal("WithTranslation changed the input or stored the wrong translation")
	}
	if !IdentityMat4().IsIdentity() || (Mat4{}).IsIdentity() {
		t.Fatal("matrix identity values are incorrect")
	}
	if got := Mat4FromRows(m.Row(0), m.Row(1), m.Row(2), m.Row(3)); got != m {
		t.Fatalf("Mat4FromRows round trip = %#v", got)
	}
}

func TestMat4AccessPanicsForInvalidIndexes(t *testing.T) {
	tests := []struct {
		name string
		call func()
	}{
		{"At negative row", func() { Mat4{}.At(-1, 1) }},
		{"At large column", func() { Mat4{}.At(0, 4) }},
		{"Row", func() { Mat4{}.Row(4) }},
		{"Column", func() { Mat4{}.Column(-1) }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("call did not panic")
				}
			}()
			test.call()
		})
	}
}

func TestCompoundIdentityAndConversions(t *testing.T) {
	if got, want := IdentityQuat(), NewQuat(0, 0, 0, 1); got != want {
		t.Fatalf("IdentityQuat() = %#v, want %#v", got, want)
	}
	if got, want := IdentityTransform(), NewTransform(Vec3{}, IdentityQuat()); got != want {
		t.Fatalf("IdentityTransform() = %#v, want %#v", got, want)
	}
	if got, want := QuatFromVector(V3(1, 2, 3), 4).Vector(), V3(1, 2, 3); got != want {
		t.Fatalf("quaternion vector round trip = %#v, want %#v", got, want)
	}
	if got, want := NewQuat(1, 2, 3, 4).DQuat().Quat(), NewQuat(1, 2, 3, 4); got != want {
		t.Fatalf("quaternion precision round trip = %#v, want %#v", got, want)
	}
	if got, want := PlaneFromVec4(V4(1, 2, 3, 4)).Vec4(), V4(1, 2, 3, 4); got != want {
		t.Fatalf("plane vector round trip = %#v, want %#v", got, want)
	}
	box := NewBox3(V3(1, 2, 3), V3(4, 5, 6))
	if got := box.DBox3().Box3(); got != box {
		t.Fatalf("box precision round trip = %#v, want %#v", got, box)
	}
}

func TestCompoundAlmostEqualUsesRepresentation(t *testing.T) {
	q := NewQuat(1, 2, 3, 4)
	if !q.AlmostEqual(NewQuat(1.01, 2, 3, 4), 0.02) {
		t.Fatal("near quaternions should compare approximately equal")
	}
	if q.AlmostEqual(NewQuat(-1, -2, -3, -4), 1) {
		t.Fatal("opposite quaternion representations must not be specially equated")
	}
	if q.AlmostEqual(NewQuat(1.25, 2, 3, 4), 0.25) {
		t.Fatal("tolerance boundary must be strict")
	}
	if NewSphere(V3(float32(math.NaN()), 0, 0), 1).AlmostEqual(NewSphere(V3(float32(math.NaN()), 0, 0), 1), 1) {
		t.Fatal("NaN must not compare approximately equal")
	}
}

func TestCompoundValuesRemainComparable(t *testing.T) {
	assertComparable := func(values ...any) {
		for _, value := range values {
			if !reflect.TypeOf(value).Comparable() {
				t.Errorf("%T is not comparable", value)
			}
		}
	}
	assertComparable(
		Quat{}, DQuat{}, AxisAngle{}, Euler{}, Mat4{}, Transform{}, Plane{}, DPlane{},
		Interval{}, DInterval{}, Box2{}, Box3{}, Box4{}, DBox2{}, DBox3{}, DBox4{},
		Ray{}, DRay{}, Sphere{}, DSphere{}, Segment2{}, Segment3{}, Triangle2{}, Triangle3{}, Quad2{}, Quad3{},
	)
}
