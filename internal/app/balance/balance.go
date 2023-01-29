package balance

import (
	"context"
	"encoding/json"

	"github.com/aririfani/wallet-service/internal/app/topicinit"
	"github.com/aririfani/wallet-service/internal/app/wallet"
	"github.com/lovoo/goka"
)

var (
	group goka.Group = "balance"
	Table goka.Table = goka.GroupTable(group)
)

type GetBalanceCodec struct{}

func (c *GetBalanceCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (c *GetBalanceCodec) Decode(data []byte) (interface{}, error) {
	var w wallet.Wallet
	err := json.Unmarshal(data, &w)

	return w, err
}

func getBalance(ctx goka.Context, data interface{}) {
	w := wallet.Wallet{
		WalletID: data.(*wallet.Wallet).WalletID,
		Amount:   data.(*wallet.Wallet).Amount,
	}

	ctx.SetValue(w)
}

func PrepareTopics(brokers []string) {
	topicinit.EnsureStreamExists(string(wallet.WalletTopic), brokers)
}

func Run(ctx context.Context, brokers []string) func() error {
	return func() error {
		g := goka.DefineGroup(group,
			goka.Input(wallet.WalletTopic, new(wallet.WalletCodec), getBalance),
			goka.Persist(new(GetBalanceCodec)),
		)

		p, err := goka.NewProcessor(brokers, g)
		if err != nil {
			return err
		}

		return p.Run(ctx)
	}
}
