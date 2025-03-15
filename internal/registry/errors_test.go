package registry

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func MockRegistryError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(400)
	w.Header().Set("Content-Type", "application/json")

	e := RegistryError{Errors: []RegistryErrorDetail{{
		"Error",
		"This is an error",
		"Something terrible happened",
	}}}

	b, _ := json.Marshal(e)
	w.Write(b)
}

func TestRegistryError(t *testing.T) {
	rec := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/err", nil)
	MockRegistryError(rec, req)

	var gotErr RegistryError
	_ = json.Unmarshal(rec.Body.Bytes(), &gotErr)

	got := gotErr.Errors[0].Detail
	want := "Something terrible happened"

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
