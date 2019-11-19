package main

import (
	"errors"
	"fmt"

	"github.com/awjh-ibm/fabric-contract-api-go-ledger-api/ledgerapi"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func genArgs(funcName string, params ...string) [][]byte {
	argBytes := [][]byte{}
	argBytes = append(argBytes, []byte(funcName))

	for _, param := range params {
		argBytes = append(argBytes, []byte(param))
	}

	return argBytes
}

func submitTransaction(stub *shimtest.MockStub, txID string, funcName string, params ...string) error {
	stub.MockTransactionStart(txID)
	response := stub.MockInvoke(txID, genArgs(funcName, params...))
	stub.MockTransactionEnd(txID)

	if response.GetStatus() != 200 {
		fmt.Println("ERROR: ", response.GetMessage())
		return errors.New(response.GetMessage())
	}

	fmt.Println("SUCCESS: ", string(response.GetPayload()))

	return nil
}

func main() {
	contract := new(Contract)
	contract.TransactionContextHandler = new(ledgerapi.TransactionContext)
	contract.Name = "org.papernet.commercialpaper"
	contract.Info.Version = "0.0.1"

	chaincode, err := contractapi.NewChaincode(contract)

	if err != nil {
		fmt.Printf("Error creating chaincode. %s", err)
		return
	}

	chaincode.Info.Title = "CommercialPaperChaincode"
	chaincode.Info.Version = "0.0.1"

	stub := shimtest.NewMockStub("myStub", chaincode)

	fmt.Println("=========== ISSUE ===========")
	err = submitTransaction(stub, "TX1", "Issue", "MAGNETOCORP", "0001", "2019-11-15", "2020-11-15", "100000")

	if err != nil {
		return
	}

	fmt.Println("=========== BUY ===========")
	err = submitTransaction(stub, "TX2", "Buy", "MAGNETOCORP", "0001", "MAGNETOCORP", "DIGIBANK")

	if err != nil {
		return
	}

	fmt.Println("=========== REDEEM ===========")
	err = submitTransaction(stub, "TX3", "Redeem", "MAGNETOCORP", "0001", "DIGIBANK")

	if err != nil {
		return
	}

	// fmt.Println("=========== HISTORY ===========")
	// err = submitTransaction(stub, "TX4", "GetPaperTransfers", "MAGNETOCORP", "0001")

	// if err != nil {
	// 	return
	// }
}
