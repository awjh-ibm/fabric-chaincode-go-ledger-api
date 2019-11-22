package ledgerapi

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

const WorldStateCollection string = "worldstate"

type Collection struct {
	Name       string
	Ctx        contractapi.TransactionContextInterface
	Serializer LedgerSerializerInterface
}

func (c *Collection) GetState(key string, v interface{}) error {
	// returns error when state doesnt exist
	return nil
}

func (c *Collection) GetStateByRange(start, end string) (*EntryIterator, error) {
	it := new(EntryIterator)
	it.Serializer = c.Serializer
	it.StateQueryIteratorInterface, _ = c.Ctx.GetStub().GetStateByRange(start, end) // or private

	return it, nil
}

func (c *Collection) GetStateByPartialCompositeKey(objectType string, attributes []string) (*EntryIterator, error) {
	it := new(EntryIterator)
	it.Serializer = c.Serializer
	it.StateQueryIteratorInterface, _ = c.Ctx.GetStub().GetStateByPartialCompositeKey(objectType, attributes) // or private

	return it, nil
}

func (c *Collection) GetQueryResult(query string) (*EntryIterator, error) {
	it := new(EntryIterator)
	it.Serializer = c.Serializer
	it.StateQueryIteratorInterface, _ = c.Ctx.GetStub().GetQueryResult(query) // or private

	return it, nil
}

func (c *Collection) GetHistoryForKey(key string) (*HistoricEntryIterator, error) {
	it := new(HistoricEntryIterator)
	it.Serializer = c.Serializer
	it.HistoryQueryIteratorInterface, _ = c.Ctx.GetStub().GetHistoryForKey(key)

	return it, nil
}

func (c *Collection) GetHash(key string) ([]byte, error) {
	return c.Ctx.GetStub().GetPrivateDataHash(c.Name, key)
}

func (c *Collection) PutState(key string, state interface{}) error {
	return nil
}

func (c *Collection) DelState(key string, state interface{}) error {
	return nil
}
