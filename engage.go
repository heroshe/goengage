// Copyright 2021 Heroshe Inc. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE.md file.

package goengage

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	userAgent string = "Heroshe - GoEngage Lib"
	apiUrl    string = "https://api.engage.so/v1"
)

type (
	service struct {
		client *Client
	}

	Client struct {
		httpClient   *http.Client
		BaseUrl      string
		UserAgent    string
		credentials  *Credentials
		commonClient service

		Users UserService
		Lists ListService
	}

	Error struct {
		Code    int
		Message string
	}
)

func New(config *Config) (*Client, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	c := &Client{
		BaseUrl:     apiUrl,
		credentials: config.Credentials,
		httpClient:  config.HTTPClient,
		UserAgent:   userAgent,
	}
	c.commonClient.client = c
	c.Users = (*Users)(&c.commonClient)
	c.Lists = (*Lists)(&c.commonClient)
	return c, nil
}

func (c *Client) newRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%v/%v", c.BaseUrl, endpoint)
	if strings.HasPrefix(endpoint, "/") {
		url = fmt.Sprintf("%v%v", c.BaseUrl, endpoint)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.credentials.publicKey, c.credentials.privateKey)
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (c *Client) makeRequest(req *http.Request, target interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	successCodes := []int{
		200, 201, 202, 203, 204, 205, 206,
	}
	for _, code := range successCodes {
		if code == resp.StatusCode {
			return json.Unmarshal(body, target)
		}
	}

	return Error{
		Code:    resp.StatusCode,
		Message: string(body),
	}
}

func validateConfig(config *Config) error {
	if config.Credentials.publicKey == "" {
		return fmt.Errorf("goengage: Public Key is required")
	}

	if config.Credentials.privateKey == "" {
		return fmt.Errorf("goengage: Private Key is required")
	}
	return nil
}

func (e Error) Error() string {
	return fmt.Sprintf("Go Engage Error - Code: %v | Message: %v", e.Code, e.Message)
}
