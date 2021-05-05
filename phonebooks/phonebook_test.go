package phonebooks_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kimuson13/phonebook-api/phonebooks"
)

var normalurl = "/api/phonebooks"
var idurl = "/api/phonebooks/1"

func TestAllHandlers_RespIsOk(t *testing.T) {
	cases := map[string]struct {
		method  string
		url     string
		handler func(w http.ResponseWriter, r *http.Request)
	}{
		"getall": {
			method:  "GET",
			url:     normalurl,
			handler: phonebooks.GetPhonebooksHandler,
		},
		"getid": {
			method:  "GET",
			url:     idurl,
			handler: phonebooks.GetPhonebookHandler,
		},
		"create": {
			method:  "POST",
			url:     normalurl,
			handler: phonebooks.CreateHandler,
		},
		"update": {
			method:  "PUT",
			url:     idurl,
			handler: phonebooks.UpdateHandler,
		},
		"delete": {
			method:  "DELETE",
			url:     idurl,
			handler: phonebooks.DeleteHandler,
		},
	}
	for n, c := range cases {
		c := c
		t.Run(n, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(c.method, c.url, nil)
			c.handler(w, r)
			rw := w.Result()
			defer rw.Body.Close()
			if rw.StatusCode != http.StatusOK {
				t.Fatalf("unexpected status code: %d", rw.StatusCode)
			}
		})
	}
}
