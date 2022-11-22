package common

import (
	"crypto/sha256"
	"fmt"
	"github.com/shopspring/decimal"
)

type Balances struct {
	BTC  string `json:"BTC"`
	ETH  string `json:"ETH"`
	USDT string `json:"USDT"`
}

func (b *Balances) String() string {
	return fmt.Sprintf("{\"BTC\":\"%s\",\"ETH\":\"%s\",\"USDT\":\"%s\"}", b.BTC, b.ETH, b.USDT)
}

func (b *Balances) Equal(b2 *Balances) bool {
	return b.ETH == b2.ETH && b.USDT == b2.USDT && b.BTC == b2.BTC
}

func (b *Balances) Add(b2 *Balances) *Balances {
	btc1, err := decimal.NewFromString(b.BTC)
	if err != nil {
		return nil
	}
	eth1, err := decimal.NewFromString(b.ETH)
	if err != nil {
		return nil
	}
	usdt1, err := decimal.NewFromString(b.USDT)
	if err != nil {
		return nil
	}
	btc2, err := decimal.NewFromString(b2.BTC)
	if err != nil {
		return nil
	}
	eth2, err := decimal.NewFromString(b2.ETH)
	if err != nil {
		return nil
	}
	usdt2, err := decimal.NewFromString(b2.USDT)
	if err != nil {
		return nil
	}
	return &Balances{BTC: fmt.Sprintf("%s", btc1.Add(btc2).RoundDown(8).String()),
		ETH:  fmt.Sprintf("%s", eth1.Add(eth2).RoundDown(8).String()),
		USDT: fmt.Sprintf("%s", usdt1.Add(usdt2).RoundDown(8).String())}
}

func (b *Balances) Validate() bool {
	if _, err := decimal.NewFromString(b.BTC); err != nil {
		return false
	}
	if _, err := decimal.NewFromString(b.ETH); err != nil {
		return false
	}

	if _, err := decimal.NewFromString(b.USDT); err != nil {
		return false
	}
	return true
}

type Path struct {
	Height   int       `json:"height"`
	Type     int       `json:"type"`
	Hash     string    `json:"hash"`
	Balances *Balances `json:"balances"`
}

func (p *Path) Validate() bool {
	if p.Balances == nil || !p.Balances.Validate() || len(p.Hash) == 0 || p.Type < 1 || p.Type > 4 {
		return false
	}
	if p.Height < 1 {
		return false
	}
	return true
}

type Self struct {
	Nonce    string    `json:"nonce"`
	Hash     string    `json:"hash"`
	Balances *Balances `json:"balances"`
	Type     int       `json:"type"`
	Height   int       `json:"height"`
}

func (s *Self) Validate() bool {
	if s == nil {
		return false
	}
	data := fmt.Sprintf("%s%s", s.Nonce, s.Balances.String())
	hash := sha256.New()
	hash.Write([]byte(data))
	res := Encode(hash.Sum(nil))
	return res == "0x"+s.Hash
}

type MerkleProof struct {
	Self *Self   `json:"self"`
	Path []*Path `json:"path"`
}

func NewPath(lHash, rHash string, balance1, balance2 *Balances, height int) *Path {
	if !balance1.Validate() || !balance2.Validate() {
		return nil
	}
	balance := balance1.Add(balance2)
	data := fmt.Sprintf("%s%s%s%s%s%d", lHash, rHash, balance.BTC, balance.ETH, balance.USDT, height)
	hash := sha256.New()
	hash.Write([]byte(data))
	res := Encode(hash.Sum(nil))
	return &Path{Balances: balance, Height: height, Hash: res[2:]}
}

func (m *MerkleProof) Validate() bool {
	if m.Self == nil || len(m.Path) == 0 {
		return false
	}
	if !m.Self.Validate() {
		return false
	}
	if !m.Path[0].Validate() {
		return false
	}
	if m.Path[0].Type == m.Self.Type {
		return false
	}
	var left, right string
	if m.Self.Type == 1 {
		left, right = m.Self.Hash, m.Path[0].Hash
	} else {
		left, right = m.Path[0].Hash, m.Self.Hash
	}
	height := 1
	node := NewPath(left, right, m.Self.Balances, m.Path[0].Balances, height+1)
	for i := 1; i < len(m.Path)-1; i++ {
		if !m.Path[i].Validate() {
			return false
		}
		if m.Path[i].Type == 1 {
			left, right = m.Path[i].Hash, node.Hash
		} else {
			left, right = node.Hash, m.Path[i].Hash
		}
		node = NewPath(left, right, node.Balances, m.Path[i].Balances, m.Path[i].Height+1)
	}
	root := m.Path[len(m.Path)-1]
	fmt.Println(fmt.Sprintf("root BTC balance : %s ,root BTC balance in file: %s ", node.Balances.BTC, root.Balances.BTC))
	fmt.Println(fmt.Sprintf("root ETH balance : %s ,root ETH balance in file: %s ", node.Balances.ETH, root.Balances.ETH))
	fmt.Println(fmt.Sprintf("root USDT balance : %s ,root USDT balance in file: %s ", node.Balances.USDT, root.Balances.USDT))
	fmt.Println(fmt.Sprintf("root hash: %s ,root hash in file: %s ", node.Hash, root.Hash))
	if node.Hash != root.Hash || !node.Balances.Equal(root.Balances) || node.Height != root.Height {
		return false
	}
	return true
}
