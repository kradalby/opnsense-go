package main

import (
	"log"

	"github.com/kradalby/opnsense-go/opnsense"
)

func main() {
	c, err := opnsense.NewClient("https://10.61.0.35", "nY+Dh3TW6JUNBbiNcaoaVYU81MX0HW85hhkxiOl2Ehg2KlTySClV4J/wymk5XiSuNnw+IWoTc6hzCvcM", "vZdE0GnMAV23EPOvf0IzV86gzMmiHsGhOpvnOy7NIS8xy1v3qz//AI9I6edd3opJbozNrgsDFnp5LYjj", true)
	if err != nil {
		log.Fatal(err)
	}

	sp, err := c.FirewallFilterSavepoint()
	if err != nil {
		log.Fatalf("Error: %#v", err)
	}

	log.Printf("Resp: %#v", sp)

	// resp, err := c.FirewallFilterRuleSearch()
	// if err != nil {
	// 	log.Fatalf("Error: %#v", err)
	// }

	// log.Printf("Resp: %#v", resp)

	// for _, rule := range resp {
	// 	c.FirewallFilterRuleDelete(*rule.UUID)
	// }

	rule := opnsense.FilterRule{
		// Enabled:        true,
		Description:    "Test rule",
		SourceNet:      "10.100.0.0/24",
		Protocol:       "TCP",
		DestinationNet: "10.101.0.0/24",
		Interface:      "wan",
	}

	err = c.FirewallFilterRuleAdd(rule)
	if err != nil {
		log.Fatalf("Error: %#v", err)
	}

	resp, err := c.FirewallFilterApply(sp.Revision)
	if err != nil {
		log.Fatalf("Error: %#v", err)
	}

	log.Printf("Resp: %#v", resp)

	// resp4, err := c.FirewallFilterCancelRollback(sp.Revision)
	// if err != nil {
	// 	log.Fatalf("Error: %#v", err)
	// }

	// log.Printf("Resp: %#v", resp4)

	resp2, err := c.FirewallFilterRuleSearch()
	if err != nil {
		log.Fatalf("Error: %#v", err)
	}

	log.Printf("Resp: %#v", resp2)

	// resp2, err := c.FirewallSourceNatRuleSearch()
	// if err != nil {
	// 	log.Fatalf("Error: %#v", err)
	// }

	// log.Printf("Resp2: %#v", resp2)
}
