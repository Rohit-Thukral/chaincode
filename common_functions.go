/*****Chaincode to perform Common Functionality*****
Methods Involved
CreateAsset
CreateCarton
CreatePallet
GetAsset
GetCarton
GetPallet
UpdateAssetDetails
UpdateCartonDetails
UpdatePalletDetails
UpdatePalletCartonAssetByWayBill
FetchShipmentWayBillIndex
FetchEWWayBillIndex
Author: Mohd Arshad
Dated: 30/7/2017
/*****************************************************/
package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

/************** Create Assets Starts ************************/

func CreateAsset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Create Assets ")

	assetDetailsRequest := parseAssetRequest(args[0])

	return saveAssetDetails(stub, assetDetailsRequest)

}
func parseAssetRequest(jsondata string) AssetDetails {
	res := AssetDetails{}
	json.Unmarshal([]byte(jsondata), &res)
	fmt.Println(res)
	return res
}
func saveAssetDetails(stub shim.ChaincodeStubInterface, createAssetDetailsRequest AssetDetails) ([]byte, error) {
	assetDetails := AssetDetails{}
	assetDetails.AssetSerialNumber = createAssetDetailsRequest.AssetSerialNumber
	assetDetails.AssetModel = createAssetDetailsRequest.AssetModel
	assetDetails.AssetType = createAssetDetailsRequest.AssetType
	assetDetails.AssetMake = createAssetDetailsRequest.AssetMake
	assetDetails.AssetCOO = createAssetDetailsRequest.AssetCOO
	assetDetails.AssetMaufacture = createAssetDetailsRequest.AssetMaufacture
	assetDetails.AssetStatus = createAssetDetailsRequest.AssetStatus
	assetDetails.CreatedBy = createAssetDetailsRequest.CreatedBy
	assetDetails.CreatedDate = createAssetDetailsRequest.CreatedDate
	assetDetails.ModifiedBy = createAssetDetailsRequest.ModifiedBy
	assetDetails.ModifiedDate = createAssetDetailsRequest.ModifiedDate
	assetDetails.PalletSerialNumber = createAssetDetailsRequest.PalletSerialNumber
	assetDetails.CartonSerialNumber = createAssetDetailsRequest.CartonSerialNumber
	assetDetails.MshipmentNumber = createAssetDetailsRequest.MshipmentNumber
	assetDetails.DcShipmentNumber = createAssetDetailsRequest.DcShipmentNumber
	assetDetails.MwayBillNumber = createAssetDetailsRequest.MwayBillNumber
	assetDetails.DcWayBillNumber = createAssetDetailsRequest.DcWayBillNumber
	assetDetails.EwWayBillNumber = createAssetDetailsRequest.EwWayBillNumber
	assetDetails.MShipmentDate = createAssetDetailsRequest.MShipmentDate
	assetDetails.DcShipmentDate = createAssetDetailsRequest.DcShipmentDate
	assetDetails.MWayBillDate = createAssetDetailsRequest.MWayBillDate
	assetDetails.DcWayBillDate = createAssetDetailsRequest.DcWayBillDate
	assetDetails.EwWayBillDate = createAssetDetailsRequest.EwWayBillDate
	fmt.Println("assetDetails data to update-->", assetDetails)

	dataToStore, _ := json.Marshal(assetDetails)

	err := DumpData(stub, assetDetails.AssetSerialNumber, string(dataToStore))
	if err != nil {
		fmt.Println("Could not save Assets Details to ledger", err)
		return nil, err
	}

	resp := BlockchainResponse{}
	resp.Err = "000"
	resp.Message = assetDetails.AssetSerialNumber

	respString, _ := json.Marshal(resp)

	fmt.Println("Successfully saved Asset Details")
	return []byte(respString), nil

}

/************** Create Assets Ends ************************/

/************** Create Carton Starts ************************/
func CreateCarton(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Create Cortons ")

	cartonDetailslRequest := parseCartonRequest(args[0])

	return saveCartonDetails(stub, cartonDetailslRequest)

}
func parseCartonRequest(jsondata string) CartonDetails {
	res := CartonDetails{}
	json.Unmarshal([]byte(jsondata), &res)
	fmt.Println(res)
	return res
}
func saveCartonDetails(stub shim.ChaincodeStubInterface, createCartonDetailsRequest CartonDetails) ([]byte, error) {
	cartonDetails := CartonDetails{}
	cartonDetails.CartonSerialNumber = createCartonDetailsRequest.CartonSerialNumber
	cartonDetails.CartonModel = createCartonDetailsRequest.CartonModel
	cartonDetails.CartonStatus = createCartonDetailsRequest.CartonStatus
	cartonDetails.CartonCreationDate = createCartonDetailsRequest.CartonCreationDate
	cartonDetails.PalletSerialNumber = createCartonDetailsRequest.PalletSerialNumber
	cartonDetails.AssetSerialNumber = createCartonDetailsRequest.AssetSerialNumber
	cartonDetails.MshipmentNumber = createCartonDetailsRequest.MshipmentNumber
	cartonDetails.DcShipmentNumber = createCartonDetailsRequest.DcShipmentNumber
	cartonDetails.MwayBillNumber = createCartonDetailsRequest.MwayBillNumber
	cartonDetails.DcWayBillNumber = createCartonDetailsRequest.DcWayBillNumber
	cartonDetails.EwWayBillNumber = createCartonDetailsRequest.EwWayBillNumber
	cartonDetails.Dimensions = createCartonDetailsRequest.Dimensions
	cartonDetails.Weight = createCartonDetailsRequest.Weight
	cartonDetails.MShipmentDate = createCartonDetailsRequest.MShipmentDate
	cartonDetails.DcShipmentDate = createCartonDetailsRequest.DcShipmentDate
	cartonDetails.MWayBillDate = createCartonDetailsRequest.MWayBillDate
	cartonDetails.DcWayBillDate = createCartonDetailsRequest.DcWayBillDate
	cartonDetails.EwWayBillDate = createCartonDetailsRequest.EwWayBillDate
	fmt.Println("cartondetails data update---->", cartonDetails)
	dataToStore, _ := json.Marshal(cartonDetails)

	err := DumpData(stub, cartonDetails.CartonSerialNumber, string(dataToStore))
	if err != nil {
		fmt.Println("Could not save Carton Details to ledger", err)
		return nil, err
	}

	resp := BlockchainResponse{}
	resp.Err = "000"
	resp.Message = cartonDetails.CartonSerialNumber

	respString, _ := json.Marshal(resp)

	fmt.Println("Successfully saved Carton Details")
	return []byte(respString), nil

}

/************** Create Carton Ends ************************/
/************** Create Pallets Starts ************************/
func CreatePallet(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Create Pallets ")

	palletDetailslRequest := parsePalletRequest(args[0])

	return savePalletDetails(stub, palletDetailslRequest)

}
func parsePalletRequest(jsondata string) PalletDetails {
	res := PalletDetails{}
	json.Unmarshal([]byte(jsondata), &res)
	fmt.Println(res)
	return res
}
func savePalletDetails(stub shim.ChaincodeStubInterface, createPalletDetailsRequest PalletDetails) ([]byte, error) {
	palletDetails := PalletDetails{}
	palletDetails.PalletSerialNumber = createPalletDetailsRequest.PalletSerialNumber
	palletDetails.PalletModel = createPalletDetailsRequest.PalletModel
	palletDetails.PalletStatus = createPalletDetailsRequest.PalletStatus
	palletDetails.CartonSerialNumber = createPalletDetailsRequest.CartonSerialNumber
	palletDetails.PalletCreationDate = createPalletDetailsRequest.PalletCreationDate
	palletDetails.AssetSerialNumber = createPalletDetailsRequest.AssetSerialNumber
	palletDetails.MshipmentNumber = createPalletDetailsRequest.MshipmentNumber
	palletDetails.DcShipmentNumber = createPalletDetailsRequest.DcShipmentNumber
	palletDetails.MwayBillNumber = createPalletDetailsRequest.MwayBillNumber
	palletDetails.DcWayBillNumber = createPalletDetailsRequest.DcWayBillNumber
	palletDetails.EwWayBillNumber = createPalletDetailsRequest.EwWayBillNumber
	palletDetails.Dimensions = createPalletDetailsRequest.Dimensions
	palletDetails.Weight = createPalletDetailsRequest.Weight
	palletDetails.MShipmentDate = createPalletDetailsRequest.MShipmentDate
	palletDetails.DcShipmentDate = createPalletDetailsRequest.DcShipmentDate
	palletDetails.MWayBillDate = createPalletDetailsRequest.MWayBillDate
	palletDetails.DcWayBillDate = createPalletDetailsRequest.DcWayBillDate
	palletDetails.EwWayBillDate = createPalletDetailsRequest.EwWayBillDate
	fmt.Println("pallet data to update-->", palletDetails)
	dataToStore, _ := json.Marshal(palletDetails)

	err := DumpData(stub, palletDetails.PalletSerialNumber, string(dataToStore))
	if err != nil {
		fmt.Println("Could not save Pallet Details to ledger", err)
		return nil, err
	}

	resp := BlockchainResponse{}
	resp.Err = "000"
	resp.Message = palletDetails.PalletSerialNumber

	respString, _ := json.Marshal(resp)

	fmt.Println("Successfully saved Pallet Details")
	return []byte(respString), nil

}

/************** Create Pallets Ends ************************/
/************** View Asset Starts ************************/
func GetAsset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering GetAsset " + args[0])

	AssetSerialNumber := args[0]

	assetData, dataerr := fetchAssetDetails(stub, AssetSerialNumber)
	if dataerr == nil {

		dataToStore, _ := json.Marshal(assetData)
		return []byte(dataToStore), nil

	}

	return nil, dataerr

}
func fetchAssetDetails(stub shim.ChaincodeStubInterface, AssetSerialNumber string) (AssetDetails, error) {
	var assetDetails AssetDetails

	indexByte, err := stub.GetState(AssetSerialNumber)
	if err != nil {
		fmt.Println("Could not retrive Asset Details ", err)
		return assetDetails, err
	}

	if marshErr := json.Unmarshal(indexByte, &assetDetails); marshErr != nil {
		fmt.Println("Could not retrieve Asset Details from ledger", marshErr)
		return assetDetails, marshErr
	}

	return assetDetails, nil

}

/************** View Asset Ends ************************/

/************** View Carton Starts ************************/
func GetCarton(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering GetPallet " + args[0])

	cartonSerialNo := args[0]

	cartonData, dataerr := fetchCartonDetails(stub, cartonSerialNo)
	if dataerr == nil {

		dataToStore, _ := json.Marshal(cartonData)
		return []byte(dataToStore), nil

	}

	return nil, dataerr

}
func fetchCartonDetails(stub shim.ChaincodeStubInterface, cartonSerialNo string) (CartonDetails, error) {
	var cartonDetails CartonDetails

	indexByte, err := stub.GetState(cartonSerialNo)
	if err != nil {
		fmt.Println("Could not retrive Carton Details ", err)
		return cartonDetails, err
	}

	if marshErr := json.Unmarshal(indexByte, &cartonDetails); marshErr != nil {
		fmt.Println("Could not retrieve Carton Details from ledger", marshErr)
		return cartonDetails, marshErr
	}

	return cartonDetails, nil

}

/************** View Carton Ends ************************/
/************** View Pallet Starts ************************/
func GetPallet(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering GetPallet " + args[0])

	palletSerialNo := args[0]

	palletData, dataerr := fetchPalletDetails(stub, palletSerialNo)
	if dataerr == nil {

		dataToStore, _ := json.Marshal(palletData)
		return []byte(dataToStore), nil

	}

	return nil, dataerr

}
func fetchPalletDetails(stub shim.ChaincodeStubInterface, palletSerialNo string) (PalletDetails, error) {
	var palletDetails PalletDetails

	indexByte, err := stub.GetState(palletSerialNo)
	if err != nil {
		fmt.Println("Could not retrive Pallet Details ", err)
		return palletDetails, err
	}

	if marshErr := json.Unmarshal(indexByte, &palletDetails); marshErr != nil {
		fmt.Println("Could not retrieve Pallet Details from ledger", marshErr)
		return palletDetails, marshErr
	}

	return palletDetails, nil

}

/************** View Pallet Ends ************************/

/************** Update Asset Details Starts ************************/
func UpdateAssetDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Update Asset Details")
	AssetSerialNumber := args[0]
	wayBillNumber := args[1]
	assetDetails, _ := fetchAssetDetails(stub, AssetSerialNumber)

	assetDetails.EwWayBillNumber = wayBillNumber

	fmt.Println("Updated Entity", assetDetails)
	dataToStore, _ := json.Marshal(assetDetails)
	err := DumpData(stub, AssetSerialNumber, string(dataToStore))
	if err != nil {
		fmt.Println("Could not save Entity WayBill Mapping to ledger", err)
		return nil, err
	}

	resp := BlockchainResponse{}
	resp.Err = "000"
	resp.Message = AssetSerialNumber

	respString, _ := json.Marshal(resp)

	fmt.Println("Successfully saved Entity WayBill Mapping")
	return []byte(respString), nil

}

/************** Update Asset Details Ends ************************/

/************** Update Carton Details Starts ************************/
func UpdateCartonDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Update Carton Details")
	cartonSerialNo := args[0]
	wayBillNumber := args[1]
	cartonDetails, _ := fetchCartonDetails(stub, cartonSerialNo)
	cartonDetails.MwayBillNumber = wayBillNumber
	fmt.Println("Updated Entity", cartonDetails)
	dataToStore, _ := json.Marshal(cartonDetails)
	err := DumpData(stub, cartonSerialNo, string(dataToStore))
	if err != nil {
		fmt.Println("Could not save Pallet Details to ledger", err)
		return nil, err
	}

	resp := BlockchainResponse{}
	resp.Err = "000"
	resp.Message = cartonSerialNo

	respString, _ := json.Marshal(resp)

	fmt.Println("Successfully saved Carton Details")
	return []byte(respString), nil
}

/************** Update Carton Details Ends ************************/

/************** Update Pallet Details Starts ************************/
func UpdatePalletDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Update Pallet Details")
	palletSerialNo := args[0]
	wayBillNumber := args[1]
	palletDetails, _ := fetchPalletDetails(stub, palletSerialNo)
	palletDetails.MwayBillNumber = wayBillNumber
	fmt.Println("Updated Entity", palletDetails)
	dataToStore, _ := json.Marshal(palletDetails)
	err := DumpData(stub, palletSerialNo, string(dataToStore))
	if err != nil {
		fmt.Println("Could not save Pallet Details to ledger", err)
		return nil, err
	}

	resp := BlockchainResponse{}
	resp.Err = "000"
	resp.Message = palletSerialNo

	respString, _ := json.Marshal(resp)

	fmt.Println("Successfully saved Pallet Details")
	return []byte(respString), nil

}

/************** Update Pallet Details Ends ************************/

/************** Update Pallet Carton Asset Details Starts ************************/
func UpdatePalletCartonAssetByWayBill(stub shim.ChaincodeStubInterface, wayBillRequest ShipmentWayBill, source string, ewWaybillId string) ([]byte, error) {
	fmt.Println("Entering Update Pallet Carton Asset Details", wayBillRequest.ShipmentNumber)
	// Start Loop for Pallet Nos
	lenOfArray := len(wayBillRequest.PalletsSerialNumber)
	fmt.Println("palletData=lenOfArray==", lenOfArray)

	for q := 0; q < lenOfArray; q++ {

		fmt.Println("iiiiiiiiiii>...............", q)
		palletData, err := fetchPalletDetails(stub, wayBillRequest.PalletsSerialNumber[q])
		fmt.Println("palletData===", palletData)
		if err != nil {
			fmt.Println("Error while retrieveing the Pallet Details", err)
			return nil, err
		}

		if source == SHIPMENT {
			palletData.MshipmentNumber = wayBillRequest.ShipmentNumber
		} else if source == WAYBILL {
			palletData.MwayBillNumber = wayBillRequest.WayBillNumber
		} else if source == DCSHIPMENT {
			palletData.DcShipmentNumber = wayBillRequest.ShipmentNumber
		} else if source == DCWAYBILL {
			palletData.DcWayBillNumber = wayBillRequest.WayBillNumber
		}
		savePalletDetails(stub, palletData)
		fmt.Println("After savepalletDetails==")

		//Start Loop for Carton Nos
		lenOfcartonArray := len(palletData.CartonSerialNumber)
		fmt.Println("Carton lenOfArray==", lenOfcartonArray)

		for j := 0; j < lenOfcartonArray; j++ {

			cartonData, err := fetchCartonDetails(stub, palletData.CartonSerialNumber[j])
			fmt.Println("cartonData===", cartonData)
			fmt.Println("jjjjjjjjjjjjjj>...............", j)

			if err != nil {
				fmt.Println("Error while retrieveing the Carton Details", err)
				return nil, err
			}
			if source == SHIPMENT {
				cartonData.MshipmentNumber = wayBillRequest.ShipmentNumber
			} else if source == WAYBILL {
				cartonData.MwayBillNumber = wayBillRequest.WayBillNumber
			} else if source == DCSHIPMENT {
				cartonData.DcShipmentNumber = wayBillRequest.ShipmentNumber
			} else if source == DCWAYBILL {
				cartonData.DcWayBillNumber = wayBillRequest.WayBillNumber
			}
			saveCartonDetails(stub, cartonData)
			fmt.Println("after save carton===")

		} //End Loop for Carton Nos
		fmt.Println("iiiiiiiiiii after carton save----...............", q)

		//Start Loop for Asset Nos
		lenOfassetArray := len(palletData.AssetSerialNumber)
		fmt.Println("assets lenOfArray===", lenOfassetArray)
		for k := 0; k < lenOfassetArray; k++ {

			assetData, err := fetchAssetDetails(stub, palletData.AssetSerialNumber[k])
			fmt.Println("assetData===", assetData)
			fmt.Println("kkkkkkkkkkkkkkkkkkkkkkkkkk>...............", k)

			if err != nil {
				fmt.Println("Error while retrieveing the Asset Details", err)
				return nil, err
			}
			if source == SHIPMENT {
				assetData.MshipmentNumber = wayBillRequest.ShipmentNumber
			} else if source == WAYBILL {
				assetData.MwayBillNumber = wayBillRequest.WayBillNumber
			} else if source == DCSHIPMENT {
				assetData.DcShipmentNumber = wayBillRequest.ShipmentNumber
			} else if source == DCWAYBILL {
				assetData.DcWayBillNumber = wayBillRequest.WayBillNumber
			}
			saveAssetDetails(stub, assetData)
			fmt.Println("after save assets===")
		} //End Loop for Asset Nos
		fmt.Println("iiiiiiiiiii>  after asset save-------------...............", q)

	}
	fmt.Println("after all save--------------->")

	resp := BlockchainResponse{}
	resp.Err = "000"
	resp.Message = ""

	respString, _ := json.Marshal(resp)

	fmt.Println("Successfully saved Pallet Carton Asset Details")
	return []byte(respString), nil
}

/************** Update Pallet Carton Asset Details  Ends ************************/

//******************fetch shipmentwaybillIndex****************//
func FetchShipmentWayBillIndex(stub shim.ChaincodeStubInterface, shipmentidkey string) (ShipmentWayBillIndex, error) {
	var shipmentWayBill ShipmentWayBillIndex

	indexByte, err := stub.GetState(shipmentidkey)
	if err != nil {
		fmt.Println("Could not retrive  Shipment WayBill ", err)
		return shipmentWayBill, err
	}

	fmt.Println("shipmentWayBill : ", string(indexByte), "index---->", indexByte)

	json.Unmarshal(indexByte, &shipmentWayBill)

	/*if marshErr := ; marshErr != nil {
		fmt.Println("Could not retrieve Shipment WayBill from ledger", marshErr)
		return shipmentWayBill, marshErr
	}*/

	fmt.Println("======================")
	fmt.Println("shipmentwaybill index after unmarshal--->", shipmentWayBill)
	fmt.Println("======================")

	return shipmentWayBill, nil

}

func FetchEWWayBillIndex(stub shim.ChaincodeStubInterface, ewwaybillKey string) (AllEWWayBill, error) {
	var shipmentWayBill AllEWWayBill

	indexByte, err := stub.GetState(ewwaybillKey)
	if err != nil {
		fmt.Println("Could not retrive  Shipment WayBill ", err)
		return shipmentWayBill, err
	}

	json.Unmarshal(indexByte, &shipmentWayBill)

	fmt.Println("======================")
	fmt.Println(shipmentWayBill)
	fmt.Println("======================")

	return shipmentWayBill, nil

}

func SaveShipmentWaybillIndex(stub shim.ChaincodeStubInterface, shipmentids ShipmentWayBillIndex) ([]byte, error) {
	dataToStore, _ := json.Marshal(shipmentids)
	fmt.Println("save shipmentwaybillids..", dataToStore)
	err := DumpData(stub, "ShipmentWayBillIndex", string(dataToStore))
	if err != nil {
		fmt.Println("Could not save shipmentindex to ledger", err)
		return nil, err
	}

	resp := BlockchainResponse{}
	resp.Err = "000"
	resp.Message = "ShipmentWayBillIndex"
	respString, _ := json.Marshal(resp)

	fmt.Println("Successfully saved ShipmentWayBillIndex IDs")
	return []byte(respString), nil

}
func SaveEWWaybillIndex(stub shim.ChaincodeStubInterface, ewwaybillids AllEWWayBill) ([]byte, error) {
	dataToStore, _ := json.Marshal(ewwaybillids)
	fmt.Println("save ewwaybillids..", dataToStore)
	err := DumpData(stub, "AllEWWayBill", string(dataToStore))
	if err != nil {
		fmt.Println("Could not save shipmentindex to ledger", err)
		return nil, err
	}

	resp := BlockchainResponse{}
	resp.Err = "000"
	resp.Message = "shipmentindex"
	respString, _ := json.Marshal(resp)

	fmt.Println("Successfully saved shipmentindex IDs")
	return []byte(respString), nil

}

/////////////////////////////////////
///Update Cuustodian

func UpdateShipmentCustodianHistoryList(stub shim.ChaincodeStubInterface, shipmentwaybill ShipmentWayBill) []CustodianHistoryDetail {
	var custodianHistoryDetail CustodianHistoryDetail
	custodianHistoryDetail.CustodianName = shipmentwaybill.Custodian
	if shipmentwaybill.WayBillNumber == "" {
		custodianHistoryDetail.Comments = shipmentwaybill.Comments
		custodianHistoryDetail.CustodianTime = shipmentwaybill.ShipmentCreationDate
	} else {
		custodianHistoryDetail.Comments = shipmentwaybill.TpComments
		custodianHistoryDetail.CustodianTime = shipmentwaybill.WayBillCreationDate
	}
	var custodianHistoryDetailList []CustodianHistoryDetail
	indexbyte, err := stub.GetState(shipmentwaybill.ShipmentNumber)
	var tmpshipmentwaybill ShipmentWayBill
	if err != nil {
		fmt.Println("custodianHistory already not available-->")
		custodianHistoryDetailList = append(custodianHistoryDetailList, custodianHistoryDetail)
	} else {
		json.Unmarshal(indexbyte, &tmpshipmentwaybill)
		fmt.Println("custodianHistory already available-->", shipmentwaybill.CustodianHistory)
		custodianHistoryDetailList = append(tmpshipmentwaybill.CustodianHistory, custodianHistoryDetail)

	}
	return custodianHistoryDetailList
}
func UpdateShipmentCustodianHistoryListForeWWaybill(stub shim.ChaincodeStubInterface, custodianName string, comment string, custodianTime string, shipmentNumber string) []CustodianHistoryDetail {
	var custodianHistoryDetail CustodianHistoryDetail
	custodianHistoryDetail.CustodianName = custodianName
	custodianHistoryDetail.Comments = comment
	custodianHistoryDetail.CustodianTime = custodianTime

	var custodianHistoryDetailList []CustodianHistoryDetail
	indexbyte, err := stub.GetState(shipmentNumber)
	var tmpshipmentwaybill ShipmentWayBill
	if err != nil {
		fmt.Println("custodianHistory already  not available-->")
		custodianHistoryDetailList = append(custodianHistoryDetailList, custodianHistoryDetail)
	} else {
		fmt.Println("custodianHistory already available-->", tmpshipmentwaybill.CustodianHistory)
		json.Unmarshal(indexbyte, &tmpshipmentwaybill)
		custodianHistoryDetailList = append(tmpshipmentwaybill.CustodianHistory, custodianHistoryDetail)

	}
	return custodianHistoryDetailList
}
func UpdateEWWaybillCustodianHistoryList(stub shim.ChaincodeStubInterface, ewwayBill EWWayBill) []CustodianHistoryDetail {
	var custodianHistoryDetail CustodianHistoryDetail
	custodianHistoryDetail.CustodianName = ewwayBill.Custodian
	if ewwayBill.EwWayBillModifiedDate == "" {
		custodianHistoryDetail.Comments = ewwayBill.Comments
		custodianHistoryDetail.CustodianTime = ewwayBill.EwWayBillCreationDate
	} else {
		custodianHistoryDetail.Comments = ewwayBill.Comments
		custodianHistoryDetail.CustodianTime = ewwayBill.EwWayBillModifiedDate
	}
	var custodianHistoryDetailList []CustodianHistoryDetail
	indexwaybill, err := stub.GetState(ewwayBill.EwWayBillNumber)
	var tmpewwaybill EWWayBill
	if err != nil {
		fmt.Println("custodianHistory for ewwaybill already not available-->")
		custodianHistoryDetailList = append(custodianHistoryDetailList, custodianHistoryDetail)
	} else {
		json.Unmarshal(indexwaybill, &tmpewwaybill)
		fmt.Println("custodianHistory for ew already available-->", ewwayBill.CustodianHistory)
		custodianHistoryDetailList = append(tmpewwaybill.CustodianHistory, custodianHistoryDetail)

	}
	return custodianHistoryDetailList
}

//End Update Custodian
