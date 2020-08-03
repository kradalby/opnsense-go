package opnsense

import (
	"io/ioutil"
	"log"
)

// Requires: os-api-backup.
func (c *Client) Backup() (string, error) {
	api := "backup/backup/download"

	resp, err := c.Get(api)
	if err != nil {
		log.Printf("[ERROR] Failed to download backup: %#v", err)

		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read GET response: %#v\n", err)

		return "", err
	}

	return string(body), nil
}
