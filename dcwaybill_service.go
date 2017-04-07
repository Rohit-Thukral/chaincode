/*****Chaincode to perform Delivery Centre Way Bill*****
Methods Involved
CreateDCWayBill : Create Delivery Centre Way Bill

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

/************** Create DC Way Bill Starts ***********************/
func CreateDCWayBill(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Create WayBill", args[0])

	dcwayBillRequest := parseShipmentWayBillRequest(args[0])

	dcshipmentDetails, err := fetchShipmentWayBillData(stub, dcwayBillRequest.ShipmentNumber)

	if err != nil {
		fmt.Println("Error while retrieveing the Shipment Details", err)
		return nil, err
	}
	dcshipmentDetails.WayBillNumber = dcwayBillRequest.WayBillNumber
	dcshipmentDetails.VehicleNumber = dcwayBillRequest.VehicleNumber
	dcshipmentDetails.VehicleType = dcwayBillRequest.VehicleType
	dcshipmentDetails.PickupDate = dcwayBillRequest.PickupDate
	dcshipmentDetails.Custodian = dcwayBillRequest.Custodian
	dcshipmentDetails.TpComments = dcwayBillRequest.TpComments
	dcshipmentDetails.WayBillCreationDate = dcwayBillRequest.WayBillCreationDate
	dcshipmentDetails.WayBillCreatedBy = dcwayBillRequest.WayBillCreatedBy
	dcshipmentDetails.VehicleType = dcwayBillRequest.VehicleType
	dcshipmentDetails.EntityName = dcwayBillRequest.EntityName

	UpdatePalletCartonAssetByWayBill(stub, dcwayBillRequest, DCWAYBILL, "")
	UpdateEntityWayBillMapping(stub, dcshipmentDetails.EntityName, dcshipmentDetails.WayBillNumber)
	err = stub.PutState(dcshipmentDetails.WayBillNumber, []byte(dcwayBillRequest.ShipmentNumber))
	if err != nil {
		fmt.Println("Could not save WayBill to ledger", err)
		return nil, err
	}
	return saveShipmentWayBill(stub, dcshipmentDetails)
}

/************** Create DC Way Bill Ends ************************/
