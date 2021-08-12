package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

func main() {
	ctx := context.Background()

	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		panic(errors.Wrap(err, "dialing eth gateway"))
	}

	bGasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		panic(errors.Wrap(err, "getting suggested gas price"))
	}
	fmt.Println(bGasPrice.String())
}
