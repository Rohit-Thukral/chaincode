/*****Chaincode to perform Delivery Centre Shipment*****
Methods Involved
CreateDCShipment : Create Delivery Centre Shipment

Author: Mohd Arshad
Dated: 30/7/2017
/*****************************************************/
package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

/*****Chaincode to perfor shipment realeted task*****
Methods Involved
CreateShipment : Used for Creating Shipment

Author: Mohd Arshad
Dated: 30/7/2017
/*****************************************************/

/************** Create DC Shipment Starts ***********************/
func CreateDCShipment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering DC Create Shipment", args[0])
	shipmentRequest := parseShipmentWayBillRequest(args[0])
	UpdatePalletCartonAssetByWayBill(stub, shipmentRequest, DCSHIPMENT, "")
	shipmentRequest.CustodianHistory = UpdateShipmentCustodianHistoryList(stub, shipmentRequest)
	saveResult, errMsg := saveShipmentWayBill(stub, shipmentRequest)

	/*********Storing Shipment number in shipmentwaybillindex array to retrieve through inbox*************/
	shipmentwaybillidsRequest := ShipmentWayBillIndex{}
	shipmentwaybillids, err := FetchShipmentWayBillIndex(stub, "ShipmentWayBillIndex")
	fmt.Println("shipment ids.....", shipmentwaybillids)
	if err != nil {
		shipmentwaybillidsRequest.ShipmentNumber = append(shipmentwaybillidsRequest.ShipmentNumber, shipmentRequest.ShipmentNumber)
		SaveShipmentWaybillIndex(stub, shipmentwaybillidsRequest)
	} else {
		shipmentwaybillidsRequest.ShipmentNumber = append(shipmentwaybillids.ShipmentNumber, shipmentRequest.ShipmentNumber)
		fmt.Println("Updated entity shipmentwaybillindex", shipmentwaybillidsRequest)
		SaveShipmentWaybillIndex(stub, shipmentwaybillidsRequest)
	}
	/********* End Storing Shipment number in shipmentwaybillindex array to retrieve through inbox*************/

	return saveResult, errMsg
}

/************** Create DC Shipment Ends ************************/
