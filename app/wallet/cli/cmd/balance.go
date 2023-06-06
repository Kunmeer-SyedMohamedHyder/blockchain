package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Kunmeer-SyedMohamedHyder/blockchain/foundation/blockchain/database"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Print your balance.",
	Run:   balanceRun,
}

var url string

func init() {
	rootCmd.AddCommand(balanceCmd)
	balanceCmd.Flags().StringVarP(&url, "url", "u", "http://localhost:8080", "Url of the node.")
}

func balanceRun(cmd *cobra.Command, args []string) {
	privateKey, err := crypto.LoadECDSA(getPrivateKeyPath())
	if err != nil {
		log.Fatal(err)
	}

	accountID := database.PublicKeyToAccountID(privateKey.PublicKey)
	fmt.Printf("Balance for Account: %v\n", accountID)

	resp, err := http.Get(fmt.Sprintf("%s/v1/accounts/list/%s", url, accountID))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	var balances = make(map[database.AccountID]database.Account)
	if err := decoder.Decode(&balances); err != nil {
		log.Fatal(err)
	}

	fmt.Println(balances[accountID].Balance)
}
