package ledgerapi

import (
	"errors"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
)

// HistoricEntryInterface functionality for
// getting details about an historic entry
type HistoricEntryInterface interface {
	GetTxId() string
	GetTimeStamp() *timestamp.Timestamp
	GetIsDelete() bool
	GetValue(v interface{}) error
}

// HistoryQueryIteratorInterface functionality for
// iterating over historic data in the chain
type HistoryQueryIteratorInterface interface {
	shim.CommonIteratorInterface

	Next() (*HistoricEntry, error)
}

// HistoricEntry implementation of
// HistoricEntryInterface
type HistoricEntry struct {
	Serializer LedgerSerializerInterface
	*queryresult.KeyModification
}

// GetValue deserializes the value from the historic
// snapshot into the provided interface
func (he *HistoricEntry) GetValue(v interface{}) error {
	bytes := he.KeyModification.GetValue()

	return he.Serializer.FromBytes(bytes, v)
}

// HistoryQueryIterator implementation of
// HistoryQueryIteratorInterface
type HistoryQueryIterator struct {
	Serializer LedgerSerializerInterface
	Iterator   shim.HistoryQueryIteratorInterface
}

// Close shuts the iterator
func (hqi *HistoryQueryIterator) Close() error {
	return hqi.Iterator.Close()
}

// HasNext are there unhandled values still
func (hqi *HistoryQueryIterator) HasNext() bool {
	return hqi.Iterator.HasNext()
}

// Next get the next value
func (hqi *HistoryQueryIterator) Next() (*HistoricEntry, error) {
	if !hqi.Iterator.HasNext() {
		return nil, errors.New("No next value")
	}

	keymod, err := hqi.Iterator.Next()

	if err != nil {
		return nil, err
	}

	he := new(HistoricEntry)
	he.KeyModification = keymod

	return he, nil
}
