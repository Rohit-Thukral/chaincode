/*****Chaincode to perform Waybill realeted task*****
Methods Involved
CreateWayBill : Used for Creating Waybill

Author: Mohd Arshad
Dated: 30/7/2017
/*****************************************************/

package main

import (
	"encoding/json"
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
	shipmentDetails.Status = wayBillRequest.Status
	shipmentDetails.WaybillImage = wayBillRequest.WaybillImage
	shipmentDetails.CustodianHistory = UpdateShipmentCustodianHistoryList(stub, shipmentDetails)

	_, cartonsSerialNumber, assetsSerialNumber, _ := UpdatePalletCartonAssetByWayBill(stub, wayBillRequest, WAYBILL, "")
	shipmentDetails.CartonsSerialNumber = cartonsSerialNumber
	shipmentDetails.AssetsSerialNumber = assetsSerialNumber
	shipmentDetails.WaybillImage = wayBillRequest.WaybillImage
	saveResult, errMsg := saveShipmentWayBill(stub, shipmentDetails)
	fmt.Println("Start of Transaction Details Store Methods............")
	saveResultRes := BlockchainResponse{}
	json.Unmarshal([]byte(saveResult), &saveResultRes)

	transactionDet := TransactionDetails{}
	transactionDet.TransactionId = saveResultRes.TxID
	transactionDet.TransactionTime = shipmentDetails.WayBillCreationDate
	if errMsg != nil {
		transactionDet.Status = "Failure"
	} else {
		transactionDet.Status = "Success"
	}
	transactionDet.FromUserId = shipmentDetails.Custodian
	transactionDet.ToUserId = append(transactionDet.ToUserId, shipmentDetails.Consignee)
	fmt.Println("Start of Transaction Details Store Methods transactionDet.TransactionId............", transactionDet.TransactionId)
	fmt.Println("Start of Transaction Details Store Methods transactionDet.status............", transactionDet.Status)
	fmt.Println("Start of Transaction Details Store Methods transactionDet.FromUserId............", transactionDet.FromUserId)
	fmt.Println("Start of Transaction Details Store Methods transactionDet.ToUserId............", transactionDet.ToUserId)
	fmt.Println("Start of Transaction Details Store Methods transactionDet.TransactionTime............", transactionDet.TransactionTime)
	_ = saveTransactionDetails(stub, transactionDet)
	fmt.Println("End of Transaction Details Store Methods............")
	return saveResult, errMsg
}

/************** Create Way Bill Ends ************************/

/************** Update Waybill Starts *********************/
func UpdateWayBill(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Update Waybill", args[0])
	shipmentRequest := parseShipmentWayBillRequest(args[0])
	wayBilldata, _ := fetchShipmentWayBillData(stub, shipmentRequest.ShipmentNumber)
	shipmentRequest.CustodianHistory = UpdateShipmentCustodianHistoryList(stub, shipmentRequest)
	shipmentRequest.SupportiveDocuments = wayBilldata.SupportiveDocuments
	shipmentRequest.ShipmentImage = wayBilldata.ShipmentImage
	shipmentRequest.WaybillImage = wayBilldata.WaybillImage

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

func ViewWayBill(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering ViewWayBill " + args[0])

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
