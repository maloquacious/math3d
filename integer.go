package math3d

// Byte2 is the upstream Byte2 value.
type Byte2 struct {
	X, Y byte
}

// NewByte2 constructs a Byte2.
func NewByte2(x, y byte) Byte2 { return Byte2{X: x, Y: y} }

// Byte3 is the upstream Byte3 value.
type Byte3 struct {
	X, Y, Z byte
}

// NewByte3 constructs a Byte3.
func NewByte3(x, y, z byte) Byte3 { return Byte3{X: x, Y: y, Z: z} }

// Byte4 is the upstream Byte4 value.
type Byte4 struct {
	X, Y, Z, W byte
}

// NewByte4 constructs a Byte4.
func NewByte4(x, y, z, w byte) Byte4 { return Byte4{X: x, Y: y, Z: z, W: w} }

// IVec2 is the upstream Int2 value. Its components retain C#'s 32-bit width.
type IVec2 struct {
	X, Y int32
}

// IV2 constructs an IVec2.
func IV2(x, y int32) IVec2 { return IVec2{X: x, Y: y} }

// SplatIV2 constructs an IVec2 whose components are value.
func SplatIV2(value int32) IVec2 { return IV2(value, value) }

// Vec2 converts v to float32 precision.
func (v IVec2) Vec2() Vec2 { return V2(float32(v.X), float32(v.Y)) }

// IVec3 is the upstream Int3 value. Its components retain C#'s 32-bit width.
type IVec3 struct {
	X, Y, Z int32
}

// IV3 constructs an IVec3.
func IV3(x, y, z int32) IVec3 { return IVec3{X: x, Y: y, Z: z} }

// SplatIV3 constructs an IVec3 whose components are value.
func SplatIV3(value int32) IVec3 { return IV3(value, value, value) }

// Vec3 converts v to float32 precision.
func (v IVec3) Vec3() Vec3 { return V3(float32(v.X), float32(v.Y), float32(v.Z)) }

// IVec4 is the upstream Int4 value. Its components retain C#'s 32-bit width.
type IVec4 struct {
	X, Y, Z, W int32
}

// IV4 constructs an IVec4.
func IV4(x, y, z, w int32) IVec4 { return IVec4{X: x, Y: y, Z: z, W: w} }

// SplatIV4 constructs an IVec4 whose components are value.
func SplatIV4(value int32) IVec4 { return IV4(value, value, value, value) }
