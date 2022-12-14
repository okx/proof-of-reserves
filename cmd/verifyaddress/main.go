package main

import (
	"bufio"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/okx/proof-of-reserves/common"
	"github.com/spf13/cobra"
	"io"
	"math/big"
	"os"
	"strings"
)

var cfgFile, csvFileName string

var rootCmd = &cobra.Command{
	Use:   "AddressVerify",
	Short: "Verify address signature",
	Long:  ``,
	Run:   AddressVerify,
}
var (
	zero = big.NewInt(0)
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&csvFileName, "por_csv_filename", "", "")
}

func initConfig() {}

var (
	coinMap = map[string]string{"USDT-ERC20": "ETH", "USDT-TRC20": "TRX", "USDT-OMNI": "BTC", "BTC": "BTC", "ETH": "ETH",
		"USDT-POLY": "ETH", "USDT-AVAXC": "ETH", "USDT-ARBITRUM": "ETH", "ETH-ARBITRUM": "ETH", "ETH-OPTIMISM": "ETH", "USDT-OPTIMISM": "ETH"}
)

func divideInt(coin string, i *big.Int) *big.Int {
	switch {
	case strings.HasPrefix(coin, "ETH"):
		return i.Div(i, big.NewInt(1e18))
	case strings.HasPrefix(coin, "USDT"):
		return i.Div(i, big.NewInt(1e6))
	case coin == "BTC":
		return i.Div(i, big.NewInt(1e8))
	default:
		panic("未知")

	}
}

func handle(i int, line string) (string, *big.Int, bool) {
	if len(line) == 0 {
		return "", zero, true
	}
	as := strings.Split(line, ",")
	coin, addr, balance, message, sign1, sign2, script := as[0], as[1], as[2], as[3], as[4], as[5], as[6]
	v := big.NewInt(0)
	val, ok := v.SetString(balance, 10)
	if !ok {
		fmt.Println(fmt.Sprintf("Fail to verify address signature.The line %d  has invalid balance number.", i+1))
		return coin, zero, false
	}
	switch coinMap[coin] {
	case "ETH":
		if err := common.VerifyETH(addr, message, sign1); err != nil {
			fmt.Println(fmt.Sprintf("Fail to verify address signature.The line %d  has error:%s.", i+1, err))
			return coin, zero, false
		}
	case "SOL":
		if err := common.VerifySol(addr, message, sign1); err != nil {
			fmt.Println(fmt.Sprintf("Fail to verify address signature.The line %d  has error:%s.", i+1, err))
			return coin, zero, false
		}
	case "TRX":
		if err := common.VerifyTRX(addr, message, sign1); err != nil {
			fmt.Println(fmt.Sprintf("Fail to verify address signature.The line %d  has error:%s.", i+1, err))
			return coin, zero, false
		}
	case "BTC":
		if len(sign2) == 0 || sign2 == "null" {
			if err := common.VerifyBTC(addr, message, sign1); err != nil {
				fmt.Println(fmt.Sprintf("Fail to verify address signature.The line %d  has error:%s.", i+1, err))
				return coin, zero, false
			}
		} else if strings.HasPrefix(addr, "bc1") {
			if !common.VerifyBTCWitness(addr, message, script, sign1, sign2) {
				fmt.Println(fmt.Sprintf("Fail to verify address signature.The line %d  has error.", i+1))
				return coin, zero, false
			}
		} else {
			addrPub, err := btcutil.NewAddressScriptHash(common.MustDecode(script), &chaincfg.MainNetParams)
			if err != nil {
				fmt.Println(fmt.Sprintf("Fail to verify address signature.The line %d  has error:%s.", i+1, err))
				return coin, zero, false
			}
			if addrPub.EncodeAddress() != addr {
				fmt.Println(fmt.Sprintf("Fail to verify address signature.The line %d  has error:%s.", i+1, err))
				return coin, zero, false
			}
			addr1, err := common.SigToAddrBTC(message, sign1)
			if err != nil {
				fmt.Println(fmt.Sprintf("Fail to verify address signature.The line %d  has error:%s.", i+1, err))
				return coin, zero, false
			}
			addr2, err := common.SigToAddrBTC(message, sign2)
			if err != nil {
				fmt.Println(fmt.Sprintf("Fail to verify address signature.The line %d  has error:%s.", i+1, err))
				return coin, zero, false
			}
			typ, addrs, _, err := txscript.ExtractPkScriptAddrs(common.MustDecode(script), &chaincfg.MainNetParams)
			if typ != txscript.MultiSigTy {
				fmt.Println(fmt.Sprintf("Fail to verify address signature.The line %d  has error:%s.", i+1, err))
				return coin, zero, false
			}
			if err != nil {
				fmt.Println(fmt.Sprintf("Fail to verify address signature.The line %d  has error:%s.", i+1, err))
				return coin, zero, false
			}
			if len(addrs) != 3 {
				fmt.Println(fmt.Sprintf("Fail to verify address signature.The line %d  has invalid address.", i+1))
				return coin, zero, false
			}
			m := map[string]struct{}{addr1: {}, addr: {}, addr2: {}}
			for _, v := range addrs {
				delete(m, v.EncodeAddress())
			}
			if len(m) > 1 {
				fmt.Println(fmt.Sprintf("Fail to verify address signature.The line %d  has invalid address.", i+1))
				return coin, zero, false
			}
		}
	default:
		fmt.Println("Fail to verify address signature,invalid coin type.")
		return coin, zero, false
	}
	return coin, val, true
}

func AddressVerify(cmd *cobra.Command, args []string) {
	fmt.Println("Verify address signature start")
	fmt.Println("Your input csv filename: " + csvFileName)
	f, err := os.Open(csvFileName)
	defer f.Close()
	if err != nil {
		fmt.Println("Fail to verify address signature.The error is ", err)
		return
	}
	buf := bufio.NewReader(f)
	btc, eth, usdt := big.NewInt(0), big.NewInt(0), big.NewInt(0)
	count, lineSize, flag := 0, 0, 2
	success, fail := make(map[string]uint64), make(map[string]uint64)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Fail to verify address signature.The error is ", err)
			return
		}
		lineSize++
		temp := string(line)
		if temp == "" {
			flag--
			continue
		}
		if lineSize == 1 {
			continue
		}
		as := strings.Split(temp, ",")
		if flag > 1 {
			fmt.Println(fmt.Sprintf("%s 's height is %s and total balance is %s.", as[0], as[1], as[2]))
		}

		if flag > 0 {
			if flag == 1 {
				flag--
			}
			continue
		}
		coin, val, ok := handle(count, strings.TrimSpace(temp))
		if !ok {
			fail[coin]++
		} else {
			success[coin]++
		}

		if strings.HasPrefix(coin, "USDT") {
			usdt = usdt.Add(usdt, val)
		} else if strings.HasPrefix(coin, "ETH") {
			eth = eth.Add(eth, val)
		} else if coin == "BTC" {
			btc = btc.Add(btc, val)
		}
		count++
	}
	if count == 0 {
		fmt.Println("Verify address signature end.The file is empty.")
	}
	for k, v := range success {
		fmt.Println(fmt.Sprintf("%s  %d accoounts, %d verified, %d failed", k, v+fail[k], v, fail[k]))
	}
	fmt.Println(fmt.Sprintf("Total balance :BTC %s,ETH(ALL) %s,USDT(ALL):%s", divideInt("BTC", btc).String(), divideInt("ETH", eth).String(), divideInt("USDT", usdt).String()))
	if len(fail) == 0 {
		fmt.Println("Verify address signature end, all address passed")
	}
}

func main() {
	Execute()
}
