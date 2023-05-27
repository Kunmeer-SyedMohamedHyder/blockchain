package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {

	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

// Tx is the transactional information between two parties.
type Tx struct {
	FromID string `json:"from"`  // Ethereum: Account sending the transaction. Will be checked against signature.
	ToID   string `json:"to"`    // Ethereum: Account receiving the benefit of the transaction.
	Value  uint64 `json:"value"` // Ethereum: Monetary value received from this transaction.
}

func run() error {

	tx := Tx{
		FromID: "Bill",
		ToID:   "Syed",
		Value:  1000,
	}

	privateKey, err := crypto.LoadECDSA("C:\\Users\\kunmeer\\go\\src\\github.com\\blockchain\\zblock\\accounts\\baba.ecdsa")
	if err != nil {
		return fmt.Errorf("unable to load private key: %w", err)
	}

	data, err := json.Marshal(tx)
	if err != nil {
		return fmt.Errorf("unable to marshal: %w", err)
	}

	// Sign the hash with the private key to produce a signature.
	sig, err := crypto.Sign(data, privateKey)
	if err != nil {
		return err
	}

	fmt.Println(hexutil.Encode(sig))

	return nil
}
