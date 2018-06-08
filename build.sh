#!/bin/sh

# pull in dependencies
go get -u -v github.com/ethereum/go-ethereum

# generate the Go ABI bindings from the Solidity source code
#abigen -sol contracts/ERC20Interface.sol -pkg main -out ./erc20.go

# build the application
go build -o toy-block-explorer
