package ledgerapi

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

func GetLedger(ctx contractapi.TransactionContextInterface) *Ledger {
	ledger := new(Ledger)
	ledger.Ctx = ctx

	return ledger
}

func CreateCompositeKey(objectType string, attributes []string) (string, error) {
	return "", nil
}

func SplitCompositeKey(key string) (string, []string, error) {
	return "", []string{}, nil
}
