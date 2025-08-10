package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerHelloEndpoint(t *testing.T) {
	ts := httptest.NewServer(New())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/hello")
	if err != nil {
		t.Fatalf("failed to get hello endpoint: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.JSONEq(t, `{"message":"hello world"}`, string(body))
}
