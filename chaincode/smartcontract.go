package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) CreateTransaction(ctx contractapi.TransactionContextInterface, assetJson string) error {

	var transaction TransactionInitiate
	err := json.Unmarshal([]byte(assetJson), &transaction)
	if err != nil {
		return err
	}

	id := strconv.Itoa(transaction.Id)

	exists, err := s.TransactionExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}
	transaction.CreatedAt, err = transaction.getTxTimeStamp(ctx)
	if err != nil {
		return err
	}
	transaction.ModifiedAt = transaction.CreatedAt

	assetJSON, err := json.Marshal(&transaction)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

func (s *SmartContract) ReadTransactionById(ctx contractapi.TransactionContextInterface, id string) (*TransactionInitiate, error) {
	transactionJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if transactionJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var transaction TransactionInitiate
	err = json.Unmarshal(transactionJSON, &transaction)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (s *SmartContract) UpdateTransactionProcess(ctx contractapi.TransactionContextInterface, tnxJson string) error {
	var updateTnx UpdateTransactionProcess
	err := json.Unmarshal([]byte(tnxJson), &updateTnx)
	if err != nil {
		return err
	}

	id := strconv.Itoa(updateTnx.Id)
	exists, err := s.TransactionExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("transaction doesn't %s  exists", id)
	}

	tnxData, err := s.ReadTransactionById(ctx, id)
	if err != nil {
		return err
	}
	tnxData.TransactionPhase = updateTnx.TransactionPhase
	for i, tnx := range tnxData.TransactionProcesses {
		for _, updatedTnxProcess := range updateTnx.TransactionProcesses {
			if tnx.Id == updatedTnxProcess.Id {
				newTrxProcess := tnx
				newTrxProcess.Remarks = updatedTnxProcess.Remarks
				newTrxProcess.StatusModifiedDate = updatedTnxProcess.StatusModifiedDate
				newTrxProcess.TransactionStatus = updatedTnxProcess.TransactionStatus
				tnxData.TransactionProcesses[i] = newTrxProcess
			}
		}
	}
	tnxData.ModifiedAt, err = tnxData.getTxTimeStamp(ctx)
	if err != nil {
		return err
	}
	trxJson, err := json.Marshal(tnxData)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(id, trxJson)
	if err != nil {
		return err
	}
	return nil

}

func (s *SmartContract) GetTransactionHistoryById(ctx contractapi.TransactionContextInterface, id string) ([]*TransactionInitiate, error) {
	transactionIterator, err := ctx.GetStub().GetHistoryForKey(id)
	if err != nil {
		return nil, err
	}
	defer transactionIterator.Close()
	var transactions []*TransactionInitiate
	for transactionIterator.HasNext() {
		queryResponse, err := transactionIterator.Next()
		if err != nil {
			return nil, err
		}
		var transaction TransactionInitiate
		err = json.Unmarshal(queryResponse.Value, &transaction)
		if err != nil {
			return nil, err

		}
		transactions = append(transactions, &transaction)
	}
	return transactions, nil
}

func (s *SmartContract) TransactionExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

func (s *SmartContract) GetAllTransactions(ctx contractapi.TransactionContextInterface) ([]*TransactionInitiate, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var transactions []*TransactionInitiate
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var transaction TransactionInitiate
		err = json.Unmarshal(queryResponse.Value, &transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

func (s *SmartContract) AddRecipientPayment(ctx contractapi.TransactionContextInterface, recipientJson string) error {

	var recipientGroupPayment RecipientGroupPayment
	err := json.Unmarshal([]byte(recipientJson), &recipientGroupPayment)
	if err != nil {
		return err
	}
	id := strconv.Itoa(recipientGroupPayment.Id)
	transaction, err := s.ReadTransactionById(ctx, id)
	if err != nil {
		return err
	}
	transaction.TransferredAmount = recipientGroupPayment.TransferredAmount
	transaction.TransactionStatus = &recipientGroupPayment.TransactionStatus
	transaction.addTransactionProcess(recipientGroupPayment.TransactionProcess)
	transaction.TransactionPayments = append(transaction.TransactionPayments, recipientGroupPayment.TransactionPayment)
	transaction.ModifiedAt, err = transaction.getTxTimeStamp(ctx)
	if err != nil {
		return err
	}
	trxJson, err := json.Marshal(transaction)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(id, trxJson)
	if err != nil {
		return err
	}
	return nil

}

func (s *SmartContract) AddMemberPayment(ctx contractapi.TransactionContextInterface, memberJson string) error {

	var updatedMemberPayment UpdateMemberTransaction
	err := json.Unmarshal([]byte(memberJson), &updatedMemberPayment)
	if err != nil {
		return err
	}
	id := strconv.Itoa(updatedMemberPayment.Id)
	transaction, err := s.ReadTransactionById(ctx, id)
	if err != nil {
		return err
	}
	transaction.TransactionPhase = updatedMemberPayment.TransactionPhase
	transaction.addTransactionProcess(updatedMemberPayment.TransactionProcess)

	for i, member := range transaction.TransactionMembers {
		if member.Id == updatedMemberPayment.TransactionMember.Id {
			newMemberTransaction := member
			newMemberTransaction.TransferredAmount = updatedMemberPayment.TransactionMember.TransferredAmount
			newMemberTransaction.TransactionStatus = updatedMemberPayment.TransactionMember.TransactionStatus

			newMemberTransaction.TransactionMemberPayments = append(newMemberTransaction.TransactionMemberPayments, updatedMemberPayment.TransactionMember.TransactionMemberPayment)
			transaction.TransactionMembers[i] = newMemberTransaction

			break
		}
	}
	transaction.ModifiedAt, err = transaction.getTxTimeStamp(ctx)
	if err != nil {
		return err
	}

	trxJson, err := json.Marshal(transaction)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(id, trxJson)
	if err != nil {
		return err
	}
	return nil

}

func (t *TransactionInitiate) addTransactionProcess(transactionProcess TransactionProcess) {

	for i, trx := range t.TransactionProcesses {
		if trx.Id == transactionProcess.Id {
			newTrxProcess := trx
			newTrxProcess.Remarks = transactionProcess.Remarks
			newTrxProcess.StatusModifiedDate = transactionProcess.StatusModifiedDate
			newTrxProcess.TransactionStatus = transactionProcess.TransactionStatus
			t.TransactionProcesses[i] = newTrxProcess
			break
		}
	}

}

func (t *TransactionInitiate) getTxTimeStamp(ctx contractapi.TransactionContextInterface) (time.Time, error) {
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return time.Now(), err
	}
	return time.Unix(timestamp.Seconds, int64(timestamp.Nanos)), nil

}
