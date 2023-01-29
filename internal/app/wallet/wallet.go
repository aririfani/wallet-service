package wallet

import (
	"encoding/json"

	"github.com/lovoo/goka"
)

var (
	WalletTopic goka.Stream = "deposits"
)

type Wallet struct {
	WalletID string
	Amount   float64
}

type WalletCodec struct{}

// Encode ...
func (w *WalletCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

// Decode ...
func (w *WalletCodec) Decode(data []byte) (interface{}, error) {
	var wallet Wallet
	return &wallet, json.Unmarshal(data, &wallet)
}
