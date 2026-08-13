# Plan: port `vimaec/math3d` to idiomatic Go

## Goal

Port the public behavior of [`vimaec/math3d`](https://github.com/vimaec/math3d) to a dependency-free Go package, preserving its float precision and geometric conventions while making its data model immutable from outside the package.

The first milestone is deliberately structural: define, document, and test the value types before porting the larger algorithms.

## Decisions to make explicit

These are the proposed defaults. Record them in package documentation before implementation so later ports do not accidentally change conventions.

1. **Preserve source precision.** The primary types use `float32`; `D`-prefixed types use `float64`. Do not replace the primary API with `float64` or expose generic numeric types. This preserves source behavior, data size, and straightforward test translation.
2. **Use immutable value types.** Store components in unexported fields, expose constructors and accessors, and give operations value receivers that return new values. Do not expose pointers from values or mutating methods.
3. **Make zero values safe where possible, not necessarily meaningful.** Zero vectors and scalar intervals are valid. A zero quaternion, ray, plane, sphere, transform, or box may not meet an algorithm's preconditions; document those cases rather than hiding them with implicit initialization.
4. **Preserve row-vector matrix semantics.** Points transform as `[x y z 1] * M`, normals as `[x y z 0] * M`, and translation occupies `M41`, `M42`, `M43`. Composition order must match the source.
5. **Prefer normal Go result forms.** Use `(value, ok)` for intersections, normalization of zero-length values, matrix inversion, and other operations that can legitimately fail. Reserve errors for malformed external input, not ordinary geometric misses.
6. **Do not reproduce C#-shaped interfaces.** Add small interfaces only when a Go algorithm has multiple real implementations. Use concrete methods for vectors and shapes.
7. **Do not serialize raw struct memory.** Supply explicit binary encoding in documented component order and endianness if binary compatibility is required. Private representation remains free to change.
8. **Pin the source revision.** Before implementation, record the exact upstream commit in `UPSTREAM.md`. The repository is a mirror and its `dev` branch can change.

## Immutable representation

Go has no `readonly struct`. Exported fields would undermine the stated immutability goal, so components should be private:

```go
type Vec3 struct {
	x, y, z float32
}

func V3(x, y, z float32) Vec3 { return Vec3{x: x, y: y, z: z} }
func (v Vec3) X() float32      { return v.x }
func (v Vec3) Y() float32      { return v.y }
func (v Vec3) Z() float32      { return v.z }
func (v Vec3) WithX(x float32) Vec3 { return V3(x, v.y, v.z) }
```

Consequences:

- Assignment and method calls copy values; callers cannot mutate package-owned state.
- Types containing only comparable fields remain usable as map keys and with `==`.
- Package-level vector and matrix variables must be avoided because callers can reassign them. Use scalar constants and value-returning functions such as `IdentityQuat()` and `IdentityMat4()`.
- JSON, text, and binary support needs explicit marshal methods because fields are private.
- Algorithms inside the package may construct values directly, but public checked constructors should distinguish preconditions from raw representation.

`Mat4` should use a private `[16]float32` in source flattening order rather than sixteen mutable public fields. Expose `At(row, col)`, `Row`, `Column`, and semantic constructors. The array keeps `Mat4` comparable and contiguous without making raw in-memory serialization part of the API.

## Public type map

Use concise Go names for common mathematical types while documenting the upstream name in each type's doc comment.

| Upstream | Proposed Go type | Storage |
|---|---|---|
| `Vector2/3/4` | `Vec2`, `Vec3`, `Vec4` | `float32` components |
| `DVector2/3/4` | `DVec2`, `DVec3`, `DVec4` | `float64` components |
| `Int2/3/4` | `IVec2`, `IVec3`, `IVec4` | `int32` components, preserving C# width |
| `Byte2/3/4` | `Byte2`, `Byte3`, `Byte4` | `byte` components |
| `Complex` | `Complex` | two `float64` components; evaluate interoperability with built-in `complex128` |
| `Quaternion`, `DQuaternion` | `Quat`, `DQuat` | four components |
| `Matrix4x4` | `Mat4` | private `[16]float32` |
| `Transform` | `Transform` | `Vec3` position and `Quat` orientation |
| `Plane`, `DPlane` | `Plane`, `DPlane` | normal and scalar offset |
| `Interval`, `DInterval` | `Interval`, `DInterval` | min and max |
| `AABox2D`, `AABox`, `AABox4D` | `Box2`, `Box3`, `Box4` | min and max vectors |
| `DAABox2D`, `DAABox`, `DAABox4D` | `DBox2`, `DBox3`, `DBox4` | min and max vectors |
| `Ray`, `DRay` | `Ray`, `DRay` | origin and direction |
| `Sphere`, `DSphere` | `Sphere`, `DSphere` | center and radius |
| `Line`, `Line2D` | `Segment3`, `Segment2` | endpoints |
| `Triangle`, `Triangle2D` | `Triangle3`, `Triangle2` | vertices |
| `Quad`, `Quad2D` | `Quad3`, `Quad2` | vertices |
| color records | `RGB`, `RGBA`, `HDRColor` | matching byte/float fields |
| coordinate records | names without `Coordinate`, e.g. `Spherical`, `Geo` | source `float64` fields |
| `AxisAngle`, `Euler`, `HorizontalCoordinate` | `AxisAngle`, `Euler`, `Horizontal` | matching source precision |
| `Stats<T>` | `Stats[T]` | immutable count/min/max/sum snapshot |
| `ValueDomain` | `Domain` | two `float64` bounds |

Before freezing the API, create a short naming prototype and run `golint`-equivalent documentation checks plus a small consumer example. If source-name familiarity is more important to intended users, retain `Vector3`, `Matrix4`, and `Quaternion`; do not provide aliases for both naming schemes because that doubles the apparent API.

The current README's `LinearMotion`, `AngularMotion`, `Motion`, `Triangle2`, and `DQuad` entries are stale and have no corresponding current source declarations; do not invent them. Treat the `Vim.Experimental` interfaces and empty `PrimitiveShapes` placeholders as explicitly out of scope for the stable port. Track non-structural helpers (`Stats[T]`, `Domain`, stateless random functions, and public math helpers) in the parity ledger and port them after the core geometry they depend on.

## Structure milestone

### 1. Establish the module and parity ledger

- Create `go.mod` after the module import path is known.
- Add package documentation containing handedness, row-vector convention, angle units, precision policy, tolerance policy, and immutability guarantees.
- Add `UPSTREAM.md` with the pinned commit, source file inventory, source-to-Go symbol map, and intentional deviations.
- Use one package initially. Split packages only if an actual import boundary emerges; these types are highly interconnected.

**Exit criteria:** `go test ./...` runs, upstream is pinned, and every upstream public type is listed as planned, deferred, or intentionally omitted.

### 2. Define scalar policy and simple records

- Port mathematical constants as typed constants where Go permits it.
- Define `Containment` and `PlaneIntersection` enums with an explicit zero value.
- Define byte vectors, colors, integer vectors, and coordinate records.
- Add exact equality naturally through comparable values; add `AlmostEqual` only to floating-point types.
- Decide how NaN affects equality, ordering, hashing/map-key use, and approximate comparisons. Match Go's normal IEEE behavior unless source tests require otherwise.

**Exit criteria:** constructors/accessors/`With…` methods and table-driven tests exist for every simple record; `go vet ./...` is clean.

### 3. Define vector values

- Implement `Vec2/3/4` and `DVec2/3/4` manually first.
- Provide named constructors (`V2`, `V3`, `V4`) and splat constructors only where useful.
- Port component access, `WithX`-style replacement, arithmetic methods, dot product, magnitude, component reductions, min/max, finite/NaN checks, and approximate equality.
- Use `Normalized() (VecN, bool)`; do not let zero normalization silently create NaNs.
- Keep component-wise multiplication and scalar multiplication visibly distinct (`Mul` and `Scale`).
- Add precision conversions explicitly (`Vec3.DVec3()` or a clearly documented constructor), never implicitly.

**Exit criteria:** translated vector tests pass; property/fuzz tests cover commutativity where valid, dot symmetry, cross orthogonality, normalization, and float edge cases.

### 4. Define compound geometry records without algorithms

Define immutable storage, constructors, accessors, replacement methods, equality, and documented invariants for:

- `Quat`, `DQuat`, `AxisAngle`, and `Euler`
- `Mat4`
- `Transform`
- `Plane` and `DPlane`
- intervals and all box dimensions/precisions
- rays and spheres
- segments, triangles, and quads

Use two constructor levels only when they have distinct semantics:

- A direct representation constructor preserves source expressiveness, e.g. `NewRay(origin, direction)`.
- A checked semantic constructor returns `(T, bool)`, e.g. `UnitRay(origin, direction)` or `SphereFromCenterRadius` if negative radii are rejected.

Do not silently normalize quaternions, ray directions, plane normals, or box endpoint order in direct constructors. Such normalization changes values and can conceal mistakes.

**Exit criteria:** every value is externally immutable and comparable; all invariants and zero-value behavior are documented; representation tests cover component order and sizes of explicit encodings.

## Algorithm milestones

### 5. Quaternion and matrix operations

- Port quaternion algebra, normalization, interpolation, axis-angle/Euler creation, and vector rotation.
- Port matrix factories, multiplication, transpose, determinant, inversion, decomposition, and projection/view matrices.
- Keep matrix/quaternion conversions together in the same package to avoid artificial dependency layers.
- Translate source tests before adding convenience APIs.

**Critical tests:** translation slots, row-vector composition order, quaternion multiplication order, quaternion/matrix round trips, singular inversion, projection handedness, and reflection detection.

### 6. Planes and transforms

- Standardize on the plane equation `dot(normal, point) + d = 0`.
- Correct the upstream inconsistency where projection and one factory use the opposite sign; record this as an intentional deviation with regression tests.
- Port rigid `Transform` behavior without adding scale to the type.
- Return failure from transforms requiring a matrix inverse when inversion fails.

### 7. Intervals, boxes, rays, and spheres

- Preserve invalid `min > max` as the empty-box representation, but provide `Valid()`/`Empty()` and clear constructors.
- Return `(intersection, ok)` for interval/box intersections rather than relying only on an invalid sentinel.
- Port box containment and corner-based transforms.
- Port ray intersections as `(t, ok)`. State whether `t` is a ray parameter or Euclidean distance; only call it distance when direction is unit length.
- Fix or explicitly preserve the source sphere-intersection assumption that ray directions are normalized. Prefer a general formula that divides by `dot(direction, direction)` and test both unit and non-unit directions.
- Port sphere containment, merging, and conservative non-uniform transforms.

### 8. Segments and polygonal shapes

- Port segment mapping and 2D segment intersection.
- Port triangle area, normal, degeneracy/sliver checks, bounds, and ray intersection.
- Port quad mapping and transforms.
- Use slices only at algorithm boundaries; never retain caller-owned slices inside immutable values.

### 9. Remaining scalar/vector math and conversions

- Port source helpers that are demonstrably used or publicly valuable.
- Port immutable `Stats[T]` and `Domain` values, then the stateless random and hashing helpers where they remain useful in Go. Use the standard library instead of copying LINQ-specific utilities.
- Prefer package functions for operations with no natural receiver and methods for operations centered on one value.
- Correct the misleading upstream unit-conversion comments/names and test numerical direction (`mm → ft`, `ft → mm`).
- Avoid mechanically recreating every C# extension method when Go's `math` package is clearer.

### 10. Encoding, compatibility, and release readiness

- Add explicit little-endian binary encoding only if VIM data interchange needs it; otherwise defer it.
- Add `encoding.TextMarshaler`/`TextUnmarshaler` and JSON support only from concrete consumer requirements.
- Benchmark hot operations without replacing readable code prematurely.
- Run `go test ./...`, `go test -race ./...`, `go vet ./...`, fuzz/property tests, and architecture-specific encoding tests.
- Produce API examples for vectors, transforms, matrix composition, bounds, and ray intersections.

## Testing strategy

1. **Translate upstream tests first.** Keep one `_test.go` file per major type and preserve source edge cases.
2. **Create convention tests before implementation.** These should fail if matrix order, translation placement, plane sign, angles, or quaternion composition changes.
3. **Use tolerances deliberately.** Approximate checks should accept an explicit tolerance; package defaults are conveniences, not hidden global state.
4. **Add property and fuzz tests.** Useful properties include round trips, inverse composition, normalized lengths, containment consistency, and intersection symmetry where applicable.
5. **Generate golden vectors from the pinned C# implementation.** Use them for difficult matrix, quaternion, and projection behavior. Store data, not a permanent .NET build dependency.
6. **Test intentional deviations separately.** Plane signs, non-unit ray/sphere intersection, failed inversion, and conversion naming must not be mistaken for accidental incompatibility.

## Porting rules

- Port behavior, not C# syntax: no `Get…` prefixes for simple accessors, no faux operator API, no broad `MathOps` namespace, and no `ITransformable` hierarchy.
- Keep methods short and formulas recognizable against upstream. Link each nontrivial port to its upstream file and pinned revision in a comment only when that aids verification.
- Do not add code generation until the handwritten vector and record APIs stabilize. If repetition later causes drift, generate implementation files from a checked-in Go template and keep generated output reviewable.
- Avoid global tolerance or coordinate-system settings; they create hidden mutable state.
- Prefer values in the public API. Use pointers only for decoders/builders or when benchmarks demonstrate that copying a specific large value is material.
- Keep compatibility fixes visible in `UPSTREAM.md`; never silently “clean up” formulas during translation.

## Initial delivery slices

Each slice should be independently reviewable and leave the package passing tests:

1. Module/package documentation, upstream pin, constants, enums, byte/color/integer/coordinate records.
2. `Vec2/3/4` and tests.
3. `DVec2/3/4`, precision conversions, and tests.
4. Immutable compound record definitions, invariant documentation, and representation tests.
5. Quaternion operations.
6. Matrix operations and vector transforms.
7. Plane and rigid transform behavior.
8. Intervals and boxes.
9. Rays and spheres.
10. Segments, triangles, quads, remaining conversions, encoding, and examples.

The API should not be declared stable until slice 6 proves that the immutable vector, quaternion, and matrix representations work ergonomically together in real transform code.
