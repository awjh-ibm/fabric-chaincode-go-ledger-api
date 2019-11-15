package ledgerapi

// LedgerInterface a holder for many collections
type LedgerInterface interface {
	// GetCollection returns collection with given name
	GetCollection(string) CollectionInterface

	// GetDefaultCollection returns the default collection
	GetDefaultCollection() CollectionInterface
}

// Ledger implements LedgerInterface using default ledgerapi
// Collection
type Ledger struct {
	Ctx TransactionContextInterface
}

// GetCollection returns a Collection with name that uses the
// JSONLedgerSerializer
func (l *Ledger) GetCollection(name string) CollectionInterface {
	c := Collection{
		Name:       name,
		Serializer: new(JSONLedgerSerializer),
		Stub:       l.Ctx.GetStub(),
	}

	return &c
}

// GetDefaultCollection returns a Collection for managing the world
// state
func (l *Ledger) GetDefaultCollection() CollectionInterface {
	return l.GetCollection(WorldStateIdentifier)
}
