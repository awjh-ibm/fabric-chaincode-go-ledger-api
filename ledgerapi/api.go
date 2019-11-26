package ledgerapi

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

// GetLedger returns a new instance of ledger with ctx set
func GetLedger(ctx contractapi.TransactionContextInterface) *Ledger {
	ledger := new(Ledger)
	ledger.Ctx = ctx

	return ledger
}

// CreateCompositeKey returns a key for use in the world state
func CreateCompositeKey(objectType string, attributes []string) (string, error) {
	return "", nil
}
