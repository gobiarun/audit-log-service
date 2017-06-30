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

// PatientChaincode - simple Chaincode implementation for Patient
type PatientChaincode struct {
}

type Patient struct {
	Id		  		string 'json:"id"'
	Name     		string 'json:"name"'    //the fieldtags are needed to keep case from bouncing around
	Gender      	string 'json:"gender"'
	Age      		string 'json:"age"'	
	PhysicianName	string 'json:"PhysicianName"'
	Status			PatientStatus 'json:"status"'	
}

type PhysicianRelation struct {
	PatientId		string 'json:"patientId"'
	PhysicianId		string 'json:"physicianId"'
	Name     		string 'json:"name"'    
	Gender      	string 'json:"gender"'
	Age      		string 'json:"age"'
	Department		string 'json:"department"'
}

type PatientStatus struct {
	StatusDate  	string 'json:"statusDate"'
	StatusTime		string 'json:"statusTime"'
	Status			string 'json:"status"'
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(PatientChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *PatientChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// ==========================================================
// Initialize Patient
// ==========================================================
func (t *PatientChaincode) initPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	
	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 8 (Patient Id, Name, Gender, Age, Physician Name, Status, StatusDate, StatusTime)")
	}

	// ==== Input sanitation ====
	fmt.Println("- start initializing patient details")
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
	if len(args[4]) <= 0 {
		return shim.Error("5th argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return shim.Error("6th argument must be a non-empty string")
	}
	if len(args[6]) <= 0 {
		return shim.Error("7th argument must be a non-empty string")
	}
	if len(args[7]) <= 0 {
		return shim.Error("8th argument must be a non-empty string")
	}
	
	patientId := args[0]
	patientName := args[1]
	gender := args[2]
	age := args[3]
	physicianName := args[4]
	
	status := args[5]
	statusDate := args[6]
	statusTime := args[7]
	
	// ==== Check if patient already exists ====
	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error("Failed to get patient: " + err.Error())
	} else if patientAsBytes != nil {
		fmt.Println("This patient already exists: " + patientId)
		return shim.Error("This patient already exists: " + patientId)
	}
	
	//patientStatus := &PatientStatus{statusDate, statusTime, status}
	var patientStatus PatientStatus
	patientStatus.StatusDate = statusDate
	patientStatus.StatusTime = statusTime
	patientStatus.Status = status
	
	// ==== Create patient object and marshal to JSON ====
	patient := &Patient{patientId, patientName, gender, age, physicianName, patientStatus}
	patientJSONasBytes, err := json.Marshal(patient)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	// === Save patient to state ===
	err = stub.PutState(patientId, patientJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	// ==== Patient saved. Return success ====
	fmt.Println("- end initializing patient details")
	return shim.Success(nil)
}

// ===========================================================
// process a patient by setting with new details
// ===========================================================
func (t *PatientChaincode) updatePatient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4 (Patient Id, Name, Gender, Age)")
	}
	
	patientId := args[0]
	patientName := args[1]
	gender := args[2]
	age := args[3]
	
	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error("Failed to get patient: " + err.Error())
	} else if patientAsBytes == nil {
		return shim.Error("Patient doesn't exists")
	}
	
	patientToProcess := Patient{}
	err = json.Unmarshal(patientAsBytes, &patientToProcess) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	patientToProcess.Name = patientName 
	patientToProcess.Gender = gender
	patientToProcess.Age = age
	
	patientJSONasBytes, _ := json.Marshal(patientToProcess)
	err = stub.PutState(patientId, patientJSONasBytes) //rewrite the patient
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end of processing patient (success)")
	return shim.Success(nil)
	
}

func (t *PatientChaincode) assignPhysician(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	
	if len(args) != 6 {
		return shim.Error("Incorrect number of Doctor details. Expecting 6 (Patient Id, Physician Id, Doctor Name, Gender, Age, Department)")
	}

	// ==== Input sanitation ====
	fmt.Println("- start assigning physician")
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
	if len(args[4]) <= 0 {
		return shim.Error("5rd argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return shim.Error("6th argument must be a non-empty string")
	}
	
	patientId := args[0]
	physicianId := args[1]
	name := args[2]
	gender := args[3]
	age := args[4]
	department := args[5]
	
	commonId := patientId + "-" + physicianId;
	
	// ==== Check if doctor already assigned ====
	physicianAsBytes, err := stub.GetState(commonId)
	if err != nil {
		return shim.Error("Failed to get assigned doctor: " + err.Error())
	} else if physicianAsBytes != nil {
		fmt.Println("This physician already assigned: " + commonId)
		return shim.Error("This physician already assigned: " + commonId)
	}
	
	// ==== Create patient object and marshal to JSON ====
	physician := &PhysicianRelation{patientId, physicianId, name, gender, age, department}
	physicianJSONasBytes, err := json.Marshal(physician)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	// === Save physician to state ===
	err = stub.PutState(commonId, physicianJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	// ==== Patient saved. Return success ====
	fmt.Println("- end assigning physician")
	return shim.Success(nil)
}

func (t *PatientChaincode) addProcessDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
}

func (t *PatientChaincode) updateStatus(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting patient id, status, current date and current time")
	}
		
	patientId := args[0]
	status := args[1]
	statusDate := args[2]
	statusTime := args[3]
	
	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error("Failed to get patient: " + err.Error())
	} else if patientAsBytes == nil {
		return shim.Error("Patient doesn't exists")
	}
	
	patientToProcess := Patient{}
	err = json.Unmarshal(patientAsBytes, &patientToProcess) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	
	var patientStatus PatientStatus
	patientStatus.StatusDate = statusDate
	patientStatus.StatusTime = statusTime
	patientStatus.Status = status
	
	patientToProcess.Status = patientStatus
		
	patientJSONasBytes, _ := json.Marshal(patientToProcess)
	err = stub.PutState(patientId, patientJSONasBytes) //rewrite the patient
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("End of updating patient status (success)")
	return shim.Success(nil)
}

// ===================================================
// readPatient - read a patient from chaincode state
// ===================================================
func (t *PatientChaincode) readPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var patientId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting patient id to query")
	}

	patientId = args[0]
	valAsbytes, err := stub.GetState(patientId) //get the patient from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for patient " + patientId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Patient does not exist: " + patientId + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

func (t *PatientChaincode) queryPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	queryString := args[0]

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func (t *PatientChaincode) getHistoryForPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	patientId := args[0]

	fmt.Printf("- start getHistoryForPatient: %s\n", patientId)

	resultsIterator, err := stub.GetHistoryForKey(patientId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the patient
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON patient)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForPatient returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

