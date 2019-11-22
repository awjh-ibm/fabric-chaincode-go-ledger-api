package blockchainapi

import "encoding/json"

// BlockchainSerializerInterface defines the functions a valid blockchain serializer
// should have. Serialzers to be used by a chaincode must implement
// this interface. Serializers are called when data is added, updated
// or read to/from a collection in the ledger.
type BlockchainSerializerInterface interface {
	// FromBytes called on reads from a collection. Should convert the
	// value retrieved from that collection back into its Go usable value
	// using the interface passed.
	FromBytes([]byte, interface{}) error

	// ToBytes called on writes (add, update) to a collection. Should
	// convert the given Go value to a byte array to be written to that
	// collection
	ToBytes(interface{}) ([]byte, error)
}

// JSONLedgerSerializer an implementation of LedgerSerializerInterface that
// handles values using encoding/json marshall unmarshal
type JSONLedgerSerializer struct{}

// FromBytes alias for encoding/json's json.Unmarshal
func (jls *JSONLedgerSerializer) FromBytes(bytes []byte, v interface{}) error {
	return json.Unmarshal(bytes, v)
}

// ToBytes alias for encoding/json's json.Marshal
func (jls *JSONLedgerSerializer) ToBytes(v interface{}) ([]byte, error) {
	bytes, err := json.Marshal(v)

	return bytes, err
}
