# Pick and intersect geometry

This guide builds a world-space ray from normalized screen coordinates and
uses it to query boxes, spheres, planes, and triangles.

## Build a world-space screen ray

Create a right-handed view matrix and a perspective projection. Local camera
forward is negative Z, projection depth is `[0, 1]`, and matrix multiplication
follows application order:

```go
view, ok := math3d.LookAtMat4(cameraPosition, target, math3d.V3(0, 1, 0))
if !ok {
	return errors.New("camera direction and up must define a view")
}
projection, ok := math3d.PerspectiveFOVMat4(
	math3d.HalfPi,
	viewportWidth/viewportHeight,
	1,
	1000,
)
if !ok {
	return errors.New("invalid perspective parameters")
}
viewProjection := view.Mul(projection)
```

`view.Mul(projection)` first transforms world space to view space and then view
space to clip space. Pass this matrix—not its inverse—to
`RayFromNormalizedScreen`; the method performs the inversion:

```go
ray, ok := viewProjection.RayFromNormalizedScreen(math3d.V2(screenX, screenY))
if !ok {
	return errors.New("screen point could not be unprojected")
}
```

Screen coordinates are normalized with `(0, 0)` at the top left and `(1, 1)`
at the bottom right. Values outside that range extrapolate beyond the viewport.
Together, the conventions are right-handed negative-Z forward, row-vector
view-projection order, normalized top-left-origin screen coordinates, and clip
depth `[0, 1]`.

The returned direction is normalized, but the ray begins on the near plane,
not at `cameraPosition`. With a near distance of 1, a camera at `(0, 0, 5)`
looking at the origin produces a center-screen ray starting at `(0, 0, 4)` and
pointing along `(0, 0, -1)`. See executable
[`ExampleMat4_RayFromNormalizedScreen_worldSpace`](../example_test.go). The
existing [`ExampleMat4_RayFromNormalizedScreen`](../example_test.go) uses only
a projection matrix and therefore returns a view-space ray.

## Recover a hit point

All ray intersection methods return a parameter `t`. Recover the point in the
same space as the ray with:

```go
hit := ray.PointAt(t)
```

`t` is Euclidean distance only when `ray.Direction` has unit length. Screen
rays do have a unit direction. Manually constructed rays do not unless you
normalize their direction.

## Intersect a box

```go
t, ok := ray.IntersectBox(box)
if !ok {
	// Either the ray missed, the box was invalid, or the ray was non-finite
	// or had zero direction.
	return
}
hit := ray.PointAt(t)
```

Box contact is inclusive. A ray originating inside or on the box returns
`t == 0`; otherwise a successful result is the entry parameter. Reversed box
endpoints are invalid and return false.

## Intersect a sphere

```go
t, ok := ray.IntersectSphere(sphere)
if !ok {
	// Either the ray missed, the sphere was invalid, or the ray was non-finite
	// or had zero direction.
	return
}
hit := ray.PointAt(t)
```

The sphere must have a finite center and a non-negative finite radius. Its
boundary is included, and a ray originating inside or on it returns `t == 0`.
For an outside origin, the result is the first non-negative entry parameter.

## Intersect a plane

```go
t, ok := ray.IntersectPlane(plane, math3d.Tolerance)
if !ok {
	// Miss, invalid input, or no distinguished intersection.
	return
}
hit := ray.PointAt(t)
```

Planes use `dot(Normal, point) + D = 0`; the normal need not be unit length.
The tolerance is an absolute threshold on the direction-normal dot product and
also permits a hit no more than that tolerance behind the origin, clamping it
to `t == 0`. It must be non-negative. Parallel rays return false, including a
ray coplanar with the plane, because a coplanar ray has no single distinguished
hit. Non-finite rays or planes and zero-normal planes also fail.

## Intersect a triangle

```go
t, ok := ray.IntersectTriangle(triangle, math3d.Tolerance)
if !ok {
	// Miss, invalid input, parallelism, or degeneracy.
	return
}
hit := ray.PointAt(t)
```

Both triangle windings are accepted, and boundary hits on edges or vertices
are included. The tolerance must be finite and non-negative; it is an absolute
parallelism threshold and the exclusive lower bound for `t`. Consequently,
hits at the ray origin or within tolerance of it are excluded. Non-finite
inputs, zero-direction rays, degenerate or parallel triangles, and ordinary
misses all return false.

Executable [`ExampleRay_IntersectTriangle`](../example_test.go) deliberately
uses a direction of length 2. Its hit has `t == 0.5` but lies one Euclidean unit
from the origin, demonstrating why `PointAt(t)` is the reliable way to recover
the intersection.

An `ok == false` result is method-specific. Do not interpret every false result
as “no intersection”; it can also report invalid input or a failed mathematical
precondition as described above and in the declaration godoc.
