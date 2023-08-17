package chaincode

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)


func (s *SmartContract)Feed(ctx contractapi.TransactionContextInterface) error{
	params := []string{"GetFeed", "0"}
	invokeargs := make([][]byte, len(params))
	for i, arg := range params {
		invokeargs[i] = []byte(arg)
	}

	respons := ctx.GetStub().InvokeChaincode("feed", invokeargs, "mychannel")
	if respons.Status != 200 {
		return fmt.Errorf("failed to InvokeChaincode: %s", respons.Payload)
	}
	return nil
}