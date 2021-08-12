package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/Insulince/jeth/pkg/convert"
	"github.com/Insulince/jeth/pkg/eth"
	"github.com/Insulince/jeth/pkg/price"

	jio "github.com/Insulince/jlib/pkg/io"
)

const (
	defaultGasLimit          = uint64(21000)
	defaultSuggestedGasPrice = 0
)

type (
	Config struct {
		privateKeyHex         string
		receiverWalletAddress string
		gateway               string
		amount                float64
		gasPrice              int64
		gasLimit              uint64
		dryRun                bool
		help                  bool
	}
)

func getConfig() (cfg Config, err error) {
	// TODO(justin): Use dry run
	// TODO(justin): Use help

	flag.StringVar(&cfg.privateKeyHex, "private-key", "", "the hexadecimal private key of the sender's wallet [required via flag or stdin at runtime]")
	flag.StringVar(&cfg.receiverWalletAddress, "receiver-address", "", "the receiver's wallet address [required via flag or stdin at runtime]")
	flag.Float64Var(&cfg.amount, "amount", 0, "the amount of ethereum to send in ether units [required via flag or stdin at runtime]")
	flag.Int64Var(&cfg.gasPrice, "gas-price", defaultSuggestedGasPrice, "the gas price for your transaction")
	flag.Uint64Var(&cfg.gasLimit, "gas-limit", defaultGasLimit, "the gas limit for your transaction")
	flag.StringVar(&cfg.gateway, "gateway", eth.DefaultGateway, "the connection to your ethereum provider")
	flag.BoolVar(&cfg.dryRun, "dry-run", false, "don't actually send the transaction, just build and display it")
	flag.BoolVar(&cfg.help, "help", false, "display help message")
	flag.Parse()

	if cfg.privateKeyHex == "" {
		cfg.privateKeyHex = jio.MustPrivateInputWithPrompt("sender's private key not given via \"-private-key\" flag, enter manually instead: ")
		jio.SilentOutputln("")
	}
	if len(cfg.privateKeyHex) != 64 {
		return Config{}, errors.New("must provide a 64 character hexadecimal private key via \"-private-key\" or at runtime via stdin")
	}
	if cfg.receiverWalletAddress == "" {
		cfg.receiverWalletAddress = jio.MustInputWithPrompt("receiver's wallet address not given via \"-receiver-address\" flag, enter manually instead: ")
	}
	if len(cfg.receiverWalletAddress) != 42 {
		return Config{}, errors.New("must provide a 42 character hexadecimal wallet address starting with \"0x\" for receiver via \"-receiver-address\" or at runtime via stdin")
	}
	if cfg.gateway == "" {
		return Config{}, fmt.Errorf("must provide a non-blank ethereum gateway via \"-gateway\", or leave blank to use the default gateway, %s", eth.DefaultGateway)
	}
	if cfg.amount == 0 {
		amountStr := jio.MustInputWithPrompt("sender's ether amount to send not given via \"-amount\" flag, enter manually instead: ")
		cfg.amount, err = strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return Config{}, errors.Wrap(err, "stdin provided eth amount is not a float64 value")
		}
	}
	if cfg.amount <= 0 {
		return Config{}, errors.New("must provide a non-negative non-zero eth amount to send via \"-amount\" or at runtime via stdin")
	}
	if cfg.gasPrice < 0 {
		return Config{}, errors.New("must provide a non-negative gas price via \"-gas-price\" in wei units, or provide \"0\" or leave blank to choose the network's suggested gas price")
	}
	if cfg.gasLimit <= 0 {
		return Config{}, fmt.Errorf("must provide a non-negative non-zero gas limit via \"gas-limit\", or leave blank to use the default of %v", defaultGasLimit)
	}
	jio.Outputf("configuration parsed successfully (private key obfuscated):\n\t-private-key=%s\n\t-receiver-address=%s\n\t-amount=%v\n\t-gas-price=%v\n\t-gas-limit=%v\n\t-gateway=%s\n\t-dry-run=%v\n\t-help=%v\n", eth.ObfuscateKey(cfg.privateKeyHex), cfg.receiverWalletAddress, cfg.amount, cfg.gasPrice, cfg.gasLimit, cfg.gateway, cfg.dryRun, cfg.help)

	return cfg, nil
}

func main() {
	jio.Outputf("send initiated at %v\n", time.Now().Format(time.RFC3339Nano))
	defer func() { jio.Outputf("send completed at %v\n", time.Now().Format(time.RFC3339Nano)) }()

	ctx := context.Background()

	cfg, err := getConfig()
	if err != nil {
		panic(errors.Wrap(err, "getting config"))
	}

	usdPerEth, err := price.UsdPerEth()
	if err != nil {
		panic(errors.Wrap(err, "fetching latest eth price"))
	}
	jio.Outputf("current usd per ether (this figure will be used in later approximations): $%v\n", usdPerEth)

	client, err := ethclient.Dial(cfg.gateway)
	if err != nil {
		panic(errors.Wrap(err, "dialing eth gateway"))
	}
	jio.Outputf("connected to gateway: %s\n", cfg.gateway)

	privateKey, err := crypto.HexToECDSA(cfg.privateKeyHex)
	if err != nil {
		panic(errors.Wrap(err, "converting private key hex to ecdsa"))
	}
	jio.Outputf("converted private key to ECDSA: [PRIVATE] %s\n", eth.ObfuscateKey(cfg.privateKeyHex))

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic(errors.Wrap(err, "casting public key to ecdsa"))
	}
	publicKeyString := hexutil.Encode(crypto.FromECDSAPub(publicKeyECDSA))[4:]
	jio.Outputf("sender's public key extracted from given private key: [PUBLIC] %s\n", publicKeyString)

	senderAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	senderWalletAddress := senderAddress.String()
	jio.Outputf("sender's wallet address extracted from public key: [WALLET] %s\n", senderWalletAddress)

	nonce, err := client.PendingNonceAt(ctx, senderAddress)
	if err != nil {
		panic(errors.Wrapf(err, "fetching latest pending nonce for sender's wallet \"%s\"", senderWalletAddress))
	}
	jio.Outputf("sender's nonce extracted from wallet address: [NONCE] %v\n", nonce)

	bAmount := big.NewFloat(cfg.amount)
	jio.Outputf("ether to be sent: %v ether ($%.2f)\n", cfg.amount, convert.F(convert.EthToUsd(bAmount, usdPerEth)))
	bWei := convert.EthToWeiI(bAmount)
	jio.Outputf("equivalent wei to be sent: %s wei ($%.2f)\n", bWei.String(), convert.F(convert.WeiIToUsd(bWei, usdPerEth)))

	bGasPrice := big.NewInt(cfg.gasPrice)
	if cfg.gasPrice == defaultSuggestedGasPrice {
		jio.Outputln("fetching suggested gas price...")
		if bGasPrice, err = client.SuggestGasPrice(ctx); err != nil {
			panic(errors.Wrap(err, "getting suggested gas price"))
		}
		jio.Outputf("suggested gas price: %s wei ($%f)\n", bGasPrice.String(), convert.F(convert.WeiIToUsd(bGasPrice, usdPerEth)))
	}
	jio.Outputf("using gas price: %s wei ($%f)\n", bGasPrice.String(), convert.F(convert.WeiIToUsd(bGasPrice, usdPerEth)))

	bGasLimit := big.NewInt(int64(cfg.gasLimit))
	jio.Outputf("using gas limit (no unit): %s\n", bGasLimit.String())

	bTotalGas := new(big.Int).Mul(bGasPrice, bGasLimit)
	jio.Outputf("total gas for this transaction: %s wei ($%.2f)\n", bTotalGas.String(), convert.F(convert.WeiIToUsd(bTotalGas, usdPerEth)))

	gasProportion := float64(convert.I(bTotalGas)) / convert.F(convert.EthToWei(big.NewFloat(cfg.amount)))
	jio.Outputf("gas prices make up %.3f%% of the original value to be sent, the receiver's final amount will be short by this same percentage compared to what you originally opted to send\n", gasProportion*100)

	bWeiMinusGas := new(big.Int).Sub(bWei, bTotalGas)
	jio.Outputf("total wei to be sent excluding gas costs: %v wei ($%.2f)\n", bWeiMinusGas.String(), convert.F(convert.WeiIToUsd(bWeiMinusGas, usdPerEth)))
	bEthMinusGas := convert.WeiIToEth(bWeiMinusGas)
	jio.Outputf("equivalent total ether to be sent excluding gas costs (this is the actual value the receiver will get): %v eth ($%.2f)\n", bEthMinusGas.String(), convert.F(convert.EthToUsd(bEthMinusGas, usdPerEth)))

	toAddress := common.HexToAddress(cfg.receiverWalletAddress)
	jio.Outputf("will send to wallet address: %s\n", toAddress)

	jio.SilentOutputln("")
	summary := summarize(bAmount, bEthMinusGas, bGasPrice, bGasLimit, bTotalGas, senderWalletAddress, cfg.receiverWalletAddress, gasProportion, usdPerEth)
	jio.Outputln("----- SUMMARY -----")
	jio.SilentOutputln(summary)

	response := jio.MustInputWithPrompt("WARNING: you are about to send the above transaction to the ethereum network, please double check the summary above for accuracy, this cannot be undone if successful. PROCEED? [y/N]: ")
	response = strings.ToLower(response)
	if response != "y" && response != "yes" {
		jio.Output("aborting...")
		os.Exit(0)
	}
	jio.Outputln("proceeding...")

	jio.Outputln("building transaction...")
	legacyTx := types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    bWeiMinusGas,
		Gas:      cfg.gasLimit,
		GasPrice: bGasPrice,
		Data:     nil,
	}
	tx := types.NewTx(&legacyTx)
	jio.Outputln("transaction built successfully")

	chainId, err := client.NetworkID(ctx)
	if err != nil {
		panic(errors.Wrap(err, "getting chain id"))
	}
	jio.Outputf("retrieved chain id from gateway: %s\n", chainId)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
	if err != nil {
		panic(errors.Wrap(err, "signing transaction"))
	}
	jio.Outputln("transaction signed successfully")

	signedTxJsonBytes, err := signedTx.MarshalJSON()
	if err != nil {
		panic(errors.Wrap(err, "marshalling signed transaction into json"))
	}
	signedTxJson := string(signedTxJsonBytes)
	jio.SilentOutputln("")
	jio.Outputln("signed transaction json:")
	jio.SilentOutputln(signedTxJson)
	jio.SilentOutputln("")

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		panic(errors.Wrap(err, "sending transaction"))
	}
	jio.Outputf("success: transaction hash: [TRANSACTION] %s\n", signedTx.Hash().Hex())
}

func summarize(bAmount, bEthMinusGas *big.Float, bGasPrice, bGasLimit, bTotalGas *big.Int, senderWalletAddress, receiverWalletAddress string, gasProportion, usdPerEth float64) string {
	return fmt.Sprintf("ORIGINAL AMOUNT SENDING: %s ether ($%.2f)\nGAS: %s wei price * %s limit = %s wei (%s ether, $%.2f)\nGAS ADJUSTED AMOUNT SENDING: %s ether ($%.2f) [↓ %.3f%%]\nFROM:\t%s\nTO:\t%s\n",
		bAmount.String(), convert.F(convert.EthToUsd(bAmount, usdPerEth)),
		bGasPrice.String(), bGasLimit.String(), bTotalGas.String(), convert.WeiIToEth(bTotalGas).String(), convert.F(convert.WeiIToUsd(bTotalGas, usdPerEth)),
		bEthMinusGas.String(), convert.F(convert.EthToUsd(bEthMinusGas, usdPerEth)), gasProportion*100,
		senderWalletAddress,
		receiverWalletAddress,
	)
}
