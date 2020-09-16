package opnsense

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
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

	allowInsecure := getEnvAsBool("OPNSENSE_ALLOW_UNVERIFIED_TLS", insecureSkipVerify)

	httpClient := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			/* #nosec */
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: allowInsecure,
			},
		},
	}

	url, err := url.Parse(envy.Get("OPNSENSE_URL", baseURL))
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

	if resp.StatusCode == 401 {
		log.Printf("[ERROR] Failed to authenticate: %#v\n", string(body))

		return fmt.Errorf("GetAndUnmarshal failed: %w", ErrOpnsense401)
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

	log.Println("[TEMP]: ", string(body))

	if resp.StatusCode == 401 {
		log.Printf("[ERROR] Failed to authenticate: %#v\n", string(body))

		return fmt.Errorf("PostAndMarshal failed: %w", ErrOpnsense401)
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

type SearchResultPartial struct {
	RowCount int `json:"rowCount"`
	Total    int `json:"total"`
	Current  int `json:"current"`
}
