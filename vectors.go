/*
 * Copyright (c) 2023 Michael D Henderson
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package math3d

import "math"

// Vector implements a vector with length (magnitude) and direction
type Vector []float64

func NewVector(f ...float64) Vector {
	return append(Vector{}, f...)
}

func StandardBasisVector(n int) []Vector {
	v := make([]Vector, n, n)
	for i := 0; i < n; i++ {
		v[i] = make(Vector, n, n)
		v[i][i] = 1
	}
	return v
}

func UnitVector(n int) Vector {
	v := make(Vector, n, n)
	return v
}

func ZeroVector(n int) Vector {
	v := make(Vector, n, n)
	return v
}

func (v Vector) Add(w Vector) Vector {
	u := make(Vector, len(v), len(v))
	for i, s := range v {
		u[i] = s + w[i]
	}
	return u
}

func (v Vector) Div(scalar float64) Vector {
	u := make(Vector, len(v), len(v))
	for i, s := range v {
		u[i] = s / scalar
	}
	return u
}

func (v Vector) Dot(w Vector) Vector {
	u := make(Vector, len(v), len(v))
	for i, s := range v {
		u[i] = s * w[i]
	}
	return u
}

func (v Vector) IsZero() bool {
	for _, s := range v {
		if s != 0 {
			return false
		}
	}
	return true
}

// Length implements the Euclidean norm of the vector.
func (v Vector) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

// LengthSquared implements the square of the Euclidean norm of the vector.
func (v Vector) LengthSquared() float64 {
	var sum float64
	for _, s := range v {
		sum = sum + s*s
	}
	return sum
}

// ManhattanDistance implements the step-wise total distance of the vector.
func (v Vector) ManhattanDistance() float64 {
	var sum float64
	for _, s := range v {
		sum = sum + math.Abs(s)
	}
	return sum
}

func (v Vector) Mul(scalar float64) Vector {
	u := make(Vector, len(v), len(v))
	for i, s := range v {
		u[i] = s * scalar
	}
	return u
}

// Normalize returns a vector with all components divided by the vector's length
func (v Vector) Normalize() Vector {
	if v.IsZero() {
		return v.ZeroVector()
	}
	// we'll multiply by the reciprocal instead of dividing as a performance hack
	reciprocal := 1.0 / v.Length()
	return v.Mul(reciprocal)
}

func (v Vector) StandardBasis() []Vector {
	return StandardBasisVector(len(v))
}

func (v Vector) Sub(w Vector) Vector {
	u := make(Vector, len(v), len(v))
	for i, s := range v {
		u[i] = s - w[i]
	}
	return u
}

func (v Vector) UnitVector() Vector {
	return UnitVector(len(v))
}

func (v Vector) ZeroVector() Vector {
	return ZeroVector(len(v))
}

func (v Vector) toVec2() Vec2 {
	return Vec2{X: v[0], Y: v[1]}
}

func (v Vector) toVec3() Vec3 {
	return Vec3{X: v[0], Y: v[1], Z: v[2]}
}

func (v Vector) toVec4() Vec4 {
	return Vec4{W: v[0], X: v[1], Y: v[2], Z: v[3]}
}

// Vec2 implements a vector with length (magnitude) and direction.
type Vec2 struct {
	X, Y float64
}

func NewVec2(x, y float64) Vec2 {
	return Vec2{X: x, Y: y}
}

func StandardBasisVec2() []Vec2 {
	sb := StandardBasisVector(2)
	return []Vec2{
		sb[0].toVec2(),
		sb[1].toVec2(),
	}
}

func UnitVec2() Vec2 {
	return UnitVector(2).toVec2()
}

func (v Vec2) Add(w Vec2) Vec2 {
	return v.toVector().Add(w.toVector()).toVec2()
}

func (v Vec2) Div(scalar float64) Vec2 {
	return v.toVector().Div(scalar).toVec2()
}

func (v Vec2) Dot(w Vec2) Vec2 {
	return v.toVector().Dot(w.toVector()).toVec2()
}

func (v Vec2) IsZero() bool {
	return v.toVector().IsZero()
}

// Length implements the Euclidean norm of the vector.
func (v Vec2) Length() float64 {
	return v.toVector().Length()
}

// LengthSquared implements the square of the Euclidean norm of the vector.
func (v Vec2) LengthSquared() float64 {
	return v.toVector().LengthSquared()
}

// ManhattanDistance implements the step-wise total distance of the vector.
func (v Vec2) ManhattanDistance() float64 {
	return v.toVector().ManhattanDistance()
}

func (v Vec2) Mul(scalar float64) Vec2 {
	return v.toVector().Mul(scalar).toVec2()
}

func (v Vec2) Normalize() Vec2 {
	return v.toVector().Normalize().toVec2()
}

func (v Vec2) StandardBasis() []Vec2 {
	return StandardBasisVec2()
}

func (v Vec2) Sub(w Vec2) Vec2 {
	return v.toVector().Sub(w.toVector()).toVec2()
}

func (v Vec2) UnitVector() Vec2 {
	return v.toVector().UnitVector().toVec2()
}

func (v Vec2) ZeroVector() Vec2 {
	return Vec2{}
}

func (v Vec2) toVector() Vector {
	return Vector{v.X, v.Y}
}

// Vec3 implements a vector with length (magnitude) and direction.
type Vec3 struct {
	X, Y, Z float64
}

func (v Vec3) toVector() Vector {
	return Vector{v.X, v.Y, v.Z}
}

func NewVec3(x, y, z float64) Vec3 {
	return Vec3{X: x, Y: y, Z: z}
}

func StandardBasisVec3() []Vec3 {
	sb := StandardBasisVector(3)
	return []Vec3{
		sb[0].toVec3(),
		sb[1].toVec3(),
		sb[2].toVec3(),
	}
}

// Vec4 implements a vector with length (magnitude) and direction.
type Vec4 struct {
	W, X, Y, Z float64
}

func (v Vec4) toVector() Vector {
	return Vector{v.W, v.X, v.Y, v.Z}
}

func NewVec4(w, x, y, z float64) Vec4 {
	return Vec4{W: w, X: x, Y: y, Z: z}
}

func StandardBasisVec4() []Vec4 {
	sb := StandardBasisVector(4)
	return []Vec4{
		sb[0].toVec4(),
		sb[1].toVec4(),
		sb[2].toVec4(),
		sb[3].toVec4(),
	}
}
