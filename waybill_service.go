/*****Chaincode to perform Waybill realeted task*****
Methods Involved
CreateWayBill : Used for Creating Waybill

Author: Mohd Arshad
Dated: 30/7/2017
/*****************************************************/

package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

/************** Create Way Bill Starts ***********************/
func CreateWayBill(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Create WayBill", args[0])

	wayBillRequest := parseShipmentWayBillRequest(args[0])

	shipmentDetails, err := fetchShipmentWayBillData(stub, wayBillRequest.ShipmentNumber)
	if err != nil {
		fmt.Println("Error while retrieveing the Shipment Details", err)
		return nil, err
	}
	shipmentDetails.WayBillNumber = wayBillRequest.WayBillNumber
	shipmentDetails.VehicleNumber = wayBillRequest.VehicleNumber
	shipmentDetails.VehicleType = wayBillRequest.VehicleType
	shipmentDetails.PickupDate = wayBillRequest.PickupDate
	shipmentDetails.Custodian = wayBillRequest.Custodian
	shipmentDetails.TpComments = wayBillRequest.TpComments
	shipmentDetails.WayBillCreationDate = wayBillRequest.WayBillCreationDate
	shipmentDetails.WayBillCreatedBy = wayBillRequest.WayBillCreatedBy
	shipmentDetails.CustodianHistory = UpdateShipmentCustodianHistoryList(stub, shipmentDetails)

	UpdatePalletCartonAssetByWayBill(stub, wayBillRequest, WAYBILL, "")
	
	return saveShipmentWayBill(stub, shipmentDetails)
}

/************** Create Way Bill Ends ************************/
