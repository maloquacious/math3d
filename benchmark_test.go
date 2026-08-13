package math3d

import "testing"

var (
	benchmarkVec3A  = V3(1.25, -2.5, 4.75)
	benchmarkVec3B  = V3(-3.5, 6.25, 0.5)
	benchmarkVec3   Vec3
	benchmarkMat4   Mat4
	benchmarkScalar float32
	benchmarkOK     bool
)

func BenchmarkVec3Cross(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkVec3 = benchmarkVec3A.Cross(benchmarkVec3B)
	}
}

func BenchmarkMat4Mul(b *testing.B) {
	a := ComposeMat4(V3(2, 3, 4), QuatFromYawPitchRoll(0.2, 0.4, 0.6), V3(10, 20, 30))
	c := ComposeMat4(V3(0.5, 0.25, 2), QuatFromYawPitchRoll(-0.1, 0.3, -0.5), V3(-3, 8, 2))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkMat4 = a.Mul(c)
	}
}

func BenchmarkMat4Inverse(b *testing.B) {
	matrix := ComposeMat4(V3(2, 3, 4), QuatFromYawPitchRoll(0.2, 0.4, 0.6), V3(10, 20, 30))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkMat4, benchmarkOK = matrix.Inverse()
	}
}

func BenchmarkRayIntersectBox(b *testing.B) {
	ray := NewRay(V3(-10, 2, 3), V3(2, -0.25, 0.5))
	box := NewBox3(V3(-2, -2, -2), V3(4, 4, 4))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkScalar, benchmarkOK = ray.IntersectBox(box)
	}
}
