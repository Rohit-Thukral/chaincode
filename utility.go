package main


import (
	"encoding/json"
	"fmt"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Utility struct {

}


func (t *Utility) fetchShipmentData(stub shim.ChaincodeStubInterface, shipmentNumber string) (ShipmentWayBill, error) {
	fmt.Println("Entering fetchShipmentData " + shipmentNumber)
	var shipmentData ShipmentWayBill

	indexByte, err := stub.GetState(shipmentNumber)
	if err != nil {
		fmt.Println("fetchShipmentData :: Could not retrive Shipment Index", err)
		return shipmentData, err
	}

	if marshErr := json.Unmarshal(indexByte, &shipmentData); marshErr != nil {
		fmt.Println("fetchShipmentData :: Could not Unmarshal the data", marshErr)
		return shipmentData, marshErr
	}
	fmt.Println("======================")
	fmt.Println(shipmentData)
	fmt.Println("======================")
	fmt.Println("Exiting fetchShipmentData ")
	return shipmentData, nil

}

func (t *Utility) fetchShipmentIndex(stub shim.ChaincodeStubInterface, callingEntityName string, status string) ([]ShipmentWayBill, error) {
	fmt.Println("Entering fetchShipmentIndex  callingEntityName : " + callingEntityName + "   ----  status : "+status)
	allShipmentIndex := ShipmentWayBillIndex{}
	var shipmentIndexArr []string
	var tmpShipmentIndex string
	
	var shipmentDataArray []ShipmentWayBill

	var util Utility
	
	indexByte, err := stub.GetState("ShipmentWayBillIndex")
	if err != nil {
		fmt.Println("Could not retrive Shipment Index", err)
		return nil, err
	}

	fmt.Println("indexByte : ======================")
	fmt.Println(indexByte)
	fmt.Println("indexByte : ======================")

	if marshErr := json.Unmarshal(indexByte, &allShipmentIndex); marshErr != nil {
		fmt.Println("Could not Parse ShipmentWayBillIndex", marshErr)
		return nil, marshErr
	}

	shipmentIndexArr = allShipmentIndex.ShipmentNumber

	lenOfArray := len(shipmentIndexArr)

	for i := 0; i < lenOfArray; i++ {
		var shipmentData ShipmentWayBill
		tmpShipmentIndex = shipmentIndexArr[i]
		shipmentData,err = util.fetchShipmentData(stub, tmpShipmentIndex);
		if(err == nil && ((shipmentData.Custodian == callingEntityName && status != "ALL") || (util.hasPermission(shipmentData.CustodianHistory , callingEntityName) && status == "ALL"))) {
			shipmentDataArray = append(shipmentDataArray, shipmentData);
		}
		
	}
	fmt.Println("shipmentDataArray : ======================")
	fmt.Println(shipmentDataArray)
	fmt.Println("shipmentDataArray : ======================")
	fmt.Println("Exiting fetchShipmentIndex ")
	return shipmentDataArray, nil
}

func (t *Utility) hasPermission(acl []string, currUser string) bool {
	lenOfArray := len(acl)

	for i := 0; i < lenOfArray; i++ {
		if acl[i] == currUser {
			return true
		}
	}

	return false
}