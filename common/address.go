package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/okex/proof-of-reserves/client"
	"github.com/oliveagle/jsonpath"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

type AddressBalanceValidator struct {
	once                        sync.Once
	coinJSONConfig              string
	confMap                     map[string]*coin
	confCoinAddressWhiteListMap map[string]*coinAddress
}

type coin struct {
	Name string `json:"name"`
	Coin string `json:"coin"`
	API  struct {
		Endpoint      string            `json:"endpoint"`
		JSONPattern   string            `json:"jsonPattern"`
		DefaultUnit   string            `json:"defaultUnit"`
		TokenAddress  string            `json:"tokenAddress"`
		CustomHeaders map[string]string `json:"customHeaders"`
		Enabled       bool              `json:"enabled"`
	} `json:"api"`
	RPC struct {
		Endpoint      string            `json:"endpoint"`
		JSONPattern   string            `json:"jsonPattern"`
		DefaultUnit   string            `json:"defaultUnit"`
		CustomHeaders map[string]string `json:"customHeaders"`
		AuthUser      string            `json:"authUser"`
		AuthPassword  string            `json:"authPassword"`
		TokenAddress  string            `json:"tokenAddress"`
		Enabled       bool              `json:"enabled"`
	} `json:"rpc"`
	WhiteList []*coinAddress `json:"whiteList"`
}

type coinAddress struct {
	Project      string `json:"project"`
	Address      string `json:"address"`
	Height       string `json:"height"`
	TokenAddress string `json:"tokenAddress"`
	Balance      string `json:"balance"`
}

func NewAddressBalanceValidator(coinJSONConfig string) (*AddressBalanceValidator, error) {
	runner := &AddressBalanceValidator{
		coinJSONConfig: coinJSONConfig,
	}

	if err := runner.loadCoinJSONByFile(); err != nil {
		log.WithField("fileName", coinJSONConfig).Error("load rpc json file failed!")
		return nil, err
	}

	return runner, nil
}

func (r *AddressBalanceValidator) loadCoinJSONByFile() error {
	content, err := ioutil.ReadFile(r.coinJSONConfig)
	if err != nil {
		log.WithField("filename", r.coinJSONConfig).WithField("error", err).Error("load rpc json file failed!")
		return err
	}

	return r.loadCoinJSON(content)
}

func (r *AddressBalanceValidator) loadCoinJSON(content []byte) error {
	data := struct {
		Coins []*coin `json:"coins"`
	}{}
	err := json.Unmarshal(content, &data)
	if err != nil {
		log.WithField("content", string(content)).WithField("error", err).Error("unmarshal json failed!")
		return err
	}
	coinMap := make(map[string]*coin)
	addressWhiteListMap := make(map[string]*coinAddress)
	for _, value := range data.Coins {
		if _, exist := coinMap[value.Name]; !exist {
			coinMap[value.Name] = value
		}
		if len(value.WhiteList) > 0 {
			for _, addr := range value.WhiteList {
				key := fmt.Sprintf("%s:%s:%s", value.Name, addr.Address, addr.Height)
				if _, exist := addressWhiteListMap[key]; !exist {
					addressWhiteListMap[key] = addr
				}
				if addr.Project == "" || addr.Balance == "" || addr.Address == "" || addr.Height == "" {
					err = errors.New(fmt.Sprintf("%s not config the whitelist project/address/height/balance filed", value.Name))
					return err
				}
			}
		}
	}
	r.confMap = coinMap
	r.confCoinAddressWhiteListMap = addressWhiteListMap

	return err
}

func (r *AddressBalanceValidator) GetCoinAddressBalanceInfo(coin, address, height string) (result string, err error) {
	pConf, exist := r.confMap[coin]
	if !exist {
		err = errors.New(fmt.Sprintf("coin %s not exist in rpc json file, please check the json file!", coin))
		log.Error(err)
		return
	}
	return r.GetCoinAddressBalanceInfoByJSONFormat(address, height, pConf)
}

func (r *AddressBalanceValidator) GetCoinAddressBalanceInfoByJSONFormat(address, height string, pConf *coin) (result string, err error) {
	if !pConf.RPC.Enabled && !pConf.API.Enabled {
		err = errors.New(fmt.Sprintf("coin %s, rpc or api method must be enabled at leasr one in rpc json file, please check the json file!", pConf.Name))
		log.Error(err)
		return result, err
	}
	var object interface{}
	var balanceRes interface{}
	var defaultUnit string
	if pConf.RPC.Enabled {
		// get address balance from white list
		key := fmt.Sprintf("%s:%s:%s", pConf.Name, address, height)
		if _, exist := r.confCoinAddressWhiteListMap[key]; exist {
			return r.confCoinAddressWhiteListMap[key].Balance, nil
		}

		var request *client.JsonRpcRequest
		switch pConf.Name {
		case "btc":
			params := make([]interface{}, 0)
			descriptorList := make([]interface{}, 0)
			descriptor, err := r.generateAddressDescriptor(pConf.Name, address)
			if err != nil {
				return result, err
			}
			descriptorList = append(descriptorList, descriptor)
			params = append(params, "start", descriptorList)
			request, _ = client.RpcClient.MakeJsonRPCRequestParams(1, "scantxoutset", params)
		case "eth", "eth-optimism", "eth-arbitrum":
			params := make([]interface{}, 0)
			var blockHeightHex string
			if height != "latest" {
				h, _ := strconv.ParseInt(height, 10, 64)
				blockHeightHex = fmt.Sprintf("0x%s", strconv.FormatInt(h, 16))
			} else {
				blockHeightHex = height
			}
			params = append(params, address, blockHeightHex)
			request, _ = client.RpcClient.MakeJsonRPCRequestParams(1, "eth_getBalance", params)
		default:
			var cutAddress string
			if strings.HasPrefix(address, "0x") {
				cutAddress = address[2:]
			}
			requestParamData := fmt.Sprintf("0x70a08231000000000000000000000000%s", cutAddress)
			requestParam := struct {
				Data string
				To   string
			}{
				Data: requestParamData,
				To:   pConf.RPC.TokenAddress,
			}
			params := make([]interface{}, 0)
			var blockHeightHex string
			if height != "latest" {
				h, _ := strconv.ParseInt(height, 10, 64)
				blockHeightHex = fmt.Sprintf("0x%s", strconv.FormatInt(h, 16))
			} else {
				blockHeightHex = height
			}
			params = append(params, requestParam, blockHeightHex)
			request, _ = client.RpcClient.MakeJsonRPCRequestParams(1, "eth_call", params)
		}
		body, err := client.RpcClient.Post(pConf.RPC.Endpoint, request, pConf.RPC.AuthUser, pConf.RPC.AuthPassword, pConf.RPC.CustomHeaders)
		if err != nil {
			err = errors.New(fmt.Sprintf("call blockchain node rpc method failed, coin:%s, address:%s, height:%s, error:%v", pConf.Name, address, height, err))
			log.Error(err)
			return result, err
		}

		var object interface{}
		err = json.Unmarshal(body, &object)
		if err != nil {
			err = errors.New(fmt.Sprintf("unmarshall json data from blockchain node failed, coin:%s, address:%s, height:%s, error:%v", pConf.Name, address, height, err))
			log.Error(err)
			return result, err
		}

		// parse address balance
		balanceRes, err = jsonpath.JsonPathLookup(object, pConf.RPC.JSONPattern)
		if err != nil {
			log.Infof("json object:%v, json path:%s", object, pConf.RPC.JSONPattern)
			err = errors.New(fmt.Sprintf("parse json data from blockchain node failed, coin:%s, address:%s, height:%s, error:%v", pConf.Name, address, height, err))
			log.Error(err)
			return result, err
		}
		defaultUnit = pConf.RPC.DefaultUnit
	} else {
		if pConf.API.Endpoint == "" {
			err = errors.New(fmt.Sprintf("coin %s, api method endpoint is null, please check the json file!", pConf.Name))
			log.Error(err)
			return result, err
		}
		var project, tokenAddress string
		// get address coin name from white list
		key := fmt.Sprintf("%s:%s:%s", pConf.Name, address, height)
		if _, exist := r.confCoinAddressWhiteListMap[key]; exist {
			project = r.confCoinAddressWhiteListMap[key].Project
			tokenAddress = r.confCoinAddressWhiteListMap[key].TokenAddress
		} else {
			project = ""
			tokenAddress = pConf.API.TokenAddress
		}

		args := make(map[string]string)
		args["address"] = address
		// api params chain must set pConf coin
		args["chainShortName"] = pConf.Coin
		args["height"] = height
		if project != "" {
			args["project"] = project
		}
		if tokenAddress != "" {
			args["tokenContractAddress"] = tokenAddress
		}
		body, err := client.HttpClient.Get(client.HttpClient.MakeGetURL(pConf.API.Endpoint, args), pConf.API.CustomHeaders)
		if err != nil {
			log.Infof("request params: %v", args)
			err = errors.New(fmt.Sprintf("call api %s failed, coin:%s, address:%s, height:%s, error:%v", pConf.API.Endpoint, pConf.Name, address, height, err))
			log.Error(err)
			return result, err
		}

		err = json.Unmarshal(body, &object)
		if err != nil {
			err = errors.New(fmt.Sprintf("unmarshall json data from api %s failed, coin:%s, address:%s, height:%s, error:%v", pConf.API.Endpoint, pConf.Name, address, height, err))
			log.Error(err)
			return result, err
		}
		code, _ := jsonpath.JsonPathLookup(object, "$.code")
		if code.(string) == "0" {
			// parse address balance
			balanceRes, err = jsonpath.JsonPathLookup(object, pConf.API.JSONPattern)
			if err != nil {
				log.Infof("json object:%v, json path:%s", object, pConf.API.JSONPattern)
				err = errors.New(fmt.Sprintf("parse json data from api %s failed, coin:%s, address:%s, height:%s, error:%v", pConf.API.Endpoint, pConf.Name, address, height, err))
				log.Error(err)
				return result, err
			}
		} else {
			log.Infof("json object:%v", object)
			msg, _ := jsonpath.JsonPathLookup(object, "$.msg")
			if code.(string) == "50040" || msg.(string) == "No data is displayed for this block height." {
				log.Infof("coin:%s, address:%s, height:%s, %s", pConf.Name, address, height, msg.(string))
				return "0", nil
			} else {
				err = errors.New(fmt.Sprintf("call api %s failed, coin:%s, address:%s, height:%s, %v", pConf.API.Endpoint, pConf.Name,
					address, height, msg.(string)))
				log.Error(err)
				return result, err

			}
		}
		defaultUnit = pConf.API.DefaultUnit
	}
	balanceStr := r.ParseBalanceValue(balanceRes)
	balanceDecimal, _ := decimal.NewFromString(balanceStr)

	if defaultUnit != "" {
		coinNameList := strings.Split(pConf.Name, "-")
		switch strings.ToLower(coinNameList[0]) {
		case "btc":
			result = balanceDecimal.Mul(decimal.NewFromFloat(math.Pow(10, 8))).String()
		case "eth":
			result = balanceDecimal.Mul(decimal.NewFromFloat(math.Pow(10, 18))).String()
		case "usdt":
			result = balanceDecimal.Mul(decimal.NewFromFloat(math.Pow(10, 6))).String()
		default:
			log.Errorf("unsupport coin name %s in rpc json file.", pConf.Name)
		}
	} else {
		result = balanceDecimal.String()
	}

	return result, nil
}

func (r *AddressBalanceValidator) GetCoinAddressTotalBalance(coin, height string, addresses []string) (result string, err error) {
	pConf, exist := r.confMap[coin]
	if !exist {
		err = errors.New(fmt.Sprintf("coin %s not exist in rpc json file, please check the json file!", coin))
		log.Error(err)
		return
	}

	var chunkSize int
	if coin == "btc" {
		chunkSize = 10000
	} else {
		chunkSize = 1000
	}
	if len(addresses) < chunkSize {
		chunkSize = len(addresses)
	}
	addressItems := make([]interface{}, 0)
	for _, v := range addresses {
		addressItems = append(addressItems, v)
	}
	// divided address list
	divided := r.DividedAddressList(addressItems, chunkSize)
	log.Infof("coin:%s, chunk amount:%d, chunk size:%d", pConf.Name, len(divided), chunkSize)

	totalBalance := big.NewInt(0)
	for i, items := range divided {
		log.Infof("chunk %d, scanning address total balance, this may take a while...", i+1)
		addressList := make([]interface{}, 0)
		if pConf.Name == "btc" && pConf.RPC.Enabled {
			for _, item := range items {
				address := item.(string)
				// ignore white list address
				key := fmt.Sprintf("%s:%s:%s", pConf.Name, address, height)
				if _, exist := r.confCoinAddressWhiteListMap[key]; exist {
					balance := r.confCoinAddressWhiteListMap[key].Balance
					balanceInt, _ := big.NewInt(0).SetString(balance, 10)
					totalBalance = totalBalance.Add(totalBalance, balanceInt)
					continue
				}

				descriptor, err := r.generateAddressDescriptor(pConf.Name, address)
				if err != nil {
					return result, err
				}
				addressList = append(addressList, descriptor)
			}
		} else {
			addressList = items
		}

		chunkBalance := big.NewInt(0)
		if pConf.Name == "btc" && pConf.RPC.Enabled {
			retryNums := 3
			chunkBalance, err = r.BatchFetchBTCTotalAddressBalanceFromNode(height, addressList, pConf)
			if err != nil {
				for retryNums > 0 {
					log.Infof("get chunk %d coin total address balance failed, retry...", i+1)
					chunkBalance, err = r.BatchFetchBTCTotalAddressBalanceFromNode(height, addressList, pConf)
					if err == nil {
						break
					}
					retryNums--
				}
				if retryNums == 0 {
					log.Infof("chunk %d, try to cut the size and rertry...", i+1)
					reChunkSize := chunkSize
					chunkBalance = big.NewInt(0)
					for {
						reChunkSize = reChunkSize / 2
						log.Infof("chunk %d, cut chunk size to %d", i+1, reChunkSize)
						// divided address list
						reDivided := r.DividedAddressList(addressList, reChunkSize)
						var errs error
						for _, item := range reDivided {
							b, err := r.BatchFetchBTCTotalAddressBalanceFromNode(height, item, pConf)
							if err != nil {
								errs = err
								chunkBalance = big.NewInt(0)
								break
							}
							chunkBalance = chunkBalance.Add(chunkBalance, b)
						}
						if errs == nil {
							break
						}
						if reChunkSize <= 100 {
							log.Errorf("get chunk %d coin total address balance from blockchain failed, please check rpc json config.", i+1)
							return result, errs
						}
					}
				}
			}
		} else {
			chunkBalance, err = r.BatchFetchCoinTotalAddressBalance(height, addressList, pConf)
		}

		totalBalance = totalBalance.Add(totalBalance, chunkBalance)
		log.Infof("chunk %d, chunk balance %s, total balance %s", i+1, chunkBalance.String(), totalBalance.String())
	}

	return totalBalance.String(), nil
}

func (r *AddressBalanceValidator) BatchFetchBTCTotalAddressBalanceFromNode(height string, addresses []interface{}, pConf *coin) (result *big.Int, err error) {
	result = big.NewInt(0)

	var request *client.JsonRpcRequest
	params := make([]interface{}, 0)
	descriptors := addresses
	params = append(params, "start", descriptors)
	request, _ = client.RpcClient.MakeJsonRPCRequestParams(1, "scantxoutset", params)

	body, err := client.RpcClient.Post(pConf.RPC.Endpoint, request, pConf.RPC.AuthUser, pConf.RPC.AuthPassword, pConf.RPC.CustomHeaders)
	if err != nil {
		err = errors.New(fmt.Sprintf("get batch address total balance, call blockchain node rpc method failed, coin:%s, height:%s, error:%v", pConf.Name, height, err))
		log.Error(err)
		return result, err
	}

	var object interface{}
	err = json.Unmarshal(body, &object)
	if err != nil {
		err = errors.New(fmt.Sprintf("get batch address total balance, unmarshall json data from blockchain node failed, coin:%s, height:%s, error:%v", pConf.Name, height, err))
		log.Error(err)
		return result, err
	}

	// parse address balance
	balance, err := jsonpath.JsonPathLookup(object, pConf.RPC.JSONPattern)
	if err != nil {
		log.Infof("json object:%v, json path:%s", object, pConf.RPC.JSONPattern)
		err = errors.New(fmt.Sprintf("get batch address total balance, parse json data from blockchain node failed, coin:%s, height:%s, error:%v", pConf.Name, height, err))
		log.Error(err)
		return result, err
	}

	balanceStr := r.ParseBalanceValue(balance)
	balanceDecimal, _ := decimal.NewFromString(balanceStr)
	if pConf.RPC.DefaultUnit != "" {
		// convert BTC to Satoshi
		balanceDecimal = balanceDecimal.Mul(decimal.NewFromFloat(math.Pow(10, 8)))
	}
	balanceInt, _ := new(big.Int).SetString(balanceDecimal.String(), 10)
	result = result.Add(result, balanceInt)

	return result, nil
}

// BatchFetchCoinTotalAddressBalance batch get coin total address balance
// not support btc rpc mode
func (r *AddressBalanceValidator) BatchFetchCoinTotalAddressBalance(height string, addresses []interface{}, pConf *coin) (result *big.Int, err error) {
	result = big.NewInt(0)

	for _, address := range addresses {
		addressStr := address.(string)
		balance, err := r.GetCoinAddressBalanceInfoByJSONFormat(addressStr, height, pConf)
		if err != nil {
			retryNum := 3
			for retryNum > 0 {
				// some api limit rate is low...
				time.Sleep(10 * time.Second)

				balance, err = r.GetCoinAddressBalanceInfoByJSONFormat(addressStr, height, pConf)
				if err == nil {
					break
				}
				retryNum--
			}

			log.Errorf("get coin %s address %s balance failed..", pConf.Name, address)
			balance = "0"
		}

		balanceInt, _ := new(big.Int).SetString(balance, 10)
		result = result.Add(result, balanceInt)

		time.Sleep(500 * time.Millisecond)
	}

	return result, nil
}

func (r *AddressBalanceValidator) ParseBalanceValue(object interface{}) (balance string) {
	balanceType := reflect.TypeOf(object)
	switch balanceType.String() {
	case "int64":
		balance = strconv.FormatInt(object.(int64), 10)
	case "string":
		balanceStr := object.(string)
		if strings.HasPrefix(balanceStr, "0x") {
			n := new(big.Int)
			n, _ = n.SetString(balanceStr[2:], 16)
			balanceStr = n.String()
		}
		balance = balanceStr
	case "float64":
		balance = strconv.FormatFloat(object.(float64), 'f', -1, 64)
	}

	return
}

func (r *AddressBalanceValidator) DividedAddressList(addresses []interface{}, chunkSize int) (divided [][]interface{}) {
	if len(addresses) <= chunkSize {
		divided = append(divided, addresses)
	} else {
		for i := 0; i < len(addresses); i += chunkSize {
			end := i + chunkSize
			if end > len(addresses) {
				end = len(addresses)
			}
			divided = append(divided, addresses[i:end])
		}
	}
	return
}

func (r *AddressBalanceValidator) generateAddressDescriptor(coin, address string) (result string, err error) {
	addrType := GuessAddressType(address)
	if addrType == "" {
		err = errors.New(fmt.Sprintf("coin:%s, invalid address %s", coin, address))
		log.Error(err)
		return result, err
	}
	var redeemScript string
	if value, exist := PorCoinDataMap[fmt.Sprintf("%s:%s", coin, address)]; exist {
		redeemScript = value.Script
	} else {
		err = errors.New(fmt.Sprintf("coin:%s, por data not support the address %s", coin, address))
		log.Error(err)
		return result, err
	}
	descriptor, err := CreateAddressDescriptor(addrType, redeemScript, 2, 3)
	if err != nil {
		err = errors.New(fmt.Sprintf("coin:%s, address %s, create address output descriptor failed.", coin, address))
		log.Error(err)
		return result, err
	}

	return descriptor, nil
}
