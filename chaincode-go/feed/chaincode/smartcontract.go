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

// Define key names for options
const feedsCountKey = "feedsCount"

// Feed Asset describes basic details
type Feed struct {
	FeedId uint   `json:"fungusid"`
	Name   string `json:"name"`
	Dna    uint   `json:"dna"`
}

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
	fungusCount, err := ctx.GetStub().GetState(feedsCountKey)

	if err != nil {
		return false, fmt.Errorf("failed to get Name: %v", err)
	}
	if fungusCount != nil {
		return false, fmt.Errorf("contract options are already set, client is not authorized to change them")
	}

	// Initialize fungusCount to zero(0)
	err = ctx.GetStub().PutState(feedsCountKey, []byte(strconv.Itoa(0)))
	if err != nil {
		return false, fmt.Errorf("failed to set token name: %v", err)
	}

	return true, nil
}

// create a new fungus API
func (s *SmartContract) CreateRandomFeed(ctx contractapi.TransactionContextInterface, name string) error {

	dna := s._generateRandomDna(name)
	err := s._createFeeds(ctx, name, dna)
	if err != nil {
		return fmt.Errorf("failed to createFeeds: %v", err)
	}
	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) _createFeeds(ctx contractapi.TransactionContextInterface, name string, dna uint) error {

	//  make feedId
	feedsCountBytes, err := s._getState(ctx, feedsCountKey)
	if err != nil {
		return fmt.Errorf("failed to get fungusCount: %v", err)
	}
	feedsIdINT, _ := strconv.Atoi(string(feedsCountBytes))
	feedsId := uint(feedsIdINT)

	// overwriting original fungus with new fungus
	feeds := Feed{
		FeedId: feedsId,
		Name:   name,
		Dna:    dna,
	}
	assetJSON, err := json.Marshal(feeds)
	if err != nil {
		return fmt.Errorf("failed to marshal fungus: %v", err)
	}

	// PutState fungusId
	err = ctx.GetStub().PutState(strconv.Itoa(int(feedsId)), assetJSON)
	if err != nil {
		return fmt.Errorf("failed to put fungus state: %v", err)
	}

	// PutState fungusCount++
	feedsId += 1
	err = ctx.GetStub().PutState(feedsCountKey, []byte(strconv.Itoa(int(feedsId))))
	if err != nil {
		return fmt.Errorf("failed to put fungusCount state: %v", err)
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

func (S *SmartContract) GetFeed(ctx contractapi.TransactionContextInterface, feedId uint) (*Feed, error) {
	feedsCountBytes, err := S._getState(ctx, strconv.Itoa(int(feedId)))
	if err != nil {
		return nil, fmt.Errorf("failed to get fungusCount: %v", err)
	}

	var feed Feed
	err = json.Unmarshal(feedsCountBytes, &feed)
	if err != nil {
		return nil, fmt.Errorf("failed to Unmarshal feed: %v", err)
	}

	return &feed, nil
}