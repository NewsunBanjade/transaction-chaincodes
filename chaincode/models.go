package chaincode

import "time"

// *Basic Data Stucts
// Transaction Models
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
