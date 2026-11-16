package railway

import (
	"errors"
	"strings"

	"github.com/flexstack/uuid"
)

// ErrInvalidToken is the only error ParseToken returns. It never carries the configured
// value, so a bad token can be reported without printing it.
var ErrInvalidToken = errors.New("RAILWAY_API_KEY is not a valid token")

// ParseToken checks and normalises a Railway API token read from configuration. Railway
// has issued tokens as UUIDs and, more recently, as prefixed strings (rlwy_acct_...,
// rlwy_wksp_...), so the token is treated as an opaque string. Surrounding whitespace is
// trimmed. A blank or all-zero token, or one holding a space, a control character or a
// non-ASCII character, is rejected. A UUID token is written in its canonical lowercase
// form, exactly as it was when the token was parsed as a UUID, so an existing
// configuration sends the same header as before.
func ParseToken(raw string) (string, error) {
	token := strings.TrimSpace(raw)
	if token == "" {
		return "", ErrInvalidToken
	}

	for i := 0; i < len(token); i++ {
		if c := token[i]; c <= ' ' || c > '~' {
			return "", ErrInvalidToken
		}
	}

	if u, err := uuid.FromString(token); err == nil {
		// The all-zero UUID was refused before as an empty token; keep refusing it.
		if u == uuid.Nil {
			return "", ErrInvalidToken
		}

		return u.String(), nil
	}

	return token, nil
}
