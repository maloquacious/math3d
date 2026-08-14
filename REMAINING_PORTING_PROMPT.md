# Prompt: complete the remaining `vimaec/math3d` port

## Goal

Complete the remaining useful public behavior from
[`vimaec/math3d`](https://github.com/vimaec/math3d) in this dependency-free Go
package. Use the source pinned in `UPSTREAM.md` at commit
`4477350aaf223ed6f8b2985419eb8da3e58f96ed`; do not compare formulas against a
newer branch.

The core port and milestones 1–10 in `PORTING_PLAN.md` are already implemented.
This work is a focused parity pass, not a second port. Preserve the package's
existing API conventions, float precision, row-vector matrix semantics, plane
equation, and `(value, ok)` failure results.

## Working rules for every step

1. Read `PORTING_PLAN.md`, `UPSTREAM.md`, and the relevant Go implementation
   and tests before editing.
2. Inspect the pinned upstream method before translating it. Port behavior and
   formulas, not C# extension-method syntax, overloads, LINQ helpers, or broad
   interfaces.
3. Prefer methods on `Vec2`, `Vec3`, `Vec4`, `Quat`, and `Mat4` when an operation
   is centered on one value. Use a package function when no receiver is natural.
4. Reuse existing scalar functions and matrix/quaternion factories. Do not add
   wrappers that merely rename an existing Go operation.
5. Preserve primary `float32` precision. Use `float64` only where the existing
   type or calculation contract requires it.
6. Package operations must not mutate inputs. Report invalid mathematical
   preconditions consistently with the existing API instead of manufacturing
   NaNs or panicking.
7. Translate relevant upstream tests first, then add tests for intentional Go
   deviations and invalid inputs. Use explicit tolerances.
8. Keep each assigned step independently reviewable and leave `go test ./...`
   and `go vet ./...` passing.
9. Update `UPSTREAM.md` in the same step when a parity item is implemented,
   intentionally omitted, or found not to exist upstream.
10. Do not implement optional steps unless they are explicitly assigned.

## Decisions already settled

- `AxisAngle` and `Euler` are primarily generated upstream records. Do not add
  speculative method suites to them. Their substantive quaternion conversion
  behavior is already represented by `QuatFromAxisAngle`,
  `QuatFromEulerAngles`, `QuatFromYawPitchRoll`, and `Quat.EulerAngles`.
- Upstream `DQuaternion` and `DPlane` have no corresponding double-precision
  algorithm suites. Their current record and conversion APIs are sufficient.
- Expanding `DSphere` to mirror all `Sphere` methods would be a Go enhancement,
  not pinned-upstream parity. Leave it unchanged unless separately requested.
- Do not recreate `Transformable3D`, `ITransformable3D`, `IMappable`, or
  `IPoints` interfaces. Concrete Go methods and package functions remain the
  chosen design.
- Do not copy standard-library math wrappers, generated component-wise
  boilerplate, tuple/deconstruction APIs, collection overloads, throwing
  inverse wrappers, or known-defective upstream helpers.
- Binary/text/custom JSON formats remain deferred without a concrete consumer
  contract.

## Step 1: establish and test the exact remaining parity ledger

**Ownership:** documentation and parity classification only.

Audit the pinned versions of `src/MathOpsPartial.cs`, `src/Quaternion.cs`, and
`src/ITransformable.cs` against the current exported Go API. Update
`UPSTREAM.md` with a compact table that classifies each remaining public helper
as one of:

- required geometry behavior;
- optional data-layout compatibility;
- existing behavior under an idiomatic Go name;
- standard-library/generated duplicate;
- intentionally omitted C# convenience API.

At minimum, explicitly classify:

- vector projection, rejection, reflection, perpendicularity, collinearity,
  coplanarity, and back-face tests;
- vector Catmull–Rom, Hermite, and smooth-step interpolation;
- quaternion transformation overloads for 2D, 3D, and 4D vectors;
- normal transformation by a matrix;
- axis-angle vector rotation;
- quaternion look-at creation;
- screen-coordinate ray construction;
- nearest-power-of-two calculation;
- matrix flattening/unflattening and packed-box decoding;
- `Transformable3D` composition conveniences.

Do not change behavior in this step.

**Exit criteria:** the ledger names every remaining candidate and gives a clear
disposition; `go test ./...` and `go vet ./...` pass unchanged.

## Step 2: port vector geometric relations

**Ownership:** `vector.go`, `dvector.go` only if pinned upstream has a matching
double implementation, focused vector tests, and the corresponding ledger
rows.

Port the useful single-precision vector relations from
`src/MathOpsPartial.cs`:

- projection of one vector onto another;
- rejection from another vector;
- reflection around a normal;
- perpendicularity and collinearity tests;
- coplanarity and back-face tests where their argument roles and tolerance
  semantics can be documented clearly.

Choose names and receivers consistent with the current vector API. Validate
zero, non-finite, and non-unit preconditions rather than silently assuming
normalization. Preserve upstream behavior where valid, but use `(value, ok)`
where an undefined projection or equivalent operation can legitimately fail.

**Critical tests:** orthogonal projection/rejection recomposition, reflection
against unit and non-unit normals, zero-vector failures, tolerance boundaries,
collinear vectors with opposite directions, coplanar and non-coplanar point
sets, and non-mutation of inputs.

**Exit criteria:** all classified vector-relation kernels are implemented or
explicitly omitted with a reason; targeted tests, `go test ./...`, and
`go vet ./...` pass.

## Step 3: port vector interpolation

**Ownership:** vector implementation and focused interpolation tests. Avoid
editing quaternion or matrix code.

Add idiomatic vector forms of the existing scalar interpolation operations:

- Catmull–Rom;
- Hermite;
- smooth-step interpolation.

Support only vector widths and precisions present in the pinned substantive
source. Reuse the existing scalar formulas when that preserves source rounding
and avoids duplicated formulas. Do not clamp Catmull–Rom or Hermite amounts;
match the existing scalar `SmoothStep` clamping contract.

**Critical tests:** endpoints, extrapolation, component-wise agreement with the
scalar functions, smooth-step clamping, and representative upstream vectors.

**Exit criteria:** vector interpolation parity is documented and tested;
`go test ./...` and `go vet ./...` pass.

## Step 4: complete quaternion and vector rotation helpers

**Ownership:** `quaternion.go`, vector/quaternion tests, and matching ledger
rows. Do not add methods to the `Euler` or `AxisAngle` records merely for API
symmetry.

Audit and port the remaining substantive behavior from `src/Quaternion.cs` and
the quaternion-transform section of `src/MathOpsPartial.cs`:

- quaternion look-at creation;
- any missing vector transformation forms for `Vec2`, `Vec3`, and `Vec4` that
  have meaning distinct from the existing `Quat.Rotate(Vec3)`;
- axis-angle vector rotation when it is not merely a redundant wrapper;
- homogeneous-coordinate behavior for 4D values, preserving W exactly as
  upstream specifies.

First verify the handedness, local-forward axis, multiplication order, and
degenerate look-at behavior against the pinned source and existing
`LookAtMat4`. Use failure results for zero or parallel direction/up inputs when
no valid rotation exists.

**Critical tests:** identity, canonical-axis rotations, agreement between
quaternion and matrix rotation, look-at canonical directions, parallel and zero
input failures, `Vec2` embedding rules, `Vec4.W` behavior, and non-mutation.

**Exit criteria:** all substantive pinned quaternion helpers are represented or
documented as aliases/omissions; targeted tests, `go test ./...`, and
`go vet ./...` pass.

## Step 5: port matrix normal transformation

**Ownership:** `matrix.go`, matrix/vector tests, and one ledger row.

Add normal transformation only if it is not equivalent to the existing
direction transform. Normals under non-uniform scale or shear must use the
inverse transpose appropriate to this package's row-vector convention. Decide
from the pinned upstream contract whether the result is normalized; do not
silently add normalization.

Return failure when the required matrix inverse does not exist or inputs are
non-finite. Reuse `Mat4.Inverse` and existing vector operations rather than
duplicating inversion.

**Critical tests:** identity, translation independence, rotation, non-uniform
scale, shear, singular failure, geometric perpendicularity to transformed
tangents, and row-vector convention.

**Exit criteria:** normal transformation is correct for the package convention
and documented; targeted tests, `go test ./...`, and `go vet ./...` pass.

## Step 6: port screen-coordinate ray construction

**Ownership:** the module that currently owns ray/matrix interaction, focused
ray tests, one example if the API is not self-evident, and the ledger row.

Port the pinned projection-matrix helper that constructs a world-space ray from
screen coordinates. Before choosing the API, document its coordinate contract:

- screen origin and axis direction;
- pixel dimensions versus normalized coordinates;
- depth range;
- whether the supplied matrix is projection-only, view-projection, or its
  inverse;
- row-vector multiplication order.

Use `(Ray, bool)` for zero viewport dimensions, singular inversion, invalid
homogeneous division, or a degenerate direction. Do not call the ray parameter
a Euclidean distance unless the returned direction is normalized.

**Critical tests:** center-screen canonical ray, viewport corners, translated
and rotated camera/view transforms if supported by the pinned contract,
perspective divide, singular matrices, invalid dimensions, and projection
round trips.

**Exit criteria:** the helper has an unambiguous documented contract and
convention tests; `go test ./...` and `go vet ./...` pass.

## Step 7: port small scalar utility candidates

**Ownership:** `scalar.go`, focused scalar tests, and ledger rows.

Implement only the small pinned helpers judged generally useful in Step 1.
Expected candidate:

- nearest power-of-two calculation.

Specify behavior for zero, negative values, overflow, and exact powers of two.
Use `math/bits` where it produces the same defined result more clearly. Do not
port scalar wrappers already supplied by Go's `math` package.

**Exit criteria:** each candidate is implemented with edge-case tests or marked
omitted with a reason; `go test ./...` and `go vet ./...` pass.

## Optional step A: data-layout compatibility helpers

**Assign only when a consumer needs the pinned layout.**

Port matrix flattening/unflattening or packed-box decoding only after the caller
specifies the exact element order, scalar representation, bounds checks, and
ownership requirements. Prefer fixed arrays when the size is intrinsic. For
slices, reject incorrect lengths explicitly and never retain caller-owned
storage.

Do not treat raw Go struct memory as a wire format. If byte encoding is needed,
specify endianness and component order and update the intentional-deviations
section of `UPSTREAM.md`.

**Exit criteria:** golden layout tests generated from the pinned C# source pass
on the supported architectures; malformed input is tested; `go test ./...` and
`go vet ./...` pass.

## Optional step B: concrete transform composition conveniences

**Assign only after concrete consumer demand.**

Do not recreate the upstream `Transformable3D` interface or extension-method
namespace. Add only composition functions or concrete methods that materially
improve a demonstrated Go call site and cannot already be expressed clearly
with `TranslationMat4`, rotation factories, `ScaleMat4`, and `Mat4.Mul`.

Preserve the rule that in `A.Mul(B)`, A is applied first. Do not copy the
upstream `ScaleX`, `ScaleY`, and `ScaleZ` defect that zeroes unaffected axes.

**Exit criteria:** every added convenience removes real call-site complexity,
has an order-of-operations test, and is recorded as a deliberate Go API
addition rather than required parity.

## Optional step C: expand double-precision geometry

**Assign only as a new feature, not as parity work.**

If consumers require it, separately design double-precision quaternion, plane,
or sphere algorithms. Mirror the validated single-precision contracts and
failure behavior, but do not claim these algorithms were ported from equivalent
pinned upstream implementations.

Keep this work separate from the parity steps because it expands the API and
requires independent numerical review.

## Step 8: final integration

After all assigned required steps are merged:

1. Review the combined diff for duplicate APIs and inconsistent names.
2. Confirm every required row in the new parity table says implemented or gives
   a specific intentional-omission reason.
3. Run `gofmt` on changed Go files.
4. Run:

   ```sh
   go test ./...
   go test -race ./...
   go vet ./...
   ```

5. Run the existing benchmarks once to catch accidental catastrophic
   regressions; do not optimize merely because timings vary.
6. Ensure examples and package documentation still describe row-vector order,
   radians, handedness, failure results, and non-mutation accurately.

The completion condition is useful pinned-source behavioral parity, not a
one-for-one reproduction of every generated C# member.
