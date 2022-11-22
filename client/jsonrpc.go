package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
)

type JsonRPCClient struct {
	client http.Client
	Debug  bool
}

var RpcClient *JsonRPCClient

func NewJsonRPCClient() *JsonRPCClient {
	tr := &http.Transport{
		MaxIdleConnsPerHost:    100,
		MaxIdleConns:           0,
		MaxConnsPerHost:        0,
		MaxResponseHeaderBytes: 32 * 1024 * 1024,
		WriteBufferSize:        16 * 1024 * 1024,
		ReadBufferSize:         16 * 1024 * 1024,
	}
	client := http.Client{
		Timeout:   60 * 60 * time.Second,
		Transport: tr,
	}
	return &JsonRPCClient{client: client}
}

func (c *JsonRPCClient) SetTimeout(t time.Duration) {
	c.client.Timeout = t
}

type JsonRpcRequest struct {
	JsonRPC string            `json:"jsonrpc"`
	ID      interface{}       `json:"id"`
	Method  string            `json:"method"`
	Params  []json.RawMessage `json:"params"`
}

func (c *JsonRPCClient) MakeJsonRPCRequestParams(id interface{}, method string, params []interface{}) (*JsonRpcRequest, error) {
	rawParams := make([]json.RawMessage, 0, len(params))
	for _, param := range params {
		marshalledParam, err := json.Marshal(param)
		if err != nil {
			return nil, err
		}
		rawMessage := json.RawMessage(marshalledParam)
		rawParams = append(rawParams, rawMessage)
	}

	return &JsonRpcRequest{
		JsonRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  rawParams,
	}, nil
}

func (c *JsonRPCClient) MakeRequest(ctx context.Context, method, endpoint string, data io.Reader, vs []interface{}) (*http.Request, error) {
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

// Post : post method
func (c *JsonRPCClient) Post(endpoint string, request *JsonRpcRequest, username, password string, vs ...interface{}) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*60*time.Second)
	defer cancel()
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	bodyReader := bytes.NewReader(data)

	req, err := c.MakeRequest(ctx, "POST", endpoint, bodyReader, vs)
	if err != nil {
		log.Error("request failed, ", err)
		return nil, err
	}

	req.SetBasicAuth(username, password)

	if c.Debug {
		bytes, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(bytes))
	}
	res, err := c.client.Do(req)
	if err != nil {
		log.Error("client do request failed, ", err)
		return nil, err
	}
	if c.Debug {
		dump, _ := httputil.DumpResponse(res, true)
		fmt.Println(string(dump))
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("client read response body failed, ", err)
		return nil, err
	}

	return body, err
}
