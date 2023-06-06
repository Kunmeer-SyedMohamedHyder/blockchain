package cmd

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Kunmeer-SyedMohamedHyder/blockchain/foundation/blockchain/database"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

var (
	nonce uint64
	from  string
	to    string
	value uint64
	tip   uint64
	data  []byte
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send transaction",
	Run:   sendRun,
}

func init() {
	rootCmd.AddCommand(sendCmd)
	sendCmd.Flags().StringVarP(&url, "url", "u", "http://localhost:8080", "Url of the node.")
	sendCmd.Flags().Uint64VarP(&nonce, "nonce", "n", 0, "id for the transaction.")
	sendCmd.Flags().StringVarP(&from, "from", "f", "", "Who is sending the transaction.")
	sendCmd.Flags().StringVarP(&to, "to", "t", "", "Who is receiving the transaction.")
	sendCmd.Flags().Uint64VarP(&value, "value", "v", 0, "Value to send.")
	sendCmd.Flags().Uint64VarP(&tip, "tip", "c", 0, "Tip to send.")
	sendCmd.Flags().BytesHexVarP(&data, "data", "d", nil, "Data to send.")
}

func sendRun(cmd *cobra.Command, args []string) {
	privateKey, err := crypto.LoadECDSA(getPrivateKeyPath())
	if err != nil {
		log.Fatal(err)
	}

	sendWithDetails(privateKey)
}

func sendWithDetails(privateKey *ecdsa.PrivateKey) {
	fromAccount, err := database.ToAccountID(from)
	if err != nil {
		log.Fatal(err)
	}

	toAccount, err := database.ToAccountID(to)
	if err != nil {
		log.Fatal(err)
	}

	const chainID = 1
	tx, err := database.NewTx(chainID, nonce, fromAccount, toAccount, value, tip, data)
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := tx.Sign(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	data, err := json.Marshal(signedTx)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/v1/tx/submit", url), "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
}
