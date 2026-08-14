# End-user documentation implementation prompt

Document `github.com/maloquacious/math3d` for Go developers using the package
for rendering, CAD/BIM, simulation, spatial queries, and geometry processing.
Assume readers understand ordinary vector math but do not know this package's
conventions.

Implement the numbered tasks below **in order**. Each task must start from the
completed result of the previous task and leave the repository in a reviewable,
passing state. Do not create empty documentation scaffolding for later tasks.

## Documentation ownership

Keep each kind of information in one authoritative place:

- `README.md` owns discovery, installation, first success, a capability map,
  critical warnings, and navigation.
- `doc.go` owns package-wide normative conventions.
- Exported declaration comments own exact API behavior, units, preconditions,
  failure conditions, tolerance semantics, and zero-value behavior.
- `example_test.go` owns canonical, compilable examples rendered by pkg.go.dev.
- `docs/getting-started.md` owns the guided introductory learning experience.
- Focused files under `docs/` own task-oriented, cross-API how-to guides.
- `UPSTREAM.md`, `PORTING_PLAN.md`, and `REMAINING_PORTING_PROMPT.md` remain
  maintainer-only records.

Do not create a generated documentation site or a Markdown copy of the API
reference. Link to godoc instead of duplicating declaration documentation.
Every complete Go listing in Markdown must have an executable counterpart in
`example_test.go` so `go test` detects drift.

## Package-wide risks to address accurately

Documentation must not obscure these behaviors:

1. The package is right-handed, local forward is negative Z, and rotations use
   radians unless a declaration explicitly says otherwise.
2. Matrices use row vectors. `A.Mul(B)` applies A and then B, and translation is
   stored in `M41`, `M42`, and `M43`.
3. Quaternion Hamilton multiplication has a different ordering implication:
   `q.Mul(other)` applies `other` first. Use `Concatenate(first, second)` when
   communicating sequential rotations.
4. Points, directions, and normals require different transformations. Normals
   use the inverse transpose and are not normalized afterward.
5. An `ok` result can report an ordinary miss, degeneracy, invalid input,
   singularity, or another failed mathematical precondition. Do not describe
   every false result as “no intersection.”
6. Constructors generally preserve the supplied representation. They do not
   silently normalize rays, quaternions, axes, or planes, and negative sphere
   radii remain representable.
7. Zero does not mean identity for matrices, quaternions, or transforms. Use
   `IdentityMat4`, `IdentityQuat`, and `IdentityTransform`.
8. The primary API uses `float32`. `D` types use `float64`, but they do not form
   a complete parallel API; in particular, there is no `DMat4`.
9. A ray intersection result `t` is a ray parameter and is a Euclidean distance
   only when the direction has unit length.
10. Approximate comparisons use strict absolute tolerances. Quaternion `q` and
    `-q` may represent the same rotation but are not representation-wise
    `AlmostEqual`.
11. Projection uses right-handed negative-Z forward, depth `[0, 1]`, normalized
    top-left-origin screen coordinates, and row-vector view-projection order.
12. Intersection boundary behavior differs by shape. Document the behavior of
    the method being demonstrated rather than implying a universal rule.

## Task 1 — Add the end-user README

Create `README.md` as the package's concise entry point.

Include:

1. A one-sentence description and dependency-free positioning.
2. Current status without inventing compatibility guarantees: Go 1.22 and
   package version `0.1.0`.
3. Installation:

   ```sh
   go get github.com/maloquacious/math3d
   ```

4. A 60-second quick start based on the existing executable
   `ExampleTransform`. Start with rigid `Transform`, not raw matrices.
5. A short “before using matrices or intersections” section that makes radians,
   row-vector multiplication order, `(value, ok)`, primary `float32` precision,
   and ray-parameter semantics visible before the capability map.
6. A capability map grouped by vectors and rotations, transforms and matrices,
   geometry queries, bounds and shapes, scalar/statistics helpers, and available
   float64 types.
7. Links to the package on pkg.go.dev and its executable examples. Do not link
   to documentation files that have not been created yet.
8. A statement that package and declaration godoc are the authoritative API
   reference.

Do not include upstream parity tables, port milestones, implementation history,
or an exhaustive API inventory.

**Acceptance criteria**

- A newcomer can decide whether the package is suitable, install it, and obtain
  a transformed point without opening a maintainer document.
- All complete snippets correspond to an executable example.
- All links resolve in the GitHub rendering.
- The README remains an entry point rather than a second API reference.

## Task 2 — Add the getting-started tutorial

Create `docs/getting-started.md` with the outcome-oriented title “Place and
query a 3D object.” Add a package-level executable `Example` in
`example_test.go` containing the tutorial's complete program.

Guide the reader through one coherent scenario:

1. Define local object vertices with `V3`.
2. Create a unit rotation with `QuatRotationZ` and a rigid transform with
   `NewTransform`.
3. Apply `TransformPoint` to obtain world-space vertices.
4. Build world bounds with `Box3FromPoints`.
5. Normalize a query direction and handle `ok`.
6. Construct a `Ray`, call `IntersectBox`, and recover the hit with `PointAt`.
7. Connect the normalized direction to the fact that `t` is a Euclidean
   distance in this example.
8. Point readers to matrices when they require scale, projection, shear, or
   general affine composition.

Keep the tutorial focused on a reliable 10–15 minute learning experience.
Move detailed API descriptions to declaration comments rather than interrupting
the workflow. Update `README.md` to link to the completed tutorial.

**Acceptance criteria**

- The tutorial produces a visible, deterministic result.
- The full listing is executable through `go test`.
- At least one failure branch is handled and described accurately.
- Space changes and ray `t` semantics are unambiguous.

## Task 3 — Document transformations and rotations

Create `docs/transforms.md` as a task-oriented guide. Cover:

1. Choosing rigid `Transform` versus general `Mat4`.
2. The fact that `NewTransform` does not normalize its orientation.
3. Building scale-rotation-translation with `ComposeMat4`.
4. Chaining `Mat4.Mul` in application order, with a numerical result that would
   change if the order were reversed.
5. Quaternion multiplication order and when to use `Concatenate`.
6. The distinct meanings of `TransformPoint`, `TransformDirection`, and
   `TransformNormal`.
7. Normal transformation by inverse transpose and the need to normalize the
   result when a unit normal is required.
8. Handling failed matrix inversion and decomposition.

Add executable examples:

- `ExampleConcatenate`
- `ExampleMat4_TransformNormal`

Audit the godoc of only the API families touched by this guide. Fix missing or
ambiguous declaration comments at their source rather than compensating in the
guide. Update the README navigation after the guide exists.

**Acceptance criteria**

- Matrix and quaternion ordering are contrasted explicitly.
- Numerical examples detect accidental ordering reversals.
- Point, direction, and normal behavior cannot be mistaken for one another.
- Unit-quaternion and fallibility requirements are stated where relevant.

## Task 4 — Document picking and intersections

Create `docs/picking-and-intersections.md`. Show users how to:

1. Build view and projection matrices with `LookAtMat4` and
   `PerspectiveFOVMat4`.
2. Compose a world-space unprojection matrix as `view.Mul(projection)`.
3. Use `RayFromNormalizedScreen` with top-left-origin normalized coordinates,
   depth `[0, 1]`, and the correct handedness.
4. Account for the returned ray starting on the near plane rather than at the
   camera.
5. Query boxes, spheres, planes, and triangles.
6. Use `PointAt(t)` and distinguish a parameter from a distance.
7. Interpret method-specific `ok` behavior without conflating misses and invalid
   mathematical inputs.
8. Account for relevant boundary behavior, including inside-origin results,
   parallel or coplanar plane rays, and the triangle method's exclusion of hits
   at or within tolerance of the ray origin.

Add executable examples:

- `ExampleMat4_RayFromNormalizedScreen_worldSpace`
- `ExampleRay_IntersectTriangle`

Keep the existing projection-only screen-ray example because it demonstrates
view-space behavior. Audit and improve only the projection, ray, and
intersection declaration comments touched by this guide. Update README
navigation.

**Acceptance criteria**

- The world-space example uses `view.Mul(projection)`.
- Screen, clip-depth, handedness, and near-plane-origin conventions appear
  together.
- Every demonstrated intersection states its relevant miss, validity,
  tolerance, origin, and `t` semantics.
- All examples have deterministic output.

## Task 5 — Document precision and validity

Create `docs/precision-and-validity.md` as a concise explanation and practical
selection guide.

Cover:

1. The primary `float32` API and the actual, asymmetric `D`/`float64` coverage.
2. Explicit widening and narrowing conversions and possible range or precision
   loss.
3. The absence of a double-precision matrix and transform pipeline.
4. `AxisAngle` as a documented float64 exception.
5. Strict absolute tolerance behavior rather than relative or ULP comparison.
6. Constructor preservation versus geometric validity.
7. Zero-value behavior by family:
   - vectors and intervals;
   - boxes and spheres;
   - matrices, quaternions, and transforms;
   - rays and planes;
   - segments, triangles, and quads;
   - `Stats`.
8. Coordinate records whose units are not established upstream; do not imply
   that every coordinate-record angle inherits the package's radians rule.

Add an executable `ExampleDVec3` showing explicit precision conversion. Audit
the related D-type, conversion, validity, and tolerance comments. Update README
navigation.

**Acceptance criteria**

- The guide does not imply complete float64 parity.
- Every zero-value statement agrees with current declaration comments and
  behavior.
- Representation validity and semantic validity are clearly distinguished.
- The conversion example documents narrowing rather than presenting it as
  lossless.

## Task 6 — Document bounds and shape workflows

Create `docs/bounds-and-shapes.md` as a collection of focused tasks:

1. Build bounds from a point slice with `Box3FromPoints`.
2. Accumulate bounds incrementally with `EmptyBox3().Expanded(point)`.
3. Use `Valid`, `Empty`, reversed endpoints, inclusive contact, `Merge`, and
   fallible `Intersection` correctly.
4. Explain that transforming an axis-aligned box returns axis-aligned bounds
   around transformed corners.
5. Use segment parameters, noting that `[0, 1]` lies on the segment and values
   outside that range extrapolate.
6. Work with triangle winding, signed area, degenerate normals, and the strict
   interior behavior of `Triangle2.Contains`.
7. Use the plane equation `dot(normal, point) + D = 0`, including non-unit plane
   classification and projection behavior.

Add executable examples:

- `ExampleEmptyBox3`
- `ExampleBox3_Intersection`
- `ExamplePlane_Project`

Audit the touched bounds, shape, and plane declaration comments. Update README
navigation.

**Acceptance criteria**

- Empty and invalid bound representations are not conflated with a point box.
- Inclusive versus strict boundary behavior is stated per operation.
- Winding, degeneracy, segment parameters, and plane normalization are accurate.
- Examples exercise ordinary workflows rather than merely listing methods.

## Task 7 — Complete the initial documentation pass

Finish the first end-user documentation release without creating low-value
pages for record-only APIs.

1. Add `ExampleStatsOf` covering component-wise minimum, maximum, average, and
   empty-average failure.
2. Review the README capability map against the current exported API and all
   completed guides.
3. Review all Markdown links and ensure each complete code listing has an
   executable counterpart.
4. Run a focused godoc audit for declarations used by any new example but not
   already audited in Tasks 2–6.
5. Keep scalar interpolation, coordinates, colors, and record types in the
   capability map and godoc unless a genuine end-user workflow justifies a
   dedicated guide.
6. Ensure maintainer history and parity details remain in the maintainer files
   and are not presented as onboarding material.

**Acceptance criteria**

- The README provides a coherent route from first use to every completed guide.
- No generated or handwritten duplicate API reference exists.
- All canonical snippets compile and execute.
- Lower-frequency APIs are discoverable without overwhelming the main learning
  path.

## Verification required after every task

Run:

```sh
go test ./...
go test -run '^Example' ./...
go vet ./...
go doc .
```

When Go examples change, also run:

```sh
gofmt -d example_test.go
```

When declaration comments change, inspect the rendered text for every modified
symbol with commands such as:

```sh
go doc . Mat4.Mul
go doc . Mat4.TransformNormal
go doc . Mat4.RayFromNormalizedScreen
go doc . Ray.IntersectSphere
```

Render changed Markdown and inspect headings, code blocks, tables, and links.
Do not claim pkg.go.dev rendering has been verified until a version containing
the documentation has actually been published.
