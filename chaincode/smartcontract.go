package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type GroupTransaction struct {
	ID        string    `json:"ID"`
	GroupId   int       `json:"GroupId"`
	ProjectId int       `json:"ProjectId"`
	DonorId   int       `json:"DonorId"`
	Amount    float64   `json:"Amount"`
	CreatedAt time.Time `json:"CreatedAt"`
}

func (s *SmartContract) CreateGroupTransaction(ctx contractapi.TransactionContextInterface, id string, groupId, projectId int, amount float64, donorId int) error {
	exists, err := s.TransactionExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("group transaction %s already exists", id)

	}
	transaction := GroupTransaction{
		ID:        id,
		GroupId:   groupId,
		ProjectId: projectId,
		Amount:    amount,
		DonorId:   donorId,
		CreatedAt: time.Now(),
	}
	transactionJson, err := json.Marshal(transaction)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, transactionJson)
}

func (s *SmartContract) GetTransactionById(ctx contractapi.TransactionContextInterface, id string) (*GroupTransaction, error) {
	transactionJson, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, err
	}
	if transactionJson == nil {
		return nil, fmt.Errorf("transaction %s does not exist", id)
	}
	transaction := new(GroupTransaction)
	err = json.Unmarshal(transactionJson, transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (s *SmartContract) TransactionExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJson, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return assetJson != nil, nil
}

func (s *SmartContract) GetAllTransactions(ctx contractapi.TransactionContextInterface) ([]GroupTransaction, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer func(resultsIterator shim.StateQueryIteratorInterface) {
		err := resultsIterator.Close()
		if err != nil {

		}
	}(resultsIterator)
	var transactions []GroupTransaction
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var transaction GroupTransaction
		err = json.Unmarshal(queryResponse.Value, &transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (s *SmartContract) GetTransactionByGroupId(ctx contractapi.TransactionContextInterface, groupId int) ([]GroupTransaction, error) {
	queryString := fmt.Sprintf(`{"selector":{"GroupId":%d}}`, groupId)
	resultIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer func(resultIterator shim.StateQueryIteratorInterface) {
		err := resultIterator.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resultIterator)
	var transactions []GroupTransaction
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}
		var transaction GroupTransaction
		err = json.Unmarshal(queryResponse.Value, &transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (s *SmartContract) GetTransactionsByProjectId(ctx contractapi.TransactionContextInterface, projectId int) ([]GroupTransaction, error) {
	queryString := fmt.Sprintf(`{"selector":{"ProjectId":%d}}`, projectId)
	resultIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer func(resultIterator shim.StateQueryIteratorInterface) {
		err := resultIterator.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resultIterator)
	var transactions []GroupTransaction
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}
		var transaction GroupTransaction
		err = json.Unmarshal(queryResponse.Value, &transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (s *SmartContract) GetTransactionsByDonorId(ctx contractapi.TransactionContextInterface, donorId int) ([]GroupTransaction, error) {
	queryString := fmt.Sprintf(`{"selector":{"DonorId":%d}}`, donorId)
	resultIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer func(resultIterator shim.StateQueryIteratorInterface) {
		err := resultIterator.Close()
		if err != nil {

		}
	}(resultIterator)
	var transactions []GroupTransaction
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}
		var transaction GroupTransaction
		err = json.Unmarshal(queryResponse.Value, &transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}
