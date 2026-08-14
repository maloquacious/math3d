# Upstream parity ledger

This port is based on [`vimaec/math3d`](https://github.com/vimaec/math3d), pinned
to commit [`4477350aaf223ed6f8b2985419eb8da3e58f96ed`](https://github.com/vimaec/math3d/commit/4477350aaf223ed6f8b2985419eb8da3e58f96ed)
(`dev`, 2023-07-28). Upstream describes the repository as a mirror, so later
work must compare against this commit rather than the moving branch.

## Source inventory

The main library is `src/Vim.Math3D.csproj`. Its compiled sources at the pinned
revision are:

```text
src/AABox.cs                  src/MathOpsPartial.cs
src/AABox2D.cs                src/Matrix4x4.cs
src/Constants.cs              src/Plane.cs
src/ContainmentType.cs        src/PlaneIntersectionType.cs
src/Experimental.cs           src/Quad.cs
src/Hash.cs                   src/Quaternion.cs
src/IMappable.cs              src/Random.cs
src/IPoints.cs                src/Ray.cs
src/ITransformable.cs         src/Sphere.cs
src/LinqUtil.cs               src/Stats.cs
src/MathOps.cs                src/Structs.cs
src/StructsPartial.cs         src/Triangle.cs
src/Triangle2D.cs             src/ValueDomain.cs
```

The checked-in generated files `MathOps.cs` and `Structs.cs` originate from
`src/MathOps.tt`, `src/Structs.tt`, and `src/TemplateHelpers.tt`. Those templates
are part of the inventory but are not a separate public API surface.

## Public symbol map

Status meanings:

- **planned**: part of the stable Go port, in a later structural or algorithm
  milestone;
- **deferred**: public functionality to evaluate after the geometry it needs;
- **omitted**: deliberately not represented in the stable Go API.

### Value types and enums

| Upstream declaration | Go declaration | Status |
|---|---|---|
| `Vector2`, `Vector3`, `Vector4` | `Vec2`, `Vec3`, `Vec4` | implemented (milestone 3) |
| `DVector2`, `DVector3`, `DVector4` | `DVec2`, `DVec3`, `DVec4` | implemented (milestone 3) |
| `Int2`, `Int3`, `Int4` | `IVec2`, `IVec3`, `IVec4` | implemented (milestone 2) |
| `Byte2`, `Byte3`, `Byte4` | unchanged | implemented (milestone 2) |
| `Complex` | `Complex` | implemented (milestone 2); two `float64` fields |
| `Quaternion`, `DQuaternion` | `Quat`, `DQuat` | quaternion algorithms implemented (milestone 5); upstream `DQuaternion` remains a record only |
| `Matrix4x4` | `Mat4` | operations and factories implemented (milestone 5) |
| `Transform` | `Transform` | rigid-transform algorithms implemented (milestone 6) |
| `Plane`, `DPlane` | unchanged | `Plane` algorithms implemented (milestone 6); upstream `DPlane` remains a record only |
| `Interval`, `DInterval` | unchanged | algorithms implemented (milestone 7) |
| `AABox2D`, `AABox`, `AABox4D` | `Box2`, `Box3`, `Box4` | algorithms implemented (milestone 7) |
| `DAABox2D`, `DAABox`, `DAABox4D` | `DBox2`, `DBox3`, `DBox4` | generated-equivalent algorithms implemented (milestone 7) |
| `Ray`, `DRay` | unchanged | box/sphere/plane and transform algorithms implemented (milestone 7); `Ray` triangle intersection implemented (milestone 8) |
| `Sphere`, `DSphere` | unchanged | milestone 7 containment, intersection, merge, bounds, and transform behavior implemented for `Sphere`; `DSphere` supports validity, point containment, bounds, sphere overlap, and ray intersection |
| `Line`, `Line2D` | `Segment3`, `Segment2` | mapping, geometry, and 2D intersection implemented (milestone 8) |
| `Triangle`, `Triangle2D` | `Triangle3`, `Triangle2` | area, normals, bounds, containment, transforms, and ray intersection implemented (milestone 8) |
| `Quad`, `Quad2D` | `Quad3`, `Quad2` | `Quad3` mapping and transforms implemented (milestone 8); `Quad2` remains a record like upstream |
| `ColorRGB`, `ColorRGBA`, `ColorHDR` | `RGB`, `RGBA`, `HDRColor` | implemented (milestone 2) |
| `SphericalCoordinate` | `Spherical` | implemented (milestone 2) |
| `PolarCoordinate` | `Polar` | implemented (milestone 2) |
| `LogPolarCoordinate` | `LogPolar` | implemented (milestone 2) |
| `CylindricalCoordinate` | `Cylindrical` | implemented (milestone 2) |
| `HorizontalCoordinate` | `Horizontal` | implemented (milestone 2) |
| `GeoCoordinate` | `Geo` | implemented (milestone 2) |
| `AxisAngle`, `Euler` | unchanged | records implemented (milestone 4); substantive conversions are represented by the quaternion API |
| `Stats<T>` | `Stats[T]` | implemented (milestone 9), with vector aggregation helpers |
| `ValueDomain` | `Domain` | implemented (milestone 9) |
| `ContainmentType` | `Containment` | implemented (milestone 2), with explicit zero value |
| `PlaneIntersectionType` | `PlaneIntersection` | implemented (milestone 2), with explicit zero value |

All entries above are actual declarations at the pinned revision. `AxisAngle`
uses `double` and `DVector3` upstream despite lacking a `D` prefix; the Go type
will therefore use `float64`. There is no double-precision matrix declaration.

### Static helpers and interfaces

| Upstream declaration | Go treatment | Status |
|---|---|---|
| `Constants` | typed constants and value-returning functions | implemented (milestone 2) |
| `MathOps` | focused methods and package functions | milestone 9 scalar interpolation, angle conversion, validation, conversion, statistics, and swizzle helpers implemented; remaining candidates classified below |
| `Hash` | `CombineHash` and `HashInts` | implemented (milestone 9) |
| `StatelessRandom` | deterministic scalar and vector package functions | implemented (milestone 9) |
| `LinqUtil` | ordinary slices and standard-library functions | omitted; no Go-specific helper remained useful |
| `Transformable3D` | concrete operations, not a namespace type | generic conveniences omitted; optional concrete conveniences require consumer demand |
| `ITransformable3D<TSelf>` | concrete methods | omitted |
| `IMappable<TContainer,TPart>` | add a small interface only if multiple Go implementations require it | omitted initially |
| `IPoints`, `IPoints2D` | slices at algorithm boundaries | omitted |
| `MovementExtensions` | no active members at the pinned revision | omitted |

The constants to port include `Pi`, `HalfPi`, `TwoPi`, `Tolerance`, `Log10E`,
`Log2E`, `E`, `RadiansToDegrees`, `DegreesToRadians`, `OneTenthOfADegree`, and
the millimetre/foot conversion factors. The upstream names/comments for the
last two are misleading; the Go port will use numerically correct directional
names and record that compatibility deviation below. Mutable plane variables
(`XYPlane`, `XZPlane`, and `YZPlane`) will become value-returning functions.

## Remaining helper parity audit

The tables below audit every public helper in the pinned
`MathOpsPartial.cs`, `Quaternion.cs`, and `ITransformable.cs` that was not
already a direct name-for-name part of the initial type inventory. “Required”
means useful geometry behavior included in the completed parity pass. Optional
layout and composition rows remain unassigned unless a consumer requests them.

### Required geometry behavior

| Pinned upstream helper | Go disposition |
|---|---|
| `Vector3.Projection`, `Rejection` | Implemented as fallible `Vec3.Projection` and `Vec3.Rejection`; zero, non-finite, overflowed, and non-finite-result cases fail. |
| `Reflect(Vector2, Vector2)`, `Reflect(Vector3, Vector3)` | Implemented as fallible `Vec2.Reflect` and `Vec3.Reflect`. Unlike upstream, non-unit normals are supported by dividing by their squared magnitude. |
| `Vector3.IsPerpendicular`, `Colinear` | Implemented as `Vec3.IsPerpendicular` and `Vec3.IsCollinear`. Perpendicularity preserves the pinned strict absolute dot tolerance; collinearity uses an unsigned 3D angular tolerance and accepts opposite directions, fixing the defective one-sided signed-angle implementation. |
| `Coplanar(Vector3, Vector3, Vector3, Vector3)` | Implemented as `Coplanar`, preserving the pinned strict absolute scalar-triple-product tolerance and documenting its scale dependence. |
| `Vector3.IsBackFace` | Implemented as `Vec3.IsBackFace`; the receiver is the surface normal and the argument is the line-of-sight direction. Zero and invalid inputs return false. |
| `Vector3.Angle`, `SignedAngle` overloads | Implemented as fallible `Vec3.Angle` and `Vec3.SignedAngle`; angles are radians, and the sign is determined by `axis·(from×to)`. Ambiguous nonzero signed angles fail. |
| `Vector3.CatmullRom`, `Hermite`, `SmoothStep` | Implemented as component-wise `Vec3.CatmullRom`, `Vec3.Hermite`, and `Vec3.SmoothStep`, reusing the scalar kernels and their clamping contracts. No pinned `Vector2`, `Vector4`, or double-vector overload exists. |
| Quaternion transforms of `Vector2`, `Vector3`, and `Vector4`, including `TransformToVector4` | Implemented as `Quat.RotateVec2`, `Quat.RotateVec4`, `Quat.RotateVec2ToVec4`, and `Quat.RotateToVec4`, alongside the existing `Quat.Rotate(Vec3)`. The 2D form embeds with Z = 0 and discards output Z, the 4D form preserves input W exactly, and the to-Vec4 forms write W = 1. |
| `Quaternion.LookAt` | Implemented as fallible `LookAtQuat`, preserving the pinned caller-supplied local-forward axis, up-plane heading followed by tilt, and `q2 * q1` multiplication order. Unlike upstream, it normalizes finite nonzero `up` and `localForward` and rejects coincident points, direction/up parallelism, and local-forward/up non-perpendicularity instead of propagating NaNs or constructing an ambiguous rotation. |
| `TransformNormal(Vector2/Vector3/Vector4, Matrix4x4)` | Implemented for the package's primary 3D geometry as fallible `Mat4.TransformNormal`. Unlike the pinned helpers, which merely apply the linear matrix, Go uses the row-vector inverse transpose so normals remain perpendicular under non-uniform scale and shear. The result is not normalized, matching the pinned magnitude contract. |
| `Matrix4x4.RayFromProjectionMatrix` | Implemented as `Mat4.RayFromNormalizedScreen`. It consumes normalized top-left-origin coordinates, uses clip depths 0 and 1, and inverts the supplied projection or row-vector view-projection matrix. The result starts on the near clip plane and has a unit direction; invalid inversion, homogeneous division, and direction report failure. |
| `ToNearestPowOf2` | Implemented as `NearestPowerOfTwo(int32) (int32, bool)`. It preserves logarithmic nearest-power rounding for positive inputs and reports zero, negative inputs, and `int32` result overflow as failures. |

The pinned source has no double-precision counterparts for projection,
rejection, reflection, perpendicularity, collinearity, angles, coplanarity, or
back-face tests, so Step 2 intentionally adds no `DVec` relation methods.

### Optional data-layout compatibility

| Pinned upstream helper | Go disposition |
|---|---|
| `Matrix4x4.ToFloats`, `Matrix4x4[].ToFloats` | Optional row-major component flattening; add only with a concrete slice/array ownership contract. |
| `float[].ToMatrix`, `ToMatrixArray` | Optional inverse of the pinned 16-float row-major layout; malformed-length behavior must be defined rather than relying on indexing or `Debug.Assert`. |
| `float[].ToAABoxArray` | Optional packed-box decoding in `Min.X, Min.Y, Min.Z, Max.X, Max.Y, Max.Z` order; add only for a consumer of that format. |
| `Quaternion.Vector4` | Optional four-component record conversion. Public quaternion fields and `NewQuat` already expose the component layout, so this is not an algorithmic parity gap. |

### Existing behavior under idiomatic Go names

| Pinned upstream helper | Existing Go behavior |
|---|---|
| `Percentage`; scalar `CatmullRom`, `Hermite`, `SmoothStep`, `WrapAngle`, `IsNonZeroAndValid` | `Percentage`, `CatmullRom`, `Hermite`, `SmoothStep`, `WrapAngle`, and `IsNonZeroAndValid`. |
| Scalar/vector splat, narrowing, widening, and component conversions | `SplatV2`/`SplatV3`/`SplatV4`, `Vec2.Vec3`, `Vec2.Vec4`, `Vec3.Vec4`, `Vec4.Vec2`, and `Vec4.Vec3`. |
| `Vector3.Transform(Quaternion)` | `Quat.Rotate(Vec3)`. |
| `Vector3.Rotate(axis, angle)` | `Mat4FromAxisAngle(axis, angle).TransformDirection(v)` or `QuatFromAxisAngle(axis, angle).Rotate(v)`; the upstream method is only this composition. |
| `Vector2.Transform(Matrix4x4)` and matrix `TransformToVector4` overloads | `Mat4.TransformPoint`/`TransformVec4` with explicit `Z = 0` and `W = 1` embedding. Keeping embedding visible avoids overload-dependent point semantics. |
| `Stats<Vector3>.ToBox`, `Stats<DVector3>.ToBox`, `IEnumerable<Vector3>.ToBox` | `NewBox3`, `NewDBox3`, and `Box3FromPoints`. |
| `Vector3.ToMatrix`, `Quaternion.ToMatrix`, `Transform.ToMatrix` | `TranslationMat4`, `Mat4FromQuat`, and `Transform.Mat4`. |
| Quaternion identity, vector/scalar construction, length, normalization, conjugate, inverse, axis/Euler/axis-specific/yaw-pitch-roll creation | `IdentityQuat`, `QuatFromVector`, quaternion methods, `QuatFromAxisAngle`, `QuatFromEulerAngles`, `QuatRotationX/Y/Z`, and `QuatFromYawPitchRoll`. |
| `Quaternion.CreateRotationFromAToB`, `CreateFromRotationMatrix` | `RotationBetween` and `Mat4.Quat`. |
| Quaternion dot, slerp, lerp, concatenate, arithmetic operators, and `ToEulerAngles` | `Quat.Dot`, `Quat.Slerp`, `Quat.Lerp`, `Concatenate`, arithmetic methods, and `Quat.EulerAngles`; fallible operations use Go's `(value, ok)` convention. |
| `Transformable3D.Multiply(params Matrix4x4[])` | Repeated `Mat4.Mul`, with the package's documented left-to-right application order. A variadic fold adds no geometry behavior. |

### Standard-library or generated duplicates

| Pinned upstream helper | Go disposition |
|---|---|
| `Vector4.Transform(Matrix4x4)` forwarding overload | Generated/self-forwarding duplicate of matrix-vector transformation; use `Mat4.TransformVec4`. |
| Static `Cross(Vector3, Vector3)` and `Cross(DVector3, DVector3)` | Generated-style duplicates of `Vec3.Cross` and `DVec3.Cross`. |
| Quaternion `IsIdentity` and equality/operator-shaped members | Direct comparison with `IdentityQuat()` and ordinary Go methods/operators already supply the useful behavior. |

### Intentionally omitted C# convenience API

| Pinned upstream helper | Reason |
|---|---|
| Throwing `Matrix4x4.Inverse()` extension | `Mat4.Inverse() (Mat4, bool)` represents singularity without exceptions. |
| `Vector3.IsNonZeroAndValid`, `IsZeroOrInvalid` | Thin combinations of `IsFinite` and squared magnitude; no second validity vocabulary is needed. |
| `Vector3.ToLine`, `Along`, and scalar `AlongX/Y/Z` | Constructor/normalize/scale conveniences are clearer explicitly. Pinned `AlongZ` also incorrectly returns the X axis. |
| `Quaternion.ToSphericalAngle` overloads, `Create(HorizontalCoordinate)`, and implicit quaternion/horizontal conversions | Specialized coordinate conveniences with a hard-coded default forward axis and silent normalization singularities. Callers can state their chosen forward convention explicitly with `Quat.Rotate` and ordinary trigonometry; implicit conversions are not idiomatic Go. |
| `ITransformable3D<TSelf>` and generic `Transform` | The C# interface and extension dispatch are intentionally not reproduced; concrete Go values own their transformation methods. |
| Generic `Translate`, `Rotate`, `Scale`, `ScaleX/Y/Z`, `LookAt`, `RotateAround`, `Reflect`, axis rotations, and `TranslateRotateScale` | Optional `Transformable3D` composition conveniences, not geometry kernels. Existing matrix factories and `Mat4.Mul` express them without a broad interface. Pinned `ScaleX/Y/Z` also zero the unaffected axes and are not copied. |

### Experimental and placeholder declarations

Everything in the `Vim.Experimental` namespace is omitted from the stable port:

```text
IAlgebraicSet<T>  IPolygon              SurfaceDistance
IInsideOutside    ILerp<T>              IPositionRotationScale
IPoints<T>        IRange<T>             IField<T>
IVectorField      ISignedDistanceField  ICurve
ISurface          IClosedShape          IVolume
IPrimitiveShape   IPrimitiveCurve       IPrimitiveVolume
IExtrudedCurve    ILoftedCurve          PrimitiveShapes
```

`SurfaceDistance` and `IInsideOutside` are empty interfaces. `PrimitiveShapes`
contains uninitialized placeholder fields rather than implementations. They do
not establish behavior to port.

The upstream README also names `LinearMotion`, `AngularMotion`, `Motion`,
`Triangle2`, and `DQuad`, but none is declared in the current source. They are
stale documentation and are omitted. The Go name `Triangle2` maps the actual
upstream `Triangle2D`; it does not port a separate README-only type.

## Intentional deviations

This ledger records decided deviations now; later compatibility fixes must be
added when implemented.

1. Go values expose public fields and use value receivers rather than generated
   C# properties, operators, and broad interfaces.
2. Fallible geometric operations use `(value, ok)` instead of sentinels or
   exceptions for ordinary misses and invalid mathematical preconditions.
3. Package operations do not mutate mathematical inputs. Callers remain free to
   modify public fields on values they own.
4. Matrix behavior preserves upstream row-vector semantics, composition order,
   and translation in `M41`, `M42`, and `M43`.
5. The source plane-sign inconsistency will be standardized on
   `dot(normal, point) + d = 0` when plane algorithms are implemented.
6. Ray/sphere intersection will support non-unit ray directions rather than
   preserving the source's implicit normalization assumption.
7. Millimetre/foot conversion APIs will be named for their numerical direction,
   correcting misleading upstream names/comments.
8. C# collection, LINQ, transformation, and mapping interfaces are not copied;
   concrete methods and slices will be used unless a real Go interface boundary
   emerges.
9. Raw struct memory is never a wire format. Any encoding will specify component
   order and endianness explicitly.
10. Vector normalization returns `ok == false` for zero, non-finite, and
    overflowed magnitudes instead of returning NaN or a misleading zero vector.
11. `Ray.Position` and `DRay.Position` are named `Origin` in Go to distinguish
    the ray's fixed starting point from points evaluated along its direction.
12. Binary encoding remains deferred until a concrete interchange requirement
    exists; milestone 4 defines component order but does not make raw Go struct
    memory a wire format.
13. Quaternion normalization, inversion, division, and linear interpolation use
    `(value, ok)` when a zero or non-finite norm prevents a meaningful result.
    Upstream returns NaN components for those invalid preconditions.
14. Matrix inversion returns `(Mat4{}, false)` for a singular or non-finite
    matrix instead of upstream's all-NaN failure sentinel.
15. Plane construction from a normal and point uses `D = -dot(normal, point)`,
    and point projection uses `dot(normal, point) + D`, consistently following
    the package plane equation. The upstream versions use the opposite sign.
16. Plane normalization, construction, projection, and matrix transformation
    report invalid mathematical preconditions with an `ok` result. Projection
    also supports non-unit normals instead of assuming a normalized plane.
17. Intervals and boxes expose `Valid`/`Empty`; intersections return an
    invalid representation together with `ok == false`. Canonical empty values
    use reversed infinities so they remain merge identities.
18. Sphere containment explicitly requires the containing radius to be at
    least the contained radius, fixing the upstream reversed-containment bug.
    Sphere/box overlap uses nearest-point clamping, fixing upstream false
    negatives for spheres whose centers are inside a box near multiple faces.
19. Sphere transformation uses the largest singular value of the linear matrix
    rather than the largest row length. This preserves rotation and non-uniform
    scale behavior while also providing a conservative enclosure under shear.
20. Triangle normals and segment length changes report failure for degenerate
    or non-finite inputs rather than allowing normalization to produce NaNs.
    `Triangle3.Degenerate` tests geometric zero area rather than upstream's
    misleading pairwise vertex inequality. `Triangle2.SignedArea` uses the
    conventional positive-for-counter-clockwise sign instead of upstream's
    reversed sign.
21. Vector statistics initialize component minima and maxima from the first
    input rather than zero, fixing biased bounds for all-positive and
    all-negative sequences. Empty averages return `ok == false` instead of NaN.
22. The upstream `Vector3.ZYX` implementation returns `(Z,Y,Z)`; Go returns the
    correctly named `(Z,Y,X)` swizzle.
23. The defective list-specific hash overload that combines list indexes rather
    than values is omitted. `HashInts` consistently folds the supplied values.
24. Binary, text, and custom JSON encodings remain deferred at milestone 10:
    the repository has no concrete VIM interchange or consumer format
    requirement. Public fields continue to use Go's standard JSON behavior.
25. `Mat4.TransformNormal` uses the mathematically correct row-vector inverse
    transpose and reports inversion or finite-value failures. Pinned upstream's
    `TransformNormal` overloads only apply the ordinary linear transform. Go
    preserves their no-normalization behavior but not that defect under
    non-uniform scale or shear.
26. `NearestPowerOfTwo` reports failure for zero, negative inputs, and a nearest
    power outside the positive `int32` range. Pinned `ToNearestPowOf2` passes
    those cases through logarithm, floating-point conversion, and unchecked
    integer conversion without a useful mathematical result.

## Release readiness

Milestone 10 adds executable package examples for vectors, rigid transforms,
row-vector matrix composition, bounds construction, and ray intersections. It
also adds repeatable benchmarks for vector cross products, matrix
multiplication and inversion, and ray/box intersection. These benchmarks are
baselines only; no readable implementation was replaced without evidence.

The required remaining-parity rows are complete. Optional data-layout,
transform-composition, and expanded double-precision work remains deferred
until requested by a concrete consumer.

## Source behavior notes

- Cartesian cross products are right-handed (`UnitX × UnitY = UnitZ`), and view
  helpers use negative Z as local forward. Upstream does not mandate a world-up
  axis.
- Matrices are row-major by field order and transform row vectors on the left.
  In `A * B`, A is applied first.
- Rotation APIs use radians and positive rotations follow the right-hand rule.
  Upstream coordinate record declarations do not consistently state angle
  units, so each Go record must document its own evidence-backed convention.
- Primary values use 32-bit floats; `D` values use 64-bit floats. Upstream's
  README descriptions of `DQuaternion` and `AxisAngle` are incorrect; their
  source fields are double precision.
- Upstream approximate equality is component-wise strict absolute difference,
  `abs(a-b) < tolerance`, with default tolerance `1e-7`. It has no relative,
  ULP, NaN, infinity, or geometric-equivalence special cases.
