package math3d

import "testing"

func TestVec3CatmullRom(t *testing.T) {
	value1 := V3(0, 0, 0)
	value2 := V3(1, 2, 3)
	value3 := V3(4, 6, 8)
	value4 := V3(9, 12, 15)

	if got := value1.CatmullRom(value2, value3, value4, 0); got != value2 {
		t.Fatalf("CatmullRom at zero = %#v, want %#v", got, value2)
	}
	if got := value1.CatmullRom(value2, value3, value4, 1); got != value3 {
		t.Fatalf("CatmullRom at one = %#v, want %#v", got, value3)
	}
	if got, want := value1.CatmullRom(value2, value3, value4, 0.5), V3(2.25, 3.75, 5.25); !got.AlmostEqual(want, 1e-6) {
		t.Fatalf("CatmullRom representative vector = %#v, want %#v", got, want)
	}
	if got, want := value1.CatmullRom(value2, value3, value4, 1.5), V3(6.25, 8.75, 11.25); !got.AlmostEqual(want, 1e-6) {
		t.Fatalf("CatmullRom extrapolation = %#v, want %#v", got, want)
	}

	amount := float32(0.37)
	got := value1.CatmullRom(value2, value3, value4, amount)
	want := V3(
		CatmullRom(value1.X, value2.X, value3.X, value4.X, amount),
		CatmullRom(value1.Y, value2.Y, value3.Y, value4.Y, amount),
		CatmullRom(value1.Z, value2.Z, value3.Z, value4.Z, amount),
	)
	if got != want {
		t.Fatalf("CatmullRom components = %#v, want %#v", got, want)
	}
}

func TestVec3Hermite(t *testing.T) {
	value1 := V3(1, 2, 3)
	tangent1 := V3(2, -3, 1)
	value2 := V3(5, 8, -2)
	tangent2 := V3(-1, 4, 2)

	if got := value1.Hermite(tangent1, value2, tangent2, 0); got != value1 {
		t.Fatalf("Hermite at zero = %#v, want %#v", got, value1)
	}
	if got := value1.Hermite(tangent1, value2, tangent2, 1); got != value2 {
		t.Fatalf("Hermite at one = %#v, want %#v", got, value2)
	}
	if got, want := value1.Hermite(tangent1, value2, tangent2, 0.25), V3(1.953125, 2.328125, 2.265625); !got.AlmostEqual(want, 1e-6) {
		t.Fatalf("Hermite representative vector = %#v, want %#v", got, want)
	}
	if got, want := value1.Hermite(tangent1, value2, tangent2, 1.5), V3(0.625, 5.375, 5.625); !got.AlmostEqual(want, 1e-6) {
		t.Fatalf("Hermite extrapolation = %#v, want %#v", got, want)
	}

	amount := float32(0.37)
	got := value1.Hermite(tangent1, value2, tangent2, amount)
	want := V3(
		Hermite(value1.X, tangent1.X, value2.X, tangent2.X, amount),
		Hermite(value1.Y, tangent1.Y, value2.Y, tangent2.Y, amount),
		Hermite(value1.Z, tangent1.Z, value2.Z, tangent2.Z, amount),
	)
	if got != want {
		t.Fatalf("Hermite components = %#v, want %#v", got, want)
	}
}

func TestVec3SmoothStep(t *testing.T) {
	start, end := V3(1, -2, 4), V3(5, 8, -6)

	for _, test := range []struct {
		name   string
		amount float32
		want   Vec3
	}{
		{"below zero", -1, start},
		{"zero", 0, start},
		{"one", 1, end},
		{"above one", 2, end},
		{"midpoint", 0.5, V3(3, 3, -1)},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := start.SmoothStep(end, test.amount); got != test.want {
				t.Fatalf("SmoothStep(%v) = %#v, want %#v", test.amount, got, test.want)
			}
		})
	}

	amount := float32(0.37)
	got := start.SmoothStep(end, amount)
	want := V3(
		SmoothStep(start.X, end.X, amount),
		SmoothStep(start.Y, end.Y, amount),
		SmoothStep(start.Z, end.Z, amount),
	)
	if got != want {
		t.Fatalf("SmoothStep components = %#v, want %#v", got, want)
	}
}

func TestVec3InterpolationDoesNotMutateInputs(t *testing.T) {
	values := [4]Vec3{V3(1, 2, 3), V3(4, 5, 6), V3(7, 8, 9), V3(10, 11, 12)}
	want := values
	_ = values[0].CatmullRom(values[1], values[2], values[3], 0.25)
	_ = values[0].Hermite(values[1], values[2], values[3], 0.25)
	_ = values[0].SmoothStep(values[1], 0.25)
	if values != want {
		t.Fatalf("interpolation inputs changed: %#v", values)
	}
}
