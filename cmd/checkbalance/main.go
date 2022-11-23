package main

import (
	"fmt"
	"github.com/okex/proof-of-reserves/client"
	"github.com/okex/proof-of-reserves/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"math"
	"math/big"
	"os"

	"strings"
	"time"
)

var coin, addr, mode, rpcJsonFileName, porCsvFileName string

var rootCmd = &cobra.Command{
	Use:   "checkbalance",
	Short: "check balance",
	Long:  ``,
	Run:   CoinAddressBalanceValidator,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&coin, "coin_name", "ETH", "")
	rootCmd.PersistentFlags().StringVar(&addr, "address", "", "")
	rootCmd.PersistentFlags().StringVar(&mode, "mode", "", "")
	rootCmd.PersistentFlags().StringVar(&rpcJsonFileName, "rpc_json_filename", "rpc.json", "")
	rootCmd.PersistentFlags().StringVar(&porCsvFileName, "por_csv_filename", "", "")
}

func CoinAddressBalanceValidator(cmd *cobra.Command, args []string) {
	start := time.Now().UTC()
	// init log
	common.ConfigLocalFilesystemLogger("./logs/", "check_balance.log", 7*time.Hour*24, time.Second*20)
	// init http client
	client.HttpClient = client.NewHTTPClient()
	// init rpc client
	client.RpcClient = client.NewJsonRPCClient()
	// init por csv data
	log.Info("loading por csv data...")
	var err error
	common.PorCoinDataMap, common.PorCoinGeneralDataMap, err = common.InitPorCsvDataMap(porCsvFileName)
	if err != nil {
		log.Errorf("load por csv data failed, error: %v", err)
		return
	}

	porCoinTotalBalance := make(map[string]*big.Int)
	for coinName := range common.PorCoinGeneralDataMap {
		if _, exist := porCoinTotalBalance[coinName]; !exist {
			porCoinTotalBalance[coinName] = big.NewInt(0)
		}
	}

	// scan por data
	for k, v := range common.PorCoinDataMap {
		keys := strings.Split(k, ":")
		coinName := keys[0]

		if _, exist := porCoinTotalBalance[coinName]; exist {
			if v.Balance == "" {
				log.Errorf("balance is null, coin: %s, address: %s", coinName, v.Address)
				continue
			}
			b, _ := big.NewInt(0).SetString(v.Balance, 10)
			porCoinTotalBalance[coinName] = porCoinTotalBalance[coinName].Add(porCoinTotalBalance[coinName], b)
		}

		// stats address count
		if item, exist := common.PorCoinGeneralDataMap[coinName]; exist {
			item.AddressCount++
		}

		// recover btc P2PKH address pubkey
		if coinName == "btc" {
			addrTye := common.GuessAddressType(v.Address)
			if addrTye == "P2PKH" {
				// recover pubKey
				pubkey := common.RecoveryPubKeyFromSign(v.Address, v.Message, v.Sign1)
				v.Script = pubkey
				if pubkey == "" {
					log.Errorf("recovery pubkey from sign msg failed, coin: %s, address: %s", coinName, v.Address)
				}
			}
		}
	}

	for coinName, balance := range porCoinTotalBalance {
		if _, exist := common.PorCoinGeneralDataMap[coinName]; exist {
			common.PorCoinGeneralDataMap[coinName].Balance = balance.String()
		}
	}

	// init Validator
	validator, err := common.NewAddressBalanceValidator(rpcJsonFileName)
	if err != nil {
		log.Errorf("init rpc.json failed, error: %v, please check the rpc json file", err)
		return
	}

	coin = strings.ToLower(coin)
	mode = strings.ToLower(mode)
	// check coin_name
	if coin != "" {
		if _, exist := common.PorCoinGeneralDataMap[coin]; !exist {
			log.Errorf("por data not support the coin %s, please set the correct one!", coin)
			return
		}
	}

	// check address
	if addr != "" {
		if coin == "" {
			log.Error("you must set the coin_name")
			return
		} else {
			if _, exist := common.PorCoinDataMap[fmt.Sprintf("%s:%s", coin, addr)]; !exist {
				log.Errorf("por data not support the coin %s, address %s, please set the correct one!", coin, addr)
				return
			}
		}
	}

	switch mode {
	case "single_address":
		if coin == "" {
			log.Error("you must choose set the coin_name")
			return
		}
		if addr == "" {
			log.Error("you must choose set the address")
			return
		}

		start = time.Now().UTC()
		log.Infof("start to verify coin %s, address %s balance...", coin, addr)
		VerifySingleAddressBalance(validator, coin, addr)
		log.Infof("verify coin %s, address %s balance finished, consume time %ds", coin, addr, time.Now().UTC().Unix()-start.Unix())

	case "single_coin":
		if coin == "" {
			log.Error("you must choose set the coin_name")
			return
		}

		start = time.Now().UTC()
		log.Infof("start to verify coin %s every signle address balance...", coin)
		VerifySingleCoinAllAddressBalance(validator, coin)
		log.Infof("verify coin %s total address balance finished, consume time %ds", coin, time.Now().UTC().Unix()-start.Unix())

	case "all_coin":
		start = time.Now().UTC()
		log.Info("start to verify all coin address balance...")
		VerifyAllCoinAddressBalance(validator)
		log.Infof("verify all coin address balance finished, consume time %ds", time.Now().UTC().Unix()-start.Unix())

	case "single_coin_total_balance":
		if coin == "" {
			log.Error("you must choose set the coin_name")
			return
		}

		start = time.Now().UTC()
		log.Infof("start to verify coin %s total address balance...", coin)
		VerifyCoinAddressTotalBalance(validator, coin)
		log.Infof("verify coin %s total address balance finished, consume time %ds", coin, time.Now().UTC().Unix()-start.Unix())

	case "all_coin_total_balance":
		start = time.Now().UTC()
		log.Info("start to verify all coin total address balance...")
		VerifyAllCoinAddressTotalBalance(validator)
		log.Infof("verify all coin total address balance finished, consume time %ds", time.Now().UTC().Unix()-start.Unix())

	default:
		// return por data info
		BTCAmount, ETHAmount, USDTAmount := big.NewFloat(0), big.NewFloat(0), big.NewFloat(0)
		for _, value := range common.PorCoinGeneralDataMap {
			coinName := strings.Split(value.Coin, "-")
			b, _ := big.NewFloat(0).SetString(value.Balance)
			switch coinName[0] {
			case "btc":
				BTCAmount = BTCAmount.Add(BTCAmount, b)
				b = b.Mul(b, big.NewFloat(math.Pow(10, -8)))
			case "eth":
				ETHAmount = ETHAmount.Add(ETHAmount, b)
				b = b.Mul(b, big.NewFloat(math.Pow(10, -18)))
			case "usdt":
				USDTAmount = USDTAmount.Add(USDTAmount, b)
				b = b.Mul(b, big.NewFloat(math.Pow(10, -6)))
			default:
				log.Errorf("unsupport coin %s in por dara.", coinName)
			}
			bFloat64, _ := b.Float64()
			log.Infof("por data, coin: %s, snapshot height: %s, address count: %d, total balance: %0.4f", strings.ToUpper(value.Coin), value.SnapshotHeight,
				value.AddressCount, bFloat64)
		}

		BTCAmountFloat64 := BTCAmount.Mul(BTCAmount, big.NewFloat(math.Pow(10, -8)))
		ETHAmountFloat64 := ETHAmount.Mul(ETHAmount, big.NewFloat(math.Pow(10, -18)))
		USDTAmountFloat64 := USDTAmount.Mul(USDTAmount, big.NewFloat(math.Pow(10, -6)))

		// BTC,ETH,USDT total balance
		log.Infof("por data: BTC total balance %0.4f, ETH(ALL) total balance %0.4f, USDT(ALL) total balance %0.4f", BTCAmountFloat64, ETHAmountFloat64, USDTAmountFloat64)
	}
}

func VerifySingleAddressBalance(validator *common.AddressBalanceValidator, coin, addr string) {
	coinInfo, exist := common.PorCoinGeneralDataMap[coin]
	if !exist {
		log.Errorf("unsupport the coin %s, please check the coin_name!", coin)
		return
	}
	height := coinInfo.SnapshotHeight
	balance, err := validator.GetCoinAddressBalanceInfo(coin, addr, height)
	if err != nil {
		log.Errorf("get coin %s, address %s balance from blockchain failed!", coin, addr)
		return
	}
	// compare
	if value, exist := common.PorCoinDataMap[fmt.Sprintf("%s:%s", coin, addr)]; exist {
		if balance == value.Balance {
			log.Infof("verify coin %s, address %s balance success, in chain balance: %s, in por balance: %s", coin, addr, balance, value.Balance)
		} else {
			log.Infof("verify coin %s, address %s balance failed, in chain balance: %s, in por balance: %s", coin, addr, balance, value.Balance)
		}
	}
}

func VerifySingleCoinAllAddressBalance(validator *common.AddressBalanceValidator, coin string) {
	coinInfo, exist := common.PorCoinGeneralDataMap[coin]
	if !exist {
		log.Errorf("unsupport the coin %s, please choose the correct one!", coin)
		return
	}
	height := coinInfo.SnapshotHeight
	coinDataList := make([]*common.CoinData, 0)
	for _, value := range common.PorCoinDataMap {
		if value.Coin == coin {
			coinDataList = append(coinDataList, value)
		}
	}

	for _, v := range coinDataList {
		balance, err := validator.GetCoinAddressBalanceInfo(coin, v.Address, height)
		if err != nil {
			log.Errorf("get address %s balance from blockchain failed, error: %v", v.Address, err)
			continue
		}
		// compare
		if value, exist := common.PorCoinDataMap[fmt.Sprintf("%s:%s", coin, v.Address)]; exist {
			if balance == value.Balance {
				log.Infof("verify coin %s, address %s balance success, in chain balance: %s, in por balance: %s", coin, v.Address, balance, value.Balance)
			} else {
				log.Infof("verify coin %s, address %s balance failed, in chain balance:%s, in por balance:%s", coin, v.Address, balance, value.Balance)
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func VerifyAllCoinAddressBalance(validator *common.AddressBalanceValidator) {
	for _, v := range common.PorCoinDataMap {
		coinInfo, exist := common.PorCoinGeneralDataMap[v.Coin]
		if !exist {
			log.Errorf("unsupport the coin %s, please choose the correct one!", coin)
			continue
		}
		height := coinInfo.SnapshotHeight
		balance, err := validator.GetCoinAddressBalanceInfo(v.Coin, v.Address, height)
		if err != nil {
			log.Error(err)
		}
		// compare
		if balance == v.Balance {
			log.Infof("verify coin %s, address %s balance success, in chain balance: %s, in por balance: %s", v.Coin, v.Address, balance, v.Balance)
		} else {
			log.Infof("verify coin %s, address %s balance failed, in chain balance: %s, in por balance: %s", v.Coin, v.Address, balance, v.Balance)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func VerifyCoinAddressTotalBalance(validator *common.AddressBalanceValidator, coin string) {
	coinInfo, exist := common.PorCoinGeneralDataMap[coin]
	if !exist {
		log.Errorf("unsupport the coin %s, please choose the correct one!", coin)
		return
	}
	height := coinInfo.SnapshotHeight
	addressList := make([]string, 0)
	for _, value := range common.PorCoinDataMap {
		if value.Coin == coin {
			if value.Address == "" {
				continue
			}

			if coin == "btc" && value.Script == "" {
				continue
			}

			addressList = append(addressList, value.Address)
		}
	}

	totalBalance, err := validator.GetCoinAddressTotalBalance(coin, height, addressList)
	if err != nil {
		log.Errorf("get coin %s total address balance from blockchain failed, error: %v", coin, err)
		return
	}
	// compare
	if totalBalance == coinInfo.Balance {
		log.Infof("verify coin %s total address balance success, in chain balance: %s, in por balance: %s", coin, totalBalance, coinInfo.Balance)
	} else {
		log.Infof("verify coin %s total address balance failed, in chain balance: %s, in por balance: %s", coin, totalBalance, coinInfo.Balance)
	}
}

func VerifyAllCoinAddressTotalBalance(validator *common.AddressBalanceValidator) {
	for _, value := range common.PorCoinGeneralDataMap {
		VerifyCoinAddressTotalBalance(validator, value.Coin)
	}
}

func main() {
	Execute()
}
