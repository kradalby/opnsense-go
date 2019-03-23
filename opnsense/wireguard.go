package opnsense

import (
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
)

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

func (c *Client) WireGuardGetSettings() (*WireGuardSettings, error) {
	api := "wireguard/general/get"

	var response WireGuardSettings
	err := c.GetAndUnmarshal(api, &response)

	return &response, err
}

func (c *Client) WireGuardSetSettings(settings WireGuardSettings) (*GenericResponse, error) {
	api := "wireguard/general/set"

	var response GenericResponse
	err := c.PostAndMarshal(api, settings, &response)
	if err != nil {
		return nil, err
	}

	if response.Response != "saved" {
		err := errors.New(
			fmt.Sprintf("Failed to save, response from server: %s", response.Response),
		)
		log.Printf("[ERROR] %#v\n", err)
		return nil, err
	}

	return &response, nil
}

func (c *Client) WireGuardEnableService() error {
	ws := WireGuardSettings{
		WireGuardSettingsGeneral{
			Enabled: "1",
		},
	}
	_, err := c.WireGuardSetSettings(ws)
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
	_, err := c.WireGuardSetSettings(ws)
	if err != nil {
		return err
	}
	return nil
}

type WireGuardClientBase struct {
	UUID          uuid.UUID `json:"uuid"`
	Enabled       string    `json:"enabled"`
	Name          string    `json:"name"`
	PubKey        string    `json:"pubkey"`
	Psk           string    `json:"psk"`
	ServerAddress string    `json:"serveraddress"`
	ServerPort    string    `json:"serverport"`
	KeepAlive     string    `json:"keepalive"`
}

type WireGuardClientSet struct {
	WireGuardClientBase
	TunnelAddress string `json:"tunneladdress"`
}

type WireGuardClientGet struct {
	WireGuardClientBase
	TunnelAddress map[string]Selected `json:"tunneladdress"`
}

func (c *Client) WireGuardGetClient(uuid uuid.UUID) (*WireGuardClientGet, error) {
	api := "wireguard/client/getclient/" + uuid.String()

	var response WireGuardClientGet
	err := c.GetAndUnmarshal(api, &response)

	return &response, err
}

func (c *Client) WireGuardGetClientUUIDs() ([]*uuid.UUID, error) {
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

func (c *Client) WireGuardGetClients() ([]*WireGuardClientGet, error) {
	uuids, err := c.WireGuardGetClientUUIDs()
	if err != nil {
		return nil, err
	}

	clients := []*WireGuardClientGet{}
	for _, uuid := range uuids {
		client, err := c.WireGuardGetClient(*uuid)
		if err == nil {
			clients = append(clients, client)
		}
	}
	return clients, nil
}

func (c *Client) WireGuardSetClient(uuid uuid.UUID, clientConf WireGuardClientSet) (*GenericResponse, error) {
	api := "wireguard/client/setclient/" + uuid.String()

	var response GenericResponse
	err := c.PostAndMarshal(api, clientConf, &response)
	if err != nil {
		return nil, err
	}

	if response.Response != "saved" {
		err := errors.New(
			fmt.Sprintf("Failed to save, response from server: %s", response.Response),
		)
		log.Printf("[ERROR] %#v\n", err)
		return nil, err
	}

	return &response, nil
}

func (c *Client) WireGuardAddClient(clientConf WireGuardClientSet) (*uuid.UUID, error) {
	api := "wireguard/client/addclient"

	var response GenericResponseUUID
	err := c.PostAndMarshal(api, clientConf, &response)
	if err != nil {
		return nil, err
	}

	if response.Response != "saved" {
		err := errors.New(
			fmt.Sprintf("Failed to save, response from server: %s", response.Response),
		)
		log.Printf("[ERROR] %#v\n", err)
		return nil, err
	}

	return &response.UUID, nil
}

func (c *Client) WireGuardDeleteClient(uuid uuid.UUID) (*GenericResponse, error) {
	api := "wireguard/client/addclient/" + uuid.String()

	var response GenericResponse
	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	if response.Response != "delete" {
		err := errors.New(
			fmt.Sprintf("Failed to delete, response from server: %s", response.Response),
		)
		log.Printf("[ERROR] %#v\n", err)
		return nil, err
	}

	return &response, nil
}

type WireGuardServerBase struct {
	UUID          uuid.UUID `json:"uuid"`
	Enabled       string    `json:"enabled"`
	Name          string    `json:"name"`
	PubKey        string    `json:"pubkey"`
	PrivKey       string    `json:"privkey"`
	Port          string    `json:"port"`
	DisableRoutes string    `json:"disableroutes"`
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
	DNS           map[string]Selected `json:"dns"`
	TunnelAddress map[string]Selected `json:"tunneladdress"`
	Peers         map[string]Selected `json:"peers"`
}

func (c *Client) WireGuardGetServer(uuid uuid.UUID) (*WireGuardServerGet, error) {
	api := "wireguard/server/getserver/" + uuid.String()

	var response WireGuardServerGet
	err := c.GetAndUnmarshal(api, &response)

	return &response, err
}

func (c *Client) WireGuardGetServerUUIDs() ([]*uuid.UUID, error) {
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

func (c *Client) WireGuardGetServers() ([]*WireGuardServerGet, error) {
	uuids, err := c.WireGuardGetServerUUIDs()
	if err != nil {
		return nil, err
	}

	servers := []*WireGuardServerGet{}
	for _, uuid := range uuids {
		server, err := c.WireGuardGetServer(*uuid)
		if err == nil {
			servers = append(servers, server)
		}
	}
	return servers, nil
}

func (c *Client) WireGuardSetServer(uuid uuid.UUID, serverConf WireGuardServerSet) (*GenericResponse, error) {
	api := "wireguard/server/setserver/" + uuid.String()

	var response GenericResponse
	err := c.PostAndMarshal(api, serverConf, &response)
	if err != nil {
		return nil, err
	}

	if response.Response != "saved" {
		err := errors.New(
			fmt.Sprintf("Failed to save, response from server: %s", response.Response),
		)
		log.Printf("[ERROR] %#v\n", err)
		return nil, err
	}

	return &response, nil
}

func (c *Client) WireGuardAddServer(serverConf WireGuardServerSet) (*uuid.UUID, error) {
	api := "wireguard/server/addserver"

	var response GenericResponseUUID
	err := c.PostAndMarshal(api, serverConf, &response)
	if err != nil {
		return nil, err
	}

	if response.Response != "saved" {
		err := errors.New(
			fmt.Sprintf("Failed to save, response from server: %s", response.Response),
		)
		log.Printf("[ERROR] %#v\n", err)
		return nil, err
	}

	return &response.UUID, nil
}

func (c *Client) WireGuardDeleteServer(uuid uuid.UUID) (*GenericResponse, error) {
	api := "wireguard/server/addserver/" + uuid.String()

	var response GenericResponse
	err := c.PostAndMarshal(api, nil, &response)
	if err != nil {
		return nil, err
	}

	if response.Response != "delete" {
		err := errors.New(
			fmt.Sprintf("Failed to delete, response from server: %s", response.Response),
		)
		log.Printf("[ERROR] %#v\n", err)
		return nil, err
	}

	return &response, nil
}
