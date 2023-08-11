package chaincode

import (
	"encoding/json"
	"time"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}
// Fungus Asset describes basic details
type Fungus struct{
	Name string			`json:"name"`
	Dna uint			`json:"dna"`
	ReadyTime uint32	`json:"readytime`
}

//fungusToOwner (key:value) -> WSDB 
//ownerFungusCount (key:value) -> WSDB 
//fungusCount (key:value) -> WSDB 

func (s *SmartContract) Initialize(ctx contractapi.TransactionContextInterface, name string, symbol string, decimals string) (bool, error) {

	// Check minter authorization - this sample assumes Org1 is the central banker with privilege to intitialize contract
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed to get MSPID: %v", err)
	}
	if clientMSPID != "Org1MSP" {
		return false, fmt.Errorf("client is not authorized to initialize contract")
	}

	// Check contract options are not already set, client is not authorized to change them once intitialized
	bytes, err := ctx.GetStub().GetState("fungusCount")
	if err != nil {
		return false, fmt.Errorf("failed to get Name: %v", err)
	}
	if bytes != nil {
		return false, fmt.Errorf("contract options are already set, client is not authorized to change them")
	}

	err = ctx.GetStub().PutState("fungusCount", []byte(strconv.Itoa(0)))
	if err != nil {
		return false, fmt.Errorf("failed to set token name: %v", err)
	}

	return true, nil
}


// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) createFungus(ctx contractapi.TransactionContextInterface, name string, dna uint) error {
	// PutState Fungus in WSDB

	// readytime
	nowTime := time.Now()
	unixTime := nowTime.Unix()

	

	// overwriting original asset with new asset
	fungus := Fungus{
		Name:		name,
		Dna:		dna,
		ReadyTime:	uint32(unixTime),
	}
	assetJSON, err := json.Marshal(fungus)
	if err != nil {
		return err
	}

	// create randDNA

	return ctx.GetStub().PutState("id", assetJSON)
}