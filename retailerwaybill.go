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
func CreateRetailerWayBill(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Create Retailer WayBill", args[0])

	retailerwayBillRequest := parseShipmentWayBillRequest(args[0])

	retailershipmentDetails, err := fetchShipmentWayBillData(stub, retailerwayBillRequest.ShipmentNumber)

	if err != nil {
		fmt.Println("Error while retrieveing the retailer Shipment Details", err)
		return nil, err
	}
	retailershipmentDetails.WayBillNumber = retailerwayBillRequest.WayBillNumber
	retailershipmentDetails.VehicleNumber = retailerwayBillRequest.VehicleNumber
	retailershipmentDetails.VehicleType = retailerwayBillRequest.VehicleType
	retailershipmentDetails.PickupDate = retailerwayBillRequest.PickupDate
	retailershipmentDetails.Custodian = retailerwayBillRequest.Custodian
	retailershipmentDetails.TpComments = retailerwayBillRequest.TpComments
	retailershipmentDetails.WayBillCreationDate = retailerwayBillRequest.WayBillCreationDate
	retailershipmentDetails.WayBillCreatedBy = retailerwayBillRequest.WayBillCreatedBy
	retailershipmentDetails.VehicleType = retailerwayBillRequest.VehicleType
	retailershipmentDetails.EntityName = retailerwayBillRequest.EntityName
	retailershipmentDetails.Status = retailerwayBillRequest.Status
	retailershipmentDetails.WaybillImage = retailerwayBillRequest.WaybillImage
	retailershipmentDetails.CustodianHistory = UpdateShipmentCustodianHistoryList(stub, retailershipmentDetails)
	retailershipmentDetails.PalletsSerialNumber, _ = getPalletSerialNoByCartonNo(stub, retailerwayBillRequest.CartonsSerialNumber)
	_, cartonsSerialNumber, assetsSerialNumber, _ := UpdatePalletCartonAssetByWayBill(stub, retailerwayBillRequest, RETAILERWAYBILL, "")
	retailershipmentDetails.CartonsSerialNumber = cartonsSerialNumber
	retailershipmentDetails.AssetsSerialNumber = assetsSerialNumber

	UpdateEntityWayBillMapping(stub, retailershipmentDetails.EntityName, retailershipmentDetails.WayBillNumber, retailershipmentDetails.CountryFrom)
	err = DumpData(stub, retailershipmentDetails.WayBillNumber, retailerwayBillRequest.ShipmentNumber)
	if err != nil {
		fmt.Println("Could not save WayBill to ledger", err)
		return nil, err
	}
	saveResult, errMsg := saveShipmentWayBill(stub, retailershipmentDetails)

	fmt.Println("Start of Transaction Details Store Methods............")
	saveResultRes := BlockchainResponse{}
	json.Unmarshal([]byte(saveResult), &saveResultRes)

	transactionDet := TransactionDetails{}
	transactionDet.TransactionId = saveResultRes.TxID
	transactionDet.TransactionTime = retailershipmentDetails.WayBillCreationDate
	if errMsg != nil {
		transactionDet.Status = "Failure"
	} else {
		transactionDet.Status = "Success"
	}
	transactionDet.FromUserId = retailershipmentDetails.Custodian
	transactionDet.ToUserId = append(transactionDet.ToUserId, retailershipmentDetails.Consignee)
	transactionDet.ToUserId = append(transactionDet.ToUserId, retailershipmentDetails.EntityName)
	fmt.Println("Start of Transaction Details Store Methods transactionDet.TransactionId............", transactionDet.TransactionId)
	//fmt.Println("Start of Transaction Details Store Methods transactionDet.status............", transactionDet.Status)
	//fmt.Println("Start of Transaction Details Store Methods transactionDet.FromUserId............", transactionDet.FromUserId)
	//fmt.Println("Start of Transaction Details Store Methods transactionDet.ToUserId............", transactionDet.ToUserId)
	//fmt.Println("Start of Transaction Details Store Methods transactionDet.TransactionTime............", transactionDet.TransactionTime)
	_ = saveTransactionDetails(stub, transactionDet)
	fmt.Println("End of Transaction Details Store Methods............")
	return saveResult, errMsg

}

/************** Create DC Way Bill Ends ************************/

/************** Update Waybill Starts *********************/
func UpdateRetailerWayBill(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Update DCWaybill", args[0])
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

func ViewRetailerWayBill(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
