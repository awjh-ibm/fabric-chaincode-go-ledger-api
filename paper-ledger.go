package main

import "github.com/awjh-ibm/fabric-contract-api-go-ledger-api/ledgerapi"

// PaperLedger implements LedgerInterface using PaperCollection
type PaperLedger struct {
	Ctx *PaperTransactionContext
}

// GetCollection returns a Collection with name that uses the
// JSONLedgerSerializer
func (l *PaperLedger) GetCollection(name string) *PaperCollection {
	c := PaperCollection{
		Name: name,
		Stub: l.Ctx.GetStub(),
	}

	return &c
}

// GetDefaultCollection returns a Collection for managing the world
// state
func (l *PaperLedger) GetDefaultCollection() *PaperCollection {
	return l.GetCollection(ledgerapi.WorldStateIdentifier)
}
