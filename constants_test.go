package math3d

import "testing"

func TestDerivedConstants(t *testing.T) {
	if HalfPi != Pi/2 || TwoPi != Pi*2 {
		t.Fatal("float32 angle constants are inconsistent")
	}
	if OneTenthOfADegree != DegreesToRadians/10 {
		t.Fatal("OneTenthOfADegree is inconsistent")
	}
	if FeetToMillimeters != 1/MillimetersToFeet {
		t.Fatal("length conversion constants are not reciprocal")
	}
}

func TestEnumValues(t *testing.T) {
	containment := []Containment{ContainmentDisjoint, ContainmentContains, ContainmentIntersects}
	for want, got := range containment {
		if got != Containment(want) {
			t.Fatalf("containment value %d = %d", want, got)
		}
	}

	intersections := []PlaneIntersection{
		PlaneIntersectionFront,
		PlaneIntersectionBack,
		PlaneIntersectionIntersecting,
	}
	for want, got := range intersections {
		if got != PlaneIntersection(want) {
			t.Fatalf("plane intersection value %d = %d", want, got)
		}
	}
}
