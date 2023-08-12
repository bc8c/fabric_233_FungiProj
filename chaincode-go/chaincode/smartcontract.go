package chaincode

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Fungus Asset describes basic details
type Fungus struct {
	FungusId  uint   `json:"fungusid"`
	Name      string `json:"name"`
	Owner     string `json:"owner"`
	Dna       uint   `json:"dna"`
	ReadyTime uint32 `json:"readytime"`
}

// Define key names for options
const fungusCountKey = "fungusCount"

// Define const value for basic setting of contract
const dnaDigits uint = 14

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
func (s *SmartContract) createFungus(ctx contractapi.TransactionContextInterface, name string) error {
	// PutState Fungus in WSDB

	// Check minter authorization - this sample assumes Org1 is the central banker with privilege to intitialize contract
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get clientID: %v", err)
	}

	// readytime
	nowTime := time.Now()
	unixTime := nowTime.Unix()
	
	// TODO: make a common getState func
	//  How to make fungusid
	fungusCountBytes, err := ctx.GetStub().GetState(fungusCountKey)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if fungusCountBytes == nil {
		return fmt.Errorf("the asset %s does not exist", fungusCountKey)
	}
	fungusidINT,_ := strconv.Atoi(string(fungusCountBytes))
	fungusid := uint(fungusidINT)

	// create randDNA
	data := uint(unixTime) + fungusid
	hash := sha256.New()
	hash.Write([]byte(strconv.Itoa(int(data))))
	dnaHash := uint(binary.BigEndian.Uint64(hash.Sum(nil)))

	// make 14digits dna
	dna := dnaHash % uint(math.Pow(10, float64(dnaDigits)))
	dna = dna - dna%100

	// overwriting original fungus with new fungus
	fungus := Fungus{
		FungusId:  fungusid,
		Name:      name,
		Owner:     clientID,
		Dna:       dna,
		ReadyTime: uint32(unixTime),
	}
	assetJSON, err := json.Marshal(fungus)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(strconv.Itoa(int(fungusid)), assetJSON)

}

func (s *SmartContract) Testfunc(fungusid uint, name string) error {

	return nil

}
