/*****Chaincode to perform Export Warehouse Way Bill*****
Methods Involved
CreateEWWayBill : Create Export Warehouse Way Bill
ViewEWWayBill: Fetch Export Warehouse Way Bill

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
/***********************************************************/

/************** Create Export Warehouse WayBill Starts ************************/
func CreateEWWayBill(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Export Warehouse WayBill ")

	ewWayBillRequest := parseEWWayBillRequest(args[0])
	lenOfArray := len(ewWayBillRequest.WayBillsNumber)
	for i := 0; i < lenOfArray; i++ {
		wayBillShipmentMapping, err := FetchWayBillShipmentMappingData(stub, ewWayBillRequest.WayBillsNumber[i])
		dcShipmentNumber := wayBillShipmentMapping.DCShipmentNumber
		dcShipmentData, _ := fetchShipmentWayBillData(stub, dcShipmentNumber)
		dcShipmentData.CustodianHistory = UpdateShipmentCustodianHistoryListForeWWaybill(stub, ewWayBillRequest.Custodian, ewWayBillRequest.Comments, ewWayBillRequest.EwWayBillCreationDate, dcShipmentNumber)
		dcShipmentData.EWWayBillNumber = ewWayBillRequest.EwWayBillNumber

		saveShipmentWayBill(stub, dcShipmentData)

		UpdatePalletCartonAssetByWayBill(stub, dcShipmentData, EWWAYBILL, ewWayBillRequest.EwWayBillNumber)
		ewWayBillRequest.ShipmentsNumber = append(ewWayBillRequest.ShipmentsNumber, dcShipmentNumber)
		lenOfArray = len(dcShipmentData.PalletsSerialNumber)
		for j := 0; j < lenOfArray; j++ {
			ewWayBillRequest.PalletsSerialNumber = append(ewWayBillRequest.ShipmentsNumber, dcShipmentData.PalletsSerialNumber[j])
		}
		if err != nil {
			fmt.Println("Could not retrive Shipment WayBill ", err)

		}
		wayBills, _ := fetchEntityWayBillMappingData(stub, ewWayBillRequest.Consigner)
		var tmpWayBillArray []string

		for k := 0; k < len(wayBills.WayBillsNumber); k++ {
			for j := 0; j < len(ewWayBillRequest.WayBillsNumber); j++ {
				if ewWayBillRequest.WayBillsNumber[j] != wayBills.WayBillsNumber[k] {
					tmpWayBillArray = append(tmpWayBillArray, wayBills.WayBillsNumber[k])
				}
			}
		}
		ewWayBillRequest.WayBillsNumber = tmpWayBillArray
	}
	ewWayBillRequest.CustodianHistory = UpdateEWWaybillCustodianHistoryList(stub, ewWayBillRequest)
	saveResult, errMsg := saveEWWayBill(stub, ewWayBillRequest)

	/*********Storing Shipment number in shipmentwaybillindex array to retrieve through inbox*************/
	allEWWaybillidsRequest := AllEWWayBill{}
	allEWWayBillids, err := FetchEWWayBillIndex(stub, "AllEWWayBill")
	fmt.Println("ew waybill ids.....", allEWWayBillids)
	if err != nil {
		allEWWaybillidsRequest.AllWayBillNumber = append(allEWWaybillidsRequest.AllWayBillNumber, ewWayBillRequest.EwWayBillNumber)
		SaveEWWaybillIndex(stub, allEWWaybillidsRequest)
	} else {
		allEWWaybillidsRequest.AllWayBillNumber = append(allEWWayBillids.AllWayBillNumber, ewWayBillRequest.EwWayBillNumber)
		fmt.Println("Updated entity ewWayBillRequest", allEWWaybillidsRequest)
		SaveEWWaybillIndex(stub, allEWWaybillidsRequest)
	}
	/********* End Storing Shipment number in ewWayBillRequest array to retrieve through inbox*************/

	return saveResult, errMsg

}

/************** Update Export Warehouse WayBill Starts ************************/
func UpdateEWWayBill(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Update Export Warehouse WayBill ")

	ewWayBillRequest := parseEWWayBillRequest(args[0])
	lenOfArray := len(ewWayBillRequest.WayBillsNumber)
	for i := 0; i < lenOfArray; i++ {
		wayBillShipmentMapping, _ := FetchWayBillShipmentMappingData(stub, ewWayBillRequest.WayBillsNumber[i])
		dcShipmentNumber := wayBillShipmentMapping.DCShipmentNumber
		dcShipmentData, _ := fetchShipmentWayBillData(stub, dcShipmentNumber)
		dcShipmentData.CustodianHistory = UpdateShipmentCustodianHistoryListForeWWaybill(stub, ewWayBillRequest.Custodian, ewWayBillRequest.Comments, ewWayBillRequest.EwWayBillCreationDate, dcShipmentNumber)

		saveShipmentWayBill(stub, dcShipmentData)
	}
	ewWayBillRequest.CustodianHistory = UpdateEWWaybillCustodianHistoryList(stub, ewWayBillRequest)
	emWayBilldata, _ := fetchEWWayBillData(stub, ewWayBillRequest.EwWayBillNumber)
	ewWayBillRequest.SupportiveDocuments = emWayBilldata.SupportiveDocuments
	saveResult, errMsg := saveEWWayBill(stub, ewWayBillRequest)
	return saveResult, errMsg

}

/***end of updateewwaybill*************/

func parseEWWayBillRequest(jsondata string) EWWayBill {
	res := EWWayBill{}
	json.Unmarshal([]byte(jsondata), &res)
	fmt.Println(res)
	return res
}
func saveEWWayBill(stub shim.ChaincodeStubInterface, createEWWayBillRequest EWWayBill) ([]byte, error) {

	ewWayBill := EWWayBill{}
	ewWayBill.EwWayBillNumber = createEWWayBillRequest.EwWayBillNumber
	ewWayBill.WayBillsNumber = createEWWayBillRequest.WayBillsNumber
	ewWayBill.ShipmentsNumber = createEWWayBillRequest.ShipmentsNumber
	ewWayBill.CountryFrom = createEWWayBillRequest.CountryFrom
	ewWayBill.CountryTo = createEWWayBillRequest.CountryTo
	ewWayBill.Consigner = createEWWayBillRequest.Consigner
	ewWayBill.Consignee = createEWWayBillRequest.Consignee
	ewWayBill.Custodian = createEWWayBillRequest.Custodian
	ewWayBill.CustodianHistory = createEWWayBillRequest.CustodianHistory
	ewWayBill.CustodianTime = createEWWayBillRequest.CustodianTime
	ewWayBill.PersonConsigningGoods = createEWWayBillRequest.PersonConsigningGoods
	ewWayBill.Comments = createEWWayBillRequest.Comments
	ewWayBill.PalletsSerialNumber = createEWWayBillRequest.PalletsSerialNumber
	ewWayBill.AddressOfConsigner = createEWWayBillRequest.AddressOfConsigner
	ewWayBill.AddressOfConsignee = createEWWayBillRequest.AddressOfConsignee
	ewWayBill.ConsignerRegNumber = createEWWayBillRequest.ConsignerRegNumber
	ewWayBill.VesselType = createEWWayBillRequest.VesselType
	ewWayBill.VesselNumber = createEWWayBillRequest.VesselNumber
	ewWayBill.ContainerNumber = createEWWayBillRequest.ContainerNumber
	ewWayBill.ServiceType = createEWWayBillRequest.ServiceType
	//ewWayBill.SupportiveDocumentsList = createEWWayBillRequest.SupportiveDocumentsList
	ewWayBill.EwWayBillCreationDate = createEWWayBillRequest.EwWayBillCreationDate
	ewWayBill.EwWayBillCreatedBy = createEWWayBillRequest.EwWayBillCreatedBy
	ewWayBill.EwWayBillModifiedDate = createEWWayBillRequest.EwWayBillModifiedDate
	ewWayBill.EwWayBillModifiedBy = createEWWayBillRequest.EwWayBillModifiedBy
	ewWayBill.Status = createEWWayBillRequest.Status

	dataToStore, _ := json.Marshal(ewWayBill)

	err := DumpData(stub, ewWayBill.EwWayBillNumber, string(dataToStore))
	if err != nil {
		fmt.Println("Could not save Export Warehouse Way Bill to ledger", err)
		return nil, err
	}

	resp := BlockchainResponse{}
	resp.Err = "000"
	resp.Message = ewWayBill.EwWayBillNumber

	respString, _ := json.Marshal(resp)

	fmt.Println("Successfully saved Export Warehouse Way Bill")
	return []byte(respString), nil

}

/************** Create Export Warehouse WayBill Ends ************************/

/************** View Export Warehouse WayBill Starts ************************/
func ViewEWWayBill(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering ViewEWWayBill " + args[0])

	ewWayBillNumber := args[0]

	emWayBilldata, dataerr := fetchEWWayBillData(stub, ewWayBillNumber)
	if dataerr == nil {

		dataToStore, _ := json.Marshal(emWayBilldata)
		return []byte(dataToStore), nil

	}

	return nil, dataerr

}
func fetchEWWayBillData(stub shim.ChaincodeStubInterface, ewWayBillNumber string) (EWWayBill, error) {
	var ewWayBill EWWayBill

	indexByte, err := stub.GetState(ewWayBillNumber)
	if err != nil {
		fmt.Println("Could not retrive Export Warehouse WayBill ", err)
		return ewWayBill, err
	}

	if marshErr := json.Unmarshal(indexByte, &ewWayBill); marshErr != nil {
		fmt.Println("Could not retrieve Export Warehouse from ledger", marshErr)
		return ewWayBill, marshErr
	}

	return ewWayBill, nil

}

/************** View Export Warehouse WayBill Ends ************************/
