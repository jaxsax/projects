package main_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jaxsax/projects/learning/go-hex/internal"
	"github.com/stretchr/testify/require"
)

func TestUpperRoute(t *testing.T) {
	registry, err := internal.NewRegistry(&internal.Config{
		SQLPath: ":memory:",
	})
	if err != nil {
		t.Fatalf("failed to initialize registry: %v", err)
	}

	app := internal.NewApplication(registry)
	ts := httptest.NewServer(app.Router)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/upper", nil)
	if err != nil {
		t.Fatalf("unexpected err %v", err)
	}

	q := req.URL.Query()
	q.Add("q", "hey")
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected err %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %v", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unexpected err %v", err)
	}

	got := string(respBody)
	want := "HEY"

	if want != got {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRetrievePets(t *testing.T) {
	registry, err := internal.NewRegistry(&internal.Config{
		SQLPath: ":memory:",
	})
	require.NoError(t, err)
	require.NoError(t, registry.Setup())

	app := internal.NewApplication(registry)
	ts := httptest.NewServer(app.Router)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/pets", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode, string(respBody))

	type Pet struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
	}
	type petsResponse struct {
		Pets []Pet `json:"pets"`
	}

	var marshalledResponse petsResponse
	if err := json.Unmarshal(respBody, &marshalledResponse); err != nil {
		t.Fatalf("unexpected err %v", err)
	}

	want := petsResponse{
		Pets: []Pet{
			{
				ID:   1,
				Name: "Jonnay",
			},
		},
	}

	require.Equal(t, want, marshalledResponse)
}
