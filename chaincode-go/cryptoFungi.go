/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"cryptoFungi/chaincode"
)

func main() {
	// assetChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	// if err != nil {
	// 	log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
	// }

	// if err := assetChaincode.Start(); err != nil {
	// 	log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
	// }

	TestSC := chaincode.SmartContract{}
	TestSC.Testfunc(1, "TestFungus")

}
