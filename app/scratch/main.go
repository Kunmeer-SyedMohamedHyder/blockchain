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

	privateKey, err := crypto.LoadECDSA("zblock/accounts/kennedy.ecdsa")
	if err != nil {
		return fmt.Errorf("unable to load private key: %w", err)
	}

	tx := Tx{
		FromID: "Bill",
		ToID:   "Syed",
		Value:  1000,
	}

	data, err := json.Marshal(tx)
	if err != nil {
		return fmt.Errorf("unable to marshal: %w", err)
	}

	v := crypto.Keccak256(data)

	// Sign the hash with the private key to produce a signature.
	sig, err := crypto.Sign(v, privateKey)
	if err != nil {
		return fmt.Errorf("unable to sign: %w", err)
	}

	fmt.Println("SIG:", hexutil.Encode(sig))

	//======================================
	// DOWN THE WIRE

	// see if the v has changed, the public key received from the signature will also change.
	// i.e if v passed here and that when we used to sign is different then only pub key will be different.
	publicKey, err := crypto.SigToPub(v, sig)
	if err != nil {
		return fmt.Errorf("unable to pub: %w", err)
	}

	fmt.Println("PUB:", crypto.PubkeyToAddress(*publicKey).String())

	//=======================================================================================

	tx = Tx{
		FromID: "Bill",
		ToID:   "Hyder",
		Value:  10,
	}

	data2, err := json.Marshal(tx)
	if err != nil {
		return fmt.Errorf("unable to marshal: %w", err)
	}

	v2 := crypto.Keccak256(data2)

	// Sign the hash with the private key to produce a signature.
	sig2, err := crypto.Sign(v2, privateKey)
	if err != nil {
		return fmt.Errorf("unable to sign: %w", err)
	}

	fmt.Println("SIG:", hexutil.Encode(sig2))

	//=======================================
	// DOWN THE WIRE

	// even after the v changed pub key remains the same as we use the same private key to sign the v
	// and the v passed to create the signature and to get the pubKey is same.
	publicKey2, err := crypto.SigToPub(v2, sig2)
	if err != nil {
		return fmt.Errorf("unable to pub: %w", err)
	}

	fmt.Println("PUB:", crypto.PubkeyToAddress(*publicKey2).String())

	//=======================================================================================

	tx3 := Tx{
		FromID: "Syed",
		ToID:   "Asna",
		Value:  99,
	}

	data3, err := json.Marshal(tx3)
	if err != nil {
		return fmt.Errorf("unable to marshal: %w", err)
	}

	v3 := crypto.Keccak256(data3)

	// Sign the hash with the private key to produce a signature.
	sig3, err := crypto.Sign(v3, privateKey)
	if err != nil {
		return fmt.Errorf("unable to sign: %w", err)
	}

	fmt.Println("SIG:", hexutil.Encode(sig3))

	//=======================================
	// DOWN THE WIRE

	// some attacker got in middle and this data was changed
	// because this is the data that will be passed over the network.
	tx3 = Tx{
		FromID: "Syed",
		ToID:   "Rasna",
		Value:  99,
	}

	data3, err = json.Marshal(tx3)
	if err != nil {
		return fmt.Errorf("unable to marshal: %w", err)
	}

	v3 = crypto.Keccak256(data3)

	// the pubKey changes and hence we will be writing to some unkown person's ID.
	publicKey3, err := crypto.SigToPub(v3, sig3)
	if err != nil {
		return fmt.Errorf("unable to pub: %w", err)
	}

	fmt.Println("PUB:", crypto.PubkeyToAddress(*publicKey3).String())

	return nil
}
