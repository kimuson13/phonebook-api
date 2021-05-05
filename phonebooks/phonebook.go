package phonebooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Phonebook struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

var Pbooks []Phonebook

func Run() {
	Pbooks = append(Pbooks, Phonebook{ID: "1", Name: "kimu1", Phone: "09012345678"})
	Pbooks = append(Pbooks, Phonebook{ID: "2", Name: "kimu2", Phone: "08012345678"})

	r := mux.NewRouter()
	r.HandleFunc("/api/phonebooks", GetPhonebooksHandler).Methods("GET")
	r.HandleFunc("/api/phonebooks/{id}", GetPhonebookHandler).Methods("GET")
	r.HandleFunc("/api/phonebooks", CreateHandler).Methods("POST")
	r.HandleFunc("/api/phonebooks/{id}", UpdateHandler).Methods("PUT")
	r.HandleFunc("/api/phonebooks/{id}", DeleteHandler).Methods("DELETE")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func GetPhonebooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for _, pbook := range Pbooks {
		if err := enc.Encode(pbook); err != nil {
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
	for _, pbook := range Pbooks {
		if pbook.ID == params["id"] {
			if err := enc.Encode(pbook); err != nil {
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

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var p Phonebook
	var buf bytes.Buffer
	_ = json.NewDecoder(r.Body).Decode(&p)
	p.ID = strconv.Itoa(len(Pbooks) + 1)
	Pbooks = append(Pbooks, p)
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(p); err != nil {
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

func UpdateHandler(w http.ResponseWriter, r *http.Request) {

}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {

}
