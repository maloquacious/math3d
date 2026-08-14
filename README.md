# math3d

`math3d` is a dependency-free Go package for 3D vectors, rotations, transforms, geometry, bounds, and spatial queries.

The current package version is `0.1.0` and requires Go 1.22.

## Install

```sh
go get github.com/maloquacious/math3d
```

## 60-second quick start

Start with `Transform` for a rigid rotation and translation:

```go
package main

import (
	"fmt"

	"github.com/maloquacious/math3d"
)

func main() {
	transform := math3d.NewTransform(
		math3d.V3(10, 0, 0),
		math3d.QuatRotationZ(math3d.HalfPi),
	)
	point := transform.TransformPoint(math3d.V3(1, 0, 0))

	fmt.Printf("%.0f %.0f %.0f\n", point.X, point.Y, point.Z)
}
```

The result is `10 1 0`: the transform rotates the point 90 degrees around Z,
then translates it 10 units along X. This program is maintained as the
executable [`ExampleTransform`](example_test.go).

## Before using matrices or intersections

- Rotation APIs use **radians** unless their declaration says otherwise.
- Matrices use **row vectors**. `A.Mul(B)` applies A first and B second.
- A returned `(value, ok)` must be checked. `ok == false` can mean an ordinary
  miss, invalid input, degeneracy, singularity, or another failed mathematical
  precondition; consult the method's godoc.
- The primary API uses **`float32`**. The available `D` types use `float64` but
  do not provide a complete parallel API; there is no `DMat4`.
- An intersection result `t` is a **ray parameter**, not necessarily a distance.
  It is a Euclidean distance only when the ray direction has unit length.

## Capabilities

### Vectors and rotations

- 2D, 3D, and 4D vectors; arithmetic, normalization, products, projections,
  reflection, angle tests, and interpolation
- Quaternions, axis-angle and Euler records, rotation construction,
  interpolation, and vector rotation

### Transforms and matrices

- Rigid `Transform` values for rotation and translation
- General 4×4 matrices for scale, rotation, translation, camera/view,
  projection, reflection, shadow, inversion, decomposition, and composition
- Separate point, direction, and normal transformations

### Geometry queries

- Ray intersections with boxes, spheres, planes, and triangles
- Point and shape containment, overlap, distance, plane classification, and
  screen-to-ray unprojection
- Segment and triangle measurements and vector relationship tests

### Bounds and shapes

- Inclusive intervals and axis-aligned 2D, 3D, and 4D boxes
- Spheres, planes, rays, segments, triangles, and quads
- Bound construction, expansion, merging, intersection, and transformation

### Scalar and statistics helpers

- Angle conversion and wrapping, interpolation, deterministic random values,
  hashing, percentages, and powers of two
- Component-wise summary statistics for vector collections

### Available float64 types

- `DVec2`, `DVec3`, `DVec4`, `DQuat`, `DPlane`, `DRay`, and `DSphere`
- `DInterval`, `DBox2`, `DBox3`, and `DBox4`
- Float64 coordinate records and selected scalar records and helpers

## Guides

- [Place and query a 3D object](docs/getting-started.md) — a guided first
  workflow using a rigid transform, world-space bounds, and a ray query
- [Transform points, directions, and rotations](docs/transforms.md) — choose
  rigid or general transforms and apply matrix and quaternion ordering correctly
- [Pick and intersect geometry](docs/picking-and-intersections.md) — construct
  world-space screen rays and interpret shape-specific intersection results

## API documentation

The [package documentation on pkg.go.dev](https://pkg.go.dev/github.com/maloquacious/math3d)
and each exported declaration's godoc are the authoritative API reference. See
the [executable examples on pkg.go.dev](https://pkg.go.dev/github.com/maloquacious/math3d#pkg-examples)
or their [source in `example_test.go`](example_test.go).
