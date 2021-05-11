package phonebooks_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kimuson13/phonebook-api/phonebooks"
	_ "github.com/lib/pq"
)

func TestGetPhonebooksHandler(t *testing.T) {
	db, id, closefunc := setDBForTest(t)
	defer closefunc()

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

	if err := deleteRecord(db, id, t); err != nil {
		t.Fatal(err)
	}
}

func TestGetPhonebookHandler(t *testing.T) {
	db, id, closefunc := setDBForTest(t)
	defer closefunc()

	strID := strconv.Itoa(id)
	req, err := http.NewRequest("GET", "/api/phonebooks/"+strID, nil)
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

	if err := deleteRecord(db, id, t); err != nil {
		t.Fatal(err)
	}
}

func TestCreateHandler(t *testing.T) {
	db, id, closefunc := setDBForTest(t)
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

	var createID int
	if err := db.QueryRow("SELECT id FROM users WHERE name = $1 AND phone = $2", p.Name, p.Phone).Scan(&createID); err != nil {
		t.Fatal(err)
	}

	if err := deleteRecord(db, id, t); err != nil {
		t.Fatal(err)
	}

	if err := deleteRecord(db, createID, t); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateHandler(t *testing.T) {
	db, id, closefunc := setDBForTest(t)
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

	if err := deleteRecord(db, id, t); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteHandler(t *testing.T) {
	_, id, closefunc := setDBForTest(t)
	defer closefunc()

	strID := strconv.Itoa(id)
	req, err := http.NewRequest("DELETE", "/api/phonebooks/"+strID, nil)
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

func setDBForTest(t *testing.T) (*sql.DB, int, func()) {
	t.Helper()
	if err := godotenv.Load("_testdata/test.env"); err != nil {
		t.Fatal(err)
	}
	db, closefunc, err := phonebooks.ConnectSQL()
	if err != nil {
		t.Fatal(err)
	}

	testRecoed := phonebooks.Phonebook{Name: "test", Phone: "09087878787"}
	id, err := phonebooks.AddRecords(db, testRecoed)
	if err != nil {
		t.Fatal(err)
	}
	return db, id, closefunc
}

func deleteRecord(db *sql.DB, id int, t *testing.T) error {
	t.Helper()
	if _, err := db.Exec("DELETE FROM users WHERE id = $1", id); err != nil {
		return err
	}

	return nil
}
