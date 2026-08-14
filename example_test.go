package math3d_test

import (
	"fmt"

	"github.com/maloquacious/math3d"
)

func ExampleVec3() {
	direction, ok := math3d.V3(3, 0, 4).Normalized()
	if !ok {
		return
	}

	fmt.Printf("unit: %.1f %.1f %.1f\n", direction.X, direction.Y, direction.Z)
	fmt.Printf("dot: %.1f\n", direction.Dot(math3d.V3(1, 2, 3)))
	// Output:
	// unit: 0.6 0.0 0.8
	// dot: 3.0
}

func ExampleTransform() {
	transform := math3d.NewTransform(
		math3d.V3(10, 0, 0),
		math3d.QuatRotationZ(math3d.HalfPi),
	)
	point := transform.TransformPoint(math3d.V3(1, 0, 0))

	fmt.Printf("%.0f %.0f %.0f\n", point.X, point.Y, point.Z)
	// Output:
	// 10 1 0
}

func ExampleMat4_Mul() {
	// With row vectors, multiplication order is application order: scale,
	// then rotate, then translate.
	matrix := math3d.ScaleMat4(math3d.V3(2, 2, 2)).
		Mul(math3d.RotationZMat4(math3d.HalfPi)).
		Mul(math3d.TranslationMat4(math3d.V3(5, -1, 0)))
	point := matrix.TransformPoint(math3d.V3(1, 0, 0))

	fmt.Printf("%.0f %.0f %.0f\n", point.X, point.Y, point.Z)
	// Output:
	// 5 1 0
}

func ExampleMat4_RayFromNormalizedScreen() {
	projection, ok := math3d.PerspectiveFOVMat4(math3d.HalfPi, 1, 1, 10)
	if !ok {
		return
	}
	// Screen coordinates are normalized with a top-left origin. A projection
	// matrix alone produces a view-space ray starting on the near clip plane.
	ray, ok := projection.RayFromNormalizedScreen(math3d.V2(0.5, 0.5))
	if !ok {
		return
	}

	fmt.Printf("origin: %.0f %.0f %.0f\n", ray.Origin.X, ray.Origin.Y, ray.Origin.Z)
	fmt.Printf("direction: %.0f %.0f %.0f\n", ray.Direction.X, ray.Direction.Y, ray.Direction.Z)
	// Output:
	// origin: 0 0 -1
	// direction: 0 0 -1
}

func ExampleBox3FromPoints() {
	box, ok := math3d.Box3FromPoints([]math3d.Vec3{
		math3d.V3(2, -1, 4),
		math3d.V3(-3, 5, 1),
		math3d.V3(0, 2, 8),
	})
	if !ok {
		return
	}

	fmt.Printf("min: %.0f %.0f %.0f\n", box.Min.X, box.Min.Y, box.Min.Z)
	fmt.Printf("max: %.0f %.0f %.0f\n", box.Max.X, box.Max.Y, box.Max.Z)
	// Output:
	// min: -3 -1 1
	// max: 2 5 8
}

func ExampleRay_IntersectSphere() {
	// Direction is not unit length, so t is a ray parameter rather than a
	// Euclidean distance.
	ray := math3d.NewRay(math3d.V3(-5, 0, 0), math3d.V3(2, 0, 0))
	t, ok := ray.IntersectSphere(math3d.NewSphere(math3d.V3(0, 0, 0), 1))
	if !ok {
		return
	}
	point := ray.PointAt(t)

	fmt.Printf("t: %.0f\n", t)
	fmt.Printf("point: %.0f %.0f %.0f\n", point.X, point.Y, point.Z)
	// Output:
	// t: 2
	// point: -1 0 0
}
