package main_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jaxsax/projects/learning/go-hex/internal"
)

func TestUpperRoute(t *testing.T) {
	app := internal.NewApplication()
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
