package main

import (
	"log"

	"asset"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	cc, err := contractapi.NewChaincode(&asset.SmartContract{})
	if err != nil {
		log.Panicf("Error creating CHoT medical data for quantum-blockchain chaincode: %v", err)
	}

	if err := cc.Start(); err != nil {
		log.Panicf("Error starting CHoT medical data for quantum-blockchain chaincode: %v", err)
	}
}
