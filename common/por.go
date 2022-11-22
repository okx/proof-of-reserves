package common

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

var PorCoinDataMap map[string]*CoinData

var PorCoinGeneralDataMap map[string]*CoinGeneralInfo

type CoinData struct {
	Coin    string
	Address string
	Balance string
	Message string
	Sign1   string
	Sign2   string
	Script  string
}

type CoinGeneralInfo struct {
	Coin           string
	SnapshotHeight string
	Balance        string
	AddressCount   int64
}

func InitPorCsvDataMap(fileName string) (coinData map[string]*CoinData, coinGeneralInfo map[string]*CoinGeneralInfo, err error) {
	fs, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("can not open the file, err is %+v", err)
		return
	}
	defer fs.Close()
	coinData = make(map[string]*CoinData)
	coinGeneralInfo = make(map[string]*CoinGeneralInfo)

	buf := bufio.NewReader(fs)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		if strings.Contains(string(line), "coin,") || strings.Contains(string(line), "ETH(ALL),") || strings.Contains(string(line), "USDT(ALL),") {
			continue
		}

		args := strings.Split(string(line), ",")
		if len(args) == 3 {
			c := &CoinGeneralInfo{
				Coin:           strings.ToLower(cleanout(args[0])),
				SnapshotHeight: cleanout(args[1]),
				Balance:        cleanout(args[2]),
			}
			if _, exist := coinGeneralInfo[c.Coin]; !exist {
				coinGeneralInfo[c.Coin] = c
			}
		} else if len(args) == 7 {
			// coin,address,balance,message,signature1,signature2,redeem_script
			d := &CoinData{
				Coin:    strings.ToLower(cleanout(args[0])),
				Address: cleanout(args[1]),
				Balance: cleanout(args[2]),
				Message: cleanout(args[3]),
				Sign1:   cleanout(args[4]),
				Sign2:   cleanout(args[5]),
				Script:  cleanout(args[6]),
			}
			if _, exist := coinData[fmt.Sprintf("%s:%s", d.Coin, d.Address)]; !exist {
				coinData[fmt.Sprintf("%s:%s", d.Coin, d.Address)] = d
			}
		}
	}

	return coinData, coinGeneralInfo, nil
}

func cleanout(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}
