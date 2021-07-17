package convert

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Ftoi(t *testing.T) {
	tests := []struct {
		f *big.Float
		i *big.Int
	}{
		{f: big.NewFloat(1), i: big.NewInt(1)},
		{f: big.NewFloat(10), i: big.NewInt(10)},
		{f: big.NewFloat(0.1), i: big.NewInt(0)},
		{f: big.NewFloat(1.5), i: big.NewInt(1)},
		{f: big.NewFloat(8.99), i: big.NewInt(8)},
		{f: big.NewFloat(123456789.987654321), i: big.NewInt(123456789)},
	}

	for _, test := range tests {
		t.Run(test.f.String(), func(t *testing.T) {
			i := Ftoi(test.f)
			assert.Equal(t, 0, i.Cmp(test.i))
		})
	}
}

func Test_Itof(t *testing.T) {
	tests := []struct {
		i *big.Int
		f *big.Float
	}{
		{i: big.NewInt(1), f: big.NewFloat(1)},
		{i: big.NewInt(10), f: big.NewFloat(10)},
		{i: big.NewInt(0), f: big.NewFloat(0)},
		{i: big.NewInt(8), f: big.NewFloat(8)},
		{i: big.NewInt(123456789), f: big.NewFloat(123456789)},
	}

	for _, test := range tests {
		t.Run(test.i.String(), func(t *testing.T) {
			f := Itof(test.i)
			assert.Equal(t, 0, f.Cmp(test.f))
		})
	}
}

func Test_EthToWeiI(t *testing.T) {
	tests := []struct {
		eth       *big.Float
		weiString string
	}{
		{eth: big.NewFloat(1), weiString: "1000000000000000000"},
		{eth: big.NewFloat(0.000000000000000001), weiString: "1"},
		{eth: big.NewFloat(10000000), weiString: "10000000000000000000000000"},
		{eth: big.NewFloat(0.000000000123456789), weiString: "123456789"},
		{eth: big.NewFloat(2134879871723987598723987), weiString: "2134879871723987598723987000000000000000000"},
		{eth: big.NewFloat(1000000000000000000000000000000000000000000), weiString: "1000000000000000000000000000000000000000000000000000000000000"},
	}

	for _, test := range tests {
		t.Run(test.eth.String(), func(t *testing.T) {
			twei, _ := new(big.Int).SetString(test.weiString, 10)
			wei := EthToWeiI(test.eth)

			// It is not sufficient to simply check that the received and expected
			// are the same values, for they could differ due to float64 imprecision.
			// We need to check if the reason for the difference is because of an actual
			// error or just this float imprecision.
			// Thus we check if the two values are within 12 significant digits of each other.
			// This means the value differ by less than 0.000000000001%
			assert.True(t, InDelta(Itof(wei), Itof(twei), 1e-12))
		})
	}
}

func Test_WeiIToEth(t *testing.T) {
	tests := []struct {
		weiString string
		eth       *big.Float
	}{
		{weiString: "1000000000000000000", eth: big.NewFloat(1)},
		{weiString: "1", eth: big.NewFloat(0.000000000000000001)},
		{weiString: "10000000000000000000000000", eth: big.NewFloat(10000000)},
		{weiString: "123456789", eth: big.NewFloat(0.000000000123456789)},
	}

	for _, test := range tests {
		t.Run(test.weiString, func(t *testing.T) {
			wei, _ := new(big.Int).SetString(test.weiString, 10)
			eth := WeiIToEth(wei)
			tethVal, _ := test.eth.Float64()
			ethVal, _ := eth.Float64()
			assert.InDelta(t, tethVal, ethVal, 0.00000001)
		})
	}
}

func InDelta(a, b *big.Float, delta float64) bool {
	diff := new(big.Float).Sub(a, b)
	diff.Abs(diff)

	a.Abs(a)
	b.Abs(b)

	max := a
	if a.Cmp(b) <= 0 {
		max = b
	}

	maxDiff := new(big.Float).Mul(big.NewFloat(delta), max)

	return diff.Cmp(maxDiff) <= 0
}
