package main

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"tutdeputraw.com/chaincode"
)

func main() {
	err := shim.Start(new(chaincode.RealEstateChaincode))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
