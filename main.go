package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Config struct {
	RPC               string   `json:"rpc"`
	SenderPrivateKeys []string `json:"senderPrivateKeys"`
	ReceiverAddress   string   `json:"receiverAddress"`
}

var (
	configFile     string
	rpcFlag        string
	receiverFlag   string
	senderKeysFlag stringSliceFlag
)

type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func loadConfig() (*Config, error) {
	flag.StringVar(&configFile, "config", "", "Path to config file")
	flag.StringVar(&rpcFlag, "rpc", "", "EVM RPC URL")
	flag.StringVar(&receiverFlag, "receiver", "", "Receiver address")
	flag.Var(&senderKeysFlag, "sender", "Sender private key (can be specified multiple times)")
	flag.Parse()

	var config Config

	if configFile != "" {
		file, err := os.ReadFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}

		err = json.Unmarshal(file, &config)
		if err != nil {
			return nil, fmt.Errorf("error parsing config file: %v", err)
		}
	}

	// Override config with command-line flags if provided
	if rpcFlag != "" {
		config.RPC = rpcFlag
	}
	if receiverFlag != "" {
		config.ReceiverAddress = receiverFlag
	}
	if len(senderKeysFlag) > 0 {
		config.SenderPrivateKeys = senderKeysFlag
	}

	// Validate config
	if config.RPC == "" {
		return nil, fmt.Errorf("RPC URL is required")
	}
	if config.ReceiverAddress == "" {
		return nil, fmt.Errorf("receiver address is required")
	}
	if !common.IsHexAddress(config.ReceiverAddress) {
		return nil, fmt.Errorf("invalid receiver address")
	}
	if len(config.SenderPrivateKeys) == 0 {
		return nil, fmt.Errorf("at least one sender private key is required")
	}

	return &config, nil
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// EVM client
	client, err := ethclient.Dial(config.RPC)
	if err != nil {
		log.Fatalf("Failed to connect to the EVM client: %v", err)
	}

	// Receiver's address
	receiverAddress := common.HexToAddress(config.ReceiverAddress)

	for _, privateKeyHex := range config.SenderPrivateKeys {
		// Transfer all ETH from each sender to the receiver
		err := transferAllETH(client, privateKeyHex, receiverAddress)
		if err != nil {
			log.Printf("Failed to transfer ETH from private key %s: %v", privateKeyHex, err)
		}
	}
}

func transferAllETH(client *ethclient.Client, privateKeyHex string, receiverAddress common.Address) error {
	// Convert the private key from hex to ECDSA format
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("invalid private key: %v", err)
	}

	// Derive the sender's address from the private key
	senderAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Println("Sender Address:", senderAddress.Hex())

	// Get the sender's nonce (transaction count)
	nonce, err := client.PendingNonceAt(context.Background(), senderAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}

	// Get the current balance of the sender
	balance, err := client.BalanceAt(context.Background(), senderAddress, nil)
	if err != nil {
		return fmt.Errorf("failed to retrieve account balance: %v", err)
	}
	fmt.Println("Sender Balance:", balance)

	// Get the gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %v", err)
	}

	// Define the gas limit for a standard ETH transfer
	gasLimit := uint64(21000)

	// Calculate the value to send (balance - gasCost)
	gasCost := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	if balance.Cmp(gasCost) <= 0 {
		return fmt.Errorf("insufficient balance to cover gas cost")
	}
	valueToSend := new(big.Int).Sub(balance, gasCost)

	// Create the transaction
	tx := types.NewTransaction(nonce, receiverAddress, valueToSend, gasLimit, gasPrice, nil)

	// Sign the transaction with the sender's private key
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get network ID: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send the transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}

	fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())
	return nil
}
