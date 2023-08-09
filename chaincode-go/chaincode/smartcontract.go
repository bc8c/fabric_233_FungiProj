package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}
// Fungus Asset describes basic details
type Fungus struct{
	Name string
	Dna uint
	ReadyTime uint32
}

//fungusToOwner (key:value) -> WSDB 
//ownerFungusCount (key:value) -> WSDB 

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) createFungus(ctx contractapi.TransactionContextInterface, name string, dna uint) error {
	// PutState Fungus in WSDB

	return ctx.GetStub().PutState(id, assetJSON)
}