package bullet_stl

import "testing"

func TestBulletIdIntToaasciAndDecode(t *testing.T) {
	tests := []struct {
		intValue int64
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{35, "z"},
		{36, "10"},
		{12345, "9ix"},
		{1679615, "zzzz"}, // 36^4 - 1, max value for 4-char base36
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			encoded, err := BulletIdIntToaasci(tt.intValue)
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
	_, err := BulletIdIntToaasci(-1)
	if err == nil {
		t.Error("expected error for negative value, got nil")
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

	for i := int64(0); i < 1000; i++ {
		encoded, err := BulletIdIntToaasci(i)
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
