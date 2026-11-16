package util

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// closedAddress returns a local address nothing listens on.
func closedAddress(t *testing.T) string {
	t.Helper()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	addr := l.Addr().String()
	l.Close()

	return addr
}

func TestRedactURLErrorDropsCredentialsFromARequestError(t *testing.T) {
	addr := closedAddress(t)

	for _, raw := range []string{
		"http://FAKEUSERaa11:FAKEPASSbb22@" + addr + "/loki/api/v1/push",
		"http://FAKEUSERaa11@" + addr + "/loki/api/v1/push",
		"http://" + addr + "/loki/api/v1/push?api_key=FAKEKEYcc33#FAKEFRAGdd44",
	} {
		u, err := url.Parse(raw)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(nil))
		if err != nil {
			t.Fatal(err)
		}

		_, err = http.DefaultClient.Do(req)
		if err == nil {
			t.Fatal("expected the request to a closed port to fail")
		}

		msg := fmt.Errorf("failed to send webhook request: %w", RedactURLError(err, *u)).Error()

		for _, secret := range []string{"FAKEUSER", "FAKEPASS", "FAKEKEY", "FAKEFRAG", "***"} {
			if strings.Contains(msg, secret) {
				t.Errorf("error %q contains %q", msg, secret)
			}
		}

		if want := "http://" + addr + "/loki/api/v1/push"; !strings.Contains(msg, want) {
			t.Errorf("error %q does not name %q", msg, want)
		}
	}
}

func TestRedactURLErrorLeavesOtherErrorsAlone(t *testing.T) {
	err := fmt.Errorf("non success status code: %d", 500)
	if got := RedactURLError(err, url.URL{}); got != err {
		t.Errorf("RedactURLError changed an error with no URL in it")
	}
}
