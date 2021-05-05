package phonebooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Phonebook struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

var books []Phonebook

func Run() {
	books = append(books, Phonebook{ID: "1", Name: "kimu1", Phone: "09012345678"})
	books = append(books, Phonebook{ID: "2", Name: "kimu2", Phone: "08012345678"})

	r := mux.NewRouter()
	r.HandleFunc("/api/phonebooks", GetPhonebooksHandler).Methods("GET")
	r.HandleFunc("/api/phonebooks/{id}", GetPhonebookHandler).Methods("GET")
	r.HandleFunc("/api/phonebooks/create", CreatePhonebookHandler).Methods("POST")
	r.HandleFunc("/api/phonebooks/update/{id}", UpdatePhonebookHandler).Methods("PUT")
	r.HandleFunc("/api/phonebooks/delete/{id}", DeletePhonebookHandler).Methods("DELETE")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func GetPhonebooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for _, book := range books {
		if err := enc.Encode(book); err != nil {
			log.Print(err)
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	}

	str := buf.String()
	_, err := fmt.Fprint(w, str)
	if err != nil {
		log.Fatal(err)
	}
}

func GetPhonebookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for _, book := range books {
		if book.ID == params["id"] {
			if err := enc.Encode(book); err != nil {
				log.Print(err)
				http.Error(w, "encode error", http.StatusInternalServerError)
			}

			str := buf.String()
			_, err := fmt.Fprint(w, str)
			if err != nil {
				log.Print(err)
				http.Error(w, "print error", http.StatusInternalServerError)
			}
			return
		}
	}

	if err := enc.Encode(&Phonebook{}); err != nil {
		log.Print(err)
		http.Error(w, "encode error", http.StatusInternalServerError)
	}

	str := buf.String()
	_, err := fmt.Fprint(w, str)
	if err != nil {
		log.Print(err)
		http.Error(w, "print error", http.StatusInternalServerError)
	}
}

func CreatePhonebookHandler(w http.ResponseWriter, r *http.Request) {

}

func UpdatePhonebookHandler(w http.ResponseWriter, r *http.Request) {

}

func DeletePhonebookHandler(w http.ResponseWriter, r *http.Request) {

}
