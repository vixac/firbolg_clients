package bullet_stl

import "testing"

func TestBulletIdIntToaasciAndDecode(t *testing.T) {
	tests := []struct {
		intValue int64
		idSize   int
		expected string
	}{
		{0, 4, "0000"},
		{1, 4, "0001"},
		{35, 4, "000z"},
		{36, 4, "0010"},
		{12345, 4, "09ix"},
		{1679615, 4, "zzzz"}, // 36^4 - 1, max value for 4-char base36
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			encoded, err := BulletIdIntToaasci(tt.intValue, tt.idSize)
			if err != nil {
				t.Fatalf("unexpected error encoding: %v", err)
			}
			if encoded != tt.expected {
				t.Errorf("expected encoded %s, got %s", tt.expected, encoded)
			}

			decoded, err := AasciBulletIdToInt(encoded)
			if err != nil {
				t.Fatalf("unexpected error decoding: %v", err)
			}
			if decoded != tt.intValue {
				t.Errorf("round-trip failed: expected %d, got %d", tt.intValue, decoded)
			}
		})
	}
}

func TestBulletIdIntToaasci_Invalid(t *testing.T) {
	_, err := BulletIdIntToaasci(-1, 4)
	if err == nil {
		t.Error("expected error for negative value, got nil")
	}

	_, err = BulletIdIntToaasci(1679616, 4) // 36^4 = 1679616, out of range
	if err == nil {
		t.Error("expected error for value too large, got nil")
	}
}

func TestAasciBulletIdToInt_Invalid(t *testing.T) {
	invalidStrings := []string{"!000", "123$", "abcd*", "12 3", "ABCD"}

	for _, s := range invalidStrings {
		_, err := AasciBulletIdToInt(s)
		if err == nil {
			t.Errorf("expected error for invalid string %q, got nil", s)
		}
	}
}

func TestBulletIdIntToaasci_RoundTrip(t *testing.T) {
	// Test a range of values for round-trip correctness
	idSize := 4
	for i := int64(0); i < 1000; i++ {
		encoded, err := BulletIdIntToaasci(i, idSize)
		if err != nil {
			t.Fatalf("unexpected error encoding %d: %v", i, err)
		}

		decoded, err := AasciBulletIdToInt(encoded)
		if err != nil {
			t.Fatalf("unexpected error decoding %s: %v", encoded, err)
		}

		if decoded != i {
			t.Errorf("round-trip mismatch: %d -> %s -> %d", i, encoded, decoded)
		}
	}
}
