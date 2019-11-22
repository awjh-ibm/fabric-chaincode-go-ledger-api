package blockchainapi

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

func GetCurrentTransaction(ctx contractapi.TransactionContextInterface) *Transaction {
	transaction := new(Transaction)
	transaction.Ctx = ctx

	return transaction
}
