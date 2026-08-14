# Place and query a 3D object

In this tutorial, you will place a square in world space, compute its
axis-aligned bounds, and cast a ray into those bounds. The final program prints
the transformed bounds and the first hit point.

## Define the object in local space

Start with four vertices centered on the object's local origin:

```go
localVertices := []math3d.Vec3{
	math3d.V3(-1, -1, 0),
	math3d.V3(1, -1, 0),
	math3d.V3(1, 1, 0),
	math3d.V3(-1, 1, 0),
}
```

These values are object-local positions. They do not yet describe where the
object belongs in the world.

## Place the object in world space

Create a unit quaternion that rotates 90 degrees around Z, then pair it with a
world-space position in a rigid `Transform`:

```go
transform := math3d.NewTransform(
	math3d.V3(5, 0, 0),
	math3d.QuatRotationZ(math3d.HalfPi),
)
worldVertices := make([]math3d.Vec3, len(localVertices))
for i, vertex := range localVertices {
	worldVertices[i] = transform.TransformPoint(vertex)
}
```

`TransformPoint` rotates each local position and then translates it. The
resulting vertices and all bounds built from them are in world space.

## Build world-space bounds

Build an axis-aligned box around the transformed vertices:

```go
worldBounds, ok := math3d.Box3FromPoints(worldVertices)
if !ok {
	fmt.Println("could not build world bounds")
	return
}
```

`Box3FromPoints` returns false for an empty slice or a non-finite point. This is
a failed input precondition, not an intersection miss.

## Cast a query ray

Normalize the query direction before constructing the ray, and handle the case
where normalization is impossible:

```go
direction, ok := math3d.V3(2, 0, 0).Normalized()
if !ok {
	fmt.Println("invalid query direction")
	return
}
ray := math3d.NewRay(math3d.V3(0, 0, 0), direction)
t, ok := ray.IntersectBox(worldBounds)
if !ok {
	fmt.Println("ray missed the object bounds")
	return
}
hit := ray.PointAt(t)
```

The ray origin, direction, bounds, and recovered hit are all in world space.
Here `direction` has unit length, so the returned ray parameter `t` is also the
Euclidean distance from the ray origin. For an unnormalized direction, `t`
would only be a parameter; `PointAt(t)` would still recover the hit position.

## Run the complete program

```go
package main

import (
	"fmt"

	"github.com/maloquacious/math3d"
)

func main() {
	localVertices := []math3d.Vec3{
		math3d.V3(-1, -1, 0),
		math3d.V3(1, -1, 0),
		math3d.V3(1, 1, 0),
		math3d.V3(-1, 1, 0),
	}

	transform := math3d.NewTransform(
		math3d.V3(5, 0, 0),
		math3d.QuatRotationZ(math3d.HalfPi),
	)
	worldVertices := make([]math3d.Vec3, len(localVertices))
	for i, vertex := range localVertices {
		worldVertices[i] = transform.TransformPoint(vertex)
	}

	worldBounds, ok := math3d.Box3FromPoints(worldVertices)
	if !ok {
		fmt.Println("could not build world bounds")
		return
	}

	direction, ok := math3d.V3(2, 0, 0).Normalized()
	if !ok {
		fmt.Println("invalid query direction")
		return
	}
	ray := math3d.NewRay(math3d.V3(0, 0, 0), direction)
	t, ok := ray.IntersectBox(worldBounds)
	if !ok {
		fmt.Println("ray missed the object bounds")
		return
	}
	hit := ray.PointAt(t)

	fmt.Printf("bounds: (%.0f, %.0f, %.0f) to (%.0f, %.0f, %.0f)\n",
		worldBounds.Min.X, worldBounds.Min.Y, worldBounds.Min.Z,
		worldBounds.Max.X, worldBounds.Max.Y, worldBounds.Max.Z)
	fmt.Printf("hit: (%.0f, %.0f, %.0f), distance: %.0f\n", hit.X, hit.Y, hit.Z, t)
}
```

The output is:

```text
bounds: (4, -1, 0) to (6, 1, 0)
hit: (4, 0, 0), distance: 4
```

The complete listing is maintained as the package-level executable
[`Example`](../example_test.go).

Use `Mat4` instead of `Transform` when you need scale, projection, shear, or
general affine composition.
