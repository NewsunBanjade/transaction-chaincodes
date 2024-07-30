package main

import (
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	chaincode2 "github.com/newsunbanjade/transaction-chaincode/chaincode"
	"log"
)

func main() {
	chaincode, err := contractapi.NewChaincode(new(chaincode2.SmartContract))
	if err != nil {
		log.Panicf("Error creating transaction chaincode: %s", err.Error())
	}
	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting transaction chaincode: %s", err.Error())
	}
}
