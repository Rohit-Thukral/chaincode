/*****Chaincode to perform Shipment realeted task*****
Methods Involved
CreateShipment : Used for Creating Shipment
ViewShipmentWayBill: Used to fetch Shipment details

Author: Mohd Arshad
Dated: 30/7/2017
/*****************************************************/
package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

/************** Create Shipment Starts *********************/
func CreateShipment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Create Shipment", args[0])
	shipmentRequest := parseShipmentWayBillRequest(args[0])
	UpdatePalletCartonAssetByWayBill(stub, shipmentRequest, SHIPMENT, "")
	fmt.Println("after updatepalletcartonasset............")
	shipmentRequest.CustodianHistory = UpdateShipmentCustodianHistoryList(stub, shipmentRequest)
	saveResult, errMsg := saveShipmentWayBill(stub, shipmentRequest)

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
	return saveResult, errMsg
}

/************** Create Shipment Ends ************************/

/************** Update Shipment Starts *********************/
func UpdateShipment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Update Shipment", args[0])
	shipmentRequest := parseShipmentWayBillRequest(args[0])
	wayBilldata, _ := fetchShipmentWayBillData(stub, shipmentRequest.ShipmentNumber)
	shipmentRequest.CustodianHistory = UpdateShipmentCustodianHistoryList(stub, shipmentRequest)
	shipmentRequest.SupportiveDocuments = wayBilldata.SupportiveDocuments
	return saveShipmentWayBill(stub, shipmentRequest)
}

/************** Update Shipment Ends ************************/

/************** Save Shipment WayBill Starts ****************/
/*This is common code for Save Shipment,WayBill,DCShipment,DCWayBill*/

func parseShipmentWayBillRequest(jsondata string) ShipmentWayBill {
	res := ShipmentWayBill{}
	json.Unmarshal([]byte(jsondata), &res)
	fmt.Println(res)
	return res
}
func saveShipmentWayBill(stub shim.ChaincodeStubInterface, createShipmentWayBillRequest ShipmentWayBill) ([]byte, error) {
	fmt.Println("way Bill no ", createShipmentWayBillRequest.WayBillNumber)
	shipmentWayBill := ShipmentWayBill{}
	shipmentWayBill.WayBillNumber = createShipmentWayBillRequest.WayBillNumber
	shipmentWayBill.ShipmentNumber = createShipmentWayBillRequest.ShipmentNumber
	shipmentWayBill.CountryFrom = createShipmentWayBillRequest.CountryFrom
	shipmentWayBill.CountryTo = createShipmentWayBillRequest.CountryTo
	shipmentWayBill.Consigner = createShipmentWayBillRequest.Consigner
	shipmentWayBill.Consignee = createShipmentWayBillRequest.Consignee
	shipmentWayBill.Custodian = createShipmentWayBillRequest.Custodian
	shipmentWayBill.CustodianHistory = createShipmentWayBillRequest.CustodianHistory
	shipmentWayBill.PersonConsigningGoods = createShipmentWayBillRequest.PersonConsigningGoods
	shipmentWayBill.Comments = createShipmentWayBillRequest.Comments
	shipmentWayBill.TpComments = createShipmentWayBillRequest.TpComments
	shipmentWayBill.VehicleNumber = createShipmentWayBillRequest.VehicleNumber
	shipmentWayBill.VehicleType = createShipmentWayBillRequest.VehicleType
	shipmentWayBill.PickupDate = createShipmentWayBillRequest.PickupDate
	shipmentWayBill.PalletsSerialNumber = createShipmentWayBillRequest.PalletsSerialNumber
	shipmentWayBill.AddressOfConsigner = createShipmentWayBillRequest.AddressOfConsigner
	shipmentWayBill.AddressOfConsignee = createShipmentWayBillRequest.AddressOfConsignee
	shipmentWayBill.ConsignerRegNumber = createShipmentWayBillRequest.ConsignerRegNumber
	shipmentWayBill.Carrier = createShipmentWayBillRequest.Carrier
	shipmentWayBill.VesselType = createShipmentWayBillRequest.VesselType
	shipmentWayBill.VesselNumber = createShipmentWayBillRequest.VesselNumber
	shipmentWayBill.ContainerNumber = createShipmentWayBillRequest.ContainerNumber
	shipmentWayBill.ServiceType = createShipmentWayBillRequest.ServiceType
	shipmentWayBill.ShipmentModel = createShipmentWayBillRequest.ShipmentModel
	shipmentWayBill.PalletsQuantity = createShipmentWayBillRequest.PalletsQuantity
	shipmentWayBill.CartonsQuantity = createShipmentWayBillRequest.CartonsQuantity
	shipmentWayBill.AssetsQuantity = createShipmentWayBillRequest.AssetsQuantity
	shipmentWayBill.ShipmentValue = createShipmentWayBillRequest.ShipmentValue
	shipmentWayBill.EntityName = createShipmentWayBillRequest.EntityName
	shipmentWayBill.ShipmentCreationDate = createShipmentWayBillRequest.ShipmentCreationDate
	shipmentWayBill.EWWayBillNumber = createShipmentWayBillRequest.EWWayBillNumber
	shipmentWayBill.SupportiveDocuments = createShipmentWayBillRequest.SupportiveDocuments
	shipmentWayBill.ShipmentCreatedBy = createShipmentWayBillRequest.ShipmentCreatedBy
	shipmentWayBill.ShipmentModifiedDate = createShipmentWayBillRequest.ShipmentModifiedDate
	shipmentWayBill.ShipmentModifiedBy = createShipmentWayBillRequest.ShipmentModifiedBy
	shipmentWayBill.WayBillCreationDate = createShipmentWayBillRequest.WayBillCreationDate
	shipmentWayBill.WayBillCreatedBy = createShipmentWayBillRequest.WayBillCreatedBy
	shipmentWayBill.WayBillModifiedDate = createShipmentWayBillRequest.WayBillModifiedDate
	shipmentWayBill.WayBillModifiedBy = createShipmentWayBillRequest.WayBillModifiedBy
	shipmentWayBill.Status = createShipmentWayBillRequest.Status
	dataToStore, _ := json.Marshal(shipmentWayBill)
	fmt.Println("shipmentWayBill============ ", shipmentWayBill)
	fmt.Println("dataToStore============ ", dataToStore)

	err := DumpData(stub, shipmentWayBill.ShipmentNumber, string(dataToStore))
	if err != nil {
		fmt.Println("Could not save WayBill to ledger", err)
		return nil, err
	}

	resp := BlockchainResponse{}
	resp.Err = "000"
	resp.Message = shipmentWayBill.ShipmentNumber
	respString, _ := json.Marshal(resp)

	fmt.Println("Successfully saved Way Bill")
	return []byte(respString), nil

}

/************** Save Shipment WayBill Ends *******************/

/************** Get Shipment WayBill Starts ******************/
/*This is common code for Get Shipment,WayBill,DCShipment,DCWayBill*/

func ViewShipmentWayBill(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering ViewWayBill " + args[0])

	shipmentNo := args[0]

	wayBilldata, dataerr := fetchShipmentWayBillData(stub, shipmentNo)
	if dataerr == nil {

		dataToStore, _ := json.Marshal(wayBilldata)
		return []byte(dataToStore), nil

	}

	return nil, dataerr

}
func fetchShipmentWayBillData(stub shim.ChaincodeStubInterface, shipmentNo string) (ShipmentWayBill, error) {
	var shipmentWayBill ShipmentWayBill

	indexByte, err := stub.GetState(shipmentNo)
	if err != nil {
		fmt.Println("Could not retrive  Shipment WayBill ", err)
		return shipmentWayBill, err
	}

	json.Unmarshal(indexByte, &shipmentWayBill)

	fmt.Println("======================shipment data-->")
	fmt.Println(shipmentWayBill)
	fmt.Println("======================")

	return shipmentWayBill, nil

}

/************** Get Shipment WayBill Ends ********************/
