package eth

import (
	"math/big"
)

const (
	DefaultGateway = "https://cloudflare-eth.com"
)

var (
	LatestBlock *big.Int = nil
)

func ObfuscateKey(key string) string {
	out := ""
	numStars := len(key[:len(key)-4])
	for i := 0; i < numStars; i++ {
		out += "*"
	}
	out += key[len(key)-4:]
	return out
}
