package math3d

import (
	"math"
	"testing"
)

func TestCoordinateConstructors(t *testing.T) {
	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Spherical", NewSpherical(1, 2, 3), Spherical{Radius: 1, Azimuth: 2, Inclination: 3}},
		{"Polar", NewPolar(1, 2), Polar{Radius: 1, Azimuth: 2}},
		{"LogPolar", NewLogPolar(1, 2), LogPolar{Rho: 1, Azimuth: 2}},
		{"Cylindrical", NewCylindrical(1, 2, 3), Cylindrical{Radius: 1, Azimuth: 2, Height: 3}},
		{"Horizontal", NewHorizontal(1, 2), Horizontal{Azimuth: 1, Inclination: 2}},
		{"Geo", NewGeo(1, 2), Geo{Latitude: 1, Longitude: 2}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.got != test.want {
				t.Fatalf("got %#v, want %#v", test.got, test.want)
			}
		})
	}
}

func TestCoordinateAlmostEqual(t *testing.T) {
	tests := []struct {
		name  string
		equal func(float64) bool
	}{
		{"Spherical", func(delta float64) bool { return NewSpherical(1, 2, 3).AlmostEqual(NewSpherical(1+delta, 2, 3), 0.1) }},
		{"Polar", func(delta float64) bool { return NewPolar(1, 2).AlmostEqual(NewPolar(1+delta, 2), 0.1) }},
		{"LogPolar", func(delta float64) bool { return NewLogPolar(1, 2).AlmostEqual(NewLogPolar(1+delta, 2), 0.1) }},
		{"Cylindrical", func(delta float64) bool {
			return NewCylindrical(1, 2, 3).AlmostEqual(NewCylindrical(1+delta, 2, 3), 0.1)
		}},
		{"Horizontal", func(delta float64) bool { return NewHorizontal(1, 2).AlmostEqual(NewHorizontal(1+delta, 2), 0.1) }},
		{"Geo", func(delta float64) bool { return NewGeo(1, 2).AlmostEqual(NewGeo(1+delta, 2), 0.1) }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !test.equal(0.05) {
				t.Fatal("near values should compare approximately equal")
			}
			if test.equal(math.NaN()) {
				t.Fatal("NaN must not compare approximately equal")
			}
		})
	}
}
