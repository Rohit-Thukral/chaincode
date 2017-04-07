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
	return saveShipmentWayBill(stub, shipmentRequest)
}

/************** Create DC Shipment Ends ************************/
