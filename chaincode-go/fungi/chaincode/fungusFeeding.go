package chaincode

import (
	"strconv"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)


func (s *SmartContract)Feed(ctx contractapi.TransactionContextInterface, feedId uint) (string,error){
	params := []string{"GetFeed", strconv.Itoa(int(feedId))}
	invokeargs := make([][]byte, len(params))
	for i, arg := range params {
		invokeargs[i] = []byte(arg)
	}

	respons := ctx.GetStub().InvokeChaincode("feed", invokeargs, "mychannel")
	if respons.Status != 200 {
		return "",fmt.Errorf("failed to InvokeChaincode: %s", respons.Payload)
	}
	return string(respons.Payload), nil
}