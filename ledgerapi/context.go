package ledgerapi

import (
	"github.com/awjh-ibm/fabric-chaincode-go/contractapi"
)

// TransactionContextInterface describes functions of transaction context
type TransactionContextInterface interface {
	contractapi.TransactionContextInterface
	GetLedger() LedgerInterface
}

// TransactionContext implementation of TransactionContextInterface
// which extends contractapi.TransactionContext
type TransactionContext struct {
	contractapi.TransactionContext
}

// GetLedger returns Ledger with ctx set as Ctx
func (ctx *TransactionContext) GetLedger() LedgerInterface {
	ledger := Ledger{
		Ctx: ctx,
	}

	return &ledger
}
