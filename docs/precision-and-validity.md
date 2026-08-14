# Choose precision and validate representations

The package's precision families are intentionally asymmetric. Choose a family
for the operations you need, convert explicitly, and distinguish a value that
can be represented from one that satisfies a geometric operation's
preconditions.

## Choose `float32` or an available `float64` type

The primary vectors, quaternions, matrices, transforms, rays, planes, bounds,
and shapes use `float32`. This is the complete transformation and geometry
pipeline.

The available `float64` types are:

- `DVec2`, `DVec3`, and `DVec4`;
- `DQuat` and `DPlane`;
- `DRay` and `DSphere`;
- `DInterval`, `DBox2`, `DBox3`, and `DBox4`.

These types do not form a complete parallel API. In particular, there is no
`DMat4` or double-precision transform pipeline. Some D types also expose fewer
operations than their float32 counterparts. Choose the primary API when your
workflow requires matrices, transforms, projection, or the full set of shape
queries.

`AxisAngle` is a named exception: its axis and angle use `float64` despite the
absence of a `D` prefix. Coordinate records and selected scalar helpers also
store float64 values, but they are records and utilities rather than a complete
double-precision geometry pipeline.

## Convert precision explicitly

Conversions are methods rather than implicit coercions:

```go
wide := point.DVec3() // float32 to float64
narrow := wide.Vec3() // float64 to float32
```

Widening preserves the exact float32 value but cannot restore information that
was already rounded away. Narrowing can round a value, underflow it, or overflow
it to an infinity. Executable [`ExampleDVec3`](../example_test.go) narrows the
exact float64 integer `16777217` to `16777216`, overflows `1e40`, and shows that
widening the result does not recover the original value.

The same narrowing concern applies to `DQuat.Quat`, `DBox3.Box3`, and the D
vector conversion methods. Check range and finiteness when conversion loss
would make the result invalid for the next operation.

## Select an absolute tolerance

`AlmostEqual` methods compare stored components with a strict absolute test:

```text
abs(a - b) < tolerance
```

A difference exactly equal to the tolerance is not equal. The comparison is
not relative and is not measured in ULPs, so choose a tolerance in the units
and scale of the values being compared. Non-positive tolerances do not make
equal finite representations compare approximately equal.

These methods compare representation, not every possible semantic
equivalence. In particular, quaternion `q` and `-q` can encode the same rotation
but are not representation-wise `AlmostEqual`.

## Separate construction from validity

Constructors generally preserve caller input. They do not silently repair or
canonicalize it:

- `NewQuat`, `NewDQuat`, and `NewTransform` do not normalize quaternions;
- `NewRay` and `NewDRay` do not normalize directions;
- `NewPlane` and `NewDPlane` do not normalize plane normals;
- `NewSphere` and `NewDSphere` preserve negative radii;
- bound constructors preserve reversed endpoints;
- `NewAxisAngle` does not normalize its axis.

This makes those values valid *representations*: their fields retain what the
caller supplied. It does not necessarily make them semantically valid for a
geometric operation. For example, a negative-radius sphere is representable
but `Sphere.Valid` reports false; a zero-direction ray is representable but
does not define a geometric ray; and a non-unit quaternion is storable but does
not satisfy APIs that require a rotation quaternion.

Use `Valid`, normalization methods, and operation-specific `(value, ok)`
results at the boundary where a stronger mathematical precondition matters.

## Know what each zero value means

Zero values are safe to store and pass, but identity, validity, and degeneracy
depend on the family:

- **Vectors and intervals:** zero vectors are finite zero vectors. Zero
  `Interval` and `DInterval` values are the valid point interval `[0, 0]`, not
  empty intervals.
- **Boxes and spheres:** zero boxes are valid point boxes at the origin. Zero
  spheres are valid radius-zero spheres at the origin. Use `EmptyBox*` or
  `EmptyDBox*` when an empty accumulation identity is required.
- **Matrices, quaternions, and transforms:** zero matrices are zero matrices,
  not identity. Zero quaternions do not define rotations, so zero transforms do
  not define rigid transforms. Use `IdentityMat4`, `IdentityQuat`,
  `IdentityDQuat`, and `IdentityTransform`.
- **Rays and planes:** zero rays have zero direction and do not define a
  geometric ray. Zero planes have zero normal and do not define a geometric
  plane.
- **Segments, triangles, and quads:** their zero values collapse all vertices
  to the origin and are degenerate.
- **`Stats`:** zero `Stats[T]` has `Count == 0` and zero `Min`, `Max`, and `Sum`.
  `StatsAverage` returns false when `Count` is not positive.

## Do not infer coordinate-record angle units

Rotation APIs use radians unless their declarations say otherwise. That rule
does not establish units for coordinate-record fields whose units are not
specified. In particular, do not assume units for angles in `Spherical`,
`Polar`, `LogPolar`, `Cylindrical`, `Horizontal`, or `Geo`. Treat those values as
records until an application-level conversion establishes a unit contract.

`Euler` and `AxisAngle` explicitly document radians and therefore are not
ambiguous in this way.
