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

type CoinData struct {
	Coin           string
	Network        string
	SnapshotHeight string
	Address        string
	Balance        string
	Message        string
	Sign1          string
	Sign2          string
	Script         string
}

func InitPorCsvDataMap(fileName string) (coinData map[string]*CoinData, err error) {
	fs, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("can not open the file, err is %+v", err)
		return
	}
	defer fs.Close()
	coinData = make(map[string]*CoinData)

	buf := bufio.NewReader(fs)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		if strings.Contains(string(line), "coin,") {
			continue
		}

		args := strings.Split(string(line), ",")
		if len(args) == 2 {
			continue
		} else if len(args) == 9 {
			// coin,network,snapshot height,address,balance,message,signature1,signature2,redeem_script
			d := &CoinData{
				Coin:           strings.ToUpper(cleanout(args[0])),
				Network:        strings.ToUpper(cleanout(args[1])),
				SnapshotHeight: cleanout(args[2]),
				Address:        cleanout(args[3]),
				Balance:        cleanout(args[4]),
				Message:        cleanout(args[5]),
				Sign1:          cleanout(args[6]),
				Sign2:          cleanout(args[7]),
				Script:         cleanout(args[8]),
			}
			if _, exist := coinData[fmt.Sprintf("%s:%s", d.Coin, d.Address)]; !exist {
				coinData[fmt.Sprintf("%s:%s", d.Coin, d.Address)] = d
			}
		}
	}

	return coinData, nil
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
