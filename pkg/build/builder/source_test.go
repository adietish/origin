package builder

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCheckRemoteGit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	var err error
	err = checkRemoteGit(server.URL, 10*time.Second)
	switch v := err.(type) {
	case gitAuthError:
	default:
		t.Errorf("expected gitAuthError, got %q", v)
	}

	t0 := time.Now()
	err = checkRemoteGit("https://254.254.254.254/foo/bar", 4*time.Second)
	t1 := time.Now()
	if err == nil || (err != nil && !strings.Contains(fmt.Sprintf("%s", err), "timeout")) {
		t.Errorf("expected timeout error, got %q", err)
	}
	if t1.Sub(t0) > 5*time.Second {
		t.Errorf("expected timeout in 4 seconds, it took %v", t1.Sub(t0))
	}

	err = checkRemoteGit("https://github.com/openshift/origin", 10*time.Second)
	if err != nil {
		t.Errorf("unexpected error %q", err)
	}
}
