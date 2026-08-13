package math3d

// Mathematical constants preserve the precision and values of the upstream
// declarations.
const (
	Pi        float32 = 3.1415927
	HalfPi    float32 = Pi / 2
	TwoPi     float32 = Pi * 2
	Tolerance float32 = 0.0000001
	Log10E    float32 = 0.4342945
	Log2E     float32 = 1.442695
	E         float32 = 2.7182817

	RadiansToDegrees  float64 = 57.295779513082320876798154814105
	DegreesToRadians  float64 = 0.017453292519943295769236907684886
	OneTenthOfADegree float64 = DegreesToRadians / 10

	MillimetersToFeet float64 = 0.00328084
	FeetToMillimeters float64 = 1 / MillimetersToFeet
)

// Containment describes the spatial relationship between two shapes.
type Containment uint8

const (
	// ContainmentDisjoint is the zero value.
	ContainmentDisjoint Containment = iota
	ContainmentContains
	ContainmentIntersects
)

// PlaneIntersection describes which side of a plane a shape occupies.
type PlaneIntersection uint8

const (
	// PlaneIntersectionFront is the zero value.
	PlaneIntersectionFront PlaneIntersection = iota
	PlaneIntersectionBack
	PlaneIntersectionIntersecting
)
