package main

import (
	"bufio"
	"fmt"
	"github.com/okx/proof-of-reserves/common"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

var (
	cfgFile, csvFileName string
	coinTotalBalance     = make(map[string]decimal.Decimal)
)

var rootCmd = &cobra.Command{
	Use:   "AddressVerify",
	Short: "Verify address signature",
	Long:  ``,
	Run:   AddressVerify,
}

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

func parseLine(line string) []string {
	var result []string
	var part string
	inBrackets := false

	for i := 0; i < len(line); i++ {
		char := line[i]

		switch char {
		case '[':
			inBrackets = true
			part += string(char)
		case ']':
			inBrackets = false
			part += string(char)
		case ',':
			if inBrackets {
				part += string(char)
			} else {
				result = append(result, part)
				part = ""
			}
		default:
			part += string(char)
		}
	}

	result = append(result, part)

	return result
}

// detectFormatOffset inspects a detail-section header and returns the column offset for the
// report format: 1 when a "Type" column is present (the 12-column format that inserts a Type
// column right after coin), otherwise 0 for the legacy 11-column format. Every column after
// coin shifts right by this offset; a missing Type column means the section is parsed exactly
// as before and each row is treated as non-staking.
func detectFormatOffset(header []string) int {
	for _, col := range header {
		if strings.EqualFold(strings.TrimSpace(col), "Type") {
			return 1
		}
	}
	return 0
}

func handle(i int, line string, off int) (coin string, success bool) {
	if len(line) == 0 {
		return "", true
	}
	as := parseLine(line)
	if len(as) < 9+off {
		fmt.Println(fmt.Sprintf("Fail to verify address signature.The line %d has fewer columns than the report header.", i+1))
		return "", false
	}
	coin, addr, balance, message, sign1, sign2, script := as[0], as[3+off], as[4+off], as[5+off], as[6+off], as[7+off], as[8+off]

	// Recover from panic and print detailed error info
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("PANIC at line %d: %v | coin=%s, addr=%s, balance=%s, message=%s, sign1=%s, sign2=%s, script=%s\n", i+1, r, coin, addr, balance, message, sign1, sign2, script)
			success = false
		}
	}()
	var eoa1, eoa2 string
	if len(as) > 10+off {
		eoa1 = as[9+off]
		eoa2 = as[10+off]
	}

	val, err := decimal.NewFromString(balance)
	if err != nil {
		fmt.Println(fmt.Sprintf("Fail to verify address signature.The line %d  has invalid balance number.", i+1))
		return coin, false
	}

	coin = strings.ToUpper(coin)
	totalCoin, exist := common.PorCoinUnitMap[coin]
	if !exist {
		fmt.Println(fmt.Sprintf("Fail to verify address signature.The line %d  has invalid coin name, %s", i+1, coin))
		return coin, false
	}
	_, exist = coinTotalBalance[totalCoin]
	if exist {
		coinTotalBalance[totalCoin] = coinTotalBalance[totalCoin].Add(val)
	} else {
		coinTotalBalance[totalCoin] = val
	}

	if common.IsVerifyAddressBannedCoin(coin) {
		return coin, true
	}

	if addr == "" || message == "" || sign1 == "" {
		fmt.Println(fmt.Sprintf("Fail to verify address signature.The line %d is missing some parameters. coin:%s, addr: %s", i+1, coin, addr))
		return coin, false
	}

	switch common.PorCoinTypeMap[coin] {
	case common.EvmCoinTye:
		if eoa1 != "" && eoa2 != "" {
			if err := common.VerifyEvmCoin(coin, eoa1, message, sign1); err != nil {
				fmt.Println(fmt.Sprintf("Fail to verify address %s signature.The line %d  has error:%s.", addr, i+1, err))
				return coin, false
			}
			if err := common.VerifyEvmCoin(coin, eoa2, message, sign2); err != nil {
				fmt.Println(fmt.Sprintf("Fail to verify address %s signature.The line %d  has error:%s.", addr, i+1, err))
				return coin, false
			}
		} else if eoa1 != "" {
			if err := common.VerifyEvmCoin(coin, eoa1, message, sign1); err != nil {
				fmt.Println(fmt.Sprintf("Fail to verify address %s signature.The line %d  has error:%s.", eoa1, i+1, err))
				return coin, false
			}
		} else {
			if err := common.VerifyEvmCoin(coin, addr, message, sign1); err != nil {
				fmt.Println(fmt.Sprintf("Fail to verify address %s signature.The line %d  has error:%s.", addr, i+1, err))
				return coin, false
			}
		}
	case common.EcdsaCoinType:
		if eoa1 != "" && eoa2 != "" {
			if err := common.VerifyEcdsaCoin(coin, eoa1, message, sign1); err != nil {
				fmt.Println(fmt.Sprintf("Fail to verify address %s signature.The line %d  has error:%s.", addr, i+1, err))
				return coin, false
			}
			if err := common.VerifyEcdsaCoin(coin, eoa2, message, sign2); err != nil {
				fmt.Println(fmt.Sprintf("Fail to verify address %s signature.The line %d  has error:%s.", addr, i+1, err))
				return coin, false
			}
		} else if eoa1 != "" {
			if err := common.VerifyEcdsaCoin(coin, eoa1, message, sign1); err != nil {
				fmt.Println(fmt.Sprintf("Fail to verify address %s signature.The line %d  has error:%s.", addr, i+1, err))
				return coin, false
			}
		} else {
			if err := common.VerifyEcdsaCoin(coin, addr, message, sign1); err != nil {
				fmt.Println(fmt.Sprintf("Fail to verify address %s signature.The line %d  has error:%s.", addr, i+1, err))
				return coin, false
			}
		}
	case common.Ed25519CoinType:
		// owner mode: when eoa1 (current authentication key, e.g. a rotated APTOS account)
		// is present, verify against eoa1 instead of the claimed address (mirrors EVM owner mode).
		if eoa1 != "" {
			if err := common.VerifyEd25519Coin(coin, eoa1, message, sign1, script); err != nil {
				fmt.Println(fmt.Sprintf("Fail to verify address %s signature.The line %d  has error:%s.", eoa1, i+1, err))
				return coin, false
			}
		} else {
			if err := common.VerifyEd25519Coin(coin, addr, message, sign1, script); err != nil {
				fmt.Println(fmt.Sprintf("Fail to verify address %s signature.The line %d  has error:%s.", addr, i+1, err))
				return coin, false
			}
		}
	case common.TrxCoinType:
		if err := common.VerifyTRX(addr, message, sign1); err != nil {
			fmt.Println(fmt.Sprintf("Fail to verify address %s signature.The line %d  has error:%s.", addr, i+1, err))
			return coin, false
		}
	case common.BethCoinType:
		if err := common.VerifyBETH(addr, message, sign1); err != nil {
			fmt.Println(fmt.Sprintf("Fail to verify address %s signature.The line %d  has error:%s.", addr, i+1, err))
			return coin, false
		}
	case common.UTXOCoinType:
		if err := common.VerifyUtxoCoin(coin, addr, message, sign1, sign2, script); err != nil {
			fmt.Println(fmt.Sprintf("Fail to verify address %s signature.The line %d  has error:%s.", addr, i+1, err))
			return coin, false
		}
	case common.StarkCoinType:
		if err := common.VerifyStarkCoin(coin, addr, message, sign1, script); err != nil {
			fmt.Println(fmt.Sprintf("Fail to verify address %s signature.The line %d  has error:%s.", addr, i+1, err))
			return coin, false
		}
	default:
		fmt.Println(fmt.Sprintf("Fail to verify address %s signature. Invaild coin type:%s", addr, coin))
		return coin, false
	}
	return coin, true
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
	count, lineSize, flag := 0, 0, 2
	off := 0
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
			fmt.Println(fmt.Sprintf("%s's total balance is %s.", as[0], as[1]))
		}

		if flag > 0 {
			if flag == 1 {
				flag--
				// Resolve the detail-section column layout from its header.
				off = detectFormatOffset(strings.Split(temp, ","))
			}
			continue
		}
		coin, ok := handle(count, strings.TrimSpace(temp), off)

		if common.IsVerifyAddressBannedCoin(coin) {
			continue
		}

		if _, exist := fail[coin]; !exist {
			fail[coin] = 0
		}

		if _, exist := success[coin]; !exist {
			success[coin] = 0
		}

		if !ok {
			fail[coin]++
		} else {
			success[coin]++
		}
		count++
	}
	if count == 0 {
		fmt.Println("Verify address signature end.The file is empty.")
	}
	var allPass = true
	for k, v := range success {
		fmt.Println(fmt.Sprintf("%s  %d accoounts, %d verified, %d failed", k, v+fail[k], v, fail[k]))
		if fail[k] != 0 {
			allPass = false
		}
	}

	coinTotalBalanceResult := make([]string, 0)
	for coin, balance := range coinTotalBalance {
		coinTotalBalanceResult = append(coinTotalBalanceResult, fmt.Sprintf("%s(%s)", coin, balance.Round(2).String()))
	}
	fmt.Printf("Total balance: [%s]\n", strings.Join(coinTotalBalanceResult, ","))

	if allPass {
		fmt.Println("Verify address signature end, all address passed")
	}
}

func main() {
	Execute()
}
