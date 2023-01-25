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

package math3d_test

import (
	"github.com/maloquacious/math3d"
	"testing"
)

func TestPoints(t *testing.T) {
	p1, p2 := math3d.Point{0, 0, 0}, math3d.Point{1, 1, 1}

	dx, dy, dz := p1.DeltaXYZ(p2)
	if dx != 1.0 {
		t.Errorf("DeltaXYZ: x: wanted %f, got %f\n", 1.0, dx)
	}
	if dy != 1.0 {
		t.Errorf("DeltaXYZ: y: wanted %f, got %f\n", 1.0, dy)
	}
	if dz != 1.0 {
		t.Errorf("DeltaXYZ: z: wanted %f, got %f\n", 1.0, dz)
	}

	for _, tt := range []struct {
		p      math3d.Point
		expect float64
	}{
		{p1, 0},
		{math3d.Point{1, 0, 0}, 1},
		{math3d.Point{0, 1, 0}, 1},
		{math3d.Point{0, 0, 1}, 1},
	} {
		d := p1.Distance(tt.p)
		if d != tt.expect {
			t.Errorf("Distance: z: wanted %f, got %f\n", tt.expect, d)
		}
	}
}
