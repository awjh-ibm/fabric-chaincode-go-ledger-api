package ledgerapi

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

type Ledger struct {
	Ctx contractapi.TransactionContextInterface
}

func (l *Ledger) GetCollection(name string) *Collection {
	return &Collection{name, l.Ctx, new(JSONLedgerSerializer)}
}

func (l *Ledger) GetDefaultCollection() *Collection {
	return l.GetCollection(WorldStateCollection)
}
