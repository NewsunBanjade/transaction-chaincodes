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
type TransactionInitiate struct {
	Id                   int                  `json:"Id"`
	Name                 string               `json:"Name"`
	Guid                 string               `json:"Guid"`
	Amount               float64              `json:"Amount"`
	TransferredAmount    *float64             `json:"TransferredAmount"`
	RecipientGroup       RecipientGroup       `json:"RecipientGroup"`
	Project              Project              `json:"Project"`
	TransactionPhase     TransactionPhase     `json:"TransactionPhase"`
	TransactionStatus    *TransactionStatus   `json:"TransactionStatus"`
	TransactionProcesses []TransactionProcess `json:"TransactionProcesses"`
	TransactionMembers   []TransactionMember  `json:"TransactionMembers"`
	TransactionPayments  []TransactionPayment `json:"TransactionPayments"`
	CreatedAt            time.Time            `json:"CreatedAt"`
}

type RecipientGroup struct {
	Id   int    `json:"Id"`
	Name string `json:"Name"`
}

type Project struct {
	Id   int    `json:"Id"`
	Name string `json:"Name"`
}

type TransactionPhase struct {
	Id   int    `json:"Id"`
	Name string `json:"Name"`
}
type TransactionStatus struct {
	Id   int    `json:"Id"`
	Name string `json:"Name"`
}
type TransactionForwardPurpose struct {
	Id   int    `json:"Id"`
	Name string `json:"Name"`
}

type ApproverType struct {
	Id   int    `json:"Id"`
	Name string `json:"Name"`
}

type PaymentMethod struct {
	Id   int    `json:"Id"`
	Name string `json:"Name"`
}
type TransactionProcess struct {
	Id                        int                       `json:"Id"`
	Name                      string                    `json:"Name"`
	TransactionId             int                       `json:"TransactionId"`
	MemberId                  *int                      `json:"MemberId"`
	Remarks                   *string                   `json:"Remarks"`
	StatusModifiedDate        *string                   `json:"StatusModifiedDate "`
	ApproverType              ApproverType              `json:"ApproverType"`
	TransactionForwardPurpose TransactionForwardPurpose `json:"TransactionForwardPurpose"`
	TransactionStatus         TransactionStatus         `json:"TransactionStatus"`
}
type Member struct {
	Id   int    `json:"Id"`
	Name string `json:"Name"`
	Code string `json:"Code"`
}

type BankAccount struct {
	Id            int    `json:"Id"`
	BankName      string `json:"BankName"`
	AccountNumber string `json:"AccountNumber"`
}

type TransactionPayment struct {
	Id                        int           `json:"Id"`
	TransactionId             int           `json:"TransactionId"`
	PaymentReference          string        `json:"PaymentReference"`
	PaymentAmount             float64       `json:"PaymentAmount"`
	PaymentDate               string        `json:"PaymentDate"`
	Remarks                   *string       `json:"Remarks "`
	PaymentMethod             PaymentMethod `json:"PaymentMethod"`
	BankAccount               BankAccount   `json:"BankAccount"`
	RecipientGroupBankAccount BankAccount   `json:"RecipientGroupBankAccount"`
}

type TransactionMemberPayment struct {
	Id                  int           `json:"Id"`
	TransactionMemberId int           `json:"TransactionMemberId"`
	PaymentReference    string        `json:"PaymentReference"`
	PaymentAmount       float64       `json:"PaymentAmount"`
	PaymentDate         string        `json:"PaymentDate"`
	Remarks             *string       `json:"Remarks"`
	PaymentMethod       PaymentMethod `json:"PaymentMethod"`
}

type TransactionMember struct {
	Id                        int                        `json:"Id"`
	TransactionId             int                        `json:"TransactionId"`
	MemberId                  int                        `json:"MemberId"`
	Amount                    float64                    `json:"Amount"`
	TransferredAmount         *float64                   `json:"TransferredAmount"`
	TransactionStatus         *TransactionStatus         `json:"TransactionStatus"`
	Member                    Member                     `json:"Member"`
	TransactionMemberPayments []TransactionMemberPayment `json:"TransactionMemberPayments"`
}

//* Update Struct's

type UpdateTransactionProcess struct {
	Id                   int                  `json:"Id"`
	TransactionPhase     TransactionPhase     `json:"TransactionPhase"`
	TransactionProcesses []TransactionProcess `json:"TransactionProcesses"`
}

type RecipientGroupPayment struct {
	Id                 int                `json:"id"`
	TransferredAmount  float64            `json:"TransferredAmount"`
	TransactionStatus  TransactionStatus  `json:"TransactionStatus"`
	TransactionProcess TransactionProcess `json:"TransactionProcess"`
	TransactionPayment TransactionPayment `json:"TransactionPayment"`
}

type UpdateMemberTransaction struct {
	Id                 int                `json:"id"`
	TransactionPhase   TransactionPhase   `json:"TransactionPhase"`
	TransactionProcess TransactionProcess `json:"TransactionProcess"`
	TransactionMember  TransactionMember  `json:"TransactionMember"`
}

//* Contract Methods

func (t *TransactionContract) InitiateTransaction(ctx contractapi.TransactionContextInterface, transactionJson []byte) error {

	var transactionInitiate TransactionInitiate

	err := json.Unmarshal(transactionJson, &transactionInitiate)
	if err != nil {
		return fmt.Errorf("failed to unmarshal Json: %v", err)
	}
	transactionId := string(transactionInitiate.Id)
	fmt.Printf("Deserialized data: %+v\n", transactionInitiate)
	fmt.Println(transactionId)

	exists, err := t.TrxExists(ctx, transactionId)
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
	fmt.Println(trxJson)

	return ctx.GetStub().PutState(transactionId, trxJson)

}

func (t *TransactionContract) ReadAllTransactions(ctx contractapi.TransactionContextInterface) ([]*TransactionInitiate, error) {
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
		fmt.Println(queryResponse.Value)
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
	var transaction *TransactionInitiate
	err = json.Unmarshal(transactionJson, transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (t *TransactionContract) UpdateTransactionProcess(ctx contractapi.TransactionContextInterface, updatedTransactionProcessJson []byte) error {

	var updateTransactionProcess UpdateTransactionProcess

	err := json.Unmarshal((updatedTransactionProcessJson), &updateTransactionProcess)

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
				newTrxProcess := trx
				newTrxProcess.Remarks = trxProcess.Remarks
				newTrxProcess.StatusModifiedDate = trxProcess.StatusModifiedDate
				newTrxProcess.TransactionStatus = trx.TransactionStatus
				transactionData.TransactionProcesses = append(transactionData.TransactionProcesses[:i], newTrxProcess)
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

func (t *TransactionContract) AddRecipientPayment(ctx contractapi.TransactionContextInterface, recipientJson []byte) error {

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

	transaction.TransactionPayments = append(transaction.TransactionPayments, recipientGroupPayment.TransactionPayment)
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

func (t *TransactionContract) AddMemberPayment(ctx contractapi.TransactionContextInterface, memberPaymentJson []byte) error {
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
