package math3d

import "testing"

func TestByteConstructors(t *testing.T) {
	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Byte2", NewByte2(1, 2), Byte2{X: 1, Y: 2}},
		{"Byte3", NewByte3(1, 2, 3), Byte3{X: 1, Y: 2, Z: 3}},
		{"Byte4", NewByte4(1, 2, 3, 4), Byte4{X: 1, Y: 2, Z: 3, W: 4}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.got != test.want {
				t.Fatalf("got %#v, want %#v", test.got, test.want)
			}
		})
	}
}

func TestIntegerVectorConstructors(t *testing.T) {
	tests := []struct {
		name string
		got  any
		want any
	}{
		{"IV2", IV2(1, 2), IVec2{X: 1, Y: 2}},
		{"IV3", IV3(1, 2, 3), IVec3{X: 1, Y: 2, Z: 3}},
		{"IV4", IV4(1, 2, 3, 4), IVec4{X: 1, Y: 2, Z: 3, W: 4}},
		{"SplatIV2", SplatIV2(7), IVec2{X: 7, Y: 7}},
		{"SplatIV3", SplatIV3(7), IVec3{X: 7, Y: 7, Z: 7}},
		{"SplatIV4", SplatIV4(7), IVec4{X: 7, Y: 7, Z: 7, W: 7}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.got != test.want {
				t.Fatalf("got %#v, want %#v", test.got, test.want)
			}
		})
	}
}

func TestIntegerRecordsAreComparable(t *testing.T) {
	byteKeys := map[Byte4]bool{NewByte4(1, 2, 3, 4): true}
	intKeys := map[IVec4]bool{IV4(1, 2, 3, 4): true}
	if !byteKeys[NewByte4(1, 2, 3, 4)] || !intKeys[IV4(1, 2, 3, 4)] {
		t.Fatal("equal values did not retrieve map entries")
	}
}
