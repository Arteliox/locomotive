package railway

import (
	"errors"
	"strings"
	"testing"
)

func TestParseTokenAcceptsBothTokenFormats(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want string
	}{
		{"legacy uuid token", "0b3f6a52-7c1e-4d7a-9a51-2f4c8e9d1a6b", "0b3f6a52-7c1e-4d7a-9a51-2f4c8e9d1a6b"},
		{"uppercase uuid token is lowercased as before", "0B3F6A52-7C1E-4D7A-9A51-2F4C8E9D1A6B", "0b3f6a52-7c1e-4d7a-9a51-2f4c8e9d1a6b"},
		{"account token", "rlwy_acct_3kQ9xV7mN2pL8rT4wY6zB1cD5fG0hJ", "rlwy_acct_3kQ9xV7mN2pL8rT4wY6zB1cD5fG0hJ"},
		{"workspace token", "rlwy_wksp_8nR2tY5uI9oP3aS7dF1gH4jK6lZ0xC", "rlwy_wksp_8nR2tY5uI9oP3aS7dF1gH4jK6lZ0xC"},
		{"surrounding whitespace is trimmed", "  rlwy_acct_abc123\n", "rlwy_acct_abc123"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseToken(tc.raw)
			if err != nil {
				t.Fatalf("ParseToken returned error: %v", err)
			}

			if got != tc.want {
				t.Errorf("ParseToken = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParseTokenRejectsABadValueWithoutPrintingIt(t *testing.T) {
	cases := []struct {
		name string
		raw  string
	}{
		{"empty", ""},
		{"blank", "   \n\t"},
		{"inner space", "rlwy_acct_secretpart1 secretpart2"},
		{"control character", "rlwy_acct_secretpart1\x00secretpart2"},
		{"inner newline", "rlwy_acct_secretpart1\nsecretpart2"},
		{"all-zero uuid", "00000000-0000-0000-0000-000000000000"},
		{"non ascii", "rlwy_acct_secretpart1ésecretpart2"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseToken(tc.raw)
			if !errors.Is(err, ErrInvalidToken) {
				t.Fatalf("ParseToken error = %v, want ErrInvalidToken", err)
			}

			if got != "" {
				t.Errorf("ParseToken returned a value alongside its error")
			}

			if msg := err.Error(); msg != "RAILWAY_API_KEY is not a valid token" {
				t.Errorf("error message = %q, want exactly the fixed sentence", msg)
			}

			if strings.Contains(err.Error(), "secretpart") {
				t.Errorf("error message contains part of the configured value")
			}
		})
	}
}
