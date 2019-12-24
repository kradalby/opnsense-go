package opnsense

import (
	"fmt"
	"log"
	"path"

	uuid "github.com/satori/go.uuid"
)

// Requires: os-frr

type BGP struct {
	Enabled   string    `json:"enabled"`
	Asnumber  string    `json:"asnumber"`
	Routerid  string    `json:"routerid"`
	Neighbors Neighbors `json:"neighbors"`

	// TODO: Fully implement these attributes
	Networks     []interface{} `json:"networks"`
	Redistribute Redistribute  `json:"redistribute"`
	Aspaths      Aspaths       `json:"aspaths"`
	Prefixlists  Prefixlists   `json:"prefixlists"`
	Routemaps    Routemaps     `json:"routemaps"`
}

type Aspaths struct {
	Aspath []interface{} `json:"aspath"`
}

type Neighbors struct {
	Neighbor map[string]BgpNeighborGet `json:"neighbor"`
}

type BgpNeighborBase struct {
	UUID             *uuid.UUID `json:"uuid,omitempty"`
	Enabled          string     `json:"enabled"`
	Address          string     `json:"address"`
	Remoteas         string     `json:"remoteas"`
	Nexthopself      string     `json:"nexthopself"`
	Defaultoriginate string     `json:"defaultoriginate"`
}

type BgpNeighborGet struct {
	BgpNeighborBase
	Updatesource        SelectedMap `json:"updatesource"`
	LinkedPrefixlistIn  SelectedMap `json:"linkedPrefixlistIn"`
	LinkedPrefixlistOut SelectedMap `json:"linkedPrefixlistOut"`
	LinkedRoutemapIn    SelectedMap `json:"linkedRoutemapIn"`
	LinkedRoutemapOut   SelectedMap `json:"linkedRoutemapOut"`
}

type BgpNeighborSet struct {
	BgpNeighborBase
	Updatesource        string `json:"updatesource"`
	LinkedPrefixlistIn  string `json:"linkedPrefixlistIn"`
	LinkedPrefixlistOut string `json:"linkedPrefixlistOut"`
	LinkedRoutemapIn    string `json:"linkedRoutemapIn"`
	LinkedRoutemapOut   string `json:"linkedRoutemapOut"`
}

type Prefixlists struct {
	Prefixlist []interface{} `json:"prefixlist"`
}

type Redistribute struct {
	OSPF      SelectedMap `json:"ospf"`
	Connected SelectedMap `json:"connected"`
	Kernel    SelectedMap `json:"kernel"`
	RIP       SelectedMap `json:"rip"`
	Static    SelectedMap `json:"static"`
}

type Routemaps struct {
	Routemap []interface{} `json:"routemap"`
}

func (c *Client) BgpNeighborGet(uuid uuid.UUID) (*BgpNeighborGet, error) {
	api := path.Join("quagga/bgp/getNeighbor", uuid.String())

	type Response struct {
		Neighbor BgpNeighborGet `json:"neighbor"`
	}
	var response Response

	err := c.GetAndUnmarshal(api, &response)

	// UUID does not exist in the JSON, so we add it since we know it.
	response.Neighbor.UUID = &uuid

	return &response.Neighbor, err
}

func (c *Client) BgpNeighborGetUUIDs() ([]*uuid.UUID, error) {
	api := "quagga/bgp/searchNeighbor"

	var response SearchResult
	err := c.GetAndUnmarshal(api, &response)
	if err != nil {
		return nil, err
	}

	uuids := []*uuid.UUID{}
	for _, row := range response.Rows {
		m := row.(map[string]interface{})
		uuid, err := uuid.FromString(m["uuid"].(string))
		if err == nil {
			uuids = append(uuids, &uuid)
		}
	}

	return uuids, err
}

func (c *Client) BgpNeighborList() ([]*BgpNeighborGet, error) {
	uuids, err := c.BgpNeighborGetUUIDs()
	if err != nil {
		return nil, err
	}

	clients := []*BgpNeighborGet{}
	for _, uuid := range uuids {
		client, err := c.BgpNeighborGet(*uuid)
		if err == nil {
			clients = append(clients, client)
		}
	}
	return clients, nil
}
func (c *Client) BgpNeighborSet(uuid uuid.UUID, clientConf BgpNeighborSet) (*GenericResponse, error) {
	api := path.Join("quagga/bgp/setNeighbor", uuid.String())

	request := map[string]interface{}{
		"neighbor": clientConf,
	}

	var response GenericResponse
	err := c.PostAndMarshal(api, request, &response)
	if err != nil {
		return nil, err
	}

	if response.Result != "saved" {
		err := fmt.Errorf("Failed to save, response from server: %#v", response)
		log.Printf("[ERROR] %#v\n", err)
		return nil, err
	}

	return &response, nil
}

func (c *Client) BgpNeighborAdd(clientConf BgpNeighborSet) (*uuid.UUID, error) {
	api := "quagga/bgp/addNeighbor"

	request := map[string]interface{}{
		"neighbor": clientConf,
	}

	var response GenericResponse
	err := c.PostAndMarshal(api, request, &response)
	if err != nil {
		return nil, err
	}

	if response.Result != "saved" {
		err := fmt.Errorf("Failed to save, response from server: %#v", response)
		log.Printf("[ERROR] %#v\n", err)
		return nil, err
	}

	return response.UUID, nil
}

func (c *Client) BgpNeighborDelete(uuid uuid.UUID) (*GenericResponse, error) {
	api := path.Join("quagga/bgp/delNeighbor", uuid.String())

	var response GenericResponse
	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	if response.Result != "deleted" {
		err := fmt.Errorf("Failed to delete, response from server: %#v", response)
		log.Printf("[ERROR] %#v\n", err)
		return nil, err
	}

	return &response, nil
}
