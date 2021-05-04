package main

import (
	"log"
	"net/http"
)

type Phonebook struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/api/phonebooks", GetPhonebooksHandler)
	r.HandleFunc("/api/phoneboooks/{id}", GetPhonebookHandler)
	r.HandleFunc("/api/phonebooks/create", CreatePhonebookHandler)
	r.HandleFunc("/api/phonebooks/update/{id}", UpdatePhonebookHandler)
	r.HandleFunc("/api/phonebooks/delete/{id}", DeletePhonebookHandler)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func GetPhonebooksHandler(w http.ResponseWriter, r *http.Request) {

}

func GetPhonebookHandler(w http.ResponseWriter, r *http.Request) {

}

func CreatePhonebookHandler(w http.ResponseWriter, r *http.Request) {

}

func UpdatePhonebookHandler(w http.ResponseWriter, r *http.Request) {

}

func DeletePhonebookHandler(w http.ResponseWriter, r *http.Request) {

}
