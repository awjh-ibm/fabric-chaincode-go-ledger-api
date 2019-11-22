package blockchainapi

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type Transaction struct {
	Ctx contractapi.TransactionContextInterface
}

func (t *Transaction) GetArgs() []string {
	return t.Ctx.GetStub().GetStringArgs()
}

func (t *Transaction) GetTransient() (map[string][]byte, error) {
	return t.Ctx.GetStub().GetTransient()
}

func (t *Transaction) GetTxID() string {
	return t.Ctx.GetStub().GetTxID()
}

func (t *Transaction) GetTxTimestamp() (*timestamp.Timestamp, error) {
	return t.Ctx.GetStub().GetTxTimestamp()
}

func (t *Transaction) GetCreator() cid.ClientIdentity {
	return t.Ctx.GetClientIdentity()
}

func (t *Transaction) GetDecorations() map[string][]byte {
	return t.Ctx.GetStub().GetDecorations()
}

func (t *Transaction) GetBinding() ([]byte, error) {
	return t.Ctx.GetStub().GetBinding()
}

func (t *Transaction) GetSignedProposal() (*peer.SignedProposal, error) {
	return t.Ctx.GetStub().GetSignedProposal()
}

func (t *Transaction) SetEvent(evt EventInterface) error {
	bytes, err := evt.GetPayloadBytes()

	if err != nil {
		return err
	}

	return t.Ctx.GetStub().SetEvent(evt.GetName(), bytes)
}
