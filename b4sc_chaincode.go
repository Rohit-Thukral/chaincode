/*****Main Chaicode to start the execution*****

/*****************************************************/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const NODATA_ERROR_CODE string = "400"

const NODATA_ERROR_MSG string = "No data found"

const INVALID_INPUT_ERROR_CODE string = "401"
const INVALID_INPUT_ERROR_MSG string = "Invalid Input"

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

type B4SCChaincode struct {
}

//custom data models

type Pallet struct {
	PalletId    string
	Modeltype   string
	CartonId    []string
	ShipmentIds []string
}

type Carton struct {
	CartonId    string
	PalletId    string
	AssetId     []string
	ShipmentIds []string
}

type Asset struct {
	AssetId     string
	Modeltype   string
	Color       string
	CartonId    string
	PalletId    string
	ShipmentIds []string
}

type WayBillHistory struct {
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Status    string  `json:"status"`
	Timestamp string  `json:"timestamp"`
	Notes     string  `json:"notes"`
	Lat       float64 `json:"lat"`
	Log       float64 `json:"log"`
}



type ShipmentIndex struct {
	ShipmentNumber string
	Status         string
	Acl            []string
}

type AllShipment struct {
	ShipmentIndexArr []ShipmentIndex
}

type AllShipmentDump struct {
	ShipmentIndexArr []string `json:"shipmentIndexArr"`
}

type Entity struct {
	EntityId        string  `json:"entityId"`
	EntityName      string  `json:"entityName"`
	EntityType      string  `json:"entityType"`
	EntityAddress   string  `json:"entityAddress"`
	EntityRegNumber string  `json:"entityRegNumber"`
	EntityCountry   string  `json:"entityCountry"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
}

//Will be avlable in the WorldStats as "ALL_ENTITIES"
type AllEntities struct {
	EntityArr []string `json:"entityArr"`
}

//Will be avlable in the WorldStats as "ASSET_MODEL_NAMES"
type AssetModelDetails struct {
	ModelNames []string `json:"modelNames"`
}

type WorkflowDetails struct {
	FromEntity  string   `json:"fromEntity"`
	ToEntity    string   `json:"toEntity"`
	Carrier     string   `json:"carrier"`
	EntityOrder []string `json:"entityOrder"`
}

//Will be available in the WorldStats as "ALL_WORKFLOWS"
type AllWorkflows struct {
	Workflows []WorkflowDetails `json:"workflows"`
}

/************** Arshad Start Code This new struct for AssetDetails , CartonDetails , PalletDetails  is added by Arshad as to incorporate new LLD published orginal structure
are not touched as of now to avoid break of any functionality devloped by Kartik 20/3/2017***************/

//Will be avlable in the WorldStats as "ShipmentWayBillIndex"
type ShipmentWayBillIndex struct {
	ShipmentNumber []string `json:"shipmentNumber"`
}

//Will be avlable in the WorldStats as "WayBillNumberIndex"
type WayBillNumberIndex struct {
	WayBillNumber []string
}

type EntityDetails struct {
	EntityName      string
	EntityType      string
	EntityAddress   string
	EntityRegNumber string
	EntityCountry   string
	Latitude        string
	Longitude       string
}

/**************Arshad End new code as per LLD***************
//START COMMNENTED BY ARSHAD AS THESE ARE NO MORE USED - KARTHIK PLEASE REVIEW AND PERMANNETLY DELETE IT IF NOT USED BY ANY OF YOU MODULE


/************** Create Shipment Starts ************************
/**
	Expected Input is
	{
		"shipmentNumber"" : "123456",
		"personConsigningGoods" : "KarthikS",
		"consigner" : "HCL",
		"consignerAddress" : "Chennai",
		"consignee" : "HCL-AM",
		"consigneeAddress" : "Dallas",
		"consigneeRegNo" : "12122222222",
		"ModelNo" : "IA1a1222",
		"quantity" : "50",
		"pallets" : ["11111111","22222222","333333"],
		"status" : "intra",
		"notes" : "ha haha ha ha",
		"CreatedBy" : "KarthikSukumaram",
		"custodian" : "HCL",
		"createdTimeStamp" : "2017-03-02"
	}
**

type CreateShipmentRequest struct {
	ShipmentNumber        string   `json:"shipmentNumber"`
	PersonConsigningGoods string   `json:"personConsigningGoods"`
	Consigner             string   `json:"consigner"`
	ConsignerAddress      string   `json:"consignerAddress"`
	Consignee             string   `json:"consignee"`
	ConsigneeAddress      string   `json:"consigneeAddress"`
	ConsigneeRegNo        string   `json:"consigneeRegNo"`
	ModelNo               string   `json:"modelNo"`
	Quantity              string   `json:"quantity"`
	Pallets               []string `json:"pallets"`
	Carrier               string   `json:"status"`
	Notes                 string   `json:"notes"`
	CreatedBy             string   `json:"createdBy"`
	Custodian             string   `json:"custodian"`
	CreatedTimeStamp      string   `json:"createdTimeStamp"`
	CallingEntityName     string   `json:"callingEntityName"`
}

type CreateShipmentResponse struct {
	Err     string `json:"err"`
	ErrMsg  string `json:"errMsg"`
	Message string `json:"message"`
}

/*
func CreateShipment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering CreateShipment")

	shipmentRequest := parseCreateShipmentRequest(args[0])

	return processShipment(stub, shipmentRequest)

}

func processShipment(stub shim.ChaincodeStubInterface, shipmentRequest CreateShipmentRequest) ([]byte, error) {
	shipment := Shipment{}
	shipmentIndex := ShipmentIndex{}

	shipment.ShipmentNumber = shipmentRequest.ShipmentNumber
	shipment.PersonConsigningGoods = shipmentRequest.PersonConsigningGoods
	shipment.Consigner = shipmentRequest.Consigner
	shipment.ConsignerAddress = shipmentRequest.ConsignerAddress
	shipment.Consignee = shipmentRequest.Consignee
	shipment.ConsigneeAddress = shipmentRequest.ConsigneeAddress
	shipment.ConsigneeRegNo = shipmentRequest.ConsigneeRegNo
	shipment.ModelNo = shipmentRequest.ModelNo
	shipment.Quantity = shipmentRequest.Quantity
	shipment.Pallets = shipmentRequest.Pallets
	shipment.Carrier = shipmentRequest.Carrier
	shipment.CreatedBy = shipmentRequest.CreatedBy
	shipment.Custodian = shipmentRequest.Custodian
	shipment.CreatedTimeStamp = shipmentRequest.CreatedTimeStamp
	shipment.Status = "Created"

	var acl []string
	acl = append(acl, shipmentRequest.CallingEntityName) //TODO: Have to take the Entity name from the Certificate
	shipment.Acl = acl

	shipmentIndex.ShipmentNumber = shipmentRequest.ShipmentNumber
	shipmentIndex.Status = shipment.Status
	shipmentIndex.Acl = acl

	dataToStore, _ := json.Marshal(shipment)

	err := stub.PutState(shipment.ShipmentNumber, []byte(dataToStore))
	if err != nil {
		fmt.Println("Could not save Shipment to ledger", err)
		return nil, err
	}

	addShipmentIndex(stub, shipmentIndex)

	resp := CreateShipmentResponse{}
	resp.Err = "000"
	resp.Message = shipment.ShipmentNumber

	respString, _ := json.Marshal(resp)

	fmt.Println("Successfully saved way bill")
	return []byte(respString), nil

}

func addShipmentIndex(stub shim.ChaincodeStubInterface, shipmentIndex ShipmentIndex) error {
	indexByte, err := stub.GetState("SHIPMENT_INDEX")
	if err != nil {
		fmt.Println("Could not retrive Shipment Index", err)
		return err
	}
	allShipmentIndex := AllShipment{}

	if marshErr := json.Unmarshal(indexByte, &allShipmentIndex); marshErr != nil {
		fmt.Println("Could not save Shipment to ledger", marshErr)
		return marshErr
	}

	allShipmentIndex.ShipmentIndexArr = append(allShipmentIndex.ShipmentIndexArr, shipmentIndex)
	dataToStore, _ := json.Marshal(allShipmentIndex)

	addErr := stub.PutState("SHIPMENT_INDEX", []byte(dataToStore))
	if addErr != nil {
		fmt.Println("Could not save Shipment to ledger", addErr)
		return addErr
	}

	return nil
}

func parseCreateShipmentRequest(jsondata string) CreateShipmentRequest {
	res := CreateShipmentRequest{}
	json.Unmarshal([]byte(jsondata), &res)
	fmt.Println(res)
	return res
}

/************** Create Shipment Ends ************************

/************** View Shipment Starts ************************

type ViewShipmentRequest struct {
	CallingEntityName string `json:"callingEntityName"`
	ShipmentNumber    string `json:"shipmentNumber"`
}

func ViewShipment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering ViewShipment " + args[0])

	request := parseViewShipmentRequest(args[0])

	shipmentData, dataerr := fetchShipmentData(stub, request.ShipmentNumber)
	if dataerr == nil {
		if hasPermission(shipmentData.Acl, request.CallingEntityName) {
			dataToStore, _ := json.Marshal(shipmentData)
			return []byte(dataToStore), nil
		} else {
			return []byte("{ \"errMsg\": \"No data found\" }"), nil
		}
	}

	return nil, nil

}

func parseViewShipmentRequest(jsondata string) ViewShipmentRequest {
	res := ViewShipmentRequest{}
	json.Unmarshal([]byte(jsondata), &res)
	fmt.Println(res)
	return res
}

/************** View Shipment Ends ************************/

/************** Inbox Service Starts ************************/

/**
	Expected Input is
	{
		"callingEntityName" : "INTEL",
		"status" : "Created"
	}
**/

//END COMMNENTED BY ARSHAD AS THESE ARE NO MORE USED - KARTHIK PLEASE REVIEW AND PERMANNETLY DELETE IT IF NOT USED BY ANY OF YOU MODULE

/*type InboxRequest struct {
	CallingEntityName string `json:"callingEntityName"`
	Status            string `json:"status"`
}

type InboxResponse struct {
	Data []Shipment `json:"data"`
}

func Inbox(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Inbox " + args[0])

	request := parseInboxRequest(args[0])

	return fetchShipmentIndex(stub, request.CallingEntityName, request.Status)

}

func parseInboxRequest(jsondata string) InboxRequest {
	res := InboxRequest{}
	json.Unmarshal([]byte(jsondata), &res)
	fmt.Println(res)
	return res
}

func hasPermission(acl []string, currUser string) bool {
	lenOfArray := len(acl)

	for i := 0; i < lenOfArray; i++ {
		if acl[i] == currUser {
			return true
		}
	}

	return false
}

func fetchShipmentData(stub shim.ChaincodeStubInterface, shipmentNumber string) (Shipment, error) {
	var shipmentData Shipment

	indexByte, err := stub.GetState(shipmentNumber)
	if err != nil {
		fmt.Println("Could not retrive Shipment Index", err)
		return shipmentData, err
	}

	if marshErr := json.Unmarshal(indexByte, &shipmentData); marshErr != nil {
		fmt.Println("Could not save Shipment to ledger", marshErr)
		return shipmentData, marshErr
	}

	return shipmentData, nil

}

func fetchShipmentIndex(stub shim.ChaincodeStubInterface, callingEntityName string, status string) ([]byte, error) {
	allShipmentIndex := AllShipment{}
	var shipmentIndexArr []ShipmentIndex
	var tmpShipmentIndex ShipmentIndex
	var shipmentDataArr []Shipment
	resp := InboxResponse{}

	indexByte, err := stub.GetState("SHIPMENT_INDEX")
	if err != nil {
		fmt.Println("Could not retrive Shipment Index", err)
		return nil, err
	}

	if marshErr := json.Unmarshal(indexByte, &allShipmentIndex); marshErr != nil {
		fmt.Println("Could not save Shipment to ledger", marshErr)
		return nil, marshErr
	}

	shipmentIndexArr = allShipmentIndex.ShipmentIndexArr

	lenOfArray := len(shipmentIndexArr)

	for i := 0; i < lenOfArray; i++ {
		tmpShipmentIndex = shipmentIndexArr[i]
		if tmpShipmentIndex.Status == status {
			if hasPermission(tmpShipmentIndex.Acl, callingEntityName) {
				shipmentData, dataerr := fetchShipmentData(stub, tmpShipmentIndex.ShipmentNumber)
				if dataerr == nil {
					shipmentDataArr = append(shipmentDataArr, shipmentData)
				}
			}
		}
	}

	resp.Data = shipmentDataArr
	dataToStore, _ := json.Marshal(resp)

	return []byte(dataToStore), nil
}*/

/************** Inbox Service Ends ************************/

//START COMMNENTED BY ARSHAD AS THESE ARE NO MORE USED - KARTHIK PLEASE REVIEW AND PERMANNETLY DELETE IT IF NOT USED BY ANY OF YOU MODULE



/************** View Data for Key Starts ************************/

func ViewDataForKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering ViewDataForKey " + args[0])

	return stub.GetState(args[0])

}

/************** View Data for Key Ends ************************/

/************** DumpData Start ************************/

type DumpDataKeysType struct {
	Keys 	[]string 	`json:"keys"`
}


func parseDumpDataKeysType(stub shim.ChaincodeStubInterface) DumpDataKeysType {
	res := DumpDataKeysType{}
	indexByte, err := stub.GetState("DumpDataKeysType")
	err = json.Unmarshal(indexByte, &res)
	if err != nil {
		fmt.Println("Could not Parse the Data", err)
	}
	fmt.Println(res)
	return res
}

func storeDumpDataKeysType(stub shim.ChaincodeStubInterface, keyString string) (error) {
	ddData := parseDumpDataKeysType(stub)
	var keys []string
	keys = ddData.Keys
	keys = append(keys, keyString)
	ddData.Keys = keys
	fmt.Println(ddData)
	dataToStore, _ := json.Marshal(ddData)
	fmt.Println(string(dataToStore))
	err :=  stub.PutState("DumpDataKeysType", dataToStore)
	if err != nil {
		fmt.Println("Could not save the DumpDataKeysType", err)
		return err
	}
	return nil
}

func DumpData(stub shim.ChaincodeStubInterface, argsKey string, argsValue string) (error) {
	fmt.Println("Entering DumpData " + argsKey + "  " + argsValue)
	

	err := stub.PutState(argsKey, []byte(argsValue))
	if err != nil {
		fmt.Println("Could not save the Data", err)
		return err
	}
	storeDumpDataKeysType(stub, argsKey)

	return nil
}

/************** DumpData Ends ************************/

func Initialize(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

// Init resets all the things
func (t *B4SCChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Inside INIT for test chaincode")
	ddData := DumpDataKeysType{}
	allShipment := AllShipment{}
	var tmpShipmentIndex []ShipmentIndex
	allShipment.ShipmentIndexArr = tmpShipmentIndex

	dataToStore, _ := json.Marshal(allShipment)
	dataToStore2, err2 := json.Marshal(ddData)

	stub.PutState("DumpDataKeysType", dataToStore2)

	err := stub.PutState("SHIPMENT_INDEX", dataToStore)
	if err != nil {
		fmt.Println("Could not save Shipment to ledger", err)
		return nil, err
	}

	if err2 != nil {
		fmt.Println("Could not save Shipment to ledger", err2)
	}

	return nil, nil
}

func (t *B4SCChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	/*if function == "Init" {
		return Init(stub, function, args)
	}else*/
	if function == "CreateShipment" {
		return CreateShipment(stub, args)
	} else if function == "DumpData" {
		DumpData(stub, args[0], args[1])
		return nil, nil
	} else if function == "CreateShipment" {
		return CreateShipment(stub, args)
	} else if function == "UpdateShipment" {
		return UpdateShipment(stub, args)
	} else if function == "CreateWayBill" {
		return CreateWayBill(stub, args)
	} else if function == "CreateDCShipment" {
		return CreateDCShipment(stub, args)
	} else if function == "CreateDCWayBill" {
		return CreateDCWayBill(stub, args)
	} else if function == "CreateEWWayBill" {
		return CreateEWWayBill(stub, args)
	} else if function == "UpdateEWWayBill" {
		return UpdateEWWayBill(stub, args)
	} else if function == "CreateEntityWayBillMapping" {
		return CreateEntityWayBillMapping(stub, args)
	} else if function == "CreateAsset" {
		return CreateAsset(stub, args)
	} else if function == "CreateCarton" {
		return CreateCarton(stub, args)
	} else if function == "CreatePallet" {
		return CreatePallet(stub, args)
	} else if function == "UpdateAssetDetails" {
		return UpdateAssetDetails(stub, args)
	} else if function == "UpdateCartonDetails" {
		return UpdateCartonDetails(stub, args)
	} else if function == "UpdatePalletDetails" {
		return UpdatePalletDetails(stub, args)
	} else if function == "uploadComplianceDocument" {
		return uploadComplianceDocument(stub, args)
	} else {
		return nil, errors.New("Invalid function name " + function)
	}
	//return nil, nil
}

func (t *B4SCChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Query : function : " + function + " args : " + args[0])

	if function == "ViewDataForKey" {
		return ViewDataForKey(stub, args)
	} else if function == "Inbox" {
		var inboxService InboxService
		return inboxService.Inbox(stub, args)
	} else if function == "SearchDateRange" {
		return SearchDateRange(stub, args)
	} else if function == "ShipmentPageLoad" {
		var pageLoadService ShipmentPageLoadService
		return pageLoadService.ShipmentPageLoad(stub, args)
	} else if function == "ViewEWWayBill" {
		return ViewEWWayBill(stub, args)
	} else if function == "ViewEWWayBill" {
		return ViewEWWayBill(stub, args)
	} else if function == "GetEntityWayBillMapping" {
		return GetEntityWayBillMapping(stub, args)
	} else if function == "GetAsset" {
		return GetAsset(stub, args)
	} else if function == "GetPallet" {
		return GetPallet(stub, args)
	} else if function == "GetCarton" {
		return GetCarton(stub, args)
	} else if function == "ViewShipmentWayBill" {
		return ViewShipmentWayBill(stub, args)
	} else if function == "getAllComplianceDocument" {
		return getAllComplianceDocument(stub, args)
	} else if function == "SearchDateRange" {
		return SearchDateRange(stub, args)
	} else if function == "SearchPallet" {
		return SearchPallet(stub, args)
	} else if function == "SearchCarton" {
		return SearchCarton(stub, args)
	} else if function == "SearchAsset" {
		return SearchAsset(stub, args)
	} else if function == "GetCountryWarehouse" {
		var pageLoadService ShipmentPageLoadService
		return pageLoadService.GetCountryWarehouse(stub, args)
	}
	return nil, errors.New("Invalid function name " + function)

}

func main() {
	Initialize(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	err := shim.Start(new(B4SCChaincode))
	if err != nil {
		fmt.Println("Could not start B4SCChaincode")
	} else {
		fmt.Println("B4SCChaincode successfully started")
	}
}
