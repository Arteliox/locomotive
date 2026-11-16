package util

import (
	"errors"
	"net/url"
)

// SafeURL returns the URL as scheme, host and path only. The user, password, query and
// fragment are dropped, since any of them can carry a credential.
func SafeURL(u url.URL) string {
	return (&url.URL{Scheme: u.Scheme, Host: u.Host, Path: u.Path, RawPath: u.RawPath}).String()
}

// RedactURLError rewrites the URL quoted by a *url.Error inside err to SafeURL(u), so the
// error can be logged without the credentials the URL carries. Other errors pass through.
func RedactURLError(err error, u url.URL) error {
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		urlErr.URL = SafeURL(u)
	}

	return err
}
