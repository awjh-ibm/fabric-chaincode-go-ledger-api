package ledgerapi

import (
	"errors"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// StateInterface functions a state will need
type StateInterface interface {
	GetValue() ([]byte, error)
	GetKey() string
	GetSplitKey() (string, []string, error)
	GetHistory() (*HistoricStateIterator, error)
	GetHash() ([]byte, error)
	GetValidationParameter() ([]byte, error)
	SetValidationParameter([]byte) error
}

// BasicState implementation of StateInterface where
// GetValue reads from world state
type BasicState struct {
	Ctx        contractapi.TransactionContextInterface
	Key        string
	Collection string
}

// GetValue returns collection value of this state
func (bs *BasicState) GetValue() ([]byte, error) {
	if bs.Collection != WorldStateCollection {
		return bs.Ctx.GetStub().GetPrivateData(bs.Collection, bs.Key)
	}

	return bs.Ctx.GetStub().GetState(bs.Key)
}

// GetKey returns the key for the state
func (bs *BasicState) GetKey() string {
	return bs.Key
}

// GetSplitKey returns the key for the state
func (bs *BasicState) GetSplitKey() (string, []string, error) {
	return bs.Ctx.GetStub().SplitCompositeKey(bs.Key)
}

// GetHistory gets the historic entries from the ledger
func (bs *BasicState) GetHistory() (*HistoricStateIterator, error) {
	if bs.Collection != WorldStateCollection {
		return nil, errors.New("Cannot get history when state not in world state collection")
	}

	return new(HistoricStateIterator), nil
}

// GetHash gets the state's hash as its stored in its collection
func (bs *BasicState) GetHash() ([]byte, error) {
	if bs.Collection == WorldStateCollection {
		return nil, errors.New("Cannot get hash when state in world state collection")
	}

	return []byte{}, nil
}

// GetValidationParameter gets the endorsement that exists on the
// state's key
func (bs *BasicState) GetValidationParameter() ([]byte, error) {
	return []byte{}, nil
}

// SetValidationParameter sets the endorsement that exists on the
// state's key
func (bs *BasicState) SetValidationParameter([]byte) error {
	return nil
}

// QueryState implementation of StateInterface where
// GetValue returns a preset value
type QueryState struct {
	BasicState
	Value []byte
}

// GetValue returns set value of this state
func (qs *QueryState) GetValue() ([]byte, error) {
	return qs.Value, nil
}

// QueryStateIterator functionality to navigate through query response
type QueryStateIterator struct {
	shim.StateQueryIteratorInterface
	Ctx        contractapi.TransactionContextInterface
	Collection string
}

// Next get the next state in the iterator
func (qsi *QueryStateIterator) Next() (StateInterface, error) {
	if qsi.HasNext() {
		return nil, errors.New("No next value")
	}

	keyVal, err := qsi.StateQueryIteratorInterface.Next()

	if err != nil {
		return nil, err
	}

	qs := new(QueryState)
	qs.Ctx = qsi.Ctx
	qs.Collection = qsi.Collection
	qs.Key = keyVal.GetKey()
	qs.Value = keyVal.GetValue()

	return qs, nil
}

// HistoricStateIterator provides way to go through
// the collection history for a key
type HistoricStateIterator struct {
	shim.HistoryQueryIteratorInterface
}
