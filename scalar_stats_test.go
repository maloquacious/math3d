package math3d

import (
	"math"
	"testing"
)

func TestStatsOfVectors(t *testing.T) {
	values := []Vec3{V3(2, 8, 4), V3(6, 4, 10), V3(4, 6, 7)}
	stats := StatsOf(values)

	if stats.Count != 3 || stats.Min != V3(2, 4, 4) || stats.Max != V3(6, 8, 10) || stats.Sum != V3(12, 18, 21) {
		t.Fatalf("unexpected stats: %+v", stats)
	}
	average, ok := StatsAverage(stats)
	if !ok || average != V3(4, 6, 7) {
		t.Fatalf("unexpected average: %v, %v", average, ok)
	}
	if extents := StatsExtents(stats); extents != V3(4, 4, 6) {
		t.Fatalf("unexpected extents: %v", extents)
	}
	if middle := StatsMiddle(stats); middle != V3(4, 6, 7) {
		t.Fatalf("unexpected middle: %v", middle)
	}
}

func TestStatsDoNotBiasBoundsTowardZero(t *testing.T) {
	positive := StatsOf([]DVec2{DV2(2, 4), DV2(3, 8)})
	if positive.Min != DV2(2, 4) || positive.Max != DV2(3, 8) {
		t.Fatalf("positive bounds were biased: %+v", positive)
	}
	negative := StatsOf([]Vec2{V2(-8, -3), V2(-2, -1)})
	if negative.Min != V2(-8, -3) || negative.Max != V2(-2, -1) {
		t.Fatalf("negative bounds were biased: %+v", negative)
	}

	empty := StatsOf([]Vec4(nil))
	if empty != (Stats[Vec4]{}) {
		t.Fatalf("empty stats = %+v", empty)
	}
	if average, ok := StatsAverage(empty); ok || average != (Vec4{}) {
		t.Fatalf("empty average = %v, %v", average, ok)
	}
}

func TestDomainAndScalarHelpers(t *testing.T) {
	domain := NewDomain(2, 10)
	if got := domain.Normalize(-1); got != 0.2 {
		t.Fatalf("Normalize below domain = %v", got)
	}
	if got := domain.Normalize(5); got != 0.5 {
		t.Fatalf("Normalize within domain = %v", got)
	}
	if got := domain.Normalize(20); got != 1 {
		t.Fatalf("Normalize above domain = %v", got)
	}

	if Percentage(8, 2) != 25 || math.Abs(ToDegrees(ToRadians(90))-90) > 1e-12 {
		t.Fatal("percentage or angle conversion failed")
	}
	if Smoothstep(0.5) != 0.5 || SmoothStep(10, 20, -1) != 10 || SmoothStep(10, 20, 2) != 20 {
		t.Fatal("smooth-step behavior failed")
	}
	if CatmullRom(0, 10, 20, 30, 0.5) != 15 {
		t.Fatal("CatmullRom midpoint failed")
	}
	if Hermite(2, 5, 8, -3, 0) != 2 || Hermite(2, 5, 8, -3, 1) != 8 {
		t.Fatal("Hermite endpoints failed")
	}
	if WrapAngle(-Pi) != Pi || WrapAngle(Pi) != Pi {
		t.Fatalf("WrapAngle boundary failed: %v %v", WrapAngle(-Pi), WrapAngle(Pi))
	}
	if !IsNonZeroAndValid(1, Tolerance) || IsNonZeroAndValid(Tolerance, Tolerance) || IsNonZeroAndValid(float32(math.Inf(1)), Tolerance) {
		t.Fatal("IsNonZeroAndValid failed")
	}
	if MillimetersToFeet*FeetToMillimeters != 1 {
		t.Fatal("unit conversion constants have the wrong direction")
	}
}

func TestNearestPowerOfTwo(t *testing.T) {
	tests := []struct {
		name  string
		value int32
		want  int32
		ok    bool
	}{
		{name: "negative", value: -1},
		{name: "zero"},
		{name: "one", value: 1, want: 1, ok: true},
		{name: "exact power", value: 1024, want: 1024, ok: true},
		{name: "below logarithmic midpoint", value: 5, want: 4, ok: true},
		{name: "above logarithmic midpoint", value: 6, want: 8, ok: true},
		{name: "largest non-overflowing input", value: 1_518_500_249, want: 1 << 30, ok: true},
		{name: "rounds beyond int32", value: 1_518_500_250},
		{name: "maximum int32", value: math.MaxInt32},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := NearestPowerOfTwo(test.value)
			if got != test.want || ok != test.ok {
				t.Fatalf("NearestPowerOfTwo(%d) = (%d, %v), want (%d, %v)", test.value, got, ok, test.want, test.ok)
			}
		})
	}
}

func TestHashAndStatelessRandom(t *testing.T) {
	if got := CombineHash(1, 2); got != 35 {
		t.Fatalf("CombineHash(1, 2) = %d", got)
	}
	if got := CombineHash(-1, 0); got != -2 {
		t.Fatalf("CombineHash overflow behavior = %d", got)
	}
	if got := HashInts(1, 2); got != 35 {
		t.Fatalf("HashInts(1, 2) = %d", got)
	}
	if HashInts() != 0 || RandomInt(2, 1) != 35 || RandomUint(2, 1) != 35 {
		t.Fatal("hash/random compatibility failed")
	}
	if RandomUnitFloat(0, 0) != 0 || RandomUnitFloat(-1, 0) != 1 {
		t.Fatal("random unit interval endpoints failed")
	}

	value := RandomFloat(-4, 6, 12, 7)
	if value < -4 || value > 6 || value != RandomFloat(-4, 6, 12, 7) {
		t.Fatalf("random float is out of range or nondeterministic: %v", value)
	}
	vector := RandomVec3(5, 9)
	if vector != V3(RandomUnitFloat(15, 9), RandomUnitFloat(16, 9), RandomUnitFloat(17, 9)) {
		t.Fatalf("random vector used wrong indices: %v", vector)
	}
}

func TestRemainingVectorConversions(t *testing.T) {
	if IV2(1, -2).Vec2() != V2(1, -2) || IV3(1, -2, 3).Vec3() != V3(1, -2, 3) {
		t.Fatal("integer vector conversion failed")
	}
	if V2(1, 2).Vec3() != V3(1, 2, 0) || V2(1, 2).Vec4() != V4(1, 2, 0, 0) {
		t.Fatal("Vec2 promotion failed")
	}
	v3 := V3(1, 2, 3)
	if v3.XY() != V2(1, 2) || v3.XZ() != V2(1, 3) || v3.YZ() != V2(2, 3) || v3.ZYX() != V3(3, 2, 1) || v3.Vec4() != V4(1, 2, 3, 0) {
		t.Fatal("Vec3 conversion or swizzle failed")
	}
	if V4(1, 2, 3, 4).Vec2() != V2(1, 2) || V4(1, 2, 3, 4).Vec3() != V3(1, 2, 3) {
		t.Fatal("Vec4 projection failed")
	}

	horizontal := NewHorizontal(1.25, -2.5)
	if HorizontalFromDVec2(horizontal.DVec2()) != horizontal || !HorizontalFromVec2(horizontal.Vec2()).AlmostEqual(horizontal, 1e-6) {
		t.Fatal("horizontal conversion failed")
	}
}
