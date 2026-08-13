package math3d

import "math"

// Coordinate angles are stored as float64. Upstream does not document their
// units; callers must not assume the package-wide rotation convention applies
// until conversion algorithms define that contract.

// Spherical is the upstream SphericalCoordinate value.
type Spherical struct {
	Radius, Azimuth, Inclination float64
}

// NewSpherical constructs a Spherical value.
func NewSpherical(radius, azimuth, inclination float64) Spherical {
	return Spherical{Radius: radius, Azimuth: azimuth, Inclination: inclination}
}

// AlmostEqual compares corresponding components using a strict absolute tolerance.
func (v Spherical) AlmostEqual(other Spherical, tolerance float64) bool {
	return almostEqual64(tolerance, v.Radius-other.Radius, v.Azimuth-other.Azimuth, v.Inclination-other.Inclination)
}

// Polar is the upstream PolarCoordinate value.
type Polar struct {
	Radius, Azimuth float64
}

// NewPolar constructs a Polar value.
func NewPolar(radius, azimuth float64) Polar { return Polar{Radius: radius, Azimuth: azimuth} }

// AlmostEqual compares corresponding components using a strict absolute tolerance.
func (v Polar) AlmostEqual(other Polar, tolerance float64) bool {
	return almostEqual64(tolerance, v.Radius-other.Radius, v.Azimuth-other.Azimuth)
}

// LogPolar is the upstream LogPolarCoordinate value.
type LogPolar struct {
	Rho, Azimuth float64
}

// NewLogPolar constructs a LogPolar value.
func NewLogPolar(rho, azimuth float64) LogPolar { return LogPolar{Rho: rho, Azimuth: azimuth} }

// AlmostEqual compares corresponding components using a strict absolute tolerance.
func (v LogPolar) AlmostEqual(other LogPolar, tolerance float64) bool {
	return almostEqual64(tolerance, v.Rho-other.Rho, v.Azimuth-other.Azimuth)
}

// Cylindrical is the upstream CylindricalCoordinate value.
type Cylindrical struct {
	Radius, Azimuth, Height float64
}

// NewCylindrical constructs a Cylindrical value.
func NewCylindrical(radius, azimuth, height float64) Cylindrical {
	return Cylindrical{Radius: radius, Azimuth: azimuth, Height: height}
}

// AlmostEqual compares corresponding components using a strict absolute tolerance.
func (v Cylindrical) AlmostEqual(other Cylindrical, tolerance float64) bool {
	return almostEqual64(tolerance, v.Radius-other.Radius, v.Azimuth-other.Azimuth, v.Height-other.Height)
}

// Horizontal is the upstream HorizontalCoordinate value.
type Horizontal struct {
	Azimuth, Inclination float64
}

// NewHorizontal constructs a Horizontal value.
func NewHorizontal(azimuth, inclination float64) Horizontal {
	return Horizontal{Azimuth: azimuth, Inclination: inclination}
}

// AlmostEqual compares corresponding components using a strict absolute tolerance.
func (v Horizontal) AlmostEqual(other Horizontal, tolerance float64) bool {
	return almostEqual64(tolerance, v.Azimuth-other.Azimuth, v.Inclination-other.Inclination)
}

// DVec2 converts azimuth and inclination to X and Y.
func (v Horizontal) DVec2() DVec2 { return DV2(v.Azimuth, v.Inclination) }

// Vec2 converts azimuth and inclination to float32 X and Y.
func (v Horizontal) Vec2() Vec2 { return V2(float32(v.Azimuth), float32(v.Inclination)) }

// HorizontalFromDVec2 maps X to azimuth and Y to inclination.
func HorizontalFromDVec2(v DVec2) Horizontal { return NewHorizontal(v.X, v.Y) }

// HorizontalFromVec2 maps X to azimuth and Y to inclination.
func HorizontalFromVec2(v Vec2) Horizontal {
	return NewHorizontal(float64(v.X), float64(v.Y))
}

// Geo is the upstream GeoCoordinate value. Upstream does not specify whether
// latitude and longitude are expressed in degrees or radians.
type Geo struct {
	Latitude, Longitude float64
}

// NewGeo constructs a Geo value.
func NewGeo(latitude, longitude float64) Geo {
	return Geo{Latitude: latitude, Longitude: longitude}
}

// AlmostEqual compares corresponding components using a strict absolute tolerance.
func (v Geo) AlmostEqual(other Geo, tolerance float64) bool {
	return almostEqual64(tolerance, v.Latitude-other.Latitude, v.Longitude-other.Longitude)
}

func almostEqual64(tolerance float64, differences ...float64) bool {
	for _, difference := range differences {
		if !(math.Abs(difference) < tolerance) {
			return false
		}
	}
	return true
}
