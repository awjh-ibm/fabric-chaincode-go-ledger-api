package ledgerapi

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

// WorldStateIdentifier the identifier for the world state collection
const WorldStateIdentifier = "worldstate"

// CollectionInterface an interface to define the functionality for
// handling entries
type CollectionInterface interface {
	Exists(...string) (bool, error)
	Add(interface{}) error
	Get(interface{}, ...string) error
	GetHistory(...string) (HistoryQueryIteratorInterface, error)
	Update(interface{}) error
	Delete(...string) error
}

// Collection implementation of CollectionInterface
type Collection struct {
	Name       string
	Serializer LedgerSerializerInterface
	Stub       shim.ChaincodeStubInterface
}

func (c *Collection) validateEntry(entry interface{}) bool {
	typ := reflect.TypeOf(entry)

	return typ.Kind() == reflect.Struct || (typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct)
}

func (c *Collection) generateKey(entry interface{}) ([]string, error) {
	// could we cache the type to avoid having to do this every time?
	typ := reflect.TypeOf(entry)
	val := reflect.ValueOf(entry)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = reflect.ValueOf(entry).Elem()
	}

	keyParts := []string{}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		tag := field.Tag.Get("ledgerapi")

		if strings.Contains(tag, "primary_key") {
			keyParts = append(keyParts, fieldVal.String())
		}
	}

	if len(keyParts) == 0 {
		return nil, errors.New("Could not generate key. No field in entry with tag ledgerapi and value primary_key")
	}

	return keyParts, nil
}

func (c *Collection) formatKey(components []string) (string, error) {
	formatted, err := c.Stub.CreateCompositeKey(components[0], components[1:])

	if err != nil {
		return "", fmt.Errorf("Failed to add to collection. %s", err.Error())
	}

	return formatted, nil
}

func (c *Collection) exists(key string) (bool, error) {
	var bytes []byte
	var err error

	if c.Name != WorldStateIdentifier {
		bytes, err = c.Stub.GetPrivateData(c.Name, key)
	} else {
		bytes, err = c.Stub.GetState(key)
	}

	if err != nil {
		return false, err
	}

	return bytes != nil, nil
}

// Exists returns true if key exists in the ledger
func (c *Collection) Exists(keyComponents ...string) (bool, error) {
	var err error

	formattedKey, err := c.formatKey(keyComponents)

	if err != nil {
		return false, fmt.Errorf("Failed to ascertain whether key exists in collection. %s", err.Error())
	}

	return c.exists(formattedKey)
}

// Add adds a serialized version of the entry to the ledger in the collection with
// collection's name. This collection will be private unless the collection
// name is WorldStateIdentifier
func (c *Collection) Add(entry interface{}) error {
	keyComponents, err := c.generateKey(entry)

	if err != nil {
		return fmt.Errorf("Failed to add to collection. %s", err.Error())
	}

	key, err := c.formatKey(keyComponents)

	if err != nil {
		return fmt.Errorf("Failed to add to collection. %s", err.Error())
	}

	exists, err := c.exists(key)

	if err != nil {
		return fmt.Errorf("Failed to add to collection. %s", err.Error())
	}

	if exists {
		return fmt.Errorf("Failed to add to collection. Key already exists")
	}

	bytes, err := c.Serializer.ToBytes(entry)

	if err != nil {
		return fmt.Errorf("Failed to add to collection. %s", err.Error())
	}

	if c.Name != WorldStateIdentifier {
		err = c.Stub.PutPrivateData(c.Name, key, bytes)
	} else {
		err = c.Stub.PutState(key, bytes)
	}

	if err != nil {
		return fmt.Errorf("Failed to add to collection. %s", err.Error())
	}

	return nil
}

// Get gets a serialized version of the entry from the ledger in the collection
// with the collection's name. the entry is then deserialized into the passed
// interface. This collection will be private unless the
// collecton name is WorldStateIdentifier
func (c *Collection) Get(entry interface{}, keyComponents ...string) error {
	key, err := c.formatKey(keyComponents)

	if err != nil {
		return fmt.Errorf("Failed to get from collection. %s", err.Error())
	}

	var bytes []byte

	if c.Name != WorldStateIdentifier {
		bytes, err = c.Stub.GetPrivateData(c.Name, key)
	} else {
		bytes, err = c.Stub.GetState(key)
	}

	if err != nil {
		return fmt.Errorf("Failed to get from collection. %s", err.Error())
	} else if bytes == nil {
		return fmt.Errorf("Failed to get from collection. Key does not exist")
	}

	err = c.Serializer.FromBytes(bytes, entry)

	if err != nil {
		return fmt.Errorf("Failed to get from collection. %s", err.Error())
	}

	return nil
}

// GetHistory returns an iterator of historic data for the key. Only available for world
// state collection
func (c *Collection) GetHistory(keyComponents ...string) (HistoryQueryIteratorInterface, error) {
	key, err := c.formatKey(keyComponents)

	if err != nil {
		return nil, fmt.Errorf("Failed to get history from collection. %s", err.Error())
	}

	var shqi shim.HistoryQueryIteratorInterface

	if c.Name != WorldStateIdentifier {
		return nil, fmt.Errorf("Failed to get history from collection. %s", "Historic data does not exist for non world state collections")
	}

	shqi, err = c.Stub.GetHistoryForKey(key)

	if err != nil {
		return nil, fmt.Errorf("Failed to get history from collection. %s", err.Error())
	}

	hqi := new(HistoryQueryIterator)
	hqi.Serializer = c.Serializer
	hqi.Iterator = shqi

	return hqi, nil
}

// Update updates an existing version of the state in the ledger in the collection
// with collection's name. This collection will be private unless the collection
// name is WorldStateIdentifier
func (c *Collection) Update(entry interface{}) error {
	keyComponents, err := c.generateKey(entry)

	if err != nil {
		return fmt.Errorf("Failed to update in collection. %s", err.Error())
	}

	key, err := c.formatKey(keyComponents)

	if err != nil {
		return fmt.Errorf("Failed to update in collection. %s", err.Error())
	}

	exists, err := c.exists(key)

	if err != nil {
		return fmt.Errorf("Failed to update in collection. %s", err.Error())
	}

	if !exists {
		return fmt.Errorf("Failed to update in collection. Key already exists")
	}

	bytes, err := c.Serializer.ToBytes(entry)

	if err != nil {
		return fmt.Errorf("Failed to update in collection. %s", err.Error())
	}

	if c.Name != WorldStateIdentifier {
		err = c.Stub.PutPrivateData(c.Name, key, bytes)
	} else {
		err = c.Stub.PutState(key, bytes)
	}

	if err != nil {
		return fmt.Errorf("Failed to update in collection. %s", err.Error())
	}

	return nil
}

// Delete removes the state from the collection in the ledger with
// the given key
func (c *Collection) Delete(keyComponents ...string) error {
	key, err := c.formatKey(keyComponents)

	if err != nil {
		return fmt.Errorf("Failed to delete from collection. %s", err.Error())
	}

	exists, err := c.exists(key)

	if err != nil {
		return fmt.Errorf("Failed to delete from collection. %s", err.Error())
	}

	if !exists {
		return fmt.Errorf("Failed to delete from collection. Key does not exist")
	}

	if c.Name != WorldStateIdentifier {
		err = c.Stub.DelPrivateData(c.Name, key)
	} else {
		err = c.Stub.DelState(key)
	}

	if err != nil {
		return fmt.Errorf("Failed to delete from collection. %s", err.Error())
	}

	return nil
}
