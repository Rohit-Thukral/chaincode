/*****Chaincode to perform retailer Shipment*****
Methods Involved
CreateRetailerShipment : Create retailer Shipment

Author: santosh
Dated: 09/05/2017
/*****************************************************/
package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

/*****Chaincode to perfor shipment realeted task*****
Methods Involved
CreateShipment : Used for Creating retailer Shipment

Author: santosh
Dated: 09/05/2017
/*****************************************************/

/************** Create DC Shipment Starts ***********************/
func CreateRetailerShipment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Retailer Create Shipment", args[0])
	shipmentRequest := parseShipmentWayBillRequest(args[0])
	shipmentRequest.PalletsSerialNumber, _ = getPalletSerialNoByCartonNo(stub, shipmentRequest.CartonsSerialNumber)
	_, cartonsSerialNumber, assetsSerialNumber, _ := UpdatePalletCartonAssetByWayBill(stub, shipmentRequest, RETAILERSHIPMENT, "")
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
func UpdateRetailerShipment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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

func ViewRetailerShipment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering ViewRetailerShipment " + args[0])

	shipmentNo := args[0]
	wayBilldata, dataerr := fetchShipmentWayBillData(stub, shipmentNo)
	if dataerr == nil {

		dataToStore, _ := json.Marshal(wayBilldata)
		return []byte(dataToStore), nil

	}

	return nil, dataerr

}

/************** Get  WayBill Ends ********************/

//*************get pallet array by carton array****************//
func getPalletSerialNoByCartonNo(stub shim.ChaincodeStubInterface, cartonNo []string) ([]string, error) {
	fmt.Println("Entering getPalletSerialNoByCartonNo ", cartonNo)
	lenOfcartonArray := len(cartonNo)
	fmt.Println("length of carton array...", lenOfcartonArray)
	var palletNo []string
	for ca := 0; ca < lenOfcartonArray; ca++ {
		cartonData, _ := fetchCartonDetails(stub, cartonNo[ca])
		if stringNotExistInArray(palletNo, cartonData.PalletSerialNumber) {
			palletNo = append(palletNo, cartonData.PalletSerialNumber)
		}
	}

	return palletNo, nil

}

