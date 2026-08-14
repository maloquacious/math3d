package math3d

import "math"

// Interval stores inclusive scalar bounds. Construction preserves bound order;
// Min > Max is the empty representation. The zero value is [0, 0].
type Interval struct {
	Min, Max float32
}

// NewInterval constructs an Interval without reordering its bounds.
func NewInterval(minimum, maximum float32) Interval {
	return Interval{Min: minimum, Max: maximum}
}

// AlmostEqual compares corresponding bounds using a strict absolute tolerance.
func (i Interval) AlmostEqual(other Interval, tolerance float32) bool {
	return abs32(i.Min-other.Min) < tolerance && abs32(i.Max-other.Max) < tolerance
}

func (i Interval) Valid() bool     { return i.Min <= i.Max }
func (i Interval) Empty() bool     { return !i.Valid() }
func (i Interval) Extent() float32 { return i.Max - i.Min }
func (i Interval) Center() float32 { return (i.Min + i.Max) / 2 }
func (i Interval) Merge(other Interval) Interval {
	return NewInterval(min(i.Min, other.Min), max(i.Max, other.Max))
}
func (i Interval) Expanded(value float32) Interval {
	return NewInterval(min(i.Min, value), max(i.Max, value))
}
func (i Interval) Intersection(other Interval) (Interval, bool) {
	result := NewInterval(max(i.Min, other.Min), min(i.Max, other.Max))
	return result, result.Valid()
}

// EmptyInterval returns the canonical empty float32 interval.
func EmptyInterval() Interval {
	return NewInterval(float32(math.Inf(1)), float32(math.Inf(-1)))
}

// DInterval stores inclusive float64 scalar bounds. Construction preserves
// bound order; Min > Max is the empty representation. The zero value is [0, 0].
type DInterval struct {
	Min, Max float64
}

// NewDInterval constructs a DInterval without reordering its bounds.
func NewDInterval(minimum, maximum float64) DInterval {
	return DInterval{Min: minimum, Max: maximum}
}

// AlmostEqual compares corresponding bounds using a strict absolute tolerance.
func (i DInterval) AlmostEqual(other DInterval, tolerance float64) bool {
	return almostEqual64(tolerance, i.Min-other.Min, i.Max-other.Max)
}

func (i DInterval) Valid() bool     { return i.Min <= i.Max }
func (i DInterval) Empty() bool     { return !i.Valid() }
func (i DInterval) Extent() float64 { return i.Max - i.Min }
func (i DInterval) Center() float64 { return (i.Min + i.Max) / 2 }
func (i DInterval) Merge(other DInterval) DInterval {
	return NewDInterval(min(i.Min, other.Min), max(i.Max, other.Max))
}
func (i DInterval) Expanded(value float64) DInterval {
	return NewDInterval(min(i.Min, value), max(i.Max, value))
}
func (i DInterval) Intersection(other DInterval) (DInterval, bool) {
	result := NewDInterval(max(i.Min, other.Min), min(i.Max, other.Max))
	return result, result.Valid()
}

// EmptyDInterval returns the canonical empty float64 interval.
func EmptyDInterval() DInterval { return NewDInterval(math.Inf(1), math.Inf(-1)) }

// Box2 is the upstream AABox2D value. Construction preserves endpoint order;
// any Min component greater than its corresponding Max component represents an
// empty box. The zero value is a point box at the origin.
type Box2 struct {
	Min, Max Vec2
}

// NewBox2 constructs a Box2 without reordering endpoints.
func NewBox2(minimum, maximum Vec2) Box2 { return Box2{Min: minimum, Max: maximum} }

// AlmostEqual compares corresponding endpoints.
func (b Box2) AlmostEqual(other Box2, tolerance float32) bool {
	return b.Min.AlmostEqual(other.Min, tolerance) && b.Max.AlmostEqual(other.Max, tolerance)
}

func (b Box2) Valid() bool {
	return b.Min.X <= b.Max.X && b.Min.Y <= b.Max.Y
}
func (b Box2) Empty() bool              { return !b.Valid() }
func (b Box2) Extent() Vec2             { return b.Max.Sub(b.Min) }
func (b Box2) Center() Vec2             { return b.Min.Add(b.Max).Scale(0.5) }
func (b Box2) Merge(other Box2) Box2    { return NewBox2(b.Min.Min(other.Min), b.Max.Max(other.Max)) }
func (b Box2) Expanded(point Vec2) Box2 { return NewBox2(b.Min.Min(point), b.Max.Max(point)) }
func (b Box2) Intersection(other Box2) (Box2, bool) {
	result := NewBox2(b.Min.Max(other.Min), b.Max.Min(other.Max))
	return result, result.Valid()
}
func (b Box2) ContainsPoint(point Vec2) bool {
	return b.Valid() && point.X >= b.Min.X && point.X <= b.Max.X && point.Y >= b.Min.Y && point.Y <= b.Max.Y
}
func (b Box2) ContainsBox(other Box2) Containment {
	if !b.Intersects(other) {
		return ContainmentDisjoint
	}
	if b.ContainsPoint(other.Min) && b.ContainsPoint(other.Max) {
		return ContainmentContains
	}
	return ContainmentIntersects
}
func (b Box2) Intersects(other Box2) bool {
	return b.Valid() && other.Valid() && b.Min.X <= other.Max.X && b.Max.X >= other.Min.X &&
		b.Min.Y <= other.Max.Y && b.Max.Y >= other.Min.Y
}
func (b Box2) Translated(offset Vec2) Box2 { return NewBox2(b.Min.Add(offset), b.Max.Add(offset)) }
func (b Box2) Corners() [4]Vec2 {
	return [4]Vec2{b.Min, V2(b.Max.X, b.Min.Y), b.Max, V2(b.Min.X, b.Max.Y)}
}

// EmptyBox2 returns the canonical empty 2D box.
func EmptyBox2() Box2 { return NewBox2(SplatV2(float32(math.Inf(1))), SplatV2(float32(math.Inf(-1)))) }

// Box2FromPoints returns the bounds of points, or false for an empty or non-finite input.
func Box2FromPoints(points []Vec2) (Box2, bool) {
	if len(points) == 0 || !points[0].IsFinite() {
		return Box2{}, false
	}
	result := NewBox2(points[0], points[0])
	for _, point := range points[1:] {
		if !point.IsFinite() {
			return Box2{}, false
		}
		result = result.Expanded(point)
	}
	return result, true
}

// Box3 is the upstream AABox value. Construction preserves endpoint order;
// any reversed component represents an empty box. The zero value is a point
// box at the origin.
type Box3 struct {
	Min, Max Vec3
}

// NewBox3 constructs a Box3 without reordering endpoints.
func NewBox3(minimum, maximum Vec3) Box3 { return Box3{Min: minimum, Max: maximum} }

// AlmostEqual compares corresponding endpoints.
func (b Box3) AlmostEqual(other Box3, tolerance float32) bool {
	return b.Min.AlmostEqual(other.Min, tolerance) && b.Max.AlmostEqual(other.Max, tolerance)
}

func (b Box3) Valid() bool {
	return b.Min.X <= b.Max.X && b.Min.Y <= b.Max.Y && b.Min.Z <= b.Max.Z
}
func (b Box3) Empty() bool              { return !b.Valid() }
func (b Box3) Extent() Vec3             { return b.Max.Sub(b.Min) }
func (b Box3) Center() Vec3             { return b.Min.Add(b.Max).Scale(0.5) }
func (b Box3) Merge(other Box3) Box3    { return NewBox3(b.Min.Min(other.Min), b.Max.Max(other.Max)) }
func (b Box3) Expanded(point Vec3) Box3 { return NewBox3(b.Min.Min(point), b.Max.Max(point)) }
func (b Box3) Intersection(other Box3) (Box3, bool) {
	result := NewBox3(b.Min.Max(other.Min), b.Max.Min(other.Max))
	return result, result.Valid()
}
func (b Box3) ContainsPoint(point Vec3) bool {
	return b.Valid() && point.X >= b.Min.X && point.X <= b.Max.X && point.Y >= b.Min.Y && point.Y <= b.Max.Y &&
		point.Z >= b.Min.Z && point.Z <= b.Max.Z
}
func (b Box3) ContainsBox(other Box3) Containment {
	if !b.Intersects(other) {
		return ContainmentDisjoint
	}
	if b.ContainsPoint(other.Min) && b.ContainsPoint(other.Max) {
		return ContainmentContains
	}
	return ContainmentIntersects
}
func (b Box3) Intersects(other Box3) bool {
	return b.Valid() && other.Valid() && b.Min.X <= other.Max.X && b.Max.X >= other.Min.X &&
		b.Min.Y <= other.Max.Y && b.Max.Y >= other.Min.Y && b.Min.Z <= other.Max.Z && b.Max.Z >= other.Min.Z
}

func (b Box3) ContainsSphere(sphere Sphere) Containment {
	if !b.Valid() || !sphere.Valid() || !b.IntersectsSphere(sphere) {
		return ContainmentDisjoint
	}
	radius := SplatV3(sphere.Radius)
	if b.ContainsPoint(sphere.Center.Sub(radius)) && b.ContainsPoint(sphere.Center.Add(radius)) {
		return ContainmentContains
	}
	return ContainmentIntersects
}

func (b Box3) IntersectsSphere(sphere Sphere) bool { return sphere.IntersectsBox(b) }

// IntersectsPlane classifies b against plane. Contact is intersecting.
func (b Box3) IntersectsPlane(plane Plane) PlaneIntersection {
	if !b.Valid() || !plane.isFinite() {
		return PlaneIntersectionIntersecting
	}
	negative, positive := b.Min, b.Max
	if plane.Normal.X < 0 {
		negative.X, positive.X = positive.X, negative.X
	}
	if plane.Normal.Y < 0 {
		negative.Y, positive.Y = positive.Y, negative.Y
	}
	if plane.Normal.Z < 0 {
		negative.Z, positive.Z = positive.Z, negative.Z
	}
	if plane.DotCoordinate(negative) > 0 {
		return PlaneIntersectionFront
	}
	if plane.DotCoordinate(positive) < 0 {
		return PlaneIntersectionBack
	}
	return PlaneIntersectionIntersecting
}

func (b Box3) Translated(offset Vec3) Box3 { return NewBox3(b.Min.Add(offset), b.Max.Add(offset)) }
func (b Box3) Corners() [8]Vec3 {
	return [8]Vec3{
		b.Min, V3(b.Max.X, b.Min.Y, b.Min.Z), V3(b.Max.X, b.Max.Y, b.Min.Z), V3(b.Min.X, b.Max.Y, b.Min.Z),
		V3(b.Min.X, b.Min.Y, b.Max.Z), V3(b.Max.X, b.Min.Y, b.Max.Z), b.Max, V3(b.Min.X, b.Max.Y, b.Max.Z),
	}
}

// Transformed returns the axis-aligned bounds of all transformed corners.
func (b Box3) Transformed(matrix Mat4) (Box3, bool) {
	if !b.Valid() {
		return Box3{}, false
	}
	corners := b.Corners()
	points := make([]Vec3, len(corners))
	for index, corner := range corners {
		points[index] = matrix.TransformPoint(corner)
	}
	return Box3FromPoints(points)
}

// EmptyBox3 returns the canonical empty 3D box.
func EmptyBox3() Box3 { return NewBox3(SplatV3(float32(math.Inf(1))), SplatV3(float32(math.Inf(-1)))) }

// Box3FromPoints returns the bounds of points, or false for an empty or non-finite input.
func Box3FromPoints(points []Vec3) (Box3, bool) {
	if len(points) == 0 || !points[0].IsFinite() {
		return Box3{}, false
	}
	result := NewBox3(points[0], points[0])
	for _, point := range points[1:] {
		if !point.IsFinite() {
			return Box3{}, false
		}
		result = result.Expanded(point)
	}
	return result, true
}

// DBox3 converts b to float64 precision.
func (b Box3) DBox3() DBox3 { return NewDBox3(b.Min.DVec3(), b.Max.DVec3()) }

// Box4 is the upstream AABox4D value. Construction preserves endpoint order;
// any reversed component represents an empty box. The zero value is a point
// box at the origin.
type Box4 struct {
	Min, Max Vec4
}

// NewBox4 constructs a Box4 without reordering endpoints.
func NewBox4(minimum, maximum Vec4) Box4 { return Box4{Min: minimum, Max: maximum} }

// AlmostEqual compares corresponding endpoints.
func (b Box4) AlmostEqual(other Box4, tolerance float32) bool {
	return b.Min.AlmostEqual(other.Min, tolerance) && b.Max.AlmostEqual(other.Max, tolerance)
}

func (b Box4) Valid() bool {
	return b.Min.X <= b.Max.X && b.Min.Y <= b.Max.Y && b.Min.Z <= b.Max.Z && b.Min.W <= b.Max.W
}
func (b Box4) Empty() bool              { return !b.Valid() }
func (b Box4) Extent() Vec4             { return b.Max.Sub(b.Min) }
func (b Box4) Center() Vec4             { return b.Min.Add(b.Max).Scale(0.5) }
func (b Box4) Merge(other Box4) Box4    { return NewBox4(b.Min.Min(other.Min), b.Max.Max(other.Max)) }
func (b Box4) Expanded(point Vec4) Box4 { return NewBox4(b.Min.Min(point), b.Max.Max(point)) }
func (b Box4) Intersection(other Box4) (Box4, bool) {
	result := NewBox4(b.Min.Max(other.Min), b.Max.Min(other.Max))
	return result, result.Valid()
}

// EmptyBox4 returns the canonical empty 4D box.
func EmptyBox4() Box4 { return NewBox4(SplatV4(float32(math.Inf(1))), SplatV4(float32(math.Inf(-1)))) }

// DBox2 is the upstream DAABox2D value. Construction preserves endpoint order;
// any reversed component represents an empty box. The zero value is a point
// box at the origin.
type DBox2 struct {
	Min, Max DVec2
}

// NewDBox2 constructs a DBox2 without reordering endpoints.
func NewDBox2(minimum, maximum DVec2) DBox2 { return DBox2{Min: minimum, Max: maximum} }

// AlmostEqual compares corresponding endpoints.
func (b DBox2) AlmostEqual(other DBox2, tolerance float64) bool {
	return b.Min.AlmostEqual(other.Min, tolerance) && b.Max.AlmostEqual(other.Max, tolerance)
}

func (b DBox2) Valid() bool {
	return b.Min.X <= b.Max.X && b.Min.Y <= b.Max.Y
}
func (b DBox2) Empty() bool                { return !b.Valid() }
func (b DBox2) Extent() DVec2              { return b.Max.Sub(b.Min) }
func (b DBox2) Center() DVec2              { return b.Min.Add(b.Max).Scale(0.5) }
func (b DBox2) Merge(other DBox2) DBox2    { return NewDBox2(b.Min.Min(other.Min), b.Max.Max(other.Max)) }
func (b DBox2) Expanded(point DVec2) DBox2 { return NewDBox2(b.Min.Min(point), b.Max.Max(point)) }
func (b DBox2) Intersection(other DBox2) (DBox2, bool) {
	result := NewDBox2(b.Min.Max(other.Min), b.Max.Min(other.Max))
	return result, result.Valid()
}

func EmptyDBox2() DBox2 { return NewDBox2(SplatDV2(math.Inf(1)), SplatDV2(math.Inf(-1))) }

// DBox3 is the upstream DAABox value. Construction preserves endpoint order;
// any reversed component represents an empty box. The zero value is a point
// box at the origin.
type DBox3 struct {
	Min, Max DVec3
}

// NewDBox3 constructs a DBox3 without reordering endpoints.
func NewDBox3(minimum, maximum DVec3) DBox3 { return DBox3{Min: minimum, Max: maximum} }

// AlmostEqual compares corresponding endpoints.
func (b DBox3) AlmostEqual(other DBox3, tolerance float64) bool {
	return b.Min.AlmostEqual(other.Min, tolerance) && b.Max.AlmostEqual(other.Max, tolerance)
}

func (b DBox3) Valid() bool {
	return b.Min.X <= b.Max.X && b.Min.Y <= b.Max.Y && b.Min.Z <= b.Max.Z
}
func (b DBox3) Empty() bool                { return !b.Valid() }
func (b DBox3) Extent() DVec3              { return b.Max.Sub(b.Min) }
func (b DBox3) Center() DVec3              { return b.Min.Add(b.Max).Scale(0.5) }
func (b DBox3) Merge(other DBox3) DBox3    { return NewDBox3(b.Min.Min(other.Min), b.Max.Max(other.Max)) }
func (b DBox3) Expanded(point DVec3) DBox3 { return NewDBox3(b.Min.Min(point), b.Max.Max(point)) }
func (b DBox3) Intersection(other DBox3) (DBox3, bool) {
	result := NewDBox3(b.Min.Max(other.Min), b.Max.Min(other.Max))
	return result, result.Valid()
}
func (b DBox3) ContainsPoint(point DVec3) bool {
	return b.Valid() && point.X >= b.Min.X && point.X <= b.Max.X && point.Y >= b.Min.Y && point.Y <= b.Max.Y && point.Z >= b.Min.Z && point.Z <= b.Max.Z
}

func EmptyDBox3() DBox3 { return NewDBox3(SplatDV3(math.Inf(1)), SplatDV3(math.Inf(-1))) }

// Box3 narrows b to float32 precision. Endpoints can be rounded or overflow to
// an infinity.
func (b DBox3) Box3() Box3 { return NewBox3(b.Min.Vec3(), b.Max.Vec3()) }

// DBox4 is the upstream DAABox4D value. Construction preserves endpoint order;
// any reversed component represents an empty box. The zero value is a point
// box at the origin.
type DBox4 struct {
	Min, Max DVec4
}

// NewDBox4 constructs a DBox4 without reordering endpoints.
func NewDBox4(minimum, maximum DVec4) DBox4 { return DBox4{Min: minimum, Max: maximum} }

// AlmostEqual compares corresponding endpoints.
func (b DBox4) AlmostEqual(other DBox4, tolerance float64) bool {
	return b.Min.AlmostEqual(other.Min, tolerance) && b.Max.AlmostEqual(other.Max, tolerance)
}

func (b DBox4) Valid() bool {
	return b.Min.X <= b.Max.X && b.Min.Y <= b.Max.Y && b.Min.Z <= b.Max.Z && b.Min.W <= b.Max.W
}
func (b DBox4) Empty() bool                { return !b.Valid() }
func (b DBox4) Extent() DVec4              { return b.Max.Sub(b.Min) }
func (b DBox4) Center() DVec4              { return b.Min.Add(b.Max).Scale(0.5) }
func (b DBox4) Merge(other DBox4) DBox4    { return NewDBox4(b.Min.Min(other.Min), b.Max.Max(other.Max)) }
func (b DBox4) Expanded(point DVec4) DBox4 { return NewDBox4(b.Min.Min(point), b.Max.Max(point)) }
func (b DBox4) Intersection(other DBox4) (DBox4, bool) {
	result := NewDBox4(b.Min.Max(other.Min), b.Max.Min(other.Max))
	return result, result.Valid()
}

func EmptyDBox4() DBox4 { return NewDBox4(SplatDV4(math.Inf(1)), SplatDV4(math.Inf(-1))) }

func finite32(value float32) bool {
	return !math.IsNaN(float64(value)) && !math.IsInf(float64(value), 0)
}
func finite64(value float64) bool { return !math.IsNaN(value) && !math.IsInf(value, 0) }
