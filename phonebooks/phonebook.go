package phonebooks

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "modernc.org/sqlite"
)

type Phonebook struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

var (
	addFail    string = "can not add record to body"
	openFail   string = "can not open the sql"
	selectFail string = "can not select the table"
	scanFail   string = "can not scan the record selected"
	encodeFail string = "can not encode this records"
	printFail  string = "can not print this records"
	decodeFail string = "can not decode body you typed"
	updateFail string = "can not update the record you selected"

	status500 int = http.StatusInternalServerError
)

func ConnectSQL(name string) (*sql.DB, func(), error) {
	db, err := sql.Open("sqlite", name)
	if err != nil {
		return nil, nil, err
	}

	const sql = `
	CREATE TABLE IF NOT EXISTS user (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		phone TEXT NOT NULL
	);
	`

	if _, err := db.Exec(sql); err != nil {
		return nil, nil, err
	}

	return db, func() {
		if err := db.Close(); err != nil {
			log.Printf("can not close %v", name)
		}
	}, nil
}

func AddRecords(db *sql.DB, p *Phonebook) (int, error) {
	const sql = "INSERT INTO user(name, phone) values (?,?)"
	r, err := db.Exec(sql, p.Name, p.Phone)
	if err != nil {
		return 0, err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func Run() {
	r := mux.NewRouter()
	r.HandleFunc("/api/phonebooks", GetPhonebooksHandler).Methods("GET")
	r.HandleFunc("/api/phonebooks/{id}", GetPhonebookHandler).Methods("GET")
	r.HandleFunc("/api/phonebooks", CreateHandler).Methods("POST")
	r.HandleFunc("/api/phonebooks/{id}", UpdateHandler).Methods("PUT")
	r.HandleFunc("/api/phonebooks/{id}", DeleteHandler).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func GetPhonebooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	db, closefunc, err := ConnectSQL("phonebook.db")
	if err != nil {
		log.Print(err)
		http.Error(w, openFail, status500)
	}
	defer closefunc()

	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
		log.Print(err)
		http.Error(w, selectFail, status500)
	}

	for rows.Next() {
		var p Phonebook
		if err := rows.Scan(&p.ID, &p.Name, &p.Phone); err != nil {
			log.Print(err)
			http.Error(w, scanFail, status500)
		}

		if err := enc.Encode(p); err != nil {
			log.Print(err)
			http.Error(w, encodeFail, status500)
		}
	}

	str := buf.String()
	if _, err := fmt.Fprint(w, str); err != nil {
		log.Print(err)
		http.Error(w, printFail, status500)
	}
}

func GetPhonebookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	db, closefunc, err := ConnectSQL("phonebook.db")
	if err != nil {
		log.Print(err)
		http.Error(w, openFail, status500)
	}
	defer closefunc()

	selectedID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Print(err)
		http.Error(w, encodeFail, status500)
	}

	rows, err := db.Query("SELECT * FROM user WHERE id = ?", selectedID)
	for rows.Next() {
		var p Phonebook
		if err := rows.Scan(&p.ID, &p.Name, &p.Phone); err != nil {
			log.Print(err)
			http.Error(w, scanFail, status500)
		}

		if err := enc.Encode(p); err != nil {
			log.Print(err)
			http.Error(w, encodeFail, status500)
		}
	}

	str := buf.String()
	if _, err := fmt.Fprint(w, str); err != nil {
		log.Print(err)
		http.Error(w, printFail, status500)
	}
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, closefunc, err := ConnectSQL("phonebook.db")
	if err != nil {
		log.Print(err)
		http.Error(w, openFail, status500)
	}
	defer closefunc()

	var (
		p   Phonebook
		buf bytes.Buffer
	)
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Print(err)
		http.Error(w, decodeFail, status500)
	}

	id, err := AddRecords(db, &p)
	if err != nil {
		log.Print(err)
		http.Error(w, addFail, status500)
	}

	enc := json.NewEncoder(&buf)
	rows, err := db.Query("SELECT * FROM user WHERE id = ?", id)
	for rows.Next() {
		if err := rows.Scan(&p.ID, &p.Name, &p.Phone); err != nil {
			log.Print(err)
			http.Error(w, scanFail, status500)
		}

		if err := enc.Encode(p); err != nil {
			log.Print(err)
			http.Error(w, encodeFail, status500)
		}
	}

	str := buf.String()
	if _, err := fmt.Fprint(w, str); err != nil {
		log.Print(err)
		http.Error(w, printFail, status500)
	}
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	db, closefunc, err := ConnectSQL("phonebook.db")
	if err != nil {
		log.Print(err)
		http.Error(w, openFail, status500)
	}
	defer closefunc()

	var p Phonebook
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Print(err)
		http.Error(w, decodeFail, status500)
	}

	selectedID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Print(err)
		http.Error(w, encodeFail, status500)
	}

	if _, err := db.Exec("UPDATE user SET (name, phone) = (?,?) WHERE id = ?", p.Name, p.Phone, selectedID); err != nil {
		log.Print(err)
		http.Error(w, updateFail, status500)
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	rows, err := db.Query("SELECT * FROM user WHERE id = ?", selectedID)
	for rows.Next() {
		var p Phonebook
		if err := rows.Scan(&p.ID, &p.Name, &p.Phone); err != nil {
			log.Print(err)
			http.Error(w, scanFail, status500)
		}

		if err := enc.Encode(p); err != nil {
			log.Print(err)
			http.Error(w, encodeFail, status500)
		}
	}

	str := buf.String()
	if _, err := fmt.Fprint(w, str); err != nil {
		log.Print(err)
		http.Error(w, printFail, status500)
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	db, closefunc, err := ConnectSQL("phonebook.db")
	if err != nil {
		log.Print(err)
		http.Error(w, openFail, status500)
	}
	defer closefunc()

	selectedID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Print(err)
		http.Error(w, encodeFail, status500)
	}

	if _, err := db.Exec("DELETE FROM user WHERE id = ?", selectedID); err != nil {
		log.Print(err)
		http.Error(w, "can not delete the record", status500)
	}

	var buf bytes.Buffer
	var p Phonebook
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(p); err != nil {
		log.Print(err)
		http.Error(w, encodeFail, status500)
	}

	str := buf.String()
	if _, err := fmt.Fprint(w, str); err != nil {
		log.Print(err)
		http.Error(w, printFail, status500)
	}
}
