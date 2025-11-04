package phone

import "testing"

func TestNormalize(t *testing.T) {
	cases := map[string]string{
		"1234567890":     "1234567890",
		"123 456 7891":   "1234567891",
		"(123) 456 7892": "1234567892",
		"(123) 456-7893": "1234567893",
		"123-456-7894":   "1234567894",
	}
	for in, want := range cases {
		if got := Normalize(in); got != want {
			t.Fatalf("Normalize(%q) = %q, want %q", in, got, want)
		}
	}
}
