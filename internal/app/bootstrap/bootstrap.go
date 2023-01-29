package bootstrap

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aririfani/wallet-service/internal/app/balance"
	"golang.org/x/sync/errgroup"
)

var (
	brokers   = []string{"localhost:9092"}
	runWorker = flag.Bool("collector", true, "run collector processor")
)

func Run() {
	flag.Parse()
	ctx, cancle := context.WithCancel(context.Background())
	grp, ctx := errgroup.WithContext(ctx)

	if *runWorker {
		balance.PrepareTopics(brokers)
	}

	if *runWorker {
		log.Println("starting worker")
		grp.Go(balance.Run(ctx, brokers))
	}

	// wait for SIGINT/SIGTERM
	waiter := make(chan os.Signal, 1)
	signal.Notify(waiter, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-waiter:
	case <-ctx.Done():
	}

	cancle()
	if err := grp.Wait(); err != nil {
		log.Println(err)
	}
	log.Println("done")
}
