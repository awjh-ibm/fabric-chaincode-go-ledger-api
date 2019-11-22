package blockchainapi

type EventInterface interface {
	GetName() string
	GetPayloadBytes() ([]byte, error)
}

type Event struct {
	Name    string
	Payload interface{}
}

func (e *Event) GetName() string {
	return e.Name
}

func (e *Event) GetPayloadBytes() ([]byte, error) {
	return new(JSONLedgerSerializer).ToBytes(e)
}
