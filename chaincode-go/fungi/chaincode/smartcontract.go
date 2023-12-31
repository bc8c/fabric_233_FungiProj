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
const fungusCountKey = "fungusCount" // the ​total number of fungi

// Define const value for basic setting of contract
const dnaDigits uint = 14

//fungusToOwner (key:value) -> WSDB
//ownerFungusCount (key:value) -> WSDB
//fungusCount (key:value) -> WSDB

func (s *SmartContract) Initialize(ctx contractapi.TransactionContextInterface) (bool, error) {

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
	fungusCount, err := ctx.GetStub().GetState(fungusCountKey)
	
	if err != nil {
		return false, fmt.Errorf("failed to get Name: %v", err)
	}
	if fungusCount != nil {
		return false, fmt.Errorf("contract options are already set, client is not authorized to change them")
	}

	// Initialize FungusCountKey to zero(0)
	err = ctx.GetStub().PutState(fungusCountKey, []byte(strconv.Itoa(0)))
	if err != nil {
		return false, fmt.Errorf("failed to set token name: %v", err)
	}

	return true, nil
}

// create a new fungus API
func (s *SmartContract) CreateRandomFungus(ctx contractapi.TransactionContextInterface, name string) error{
	// Check ClientId
	clientId, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get clientId: %v", err)
	}

	exists, err := s._assetExists(ctx, clientId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("client has already created an initial fungus")
	}
	ctx.GetStub().PutState(clientId, []byte(strconv.Itoa(0)))
	if err != nil {
		return fmt.Errorf("failed to put fungus state: %v", err)
	}

	dna:=s._generateRandomDna(name)
	err = s._createFungus(ctx, name, dna)
	if err != nil {
		return fmt.Errorf("failed to createFungus: %v", err)
	}
	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) _createFungus(ctx contractapi.TransactionContextInterface, name string, dna uint) error {

	// Check ClientId
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get clientID: %v", err)
	}

	// for readytime
	nowTime := time.Now()
	unixTime := nowTime.Unix()
	
	//  make fungusid
	fungusCountBytes, err := s._getState(ctx, fungusCountKey)	
	if err != nil {
		return fmt.Errorf("failed to get fungusCount: %v", err)
	}
	fungusIdINT,_ := strconv.Atoi(string(fungusCountBytes))
	fungusId := uint(fungusIdINT)

	// overwriting original fungus with new fungus
	fungus := Fungus{
		FungusId:  fungusId,
		Name:      name,
		Owner:     clientID,
		Dna:       dna,
		ReadyTime: uint32(unixTime)+60,
	}
	assetJSON, err := json.Marshal(fungus)
	if err != nil {
		return fmt.Errorf("failed to marshal fungus: %v", err)
	}
	
	// PutState fungusId
	err = ctx.GetStub().PutState(strconv.Itoa(int(fungusId)), assetJSON)
	if err != nil {
		return fmt.Errorf("failed to put fungus state: %v", err)
	}

	// PutState fungusCount++
	fungusId += 1
	err = ctx.GetStub().PutState(fungusCountKey, []byte(strconv.Itoa(int(fungusId))))
	if err != nil {
		return fmt.Errorf("failed to put fungusCount state: %v", err)
	}

	//  update ownerFungusCount
	err = s._updateOwnerFungusCount(ctx, clientID, 1)
	if err != nil {
		return fmt.Errorf("failed to put fungus state: %v", err)
	}

	return nil
}

// generate random dna func
func (S *SmartContract) _generateRandomDna(name string) uint {
	nowTime := time.Now()
	unixTime := nowTime.Unix()
	data := strconv.Itoa(int(unixTime)) + name
	hash := sha256.New()
	hash.Write([]byte(data))
	dnaHash := uint(binary.BigEndian.Uint64(hash.Sum(nil)))

	// make 14digits dna
	dna := dnaHash % uint(math.Pow(10, float64(dnaDigits)))
	dna = dna - dna%100

	return dna
}

func (S *SmartContract) GetFungiByOwner(ctx contractapi.TransactionContextInterface) ([]*Fungus, error){
	// Check ClientId
	clientId, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return nil, fmt.Errorf("failed to get clientID: %v", err)
	}
	queryString := fmt.Sprintf(`{"selector":{"owner":"%s"}}`, clientId)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var fungi []*Fungus
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var fungus Fungus
		err = json.Unmarshal(queryResult.Value, &fungus)
		if err != nil {
			return nil, err
		}
		fungi = append(fungi, &fungus)
	}
	return fungi, nil
}

func (S *SmartContract) BalanceOf(ctx contractapi.TransactionContextInterface, owner string) (int, error){

	results, err := ctx.GetStub().GetState(owner)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %v", err)
	}

	balance, _ := strconv.Atoi(string(results))

	return balance, nil
}

func (S *SmartContract) OwnerOf(ctx contractapi.TransactionContextInterface, fungusid uint) (string, error){
	
	fungusBytes, err := S._getState(ctx, strconv.Itoa(int(fungusid)))
	if err != nil {
		return "", fmt.Errorf("failed to get fungus: %v", err)
	}
	var fungus Fungus
	err = json.Unmarshal(fungusBytes, &fungus)
	if err != nil {
		return "", fmt.Errorf("failed to Unmarshal fungus: %v", err)
	}
	return fungus.Owner, nil
}

func (S *SmartContract) TransferFrom(ctx contractapi.TransactionContextInterface, from string, to string, fungusid uint) (bool, error) {

	fungusBytes, err := S._getState(ctx, strconv.Itoa(int(fungusid)))
	if err != nil {
		return false, fmt.Errorf("failed to get fungus: %v", err)
	}
	if fungusBytes == nil {
		return false, fmt.Errorf("not exists fungus: %v", err)
	}
	
	var fungus Fungus
	err = json.Unmarshal(fungusBytes, &fungus)

	if fungus.Owner != from {
		return false, fmt.Errorf("not a fungus owner's request.: %v", err)
	}

	fungus.Owner = to

	assetJSON, err := json.Marshal(fungus)
	if err != nil {
		return false, fmt.Errorf("failed to marshal fungus: %v", err)
	}
	
	// PutState fungusId
	err = ctx.GetStub().PutState(strconv.Itoa(int(fungusid)), assetJSON)
	if err != nil {
		return false, fmt.Errorf("failed to put fungus state: %v", err)
	}

	//  update ownerFungusCount for from
	err = S._updateOwnerFungusCount(ctx, from, -1)
	if err != nil {
		return false, fmt.Errorf("failed to put fungus state: %v", err)
	}
	//  update ownerFungusCount for to
	err = S._updateOwnerFungusCount(ctx, to, 1)
	if err != nil {
		return false, fmt.Errorf("failed to put fungus state: %v", err)
	}

	return true, nil
}