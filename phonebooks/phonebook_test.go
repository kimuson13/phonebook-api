package phonebooks_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kimuson13/phonebook-api/phonebooks"
)

func TestGetPhonebooksHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	phonebooks.GetPhonebooksHandler(w, r)
	rw := w.Result()
	defer rw.Body.Close()
	if rw.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", rw.StatusCode)
	}
}
