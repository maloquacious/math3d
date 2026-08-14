# Work with bounds and shapes

Use these focused workflows to build bounds, combine them, transform them, and
interpret segment, triangle, and plane representations.

## Build bounds from a point slice

Use `Box3FromPoints` when all points are already available:

```go
box, ok := math3d.Box3FromPoints(points)
if !ok {
	return errors.New("points are empty or contain a non-finite value")
}
```

The result contains every point with inclusive minimum and maximum endpoints.
An empty slice or any non-finite point returns false. See executable
[`ExampleBox3FromPoints`](../example_test.go).

## Accumulate bounds incrementally

Start an incremental accumulation with `EmptyBox3`, not `Box3{}`:

```go
box := math3d.EmptyBox3()
for _, point := range points {
	box = box.Expanded(point)
}
```

`EmptyBox3` uses reversed infinite endpoints and is the canonical empty
accumulation identity. The zero `Box3` is instead a valid point box at the
origin, so starting with it would always include the origin. Executable
[`ExampleEmptyBox3`](../example_test.go) demonstrates accumulation.

## Validate and combine boxes

`NewBox3` preserves endpoint order. A box is valid when every `Min` component
is at most its corresponding `Max` component. Any reversed component makes the
box empty according to `Empty`; a flat box or point box remains valid.

Containment and intersection boundaries are inclusive:

- `ContainsPoint` includes faces, edges, and corners;
- `Intersects` reports true when two valid boxes merely touch;
- `Intersection` returns true for a touching point, edge, or face;
- `Intersection` returns false for invalid or disjoint boxes.

Always check the result when an overlap can be absent:

```go
overlap, ok := first.Intersection(second)
if !ok {
	return errors.New("boxes are invalid or disjoint")
}
```

See executable [`ExampleBox3_Intersection`](../example_test.go). `Merge` returns
the componentwise bounds enclosing both boxes; `EmptyBox3` is its identity.

## Transform axis-aligned bounds

```go
worldBounds, ok := localBounds.Transformed(matrix)
if !ok {
	return errors.New("bounds or transformed corners are invalid")
}
```

The result is still axis-aligned. It is the axis-aligned box around all eight
transformed corners, not an oriented box that retains the original edge axes.
Rotation or shear can therefore make it larger than the transformed shape.

## Evaluate a segment parameter

`Segment3.PointAt(t)` evaluates the segment's supporting line:

```go
midpoint := segment.PointAt(0.5)
beyondB := segment.PointAt(1.5)
```

Parameters in `[0, 1]` lie on the closed segment: 0 is A and 1 is B. Values
outside that range extrapolate beyond the endpoints. `Segment3.Ray` uses the
unnormalized displacement from A to B as its direction, so its ray parameter is
not generally a Euclidean distance.

## Preserve and inspect triangle winding

Triangle constructors retain vertex order. For `Triangle2`, `SignedArea` is
positive for counter-clockwise winding, negative for clockwise winding, and
zero for a degenerate triangle. `Area` returns the non-negative magnitude.

`Triangle2.Contains` tests the strict interior. Points on an edge or vertex are
excluded, even though other package operations may use inclusive boundaries.

For `Triangle3`, `NormalDirection` returns the unnormalized, winding-dependent
cross product. `Normal` normalizes it and returns false for non-finite or
zero-area triangles:

```go
normal, ok := triangle.Normal()
if !ok {
	return errors.New("triangle has no finite surface normal")
}
```

Reversing the winding reverses the normal direction.

## Classify and project with planes

Planes use this stored equation:

```text
dot(Normal, point) + D = 0
```

`NewPlane` preserves `Normal` and `D`; it does not normalize them. A zero normal
does not define a geometric plane.

`ClassifyPoint` returns the equation value. Negative values are behind the
normal, positive values are in front, and zero lies on the plane. The value is
a signed distance only when `Normal` has unit length.

`Project` handles non-unit finite normals correctly by dividing by the normal's
squared length. It returns false for a zero or non-finite normal or a non-finite
result:

```go
projected, ok := plane.Project(point)
if !ok {
	return errors.New("plane cannot project the point")
}
```

Executable [`ExamplePlane_Project`](../example_test.go) uses the non-unit plane
`2*x - 4 = 0`. Its classification value is 6 for a point three units from the
plane, while projection still returns the correct point on `x = 2`.
