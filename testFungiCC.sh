#!/bin/bash

C_YELLOW='\033[1;33m'
C_BLUE='\033[0;34m'
C_RESET='\033[0m'

# subinfoln echos in blue color
function infoln() {
  echo -e "${C_YELLOW}${1}${C_RESET}"
}

function subinfoln() {
  echo -e "${C_BLUE}${1}${C_RESET}"
}

# add PATH to ensure we are picking up the correct binaries
export PATH=${HOME}/fabric-samples/bin:$PATH
export FABRIC_CFG_PATH=${PWD}/network/config

# Test Setting
CC_NAME="fungi"
CHANNEL_NAME="mychannel"
PEER_CONN_PARMS="--peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt --peerAddresses localhost:11051 --tlsRootCertFiles ${PWD}/network/organizations/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt"
ORDERER_CA=${PWD}/network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

## Set the peer on peer0.org1
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051


## TEST1 : Invoking the chaincode
infoln "TEST1-1 : Invoking the chaincode (Initialize)"
set -x
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n ${CC_NAME} $PEER_CONN_PARMS -c '{"function":"Initialize","Args":[]}' >&log.txt
{ set +x; } 2>/dev/null
cat log.txt
sleep 3

# ## TEST2 : Query the chaincode
# infoln "TEST2 : Query the chaincode"
# set -x
# peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"Args":["GetAllAssets"]}' >&log.txt
# { set +x; } 2>/dev/null
# cat log.txt

## TEST2 : Invoking the chaincode
infoln "TEST2-1 : Invoking the chaincode (CreateRandomFungus)"
set -x
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n ${CC_NAME} $PEER_CONN_PARMS -c '{"function":"CreateRandomFungus","Args":["TestFungus1"]}' >&log.txt
{ set +x; } 2>/dev/null
cat log.txt
sleep 3

## TEST3 : Query the chaincode
infoln "TEST2 : Query the chaincode (GetFungiByOwner)"
set -x
peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"Args":["GetFungiByOwner"]}' >&log.txt
{ set +x; } 2>/dev/null
cat log.txt

