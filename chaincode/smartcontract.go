package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type TransactionContract struct {
	contractapi.Contract
}

//* Contract Methods

func (t *TransactionContract) InitiateTransaction(ctx contractapi.TransactionContextInterface, transactionJson string) error {

	var transactionInitiate TransactionInitiate

	err := json.Unmarshal([]byte(transactionJson), &transactionInitiate)
	if err != nil {
		return fmt.Errorf("failed to unmarshal Json: %v", err)
	}
	transactionId := string(transactionInitiate.Id)
	fmt.Printf("Deserialized data: %+v\n", transactionInitiate)

	exists, err := t.TrxExists(ctx, string(transactionId))
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("transaction %s already exists", transactionId)

	}
	err = transactionInitiate.setCreatedAt(ctx)
	if err != nil {
		return err
	}
	trxJson, err := json.Marshal(transactionInitiate)

	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(transactionJson, trxJson)

}

func (t *TransactionContract) GetAllTransactions(ctx contractapi.TransactionContextInterface) ([]*TransactionInitiate, error) {
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
func (t *TransactionContract) GetTransactionById(ctx contractapi.TransactionContextInterface,
	id string) (*TransactionInitiate, error) {
	transactionJson, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, err
	}
	if transactionJson == nil {
		return nil, fmt.Errorf("transaction %s does not exist", id)
	}
	transaction := new(TransactionInitiate)
	err = json.Unmarshal(transactionJson, transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (t *TransactionContract) UpdateTransactionProcess(ctx contractapi.TransactionContextInterface, updatedTransactionProcessJson string) error {

	var updateTransactionProcess UpdateTransactionProcess

	err := json.Unmarshal([]byte(updatedTransactionProcessJson), &updateTransactionProcess)

	if err != nil {
		return err
	}

	transactionId := string(updateTransactionProcess.Id)

	transactionData, err := t.GetTransactionById(ctx, transactionId)
	if err != nil {
		return err
	}

	transactionData.TransactionPhase = updateTransactionProcess.TransactionPhase

	for i, trx := range transactionData.TransactionProcesses {
		for _, trxProcess := range updateTransactionProcess.TransactionProcesses {
			if trx.Id == trxProcess.Id {
				transactionData.TransactionProcesses = append(transactionData.TransactionProcesses[:i], trx)
			}
		}
	}
	err = transactionData.setCreatedAt(ctx)
	if err != nil {
		return err
	}
	trxJson, err := json.Marshal(transactionData)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(transactionId, trxJson)

}

func (t *TransactionContract) GetTransactionHistory(ctx contractapi.TransactionContextInterface, id string) ([]*TransactionInitiate, error) {
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

func (t *TransactionContract) AddRecipientPayment(ctx contractapi.TransactionContextInterface, recipientJson string) error {

	var recipientGroupPayment RecipientGroupPayment
	err := json.Unmarshal([]byte(recipientJson), &recipientGroupPayment)
	if err != nil {
		return err
	}
	transactionId := string(recipientGroupPayment.Id)
	transaction, err := t.GetTransactionById(ctx, transactionId)
	if err != nil {
		return err
	}
	transaction.TransferredAmount = &recipientGroupPayment.TransferredAmount
	transaction.TransactionStatus = &recipientGroupPayment.TransactionStatus

	transaction.addTransactionProcess(recipientGroupPayment.TransactionProcess)

	transaction.TransactionPayment = append(transaction.TransactionPayment, recipientGroupPayment.TransactionPayment)
	err = transaction.setCreatedAt(ctx)
	if err != nil {
		return err
	}
	trxJson, err := json.Marshal(transaction)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(transactionId, trxJson)

}

func (t *TransactionContract) AddMemberPayment(ctx contractapi.TransactionContextInterface, memberPaymentJson string) error {
	var memberPayment UpdateMemberTransaction
	err := json.Unmarshal([]byte(memberPaymentJson), &memberPayment)
	if err != nil {
		return err
	}
	transactionId := string(memberPayment.Id)
	transaction, err := t.GetTransactionById(ctx, transactionId)
	if err != nil {
		return err
	}
	transaction.TransactionPhase = memberPayment.TransactionPhase
	transaction.addTransactionProcess(memberPayment.TransactionProcess)

	for _, memberTrx := range transaction.TransactionMembers {

		if memberTrx.Id == memberPayment.Id {
			transaction.TransactionMembers = append(transaction.TransactionMembers, memberPayment.TransactionMember)
			break
		}
	}
	err = transaction.setCreatedAt(ctx)
	if err != nil {
		return err
	}
	txnJson, err := json.Marshal(transaction)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(transactionId, txnJson)

}

func (t *TransactionInitiate) addTransactionProcess(transactionProcess TransactionProcess) {

	for i, trx := range t.TransactionProcesses {
		if trx.Id == transactionProcess.Id {
			t.TransactionProcesses = append(t.TransactionProcesses[:i], transactionProcess)
			break
		}
	}

}

func (t *TransactionInitiate) setCreatedAt(ctx contractapi.TransactionContextInterface) error {
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return err
	}
	t.CreatedAt = time.Unix(timestamp.Seconds, int64(timestamp.Nanos))
	return nil
}

func (t *TransactionContract) TrxExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJson, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read  state: %v", err)
	}

	return assetJson != nil, nil
}
