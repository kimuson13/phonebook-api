package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kimuson13/phonebook-api/phonebooks"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	phonebooks.Run()
}
