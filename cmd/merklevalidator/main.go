package main

import (
	"encoding/json"
	"fmt"
	"github.com/okx/proof-of-reserves/common"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var cfgFile, merkleJsonFile string

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
	rootCmd.PersistentFlags().StringVar(&merkleJsonFile, "merkle_proof_file", "", "")
}

func initConfig() {}

func MerkleValidator(cmd *cobra.Command, args []string) {
	fmt.Println("merkle tree path  validation start")
	if merkleJsonFile == "" {
		fmt.Println("Merkle tree path validation failed, invalid merkle proof file")
		return
	}
	b, err := ioutil.ReadFile(merkleJsonFile)
	if err != nil {
		fmt.Println("Merkle tree path validation failed, invalid merkle proof file", err)
		return
	}
	if len(b) == 0 {
		fmt.Println("Merkle tree path validation failed, empty merkle proof file")
		return
	}
	var pf common.MerkleProof
	if err := json.Unmarshal(b, &pf); err != nil {
		fmt.Println(fmt.Sprintf("Merkle tree path validation failed, error:%s", err))
		return
	}
	if pf.Validate() {
		fmt.Println("Merkle tree path validation passed.")
	} else {
		fmt.Println("Merkle tree path validation failed.")
	}
}
func main() {
	Execute()
}
