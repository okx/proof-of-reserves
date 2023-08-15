package main

import (
	"fmt"
	"github.com/okx/proof-of-reserves/client"
	"github.com/okx/proof-of-reserves/common"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"math"
	"os"

	"strings"
	"time"
)

var coin, addr, mode, rpcJsonFileName, porCsvFileName string

var porCoinTotalBalance map[string]decimal.Decimal

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
	// set decimal precision
	decimal.DivisionPrecision = 18
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
	common.PorCoinDataMap, err = common.InitPorCsvDataMap(porCsvFileName)
	if err != nil {
		log.Errorf("load por csv data failed, error: %v", err)
		return
	}

	porCoinTotalBalance = make(map[string]decimal.Decimal)
	// scan por data
	for k, v := range common.PorCoinDataMap {
		keys := strings.Split(k, ":")
		coinName := keys[0]
		if v.Balance == "" {
			log.Errorf("balance is null, coin: %s, address: %s", coinName, v.Address)
			continue
		}
		b, _ := decimal.NewFromString(v.Balance)
		if _, exist := porCoinTotalBalance[coinName]; exist {
			porCoinTotalBalance[coinName] = porCoinTotalBalance[coinName].Add(b)
		} else {
			porCoinTotalBalance[coinName] = b
		}

		// recover btc P2PKH address pubkey
		if coinName == "BTC" {
			addrTye := common.GuessUtxoCoinAddressType(v.Address)
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

	// init Validator
	validator, err := common.NewAddressBalanceValidator(rpcJsonFileName)
	if err != nil {
		log.Errorf("init rpc.json failed, error: %v, please check the rpc json file", err)
		return
	}

	coin = strings.ToUpper(coin)
	mode = strings.ToLower(mode)
	// check coin_name
	if coin != "" {
		if _, exist := common.PorCoinUnitMap[coin]; !exist {
			log.Errorf("por data not support the coin %s, please set the correct one!", coin)
			return
		}
		// coin black list
		if common.IsCheckBalanceBannedCoin(coin) {
			log.Errorf("check balance not support the coin %s, please set the correct one!", coin)
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

	start = time.Now().UTC()
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

		log.Infof("start to verify coin %s, address %s balance...", coin, addr)
		VerifySingleAddressBalance(validator, coin, addr)
		log.Infof("verify coin %s, address %s balance finished, consume time %ds", coin, addr, time.Now().UTC().Unix()-start.Unix())

	case "single_coin":
		if coin == "" {
			log.Error("you must choose set the coin_name")
			return
		}

		log.Infof("start to verify coin %s every signle address balance...", coin)
		VerifySingleCoinAllAddressBalance(validator, coin)
		log.Infof("verify coin %s total address balance finished, consume time %ds", coin, time.Now().UTC().Unix()-start.Unix())

	case "all_coin":
		log.Info("start to verify all coin address balance...")
		VerifyAllCoinAddressBalance(validator)
		log.Infof("verify all coin address balance finished, consume time %ds", time.Now().UTC().Unix()-start.Unix())

	case "single_coin_total_balance":
		if coin == "" {
			log.Error("you must choose set the coin_name")
			return
		}

		log.Infof("start to verify coin %s total address balance...", coin)
		VerifyCoinAddressTotalBalance(validator, coin)
		log.Infof("verify coin %s total address balance finished, consume time %ds", coin, time.Now().UTC().Unix()-start.Unix())

	case "all_coin_total_balance":
		log.Info("start to verify all coin total address balance...")
		VerifyAllCoinAddressTotalBalance(validator)
		log.Infof("verify all coin total address balance finished, consume time %ds", time.Now().UTC().Unix()-start.Unix())

	default:
		// return por data info
		log.Info("por coin total balance:")
		for coin, value := range porCoinTotalBalance {
			log.Infof("coin %s, por total balance %s.", coin, value.String())
		}
	}
}

func VerifySingleAddressBalance(validator *common.AddressBalanceValidator, coin, addr string) {
	value, exist := common.PorCoinDataMap[fmt.Sprintf("%s:%s", coin, addr)]
	if !exist {
		log.Errorf("unsupport the coin %s, please check the coin_name!", coin)
		return
	}
	height := value.SnapshotHeight
	balance, err := validator.GetCoinAddressBalanceInfo(strings.ToLower(coin), addr, height)
	if err != nil {
		log.Errorf("get coin %s, address %s balance from blockchain failed!", coin, addr)
		return
	}

	balance = convertCoinBalanceToBaseUnit(coin, balance, -1)
	// compare
	if isCoinBalanceEqual(balance, value.Balance) {
		log.Infof("verify coin %s, address %s balance success, in chain balance: %s, in por balance: %s", coin, addr, balance, value.Balance)
	} else {
		log.Infof("verify coin %s, address %s balance failed, in chain balance: %s, in por balance: %s", coin, addr, balance, value.Balance)
	}
}

func VerifySingleCoinAllAddressBalance(validator *common.AddressBalanceValidator, coin string) {
	destCoins := getDestCoinList(coin)
	coinDataList := make([]*common.CoinData, 0)
	for _, value := range common.PorCoinDataMap {
		for _, destCoin := range destCoins {
			if destCoin == value.Coin {
				coinDataList = append(coinDataList, value)
			}
		}
	}

	for _, v := range coinDataList {
		balance, err := validator.GetCoinAddressBalanceInfo(strings.ToLower(v.Coin), v.Address, v.SnapshotHeight)
		if err != nil {
			log.Errorf("get address %s balance from blockchain failed, error: %v", v.Address, err)
			continue
		}
		balance = convertCoinBalanceToBaseUnit(v.Coin, balance, -1)
		// compare
		if value, exist := common.PorCoinDataMap[fmt.Sprintf("%s:%s", v.Coin, v.Address)]; exist {
			if isCoinBalanceEqual(balance, value.Balance) {
				log.Infof("verify coin %s, address %s balance success, in chain balance: %s, in por balance: %s", v.Coin, v.Address, balance, value.Balance)
			} else {
				log.Infof("verify coin %s, address %s balance failed, in chain balance:%s, in por balance:%s", v.Coin, v.Address, balance, value.Balance)
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func VerifyAllCoinAddressBalance(validator *common.AddressBalanceValidator) {
	for _, v := range common.PorCoinDataMap {
		// check coin black list
		if common.IsCheckBalanceBannedCoin(coin) {
			continue
		}
		balance, err := validator.GetCoinAddressBalanceInfo(strings.ToLower(v.Coin), v.Address, v.SnapshotHeight)
		if err != nil {
			log.Errorf("get address %s balance from blockchain failed, error: %v", v.Address, err)
			continue
		}
		balance = convertCoinBalanceToBaseUnit(v.Coin, balance, -1)
		// compare
		if isCoinBalanceEqual(balance, v.Balance) {
			log.Infof("verify coin %s, address %s balance success, in chain balance: %s, in por balance: %s", v.Coin, v.Address, balance, v.Balance)
		} else {
			log.Infof("verify coin %s, address %s balance failed, in chain balance: %s, in por balance: %s", v.Coin, v.Address, balance, v.Balance)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func VerifyCoinAddressTotalBalance(validator *common.AddressBalanceValidator, coin string) {
	totalBalance, totalPorBalance := decimal.NewFromInt(0), decimal.NewFromInt(0)
	coinAddressListMap, coinSnapshotHeightMap := make(map[string][]string), make(map[string]string)
	// get dest coins
	destCoins := getDestCoinList(coin)
	for _, value := range common.PorCoinDataMap {
		for index := range destCoins {
			if value.Coin == destCoins[index] {
				if value.Address == "" {
					continue
				}

				if coin == "BTC" && value.Script == "" {
					continue
				}

				if _, exist := coinAddressListMap[value.Coin]; exist {
					coinAddressListMap[value.Coin] = append(coinAddressListMap[value.Coin], value.Address)
				} else {
					coinAddressListMap[value.Coin] = []string{value.Address}
				}

				if _, exist := coinSnapshotHeightMap[value.Coin]; !exist {
					coinSnapshotHeightMap[value.Coin] = value.SnapshotHeight
				}
			}
		}
	}

	for coinTemp, addressList := range coinAddressListMap {
		height := coinSnapshotHeightMap[coinTemp]
		if len(addressList) == 0 {
			log.Errorf("no address to verify coin %s total balance", coin)
			return
		}
		amount, err := validator.GetCoinAddressTotalBalance(strings.ToLower(coinTemp), height, addressList)
		if err != nil {
			log.Errorf("get coin %s total address balance from blockchain failed, error: %v", coin, err)
			return
		}
		coinAmount, _ := decimal.NewFromString(amount)
		// USDC-OKC20/USDT-OKC20 precision is 18, convert to 6
		if coinTemp == "USDC-OKC20" || coinTemp == "USDT-OKC20" {
			coinAmountDecimal, _ := decimal.NewFromString(convertCoinBalanceToBaseUnit(coinTemp, coinAmount.String(), -1))
			totalBalance = totalBalance.Add(coinAmountDecimal.Mul(decimal.NewFromInt(1000000)))
		} else {
			totalBalance = totalBalance.Add(coinAmount)
		}

		porAmount := porCoinTotalBalance[coinTemp]
		totalPorBalance = totalPorBalance.Add(porAmount)

		log.Infof("coin %s, in chain balance: %s, in por balance: %s", coinTemp,
			convertCoinBalanceToBaseUnit(coinTemp, coinAmount.String(), -1), porAmount.String())
	}

	// convert coin balance to base unit
	dCoin := common.PorCoinUnitMap[strings.ToUpper(coin)]
	toalBalance := convertCoinBalanceToBaseUnit(dCoin, totalBalance.String(), -1)
	// compare
	if isCoinBalanceEqual(toalBalance, totalPorBalance.String()) {
		log.Infof("verify coin %s total address balance success, in chain balance: %s, in por balance: %s", coin, toalBalance, totalPorBalance.String())
	} else {
		log.Infof("verify coin %s total address balance failed, in chain balance: %s, in por balance: %s", coin, toalBalance, totalPorBalance.String())
	}
}

func VerifyAllCoinAddressTotalBalance(validator *common.AddressBalanceValidator) {
	for coin := range porCoinTotalBalance {
		VerifyCoinAddressTotalBalance(validator, coin)
	}
}

func convertCoinBalanceToBaseUnit(coin, amount string, round int32) string {
	precision, exist := common.PorCoinBaseUnitPrecisionMap[coin]
	if !exist {
		log.Errorf("coin %s not exist in base unit precision map", coin)
		precision = 8
	}
	amountDecimal, _ := decimal.NewFromString(amount)
	amountDecimal = amountDecimal.Div(decimal.NewFromInt(int64(math.Pow10(precision))))
	if round < 0 {
		return amountDecimal.String()
	} else {
		return amountDecimal.Round(round).String()
	}
}

func isCoinBalanceEqual(amount01, amount02 string) bool {
	amountDecimal01, _ := decimal.NewFromString(amount01)
	amountDecimal02, _ := decimal.NewFromString(amount02)
	return amountDecimal01.Equals(amountDecimal02)
}

func getDestCoinList(coin string) []string {
	coinList := make([]string, 0)
	destCoin := common.PorCoinUnitMap[strings.ToUpper(coin)]
	if strings.ToUpper(coin) != destCoin {
		coinList = append(coinList, coin)
	} else {
		for k, v := range common.PorCoinUnitMap {
			if v == destCoin {
				// check coin black list
				if common.IsCheckBalanceBannedCoin(k) {
					if _, exist := porCoinTotalBalance[k]; exist {
						log.Errorf("check balance not support the coin %s, ignore por total balance: %s", k, porCoinTotalBalance[k].String())
						continue
					}
				}
				coinList = append(coinList, k)
			}
		}
	}

	return coinList
}

func main() {
	Execute()
}
