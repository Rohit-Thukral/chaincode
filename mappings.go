/*****Chaincode to perform Mappings*****
Methods Involved
CreateEntityWayBillMapping : Mapping for Entity Name and Waybill
UpdateEntityWayBillMapping: Update Mapping for Entity Name and Waybill

Author: Mohd Arshad
Dated: 30/7/2017
/*****************************************************/
package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

/************** Create Entity WayBill Mapping Starts ************************/
func CreateEntityWayBillMapping(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Create Entity WayBill Mapping")
	entityWayBillMappingRequest := parseEntityWayBillMapping(args[0])

	return saveEntityWayBillMapping(stub, entityWayBillMappingRequest)

}
func parseEntityWayBillMapping(jsondata string) CreateEntityWayBillMappingRequest {
	res := CreateEntityWayBillMappingRequest{}
	json.Unmarshal([]byte(jsondata), &res)
	fmt.Println(res)
	return res
}
func saveEntityWayBillMapping(stub shim.ChaincodeStubInterface, createEntityWayBillMappingRequest CreateEntityWayBillMappingRequest) ([]byte, error) {

	entityWayBillMapping := EntityWayBillMapping{}
	entityWayBillMapping.WayBillsNumber = createEntityWayBillMappingRequest.WayBillsNumber

	dataToStore, _ := json.Marshal(entityWayBillMapping)

	err := stub.PutState(createEntityWayBillMappingRequest.EntityName, []byte(dataToStore))
	if err != nil {
		fmt.Println("Could not save Entity WayBill Mapping to ledger", err)
		return nil, err
	}

	resp := BlockchainResponse{}
	resp.Err = "000"
	resp.Message = createEntityWayBillMappingRequest.EntityName

	respString, _ := json.Marshal(resp)

	fmt.Println("Successfully saved Entity WayBill Mapping")
	return []byte(respString), nil

}

/************** Create Entity WayBill Mapping Ends ************************/

/************** Update Entity WayBill Mapping Starts ************************/
func UpdateEntityWayBillMapping(stub shim.ChaincodeStubInterface, entityName string, wayBillsNumber string) ([]byte, error) {
	fmt.Println("Entering Update Entity WayBill Mapping")
	entityWayBillMappingRequest := CreateEntityWayBillMappingRequest{}
	entityWayBillMapping, err := fetchEntityWayBillMappingData(stub, entityName)

	if err != nil {
		entityWayBillMappingRequest.EntityName = entityName
		entityWayBillMappingRequest.WayBillsNumber = append(entityWayBillMappingRequest.WayBillsNumber, wayBillsNumber)
		saveEntityWayBillMapping(stub, entityWayBillMappingRequest)
	} else {
		entityWayBillMappingRequest.WayBillsNumber = append(entityWayBillMapping.WayBillsNumber, wayBillsNumber)
		fmt.Println("Updated Entity", entityWayBillMappingRequest)
		dataToStore, _ := json.Marshal(entityWayBillMappingRequest)
		err := stub.PutState(entityName, []byte(dataToStore))
		if err != nil {
			fmt.Println("Could not save Entity WayBill Mapping to ledger", err)
			return nil, err
		}
	}
	resp := BlockchainResponse{}
	resp.Err = "000"
	resp.Message = entityName

	respString, _ := json.Marshal(resp)

	fmt.Println("Successfully saved Entity WayBill Mapping")
	return []byte(respString), nil

}

/************** Get Entity WayBill Mapping Starts ************************/
func GetEntityWayBillMapping(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Get Entity WayBill Mapping")
	entityName := args[0]
	wayBillEntityMappingData, dataerr := fetchEntityWayBillMappingData(stub, entityName)
	if dataerr == nil {

		dataToStore, _ := json.Marshal(wayBillEntityMappingData)
		return []byte(dataToStore), nil

	}

	return nil, dataerr
}

func fetchEntityWayBillMappingData(stub shim.ChaincodeStubInterface, entityName string) (EntityWayBillMapping, error) {
	var entityWayBillMapping EntityWayBillMapping

	indexByte, err := stub.GetState(entityName)
	if err != nil {
		fmt.Println("Could not retrive Entity WayBill Mapping ", err)
		return entityWayBillMapping, err
	}

	if marshErr := json.Unmarshal(indexByte, &entityWayBillMapping); marshErr != nil {
		fmt.Println("Could not retrieve Entity WayBill Mapping from ledger", marshErr)
		return entityWayBillMapping, marshErr
	}

	return entityWayBillMapping, nil

}

/************** Get Entity WayBill Mapping Ends ************************/

/************** Create WayBill Shipment Mapping Starts ************************/

func CreateWayBillShipmentMapping(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Create WayBill Shipment Mapping")
	wayBillShipmentMapping := parseWayBillShipmentMapping(args[0])

	return saveWayBillShipmentMapping(stub, wayBillShipmentMapping)

}
func parseWayBillShipmentMapping(jsondata string) WayBillShipmentMapping {
	res := WayBillShipmentMapping{}
	json.Unmarshal([]byte(jsondata), &res)
	fmt.Println(res)
	return res
}
func saveWayBillShipmentMapping(stub shim.ChaincodeStubInterface, craeteWayBillShipmentMappingRequest WayBillShipmentMapping) ([]byte, error) {

	wayBillShipmentMapping := WayBillShipmentMapping{}
	wayBillShipmentMapping.DCWayBillsNumber = craeteWayBillShipmentMappingRequest.DCWayBillsNumber
	wayBillShipmentMapping.DCShipmentNumber = craeteWayBillShipmentMappingRequest.DCShipmentNumber
	dataToStore, _ := json.Marshal(wayBillShipmentMapping)

	err := stub.PutState(wayBillShipmentMapping.DCWayBillsNumber, []byte(dataToStore))
	if err != nil {
		fmt.Println("Could not save WayBill Shipment Mapping to ledger", err)
		return nil, err
	}

	resp := BlockchainResponse{}
	resp.Err = "000"
	resp.Message = wayBillShipmentMapping.DCWayBillsNumber

	respString, _ := json.Marshal(resp)

	fmt.Println("Successfully saved Entity WayBill Mapping")
	return []byte(respString), nil

}

/************** Create WayBill Shipment Ends ************************/

/************** Get  WayBill Shipment Mapping Starts ************************/
func GetWayBillShipmentMapping(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Get Entity WayBill Mapping")
	wayBillNumber := args[0]
	wayBillShippingMappingData, dataerr := fetchEntityWayBillMappingData(stub, wayBillNumber)
	if dataerr == nil {

		dataToStore, _ := json.Marshal(wayBillShippingMappingData)
		return []byte(dataToStore), nil

	}

	return nil, dataerr
}

func fetchWayBillShipmentMappingData(stub shim.ChaincodeStubInterface, wayBillNumber string) (WayBillShipmentMapping, error) {
	var wayBillShipmentMapping WayBillShipmentMapping

	indexByte, err := stub.GetState(wayBillNumber)
	if err != nil {
		fmt.Println("Could not retrive WayBill Shipping Mapping ", err)
		return wayBillShipmentMapping, err
	}

	if marshErr := json.Unmarshal(indexByte, &wayBillShipmentMapping); marshErr != nil {
		fmt.Println("Could not retrieve Entity WayBill Mapping from ledger", marshErr)
		return wayBillShipmentMapping, marshErr
	}

	return wayBillShipmentMapping, nil

}

/************** Get  WayBill Shipment Mapping Ends ************************/
