package main

import (
	"github.com/kradalby/opnsense-go/opnsense"
	"log"
)

func main() {
	c, err := opnsense.NewClient("https://172.16.207.143/api/", "", "", true)
	if err != nil {
		log.Fatal(err)
	}

	r, err := c.GetInformation()
	log.Printf("%#v,", r)
	log.Printf("Error: %#v", err)

}
