package opnsense

import (
	"fmt"
	"log"
	"path"

	uuid "github.com/satori/go.uuid"
)

// Requires: os-wireguard-devel

func (c *Client) WireGuardRestart() (*GenericResponse, error) {
	api := "wireguard/service/restart"

	var response GenericResponse

	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) WireGuardStart() (*GenericResponse, error) {
	api := "wireguard/service/start"

	var response GenericResponse

	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) WireGuardStop() (*GenericResponse, error) {
	api := "wireguard/service/stop"

	var response GenericResponse

	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) WireGuardShowConfig() (*GenericResponse, error) {
	api := "wireguard/service/showconf"

	var response GenericResponse
	err := c.GetAndUnmarshal(api, &response)

	return &response, err
}

func (c *Client) WireGuardShowHandshake() (*GenericResponse, error) {
	api := "wireguard/service/showhandshake"

	var response GenericResponse
	err := c.GetAndUnmarshal(api, &response)

	return &response, err
}

type WireGuardSettings struct {
	General WireGuardSettingsGeneral `json:"general"`
}

type WireGuardSettingsGeneral struct {
	Enabled string `json:"enabled"`
}

func (c *Client) WireGuardSettingsGet() (*WireGuardSettings, error) {
	api := "wireguard/general/get"

	var response WireGuardSettings
	err := c.GetAndUnmarshal(api, &response)

	return &response, err
}

func (c *Client) WireGuardSettingsSet(settings WireGuardSettings) (*GenericResponse, error) {
	api := "wireguard/general/set"

	var response GenericResponse

	err := c.PostAndMarshal(api, settings, &response)
	if err != nil {
		return nil, err
	}

	if response.Result != saved {
		log.Printf("[TRACE] WireGuardSettingsSet response: %#v", response)
		return nil, fmt.Errorf("WireGuardSettingsSet failed: %w", ErrOpnsenseSave)
	}

	return &response, nil
}

func (c *Client) WireGuardEnableService() error {
	ws := WireGuardSettings{
		WireGuardSettingsGeneral{
			Enabled: "1",
		},
	}

	_, err := c.WireGuardSettingsSet(ws)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) WireGuardDisableService() error {
	ws := WireGuardSettings{
		WireGuardSettingsGeneral{
			Enabled: "0",
		},
	}

	_, err := c.WireGuardSettingsSet(ws)
	if err != nil {
		return err
	}

	return nil
}

type WireGuardClientBase struct {
	UUID          *uuid.UUID `json:"uuid,omitempty"`
	Enabled       string     `json:"enabled"`
	Name          string     `json:"name"`
	PubKey        string     `json:"pubkey"`
	Psk           string     `json:"psk"`
	ServerAddress string     `json:"serveraddress"`
	ServerPort    string     `json:"serverport"`
	KeepAlive     string     `json:"keepalive"`
}

type WireGuardClientSet struct {
	WireGuardClientBase
	TunnelAddress string `json:"tunneladdress"`
}

type WireGuardClientGet struct {
	WireGuardClientBase
	TunnelAddress SelectedMap `json:"tunneladdress"`
}

func (c *Client) WireGuardClientGet(uuid uuid.UUID) (*WireGuardClientGet, error) {
	api := path.Join("wireguard/client/getclient", uuid.String())

	type Response struct {
		Client WireGuardClientGet `json:"client"`
	}

	var response Response

	err := c.GetAndUnmarshal(api, &response)

	// UUID does not exist in the JSON, so we add it since we know it.
	response.Client.UUID = &uuid

	return &response.Client, err
}

func (c *Client) WireGuardClientGetUUIDs() ([]*uuid.UUID, error) {
	api := "wireguard/client/searchclient"

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

func (c *Client) WireGuardClientList() ([]*WireGuardClientGet, error) {
	uuids, err := c.WireGuardClientGetUUIDs()
	if err != nil {
		return nil, err
	}

	clients := []*WireGuardClientGet{}

	for _, uuid := range uuids {
		client, err := c.WireGuardClientGet(*uuid)
		if err == nil {
			clients = append(clients, client)
		}
	}

	return clients, nil
}

func (c *Client) WireGuardClientSet(uuid uuid.UUID, clientConf WireGuardClientSet) (*GenericResponse, error) {
	api := path.Join("wireguard/client/setclient", uuid.String())

	request := map[string]interface{}{
		"client": clientConf,
	}

	var response GenericResponse

	err := c.PostAndMarshal(api, request, &response)
	if err != nil {
		return nil, err
	}

	if response.Result != saved {
		log.Printf("[TRACE] WireGuardClientSet response: %#v", response)
		return nil, fmt.Errorf("WireGuardClientSet failed: %w", ErrOpnsenseSave)
	}

	return &response, nil
}

func (c *Client) WireGuardClientAdd(clientConf WireGuardClientSet) (*uuid.UUID, error) {
	api := "wireguard/client/addclient"

	request := map[string]interface{}{
		"client": clientConf,
	}

	var response GenericResponse

	err := c.PostAndMarshal(api, request, &response)
	if err != nil {
		return nil, err
	}

	if response.Result != saved {
		log.Printf("[TRACE] WireGuardClientAdd response: %#v", response)
		return nil, fmt.Errorf("WireGuardClientAdd failed: %w", ErrOpnsenseSave)
	}

	return response.UUID, nil
}

func (c *Client) WireGuardClientDelete(uuid uuid.UUID) (*GenericResponse, error) {
	api := path.Join("wireguard/client/delclient", uuid.String())

	var response GenericResponse

	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	if response.Result != deleted {
		log.Printf("[TRACE] WireGuardClientDelete response: %#v", response)
		return nil, fmt.Errorf("WireGuardClientDelete failed: %w", ErrOpnsenseDelete)
	}

	return &response, nil
}

type WireGuardServerBase struct {
	UUID          *uuid.UUID `json:"uuid,omitempty"`
	Enabled       string     `json:"enabled"`
	Name          string     `json:"name"`
	PubKey        string     `json:"pubkey"`
	PrivKey       string     `json:"privkey"`
	Port          string     `json:"port"`
	MTU           string     `json:"mtu"`
	DisableRoutes string     `json:"disableroutes"`
	Instance      *string    `json:"instance,omitempty"`
}

type WireGuardServerSet struct {
	WireGuardServerBase
	DNS           string `json:"dns"`
	TunnelAddress string `json:"tunneladdress"`
	// comma list of UUID
	Peers string `json:"peers"`
}

type WireGuardServerGet struct {
	WireGuardServerBase
	DNS           SelectedMap `json:"dns"`
	TunnelAddress SelectedMap `json:"tunneladdress"`
	Peers         SelectedMap `json:"peers"`
}

func (c *Client) WireGuardServerGet(uuid uuid.UUID) (*WireGuardServerGet, error) {
	api := path.Join("wireguard/server/getserver", uuid.String())

	type Response struct {
		Server WireGuardServerGet `json:"server"`
	}

	var response Response

	err := c.GetAndUnmarshal(api, &response)

	// UUID does not exist in the JSON, so we add it since we know it.
	response.Server.UUID = &uuid

	return &response.Server, err
}

func (c *Client) WireGuardServerGetUUIDs() ([]*uuid.UUID, error) {
	api := "wireguard/server/searchserver"

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

func (c *Client) WireGuardServerFindUUIDByName(name string) ([]*uuid.UUID, error) {
	api := "wireguard/server/searchserver"

	var response SearchResult

	err := c.GetAndUnmarshal(api, &response)
	if err != nil {
		return nil, err
	}

	uuids := []*uuid.UUID{}

	for _, row := range response.Rows {
		m := row.(map[string]interface{})
		if m["name"].(string) == name {
			uuid, err := uuid.FromString(m["uuid"].(string))
			if err == nil {
				uuids = append(uuids, &uuid)
			}
		}
	}

	return uuids, err
}

func (c *Client) WireGuardServerList() ([]*WireGuardServerGet, error) {
	uuids, err := c.WireGuardServerGetUUIDs()
	if err != nil {
		return nil, err
	}

	servers := []*WireGuardServerGet{}

	for _, uuid := range uuids {
		server, err := c.WireGuardServerGet(*uuid)
		if err == nil {
			servers = append(servers, server)
		}
	}

	return servers, nil
}

func (c *Client) WireGuardServerSet(uuid uuid.UUID, serverConf WireGuardServerSet) (*GenericResponse, error) {
	api := path.Join("wireguard/server/setserver", uuid.String())

	request := map[string]interface{}{
		"server": serverConf,
	}

	var response GenericResponse

	err := c.PostAndMarshal(api, request, &response)
	if err != nil {
		return nil, err
	}

	if response.Result != saved {
		log.Printf("[TRACE] WireGuardServerSet response: %#v", response)
		return nil, fmt.Errorf("WireGuardServerSet failed: %w", ErrOpnsenseSave)
	}

	return &response, nil
}

func (c *Client) WireGuardServerAdd(serverConf WireGuardServerSet) error {
	api := "wireguard/server/addserver"

	request := map[string]interface{}{
		"server": serverConf,
	}

	var response GenericResponse

	err := c.PostAndMarshal(api, request, &response)
	if err != nil {
		return err
	}

	if response.Result != saved {
		log.Printf("[TRACE] WireGuardServerAdd response: %#v", response)
		return fmt.Errorf("WireGuardServerAdd failed: %w", ErrOpnsenseSave)
	}

	return nil
}

func (c *Client) WireGuardServerDelete(uuid uuid.UUID) (*GenericResponse, error) {
	api := path.Join("wireguard/server/delserver", uuid.String())

	var response GenericResponse

	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	if response.Result != deleted {
		log.Printf("[TRACE] WireGuardServerDelete response: %#v", response)
		return nil, fmt.Errorf("WireGuardServerDelete failed: %w", ErrOpnsenseDelete)
	}

	return &response, nil
}
