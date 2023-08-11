package chaincode

import (
	"encoding/json"
	"encoding/binary"
	"time"
	"fmt"
	"strconv"
	"crypto/sha256"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}
// Fungus Asset describes basic details
type Fungus struct{
	FungusId uint		`json:"fungusid"`
	Name string			`json:"name"`
	Owner string		`json:"owner"`
	Dna uint			`json:"dna"`
	ReadyTime uint32	`json:"readytime"`
}

// Define key names for options
const fungusCountKey = "fungusCount"


//fungusToOwner (key:value) -> WSDB 
//ownerFungusCount (key:value) -> WSDB 
//fungusCount (key:value) -> WSDB 

func (s *SmartContract) Initialize(ctx contractapi.TransactionContextInterface, name string, symbol string, decimals string) (bool, error) {

	// Check minter authorization - this sample assumes Org1 is the central banker with privilege to intitialize contract
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed to get MSPID: %v", err)
	}
	// only Org1MSG members can call
	if clientMSPID != "Org1MSP" {
		return false, fmt.Errorf("client is not authorized to initialize contract")
	}

	// Check contract options are not already set, client is not authorized to change them once intitialized
	bytes, err := ctx.GetStub().GetState(fungusCountKey)
	if err != nil {
		return false, fmt.Errorf("failed to get Name: %v", err)
	}
	if bytes != nil {
		return false, fmt.Errorf("contract options are already set, client is not authorized to change them")
	}

	// Initialize FungusCountKey to zero(0)
	err = ctx.GetStub().PutState(fungusCountKey, []byte(strconv.Itoa(0)))
	if err != nil {
		return false, fmt.Errorf("failed to set token name: %v", err)
	}

	return true, nil
}


// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) createFungus(ctx contractapi.TransactionContextInterface, fungusid uint, name string, dna uint) error {
	// PutState Fungus in WSDB

	// Check minter authorization - this sample assumes Org1 is the central banker with privilege to intitialize contract
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get clientID: %v", err)
	}

	// readytime
	nowTime := time.Now()
	unixTime := nowTime.Unix()	
	
	// create randDNA
	data := uint(unixTime)+fungusid
	hash := sha256.New()
	hash.Write([]byte(strconv.Itoa(int(data))))
	hashDna := uint(binary.BigEndian.Uint64(hash.Sum(nil)))
	// @need ( digit 14...... )

	
	// hashDna := uint(hash.Sum(nil))

	// overwriting original fungus with new fungus
	fungus := Fungus{
		FungusId: 	fungusid,
		Name:		name,
		Owner:		clientID,
		Dna:		hashDna,
		ReadyTime:	uint32(unixTime),
	}
	assetJSON, err := json.Marshal(fungus)
	if err != nil {
		return err
	}

	

	return ctx.GetStub().PutState(fungusCountKey+name, assetJSON)
	
	
}