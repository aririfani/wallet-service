package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aririfani/wallet-service/internal/app/balance"
	"github.com/aririfani/wallet-service/internal/app/wallet"
	"github.com/gorilla/mux"
	"github.com/lovoo/goka"
)

func Run(brokers []string, stream goka.Stream) {
	view, err := goka.NewView(brokers, balance.Table, new(balance.GetBalanceCodec))
	if err != nil {
		panic(err)
	}

	go view.Run(context.Background())

	emitter, err := goka.NewEmitter(brokers, stream, new(wallet.WalletCodec))
	if err != nil {
		panic(err)
	}

	defer emitter.Finish()

	router := mux.NewRouter()
	router.HandleFunc("/deposit", deposit(emitter, stream)).Methods("POST")
	router.HandleFunc("/{wallet_id}/wallet", getBalance(view)).Methods("GET")

	log.Printf("Listen port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func deposit(emmiter *goka.Emitter, stream goka.Stream) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var wallet wallet.Wallet

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "error: %v", err)
			return
		}

		err = json.Unmarshal(b, &wallet)
		if err != nil {
			fmt.Fprintf(w, "error: %v", err)
			return
		}

		err = emmiter.EmitSync(wallet.WalletID, &wallet)
		if err != nil {
			fmt.Fprintf(w, "error: %v", err)
			return
		}

		log.Printf("Deposit send:\n %v\n", wallet)
		fmt.Println(w, "Deposit send:\n %v\n", wallet)
	}
}

func getBalance(view *goka.View) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walletID := mux.Vars(r)["wallet_id"]
		val, _ := view.Get(walletID)
		fmt.Println("view topic", val)

		if val == nil {
			fmt.Fprintf(w, "%s not found!", walletID)
			return
		}

		wall := val.(wallet.Wallet)
		fmt.Println(wall)
		w.Header().Add("Content-Type", "application/json")
	}
}
