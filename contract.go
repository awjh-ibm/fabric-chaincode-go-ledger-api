package main

import (
	"fmt"

	"github.com/awjh-ibm/fabric-contract-api-go-ledger-api/ledgerapi"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Contract chaincode that defines
// the business logic for managing commercial
// paper
type Contract struct {
	contractapi.Contract
}

// Instantiate does nothing
func (c *Contract) Instantiate() {
	fmt.Println("Instantiated")
}

// Issue creates a new commercial paper and stores it in the world state
func (c *Contract) Issue(ctx ledgerapi.TransactionContextInterface, issuer string, paperNumber string, issueDateTime string, maturityDateTime string, faceValue int) (*CommercialPaper, error) {
	paper := CommercialPaper{PaperNumber: paperNumber, Issuer: issuer, IssueDateTime: issueDateTime, FaceValue: faceValue, MaturityDateTime: maturityDateTime, Owner: issuer}
	paper.SetIssued()

	err := ctx.GetLedger().GetDefaultCollection().Add(paper)

	if err != nil {
		return nil, err
	}

	return &paper, nil
}

// Buy updates a commercial paper to be in trading status and sets the new owner
func (c *Contract) Buy(ctx ledgerapi.TransactionContextInterface, issuer string, paperNumber string, currentOwner string, newOwner string) (*CommercialPaper, error) {
	paper := new(CommercialPaper)

	err := ctx.GetLedger().GetDefaultCollection().Get(paper, issuer, paperNumber)

	if err != nil {
		return nil, err
	}

	if paper.Owner != currentOwner {
		return nil, fmt.Errorf("Paper %s:%s is not owned by %s", issuer, paperNumber, currentOwner)
	}

	if paper.IsIssued() {
		paper.SetTrading()
	}

	if !paper.IsTrading() {
		return nil, fmt.Errorf("Paper %s:%s is not trading. Current state = %s", issuer, paperNumber, paper.GetState())
	}

	paper.Owner = newOwner

	err = ctx.GetLedger().GetDefaultCollection().Update(paper)

	if err != nil {
		return nil, err
	}

	return paper, nil
}

// Redeem updates a commercial paper status to be redeemed
func (c *Contract) Redeem(ctx ledgerapi.TransactionContextInterface, issuer string, paperNumber string, redeemingOwner string) (*CommercialPaper, error) {
	paper := new(CommercialPaper)

	err := ctx.GetLedger().GetDefaultCollection().Get(paper, issuer, paperNumber)

	if err != nil {
		return nil, err
	}

	if paper.Owner != redeemingOwner {
		return nil, fmt.Errorf("Paper %s:%s is not owned by %s", issuer, paperNumber, redeemingOwner)
	}

	if paper.IsRedeemed() {
		return nil, fmt.Errorf("Paper %s:%s is already redeemed", issuer, paperNumber)
	}

	paper.Owner = paper.Issuer
	paper.SetRedeemed()

	err = ctx.GetLedger().GetDefaultCollection().Update(paper)

	if err != nil {
		return nil, err
	}

	return paper, nil
}

// GetPaperTransfers returns the array of transaction IDs of when transfers happened to the paper
func (c *Contract) GetPaperTransfers(ctx ledgerapi.TransactionContextInterface, issuer string, paperNumber string) ([]string, error) {
	history := []string{}

	it, err := ctx.GetLedger().GetDefaultCollection().GetHistory(issuer, paperNumber)

	if err != nil {
		return nil, err
	}

	defer it.Close()
	for it.HasNext() {
		nxt, _ := it.Next()

		paper := new(CommercialPaper)

		err := nxt.GetValue(paper)

		if err != nil {
			return nil, err
		}

		if paper.IsTrading() {
			history = append(history, nxt.GetTxID())
		}
	}

	return history, nil
}
