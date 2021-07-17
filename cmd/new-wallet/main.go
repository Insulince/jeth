package main

import (
	"fmt"
	"time"

	"github.com/Insulince/jeth/pkg/wallet"

	"github.com/pkg/errors"
)

func main() {
	fmt.Printf("Generating a new Ethereum private key, public key, and wallet address...\n")
	start := time.Now()

	w := wallet.New()

	if err := w.Validate(); err != nil {
		panic(errors.Wrap(err, "generated wallet is invalid"))
	}

	duration := time.Since(start)
	fmt.Printf("\nPRIVATE KEY:\n%s\n\nPUBLIC KEY:\n%s\n\nWALLET ADDRESS:\n%s\n\nWallet has been validated and is well-formed and correct.\n\nSuccess! (%v)\n", w.PrivateKeyHex(), w.PublicKeyHex(), w.Address(), duration)
}
