package phonebooks_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kimuson13/phonebook-api/phonebooks"
	_ "github.com/lib/pq"
)

func TestGetPhonebooksHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/phonebooks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(phonebooks.GetPhonebooksHandler)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestGetPhonebookHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/phonebooks/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/phonebooks/{id}", phonebooks.GetPhonebookHandler)
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("handler should have failed on routeVariavle 1: got %v want %v", rr.Code, http.StatusOK)
	}
}

func TestCreateHandler(t *testing.T) {
	_, closefunc := createDBForTest(t)
	defer closefunc()

	p := phonebooks.Phonebook{Name: "test1", Phone: "09012334567"}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(p); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/phonebooks", strings.NewReader(buf.String()))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(phonebooks.CreateHandler)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestUpdateHandler(t *testing.T) {
	_, closefunc := createDBForTest(t)
	defer closefunc()

	p := phonebooks.Phonebook{Name: "test1", Phone: "09012334567"}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(p); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/api/phonebooks/1", strings.NewReader(buf.String()))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/phonebooks/{id}", phonebooks.UpdateHandler)
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("handler should have failed on routeVariavle 1: got %v want %v", rr.Code, http.StatusOK)
	}
}

func TestDeleteHandler(t *testing.T) {
	_, closefunc := createDBForTest(t)
	defer closefunc()

	req, err := http.NewRequest("DELETE", "/api/phonebooks/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/phonebooks/{id}", phonebooks.DeleteHandler)
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("handler should have failed on routeVariavle 1: got %v want %v", rr.Code, http.StatusOK)
	}
}

func createDBForTest(t *testing.T) (*sql.DB, func()) {
	t.Helper()
	db, closefunc, err := phonebooks.ConnectSQL()
	if err != nil {
		t.Fatal(err)
	}

	testRecoed := phonebooks.Phonebook{Name: "test", Phone: "09087878787"}
	if _, err := phonebooks.AddRecords(db, testRecoed); err != nil {
		t.Fatal(err)
	}
	return db, closefunc
}
