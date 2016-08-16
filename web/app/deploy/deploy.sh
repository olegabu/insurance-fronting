#!/usr/bin/env bash

curl -XPOST -d @deployRemote.json http://vp0:5000/chaincode
curl -XPOST -d @demoInit.json http://vp0:5000/chaincode
