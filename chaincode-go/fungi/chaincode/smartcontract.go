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
type Feed struct {
	FeedId  uint   `json:"fungusid"`
	Name      string `json:"name"`
	Dna       uint   `json:"dna"`
}

// Define key names for options
const feedCountKey = "feedCount" // the ​total number of fungi

// Define const value for basic setting of contract
const dnaDigits uint = 14

func (s *SmartContract) Initialize(ctx contractapi.TransactionContextInterface) (bool, error) {

	// Check minter authorization - this sample assumes Org1 is the central banker with privilege to intitialize contract
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed to get MSPID: %v", err)
	}
	// only Org1MSG members can call
	if clientMSPID != "Org2MSP" {
		return false, fmt.Errorf("client is not authorized to initialize contract")
	}

	// Check contract options are not already set, client is not authorized to change them once intitialized
	fungusCount, err := ctx.GetStub().GetState(feedCountKey)
	
	if err != nil {
		return false, fmt.Errorf("failed to get Name: %v", err)
	}
	if fungusCount != nil {
		return false, fmt.Errorf("contract options are already set, client is not authorized to change them")
	}

	// Initialize FungusCountKey to zero(0)
	err = ctx.GetStub().PutState(feedCountKey, []byte(strconv.Itoa(0)))
	if err != nil {
		return false, fmt.Errorf("failed to set token name: %v", err)
	}

	return true, nil
}
// TODO: mod feedFactory function
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
		ReadyTime: uint32(unixTime),
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
	err = s._updateOwnerFungusCount(ctx, clientID, -1)
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