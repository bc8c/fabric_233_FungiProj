package chaincode

import (
	"strconv"
	"fmt"
	"encoding/json"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)


func (s *SmartContract)Feed(ctx contractapi.TransactionContextInterface, fungusId uint, feedId uint) error{
	params := []string{"GetFeed", strconv.Itoa(int(feedId))}
	invokeargs := make([][]byte, len(params))
	for i, arg := range params {
		invokeargs[i] = []byte(arg)
	}

	// Get FeedDna
	respons := ctx.GetStub().InvokeChaincode("feed", invokeargs, "mychannel")
	if respons.Status != 200 {
		return fmt.Errorf("failed to InvokeChaincode: %s", respons.Payload)
	}
	var feed struct {		
		Dna    uint   `json:"dna"`
	}
	err := json.Unmarshal(respons.Payload, &feed)
	if err != nil {
		return err
	}

	// Get fungiDna
	fungusJSON, err := s._getState(ctx, strconv.Itoa(int(fungusId)))
	if err != nil {
		return fmt.Errorf("failed to get fungus: %v", err)
	}

	var fungus Fungus
	err = json.Unmarshal(fungusJSON, &fungus)
	if err != nil {
		return err
	}
	
	// check readytime
	nowTime := time.Now()
	unixTime := nowTime.Unix()
	if uint32(unixTime) < fungus.ReadyTime {
		return fmt.Errorf("failed to get Multiply : Not Ready")
	}
	
	// call multiply func
	err = s._feedAndMultiply(ctx, fungus.Dna, feed.Dna)
	if err != nil {
		return fmt.Errorf("failed to get Multiply : %v", err)
	}

	return nil
}

func (s *SmartContract)_feedAndMultiply(ctx contractapi.TransactionContextInterface, fungusId uint, feedDna uint) error {
	var newDna uint = (fungusId + feedDna) /2 
	newDna = newDna - (newDna % 100) + 1
	
	err := s._createFungus(ctx, "noname", newDna)
	if err != nil {
		return fmt.Errorf("failed to createFungus: %v", err)
	}

	return nil
}
