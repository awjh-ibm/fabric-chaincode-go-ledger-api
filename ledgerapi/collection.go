package ledgerapi

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

// WorldStateCollection identifier of the world state collection
const WorldStateCollection string = "worldstate"

// Collection representation of db
type Collection struct {
	Name string
	Ctx  contractapi.TransactionContextInterface
}

// GetState return interface for accessing entry to db
func (c *Collection) GetState(key string) StateInterface {
	bs := new(BasicState)
	bs.Ctx = c.Ctx
	bs.Key = key
	bs.Collection = c.Name

	return bs
}

// PutState write value to collection
func (c *Collection) PutState(key string, data []byte) error {
	if c.Name != WorldStateCollection {
		return c.Ctx.GetStub().PutPrivateData(c.Name, key, data)
	}

	return c.Ctx.GetStub().PutState(key, data)
}

// DeleteState deletes value from collection
func (c *Collection) DeleteState(key string) error {
	if c.Name != WorldStateCollection {
		return c.Ctx.GetStub().DelPrivateData(c.Name, key)
	}

	return c.Ctx.GetStub().DelState(key)
}

// GetStates return iterator for accessing stored states
func (c *Collection) GetStates(query QueryInterface) (*QueryStateIterator, error) {
	return query.Query(c.Ctx, c.Name)
}
