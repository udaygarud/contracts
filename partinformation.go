package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	

	"github.com/hyperledger/fabric/core/chaincode/shim"
)


type PartInfo struct {
	ContractNumber   string  `json:"contractNumber"`
	ContractDescription   string  `json:"contractDescription"`
	ContractOwnedby   string  `json:"contractOwnedby"`
	FilehashKey   string  `json:"filehashKey"`
	NumberofParts float64 `json:"numberofParts"`
	ContractAmount float64 `json:"contractAmount"`
	ContractDuedate   string  `json:"contractDuedate"`
}


type Assigneeinfo struct {
	UserID float64 `json:"userID"`
	IsSigned   bool  `json:"isSigned"`
	SignedDate float64 `json:"signedDate"`
	Status   string  `json:"status"`
	
}
  


// PartInformation example simple Chaincode implementation
type PartInformation struct {
}

func main() {
	err := shim.Start(new(PartInformation))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *PartInformation) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *PartInformation) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "addpartInformation" {
		return t.addpartInformation(stub, args)
	} else if function == "addAssignee" {
		return t.addAssignee(stub, args)
	} else if function == "signbyAssignee" {
		return t.signbyAssignee(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *PartInformation) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "readpartInformation" {
		return t.readpartInformation(stub, args)
	} else if function == "readAssigneeInformation" {
		return t.readAssigneeInformation(stub, args)
	} else if function == "readAssigneeStatus" {
		return t.readAssigneeStatus(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *PartInformation) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *PartInformation) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func (t *PartInformation) addpartInformation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("adding part information")
	if len(args) != 7 {
		return nil, errors.New("Incorrect Number of arguments.Expecting 7 for addpartInformation")
	}
	amt, err := strconv.ParseFloat(args[4], 64)
	nup, err := strconv.ParseFloat(args[5], 64)

	part := PartInfo{
		ContractNumber:   args[0],
		ContractDescription: args[1],
		ContractOwnedby: args[2],
		FilehashKey: args[3],
		NumberofParts: amt,
		ContractAmount: nup,
		ContractDuedate: args[6],
	}

	bytes, err := json.Marshal(part)
	if err != nil {
		fmt.Println("Error marshaling part")
		return nil, errors.New("Error marshaling part")
	}

	err = stub.PutState(part.ContractNumber, bytes)
	if err != nil {
		return nil, err
}
return nil, nil
}

func (t *PartInformation) addAssignee(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("adding part information")
	if len(args) != 4 {
		return nil, errors.New("Incorrect Number of arguments.Expecting 4 for addAssignee")
	}
	id, err := strconv.ParseFloat(args[0], 64)
	da, err := strconv.ParseBool(args[1])

	assign := Assigneeinfo{
		UserID:   id,
		IsSigned: da,
		SignedDate: args[2],
		Status: args[3],
	}
	
	
	bytes, err := json.Marshal(assign)
	if err != nil {
		fmt.Println("Error marshaling assign")
		return nil, errors.New("Error marshaling assign")
	}

	err = stub.PutState(assign.UserID, bytes)
	if err != nil {
		return nil, err
}
return nil, nil
}

func (t *PartInformation) signbyAssignee(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("adding part information")
	if len(args) != 4 {
		return nil, errors.New("Incorrect Number of arguments.Expecting 4 for addAssignee")
	}
	key := args[0]
	assignee, err := stub.GetState(args[0])
	if err != nil {
		fmt.Println("Error retrieving " + key)
		return nil, errors.New("Error retrieving " + key)
	}
	
	// Delete the key from the state in ledger
	newbytes, err := stub.DelState(key)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}
	
	id, err := strconv.ParseFloat(args[0], 64)
	da, err := strconv.ParseBool(args[2])

	assign := Assigneeinfo{
		UserID:   id,
		IsSigned: args[1],
		SignedDate: da,
		Status: args[3],
	}
	
	
	bytes, err := json.Marshal(assign)
	if err != nil {
		fmt.Println("Error marshaling assign")
		return nil, errors.New("Error marshaling assign")
	}

	err = stub.PutState(assign.UserID, bytes)
	if err != nil {
		return nil, err
}
return nil, nil
}

func (t *PartInformation) readpartInformation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("read() is running")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. expecting 1")
	}

	key := args[0] // name of Entity
	fmt.Println("key is ")
	fmt.Println(key)
	bytes, err := stub.GetState(args[0])
	fmt.Println(bytes)
	if err != nil {
		fmt.Println("Error retrieving " + key)
		return nil, errors.New("Error retrieving " + key)
	}
	return bytes, nil
}

func (t *PartInformation) readAssigneeInformation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("read() is running")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. expecting 1")
	}

	key := args[0] // name of Entity
	fmt.Println("key is ")
	fmt.Println(key)
	bytes, err := stub.GetState(args[0])
	fmt.Println(bytes)
	if err != nil {
		fmt.Println("Error retrieving " + key)
		return nil, errors.New("Error retrieving " + key)
	}
	return bytes, nil
}

func (t *PartInformation) readAssigneeStatus(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("read() is running")

	if len(args) != 1 {
		return "false", errors.New("Incorrect number of arguments. expecting 1")
	}

	key := args[0] // name of Entity
	fmt.Println("key is ")
	fmt.Println(key)
	bytes, err := stub.GetState(args[0])
	fmt.Println(bytes)
	if err != nil {
		fmt.Println("Error retrieving " + key)
		return "false", errors.New("Error retrieving " + key)
	}
	
	res := Assigneeinfo{}
	json.Unmarshal(bytes, &res)
	if res.IsSigned == true{
		
		return "true", nil				//all stop a marble by this name exists
	}
	
	return "false", nil
}