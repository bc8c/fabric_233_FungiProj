package chaincode

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// returns the asset stored in the world state with given id.
func (s *SmartContract) _getState(ctx contractapi.TransactionContextInterface, id string) ([]byte, error) {	
	bytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if bytes == nil {
		return nil, fmt.Errorf("the asset %s does not exist", fungusCountKey)
	}
	return bytes, nil
}