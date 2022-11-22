package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

type HTTPClient struct {
	client http.Client
	Debug  bool
}

var HttpClient *HTTPClient

// NewHTTPClient init http client
func NewHTTPClient() *HTTPClient {
	tr := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConnsPerHost: 100,
		MaxIdleConns:        0,
		MaxConnsPerHost:     0,
	}
	client := http.Client{
		Timeout:   10 * 60 * time.Second,
		Transport: tr,
	}
	return &HTTPClient{client: client}
}

func (c *HTTPClient) SetTimeout(t time.Duration) {
	c.client.Timeout = t
}

func (c *HTTPClient) SetEndpointType(a bool) string {
	if a {
		return "https"
	}
	return "http"
}

func (c *HTTPClient) DisableKeepAlive() {
	c.client.Transport.(*http.Transport).DisableKeepAlives = true
}

func (c *HTTPClient) MakeGetURL(endpoint string, args map[string]string) string {
	url := endpoint
	first := true
	if strings.Contains(url, "?") {
		first = false
	}
	for k, v := range args {
		if first {
			url = fmt.Sprintf("%s?%s=%s", url, k, v)
			first = false
		} else {
			url = fmt.Sprintf("%s&%s=%s", url, k, v)
		}
	}
	return url
}

func (c *HTTPClient) MakeRequest(ctx context.Context, method, endpoint string, data io.Reader, vs []interface{}) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, endpoint, data)
	if err != nil {
		return nil, err
	}
	for _, v := range vs {
		switch vv := v.(type) {
		case http.Header:
			for key, values := range vv {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}
		case map[string]string:
			for key, value := range vv {
				req.Header[key] = []string{value}
			}
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, err
}

func (c *HTTPClient) Get(endpoint string, vs ...interface{}) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*60*time.Second)
	defer cancel()
	req, err := c.MakeRequest(ctx, "GET", endpoint, nil, vs)
	if err != nil {
		return nil, err
	}
	if c.Debug {
		bytes, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(bytes))
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if c.Debug {
		dump, _ := httputil.DumpResponse(res, true)
		fmt.Println(string(dump))
	}
	return body, nil
}

func (c *HTTPClient) Post(endpoint string, data io.Reader, vs ...interface{}) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*60*time.Second)
	defer cancel()
	req, err := c.MakeRequest(ctx, "POST", endpoint, data, vs)
	if err != nil {
		return nil, err
	}
	if c.Debug {
		bytes, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(bytes))
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if c.Debug {
		dump, _ := httputil.DumpResponse(res, true)
		fmt.Println(string(dump))
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}
