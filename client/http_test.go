package client

import (
	"encoding/json"
	"github.com/oliveagle/jsonpath"
	"testing"
)

func TestHTTPClient_ETH_GetBalance(t *testing.T) {
	args := make(map[string]string)
	args["address"] = "0x6cC5F688a315f3dC28A7781717a9A798a59fDA7b"
	args["chain"] = "eth"
	args["height"] = "15958559"
	header := map[string]string{
		"x-apiKey":     "xxx",
		"Content-Type": "application/x-www-form-urlencoded",
		"Accept":       "*/*",
	}
	client := NewHTTPClient()
	body, err := client.Get(client.MakeGetURL("https://www.oklink.com/api/explorer/v1/por/getArchiveBalance", args), header)
	if err != nil {
		t.Error(err)
		return
	}

	var object interface{}
	err = json.Unmarshal(body, &object)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("body:", object)
	// parse address balance
	balance, err := jsonpath.JsonPathLookup(object, "$.data.balance")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(balance)
}

func TestHTTPClient_ETH_GetTokenBalance(t *testing.T) {
	args := make(map[string]string)
	args["address"] = "0xF977814e90dA44bFA03b6295A0616a897441aceC"
	args["chain"] = "polygon"
	args["height"] = "35520001"
	args["tokenContractAddress"] = "0xc2132D05D31c914a87C6611C10748AEb04B58e8F"
	header := map[string]string{
		"x-apiKey":     "xxx",
		"Content-Type": "application/x-www-form-urlencoded",
		"Accept":       "*/*",
	}
	client := NewHTTPClient()
	body, err := client.Get(client.MakeGetURL("https://www.oklink.com/api/explorer/v1/por/getArchiveBalance", args), header)
	if err != nil {
		t.Error(err)
		return
	}

	var object interface{}
	err = json.Unmarshal(body, &object)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("body:", object)
	// parse address balance
	balance, err := jsonpath.JsonPathLookup(object, "$.data.balance")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(balance)
}
