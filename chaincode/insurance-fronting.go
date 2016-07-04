package main

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("insurance-fronting")

// InsuranceFrontingChaincode implements standard Chaincode interface
type InsuranceFrontingChaincode struct {
}

type contract struct {
	ID                   string
	MaxCoverage          uint64
	MaxPremium           uint64
	Captive              string
	Reinsurer            string
	CurrentTotalCoverage uint64
	CurrentTotalPremium  uint64
	CurrentPaidClaim     uint64
	CurrentPaidPremium   uint64
}

// frontingChain defines mapping between roles and orgs in insurance
type frontingChain struct {
	Captive   string
	Reinsurer string
	Fronter   string
	Affiliate string
}

// claim ..
type claim struct {
	PolicyID            string
	ClaimID             string
	Amt                 uint64
	ApprovedByCaptive   bool
	ApprovedByReinsurer bool
	ApprovedByFronter   bool
}

type policy struct {
	ContractID    string
	PolicyID      string
	Coverage      uint64
	Premium       uint64
	PaidClaim     uint64
	PaidPremium   uint64
	FrontingChain frontingChain
}

type transaction struct {
	ID      string
	From    string
	To      string
	Amt     uint64
	Purpose string
}

// Init creates Policy and Transaction tables in Ledger
func (t *InsuranceFrontingChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	log.Debugf("function: %s, args: %s", function, args)

	// Create contracts table
	err := stub.CreateTable("Contracts", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "ID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "MaxCoverage", Type: shim.ColumnDefinition_UINT64, Key: false},
		&shim.ColumnDefinition{Name: "MaxPremium", Type: shim.ColumnDefinition_UINT64, Key: false},
		&shim.ColumnDefinition{Name: "Captive", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Reinsurer", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "CurrentTotalCoverage", Type: shim.ColumnDefinition_UINT64, Key: false},
		&shim.ColumnDefinition{Name: "CurrentTotalPremium", Type: shim.ColumnDefinition_UINT64, Key: false},
		&shim.ColumnDefinition{Name: "CurrentPaidClaim", Type: shim.ColumnDefinition_UINT64, Key: false},
		&shim.ColumnDefinition{Name: "CurrentPaidPremium", Type: shim.ColumnDefinition_UINT64, Key: false},
	})
	if err != nil {
		log.Criticalf("function: %s, args: %s", function, args)
		return nil, errors.New("Failed creating Contracts table.")
	}

	// Create policies table
	err = stub.CreateTable("Policies", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "ContractID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "PolicyID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Coverage", Type: shim.ColumnDefinition_UINT64, Key: false},
		&shim.ColumnDefinition{Name: "Premium", Type: shim.ColumnDefinition_UINT64, Key: false},
		&shim.ColumnDefinition{Name: "PaidClaim", Type: shim.ColumnDefinition_UINT64, Key: false},
		&shim.ColumnDefinition{Name: "PaidPremium", Type: shim.ColumnDefinition_UINT64, Key: false},
		&shim.ColumnDefinition{Name: "FrontingChain", Type: shim.ColumnDefinition_BYTES, Key: false},
	})
	if err != nil {
		log.Criticalf("function: %s, args: %s", function, args)
		return nil, errors.New("Failed creating Policies table.")
	}

	// Create claims table
	err = stub.CreateTable("Claims", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "PolicyID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ClaimID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Amt", Type: shim.ColumnDefinition_UINT64, Key: false},
		&shim.ColumnDefinition{Name: "ApprovedByCaptive", Type: shim.ColumnDefinition_BOOL, Key: false},
		&shim.ColumnDefinition{Name: "ApprovedByReinsurer", Type: shim.ColumnDefinition_BOOL, Key: false},
		&shim.ColumnDefinition{Name: "ApprovedByFronter", Type: shim.ColumnDefinition_BOOL, Key: false},
	})
	if err != nil {
		log.Criticalf("function: %s, args: %s", function, args)
		return nil, errors.New("Failed creating Claims table.")
	}

	// Create transactions table
	err = stub.CreateTable("Transactions", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "ID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "From", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "To", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Amt", Type: shim.ColumnDefinition_UINT64, Key: false},
		&shim.ColumnDefinition{Name: "Purpose", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		log.Criticalf("function: %s, args: %s", function, args)
		return nil, errors.New("Failed creating Transactions table.")
	}

	err = stub.PutState("ContractsCounter", []byte(strconv.FormatUint(0, 10)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState("PoliciesCounter", []byte(strconv.FormatUint(0, 10)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState("TransactionsCounter", []byte(strconv.FormatUint(0, 10)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState("ClaimsCounter", []byte(strconv.FormatUint(0, 10)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke executes function on Ledger
func (t *InsuranceFrontingChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	log.Debugf("function: %s, args: %s", function, args)

	// Handle different functions
	if function == "createContract" {
		if len(args) != 4 {
			return nil, errors.New("Incorrect arguments. Expecting Captive, Reinsurer, coverage and premium.")
		}

		coverage, err := strconv.ParseUint(args[2], 10, 64)
		if err != nil {
			return nil, errors.New("Incorrect coverage. Uint64 expected.")
		}

		premium, err := strconv.ParseUint(args[3], 10, 64)
		if err != nil {
			return nil, errors.New("Incorrect premium. Uint64 expected.")
		}

		newContract := contract{Captive: args[0], Reinsurer: args[1], MaxCoverage: coverage, MaxPremium: premium}

		return t.createContract(stub, newContract)

	} else if function == "createPolicy" {
		if len(args) != 3 {
			return nil, errors.New("Incorrect arguments. Expecting contract ID, coverage and premium.")
		}

		coverage, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return nil, errors.New("Incorrect coverage. Uint64 expected.")
		}

		premium, err := strconv.ParseUint(args[2], 10, 64)
		if err != nil {
			return nil, errors.New("Incorrect premium. Uint64 expected.")
		}

		_, err = t.createPolicy(stub, args[0], coverage, premium)
		return nil, err

	} else if function == "join" {
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting policyID.")
		}

		return t.joinChain(stub, args[0])

	} else if function == "pay" {
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting policyID.")
		}
		return t.payPremium(stub, args[0])

	} else if function == "claim" {
		if len(args) != 2 {
			return nil, errors.New("Incorrect number of arguments. Expecting policyID and amount.")
		}
		n, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return nil, errors.New("Incorrect argument. Uint64 expected.")
		}
		return t.createClaim(stub, args[0], n)

	} else if function == "approve" {
		if len(args) != 2 {
			return nil, errors.New("Incorrect number of arguments. Expecting policyID and claimID.")
		}
		return t.approveClaim(stub, args[0], args[1])

	} else if function == "emitCoins" {
		if len(args) != 2 {
			return nil, errors.New("Incorrect number of arguments. Expecting company name and amount.")
		}
		n, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return nil, errors.New("Incorrect argument. Uint64 expected.")
		}
		return t.emitCoins(stub, args[0], n)
	} else if function == "demoInit" {
		return t.demoInit(stub)

	} else {
		log.Errorf("function: %s, args: %s", function, args)
		return nil, errors.New("Received unknown function invocation")
	}
}

// Query callback representing the query of a chaincode
func (t *InsuranceFrontingChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	log.Debugf("function: %s, args: %s", function, args)

	// Handle different functions
	if function == "getPolicies" {
		var contractID string
		if len(args) == 1 {
			contractID = args[0]
		}

		policies, err := t.getPolicies(stub, contractID)
		if err != nil {
			return nil, err
		}

		return json.Marshal(policies)

	} else if function == "getPolicy" {
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting policyID.")
		}

		policy, err := t.getPolicy(stub, args[0])
		if err != nil {
			return nil, err
		}

		return json.Marshal(policy)

	} else if function == "getContracts" {
		contracts, err := t.getContracts(stub)
		if err != nil {
			return nil, err
		}

		return json.Marshal(contracts)

	} else if function == "getContract" {
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting contractID.")
		}

		contract, err := t.getContract(stub, args[0])
		if err != nil {
			return nil, err
		}

		return json.Marshal(contract)

	} else if function == "getTransactions" {
		transactions, err := t.getTransactions(stub)
		if err != nil {
			return nil, err
		}

		return json.Marshal(transactions)

	} else if function == "getClaims" {
		claims, err := t.getClaims(stub)
		if err != nil {
			return nil, err
		}

		return json.Marshal(claims)

	} else if function == "getClaim" {
		if len(args) != 2 {
			return nil, errors.New("Incorrect number of arguments. Expecting policyID and claimID.")
		}

		claim, err := t.getClaim(stub, args[0], args[1])
		if err != nil {
			return nil, err
		}

		return json.Marshal(claim)

	} else if function == "getBalance" {
		return t.getBalance(stub)

	} else {
		log.Errorf("function: %s, args: %s", function, args)
		return nil, errors.New("Received unknown function invocation")
	}
}

func main() {
	err := shim.Start(new(InsuranceFrontingChaincode))
	if err != nil {
		log.Critical("Error starting InsuranceFrontingChaincode: %s", err)
	}
}

func (t *InsuranceFrontingChaincode) demoInit(stub *shim.ChaincodeStub) ([]byte, error) {
	/*	Create example policies 	*/
	log.Debug("Create example contract.")

	exampleContract := contract{MaxCoverage: 10000000, MaxPremium: 1000, Captive: "Bermuda", Reinsurer: "Art"}
	if _, err := t.createContract(stub, exampleContract); err != nil {
		return nil, err
	}

	/*	Create example coins 		*/
	log.Debug("Create example coins.")

	if _, err := t.emitCoins(stub, "Allianz", 0); err != nil {
		return nil, err
	}

	if _, err := t.emitCoins(stub, "Bermuda", 6000000); err != nil {
		return nil, err
	}

	if _, err := t.emitCoins(stub, "Art", 0); err != nil {
		return nil, err
	}

	if _, err := t.emitCoins(stub, "Nigeria", 10000); err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *InsuranceFrontingChaincode) incrementAndGetCounter(stub *shim.ChaincodeStub, counterName string) (result uint64, err error) {
	if contractIDBytes, err := stub.GetState(counterName); err != nil {
		log.Errorf("Failed retrieving %s.", counterName)
		return result, err
	} else {
		result, _ = strconv.ParseUint(string(contractIDBytes), 10, 64)
	}
	result++
	if err = stub.PutState(counterName, []byte(strconv.FormatUint(result, 10))); err != nil {
		log.Errorf("Failed saving %s!", counterName)
		return result, err
	}
	return result, err
}

func (t *InsuranceFrontingChaincode) createContract(stub *shim.ChaincodeStub, contract_ contract) ([]byte, error) {
	counter, err := t.incrementAndGetCounter(stub, "ContractsCounter")
	if err != nil {
		return nil, err
	}

	contract_.ID = strconv.FormatUint(counter, 10)

	if ok, err := stub.InsertRow("Contracts", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: contract_.ID}},
			&shim.Column{Value: &shim.Column_Uint64{Uint64: contract_.MaxCoverage}},
			&shim.Column{Value: &shim.Column_Uint64{Uint64: contract_.MaxPremium}},
			&shim.Column{Value: &shim.Column_String_{String_: contract_.Captive}},
			&shim.Column{Value: &shim.Column_String_{String_: contract_.Reinsurer}},
			&shim.Column{Value: &shim.Column_Uint64{Uint64: 0}},
			&shim.Column{Value: &shim.Column_Uint64{Uint64: 0}},
			&shim.Column{Value: &shim.Column_Uint64{Uint64: 0}},
			&shim.Column{Value: &shim.Column_Uint64{Uint64: 0}}},
	}); !ok {
		log.Error("Failed inserting new contract!")
		return nil, err
	}

	return nil, nil
}

func (t *InsuranceFrontingChaincode) updateContract(stub *shim.ChaincodeStub, contract_ contract) ([]byte, error) {
	if ok, err := stub.ReplaceRow("Contracts", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: contract_.ID}},
			&shim.Column{Value: &shim.Column_Uint64{Uint64: contract_.MaxCoverage}},
			&shim.Column{Value: &shim.Column_Uint64{Uint64: contract_.MaxPremium}},
			&shim.Column{Value: &shim.Column_String_{String_: contract_.Captive}},
			&shim.Column{Value: &shim.Column_String_{String_: contract_.Reinsurer}},
			&shim.Column{Value: &shim.Column_Uint64{Uint64: contract_.CurrentTotalCoverage}},
			&shim.Column{Value: &shim.Column_Uint64{Uint64: contract_.CurrentTotalPremium}},
			&shim.Column{Value: &shim.Column_Uint64{Uint64: contract_.CurrentPaidClaim}},
			&shim.Column{Value: &shim.Column_Uint64{Uint64: contract_.CurrentPaidPremium}}},
	}); !ok {
		log.Error("Failed updating contract!")
		return nil, err
	}

	return nil, nil
}

func (t *InsuranceFrontingChaincode) getContracts(stub *shim.ChaincodeStub) (contracts []contract, err error) {
	rows, err := stub.GetRows("Contracts", []shim.Column{})
	if err != nil {
		message := "Failed retrieving contracts. Error: " + err.Error()
		log.Error(message)
		return nil, errors.New(message)
	}

	callerCompany, err := t.getCallerCompany(stub)
	if err != nil {
		return nil, err
	}

	for row := range rows {
		var result contract
		result.ID = row.Columns[0].GetString_()
		result.MaxCoverage = row.Columns[1].GetUint64()
		result.MaxPremium = row.Columns[2].GetUint64()
		result.Captive = row.Columns[3].GetString_()
		result.Reinsurer = row.Columns[4].GetString_()
		result.CurrentTotalCoverage = row.Columns[5].GetUint64()
		result.CurrentTotalPremium = row.Columns[6].GetUint64()
		result.CurrentPaidClaim = row.Columns[7].GetUint64()
		result.CurrentPaidPremium = row.Columns[8].GetUint64()

		log.Debugf("getContracts result includes: %+v", result)

		if callerCompany == result.Captive || callerCompany == result.Reinsurer {
			contracts = append(contracts, result)
		}
	}

	return contracts, nil
}

func (t *InsuranceFrontingChaincode) getContract(stub *shim.ChaincodeStub, contractID string) (contract, error) {
	log.Debugf("Getting contract for contractID %v", contractID)

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: contractID}}
	columns = append(columns, col1)

	row, err := stub.GetRow("Contracts", columns)
	if err != nil {
		message := "Failed retrieving contract ID " + string(contractID) + ". Error: " + err.Error()
		log.Error(message)
		return contract{}, errors.New(message)
	}

	var result contract
	result.ID = row.Columns[0].GetString_()
	result.MaxCoverage = row.Columns[1].GetUint64()
	result.MaxPremium = row.Columns[2].GetUint64()
	result.Captive = row.Columns[3].GetString_()
	result.Reinsurer = row.Columns[4].GetString_()
	result.CurrentTotalCoverage = row.Columns[5].GetUint64()
	result.CurrentTotalPremium = row.Columns[6].GetUint64()
	result.CurrentPaidClaim = row.Columns[7].GetUint64()
	result.CurrentPaidPremium = row.Columns[8].GetUint64()

	log.Debugf("getContracts returns: %+v", result)

	return result, nil
}

func (t *InsuranceFrontingChaincode) getCallerCompany(stub *shim.ChaincodeStub) (string, error) {
	callerCompany, err := stub.ReadCertAttribute("company")
	if err != nil {
		log.Error("Failed fetching caller's company. Error: " + err.Error())
		return "", err
	}
	log.Debugf("Caller company is: %s", callerCompany)
	return string(callerCompany), nil
}

func (t *InsuranceFrontingChaincode) getCallerRole(stub *shim.ChaincodeStub) (string, error) {
	callerRole, err := stub.ReadCertAttribute("role")
	if err != nil {
		log.Error("Failed fetching caller role. Error: " + err.Error())
		return "", err
	}
	log.Debugf("Caller role is: %s", callerRole)
	return string(callerRole), nil
}

func (t *InsuranceFrontingChaincode) joinChain(stub *shim.ChaincodeStub, policyID string) ([]byte, error) {
	log.Debugf("function: %s, args: %s", "join", policyID)

	role, err := t.getCallerRole(stub)
	if err != nil {
		return nil, err
	}

	company, err := t.getCallerCompany(stub)
	if err != nil {
		return nil, err
	}

	policy_, err := t.getPolicy(stub, policyID)
	if err != nil {
		return nil, err
	}
	log.Debugf("Policy is: %+v", policy_)

	switch {
	case role == "fronter" && policy_.FrontingChain.Fronter == "":
		policy_.FrontingChain.Fronter = company

	case role == "affiliate" && policy_.FrontingChain.Affiliate == "":
		policy_.FrontingChain.Affiliate = company

	default:
		log.Error("Caller's role is incorrect!")
	}

	if ok, err := t.updatePolicy(stub, policy_); !ok || err != nil {
		return nil, err
	}

	return json.Marshal(policy_)
}

func (t *InsuranceFrontingChaincode) emitCoins(stub *shim.ChaincodeStub, companyName string, amount uint64) ([]byte, error) {
	log.Debug("Emitting coins for " + companyName)

	callerRole, err := t.getCallerRole(stub)
	if err != nil {
		return nil, err
	}

	if callerRole != "bank" {
		message := "Caller role: " + callerRole + ", required 'bank'."
		log.Error(message)
		return nil, errors.New(message)
	}

	bankName, err := t.getCallerCompany(stub)
	if err != nil {
		return nil, err
	}

	err = stub.PutState(bankName, []byte(strconv.FormatUint(amount, 10)))
	if err != nil {
		log.Debug("Failed initializing coins at " + bankName)
		return nil, err
	}

	transaction_ := transaction{ID: "", From: bankName, To: companyName, Amt: amount, Purpose: "coins emission"}
	if _, err := t.transact(stub, transaction_); err != nil {
		log.Error("Emitting transaction failed!")
		return nil, err
	}

	return nil, nil
}

func (t *InsuranceFrontingChaincode) getBalance(stub *shim.ChaincodeStub) ([]byte, error) {
	companyName, err := t.getCallerCompany(stub)
	if err != nil {
		return nil, err
	}

	bytes, err := stub.GetState(companyName)
	if err != nil {
		log.Debug("Failed reading balance for " + companyName)
		return nil, err
	}

	balance, err := strconv.ParseUint(string(bytes), 10, 64)
	if err != nil {
		log.Debug("Failed parsing balance for " + companyName)
		return nil, err
	}

	return json.Marshal(balance)
}

func (t *InsuranceFrontingChaincode) transact(stub *shim.ChaincodeStub, transaction_ transaction) ([]byte, error) {
	b, _ := json.Marshal(transaction_)
	log.Debug("Started transaction: " + string(b))

	var fromValue uint64
	if fromBytes, err := stub.GetState(transaction_.From); err != nil {
		log.Error("Failed retrieving fromValue for " + transaction_.From)
		return nil, errors.New("Failed to get state")
	} else {
		fromValue, _ = strconv.ParseUint(string(fromBytes), 10, 64)
	}

	var toValue uint64
	if toBytes, err := stub.GetState(transaction_.To); err != nil {
		log.Error("Failed retrieving toValue for " + transaction_.To)
		return nil, errors.New("Failed to get state")
	} else {
		toValue, _ = strconv.ParseUint(string(toBytes), 10, 64)
	}

	if fromValue < transaction_.Amt {
		log.Error("Not enough value!")
		return nil, errors.New("Not enough value!")
	}

	fromValue = fromValue - transaction_.Amt
	toValue = toValue + transaction_.Amt

	err := stub.PutState(transaction_.From, []byte(strconv.FormatUint(fromValue, 10)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(transaction_.To, []byte(strconv.FormatUint(toValue, 10)))
	if err != nil {
		return nil, err
	}

	transactionID, err := t.incrementAndGetCounter(stub, "TransactionsCounter")
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	transaction_.ID = strconv.FormatUint(transactionID, 10)

	if ok, err := stub.InsertRow("Transactions", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: transaction_.ID}},
			&shim.Column{Value: &shim.Column_String_{String_: transaction_.From}},
			&shim.Column{Value: &shim.Column_String_{String_: transaction_.To}},
			&shim.Column{Value: &shim.Column_Uint64{Uint64: transaction_.Amt}},
			&shim.Column{Value: &shim.Column_String_{String_: transaction_.Purpose}}},
	}); !ok {
		return nil, err
	}

	return nil, nil
}

func (t *InsuranceFrontingChaincode) payPremium(stub *shim.ChaincodeStub, policyID string) ([]byte, error) {
	/* Get the policy to pay premium for */
	policy_, err := t.getPolicy(stub, policyID)
	if err != nil {
		log.Error("Pay failure. PolicyID: " + policyID)
		return nil, err
	}
	chain := &policy_.FrontingChain

	/* Get the corresponding contract */
	contract_, err := t.getContract(stub, policy_.ContractID)
	if err != nil {
		log.Error("Pay failure. ContractID: " + contract_.ID)
		return nil, err
	}

	/* Identify caller Company - only Affiliate company can initiate the payment */
	callerCompany, err := t.getCallerCompany(stub)
	if err != nil {
		return nil, err
	}
	if callerCompany != chain.Affiliate {
		message := "Caller company: " + callerCompany + ", required: " + chain.Affiliate
		log.Error(message)
		return nil, errors.New(message)
	}

	/* Check if all the members of the fronting chain have been initialized */
	if chain.Captive == "" || chain.Reinsurer == "" || chain.Fronter == "" || chain.Affiliate == "" {
		chainJson, _ := json.Marshal(chain)
		message := "Fronting chain is missing a member! Check it: " + string(chainJson)
		log.Error(message)
		return chainJson, errors.New(message)
	}

	/*  That is how the money are going to be transferred:  */
	/*	Captive <- Reinsurer <- Fronter <- Affiliate 	*/
	amount := policy_.Premium
	purpose := "premium." + policy_.PolicyID

	aff_fro := transaction{ID: "", From: chain.Affiliate, To: chain.Fronter, Amt: amount, Purpose: purpose}
	if _, err := t.transact(stub, aff_fro); err != nil {
		return nil, err
	}

	fro_rei := transaction{ID: "", From: chain.Fronter, To: chain.Reinsurer, Amt: amount, Purpose: purpose}
	if _, err := t.transact(stub, fro_rei); err != nil {
		return nil, err
	}

	rei_cap := transaction{ID: "", From: chain.Reinsurer, To: chain.Captive, Amt: amount, Purpose: purpose}
	if _, err := t.transact(stub, rei_cap); err != nil {
		return nil, err
	}

	/* Update premium counters in policy and contract */
	policy_.PaidPremium += amount
	contract_.CurrentPaidPremium += amount

	if _, err := t.updatePolicy(stub, policy_); err != nil {
		log.Errorf("Failed updating policy : %s", err.Error())
		return nil, err
	}
	if _, err := t.updateContract(stub, contract_); err != nil {
		log.Errorf("Failed updating contract : %s", err.Error())
		return nil, err
	}

	return nil, nil
}

func (t *InsuranceFrontingChaincode) createClaim(stub *shim.ChaincodeStub, policyID string, amount uint64) ([]byte, error) {
	/* Get the policy to claim for */
	policy_, err := t.getPolicy(stub, policyID)
	if err != nil {
		log.Error("Claim failure. PolicyID: " + policyID)
		return nil, err
	}
	chain := &policy_.FrontingChain

	/* Get the corresponding contract */
	contract_, err := t.getContract(stub, policy_.ContractID)
	if err != nil {
		log.Error("Claim failure. ContractID: " + contract_.ID)
		return nil, err
	}

	/* Identify caller Company - only Affiliate company can initiate the claim */
	callerCompany, err := t.getCallerCompany(stub)
	if err != nil {
		return nil, err
	}
	if callerCompany != chain.Affiliate {
		message := "Caller company: " + callerCompany + ", required: " + chain.Affiliate
		log.Error(message)
		return nil, errors.New(message)
	}

	/* Check if all the members of the fronting chain have been initialized */
	if chain.Captive == "" || chain.Reinsurer == "" || chain.Fronter == "" || chain.Affiliate == "" {
		chainJson, _ := json.Marshal(chain)
		message := "Fronting chain is missing a member! Check it: " + string(chainJson)
		log.Error(message)
		return chainJson, errors.New(message)
	}

	/* Update premium counters in policy and contract */
	if policy_.PaidClaim += amount; policy_.PaidClaim > policy_.Coverage {
		message := "Claim is too high for policy " + policy_.PolicyID
		log.Error(message)
		return nil, errors.New(message)
	}
	if contract_.CurrentPaidClaim += amount; contract_.CurrentPaidClaim > contract_.MaxCoverage {
		message := "Claim is too high for contract " + contract_.ID
		log.Error(message)
		return nil, errors.New(message)
	}

	claimID, err := t.incrementAndGetCounter(stub, "ClaimsCounter")
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	claim_ := claim{PolicyID: policy_.PolicyID, ClaimID: strconv.FormatUint(claimID, 10), Amt: amount}
	if ok, err := stub.InsertRow("Claims", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: claim_.PolicyID}},
			&shim.Column{Value: &shim.Column_String_{String_: claim_.ClaimID}},
			&shim.Column{Value: &shim.Column_Uint64{Uint64: claim_.Amt}},
			&shim.Column{Value: &shim.Column_Bool{Bool: claim_.ApprovedByCaptive}},
			&shim.Column{Value: &shim.Column_Bool{Bool: claim_.ApprovedByReinsurer}},
			&shim.Column{Value: &shim.Column_Bool{Bool: claim_.ApprovedByFronter}}},
	}); !ok {
		return nil, err
	}

	if _, err := t.updatePolicy(stub, policy_); err != nil {
		log.Errorf("Failed updating policy : %s", err.Error())
		return nil, err
	}
	if _, err := t.updateContract(stub, contract_); err != nil {
		log.Errorf("Failed updating contract : %s", err.Error())
		return nil, err
	}

	return nil, nil
}

func (t *InsuranceFrontingChaincode) updateClaim(stub *shim.ChaincodeStub, claim_ claim) (bool, error) {
	return stub.ReplaceRow("Claims", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: claim_.PolicyID}},
			&shim.Column{Value: &shim.Column_String_{String_: claim_.ClaimID}},
			&shim.Column{Value: &shim.Column_Uint64{Uint64: claim_.Amt}},
			&shim.Column{Value: &shim.Column_Bool{Bool: claim_.ApprovedByCaptive}},
			&shim.Column{Value: &shim.Column_Bool{Bool: claim_.ApprovedByReinsurer}},
			&shim.Column{Value: &shim.Column_Bool{Bool: claim_.ApprovedByFronter}}},
	})
}

func (t *InsuranceFrontingChaincode) getClaim(stub *shim.ChaincodeStub, policyID string, claimID string) (claim, error) {
	log.Debugf("Getting claim for claimID %v", claimID)

	var columns []shim.Column
	columnPolicyID := shim.Column{Value: &shim.Column_String_{String_: policyID}}
	columns = append(columns, columnPolicyID)
	columnClaimID := shim.Column{Value: &shim.Column_String_{String_: claimID}}
	columns = append(columns, columnClaimID)

	rows, err := stub.GetRows("Claims", columns)
	if err != nil || len(rows) < 1 {
		message := "Failed retrieving claim ID " + string(claimID) + ". Error: " + err.Error()
		log.Error(message)
		return claim{}, errors.New(message)
	}

	row := <-rows
	result := claim{
		PolicyID:            row.Columns[0].GetString_(),
		ClaimID:             row.Columns[1].GetString_(),
		Amt:                 row.Columns[2].GetUint64(),
		ApprovedByCaptive:   row.Columns[3].GetBool(),
		ApprovedByReinsurer: row.Columns[4].GetBool(),
		ApprovedByFronter:   row.Columns[5].GetBool()}

	return result, nil
}

func (t *InsuranceFrontingChaincode) approveClaim(stub *shim.ChaincodeStub, policyID string, claimID string) ([]byte, error) {
	/* Get the policy */
	policy_, err := t.getPolicy(stub, policyID)
	if err != nil {
		log.Error("Approve claim failure. PolicyID: " + policyID)
		return nil, err
	}

	/* Get the corresponding contract */
	contract_, err := t.getContract(stub, policy_.ContractID)
	if err != nil {
		log.Error("Approve claim failure. ContractID: " + contract_.ID)
		return nil, err
	}

	/* Get the claim */
	claim_, err := t.getClaim(stub, policyID, claimID)
	if err != nil {
		log.Error("Approve claim failure. PolicyID: " + policyID + ", ClaimID: " + claimID)
		return nil, err
	}

	/* Identify caller Company - only policy Captive, Reinsurer and Fronter companies allowed */
	caller, err := t.getCallerCompany(stub)
	if err != nil {
		return nil, err
	}

	/*   That is how the claim should be approved:  */
	/*	Captive -> Reinsurer -> Fronter 	*/
	c := &policy_.FrontingChain
	if !claim_.ApprovedByCaptive {
		if c.Captive == caller {
			claim_.ApprovedByCaptive = true
		} else {
			message := "Claim approve failure. Expected captive: " + c.Captive + ", got: " + caller
			log.Error(message)
			return nil, errors.New(message)
		}
	} else if !claim_.ApprovedByReinsurer {
		if c.Reinsurer == caller {
			claim_.ApprovedByReinsurer = true
		} else {
			message := "Claim approve failure. Expected reinsurer: " + c.Reinsurer + ", got: " + caller
			log.Error(message)
			return nil, errors.New(message)
		}
	} else if !claim_.ApprovedByFronter {
		if c.Fronter == caller {
			claim_.ApprovedByFronter = true
		} else {
			message := "Claim approve failure. Expected fronter: " + c.Fronter + ", got: " + caller
			log.Error(message)
			return nil, errors.New(message)
		}
	} else {
		message := "Claim approve failure. Expected " +
			c.Captive + ", " + c.Reinsurer + " or " + c.Fronter + ", got: " + caller
		log.Error(message)
		return nil, errors.New(message)
	}

	if _, err := t.updateClaim(stub, claim_); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if claim_.ApprovedByFronter && claim_.ApprovedByReinsurer && claim_.ApprovedByCaptive {
		amount := claim_.Amt
		/*  That is how the money are going to be transferred:  */
		/*	Captive -> Reinsurer -> Fronter -> Affiliate 	*/

		cap_rei := transaction{
			From:    c.Captive,
			To:      c.Reinsurer,
			Amt:     amount,
			Purpose: claim_.PolicyID + claim_.ClaimID}
		if _, err := t.transact(stub, cap_rei); err != nil {
			return nil, err
		}

		rei_fro := transaction{
			From:    c.Reinsurer,
			To:      c.Fronter,
			Amt:     amount,
			Purpose: claim_.PolicyID + claim_.ClaimID}
		if _, err := t.transact(stub, rei_fro); err != nil {
			return nil, err
		}

		fro_aff := transaction{
			From:    c.Fronter,
			To:      c.Affiliate,
			Amt:     amount,
			Purpose: claim_.PolicyID + claim_.ClaimID}
		if _, err := t.transact(stub, fro_aff); err != nil {
			return nil, err
		}

		/* Update claim counters in policy and contract */
		policy_.PaidClaim += amount
		contract_.CurrentPaidClaim += amount

		if _, err := t.updatePolicy(stub, policy_); err != nil {
			log.Errorf("Failed updating policy : %s", err.Error())
			return nil, err
		}
		if _, err := t.updateContract(stub, contract_); err != nil {
			log.Errorf("Failed updating contract : %s", err.Error())
			return nil, err
		}
	}

	return nil, nil
}

func (t *InsuranceFrontingChaincode) getClaims(stub *shim.ChaincodeStub) (claims []claim, err error) {
	rows, err := stub.GetRows("Claims", []shim.Column{})
	if err != nil {
		message := "Failed retrieving claims. Error: " + err.Error()
		log.Error(message)
		return nil, errors.New(message)
	}

	caller, err := t.getCallerCompany(stub)
	if err != nil {
		return nil, err
	}

	for row := range rows {
		result := claim{
			PolicyID:            row.Columns[0].GetString_(),
			ClaimID:             row.Columns[1].GetString_(),
			Amt:                 row.Columns[2].GetUint64(),
			ApprovedByCaptive:   row.Columns[3].GetBool(),
			ApprovedByReinsurer: row.Columns[4].GetBool(),
			ApprovedByFronter:   row.Columns[5].GetBool()}

		policy_, err := t.getPolicy(stub, result.PolicyID)
		if err != nil {
			log.Error("Failed retrieving policy for claim " + result.ClaimID)
			continue
		}

		log.Debugf("getClaims result includes: %+v", result)

		c := &policy_.FrontingChain
		if caller == c.Captive || caller == c.Reinsurer || caller == c.Fronter || caller == c.Affiliate {
			claims = append(claims, result)
		}
	}

	return claims, nil
}

func (t *InsuranceFrontingChaincode) getPolicies(stub *shim.ChaincodeStub, contractID string) ([]policy, error) {
	var columns []shim.Column
	if contractID != "" {
		columnContractIDs := shim.Column{Value: &shim.Column_String_{String_: contractID}}
		columns = append(columns, columnContractIDs)
	}

	rows, err := stub.GetRows("Policies", columns)
	if err != nil {
		message := "Failed retrieving policies. Error: " + err.Error()
		log.Error(message)
		return nil, errors.New(message)
	}

	var policies []policy

	for row := range rows {
		result := policy{
			ContractID:  row.Columns[0].GetString_(),
			PolicyID:    row.Columns[1].GetString_(),
			Coverage:    row.Columns[2].GetUint64(),
			Premium:     row.Columns[3].GetUint64(),
			PaidClaim:   row.Columns[4].GetUint64(),
			PaidPremium: row.Columns[5].GetUint64()}

		bytes := row.Columns[6].GetBytes()
		if bytes == nil {
			message := "Failed retrieving fronting chain for policy ID: " + result.PolicyID + "."
			log.Error(message)
			return nil, errors.New(message)
		}

		err = json.Unmarshal(bytes, &result.FrontingChain)
		if err != nil {
			message := "Failed unmarshalling fronting chain for policy ID: " + result.PolicyID + "."
			log.Error(message)
			return nil, errors.New(message)
		}

		log.Debugf("getPolicies result includes: %+v", result)
		policies = append(policies, result)
	}

	return policies, nil
}

func (t *InsuranceFrontingChaincode) getTransactions(stub *shim.ChaincodeStub) ([]transaction, error) {
	callerCompany, err := t.getCallerCompany(stub)
	if err != nil {
		return nil, err
	}

	callerRole, err := t.getCallerRole(stub)
	if err != nil {
		return nil, err
	}

	rows, err := stub.GetRows("Transactions", []shim.Column{})
	if err != nil {
		message := "Failed retrieving transactions. Error: " + err.Error()
		log.Error(message)
		return nil, errors.New(message)
	}

	var transactions []transaction

	for row := range rows {
		var result transaction
		result.ID = row.Columns[0].GetString_()
		result.From = row.Columns[1].GetString_()
		result.To = row.Columns[2].GetString_()
		result.Amt = row.Columns[3].GetUint64()
		result.Purpose = row.Columns[4].GetString_()

		if callerCompany == result.From || callerCompany == result.To || callerRole == "auditor" {
			transactions = append(transactions, result)
		}
	}

	return transactions, nil
}

func (t *InsuranceFrontingChaincode) createPolicy(stub *shim.ChaincodeStub, contractID string, coverage uint64, premium uint64) (bool, error) {
	contract_, err := t.getContract(stub, contractID)
	if err != nil {
		return false, err
	}

	policyID, err := t.incrementAndGetCounter(stub, "PoliciesCounter")
	if err != nil {
		return false, err
	}

	policy_ := policy{
		ContractID:    contract_.ID,
		PolicyID:      strconv.FormatUint(policyID, 10),
		Coverage:      coverage,
		Premium:       premium,
		FrontingChain: frontingChain{Captive: contract_.Captive, Reinsurer: contract_.Reinsurer}}
	contract_.CurrentTotalCoverage += coverage
	contract_.CurrentTotalPremium += premium
	if contract_.CurrentTotalCoverage > contract_.MaxCoverage ||
		contract_.CurrentTotalPremium > contract_.MaxPremium {
		message := "Coverage or premium overlimit."
		log.Error(message)
		return false, errors.New(message)
	}

	frontingChain_, err := json.Marshal(policy_.FrontingChain)
	if err != nil {
		message := "Failed marshalling fronting chain. Error: " + err.Error()
		log.Error(message)
		return false, errors.New(message)
	}

	if ok, err := stub.InsertRow("Policies", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: policy_.ContractID}}, // ContractID    string
			&shim.Column{Value: &shim.Column_String_{String_: policy_.PolicyID}},   // PolicyID      string
			&shim.Column{Value: &shim.Column_Uint64{Uint64: policy_.Coverage}},     // Coverage      uint64
			&shim.Column{Value: &shim.Column_Uint64{Uint64: policy_.Premium}},      // Premium       uint64
			&shim.Column{Value: &shim.Column_Uint64{Uint64: 0}},                    // PaidClaim     uint64
			&shim.Column{Value: &shim.Column_Uint64{Uint64: 0}},                    // PaidPremium   uint64
			&shim.Column{Value: &shim.Column_Bytes{Bytes: frontingChain_}}},        // FrontingChain frontingChain
	}); !ok {
		return ok, err
	}

	if _, err := t.updateContract(stub, contract_); err != nil {
		return false, err
	}

	return true, nil
}

func (t *InsuranceFrontingChaincode) updatePolicy(stub *shim.ChaincodeStub, policy_ policy) (bool, error) {
	frontingChain_, err := json.Marshal(policy_.FrontingChain)
	if err != nil {
		message := "Failed marshalling fronting chain for policy ID: " + policy_.PolicyID + ". Error: " + err.Error()
		log.Error(message)
		return false, errors.New(message)
	}

	if ok, err := stub.ReplaceRow("Policies", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: policy_.ContractID}}, // ContractID    string
			&shim.Column{Value: &shim.Column_String_{String_: policy_.PolicyID}},   // PolicyID      string
			&shim.Column{Value: &shim.Column_Uint64{Uint64: policy_.Coverage}},     // Coverage      uint64
			&shim.Column{Value: &shim.Column_Uint64{Uint64: policy_.Premium}},      // Premium       uint64
			&shim.Column{Value: &shim.Column_Uint64{Uint64: 0}},                    // PaidClaim     uint64
			&shim.Column{Value: &shim.Column_Uint64{Uint64: 0}},                    // PaidPremium   uint64
			&shim.Column{Value: &shim.Column_Bytes{Bytes: frontingChain_}}},        // FrontingChain frontingChain
	}); !ok {
		return false, err
	} else {
		return true, err
	}
}

func (t *InsuranceFrontingChaincode) getPolicy(stub *shim.ChaincodeStub, policyID string) (policy, error) {
	log.Debugf("Getting policy for policyID %v", policyID)

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: policyID}}
	columns = append(columns, col1)

	row, err := stub.GetRow("Policies", columns)
	if err != nil {
		message := "Failed retrieving policy ID " + string(policyID) + ". Error: " + err.Error()
		log.Error(message)
		return policy{}, errors.New(message)
	}

	result := policy{
		ContractID:  row.Columns[0].GetString_(),
		PolicyID:    row.Columns[1].GetString_(),
		Coverage:    row.Columns[2].GetUint64(),
		Premium:     row.Columns[3].GetUint64(),
		PaidClaim:   row.Columns[4].GetUint64(),
		PaidPremium: row.Columns[5].GetUint64()}

	bytes := row.Columns[6].GetBytes()
	if bytes == nil {
		message := "Failed retrieving fronting chain for policy ID: " + policyID + "."
		log.Error(message)
		return policy{}, errors.New(message)
	}

	err = json.Unmarshal(bytes, &result.FrontingChain)
	if err != nil {
		message := "Failed unmarshalling fronting chain for policy ID: " + policyID + "."
		log.Error(message)
		return policy{}, errors.New(message)
	}

	log.Debugf("getPolicy returns: %+v", result)

	return result, nil
}

