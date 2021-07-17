package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	// These are known to be valid according to MEW. Use them to test. Do NOT send funds to this address, it is not secure and was intentionally created to test with.
	privateKeyHex = "7cd7d434407526ad4c7a64d4f7d26a2a45bb0da1cc7406c166e1e3ddfcce03ed"
	publicKeyHex  = "ac7a41fcbb11cb057a1f0bf1710e7f0aab1c93a468c54f7472695b5b97c1af687b8a0692b1d37ff1d5e65d2ecd2f6befdf7d0c89a403a5bbafcd4a9143bb9de7"
	address       = "0x19325d2D5c17AF1096D28A12850D27bD182612F6"
)

func Test_New(t *testing.T) {
	// Cannot be unit tested in a meaningful way. Validate manually using offline MEW.
}

func Test_ManualHex(t *testing.T) {
	w := ManualHex(privateKeyHex, publicKeyHex, address)

	assert.Equal(t, privateKeyHex, w.PrivateKeyHex())
	assert.Equal(t, publicKeyHex, w.PublicKeyHex())
	assert.Equal(t, address, w.Address())
}

func Test_Wallet_FromPrivateKeyHex(t *testing.T) {
	w := FromPrivateKeyHex(privateKeyHex)

	assert.Equal(t, privateKeyHex, w.PrivateKeyHex())
	assert.Equal(t, publicKeyHex, w.PublicKeyHex())
	assert.Equal(t, address, w.Address())
}

func Test_FromPublicKeyHex(t *testing.T) {
	w := FromPublicKeyHex(privateKeyHex, publicKeyHex)

	assert.Equal(t, privateKeyHex, w.PrivateKeyHex())
	assert.Equal(t, publicKeyHex, w.PublicKeyHex())
	assert.Equal(t, address, w.Address())
}

func Test_Wallet_PrivateKeyHex(t *testing.T) {
	w := ManualHex(privateKeyHex, publicKeyHex, address)

	assert.Equal(t, privateKeyHex, w.PrivateKeyHex())
}

func Test_Wallet_PublicKeyHex(t *testing.T) {
	w := ManualHex(privateKeyHex, publicKeyHex, address)

	assert.Equal(t, publicKeyHex, w.PublicKeyHex())
}

func Test_Wallet_Address(t *testing.T) {
	w := ManualHex(privateKeyHex, publicKeyHex, address)

	assert.Equal(t, address, w.Address())
}

func Test_Wallet_Clone(t *testing.T) {
	w := ManualHex(privateKeyHex, publicKeyHex, address)
	w2 := w.Clone()

	// Assert w2 has the original values
	assert.Equal(t, privateKeyHex, w2.PrivateKeyHex())
	assert.Equal(t, publicKeyHex, w2.PublicKeyHex())
	assert.Equal(t, address, w2.Address())

	// Assert that w2 is equal to w
	assert.Equal(t, w.PrivateKeyHex(), w2.PrivateKeyHex())
	assert.Equal(t, w.PublicKeyHex(), w2.PublicKeyHex())
	assert.Equal(t, w.address, w2.Address())

	// Assert that w is equal to w2
	assert.Equal(t, w2.PrivateKeyHex(), w.PrivateKeyHex())
	assert.Equal(t, w2.PublicKeyHex(), w.PublicKeyHex())
	assert.Equal(t, w2.address, w.Address())
}

func Test_Wallet_Equals(t *testing.T) {
	w := ManualHex(privateKeyHex, publicKeyHex, address)
	w2 := ManualHex(privateKeyHex, publicKeyHex, address)

	err := w.Equals(w2)
	err2 := w2.Equals(w)

	assert.NoError(t, err)
	assert.NoError(t, err2)
}

func Test_Wallet_Validate(t *testing.T) {
	w := ManualHex(privateKeyHex, publicKeyHex, address)

	err := w.Validate()

	assert.NoError(t, err)
}
