package wallet

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

// Wallet is a struct containing the 3 identifiers of an ethereum wallet
// 	1. The private key, generated randomly.
// 	2. The public key, derived from the private key.
// 	3. The address, derived from the public key.
// To access these fields from another package in hexadecimal format, use the receiver functions PrivateKeyHex, PublicKeyHex, and Address.
type Wallet struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	address    string
}

// New creates a new Wallet from a randomly generated private key.
func New() *Wallet {
	// PRIVATE KEY GENERATION
	// Generate new private key.
	privateKey, err := ethcrypto.GenerateKey()
	if err != nil {
		panic(errors.Wrap(err, "generating private key"))
	}

	// Derive the rest of the wallet (public key + wallet) from this private key.
	return FromPrivateKey(privateKey)
}

// Manual creates a new Wallet with privateKey, publicKey, and address.
// It would be wise to Wallet.Validate this wallet after creating it.
func Manual(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, address string) *Wallet {
	w := &Wallet{
		privateKey: privateKey,
		publicKey:  publicKey,
		address:    address,
	}

	return w
}

// ManualHex creates a new Wallet with privateKeyHex and publicKeyHex parsed into ECDSA format and address as-is.
func ManualHex(privateKeyHex, publicKeyHex, address string) *Wallet {
	privateKey := PrivateKeyHexToECDSA(privateKeyHex)
	publicKey := PublicKeyHexToECDSA(publicKeyHex)

	return Manual(privateKey, publicKey, address)
}

// FromPrivateKey derives a new Wallet from privateKey by deriving its public key.
// This should be called from New when a new wallet is being created from scratch.
func FromPrivateKey(privateKey *ecdsa.PrivateKey) *Wallet {
	// PUBLIC KEY DERIVATION
	// Derive public key from private key.
	publicKeyCrypto := privateKey.Public()
	// Cast public key to *ecdsa.PublicKey type.
	publicKey, ok := publicKeyCrypto.(*ecdsa.PublicKey)
	if !ok {
		panic(fmt.Errorf("cannot case public key to ECDSA, is of type \"%T\"\n", publicKeyCrypto))
	}

	// Derive the rest of the wallet (the address) from the given private key and this public key.
	return FromPublicKey(privateKey, publicKey)
}

// FromPrivateKeyHex calls FromPrivateKey after parsing privateKeyHex into ECDSA format.
func FromPrivateKeyHex(privateKeyHex string) *Wallet {
	privateKey := PrivateKeyHexToECDSA(privateKeyHex)

	return FromPrivateKey(privateKey)
}

// FromPublicKey derives a new Wallet from privateKey and publicKey by deriving its address.
// This should be called from FromPrivateKey when a new wallet is being created from scratch.
func FromPublicKey(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) *Wallet {
	// ADDRESS DERIVATION
	// Simply use go-ethereum's crypto package's PubkeyToAddress function to derive the wallet address from the public key.
	address := ethcrypto.PubkeyToAddress(*publicKey).Hex()

	w := &Wallet{
		privateKey: privateKey,
		publicKey:  publicKey,
		address:    address,
	}

	return w
}

// FromPublicKeyHex calls FromPublicKey after parsing privateKeyHex and publicKeyHex into ECDSA format.
func FromPublicKeyHex(privateKeyHex, publicKeyHex string) *Wallet {
	privateKey := PrivateKeyHexToECDSA(privateKeyHex)
	publicKey := PublicKeyHexToECDSA(publicKeyHex)

	return FromPublicKey(privateKey, publicKey)
}

// PrivateKeyHex returns w's privateKey in hexadecimal format.
func (w *Wallet) PrivateKeyHex() string {
	// Dump private key to bytes.
	privateKeyBytes := ethcrypto.FromECDSA(w.privateKey)
	// Encode private key bytes into hexadecimal.
	privateKeyHex := hexutil.Encode(privateKeyBytes)
	// Shave off the first 2 characters of the hexadecimal encoding since that is not part of the private key format and is only a hex identifier, "0x".
	privateKeyHex = privateKeyHex[2:]

	return privateKeyHex
}

// PublicKeyHex returns w's publicKey in hexadecimal format.
func (w *Wallet) PublicKeyHex() string {
	// Dump public key to bytes.
	publicKeyBytes := ethcrypto.FromECDSAPub(w.publicKey)
	// Encode public key bytes into hexadecimal.
	publicKeyHex := hexutil.Encode(publicKeyBytes)
	// Shave off the first 4 characters, 2 characters of the hexadecimal encoding, "0x", AND the next 2 characters, which is always "04" since it is a constant prefix it is not necessary.
	publicKeyHex = publicKeyHex[4:]

	return publicKeyHex
}

// Address returns w's address as-is (hexadecimal format).
func (w *Wallet) Address() string {
	// The address requires no massaging/formatting, simply return it.
	return w.address
}

// Clone duplicates w into a new Wallet.
func (w *Wallet) Clone() *Wallet {
	w2 := ManualHex(w.PrivateKeyHex(), w.PublicKeyHex(), w.Address())

	return w2
}

// Equals checks if w contains the same wallet data as w2 by comparing their hexadecimal encodings.
// If no error is returned then w and w2 are equal.
func (w *Wallet) Equals(w2 *Wallet) error {
	if w.PrivateKeyHex() != w2.PrivateKeyHex() {
		return fmt.Errorf("private keys are not equal: original: \"%s\", comparable: \"%s\"", w.PrivateKeyHex(), w2.PublicKeyHex())
	}

	if w.PublicKeyHex() != w2.PublicKeyHex() {
		return fmt.Errorf("public keys are not equal: original: \"%s\", comparable: \"%s\"", w.PublicKeyHex(), w2.PublicKeyHex())
	}

	if w.Address() != w2.Address() {
		return fmt.Errorf("addresses are not equal: original: \"%s\", comparable: \"%s\"", w.Address(), w2.Address())
	}

	return nil
}

// Validate checks if w represents a valid ethereum wallet by deriving the public key and address from w.privateKey and checking if they match what was already in w.
// This is valuable because take, for example, if you had a wallet with random hexadecimal information it would look the same as a "true" wallet, but if the address
// does not match the private key then you do not have access to the wallet.
// Make HEAVY USE of this function to ensure your wallet is accurate and you still have access to it.
func (w *Wallet) Validate() error {
	// Derive a new wallet from w.privateKey
	w2 := FromPrivateKey(w.privateKey)

	// Check if the derived wallet, w2, is equal to the original wallet, w.
	if err := w.Equals(w2); err != nil {
		// If they are not equivalent, then w is an invalid wallet.
		return errors.Wrap(err, "the generated wallet which was derived from provided wallet's private key is not equal to the original provided wallet")
	}

	// If they are equivalent, then w is a valid wallet.
	return nil
}

// PrivateKeyHexToECDSA is a helper function for converting a hexadecimal representation of a private key into ECDSA format.
func PrivateKeyHexToECDSA(privateKeyHex string) *ecdsa.PrivateKey {
	privateKey, err := ethcrypto.HexToECDSA(privateKeyHex)
	if err != nil {
		panic(errors.Wrap(err, "converting private key hex to ecdsa"))
	}

	return privateKey
}

// PublicKeyHexToECDSA is a helper function for converting a hexadecimal representation of a public key into ECDSA format.
func PublicKeyHexToECDSA(publicKeyHex string) *ecdsa.PublicKey {
	publicKeyBytes, err := hexutil.Decode("0x04" + publicKeyHex)
	if err != nil {
		panic(errors.Wrap(err, "converting public key hex to bytes"))
	}

	publicKey, err := ethcrypto.UnmarshalPubkey(publicKeyBytes)
	if err != nil {
		panic(errors.Wrap(err, "converting public key bytes to ecdsa"))
	}

	return publicKey
}
