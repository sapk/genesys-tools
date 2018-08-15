// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/sirupsen/logrus"

	"github.com/sapk/go-genesys/api/object"
)

//Gensys GAX Api client
type Client struct {
	BaseURL   *url.URL
	UserAgent string

	httpClient *http.Client
}

func NewClient(host string) *Client {
	//TODO manage https
	cookieJar, _ := cookiejar.New(nil)
	c := &Client{BaseURL: &url.URL{Host: host, Scheme: "http", Path: "/gax/api/"}, UserAgent: "NXOTestLab/0.0", httpClient: &http.Client{
		Jar: cookieJar,
	}}
	return c
}

func (c *Client) Login(user, pass string) (*object.LoginResponse, error) {
	req, err := c.newRequest("POST", "session/login", object.LoginRequest{user, pass, false})
	//{"username":"root","password":"","isPasswordEncrypted":true}
	if err != nil {
		return nil, err
	}
	_, err = c.do(req, nil)

	//Check logged user
	req, err = c.newRequest("GET", "user/info", nil)
	if err != nil {
		return nil, err
	}
	var u object.LoginResponse
	_, err = c.do(req, &u)
	return &u, err
}

func (c *Client) ListObject(t string, v interface{}) (*http.Response, error) {
	req, err := c.newRequest("GET", "cfg/objects", nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = "brief=false&type=" + t
	return c.do(req, v)
}

func (c *Client) ListApplication() ([]object.CfgApplication, error) {
	var apps []object.CfgApplication
	_, err := c.ListObject("CfgApplication", &apps)
	return apps, err
}

func (c *Client) ListHost() ([]object.CfgHost, error) {
	var apps []object.CfgHost
	_, err := c.ListObject("CfgHost", &apps)
	return apps, err
}

func (c *Client) ListDN() ([]object.CfgDN, error) {
	var apps []object.CfgDN
	_, err := c.ListObject("CfgDN", &apps)
	return apps, err
}

func (c *Client) ListSwitch() ([]object.CfgSwitch, error) {
	var apps []object.CfgSwitch
	_, err := c.ListObject("CfgSwitch", &apps)
	return apps, err
}

func (c *Client) ListPlace() ([]object.CfgPlace, error) {
	var apps []object.CfgPlace
	_, err := c.ListObject("CfgPlace", &apps)
	return apps, err
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	logrus.WithFields(logrus.Fields{
		"Method":  req.Method,
		"Path":    req.URL.Path,
		"Cookies": req.Cookies(),
		"Body":    req.Body,
	}).Debug("Executing request")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Method": req.Method,
			"Path":   req.URL.Path,
			"Error":  err,
		}).Debug("Request failed")
		return nil, err
	}
	defer resp.Body.Close()

	logrus.WithFields(logrus.Fields{
		"Method": req.Method,
		"Path":   req.URL.Path,
		"Body":   resp.Body,
	}).Debug("Request response")
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}
