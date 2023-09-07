#!/bin/bash

pushd application

rm -rf wallet/*

npm install

./public/ccp/ccp-generate.sh

sleep 2

npm start

popd