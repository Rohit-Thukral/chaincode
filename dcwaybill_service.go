/*****Chaincode to perform Delivery Centre Way Bill*****
Methods Involved
CreateDCWayBill : Create Delivery Centre Way Bill

Author: Mohd Arshad
Dated: 30/7/2017
/*****************************************************/
package main

import (
	"encoding/json"
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
	dcshipmentDetails.Status = dcwayBillRequest.Status
	dcshipmentDetails.CustodianHistory = UpdateShipmentCustodianHistoryList(stub, dcshipmentDetails)

	_, cartonsSerialNumber, assetsSerialNumber, _ := UpdatePalletCartonAssetByWayBill(stub, dcwayBillRequest, DCWAYBILL, "")
	dcshipmentDetails.CartonsSerialNumber = cartonsSerialNumber
	dcshipmentDetails.AssetsSerialNumber = assetsSerialNumber

	UpdateEntityWayBillMapping(stub, dcshipmentDetails.EntityName, dcshipmentDetails.WayBillNumber, dcshipmentDetails.CountryFrom)
	err = DumpData(stub, dcshipmentDetails.WayBillNumber, dcwayBillRequest.ShipmentNumber)
	if err != nil {
		fmt.Println("Could not save WayBill to ledger", err)
		return nil, err
	}
	saveResult, errMsg := saveShipmentWayBill(stub, dcshipmentDetails)

	fmt.Println("Start of Transaction Details Store Methods............")
	saveResultRes := BlockchainResponse{}
	json.Unmarshal([]byte(saveResult), &saveResultRes)

	transactionDet := TransactionDetails{}
	transactionDet.TransactionId = saveResultRes.TxID
	transactionDet.TransactionTime = dcshipmentDetails.WayBillCreationDate
	if errMsg != nil {
		transactionDet.Status = "Failure"
	} else {
		transactionDet.Status = "Success"
	}
	transactionDet.FromUserId = dcshipmentDetails.Custodian
	transactionDet.ToUserId = append(transactionDet.ToUserId, dcshipmentDetails.Consignee)
	transactionDet.ToUserId = append(transactionDet.ToUserId, dcshipmentDetails.EntityName)
	fmt.Println("Start of Transaction Details Store Methods transactionDet.TransactionId............", transactionDet.TransactionId)
	fmt.Println("Start of Transaction Details Store Methods transactionDet.status............", transactionDet.Status)
	fmt.Println("Start of Transaction Details Store Methods transactionDet.FromUserId............", transactionDet.FromUserId)
	fmt.Println("Start of Transaction Details Store Methods transactionDet.ToUserId............", transactionDet.ToUserId)
	fmt.Println("Start of Transaction Details Store Methods transactionDet.TransactionTime............", transactionDet.TransactionTime)
	_ = saveTransactionDetails(stub, transactionDet)
	fmt.Println("End of Transaction Details Store Methods............")

	var waybillShipmentMapping WayBillShipmentMapping
	waybillShipmentMapping.DCShipmentNumber = dcshipmentDetails.ShipmentNumber
	waybillShipmentMapping.DCWayBillsNumber = dcshipmentDetails.WayBillNumber
	saveWayBillShipmentMapping(stub, waybillShipmentMapping)
	fmt.Println("Successfully saved waybill shipment mapping details")
	return saveResult, errMsg

}

/************** Create DC Way Bill Ends ************************/

/************** Update Waybill Starts *********************/
func UpdateDCWayBill(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Update DCWaybill", args[0])
	shipmentRequest := parseShipmentWayBillRequest(args[0])
	wayBilldata, _ := fetchShipmentWayBillData(stub, shipmentRequest.ShipmentNumber)
	shipmentRequest.CustodianHistory = UpdateShipmentCustodianHistoryList(stub, shipmentRequest)
	shipmentRequest.SupportiveDocuments = wayBilldata.SupportiveDocuments
	saveResult, errMsg := saveShipmentWayBill(stub, shipmentRequest)

	fmt.Println("Start of Transaction Details Store Methods............")
	saveResultRes := BlockchainResponse{}
	json.Unmarshal([]byte(saveResult), &saveResultRes)

	transactionDet := TransactionDetails{}
	transactionDet.TransactionId = saveResultRes.TxID
	transactionDet.TransactionTime = shipmentRequest.WayBillModifiedDate
	if errMsg != nil {
		transactionDet.Status = "Failure"
	} else {
		transactionDet.Status = "Success"
	}
	transactionDet.FromUserId = shipmentRequest.Custodian
	transactionDet.ToUserId = append(transactionDet.ToUserId, shipmentRequest.Consignee)
	transactionDet.ToUserId = append(transactionDet.ToUserId, shipmentRequest.EntityName)
	fmt.Println("Start of Transaction Details Store Methods transactionDet.TransactionId............", transactionDet.TransactionId)
	fmt.Println("Start of Transaction Details Store Methods transactionDet.status............", transactionDet.Status)
	fmt.Println("Start of Transaction Details Store Methods transactionDet.FromUserId............", transactionDet.FromUserId)
	fmt.Println("Start of Transaction Details Store Methods transactionDet.ToUserId............", transactionDet.ToUserId)
	fmt.Println("Start of Transaction Details Store Methods transactionDet.TransactionTime............", transactionDet.TransactionTime)
	_ = saveTransactionDetails(stub, transactionDet)
	fmt.Println("End of Transaction Details Store Methods............")
	return saveResult, errMsg
}

/************** Update waybill Ends ************************/

/************** Get waybill WayBill Starts ******************/
/*This is common code for Get Shipment,WayBill,DCShipment,DCWayBill*/

func ViewDCWayBill(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering ViewDCWayBill " + args[0])

	waybillno := args[0]
	wayBillShipmentMapping, _ := FetchWayBillShipmentMappingData(stub, waybillno)
	shipmentNo := wayBillShipmentMapping.DCShipmentNumber
	wayBilldata, dataerr := fetchShipmentWayBillData(stub, shipmentNo)
	if dataerr == nil {

		dataToStore, _ := json.Marshal(wayBilldata)
		return []byte(dataToStore), nil

	}

	return nil, dataerr

}

/************** Get  WayBill Ends ********************/
