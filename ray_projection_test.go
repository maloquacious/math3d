package math3d

import (
	"math"
	"testing"
)

func TestRayFromNormalizedScreenPerspectiveConventions(t *testing.T) {
	projection, ok := PerspectiveFOVMat4(HalfPi, 1, 1, 10)
	if !ok {
		t.Fatal("valid perspective rejected")
	}
	tests := []struct {
		name      string
		screen    Vec2
		origin    Vec3
		direction Vec3
	}{
		{"center", V2(0.5, 0.5), V3(0, 0, -1), V3(0, 0, -1)},
		{"top-left", V2(0, 0), V3(-1, 1, -1), V3(-1, 1, -1).Scale(1 / float32(math.Sqrt(3)))},
		{"bottom-right", V2(1, 1), V3(1, -1, -1), V3(1, -1, -1).Scale(1 / float32(math.Sqrt(3)))},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ray, ok := projection.RayFromNormalizedScreen(test.screen)
			if !ok || !ray.Origin.AlmostEqual(test.origin, algorithmTolerance) ||
				!ray.Direction.AlmostEqual(test.direction, algorithmTolerance) {
				t.Fatalf("RayFromNormalizedScreen(%v) = %#v, %v", test.screen, ray, ok)
			}
			if abs32(ray.Direction.Dot(ray.Direction)-1) > algorithmTolerance {
				t.Fatalf("direction is not unit length: %#v", ray.Direction)
			}
		})
	}
}

func TestRayFromNormalizedScreenOffCenterPerspectiveDivide(t *testing.T) {
	projection, ok := PerspectiveOffCenterMat4(-2, 1, -1, 3, 2, 20)
	if !ok {
		t.Fatal("valid off-center perspective rejected")
	}
	tests := []struct {
		screen Vec2
		near   Vec3
	}{{
		V2(0, 0), V3(-2, 3, -2),
	}, {
		V2(0.5, 0.5), V3(-0.5, 1, -2),
	}, {
		V2(1, 1), V3(1, -1, -2),
	}}
	for _, test := range tests {
		ray, ok := projection.RayFromNormalizedScreen(test.screen)
		wantDirection, _ := test.near.Normalized()
		if !ok || !ray.Origin.AlmostEqual(test.near, 2*algorithmTolerance) ||
			!ray.Direction.AlmostEqual(wantDirection, 2*algorithmTolerance) {
			t.Fatalf("RayFromNormalizedScreen(%v) = %#v, %v", test.screen, ray, ok)
		}
	}
}

func TestRayFromNormalizedScreenViewProjectionAndRoundTrip(t *testing.T) {
	projection, ok := PerspectiveFOVMat4(HalfPi, 1, 1, 10)
	if !ok {
		t.Fatal("valid perspective rejected")
	}
	tests := []struct {
		name      string
		position  Vec3
		target    Vec3
		origin    Vec3
		direction Vec3
	}{
		{"translated", V3(10, 0, 0), V3(10, 0, -1), V3(10, 0, -1), V3(0, 0, -1)},
		{"rotated", Vec3{}, V3(1, 0, 0), V3(1, 0, 0), V3(1, 0, 0)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			view, ok := LookAtMat4(test.position, test.target, V3(0, 1, 0))
			if !ok {
				t.Fatal("valid view rejected")
			}
			viewProjection := view.Mul(projection)
			ray, ok := viewProjection.RayFromNormalizedScreen(V2(0.5, 0.5))
			if !ok || !ray.Origin.AlmostEqual(test.origin, 2*algorithmTolerance) ||
				!ray.Direction.AlmostEqual(test.direction, 2*algorithmTolerance) {
				t.Fatalf("world ray = %#v, %v", ray, ok)
			}

			worldPoint := ray.PointAt(3)
			clip := viewProjection.TransformVec4(V4(worldPoint.X, worldPoint.Y, worldPoint.Z, 1))
			if clip.W == 0 {
				t.Fatal("round-trip point has zero homogeneous W")
			}
			screen := V2((clip.X/clip.W+1)/2, (1-clip.Y/clip.W)/2)
			if !screen.AlmostEqual(V2(0.5, 0.5), 2*algorithmTolerance) {
				t.Fatalf("projected screen coordinate = %#v", screen)
			}
		})
	}
}

func TestRayFromNormalizedScreenFailuresAndNonMutation(t *testing.T) {
	screen := V2(0.25, 0.75)
	projection, ok := PerspectiveFOVMat4(HalfPi, 1, 1, 10)
	if !ok {
		t.Fatal("valid perspective rejected")
	}
	originalProjection, originalScreen := projection, screen
	if _, ok := projection.RayFromNormalizedScreen(screen); !ok {
		t.Fatal("valid screen coordinate rejected")
	}
	if projection != originalProjection || screen != originalScreen {
		t.Fatal("RayFromNormalizedScreen mutated an input")
	}

	if _, ok := (Mat4{}).RayFromNormalizedScreen(screen); ok {
		t.Fatal("singular matrix produced a ray")
	}
	if _, ok := projection.RayFromNormalizedScreen(V2(float32(math.NaN()), 0)); ok {
		t.Fatal("non-finite screen coordinate produced a ray")
	}
	// This invertible matrix swaps homogeneous Z and W. At clip depth zero,
	// unprojection therefore has W=0 and cannot be divided into a point.
	zeroNearW := NewMat4(
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 0, 1,
		0, 0, 1, 0,
	)
	if _, ok := zeroNearW.RayFromNormalizedScreen(screen); ok {
		t.Fatal("zero homogeneous W produced a ray")
	}
}
