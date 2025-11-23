package bullet_stl

import (
	"testing"
)

func TestMakeAndParse(t *testing.T) {
	tests := []struct {
		bucket int32
		local  int32
	}{
		{0, 0},
		{1, 42},
		{123, 999},
		{2147483647, 2147483647}, // max int32
		{-1, -1},                 // negative values also preserved
		{100, -50},
	}

	for _, tt := range tests {
		id := MakeNamespacedId(tt.bucket, tt.local)
		spaced := ParseNamespacedId(id)

		if spaced.Namespace != tt.bucket || spaced.Id != tt.local {
			t.Fatalf("roundtrip failed: input (%d,%d) got (%d,%d)",
				tt.bucket, tt.local, spaced.Namespace, spaced.Id)
		}
	}
}

func TestLocalOverflow(t *testing.T) {
	_, err := NextLocalID(0x7fffffff)
	if err == nil {
		t.Fatalf("expected overflow error")
	}
}

func FuzzMakeParse(f *testing.F) {
	seeds := []int64{
		0, 1, -1, 123456, 1 << 40, -1 << 40,
	}
	for _, s := range seeds {
		f.Add(int32(s), int32(s>>1))
	}

	f.Fuzz(func(t *testing.T, b int32, l int32) {
		spaced := MakeNamespacedId(b, l)
		parsed := ParseNamespacedId(spaced)
		if parsed.Namespace != b || parsed.Id != l {
			t.Fatalf("mismatch after roundtrip: (%d,%d) -> (%d,%d)", b, l, parsed.Namespace, parsed.Id)
		}
	})
}
