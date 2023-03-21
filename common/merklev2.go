package common

//
import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
)

type TreeNode struct {
	Height   int       `json:"height"`
	Type     int       `json:"type"`
	Hash     string    `json:"hash"`
	Balances *Balances `json:"balances"`
}
type SelfInfo struct {
	Hash          string     `json:"hash"`
	Nodes         []TreeNode `json:"nodes"`
	Nonce         string     `json:"nonce"`
	TotalBalances Balances   `json:"totalBalances"`
}

func (s *SelfInfo) Check() bool {
	if len(s.Nodes) == 0 {
		return false
	}
	for _, p := range s.Nodes {
		data := fmt.Sprintf("%s%s%s%s", s.Hash, p.Balances.BTC, p.Balances.ETH, p.Balances.USDT)
		hash := sha256.New()
		hash.Write([]byte(data))
		if !strings.EqualFold(Encode(hash.Sum(nil))[2:], p.Hash) {
			return false
		}
	}

	balancesJson, err := json.Marshal(s.TotalBalances)
	if err != nil {
		return false
	}
	data := fmt.Sprintf("%s%s", s.Nonce, balancesJson)
	hash := sha256.New()
	hash.Write([]byte(data))
	if !strings.EqualFold(Encode(hash.Sum(nil))[2:], s.Hash) {
		return false
	}

	if s.Nodes == nil && len(s.Nodes) == 0 {
		return false
	}
	totalBalances := s.Nodes[0].Balances
	for i := 1; i < len(s.Nodes); i++ {
		totalBalances = totalBalances.Add(s.Nodes[i].Balances)
	}

	return totalBalances.Equal(&s.TotalBalances)
}

func (p *TreeNode) Check(left, right string) bool {
	if !p.Balances.Validate() {
		return false
	}

	data := fmt.Sprintf("%s%s%s%s%s%d", left, right, p.Balances.BTC, p.Balances.ETH, p.Balances.USDT, p.Height)
	hash := sha256.New()
	hash.Write([]byte(data))
	return Encode(hash.Sum(nil)) == p.Hash
}

func (p *TreeNode) Equal(node TreeNode) bool {
	if !p.Balances.Validate() {
		return false
	}
	return p.Hash == node.Hash && p.Balances.Equal(node.Balances)
}
func (p *TreeNode) CheckSelf() bool {
	if p.Balances == nil || !p.Balances.Validate() || len(p.Hash) == 0 || p.Type < 1 || p.Type > 4 {
		return false
	}
	if p.Height < 1 {
		return false
	}
	return true
}
