package chaincode

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// returns the asset stored in the world state with given id.
func (s *SmartContract) _getState(ctx contractapi.TransactionContextInterface, id string) ([]byte, error) {	
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", fungusCountKey)
	}
	return assetJSON, nil
}

func (s *SmartContract) _assetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

func (s *SmartContract) _updateOwnerFungusCount(ctx contractapi.TransactionContextInterface, clientID string, increment int) error {
	countByte, err := s._getState(ctx, clientID)	
	if err != nil {
		return fmt.Errorf("failed to get fungusCount: %v", err)
	}
	ownerFungusCount,_ := strconv.Atoi(string(countByte[:]))
	ownerFungusCount += increment
	ctx.GetStub().PutState(clientID, []byte(strconv.Itoa(ownerFungusCount)))
	if err != nil {
		return fmt.Errorf("failed to put fungus state: %v", err)
	}
	return nil
}