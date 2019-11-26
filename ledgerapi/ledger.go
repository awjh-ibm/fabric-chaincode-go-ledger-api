package ledgerapi

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

// Ledger provides way to get into a collection
type Ledger struct {
	Ctx contractapi.TransactionContextInterface
}

// GetCollection returns a collection with the given name
func (l *Ledger) GetCollection(name string) *Collection {
	return &Collection{name, l.Ctx}
}

// GetDefaultCollection returns a collection with WorldStateCollection
// as its name
func (l *Ledger) GetDefaultCollection() *Collection {
	return l.GetCollection(WorldStateCollection)
}
