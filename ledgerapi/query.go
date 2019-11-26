package ledgerapi

import (
	"errors"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// QueryType types of queries that can be done
type QueryType int

// QueryInterface functions that a valid query must have
type QueryInterface interface {
	Query(contractapi.TransactionContextInterface, string) (*QueryStateIterator, error)
}

// QueryPagination stores pagination settings
type QueryPagination struct {
	PageSize int32
	Bookmark string
}

// RangeQuery inplementation of QueryInterface for
// range queries
type RangeQuery struct {
	FromKey    string
	ToKey      string
	Pagination *QueryPagination
}

// Query runs a range query. If using pagination bookmark gets updated after query
func (rq *RangeQuery) Query(ctx contractapi.TransactionContextInterface, collection string) (*QueryStateIterator, error) {
	var result shim.StateQueryIteratorInterface
	var metadata *peer.QueryResponseMetadata
	var err error

	if collection == WorldStateCollection {
		if rq.Pagination != nil {
			result, metadata, err = ctx.GetStub().GetStateByRangeWithPagination(rq.FromKey, rq.ToKey, rq.Pagination.PageSize, rq.Pagination.Bookmark)

			rq.Pagination.Bookmark = metadata.GetBookmark()
		} else {
			result, err = ctx.GetStub().GetStateByRange(rq.FromKey, rq.ToKey)
		}
	} else {
		if rq.Pagination != nil {
			return nil, errors.New("Pagination not implemented for private collections")
		}

		result, err = ctx.GetStub().GetPrivateDataByRange(collection, rq.FromKey, rq.ToKey)
	}

	if err != nil {
		return nil, err
	}

	qsi := new(QueryStateIterator)
	qsi.StateQueryIteratorInterface = result
	qsi.Ctx = ctx
	qsi.Collection = collection

	return qsi, nil
}

// TODO Other query types
