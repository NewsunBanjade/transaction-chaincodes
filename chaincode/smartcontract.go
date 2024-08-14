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

type ValueId struct {
	Id   int    `json:"Id"`
	Name string `json:"Name"`
}

type TransactionProcess struct {
	ValueId
	TransactionId      int     `json:"TransactionId"`
	Remarks            int     `json:"Remarks"`
	StatusModifiedDate int     `json:"StatusModifiedDate"`
	ApproverType       ValueId `json:"ApproverType"`
	TransactionStatus  ValueId `json:"TransactionStatus"`
}

type TransactionMember struct {
	Id            int     `json:"Id"`
	TransactionId int     `json:"TransactionId"`
	MemberId      int     `json:"MemberId"`
	Amount        float64 `json:"Amount"`
	Member        Member  `json:"Member"`
}

type Member struct {
	ValueId
	Code string `json:"code"`
}

type BankAccount struct {
	Id            int    `json:"Id"`
	BankName      string `json:"BankName"`
	AccountNumber string `json:"AccountNumber"`
}

type TransactionInitiate struct {
	Id                        string               `json:"Id"`
	Guid                      string               `json:"Guid"`
	Amount                    float64              `json:"Amount"`
	Remarks                   string               `json:"Remarks"`
	RecipientGroup            ValueId              `json:"RecipientGroup"`
	Project                   ValueId              `json:"Project"`
	BankAccount               BankAccount          `json:"BankAccount"`
	RecipientGroupBankAccount BankAccount          `json:"RecipientGroupBankAccount"`
	TransactionPhase          ValueId              `json:"TransactionPhase"`
	TransactionProcesses      []TransactionProcess `json:"TransactionProcesses"`
	TransactionMembers        []TransactionMember  `json:"TransactionMembers"`
	CreatedAt                 time.Time            `json:"CreatedAt"`
}

func (t *TransactionContract) CreateTransaction(ctx contractapi.TransactionContextInterface, id string, guid, remarks string,
	amount float64, recipientGroup, project, transactionPhase ValueId, bankAcc, recipientGroupBankAcc BankAccount, transactionProcesses []TransactionProcess, transactionMembers []TransactionMember) error {

	exists, err := t.TrxExists(ctx, id)

	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("transaction %s already exists", id)

	}
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return err
	}
	createdAt := time.Unix(timestamp.Seconds, int64(timestamp.Nanos))

	transaction := TransactionInitiate{
		Id:                        id,
		Guid:                      guid,
		Amount:                    amount,
		Remarks:                   remarks,
		RecipientGroup:            recipientGroup,
		Project:                   project,
		BankAccount:               bankAcc,
		RecipientGroupBankAccount: recipientGroupBankAcc,
		TransactionPhase:          transactionPhase,
		TransactionProcesses:      transactionProcesses,
		TransactionMembers:        transactionMembers,
		CreatedAt:                 createdAt,
	}

	trxJson, err := json.Marshal(transaction)

	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, trxJson)

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

func (t *TransactionContract) TrxExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJson, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return assetJson != nil, nil
}

func (t *TransactionContract) UpdateTransactionProcess(ctx contractapi.TransactionContextInterface, id string, transactionPhase ValueId, transactionProcess TransactionProcess) error {
	transactionData, err := t.GetTransactionById(ctx, id)
	if err != nil {
		return err
	}
	transactionData.TransactionPhase = transactionPhase
	transactionProcessData := transactionData.TransactionProcesses

	for i, trx := range transactionProcessData {
		if trx.Id == transactionProcess.Id {
			transactionProcessData = append(transactionProcessData[:i], transactionProcess)
			break
		}
	}
	transactionData.TransactionProcesses = transactionProcessData
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return err
	}
	createdAt := time.Unix(timestamp.Seconds, int64(timestamp.Nanos))
	transactionData.CreatedAt = createdAt

	trxJson, err := json.Marshal(transactionData)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, trxJson)

}
