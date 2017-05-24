package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)


type stonePaper struct {
	DocHash 			string `json:"docHash"`
	Database  		int    `json:"database"`
	Time          string `json:"time"`
	Creator       string `json:"creator"`
	SubContract   string `json:"subContract"`
	ContractType  int `json:"contractType"`
}


// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "createDoc" { //create a new doc
		return t.createDoc(stub, args)
	} else if function == "setDatabase" { //change owner of a specific marble
		return t.setDatabase(stub, args)
	} else if function == "setName" { //transfer all marbles of a certain color
		return t.setName(stub, args)
	} else if function == "getDoc"{
		return t.getDoc(stub, args)
	}
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}


// ============================================================
// initMarble - create a new marble, store into chaincode state
// ============================================================
func (t *SimpleChaincode) createDoc(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	//   0       	1       		2     					3
	// "DocHash", "Database", "SubContract", "ContractType"
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init marble")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}

	DocHash := args[0]
	Database, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("2rd argument must be a numeric string")
	}
	SubContract := args[2]
	ContractType, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("4rd argument must be a numeric string")
	}
	t := time.Now()
	TimeV := t.String()
	Creator := shim.GetCreator()


	// ==== Check if doc with matching hash exists already exists ====
	docAsBytes, err := stub.GetState(DocHash)
	if err != nil {
		return shim.Error("Failed to get doc: " + err.Error())
	} else if docAsBytes != nil {
		fmt.Println("This Hash already exists: " + DocHash)
		return shim.Error("This Hash already exists: " + DocHash)
	}

/*
	// ==== Create stonePaper object and marshal to JSON ====
	objectType := "stonePaper"
	marble := &stonePaper{DocHash, Database, TimeV, Creator, SubContract, ContractType}
	marbleJSONasBytes, err := json.Marshal(marble)
	if err != nil {
		return shim.Error(err.Error())
	}
	*/

	stonePaperJSONasString := `{"DocHash":"`DocHash`","Database": "`+Database+`", "Time":"`+TimeV+`","Creator":"`+Creator+`", "SubContract": "`+SubContract+`","ContractType": "`+ContractType+`"}`
	stonePaperJSONasBytes := []byte(str)

	//Alternatively, build the marble json string manually if you don't want to use struct marshalling
	//marbleJSONasString := `{"docType":"Marble",  "name": "` + marbleName + `", "color": "` + color + `", "size": ` + strconv.Itoa(size) + `, "owner": "` + owner + `"}`
	//marbleJSONasBytes := []byte(str)

	// === Save marble to state ===
	err = stub.PutState(DocHash, stonePaperJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	// ==== Marble saved and indexed. Return success ====
	fmt.Println("- end init document")
	return shim.Success(nil)
}

func (t *SimpleChaincode) setDatabase(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

func (t *SimpleChaincode) setName(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

func (t *SimpleChaincode) getDoc(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "DocHash"
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	DocHash := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"stonePaper\",\"DocHash\":\"%s\"}}", DocHash)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}
