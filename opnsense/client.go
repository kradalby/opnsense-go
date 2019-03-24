package opnsense

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/gobuffalo/envy"
	"github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	// "path"
	"time"
)

type Client struct {
	baseURL *url.URL
	key     string
	secret  string
	c       *http.Client
}

func NewClient(baseUrl, key, secret string, InsecureSkipVerify bool) (*Client, error) {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: InsecureSkipVerify,
			},
		},
	}

	url, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	client := &Client{
		baseURL: url,
		key:     envy.Get("OPNSENSE_KEY", key),
		secret:  envy.Get("OPNSENSE_SECRET", secret),
		c:       httpClient,
	}

	return client, nil
}

func (c *Client) Get(api string) (resp *http.Response, err error) {
	// url := path.Join(c.baseURL.String(), api)
	url := c.baseURL.String() + "/" + api
	log.Printf("[TRACE] GET to %s", url)

	request, err := http.NewRequest("GET", url, nil)
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
		log.Printf("[ERROR] Failed to GET request: %#v\n", err)
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read GET response: %#v\n", err)
		return err
	}

	err = json.Unmarshal(body, responseData)
	log.Printf("[TRACE] Response body: %s\n", string(body))
	if err != nil {
		log.Printf("[ERROR] Failed to unmarshal GET response: %#v\n", err)
		return err
	}

	return nil
}

func (c *Client) Post(api string, body io.Reader) (resp *http.Response, err error) {
	// url := path.Join(c.baseURL.String(), api)
	url := c.baseURL.String() + "/" + api
	log.Printf("[TRACE] POST to %s", url)

	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Printf("[ERROR] Failed to create POST request: %#v\n", err)
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(c.key, c.secret)

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

// Helpers
type Selected struct {
	Value    string `json:"value"`
	Selected int    `json:"selected"`
}

func listSelected(m map[string]Selected) []string {
	s := []string{}
	for _, value := range m {
		if value.Selected == 1 {
			s = append(s, value.Value)
		}
	}
	return s
}
