#!/usr/bin/env bash

curl -XPOST -d @deployRemote.json http://vp0:7050/chaincode
curl -XPOST -d @demoInit.json http://vp0:7050/chaincode


curl -XPOST -d  '{"jsonrpc": "2.0", "method": "deploy",  "params": {"type": 1,"chaincodeID": {"path": "github.com/olegabu/insurance-fronting/chaincode","language": "GOLANG"}, "ctorMsg": { "args": ["aW5pdA=="] },"secureContext": "citi", "attributes": ["role", "company"]},"id": 0}' http://vp0:7050/chaincode
curl -XPOST -d  '{"jsonrpc": "2.0", "method": "invoke", "params": {"type": 1, "chaincodeID": {"name": "'"$HASH"'"}, "ctorMsg": {"args": ["ZGVtb0luaXQ="]}, "secureContext": "citi", "attributes": ["role", "company"]}, "id": 1}' http://vp0:7050/chaincode
