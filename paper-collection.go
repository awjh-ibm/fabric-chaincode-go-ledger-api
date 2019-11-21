package main

import (
	"github.com/awjh-ibm/fabric-contract-api-go-ledger-api/ledgerapi"
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

// PaperCollection abstraction of a collection for managing papers
type PaperCollection struct {
	Name             string
	Stub             shim.ChaincodeStubInterface
	ledgerCollection *ledgerapi.Collection
}

func (pc *PaperCollection) setupLedgerCollection() {
	ledgerCollection := new(ledgerapi.Collection)
	ledgerCollection.Name = pc.Name
	ledgerCollection.Stub = pc.Stub
	ledgerCollection.Serializer = new(ledgerapi.JSONLedgerSerializer)

	pc.ledgerCollection = ledgerCollection
}

// Exists returns true if paper exists in the world state
func (pc *PaperCollection) Exists(issuer, paperNumber string) (bool, error) {
	pc.setupLedgerCollection()

	return pc.ledgerCollection.Exists(issuer, paperNumber)
}

// Add adds a commercial paper to the collection
func (pc *PaperCollection) Add(paper *CommercialPaper) error {
	pc.setupLedgerCollection()

	return pc.ledgerCollection.Add(paper)
}

// Get returns a commercial paper from the collection
func (pc *PaperCollection) Get(issuer, paperNumber string) (*CommercialPaper, error) {
	pc.setupLedgerCollection()

	paper := new(CommercialPaper)

	err := pc.ledgerCollection.Get(paper, issuer, paperNumber)

	if err != nil {
		return nil, err
	}

	return paper, nil
}

// Update updates a commercial paper in the collection
func (pc *PaperCollection) Update(paper *CommercialPaper) error {
	pc.setupLedgerCollection()

	return pc.ledgerCollection.Update(paper)
}

// Delete updates a commercial paper in the collection
func (pc *PaperCollection) Delete(issuer, paperNumber string) error {
	pc.setupLedgerCollection()

	return pc.ledgerCollection.Delete(issuer, paperNumber)
}
