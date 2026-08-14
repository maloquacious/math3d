# Transform points, directions, and rotations

Use these patterns to choose and compose transformations without mixing the
package's matrix and quaternion ordering rules.

## Choose `Transform` or `Mat4`

Use `Transform` for rigid placement: one rotation followed by one translation.
It cannot represent scale, shear, or projection. Create identity explicitly
with `IdentityTransform`; the zero value is not identity.

```go
orientation := math3d.QuatRotationY(math3d.HalfPi)
placement := math3d.NewTransform(math3d.V3(10, 0, 0), orientation)
worldPoint := placement.TransformPoint(localPoint)
```

`NewTransform` preserves the supplied quaternion. It does not normalize it.
The orientation must be unit length when `TransformPoint`,
`TransformDirection`, or `Mat4` uses it as a rotation. Normalize a quaternion
and handle failure when it did not come from a unit-producing rotation helper:

```go
orientation, ok := inputOrientation.Normalized()
if !ok {
	return errors.New("orientation is zero or non-finite")
}
placement := math3d.NewTransform(position, orientation)
```

Use `Mat4` for scale, projection, shear, general affine composition, or APIs
that require a matrix. `IdentityMat4` is matrix identity; a zero `Mat4` is not.

## Build scale, rotation, then translation

`ComposeMat4` creates the common SRT sequence directly:

```go
matrix := math3d.ComposeMat4(
	math3d.V3(2, 3, 4),
	math3d.QuatRotationZ(math3d.HalfPi),
	math3d.V3(5, -1, 0),
)
```

The quaternion must be unit length for the linear part to contain only the
requested scale and rotation.

To build the same sequence from individual matrices, chain `Mul` in application
order. Matrices use row vectors, so `A.Mul(B)` applies A and then B:

```go
matrix := math3d.ScaleMat4(math3d.V3(2, 2, 2)).
	Mul(math3d.RotationZMat4(math3d.HalfPi)).
	Mul(math3d.TranslationMat4(math3d.V3(5, -1, 0)))
point := matrix.TransformPoint(math3d.V3(1, 0, 0))
```

This produces `(5, 1, 0)`. Reversing the chain applies translation, rotation,
then scale and produces `(2, 12, 0)`, so the numerical result exposes an
ordering reversal. The first sequence is maintained as executable
[`ExampleMat4_Mul`](../example_test.go).

## Sequence quaternion rotations

Quaternion Hamilton multiplication has the opposite reading from matrix
chaining: `q.Mul(other)` applies `other` first and `q` second. Prefer
`Concatenate` when the code should state application order:

```go
first := math3d.QuatRotationZ(math3d.HalfPi)
second := math3d.QuatRotationX(math3d.HalfPi)
rotation := math3d.Concatenate(first, second)
result := rotation.Rotate(math3d.V3(1, 0, 0)) // (0, 0, 1)
```

Both quaternions must be unit length when used as rotations. See the executable
[`ExampleConcatenate`](../example_test.go).

## Transform the right geometric quantity

Use the operation that matches the input's meaning:

- `TransformPoint` uses homogeneous W=1, so translation affects positions.
- `TransformDirection` uses homogeneous W=0, so translation has no effect.
- `TransformNormal` applies the inverse transpose of the matrix so a surface
  normal remains perpendicular under non-uniform scale or shear. It returns
  `(normal, ok)` because this requires an invertible finite matrix.

`TransformNormal` does not normalize its result. Normalize it separately when
the consumer requires a unit normal, and handle both operations:

```go
normal, ok := matrix.TransformNormal(localNormal)
if !ok {
	return errors.New("matrix cannot transform normals")
}
unitNormal, ok := normal.Normalized()
if !ok {
	return errors.New("transformed normal has no finite direction")
}
```

The executable [`ExampleMat4_TransformNormal`](../example_test.go) applies the
same non-uniform scale to a point, direction, and normal to make their different
results visible.

## Handle inversion and decomposition failures

Matrix inversion can fail for a singular or non-finite matrix:

```go
inverse, ok := matrix.Inverse()
if !ok {
	return errors.New("matrix is singular or non-finite")
}
```

Decomposition expects an affine scale-rotation-translation matrix. It can fail
when the linear part is non-finite or cannot be represented as scale followed
by an orthonormal rotation within the method's tolerance, including matrices
with shear:

```go
scale, rotation, translation, ok := matrix.Decompose()
if !ok {
	return errors.New("matrix is not decomposable as SRT")
}
```

The returned rotation is a unit quaternion when decomposition succeeds.
