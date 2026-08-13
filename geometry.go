package math3d

import "math"

// Transform is a rigid transform represented by a position and orientation.
// It has no scale. Its zero value has a zero quaternion and is therefore not a
// valid rigid transform; use IdentityTransform when identity is required.
type Transform struct {
	Position    Vec3
	Orientation Quat
}

// NewTransform constructs a Transform without normalizing orientation.
func NewTransform(position Vec3, orientation Quat) Transform {
	return Transform{Position: position, Orientation: orientation}
}

// IdentityTransform returns the identity rigid transform.
func IdentityTransform() Transform { return NewTransform(Vec3{}, IdentityQuat()) }

// AlmostEqual compares the stored representation, not geometric equivalence.
func (t Transform) AlmostEqual(other Transform, tolerance float32) bool {
	return t.Position.AlmostEqual(other.Position, tolerance) &&
		t.Orientation.AlmostEqual(other.Orientation, tolerance)
}

// Mat4 returns the row-vector matrix that rotates by Orientation and then
// translates by Position. Orientation must be a unit quaternion for the result
// to be a rigid transform.
func (t Transform) Mat4() Mat4 {
	return ComposeMat4(V3(1, 1, 1), t.Orientation, t.Position)
}

// TransformPoint rotates point by Orientation and then adds Position.
// Orientation must be a unit quaternion.
func (t Transform) TransformPoint(point Vec3) Vec3 {
	return t.Orientation.Rotate(point).Add(t.Position)
}

// TransformDirection rotates direction by Orientation without translating it.
// Orientation must be a unit quaternion.
func (t Transform) TransformDirection(direction Vec3) Vec3 {
	return t.Orientation.Rotate(direction)
}

// Plane represents dot(Normal, point) + D = 0. The normal is not implicitly
// normalized. The zero value does not define a geometric plane.
type Plane struct {
	Normal Vec3
	D      float32
}

// NewPlane constructs a Plane without normalizing normal or changing D.
func NewPlane(normal Vec3, d float32) Plane { return Plane{Normal: normal, D: d} }

// PlaneFromVec4 constructs a Plane from normal components followed by D.
func PlaneFromVec4(value Vec4) Plane { return NewPlane(V3(value.X, value.Y, value.Z), value.W) }

// Vec4 returns the plane's normal components followed by D.
func (p Plane) Vec4() Vec4 { return V4(p.Normal.X, p.Normal.Y, p.Normal.Z, p.D) }

// AlmostEqual compares corresponding stored components.
func (p Plane) AlmostEqual(other Plane, tolerance float32) bool {
	return p.Normal.AlmostEqual(other.Normal, tolerance) && abs32(p.D-other.D) < tolerance
}

func (p Plane) isFinite() bool {
	return p.Normal.IsFinite() && !math.IsNaN(float64(p.D)) && !math.IsInf(float64(p.D), 0)
}

// PlaneFromPoints constructs a normalized plane through a, b, and c. Vertex
// order determines the normal direction. It fails when the points do not
// define a finite, non-degenerate plane.
func PlaneFromPoints(a, b, c Vec3) (Plane, bool) {
	normal, ok := b.Sub(a).Cross(c.Sub(a)).Normalized()
	if !ok {
		return Plane{}, false
	}
	result := NewPlane(normal, -normal.Dot(a))
	if !result.isFinite() {
		return Plane{}, false
	}
	return result, true
}

// PlaneFromNormalPoint constructs a normalized plane containing point. It
// fails when normal cannot be normalized. Unlike the inconsistent upstream
// factory, D follows dot(Normal, point)+D=0.
func PlaneFromNormalPoint(normal, point Vec3) (Plane, bool) {
	normal, ok := normal.Normalized()
	if !ok {
		return Plane{}, false
	}
	result := NewPlane(normal, -normal.Dot(point))
	if !result.isFinite() {
		return Plane{}, false
	}
	return result, true
}

// XYPlane returns the plane z=0.
func XYPlane() Plane { return NewPlane(V3(0, 0, 1), 0) }

// XZPlane returns the plane y=0.
func XZPlane() Plane { return NewPlane(V3(0, 1, 0), 0) }

// YZPlane returns the plane x=0.
func YZPlane() Plane { return NewPlane(V3(1, 0, 0), 0) }

// Normalized returns an equivalent plane with a unit normal. It fails for a
// zero or non-finite normal.
func (p Plane) Normalized() (Plane, bool) {
	lengthSquared := p.Normal.Dot(p.Normal)
	if lengthSquared == 0 || math.IsNaN(float64(lengthSquared)) || math.IsInf(float64(lengthSquared), 0) || !p.isFinite() {
		return Plane{}, false
	}
	// Match upstream's exact-preservation fast path for an already-unit plane.
	if abs32(lengthSquared-1) < 1.192092896e-7 {
		return p, true
	}
	length := float32(math.Sqrt(float64(lengthSquared)))
	result := NewPlane(p.Normal.Scale(1/length), p.D/length)
	if !result.isFinite() {
		return Plane{}, false
	}
	return result, true
}

// Transformed applies matrix's inverse transpose to p. It fails when matrix
// is singular or when the resulting plane has non-finite components.
func (p Plane) Transformed(matrix Mat4) (Plane, bool) {
	inverse, ok := matrix.Inverse()
	if !ok {
		return Plane{}, false
	}
	result := PlaneFromVec4(inverse.Transposed().TransformVec4(p.Vec4()))
	if !result.isFinite() {
		return Plane{}, false
	}
	return result, true
}

// Rotated rotates p's normal and leaves D unchanged. Rotation must be a unit
// quaternion.
func (p Plane) Rotated(rotation Quat) Plane {
	return NewPlane(rotation.Rotate(p.Normal), p.D)
}

// Project projects point orthogonally onto p. It supports non-unit normals and
// fails for a zero or non-finite normal or a non-finite result.
func (p Plane) Project(point Vec3) (Vec3, bool) {
	denominator := p.Normal.Dot(p.Normal)
	if denominator == 0 || math.IsNaN(float64(denominator)) || math.IsInf(float64(denominator), 0) || !p.isFinite() || !point.IsFinite() {
		return Vec3{}, false
	}
	result := point.Sub(p.Normal.Scale(p.DotCoordinate(point) / denominator))
	if !result.IsFinite() {
		return Vec3{}, false
	}
	return result, true
}

// Dot returns the four-component dot product of p and value.
func (p Plane) Dot(value Vec4) float32 { return p.Vec4().Dot(value) }

// DotCoordinate evaluates dot(Normal, point)+D.
func (p Plane) DotCoordinate(point Vec3) float32 { return p.Normal.Dot(point) + p.D }

// DotNormal evaluates dot(Normal, direction), ignoring D.
func (p Plane) DotNormal(direction Vec3) float32 { return p.Normal.Dot(direction) }

// ClassifyPoint returns the signed plane equation value for point. Negative
// values are behind the normal, positive values are in front, and zero is on
// the plane. It is a signed distance only when Normal is unit length.
func (p Plane) ClassifyPoint(point Vec3) float32 { return p.DotCoordinate(point) }

// DPlane is the float64 plane representation using the equation
// dot(Normal, point) + D = 0. Its zero value does not define a plane.
type DPlane struct {
	Normal DVec3
	D      float64
}

// NewDPlane constructs a DPlane without normalizing normal or changing D.
func NewDPlane(normal DVec3, d float64) DPlane { return DPlane{Normal: normal, D: d} }

// DPlaneFromVec4 constructs a DPlane from normal components followed by D.
func DPlaneFromVec4(value DVec4) DPlane {
	return NewDPlane(DV3(value.X, value.Y, value.Z), value.W)
}

// DVec4 returns the plane's normal components followed by D.
func (p DPlane) DVec4() DVec4 { return DV4(p.Normal.X, p.Normal.Y, p.Normal.Z, p.D) }

// AlmostEqual compares corresponding stored components.
func (p DPlane) AlmostEqual(other DPlane, tolerance float64) bool {
	return p.Normal.AlmostEqual(other.Normal, tolerance) &&
		almostEqual64(tolerance, p.D-other.D)
}

// Ray stores an origin and direction. Direction is not implicitly normalized,
// and its length scales the ray parameter. A zero direction is representable
// but does not define a geometric ray.
type Ray struct {
	Origin, Direction Vec3
}

// NewRay constructs a Ray without normalizing direction.
func NewRay(origin, direction Vec3) Ray { return Ray{Origin: origin, Direction: direction} }

// AlmostEqual compares corresponding stored components.
func (r Ray) AlmostEqual(other Ray, tolerance float32) bool {
	return r.Origin.AlmostEqual(other.Origin, tolerance) &&
		r.Direction.AlmostEqual(other.Direction, tolerance)
}

// PointAt evaluates the ray at parameter t. The parameter is a Euclidean
// distance only when Direction has unit length.
func (r Ray) PointAt(t float32) Vec3 { return r.Origin.Add(r.Direction.Scale(t)) }

// IntersectBox returns the first non-negative parameter at which r meets box.
func (r Ray) IntersectBox(box Box3) (float32, bool) {
	t, ok := intersectSlabs(
		[3]float64{float64(r.Origin.X), float64(r.Origin.Y), float64(r.Origin.Z)},
		[3]float64{float64(r.Direction.X), float64(r.Direction.Y), float64(r.Direction.Z)},
		[3]float64{float64(box.Min.X), float64(box.Min.Y), float64(box.Min.Z)},
		[3]float64{float64(box.Max.X), float64(box.Max.Y), float64(box.Max.Z)},
		box.Valid(),
	)
	result := float32(t)
	return result, ok && finite32(result)
}

// IntersectPlane returns the non-negative parameter at which r meets plane.
// Parallel rays, including rays contained in the plane, do not have one
// distinguished intersection and return false.
func (r Ray) IntersectPlane(plane Plane, tolerance float32) (float32, bool) {
	denominator := r.Direction.Dot(plane.Normal)
	if tolerance < 0 || !r.Origin.IsFinite() || !r.Direction.IsFinite() || !plane.isFinite() ||
		abs32(denominator) < tolerance {
		return 0, false
	}
	t := -plane.DotCoordinate(r.Origin) / denominator
	if !finite32(t) || t < -tolerance {
		return 0, false
	}
	if t < 0 {
		t = 0
	}
	return t, true
}

// IntersectSphere returns the first non-negative ray parameter at which r
// meets sphere. It supports directions of any non-zero finite length.
func (r Ray) IntersectSphere(sphere Sphere) (float32, bool) {
	t, ok := intersectSphere64(r.Origin.DVec3(), r.Direction.DVec3(), sphere.Center.DVec3(), float64(sphere.Radius), sphere.Valid())
	result := float32(t)
	return result, ok && finite32(result)
}

// IntersectTriangle returns the positive ray parameter at which r meets
// triangle. Both windings are accepted, and edges and vertices are included.
// Hits at the ray origin or within tolerance of it are excluded.
func (r Ray) IntersectTriangle(triangle Triangle3, tolerance float32) (float32, bool) {
	if tolerance < 0 || !finite32(tolerance) || !r.Origin.IsFinite() || !r.Direction.IsFinite() ||
		!triangle.A.IsFinite() || !triangle.B.IsFinite() || !triangle.C.IsFinite() {
		return 0, false
	}
	edge1 := triangle.B.Sub(triangle.A)
	edge2 := triangle.C.Sub(triangle.A)
	h := r.Direction.Cross(edge2)
	determinant := edge1.Dot(h)
	if abs32(determinant) < tolerance {
		return 0, false
	}
	inverse := 1 / determinant
	s := r.Origin.Sub(triangle.A)
	u := inverse * s.Dot(h)
	if u < 0 || u > 1 {
		return 0, false
	}
	q := s.Cross(edge1)
	v := inverse * r.Direction.Dot(q)
	if v < 0 || u+v > 1 {
		return 0, false
	}
	parameter := inverse * edge2.Dot(q)
	return parameter, finite32(parameter) && parameter > tolerance
}

// Transformed applies matrix to the origin and its linear part to Direction.
func (r Ray) Transformed(matrix Mat4) (Ray, bool) {
	result := NewRay(matrix.TransformPoint(r.Origin), matrix.TransformDirection(r.Direction))
	if !result.Origin.IsFinite() || !result.Direction.IsFinite() || result.Direction.Dot(result.Direction) == 0 {
		return Ray{}, false
	}
	return result, true
}

// DRay is the float64 ray representation. Direction is not implicitly
// normalized. A zero direction does not define a geometric ray.
type DRay struct {
	Origin, Direction DVec3
}

// NewDRay constructs a DRay without normalizing direction.
func NewDRay(origin, direction DVec3) DRay { return DRay{Origin: origin, Direction: direction} }

// AlmostEqual compares corresponding stored components.
func (r DRay) AlmostEqual(other DRay, tolerance float64) bool {
	return r.Origin.AlmostEqual(other.Origin, tolerance) &&
		r.Direction.AlmostEqual(other.Direction, tolerance)
}

func (r DRay) PointAt(t float64) DVec3 { return r.Origin.Add(r.Direction.Scale(t)) }

func (r DRay) IntersectBox(box DBox3) (float64, bool) {
	return intersectSlabs(
		[3]float64{r.Origin.X, r.Origin.Y, r.Origin.Z},
		[3]float64{r.Direction.X, r.Direction.Y, r.Direction.Z},
		[3]float64{box.Min.X, box.Min.Y, box.Min.Z},
		[3]float64{box.Max.X, box.Max.Y, box.Max.Z},
		box.Valid(),
	)
}

func (r DRay) IntersectSphere(sphere DSphere) (float64, bool) {
	return intersectSphere64(r.Origin, r.Direction, sphere.Center, sphere.Radius, sphere.Valid())
}

// Sphere stores a center and radius. NewSphere preserves any radius, including
// a negative one; algorithms that require a non-negative finite radius state
// that precondition. The zero value is a radius-zero sphere at the origin.
type Sphere struct {
	Center Vec3
	Radius float32
}

// NewSphere constructs a Sphere without validating radius.
func NewSphere(center Vec3, radius float32) Sphere { return Sphere{Center: center, Radius: radius} }

// AlmostEqual compares corresponding stored components.
func (s Sphere) AlmostEqual(other Sphere, tolerance float32) bool {
	return s.Center.AlmostEqual(other.Center, tolerance) && abs32(s.Radius-other.Radius) < tolerance
}

// Valid reports whether s has a finite center and non-negative finite radius.
func (s Sphere) Valid() bool { return s.Center.IsFinite() && finite32(s.Radius) && s.Radius >= 0 }

func (s Sphere) ContainsPoint(point Vec3) Containment {
	if !s.Valid() || !point.IsFinite() {
		return ContainmentDisjoint
	}
	delta := point.DVec3().Sub(s.Center.DVec3())
	distanceSquared := delta.Dot(delta)
	radiusSquared := float64(s.Radius) * float64(s.Radius)
	if distanceSquared < radiusSquared {
		return ContainmentContains
	}
	if distanceSquared == radiusSquared {
		return ContainmentIntersects
	}
	return ContainmentDisjoint
}

func (s Sphere) ContainsSphere(other Sphere) Containment {
	if !s.Valid() || !other.Valid() {
		return ContainmentDisjoint
	}
	distance := s.Center.DVec3().Sub(other.Center.DVec3()).Magnitude()
	if distance > float64(s.Radius)+float64(other.Radius) {
		return ContainmentDisjoint
	}
	if s.Radius >= other.Radius && distance <= float64(s.Radius-other.Radius) {
		return ContainmentContains
	}
	return ContainmentIntersects
}

func (s Sphere) ContainsBox(box Box3) Containment {
	if !s.Valid() || !box.Valid() {
		return ContainmentDisjoint
	}
	allInside := true
	for _, corner := range box.Corners() {
		if s.ContainsPoint(corner) != ContainmentContains {
			allInside = false
			break
		}
	}
	if allInside {
		return ContainmentContains
	}
	if s.IntersectsBox(box) {
		return ContainmentIntersects
	}
	return ContainmentDisjoint
}

func (s Sphere) Intersects(other Sphere) bool {
	return s.Valid() && other.Valid() &&
		s.Center.DVec3().Sub(other.Center.DVec3()).Magnitude() <= float64(s.Radius)+float64(other.Radius)
}

func (s Sphere) IntersectsBox(box Box3) bool {
	if !s.Valid() || !box.Valid() {
		return false
	}
	closest := s.Center.Max(box.Min).Min(box.Max)
	delta := s.Center.DVec3().Sub(closest.DVec3())
	return delta.Dot(delta) <= float64(s.Radius)*float64(s.Radius)
}

func (s Sphere) IntersectsPlane(plane Plane) PlaneIntersection {
	if !s.Valid() || !plane.isFinite() {
		return PlaneIntersectionIntersecting
	}
	length := plane.Normal.Magnitude()
	if length == 0 || math.IsNaN(length) || math.IsInf(length, 0) {
		return PlaneIntersectionIntersecting
	}
	distance := float64(plane.DotCoordinate(s.Center)) / length
	if distance > float64(s.Radius) {
		return PlaneIntersectionFront
	}
	if distance < -float64(s.Radius) {
		return PlaneIntersectionBack
	}
	return PlaneIntersectionIntersecting
}

func (s Sphere) Bounds() Box3 {
	radius := SplatV3(s.Radius)
	return NewBox3(s.Center.Sub(radius), s.Center.Add(radius))
}

func (s Sphere) DistanceToPoint(point Vec3) float32 {
	return max(float32(s.Center.Sub(point).Magnitude())-s.Radius, 0)
}

func (s Sphere) DistanceToSphere(other Sphere) float32 {
	return max(float32(s.Center.Sub(other.Center).Magnitude())-s.Radius-other.Radius, 0)
}

func (s Sphere) Translated(offset Vec3) Sphere { return NewSphere(s.Center.Add(offset), s.Radius) }

// Merge returns the smallest sphere on the line between the centers that
// encloses both spheres. Both inputs must be valid.
func (s Sphere) Merge(other Sphere) (Sphere, bool) {
	if !s.Valid() || !other.Valid() {
		return Sphere{}, false
	}
	delta := other.Center.DVec3().Sub(s.Center.DVec3())
	distance := delta.Magnitude()
	if float64(s.Radius) >= distance+float64(other.Radius) {
		return s, true
	}
	if float64(other.Radius) >= distance+float64(s.Radius) {
		return other, true
	}
	radius := (distance + float64(s.Radius) + float64(other.Radius)) / 2
	center := s.Center.DVec3().Add(delta.Scale((radius - float64(s.Radius)) / distance)).Vec3()
	result := NewSphere(center, float32(radius))
	if radius > math.MaxFloat32 || !result.Valid() {
		return Sphere{}, false
	}
	return result, true
}

// Transformed returns a sphere enclosing s after an affine transform. Its
// radius is scaled by the largest singular value of the linear matrix, which
// is exact for rotation and non-uniform scale and conservative for shear.
func (s Sphere) Transformed(matrix Mat4) (Sphere, bool) {
	if !s.Valid() || matrix.M14 != 0 || matrix.M24 != 0 || matrix.M34 != 0 || matrix.M44 != 1 {
		return Sphere{}, false
	}
	scale := maxLinearScale(matrix)
	center := matrix.TransformPoint(s.Center)
	radius := float64(s.Radius) * scale
	if !center.IsFinite() || math.IsNaN(radius) || math.IsInf(radius, 0) || radius > math.MaxFloat32 {
		return Sphere{}, false
	}
	return NewSphere(center, float32(radius)), true
}

// DSphere is the float64 sphere representation. NewDSphere preserves negative
// radii. The zero value is a radius-zero sphere at the origin.
type DSphere struct {
	Center DVec3
	Radius float64
}

// NewDSphere constructs a DSphere without validating radius.
func NewDSphere(center DVec3, radius float64) DSphere {
	return DSphere{Center: center, Radius: radius}
}

// AlmostEqual compares corresponding stored components.
func (s DSphere) AlmostEqual(other DSphere, tolerance float64) bool {
	return s.Center.AlmostEqual(other.Center, tolerance) &&
		almostEqual64(tolerance, s.Radius-other.Radius)
}

func (s DSphere) Valid() bool { return s.Center.IsFinite() && finite64(s.Radius) && s.Radius >= 0 }
func (s DSphere) ContainsPoint(point DVec3) Containment {
	if !s.Valid() || !point.IsFinite() {
		return ContainmentDisjoint
	}
	delta := point.Sub(s.Center)
	distanceSquared, radiusSquared := delta.Dot(delta), s.Radius*s.Radius
	if distanceSquared < radiusSquared {
		return ContainmentContains
	}
	if distanceSquared == radiusSquared {
		return ContainmentIntersects
	}
	return ContainmentDisjoint
}
func (s DSphere) Intersects(other DSphere) bool {
	return s.Valid() && other.Valid() && s.Center.Sub(other.Center).Magnitude() <= s.Radius+other.Radius
}
func (s DSphere) Bounds() DBox3 {
	radius := SplatDV3(s.Radius)
	return NewDBox3(s.Center.Sub(radius), s.Center.Add(radius))
}

func intersectSlabs(origin, direction, minimum, maximum [3]float64, validBox bool) (float64, bool) {
	if !validBox {
		return 0, false
	}
	tMin, tMax := 0.0, math.Inf(1)
	hasDirection := false
	for axis := range origin {
		if !finite64(origin[axis]) || !finite64(direction[axis]) {
			return 0, false
		}
		if direction[axis] == 0 {
			if origin[axis] < minimum[axis] || origin[axis] > maximum[axis] {
				return 0, false
			}
			continue
		}
		hasDirection = true
		near := (minimum[axis] - origin[axis]) / direction[axis]
		far := (maximum[axis] - origin[axis]) / direction[axis]
		if near > far {
			near, far = far, near
		}
		tMin, tMax = max(tMin, near), min(tMax, far)
		if tMin > tMax {
			return 0, false
		}
	}
	return tMin, hasDirection && finite64(tMin)
}

func intersectSphere64(origin, direction, center DVec3, radius float64, validSphere bool) (float64, bool) {
	if !validSphere || !origin.IsFinite() || !direction.IsFinite() {
		return 0, false
	}
	a := direction.Dot(direction)
	if a == 0 || !finite64(a) {
		return 0, false
	}
	offset := origin.Sub(center)
	c := offset.Dot(offset) - radius*radius
	if c <= 0 {
		return 0, true
	}
	b := offset.Dot(direction)
	if b >= 0 {
		return 0, false
	}
	discriminant := b*b - a*c
	if discriminant < 0 || !finite64(discriminant) {
		return 0, false
	}
	t := (-b - math.Sqrt(discriminant)) / a
	return t, t >= 0 && finite64(t)
}

// maxLinearScale returns the spectral norm of matrix's upper-left 3x3 block.
// It computes the largest eigenvalue of L*Lᵀ using the stable closed form for
// a real symmetric 3x3 matrix.
func maxLinearScale(matrix Mat4) float64 {
	r0 := [3]float64{float64(matrix.M11), float64(matrix.M12), float64(matrix.M13)}
	r1 := [3]float64{float64(matrix.M21), float64(matrix.M22), float64(matrix.M23)}
	r2 := [3]float64{float64(matrix.M31), float64(matrix.M32), float64(matrix.M33)}
	dot := func(a, b [3]float64) float64 { return a[0]*b[0] + a[1]*b[1] + a[2]*b[2] }
	a00, a11, a22 := dot(r0, r0), dot(r1, r1), dot(r2, r2)
	a01, a02, a12 := dot(r0, r1), dot(r0, r2), dot(r1, r2)
	offDiagonal := a01*a01 + a02*a02 + a12*a12
	if offDiagonal == 0 {
		return math.Sqrt(max(a00, a11, a22))
	}
	mean := (a00 + a11 + a22) / 3
	spread := math.Sqrt(((a00-mean)*(a00-mean) + (a11-mean)*(a11-mean) + (a22-mean)*(a22-mean) + 2*offDiagonal) / 6)
	b00, b11, b22 := (a00-mean)/spread, (a11-mean)/spread, (a22-mean)/spread
	b01, b02, b12 := a01/spread, a02/spread, a12/spread
	determinant := b00*b11*b22 + 2*b01*b02*b12 - b00*b12*b12 - b11*b02*b02 - b22*b01*b01
	angle := math.Acos(max(-1, min(1, determinant/2))) / 3
	largestEigenvalue := mean + 2*spread*math.Cos(angle)
	return math.Sqrt(max(largestEigenvalue, 0))
}

// Segment3 is the upstream Line value, represented by endpoints A and B. Its
// zero value is a degenerate segment at the origin.
type Segment3 struct{ A, B Vec3 }

// NewSegment3 constructs a Segment3 from its endpoints.
func NewSegment3(a, b Vec3) Segment3 { return Segment3{A: a, B: b} }

// AlmostEqual compares corresponding endpoints.
func (s Segment3) AlmostEqual(other Segment3, tolerance float32) bool {
	return s.A.AlmostEqual(other.A, tolerance) && s.B.AlmostEqual(other.B, tolerance)
}

// Vector returns the directed displacement from A to B.
func (s Segment3) Vector() Vec3 { return s.B.Sub(s.A) }

// Length returns the Euclidean distance between the endpoints.
func (s Segment3) Length() float64 { return s.Vector().Magnitude() }

// LengthSquared returns the squared distance between the endpoints.
func (s Segment3) LengthSquared() float32 { return s.Vector().Dot(s.Vector()) }

// Midpoint returns the point halfway between the endpoints.
func (s Segment3) Midpoint() Vec3 { return s.A.Add(s.B).Scale(0.5) }

// PointAt linearly maps parameter t onto the segment's supporting line. Values
// outside [0, 1] extrapolate beyond the endpoints.
func (s Segment3) PointAt(t float32) Vec3 { return s.A.Add(s.Vector().Scale(t)) }

// Reversed returns the same segment with opposite endpoint order.
func (s Segment3) Reversed() Segment3 { return NewSegment3(s.B, s.A) }

// Ray returns a ray starting at A with the unnormalized displacement to B.
func (s Segment3) Ray() Ray { return NewRay(s.A, s.Vector()) }

// WithLength returns a segment starting at A in the same direction with the
// requested length. It fails when the direction cannot be normalized or length
// is negative or non-finite.
func (s Segment3) WithLength(length float32) (Segment3, bool) {
	direction, ok := s.Vector().Normalized()
	if !ok || length < 0 || !finite32(length) {
		return Segment3{}, false
	}
	return NewSegment3(s.A, s.A.Add(direction.Scale(length))), true
}

// Map applies mapper to both endpoints while preserving their order.
func (s Segment3) Map(mapper func(Vec3) Vec3) Segment3 {
	return NewSegment3(mapper(s.A), mapper(s.B))
}

// Transformed applies matrix to both endpoints as points.
func (s Segment3) Transformed(matrix Mat4) Segment3 {
	return s.Map(matrix.TransformPoint)
}

// Segment2 is the upstream Line2D value. Its zero value is degenerate.
type Segment2 struct{ A, B Vec2 }

// NewSegment2 constructs a Segment2 from its endpoints.
func NewSegment2(a, b Vec2) Segment2 { return Segment2{A: a, B: b} }

// AlmostEqual compares corresponding endpoints.
func (s Segment2) AlmostEqual(other Segment2, tolerance float32) bool {
	return s.A.AlmostEqual(other.A, tolerance) && s.B.AlmostEqual(other.B, tolerance)
}

// Bounds returns the componentwise bounds of the endpoints.
func (s Segment2) Bounds() Box2 { return NewBox2(s.A.Min(s.B), s.A.Max(s.B)) }

func (s Segment2) linePointCross(point Vec2) float32 {
	return s.Vector().Cross(point.Sub(s.A))
}

// Vector returns the directed displacement from A to B.
func (s Segment2) Vector() Vec2 { return s.B.Sub(s.A) }

// ContainsOnLine reports whether point is within tolerance of the infinite
// supporting line. It does not restrict the point to the segment's bounds.
func (s Segment2) ContainsOnLine(point Vec2, tolerance float32) bool {
	return tolerance >= 0 && abs32(s.linePointCross(point)) < tolerance
}

// Intersects reports whether two closed 2D segments touch or cross. Collinear
// overlap and shared endpoints count as intersections. Tolerance is an
// absolute cross-product tolerance, matching the upstream operation.
func (s Segment2) Intersects(other Segment2, tolerance float32) bool {
	if tolerance < 0 || !s.A.IsFinite() || !s.B.IsFinite() || !other.A.IsFinite() || !other.B.IsFinite() ||
		!s.Bounds().Intersects(other.Bounds()) {
		return false
	}
	touchesOrCrosses := func(a, b Segment2) bool {
		crossA := a.linePointCross(b.A)
		crossB := a.linePointCross(b.B)
		return abs32(crossA) < tolerance || abs32(crossB) < tolerance || (crossA < 0) != (crossB < 0)
	}
	return touchesOrCrosses(s, other) && touchesOrCrosses(other, s)
}

// Triangle3 is the upstream Triangle value. Vertices retain caller-supplied
// order. The zero value is degenerate.
type Triangle3 struct{ A, B, C Vec3 }

// NewTriangle3 constructs a Triangle3 in A, B, C order.
func NewTriangle3(a, b, c Vec3) Triangle3 { return Triangle3{A: a, B: b, C: c} }

// AlmostEqual compares corresponding vertices without disregarding winding.
func (t Triangle3) AlmostEqual(other Triangle3, tolerance float32) bool {
	return t.A.AlmostEqual(other.A, tolerance) && t.B.AlmostEqual(other.B, tolerance) &&
		t.C.AlmostEqual(other.C, tolerance)
}

// Map applies mapper to all vertices while preserving their order.
func (t Triangle3) Map(mapper func(Vec3) Vec3) Triangle3 {
	return NewTriangle3(mapper(t.A), mapper(t.B), mapper(t.C))
}

// Transformed applies matrix to every vertex as a point.
func (t Triangle3) Transformed(matrix Mat4) Triangle3 {
	return t.Map(matrix.TransformPoint)
}

// NormalDirection returns the winding-dependent, unnormalized surface normal.
func (t Triangle3) NormalDirection() Vec3 { return t.B.Sub(t.A).Cross(t.C.Sub(t.A)) }

// Normal returns the unit surface normal. It fails for a non-finite or
// zero-area triangle.
func (t Triangle3) Normal() (Vec3, bool) { return t.NormalDirection().Normalized() }

// Area returns the non-negative geometric area.
func (t Triangle3) Area() float32 { return float32(t.NormalDirection().Magnitude() * 0.5) }

// Perimeter returns the sum of the three edge lengths.
func (t Triangle3) Perimeter() float64 {
	return t.A.Sub(t.B).Magnitude() + t.B.Sub(t.C).Magnitude() + t.C.Sub(t.A).Magnitude()
}

// Centroid returns the arithmetic mean of the vertices.
func (t Triangle3) Centroid() Vec3 { return t.A.Add(t.B).Add(t.C).Scale(1.0 / 3.0) }

// Degenerate reports whether the vertices are non-finite or have exactly zero
// geometric area. Use IsSliver when a non-zero minimum edge length is needed.
func (t Triangle3) Degenerate() bool {
	normal := t.NormalDirection()
	return !t.A.IsFinite() || !t.B.IsFinite() || !t.C.IsFinite() || !normal.IsFinite() || normal.Dot(normal) == 0
}

// IsSliver reports whether any edge length is at most tolerance. A negative or
// non-finite tolerance, or non-finite vertices, is treated as invalid.
func (t Triangle3) IsSliver(tolerance float32) bool {
	if tolerance < 0 || !finite32(tolerance) || !t.A.IsFinite() || !t.B.IsFinite() || !t.C.IsFinite() {
		return true
	}
	threshold := float64(tolerance)
	return t.A.Sub(t.B).Magnitude() <= threshold || t.B.Sub(t.C).Magnitude() <= threshold ||
		t.C.Sub(t.A).Magnitude() <= threshold
}

// Bounds returns the componentwise bounds of the vertices.
func (t Triangle3) Bounds() Box3 {
	return NewBox3(t.A.Min(t.B).Min(t.C), t.A.Max(t.B).Max(t.C))
}

// Edges returns the directed AB, BC, and CA segments.
func (t Triangle3) Edges() [3]Segment3 {
	return [3]Segment3{NewSegment3(t.A, t.B), NewSegment3(t.B, t.C), NewSegment3(t.C, t.A)}
}

// Triangle2 is the upstream Triangle2D value. Vertices retain caller-supplied
// order. The zero value is degenerate.
type Triangle2 struct{ A, B, C Vec2 }

// NewTriangle2 constructs a Triangle2 in A, B, C order.
func NewTriangle2(a, b, c Vec2) Triangle2 { return Triangle2{A: a, B: b, C: c} }

// AlmostEqual compares corresponding vertices without disregarding winding.
func (t Triangle2) AlmostEqual(other Triangle2, tolerance float32) bool {
	return t.A.AlmostEqual(other.A, tolerance) && t.B.AlmostEqual(other.B, tolerance) &&
		t.C.AlmostEqual(other.C, tolerance)
}

// SignedArea returns the oriented area. Positive values indicate
// counter-clockwise winding.
func (t Triangle2) SignedArea() float32 {
	return t.B.Sub(t.A).Cross(t.C.Sub(t.A)) * 0.5
}

// Area returns the non-negative geometric area.
func (t Triangle2) Area() float32 { return abs32(t.SignedArea()) }

// Contains reports whether point lies in the strict interior. Points on an
// edge are excluded, matching the upstream operation.
func (t Triangle2) Contains(point Vec2) bool {
	v0, v1, v2 := t.B.Sub(t.A), t.C.Sub(t.A), point.Sub(t.A)
	d00, d01, d11 := v0.Dot(v0), v0.Dot(v1), v1.Dot(v1)
	denominator := d00*d11 - d01*d01
	if denominator == 0 || !finite32(denominator) || !point.IsFinite() {
		return false
	}
	u := (d11*v2.Dot(v0) - d01*v2.Dot(v1)) / denominator
	v := (d00*v2.Dot(v1) - d01*v2.Dot(v0)) / denominator
	return u > 0 && v > 0 && u+v < 1
}

// Quad3 is the upstream Quad value. Vertices retain caller-supplied order. The
// zero value is degenerate.
type Quad3 struct{ A, B, C, D Vec3 }

// NewQuad3 constructs a Quad3 in A, B, C, D order.
func NewQuad3(a, b, c, d Vec3) Quad3 { return Quad3{A: a, B: b, C: c, D: d} }

// AlmostEqual compares corresponding vertices without disregarding winding.
func (q Quad3) AlmostEqual(other Quad3, tolerance float32) bool {
	return q.A.AlmostEqual(other.A, tolerance) && q.B.AlmostEqual(other.B, tolerance) &&
		q.C.AlmostEqual(other.C, tolerance) && q.D.AlmostEqual(other.D, tolerance)
}

// Map applies mapper to all vertices while preserving their order.
func (q Quad3) Map(mapper func(Vec3) Vec3) Quad3 {
	return NewQuad3(mapper(q.A), mapper(q.B), mapper(q.C), mapper(q.D))
}

// Transformed applies matrix to every vertex as a point.
func (q Quad3) Transformed(matrix Mat4) Quad3 { return q.Map(matrix.TransformPoint) }

// Quad2 is the upstream Quad2D value. Vertices retain caller-supplied order.
// The zero value is degenerate.
type Quad2 struct{ A, B, C, D Vec2 }

// NewQuad2 constructs a Quad2 in A, B, C, D order.
func NewQuad2(a, b, c, d Vec2) Quad2 { return Quad2{A: a, B: b, C: c, D: d} }

// AlmostEqual compares corresponding vertices without disregarding winding.
func (q Quad2) AlmostEqual(other Quad2, tolerance float32) bool {
	return q.A.AlmostEqual(other.A, tolerance) && q.B.AlmostEqual(other.B, tolerance) &&
		q.C.AlmostEqual(other.C, tolerance) && q.D.AlmostEqual(other.D, tolerance)
}
