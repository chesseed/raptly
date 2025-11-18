// Package aptly provides a client to access the aptly REST API in go
package aptly

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL string

	client *http.Client

	basic_user *string
	basic_pw   *string
}

func NewClient(url string, config *tls.Config) *Client {
	client := new(Client)
	client.BaseURL = url

	tr := &http.Transport{
		TLSClientConfig: config,
	}
	client.client = &http.Client{Transport: tr}
	return client
}

// GetClient get http client used for advanced use cases like testing
func (c *Client) GetClient() *http.Client {
	return c.client
}

func (c *Client) get(url string) *request {
	return c.newRequest(http.MethodGet, url)
}

func (c *Client) post(url string) *request {
	return c.newRequest(http.MethodPost, url)
}

func (c *Client) put(url string) *request {
	return c.newRequest(http.MethodPut, url)
}

func (c *Client) delete(url string) *request {
	return c.newRequest(http.MethodDelete, url)
}

func (c *Client) SetBasicAuth(user string, password string) *Client {
	c.basic_user = &user
	c.basic_pw = &password
	return c
}

func (c *Client) newRequest(method string, url string) *request {
	r := initRequest(method, url)
	return r
}

// send prepared request and get error
func (c *Client) Send(r *request) (*http.Response, error) {
	req, err := r.GetRawRequest(c.BaseURL)
	if err != nil {
		return nil, err
	}
	if c.basic_user != nil {
		req.SetBasicAuth(*c.basic_user, *c.basic_pw)
	}
	return c.client.Do(req)
}

func checkResponseForError(res *http.Response) error {
	if res.StatusCode < 400 {
		return nil
	}
	var apiErr APIError
	decodeErr := json.NewDecoder(res.Body).Decode(&apiErr)
	if decodeErr != nil || apiErr.ErrorMsg == nil {
		return fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	return &apiErr
}

func callAPIWithResponse(c *Client, r *request) (*http.Response, error) {
	res, err := c.Send(r)
	if err != nil {
		return nil, err
	}
	err = checkResponseForError(res)
	if err != nil {
		res.Body.Close()
		return nil, err
	}
	return res, nil
}

func callAPI(c *Client, r *request) error {
	res, err := c.Send(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	err = checkResponseForError(res)
	if err != nil {
		return err
	}
	return nil
}

// T can be any type
func callAPIwithResult[T any](c *Client, r *request) (T, error) {
	var ret T

	res, err := c.Send(r)
	if err != nil {
		return ret, err
	}
	defer res.Body.Close()
	err = checkResponseForError(res)
	if err != nil {
		return ret, err
	}

	decodeErr := json.NewDecoder(res.Body).Decode(&ret)
	if decodeErr != nil {
		return ret, err
	}
	return ret, nil
}

type APIError struct {
	// as pointer to distinguish between valid error and failed parsing errors (like empty bodies)
	ErrorMsg *string `json:"error"`
}

func (e *APIError) Error() string {
	if e.Valid() {
		return *e.ErrorMsg
	}
	return ""
}

func (e *APIError) Valid() bool {
	return e.ErrorMsg != nil
}
