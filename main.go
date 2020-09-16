// +build !codeanalysis

package main

import (
	"log"

	"github.com/kradalby/opnsense-go/opnsense"
)

func main() {
	_, err := opnsense.NewClient("http://localhost:8080", "", "", true)
	if err != nil {
		log.Fatal(err)
	}
}

