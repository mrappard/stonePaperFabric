package main

import (
	"fmt"
	"strconv"
	"time"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/crypto/attr"

)


type stonePaper struct {
	DocHash 			string `json:"docHash"`
	Database  		int    `json:"database"`
	Time          string `json:"time"`
	Creator       string `json:"creator"`
	SubContract   string `json:"subContract"`
	ContractType  int `json:"contractType"`
}

// StonePaperChaincode
// ===========================
type StonePaperChaincode struct {
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(StonePaperChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}


func (t *StonePaperChaincode) setStateToAttributes(stub shim.ChaincodeStubInterface, args []string) error {
	attrHandler, err := attr.NewAttributesHandlerImpl(stub)
	if err != nil {
		return err
	}
	for _, att := range args {
		fmt.Println("Writing attribute " + att)
		attVal, err := attrHandler.GetValue(att)
		if err != nil {
			return err
		}
		err = stub.PutState(att, attVal)
		if err != nil {
			return err
		}
	}
	return nil
}


// Init initializes chaincode
// ===========================
func (t *StonePaperChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	return nil, nil
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *StonePaperChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "createDoc" { //create a new doc
		return t.createDoc(stub, args)
	} else if function == "setDatabase" { //change owner of a specific marble
		return t.setDatabase(stub, args)
	} else if function == "setName" { //transfer all marbles of a certain color
		return t.setName(stub, args)
	}
	fmt.Println("invoke did not find func: " + function) //error
	return nil, errors.New("Received unknown function invocation")
}


// Query callback representing the query of a chaincode
func (t *StonePaperChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "getDoc" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var A string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting Hash of document")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get document for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"No document for " + A + " in Blockchain\"}"
		return nil, errors.New(jsonResp)
	}

	//jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	//fmt.Printf("Query Response:%s\n", jsonResp)
	return Avalbytes, nil
}


func GetCertAttribute(stub shim.ChaincodeStubInterface, attributeName string) (string, error) {
 fmt.Println("Entering GetCertAttribute")
 attr, err := stub.ReadCertAttribute(attributeName)
 if err != nil {
 return "", errors.New("Couldn't get attribute " + attributeName + ". Error: " + err.Error())
 }
 attrString := string(attr)
 return attrString, nil
}


// ============================================================
// initStonepaper - create a paper, store into chaincode state
// ============================================================
func (t *StonePaperChaincode) createDoc(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	//   0       	1       		2     					3
	// "DocHash", "Database", "SubContract", "ContractType"
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init marble")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}

	if len(args[4]) <= 0 {
		return nil, errors.New("5th argument must be a non-empty string")
	}

	DocHash := args[0]
	Database, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("2rd argument must be a numeric string")
	}
	SubContract := args[2]
	ContractType, err := strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("4rd argument must be a numeric string")
	}
	timerValue := time.Now()
	TimeV := timerValue.String()

	Creator := args[4]

	TestInfo,err := GetCertAttribute(stub,"username")
 	if err != nil {
 		TestInfo = "Failed " + err.Error()
 	}

	Creator = Creator + "-" + TestInfo

	// ==== Check if doc with matching hash exists already exists ====
	docAsBytes, err := stub.GetState(DocHash)
	if err != nil {
		return nil, errors.New("Failed to get doc: " + err.Error())
	} else if docAsBytes != nil {
		fmt.Println("This Hash already exists: " + DocHash)
		return nil, errors.New("This Hash already exists: " + DocHash)
	}

/*
	// ==== Create stonePaper object and marshal to JSON ====
	objectType := "stonePaper"
	marble := &stonePaper{DocHash, Database, TimeV, Creator, SubContract, ContractType}
	marbleJSONasBytes, err := json.Marshal(marble)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	*/

	stonePaperJSONasString := `{"DocHash":"`+DocHash+`","Database":`+strconv.Itoa(Database)+`,"Time":"`+TimeV+`","Creator":"`+Creator+`","SubContract": "`+SubContract+`","ContractType":`+strconv.Itoa(ContractType)+`}`
	stonePaperJSONasBytes := []byte(stonePaperJSONasString)

	//Alternatively, build the marble json string manually if you don't want to use struct marshalling
	//marbleJSONasString := `{"docType":"Marble",  "name": "` + marbleName + `", "color": "` + color + `", "size": ` + strconv.Itoa(size) + `, "owner": "` + owner + `"}`
	//marbleJSONasBytes := []byte(str)

	// === Save marble to state ===
	err = stub.PutState(DocHash, stonePaperJSONasBytes)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	// ==== Marble saved and indexed. Return success ====
	fmt.Println("- end init document")
	return nil, err
}

func (t *StonePaperChaincode) setDatabase(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}

func (t *StonePaperChaincode) setName(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
