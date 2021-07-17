package convert

import (
	"math/big"
)

const (
	weiPerEth = 1000000000000000000.0
	ethPerWei = 1.0 / weiPerEth

	gweiPerEth = 1000000000.0
	ethPerGwei = 1.0 / gweiPerEth

	weiPerGwei = 1000000000.0
	gweiPerWei = 1.0 / weiPerGwei
)

var (
	WeiPerEth = big.NewFloat(weiPerEth)
	EthPerWei = big.NewFloat(ethPerWei)

	GweiPerEth = big.NewFloat(gweiPerEth)
	EthPerGwei = big.NewFloat(ethPerGwei)

	WeiPerGwei = big.NewFloat(weiPerGwei)
	GweiPerWei = big.NewFloat(gweiPerWei)
)

// Big type converters

func Ftoi(f *big.Float) (i *big.Int) {
	i, _ = f.Int(nil)
	return i
}

func Itof(i *big.Int) (f *big.Float) {
	return new(big.Float).SetInt(i)
}

func F(bF *big.Float) (f float64) {
	v, _ := bF.Float64()
	return v
}

func I(bI *big.Int) (i int64) {
	return bI.Int64()
}

// USD converters

func EthToUsd(eth *big.Float, usdPerEth float64) (usd *big.Float) {
	return new(big.Float).Mul(eth, big.NewFloat(usdPerEth))
}

func UsdToEth(usd *big.Float, usdPerEth float64) (eth *big.Float) {
	ethPerUsd := 1.0 / usdPerEth
	return new(big.Float).Mul(usd, big.NewFloat(ethPerUsd))
}

func WeiToUsd(wei *big.Float, usdPerEth float64) (usd *big.Float) {
	eth := WeiToEth(wei)
	return EthToUsd(eth, usdPerEth)
}

func UsdToWei(usd *big.Float, usdPerEth float64) (wei *big.Float) {
	ethPerUsd := 1.0 / usdPerEth
	weiPerUsd, _ := EthToWei(big.NewFloat(ethPerUsd)).Float64()
	return new(big.Float).Mul(usd, big.NewFloat(weiPerUsd))
}

// ETH converters

func EthToWei(eth *big.Float) (wei *big.Float) {
	return new(big.Float).Mul(eth, WeiPerEth)
}

func WeiToEth(wei *big.Float) (eth *big.Float) {
	return new(big.Float).Mul(wei, EthPerWei)
}

// Helper functions to deal with big type converters and currency converters together

func EthToWeiI(eth *big.Float) (wei *big.Int) {
	return Ftoi(EthToWei(eth))
}

func WeiIToEth(wei *big.Int) (eth *big.Float) {
	return WeiToEth(Itof(wei))
}

func WeiIToUsd(wei *big.Int, usdPerEth float64) (usd *big.Float) {
	return WeiToUsd(Itof(wei), usdPerEth)
}

func UsdToWeiI(usd *big.Float, usdPerEth float64) (wei *big.Int) {
	return Ftoi(UsdToWei(usd, usdPerEth))
}
