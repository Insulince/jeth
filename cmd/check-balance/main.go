package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/Insulince/jeth/pkg/convert"
	"github.com/Insulince/jeth/pkg/eth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

func main() {
	ctx := context.Background()

	var address string

	flag.StringVar(&address, "address", "", "the wallet address whose balance you wish to check")
	flag.Parse()

	if address == "" {
		panic(errors.New("address cannot be blank, please provide a valid address via -address"))
	}

	client, err := ethclient.Dial(eth.DefaultGateway)
	if err != nil {
		panic(errors.Wrap(err, "dialing eth gateway"))
	}

	account := common.HexToAddress(address)
	balance, err := client.BalanceAt(ctx, account, eth.LatestBlock)
	if err != nil {
		panic(errors.Wrap(err, "fetching account balance"))
	}

	eth := convert.WeiIToEth(balance)
	fmt.Println(balance)
	fmt.Println(eth.String())
	usd := convert.EthToUsd(eth, 1536.31)
	fmt.Println(usd.String())
}
