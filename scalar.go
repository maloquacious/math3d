package math3d

import (
	"math"
	"math/bits"
)

// Domain is the upstream ValueDomain range.
type Domain struct {
	Lower, Upper float64
}

// NewDomain constructs a Domain without reordering its bounds.
func NewDomain(lower, upper float64) Domain { return Domain{Lower: lower, Upper: upper} }

// Normalize preserves the upstream ValueDomain behavior: it clamps value to
// the domain and divides by Upper. It is not conventional inverse lerp unless
// Lower is zero.
func (d Domain) Normalize(value float64) float64 {
	return max(min(value, d.Upper), d.Lower) / d.Upper
}

// Percentage returns numerator as a percentage of denominator.
func Percentage(denominator, numerator float64) float64 { return numerator / denominator * 100 }

// ToRadians converts degrees to radians.
func ToRadians(degrees float64) float64 { return degrees * DegreesToRadians }

// ToDegrees converts radians to degrees.
func ToDegrees(radians float64) float64 { return radians * RadiansToDegrees }

// Smoothstep evaluates x²(3-2x) without clamping x.
func Smoothstep(x float64) float64 { return x * x * (3 - 2*x) }

// CatmullRom interpolates between value2 and value3 using value1 and value4 as
// neighboring control values. Amount is not clamped.
func CatmullRom(value1, value2, value3, value4, amount float32) float32 {
	t := float64(amount)
	t2 := t * t
	t3 := t2 * t
	return float32(0.5 * (2*float64(value2) +
		(float64(value3)-float64(value1))*t +
		(2*float64(value1)-5*float64(value2)+4*float64(value3)-float64(value4))*t2 +
		(3*float64(value2)-float64(value1)-3*float64(value3)+float64(value4))*t3))
}

// Hermite performs cubic Hermite interpolation. Amount is not clamped.
func Hermite(value1, tangent1, value2, tangent2, amount float32) float32 {
	if amount == 0 {
		return value1
	}
	if amount == 1 {
		return value2
	}
	t := float64(amount)
	t2 := t * t
	t3 := t2 * t
	v1, m1 := float64(value1), float64(tangent1)
	v2, m2 := float64(value2), float64(tangent2)
	return float32((2*v1-2*v2+m2+m1)*t3 + (3*v2-3*v1-2*m1-m2)*t2 + m1*t + v1)
}

// SmoothStep smoothly interpolates from value1 to value2, clamping amount to
// [0, 1].
func SmoothStep(value1, value2, amount float32) float32 {
	return Hermite(value1, 0, value2, 0, max(float32(0), min(amount, float32(1))))
}

// WrapAngle wraps an angle in radians to (-Pi, Pi].
func WrapAngle(angle float32) float32 {
	angle = float32(math.Mod(float64(angle), float64(TwoPi)))
	if angle <= -Pi {
		return angle + TwoPi
	}
	if angle > Pi {
		return angle - TwoPi
	}
	return angle
}

// IsNonZeroAndValid reports whether value is finite and its absolute value is
// greater than tolerance.
func IsNonZeroAndValid(value, tolerance float32) bool {
	return finite32(value) && abs32(value) > tolerance
}

// CombineHash combines two source-compatible signed 32-bit hash values.
func CombineHash(first, second int32) int32 {
	rotated := bits.RotateLeft32(uint32(first), 5)
	return int32((rotated + uint32(first)) ^ uint32(second))
}

// HashInts combines all values in order from a zero seed.
func HashInts(values ...int32) int32 {
	var hash int32
	for _, value := range values {
		hash = CombineHash(hash, value)
	}
	return hash
}

// RandomInt returns a deterministic source-compatible value for index and seed.
func RandomInt(index, seed int32) int32 { return CombineHash(seed, index) }

// RandomUint returns the unsigned representation of RandomInt.
func RandomUint(index, seed int32) uint32 { return uint32(RandomInt(index, seed)) }

// RandomFloat returns a deterministic value in the inclusive range [minimum,
// maximum] when minimum <= maximum.
func RandomFloat(minimum, maximum float32, index, seed int32) float32 {
	unit := float32(RandomUint(index, seed)) / float32(math.MaxUint32)
	return unit*(maximum-minimum) + minimum
}

// RandomUnitFloat returns a deterministic value in [0, 1].
func RandomUnitFloat(index, seed int32) float32 { return RandomFloat(0, 1, index, seed) }

// RandomVec2 returns a vector whose components use consecutive derived indices.
func RandomVec2(index, seed int32) Vec2 {
	base := int32(uint32(index) * 2)
	return V2(RandomUnitFloat(base, seed), RandomUnitFloat(int32(uint32(base)+1), seed))
}

// RandomVec3 returns a vector whose components use consecutive derived indices.
func RandomVec3(index, seed int32) Vec3 {
	base := int32(uint32(index) * 3)
	return V3(RandomUnitFloat(base, seed), RandomUnitFloat(int32(uint32(base)+1), seed), RandomUnitFloat(int32(uint32(base)+2), seed))
}

// RandomVec4 returns a vector whose components use consecutive derived indices.
func RandomVec4(index, seed int32) Vec4 {
	base := int32(uint32(index) * 4)
	return V4(RandomUnitFloat(base, seed), RandomUnitFloat(int32(uint32(base)+1), seed), RandomUnitFloat(int32(uint32(base)+2), seed), RandomUnitFloat(int32(uint32(base)+3), seed))
}
