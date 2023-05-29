package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/Kunmeer-SyedMohamedHyder/blockchain/foundation/blockchain/database"
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
		FromID: "0xF01813E4B85e178A83e29B8E7bF26BD830a25f32",
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
		FromID: "0xF01813E4B85e178A83e29B8E7bF26BD830a25f32",
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
	// Hence always pass FromID as the pubKey ID to validate it.
	tx3 = Tx{
		FromID: "0xF01813E4B85e178A83e29B8E7bF26BD830a25f32",
		ToID:   "Asna",
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

	pubID := crypto.PubkeyToAddress(*publicKey3).String()

	fmt.Println("PUB:", pubID)

	if pubID != tx.FromID {
		return fmt.Errorf("publickey ID mismatch")
	}

	vv, r, s, err := ToVRSFromHexSignature(hexutil.Encode(sig3))
	if err != nil {
		return fmt.Errorf("unable to vrs (%v): %w", hexutil.Encode(sig3), err)
	}

	fmt.Println("V|R|S", vv, r, s)

	//=======================================================================================

	newTx, err := database.NewTx(
		1,
		1,
		"0xF01813E4B85e178A83e29B8E7bF26BD830a25f32",
		"0xF01813E4B85e178A83e29B8E7bF26BD830a25f32",
		100,
		0,
		nil)
	if err != nil {
		return fmt.Errorf("unable to initialize a newTx: %w", err)
	}

	signedTx, err := newTx.Sign(privateKey)
	if err != nil {
		return fmt.Errorf("unable to sign a transaction: %w", err)
	}

	fmt.Println("======================")
	fmt.Printf("signedTx: %#v", signedTx)

	return nil
}

// ToVRSFromHexSignature converts a hex representation of the signature into
// its R, S and V parts.
func ToVRSFromHexSignature(sigStr string) (v, r, s *big.Int, err error) {
	sig, err := hex.DecodeString(sigStr[2:])
	if err != nil {
		return nil, nil, nil, err
	}

	r = big.NewInt(0).SetBytes(sig[:32])
	s = big.NewInt(0).SetBytes(sig[32:64])
	v = big.NewInt(0).SetBytes([]byte{sig[64]})

	return v, r, s, nil
}
