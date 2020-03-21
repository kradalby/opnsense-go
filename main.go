package main

import (
	"github.com/kradalby/opnsense-go/opnsense"

	"log"
)

func main() {
	_, err := opnsense.NewClient("http://localhost:8080", "", "", true)

	if err != nil {
		log.Fatal(err)
	}
}
