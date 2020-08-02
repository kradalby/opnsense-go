package opnsense

import (
	"fmt"
	"log"
	"path"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type AliasBase struct {
	Enabled     string `json:"enabled"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Updatefreq  string `json:"updatefreq"`
	Counters    string `json:"counters"`
}

type AliasSet struct {
	AliasBase
	Type    string `json:"type"`
	Proto   string `json:"proto"`
	Content string `json:"content"`
}

type AliasGet struct {
	Type    map[string]AliasNestedValue `json:"type"`
	Proto   map[string]AliasNestedValue `json:"proto"`
	Content map[string]AliasNestedValue `json:"content"`
	AliasBase
}

type AliasNestedValue struct {
	Value    string `json:"value"`
	Selected int    `json:"selected"`
}

type AliasList struct {
	Rows     []AliasListItem `json:"rows"`
	RowCount int             `json:"rowCount"`
	Total    int             `json:"total"`
	Current  int             `json:"current"`
}

type AliasListItem struct {
	UUID        string `json:"uuid"`
	Enabled     string `json:"enabled"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Content     string `json:"content"`
}

type AliasFormat struct {
	UUID        *uuid.UUID
	Enabled     bool
	Name        string
	Description string
	Updatefreq  string
	Counters    string
	Type        string
	Proto       string
	Content     []string
}

type AliasReconfigureResponse struct {
	Status string `json:"status"`
}

func (c *Client) AliasGet(uuid uuid.UUID) (*AliasFormat, error) {
	type Response struct {
		Alias AliasGet `json:"alias"`
	}

	var rawResponse Response

	err := c.GetAndUnmarshal(path.Join("firewall/alias/getItem", uuid.String()), &rawResponse)
	if err != nil {
		return nil, err
	}

	var response AliasFormat
	response.UUID = &uuid
	response.Enabled = rawResponse.Alias.Enabled == "1"
	response.Name = rawResponse.Alias.Name
	response.Description = rawResponse.Alias.Description
	response.Updatefreq = rawResponse.Alias.Updatefreq
	response.Counters = rawResponse.Alias.Counters

	for k, v := range rawResponse.Alias.Type {
		if v.Selected == 1 && k != "" {
			response.Type = k
			break
		}
	}

	for k, v := range rawResponse.Alias.Content {
		if v.Selected == 1 && k != "" {
			response.Content = append(response.Content, v.Value)
		}
	}

	return &response, err
}

func (c *Client) AliasGetList() (*AliasList, error) {
	var response AliasList

	err := c.GetAndUnmarshal("firewall/alias/searchItem", &response)
	if err != nil {
		return nil, err
	}

	return &response, err
}

func (c *Client) AliasUpdate(uuid uuid.UUID, conf AliasFormat) (*GenericResponse, error) {
	type Request struct {
		Alias AliasSet `json:"alias"`
	}

	var request Request
	request.Alias = AliasFormatToSet(conf)

	var response GenericResponse

	err := c.PostAndMarshal(path.Join("firewall/alias/setItem", uuid.String()), request, &response)
	if err != nil {
		return nil, err
	}

	if response.Result != saved {
		log.Printf("[TRACE] AliasUpdate response: %#v", response)
		return nil, fmt.Errorf("AliasUpdate failed: %w", ErrOpnsenseSave)
	}

	return &response, nil
}

func (c *Client) AliasAdd(conf AliasFormat) (*uuid.UUID, error) {
	type Request struct {
		Alias AliasSet `json:"alias"`
	}

	var request Request

	request.Alias = AliasFormatToSet(conf)

	var response GenericResponse

	err := c.PostAndMarshal("firewall/alias/addItem", request, &response)
	if err != nil {
		return nil, err
	}

	if response.Result != saved {
		log.Printf("[TRACE] AliasAdd response: %#v", response)
		return nil, fmt.Errorf("AliasAdd failed: %w", ErrOpnsenseSave)
	}

	return response.UUID, nil
}

func AliasFormatToSet(conf AliasFormat) AliasSet {
	var set AliasSet

	if conf.Enabled {
		set.Enabled = "1"
	} else {
		set.Enabled = "0"
	}

	set.Name = conf.Name
	set.Description = conf.Description
	set.Updatefreq = conf.Updatefreq
	set.Counters = conf.Counters
	set.Type = conf.Type
	set.Proto = conf.Proto
	set.Content = strings.Join(conf.Content, "\n")

	return set
}

func (c *Client) AliasDelete(uuid uuid.UUID) (*GenericResponse, error) {
	var response GenericResponse

	err := c.PostAndMarshal(path.Join("firewall/alias/delItem", uuid.String()), nil, &response)
	if err != nil {
		return nil, err
	}

	if response.Result != deleted {
		log.Printf("[TRACE] AliasDelete response: %#v", response)
		return nil, fmt.Errorf("AliasDelete failed: %w", ErrOpnsenseDelete)
	}

	return &response, nil
}

func (c *Client) AliasReconfigure() (*AliasReconfigureResponse, error) {
	var response AliasReconfigureResponse

	request := map[string]interface{}{}

	err := c.PostAndMarshal("firewall/alias/reconfigure", request, &response)
	if err != nil {
		return nil, err
	}

	if response.Status != "ok" {
		log.Printf("[TRACE] AliasReconfigure response: %#v", response)
		return nil, fmt.Errorf("AliasReconfigure failed: %w", ErrOpnsenseStatusNotOk)
	}

	return &response, nil
}

//// ALIAS UILS SECTION ////.
type AliasUtilsGet struct {
	Name     string           `json:"name"`
	Total    int              `json:"total"`
	RowCount int              `json:"rowCount"`
	Current  int              `json:"current"`
	Rows     []AliasUtilsInfo `json:"rows"`
}

type AliasUtilsInfo struct {
	Address string `json:"ip"`
}

type AliasUtilsSet struct {
	Address string `json:"address"`
}

type AliasUtilsResponse struct {
	Status string `json:"status"`
}

func (c *Client) AliasUtilsGet(name string) (*AliasUtilsGet, error) {
	var response AliasUtilsGet

	err := c.GetAndUnmarshal(path.Join("firewall/alias_util/list", name), &response)
	if err != nil {
		return nil, err
	}

	response.Name = name

	return &response, nil
}

func (c *Client) AliasUtilsAdd(name string, request AliasUtilsSet) (*AliasUtilsResponse, error) {
	var response AliasUtilsResponse

	err := c.PostAndMarshal(path.Join("firewall/alias_util/add", name), request, &response)
	if err != nil {
		return nil, err
	}

	if response.Status != done {
		log.Printf("[TRACE] AliasUtilsGet response: %#v", response)
		return nil, fmt.Errorf("AliasUtilsGet failed: %w", ErrOpnsenseDone)
	}

	return &response, nil
}

func (c *Client) AliasUtilsDel(name string, request AliasUtilsSet) (*AliasUtilsResponse, error) {
	var response AliasUtilsResponse

	err := c.PostAndMarshal(path.Join("firewall/alias_util/delete", name), request, &response)
	if err != nil {
		return nil, err
	}

	if response.Status != done {
		log.Printf("[TRACE] AliasUtilsDel response: %#v", response)
		return nil, fmt.Errorf("AliasUtilsDel failed: %w", ErrOpnsenseDone)
	}

	return &response, nil
}
