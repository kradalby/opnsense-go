package opnsense

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gobuffalo/envy"
	uuid "github.com/satori/go.uuid"
)

const (
	saved   = "saved"
	deleted = "deleted"
	done    = "done"
)

var (
	ErrOpnsenseSave              = errors.New("failed to save")
	ErrOpnsenseDelete            = errors.New("failed to delete")
	ErrOpnsenseDone              = errors.New("did not finish")
	ErrOpnsenseStatusNotOk       = errors.New("status did not return ok")
	ErrOpnsenseEmptyListNotFound = errors.New("found empty array, most likely 404")
	ErrOpnsense500               = errors.New("internal server error")
)

type Client struct {
	baseURL *url.URL
	key     string
	secret  string
	c       *http.Client
}

func NewClient(baseURL, key, secret string, insecureSkipVerify bool) (*Client, error) {
	log.Printf(
		"[TRACE] Creating new OPNsense client with url: %s, key: %s, secret: %s, insecure: %t",
		baseURL,
		key,
		secret,
		insecureSkipVerify,
	)

	httpClient := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			/* #nosec */
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecureSkipVerify,
			},
		},
	}

	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	log.Printf("[TRACE] Parsed URL: %s", url.String())

	client := &Client{
		baseURL: url,
		key:     envy.Get("OPNSENSE_KEY", key),
		secret:  envy.Get("OPNSENSE_SECRET", secret),
		c:       httpClient,
	}

	log.Printf("[TRACE] Finished setting up client: %#v", client)

	return client, nil
}

func (c *Client) Get(api string) (resp *http.Response, err error) {
	// url := path.Join(c.baseURL.String(), api)
	url := c.baseURL.String() + "/api/" + api
	log.Printf("[TRACE] GET to %s", url)

	request, err := http.NewRequestWithContext(context.TODO(), "GET", url, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to create GET request: %#v\n\n", err)
		return nil, err
	}

	// request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(c.key, c.secret)

	return c.c.Do(request)
}

func (c *Client) GetAndUnmarshal(api string, responseData interface{}) error {
	resp, err := c.Get(api)
	if err != nil {
		log.Printf("[ERROR] Failed to GET request: %s\n", err)
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read GET response: %#v\n", err)
		return err
	}

	// add possible internal error from the OPNSense API
	if resp.StatusCode == 500 {
		log.Printf("[ERROR] Internal Error status code received: %#v\n", string(body))
		return fmt.Errorf("GetAndUnmarshal failed: %w", ErrOpnsense500)
	}

	// The OPNsense API does not return 404 when you fetch something that does
	// not exist, but returns an empty list instead. Check for the empty list
	// and return a 404 error instead so implmenters could handle that error
	// differently
	if string(body) == "[]" {
		return ErrOpnsenseEmptyListNotFound
	}

	err = json.Unmarshal(body, responseData)

	log.Printf("[TRACE] Response for URL: %s\n", api)
	log.Printf("[TRACE] Response body: %s\n", string(body))

	if err != nil {
		log.Printf("[ERROR] Failed to unmarshal GET response: %#v\n", err)
		return err
	}

	return nil
}

func (c *Client) Post(api string, body io.Reader) (resp *http.Response, err error) {
	// url := path.Join(c.baseURL.String(), api)
	url := c.baseURL.String() + "/api/" + api
	log.Printf("[TRACE] POST to %s", url)

	request, err := http.NewRequestWithContext(context.TODO(), "POST", url, body)
	if err != nil {
		log.Printf("[ERROR] Failed to create POST request: %#v\n", err)
		return nil, err
	}

	request.SetBasicAuth(c.key, c.secret)
	request.Header.Set("Content-Type", "application/json")

	return c.c.Do(request)
}

func (c *Client) PostAndMarshal(api string, requestData interface{}, responseData interface{}) error {
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal requestData for POST request: %#v\n", err)
		return err
	}

	log.Printf("[TRACE] Request payload: %s", string(requestBody))

	resp, err := c.Post(api, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("[ERROR] Failed to POST request: %#v\n", err)
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read POST response: %#v\n", err)
		return err
	}

	err = json.Unmarshal(body, responseData)
	if err != nil {
		log.Printf("[ERROR] Failed to unmarshal POST response: %#v\n", err)
		log.Printf("[ERROR] Failed to unmarshal POST response: %s\n", string(body))

		return err
	}

	return nil
}

// Generic types

type StatusMessage struct {
	Status string    `json:"status"`
	MsgID  uuid.UUID `json:"msg_uuid"`
}

type GenericResponse struct {
	Result      string            `json:"result"`
	UUID        *uuid.UUID        `json:"uuid,omitempty"`
	Validations map[string]string `json:"validations,omitempty"`
}

type SearchResult struct {
	Rows     []interface{} `json:"rows"`
	RowCount int           `json:"rowCount"`
	Total    int           `json:"total"`
	Current  int           `json:"current"`
}

// Helpers.
type SelectedMap map[string]Selected

// The OPNsense API returns a [] when there is no
// objects in the list of selected items. This is
// very inconvinient and this function tries to work
// around this by making the map pointer an empty map
// if the there is an empty array.
func (sm *SelectedMap) UnmarshalJSON(b []byte) error {
	*sm = SelectedMap{}

	type Alias SelectedMap

	var temp2 Alias

	err := json.Unmarshal(b, &temp2)
	if err != nil {
		var temp []string

		err := json.Unmarshal(b, &temp)
		if err != nil {
			return err
		}

		return nil
	}

	for key, value := range temp2 {
		(*sm)[key] = value
	}

	return nil
}

type Selected struct {
	Value    string `json:"value"`
	Selected int    `json:"selected"`
}

func (s *Selected) UnmarshalJSON(b []byte) error {
	*s = Selected{}

	type Alias Selected

	var temp Alias

	err := json.Unmarshal(b, &temp)
	if err != nil {
		type Selected2 struct {
			Value    string `json:"value"`
			Selected bool   `json:"selected"`
		}

		var temp2 Selected2

		err := json.Unmarshal(b, &temp2)
		if err != nil {
			return err
		}

		s.Value = temp2.Value
		if temp2.Selected {
			s.Selected = 1
		} else {
			s.Selected = 0
		}
	}

	s.Value = temp.Value
	s.Selected = temp.Selected

	return nil
}

func ListSelectedValues(m SelectedMap) []string {
	s := []string{}

	for _, value := range m {
		if value.Selected == 1 {
			s = append(s, value.Value)
		}
	}

	return s
}

func ListSelectedKeys(m SelectedMap) []string {
	s := []string{}

	for key, value := range m {
		if value.Selected == 1 {
			s = append(s, key)
		}
	}

	return s
}
