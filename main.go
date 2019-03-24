package main

import (
	"github.com/kradalby/opnsense-go/opnsense"
	"log"
)

func main() {
	// c, err := opnsense.NewClient("http://127.0.0.1:8080/api/", "", "", true)
	c, err := opnsense.NewClient("https://172.16.207.143/api", "", "", true)
	if err != nil {
		log.Fatal(err)
	}

	// resp, err := c.WireGuardGetSettings()
	// log.Printf("Error: %#v", err)
	// log.Printf("%#v,", resp)
	// err = c.WireGuardEnableService()
	// log.Printf("Error: %#v", err)
	// resp, err = c.WireGuardGetSettings()
	// log.Printf("Error: %#v", err)
	// log.Printf("%#v,", resp)
	// err = c.WireGuardDisableService()
	// log.Printf("Error: %#v", err)
	// resp, err = c.WireGuardGetSettings()
	// log.Printf("%#v,", resp)
	// log.Printf("Error: %#v", err)
	// err = c.WireGuardEnableService()
	// log.Printf("Error: %#v", err)

	b, err := c.Backup()
	if err != nil {
		log.Printf("Error: %#v", err)
	}
	log.Printf("%s", b)

}
