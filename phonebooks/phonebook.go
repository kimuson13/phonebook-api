package phonebooks

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
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

func ConnectSQL() (*sql.DB, func(), error) {
	s := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_MODE"))
	db, err := sql.Open("postgres", s)
	if err != nil {
		return nil, nil, err
	}

	const sql = `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(12) NOT NULL,
		phone VARCHAR(12) NOT NULL
	);
	`

	if _, err := db.Exec(sql); err != nil {
		return nil, nil, err
	}

	return db, func() {
		if err := db.Close(); err != nil {
			log.Print("can not close database")
		}
	}, nil
}

func AddRecords(db *sql.DB, p Phonebook) (int, error) {
	var id int
	err := db.QueryRow("INSERT INTO users(name, phone) values ($1, $2) RETURNING id", p.Name, p.Phone).Scan(&id)
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

	db, closefunc, err := ConnectSQL()
	if err != nil {
		log.Print(err)
		http.Error(w, openFail, status500)
	}
	defer closefunc()
	rows, err := db.Query("SELECT * FROM users")
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

	db, closefunc, err := ConnectSQL()
	if err != nil {
		log.Print(err)
		http.Error(w, openFail, status500)
	}
	defer closefunc()

	rows, err := db.Query("SELECT * FROM users WHERE id = " + params["id"])
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

	db, closefunc, err := ConnectSQL()
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

	id, err := AddRecords(db, p)
	if err != nil {
		log.Print(err)
		http.Error(w, addFail, status500)
	}

	enc := json.NewEncoder(&buf)
	rows, err := db.Query("SELECT * FROM users WHERE id = " + strconv.Itoa(id))
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

	db, closefunc, err := ConnectSQL()
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

	if _, err := db.Exec("UPDATE users SET (name, phone) = ($1, $2) WHERE id = "+params["id"], p.Name, p.Phone); err != nil {
		log.Print(err)
		http.Error(w, updateFail, status500)
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	rows, err := db.Query("SELECT * FROM users WHERE id = " + params["id"])
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

	db, closefunc, err := ConnectSQL()
	if err != nil {
		log.Print(err)
		http.Error(w, openFail, status500)
	}
	defer closefunc()

	if _, err := db.Exec("DELETE FROM users WHERE id = " + params["id"]); err != nil {
		log.Print(err)
		http.Error(w, "can not delete the record", status500)
		return
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
