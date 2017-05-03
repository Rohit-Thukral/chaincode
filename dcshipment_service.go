/*****Chaincode to perform Delivery Centre Shipment*****
Methods Involved
CreateDCShipment : Create Delivery Centre Shipment

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

/************** Create DC Shipment Starts ***********************/
func CreateDCShipment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering DC Create Shipment", args[0])
	shipmentRequest := parseShipmentWayBillRequest(args[0])
	_, cartonsSerialNumber, assetsSerialNumber, _ := UpdatePalletCartonAssetByWayBill(stub, shipmentRequest, DCSHIPMENT, "")
	shipmentRequest.CartonsSerialNumber = cartonsSerialNumber
	shipmentRequest.AssetsSerialNumber = assetsSerialNumber
	shipmentRequest.CustodianHistory = UpdateShipmentCustodianHistoryList(stub, shipmentRequest)
	saveResult, errMsg := saveShipmentWayBill(stub, shipmentRequest)

	fmt.Println("Start of Transaction Details Store Methods............")
	saveResultRes := BlockchainResponse{}
	json.Unmarshal([]byte(saveResult), &saveResultRes)

	transactionDet := TransactionDetails{}
	transactionDet.TransactionId = saveResultRes.TxID
	transactionDet.TransactionTime = shipmentRequest.ShipmentCreationDate
	if errMsg != nil {
		transactionDet.Status = "Failure"
	} else {
		transactionDet.Status = "Success"
	}
	transactionDet.FromUserId = shipmentRequest.Consigner
	transactionDet.ToUserId = append(transactionDet.ToUserId, shipmentRequest.Consignee)
	transactionDet.ToUserId = append(transactionDet.ToUserId, shipmentRequest.Carrier)
	fmt.Println("Start of Transaction Details Store Methods transactionDet.TransactionId............", transactionDet.TransactionId)
	fmt.Println("Start of Transaction Details Store Methods transactionDet.status............", transactionDet.Status)
	fmt.Println("Start of Transaction Details Store Methods transactionDet.FromUserId............", transactionDet.FromUserId)
	fmt.Println("Start of Transaction Details Store Methods transactionDet.ToUserId............", transactionDet.ToUserId)
	fmt.Println("Start of Transaction Details Store Methods transactionDet.TransactionTime............", transactionDet.TransactionTime)
	_ = saveTransactionDetails(stub, transactionDet)
	fmt.Println("End of Transaction Details Store Methods............")

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

/************** Update Waybill Starts *********************/
func UpdateDCShipment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Update Waybill", args[0])
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
	transactionDet.TransactionTime = shipmentRequest.ShipmentModifiedDate
	if errMsg != nil {
		transactionDet.Status = "Failure"
	} else {
		transactionDet.Status = "Success"
	}
	transactionDet.FromUserId = shipmentRequest.Consigner
	transactionDet.ToUserId = append(transactionDet.ToUserId, shipmentRequest.Consignee)
	transactionDet.ToUserId = append(transactionDet.ToUserId, shipmentRequest.Carrier)
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

func ViewDCShipment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering ViewDCShipment " + args[0])

	dcshipmentNo := args[0]
	wayBilldata, dataerr := fetchShipmentWayBillData(stub, dcshipmentNo)
	if dataerr == nil {

		dataToStore, _ := json.Marshal(wayBilldata)
		return []byte(dataToStore), nil

	}

	return nil, dataerr

}

/************** Get  WayBill Ends ********************/
