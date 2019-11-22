package ledgerapi

import (
	"errors"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
)

type QueryEntry struct {
	*queryresult.KV
	Serializer LedgerSerializerInterface
}

func (e *QueryEntry) GetValue(v interface{}) error {
	bytes := e.KV.GetValue()

	return e.Serializer.FromBytes(bytes, v)
}

type HistoricEntry struct {
	*queryresult.KeyModification
	Serializer LedgerSerializerInterface
}

func (e *HistoricEntry) GetValue(v interface{}) error {
	bytes := e.KeyModification.GetValue()

	return e.Serializer.FromBytes(bytes, v)
}

type EntryIterator struct {
	shim.StateQueryIteratorInterface
	Serializer LedgerSerializerInterface
}

func (ei *EntryIterator) Next() (*QueryEntry, error) {
	if ei.HasNext() {
		return nil, errors.New("No next value")
	}

	keyVal, err := ei.StateQueryIteratorInterface.Next()

	if err != nil {
		return nil, err
	}

	qe := new(QueryEntry)
	qe.Serializer = ei.Serializer
	qe.KV = keyVal

	return qe, nil
}

type HistoricEntryIterator struct {
	shim.HistoryQueryIteratorInterface
	Serializer LedgerSerializerInterface
}

func (hei *HistoricEntryIterator) Next() (*HistoricEntry, error) {
	if hei.HasNext() {
		return nil, errors.New("No next value")
	}

	keyMod, err := hei.HistoryQueryIteratorInterface.Next()

	if err != nil {
		return nil, err
	}

	he := new(HistoricEntry)
	he.Serializer = hei.Serializer
	he.KeyModification = keyMod

	return he, nil
}
