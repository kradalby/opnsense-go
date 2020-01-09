package main

import (
	"log"

	"github.com/kradalby/opnsense-go/opnsense"
	"github.com/satori/go.uuid"
)

func main() {
	// c, err := opnsense.NewClient("http://localhost:8080", "", "", true)
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

	// b, err := c.Backup()
	// if err != nil {
	// 	log.Printf("Error: %#v", err)
	// }
	// log.Printf("%s", b)

	// clients, err := c.BgpNeighborList()
	// if err != nil {
	// 	log.Printf("Error: %#v", err)
	// }

	u, _ := uuid.FromString("9cc6398c-26ca-4cae-8ff0-8bb97114bdd0")
	resp, err := c.WireGuardServerGet(u)
	log.Println(err)
	log.Printf("%#v", resp)

}
