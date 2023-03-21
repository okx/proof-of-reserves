package common

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/schollz/progressbar/v3"

	//"github.com/gosuri/uiprogress"

	//"github.com/schollz/progressbar/v3"
	"os"
	"strconv"
	"strings"
)

func FindUserLeafNodesCountInMerkle(merkleJsonFileName string, leafNodes []TreeNode) (int, *Balances, *Balances, []*TreeNode, error) {

	if len(leafNodes) == 0 {
		return 0, nil, nil, nil, errors.New("no user info")
	}

	f, err := os.Open(merkleJsonFileName)
	if err != nil {
		return 0, nil, nil, nil, err
	}
	defer f.Close()

	count := 0
	reader := bufio.NewReader(f)
	var totalBalances = &Balances{
		BTC:  "0",
		ETH:  "0",
		USDT: "0",
	}

	var userTotalBalances = &Balances{
		BTC:  "0",
		ETH:  "0",
		USDT: "0",
	}
	var userLeafNodes = make([]*TreeNode, 0)
	lineCount := 0
	bar := progressbar.NewOptions(100, progressbar.OptionSetRenderBlankState(true))
	bar.Describe("validating")
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		detailInfo := strings.Split(string(line), ",")
		if len(detailInfo) < 3 {
			continue
		}
		// processor bar

		hash := detailInfo[0]
		height, _ := strconv.ParseInt(detailInfo[1], 10, 32)
		lineCount++
		if lineCount%300000 == 0 && lineCount/300000 < 100 {
			bar.Add(1)
		}
		if height == 1 {
			start := len(detailInfo[0]) + len(detailInfo[1]) + 2*len(",")
			balancesJson := line[start:len(line)]
			var balancesObj Balances
			json.Unmarshal([]byte(balancesJson), &balancesObj)
			totalBalances = totalBalances.Add(&balancesObj)

			ok, leafNode := isLeafNodeHash(hash, leafNodes)
			if !ok {
				continue
			}
			var curNode = TreeNode{
				Height:   int(height),
				Type:     2,
				Hash:     hash,
				Balances: &balancesObj,
			}

			if curNode.Equal(leafNode) {
				userTotalBalances = userTotalBalances.Add(&balancesObj)
				userLeafNodes = append(userLeafNodes, &curNode)
				count++
			}
		}

	}
	bar.Set(100)
	return count, userTotalBalances, totalBalances, userLeafNodes, nil
}

func isLeafNodeHash(curNodeHash string, leafNodes []TreeNode) (bool, TreeNode) {
	for _, leafNode := range leafNodes {
		if strings.EqualFold(curNodeHash, leafNode.Hash) {
			return true, leafNode
		}
	}

	return false, TreeNode{}
}

func NotFoundLeafNode(userNodes []TreeNode, findNodes []*TreeNode) []*TreeNode {
	result := make([]*TreeNode, 0)

	for _, userNode := range userNodes {
		isFind := false
		for _, findNode := range findNodes {
			if userNode.Hash == findNode.Hash {
				isFind = true
				break
			}
		}

		if !isFind {
			result = append(result, &userNode)
		}
	}

	return result
}
