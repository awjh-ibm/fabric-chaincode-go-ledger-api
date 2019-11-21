package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// PaperTransactionContext implementation of ledgerapi.TransactionContextInterface
// which extends contractapi.TransactionContext
type PaperTransactionContext struct {
	contractapi.TransactionContext
}

// GetLedger returns Ledger with ctx set as Ctx
func (ctx *PaperTransactionContext) GetLedger() *PaperLedger {
	ledger := PaperLedger{
		Ctx: ctx,
	}

	return &ledger
}
