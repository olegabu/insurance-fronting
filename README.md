# Insurance Fronting Proof Of Concept
Proof of concept demo illustrating the practice of international insurance fronting on Hyperledger blockchain.

# Web Application
Demo is served by an Angular single page web application. Please install and run in `web` directory.

## Install
```
npm install
bower install
```
Will download developer dependencies which may take a little while.

## Run
The web app is built and run by `gulp`:

```
gulp serve
```

# Changes to core.yaml 

    security.enabled = true

# Changes to membersrvc.yaml 
Add 6 more users

                citi: 1 4nXSrfoYGFCP bank_a
                auditor: 1 yg5DVhm0er1z institution_a
                bermuda: 1 b7pmSxzKNFiw institution_a
                art: 1 YsWZD4qQmYxo institution_a
                allianz: 1 W8G0usrU7jRk institution_a
                nigeria: 1 H80SiB5ODKKQ institution_a

And add correct atributes to the aca.attributes section

              attribute-entry-17: citi;bank_a;role;bank;2016-01-01T00:00:00-03:00;;
              attribute-entry-18: auditor;institution_a;role;auditor;2016-01-01T00:00:00-03:00;;
              attribute-entry-19: bermuda;institution_a;role;captive;2015-02-02T00:00:00-03:00;;
              attribute-entry-20: art;institution_a;role;reinsurer;2015-02-02T00:00:00-03:00;;
              attribute-entry-21: allianz;institution_a;role;fronter;2015-01-01T00:00:00-03:00;;
              attribute-entry-22: nigeria;institution_a;role;affiliate;2015-01-01T00:00:00-03:00;;
              
              attribute-entry-23: citi;bank_a;company;Citi;2016-01-01T00:00:00-03:00;;
              attribute-entry-24: auditor;institution_a;company;Auditor;2016-01-01T00:00:00-03:00;;
              attribute-entry-25: bermuda;institution_a;company;Bermuda;2015-02-02T00:00:00-03:00;;
              attribute-entry-26: art;institution_a;company;Art;2015-02-02T00:00:00-03:00;;
              attribute-entry-27: allianz;institution_a;company;Allianz;2015-01-01T00:00:00-03:00;;
              attribute-entry-28: nigeria;institution_a;company;Nigeria;2015-01-01T00:00:00-03:00;;
            
# Deploy chaincode
              
1) Deploy chain
              
    curl -XPOST -d  '{"jsonrpc": "2.0", "method": "deploy",  "params": {"type": 1,"chaincodeID": {"path": "github.com/olegabu/insurance-fronting/chaincode","language": "GOLANG"}, "ctorMsg": { "args": ["aW5pdA=="] },"secureContext": "citi", "attributes": ["role", "company"]},"id": 0}' http://vp0:7050/chaincode

2) copy new HASH into config.js

3) and setup initial configuration

    curl -XPOST -d  '{"jsonrpc": "2.0", "method": "invoke", "params": {"type": 1, "chaincodeID": {"name": "'"$HASH"'"}, "ctorMsg": {"args": ["ZGVtb0luaXQ="]}, "secureContext": "citi", "attributes": ["role", "company"]}, "id": 1}' http://vp0:7050/chaincode
    