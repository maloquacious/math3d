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

// Point is a three-dimensional coordinate.
type Point struct {
	X, Y, Z float64
}

// DeltaXYZ returns the changes in x, y, and z between two points.
func (p Point) DeltaXYZ(p2 Point) (dx, dy, dz float64) {
	return p2.X - p.X, p2.Y - p.Y, p2.Z - p.Z
}

// Distance returns the distance between two points
func (p Point) Distance(p2 Point) float64 {
	dx, dy, dz := p.DeltaXYZ(p2)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// PointSlope returns a function to produce points on the line connecting two points.
//
//	⟨mx,my,mz⟩ = ⟨x1,y1,z1⟩ −  ⟨x0,y0,z0⟩
//	⟨x ,y ,z ⟩ = ⟨x0,y0,z0⟩ + t⟨mx,my,mz⟩
func (p Point) PointSlope(p2 Point) func(t float64) (tx, ty, tz float64) {
	// https://math.stackexchange.com/questions/799783/slope-of-a-line-in-3d-coordinate-system
	mx, my, mz := p.DeltaXYZ(p2)
	return func(t float64) (float64, float64, float64) {
		return p.X + t*mx, p.Y + t*my, p.Z + t*mz
	}
}

// Slope returns the slope (really, the direction cosines) of the line connecting two points.
func (p Point) Slope(p2 Point) (xy, xz, yz float64) {
	// https://math.stackexchange.com/questions/799783/slope-of-a-line-in-3d-coordinate-system
	dx, dy, dz := p.DeltaXYZ(p2)
	d := p.Distance(p2)
	return dz / d, dy / d, dx / d
}
