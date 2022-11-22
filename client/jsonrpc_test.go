package client

import (
	"encoding/json"
	"github.com/oliveagle/jsonpath"
	"testing"
)

func TestJsonRPCClient_Eth_Call(t *testing.T) {
	/*
		curl --location --request POST 'https://arb1.arbitrum.io/rpc' \
		--header 'Content-Type: application/json' \
		--data-raw '{"jsonrpc":"2.0","method":"eth_call","params":
		[{"data":"0x70a0823100000000000000000000000062383739D68Dd0F844103Db8dFb05a7EdED5BBE6","to":"0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9"},"0x2305876"],"id":1}'
	*/
	c := NewJsonRPCClient()
	requestParam := struct {
		Data string
		To   string
	}{
		Data: "0x70a0823100000000000000000000000062383739D68Dd0F844103Db8dFb05a7EdED5BBE6",
		To:   "0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9",
	}
	params := make([]interface{}, 0)
	blockHeight := "0x2305876"
	params = append(params, requestParam, blockHeight)
	request, _ := c.MakeJsonRPCRequestParams(1, "eth_call", params)
	// post
	rsp, err := c.Post("https://arb1.arbitrum.io/rpc", request, "", "", nil)
	if err != nil {
		t.Logf("get data from rpc client failed,error:%v", err)
	}
	t.Log(string(rsp))
	// parse json data
	var jsonData interface{}
	_ = json.Unmarshal(rsp, &jsonData)

	res, err := jsonpath.JsonPathLookup(jsonData, "$.result")
	t.Log(res)
}

func TestJsonRPCClient_Btc_Scantxoutset(t *testing.T) {
	/*
		curl --user myusername --data-binary '{"jsonrpc": "1.0", "id": "curltest", "method": "scantxoutset", "params": ["start", ["sh(multi(2,022bc0ca1d6aea1c1e523bfcb33f46131bd1a3240aa04f71
		c34b1a177cfd5ff933,035dc63727e7719824978161cdd94609db5235537bc8339a07b6838a6075f02530,033eeee979afb70450d2aebb17ace1b170a96199b495cdf3dd0631e
		b96aa21e6a8))"]]}' -H 'content-type: text/plain;' http://127.0.0.1:8332/
	*/
	c := NewJsonRPCClient()
	params := make([]interface{}, 0)
	descriptor01 := "wsh(multi(2,0375e00eb72e29da82b89367947f29ef34afb75e8654f6ea368e0acdfd92976b7c,03a1b26313f430c4b15bb1fdce663207659d8cac749a0e53d70eff01874496feff,03c96d495bfdd5ba4145e3e046fee45e84a8a48ad05bd8dbb395c011a32cf9f880))"
	descriptor02 := "sh(multi(2,02f1ff11c3b2f8f6a7a636043adde524e2e130ce23ca7b364ed868a69d3980da5c,034f0457f9ceb8ba7a6a49f242f565f12f7039c996b9dd9a63aa812d75e51638d1))"
	descriptorList := make([]interface{}, 0)
	descriptorList = append(descriptorList, descriptor01, descriptor02)
	params = append(params, "start", descriptorList)
	request, _ := c.MakeJsonRPCRequestParams(1, "scantxoutset", params)
	// post
	rsp, err := c.Post("http://127.0.0.1:8332", request, "123456", "123456", nil)
	if err != nil {
		t.Logf("get data from rpc client failed,error:%v", err)
	}
	t.Log(string(rsp))
	// parse json data
	var jsonData interface{}
	_ = json.Unmarshal(rsp, &jsonData)

	res, err := jsonpath.JsonPathLookup(jsonData, "$.result.total_amount")
	t.Log(res)
}

func TestJsonRPCClient_Btc_Getblockcount(t *testing.T) {
	c := NewJsonRPCClient()
	params := make([]interface{}, 0)
	request, _ := c.MakeJsonRPCRequestParams(1, "getblockcount", params)
	// post
	rsp, err := c.Post("http://127.0.0.1:8332", request, "123456", "123456", nil)
	if err != nil {
		t.Logf("get data from rpc client failed,error:%v", err)
	}
	t.Log(string(rsp))
	// parse json data
	var jsonData interface{}
	_ = json.Unmarshal(rsp, &jsonData)

	res, err := jsonpath.JsonPathLookup(jsonData, "$.result")
	t.Log(res)
}

func TestJsonRPCClient_Eth_Getbalance(t *testing.T) {
	c := NewJsonRPCClient()
	c.Debug = true
	params := make([]interface{}, 0)
	blockHeight := "0xf3dec4"
	params = append(params, "0x5E2B6c6B2240d582995537D3FafdB556E4A3822F", blockHeight)
	request, _ := c.MakeJsonRPCRequestParams(1, "eth_getBalance", params)
	// post
	rsp, err := c.Post("https://mainnet.infura.io/v3/xxx", request, "", "", nil)
	if err != nil {
		t.Logf("get data from rpc client failed,error:%v", err)
	}
	t.Log(string(rsp))
	// parse json data
	var jsonData interface{}
	_ = json.Unmarshal(rsp, &jsonData)

	res, err := jsonpath.JsonPathLookup(jsonData, "$.result")
	t.Log(res)
}
