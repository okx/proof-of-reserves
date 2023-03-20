package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/okx/proof-of-reserves/common"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var cfgFile, merkleJsonFileV2, userInfoFile, merkleJsonFile string

var rootCmd = &cobra.Command{
	Use:   "MerkleValidator",
	Short: "merkle tree path  validation",
	Long:  ``,
	Run:   MerkleValidator,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&merkleJsonFileV2, "merkle_file", "", "")
	rootCmd.PersistentFlags().StringVar(&userInfoFile, "user_info_file", "", "")
	rootCmd.PersistentFlags().StringVar(&merkleJsonFile, "merkle_proof_file", "", "")
}

func initConfig() {}

func MerkleValidator(cmd *cobra.Command, args []string) {
	// v1 verify
	if len(userInfoFile) == 0 {
		if merkleJsonFile == "" {
			fmt.Println("Merkle tree path validation failed, invalid merkle proof file")
			return
		}
		passed := execV1(merkleJsonFile)
		if passed {
			fmt.Println("Merkle tree path validation passed.")
		} else {
			fmt.Println("Merkle tree path validation failed.")
		}
		return
	}

	// v2 verify
	if merkleJsonFileV2 == "" {
		fmt.Println("Merkle tree path validation failed, invalid merkle proof file")
		return
	}
	passed, userBalance, totalBalance, userLeafNodes, err := execV2(userInfoFile, merkleJsonFileV2)
	if err != nil {
		fmt.Println(fmt.Sprintf("Merkle tree path validation failed, %s", err))
		return
	}

	if passed {
		fmt.Println(fmt.Sprintf("Merkle tree path validation passed"))
		for _, leafNode := range userLeafNodes {
			fmt.Println(fmt.Sprintf("%s is found in the merkle tree. %s BTC , %s ETH , %s USDT", leafNode.Hash, leafNode.Balances.BTC, leafNode.Balances.ETH, leafNode.Balances.USDT))
		}
		fmt.Println(fmt.Sprintf("Your asset holdings: %s BTC , %s ETH , %s USDT", userBalance.BTC, userBalance.ETH, userBalance.USDT))
		fmt.Println(fmt.Sprintf("Total OKX user asset holdings:  %s BTC , %s ETH , %s USDT", totalBalance.BTC, totalBalance.ETH, totalBalance.USDT))
	} else {
		hash := ""
		for _, notFountLeafNode := range userLeafNodes {
			hash = hash + notFountLeafNode.Hash + ","
		}
		hash = hash[0 : len(hash)-1]
		fmt.Println(fmt.Sprintf("Merkle tree path validation failed, %s is not found in merkle tree or invalid data", hash)) //todo
	}
}
func main() {
	Execute()
}
func execV1(userInfoFile string) bool {
	b, err := ioutil.ReadFile(userInfoFile)
	if err != nil {
		fmt.Println("Merkle tree path validation failed, invalid merkle proof file", err)
		return false
	}
	if len(b) == 0 {
		fmt.Println("Merkle tree path validation failed, empty merkle proof file")
		return false
	}

	var pf common.MerkleProof
	if err := json.Unmarshal(b, &pf); err != nil {
		fmt.Println(fmt.Sprintf("Merkle tree path validation failed, error:%s", err))
		return false
	}
	return pf.Validate()
}

func execV2(userInfoFile string, merkleJsonFile string) (bool, *common.Balances, *common.Balances, []*common.TreeNode, error) {
	b, err := ioutil.ReadFile(userInfoFile)
	if err != nil {
		return false, nil, nil, nil, errors.New("user info path validation failed, invalid user info file")
	}
	if len(b) == 0 {
		return false, nil, nil, nil, errors.New("user info path validation failed, invalid user info file")
	}

	var selfInfo common.SelfInfo
	if err := json.Unmarshal(b, &selfInfo); err != nil {
		return false, nil, nil, nil, errors.New(fmt.Sprintf("provided data is invalid, error:%s", err))
	}
	if !selfInfo.Check() {
		return false, nil, nil, nil, errors.New("provided data is invalid")
	}
	count, userBalances, totalBalances, userLeafNodes, err := common.FindUserLeafNodesCountInMerkle(merkleJsonFile, selfInfo.Nodes)
	if err != nil {
		return false, nil, nil, nil, err
	}
	isFind := count == len(selfInfo.Nodes)
	if !isFind {
		notFoundLeafNode := common.NotFoundLeafNode(selfInfo.Nodes, userLeafNodes)
		return isFind, userBalances, totalBalances, notFoundLeafNode, nil
	}
	return isFind, userBalances, totalBalances, userLeafNodes, nil
}
