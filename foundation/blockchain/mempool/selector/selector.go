// Package selector provides different transaction selecting algorithms.
package selector

import (
	"fmt"
	"strings"

	"github.com/Kunmeer-SyedMohamedHyder/blockchain/foundation/blockchain/database"
)

// List of different select strategies.
const (
	StrategyTip         = "tip"
	StrategyTipAdvanced = "tip_advanced"
)

// Map of different select strategies with functions.
var strategies = map[string]Func{
	StrategyTip:         tipSelect,
	StrategyTipAdvanced: advancedTipSelect,
}

// Func defines a function that takes a mempool of transactions grouped by
// account and selects howMany of them in an order based on the functions
// strategy. All selector functions MUST respect nonce ordering. Receiving 0
// for howMany must return all the transactions in the strategies ordering.
type Func func(transactions map[database.AccountID][]database.BlockTx, howMany int) []database.BlockTx

// Retrieve returns the specified select strategy function.
func Retrieve(strategy string) (Func, error) {
	fn, exists := strategies[strings.ToLower(strategy)]
	if !exists {
		return nil, fmt.Errorf("strategy %q does not exist", strategy)
	}
	return fn, nil
}

// =============================================================================

// byNonce provides sorting support by the transaction id value.
type byNonce []database.BlockTx

// Len returns the number of transactions in the list.
func (bn byNonce) Len() int {
	return len(bn)
}

// Less helps to sort the list by nonce in ascending order to keep the
// transactions in the right order of processing.
func (bn byNonce) Less(i, j int) bool {
	return bn[i].Nonce < bn[j].Nonce
}

// Swap moves transactions in the order of the nonce value.
func (bn byNonce) Swap(i, j int) {
	bn[i], bn[j] = bn[j], bn[i]
}

// =============================================================================

// byTip provides sorting support by the transaction tip value.
type byTip []database.BlockTx

// Len returns the number of transactions in the list.
func (bt byTip) Len() int {
	return len(bt)
}

// Less helps to sort the list by tip in decending order to pick the
// transactions that provide the best reward.
func (bt byTip) Less(i, j int) bool {
	return bt[i].Tip > bt[j].Tip
}

// Swap moves transactions in the order of the tip value.
func (bt byTip) Swap(i, j int) {
	bt[i], bt[j] = bt[j], bt[i]
}
