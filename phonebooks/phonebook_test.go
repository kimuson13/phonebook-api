package phonebooks_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kimuson13/phonebook-api/phonebooks"
)

func TestGetPhonebooksHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/phonebooks", nil)
	phonebooks.GetPhonebooksHandler(w, r)
	rw := w.Result()
	defer rw.Body.Close()
	if rw.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", rw.StatusCode)
	}
}

func TestGetPhonebookHandler(t *testing.T) {
	cases := map[string]string{
		"id=1": "1",
		"id=2": "2",
	}
	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			c := c
			path := fmt.Sprintf("/api/phonebooks/%s", c)
			r := httptest.NewRequest("GET", path, nil)
			w := httptest.NewRecorder()
			phonebooks.GetPhonebookHandler(w, r)
			rw := w.Result()
			defer rw.Body.Close()
			if rw.StatusCode != http.StatusOK {
				t.Fatalf("unexpected status code: %d", rw.StatusCode)
			}
		})
	}
}

func TestCreateHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/phonebooks", nil)
	phonebooks.CreateHandler(w, r)
	rw := w.Result()
	defer rw.Body.Close()
	if rw.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", rw.StatusCode)
	}
}
