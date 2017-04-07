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
	assetDetails.AssetSerialNo = createAssetDetailsRequest.AssetSerialNo
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

	dataToStore, _ := json.Marshal(assetDetails)

	err := stub.PutState(assetDetails.AssetSerialNo, []byte(dataToStore))
	if err != nil {
		fmt.Println("Could not save Assets Details to ledger", err)
		return nil, err
	}

	resp := BlockchainResponse{}
	resp.Err = "000"
	resp.Message = assetDetails.AssetSerialNo

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
	cartonDetails.CartonSerialNo = createCartonDetailsRequest.CartonSerialNo
	cartonDetails.CartonModel = createCartonDetailsRequest.CartonModel
	cartonDetails.CartonStatus = createCartonDetailsRequest.CartonStatus
	cartonDetails.CartonCreationDate = createCartonDetailsRequest.CartonCreationDate
	cartonDetails.PalletSerialNumber = createCartonDetailsRequest.PalletSerialNumber
	cartonDetails.AssetsSerialNumber = createCartonDetailsRequest.AssetsSerialNumber
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
	dataToStore, _ := json.Marshal(cartonDetails)

	err := stub.PutState(cartonDetails.CartonSerialNo, []byte(dataToStore))
	if err != nil {
		fmt.Println("Could not save Carton Details to ledger", err)
		return nil, err
	}

	resp := BlockchainResponse{}
	resp.Err = "000"
	resp.Message = cartonDetails.CartonSerialNo

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
	palletDetails.PalletSerialNo = createPalletDetailsRequest.PalletSerialNo
	palletDetails.PalletModel = createPalletDetailsRequest.PalletModel
	palletDetails.PalletStatus = createPalletDetailsRequest.PalletStatus
	palletDetails.CartonSerialNumber = createPalletDetailsRequest.CartonSerialNumber
	palletDetails.PalletCreationDate = createPalletDetailsRequest.PalletCreationDate
	palletDetails.AssetsSerialNumber = createPalletDetailsRequest.AssetsSerialNumber
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
	dataToStore, _ := json.Marshal(palletDetails)

	err := stub.PutState(palletDetails.PalletSerialNo, []byte(dataToStore))
	if err != nil {
		fmt.Println("Could not save Pallet Details to ledger", err)
		return nil, err
	}

	resp := BlockchainResponse{}
	resp.Err = "000"
	resp.Message = palletDetails.PalletSerialNo

	respString, _ := json.Marshal(resp)

	fmt.Println("Successfully saved Pallet Details")
	return []byte(respString), nil

}

/************** Create Pallets Ends ************************/
/************** View Asset Starts ************************/
func GetAsset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering GetAsset " + args[0])

	assetSerialNo := args[0]

	assetData, dataerr := fetchAssetDetails(stub, assetSerialNo)
	if dataerr == nil {

		dataToStore, _ := json.Marshal(assetData)
		return []byte(dataToStore), nil

	}

	return nil, dataerr

}
func fetchAssetDetails(stub shim.ChaincodeStubInterface, assetSerialNo string) (AssetDetails, error) {
	var assetDetails AssetDetails

	indexByte, err := stub.GetState(assetSerialNo)
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
	assetSerialNo := args[0]
	wayBillNumber := args[1]
	assetDetails, _ := fetchAssetDetails(stub, assetSerialNo)

	assetDetails.EwWayBillNumber = wayBillNumber

	fmt.Println("Updated Entity", assetDetails)
	dataToStore, _ := json.Marshal(assetDetails)
	err := stub.PutState(assetSerialNo, []byte(dataToStore))
	if err != nil {
		fmt.Println("Could not save Entity WayBill Mapping to ledger", err)
		return nil, err
	}

	resp := BlockchainResponse{}
	resp.Err = "000"
	resp.Message = assetSerialNo

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
	err := stub.PutState(cartonSerialNo, []byte(dataToStore))
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
	err := stub.PutState(palletSerialNo, []byte(dataToStore))
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
	fmt.Println("Entering Update Pallet Carton Asset Details")
	// Start Loop for Pallet Nos
	lenOfArray := len(wayBillRequest.PalletsSerialNumber)
	for i := 0; i < lenOfArray; i++ {

		palletData, err := fetchPalletDetails(stub, wayBillRequest.PalletsSerialNumber[i])

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

		//Start Loop for Carton Nos
		lenOfArray = len(palletData.CartonSerialNumber)
		for i := 0; i < lenOfArray; i++ {

			cartonData, err := fetchCartonDetails(stub, palletData.CartonSerialNumber[i])

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
		} //End Loop for Carton Nos

		//Start Loop for Asset Nos
		lenOfArray = len(palletData.AssetsSerialNumber)
		for i := 0; i < lenOfArray; i++ {

			assetData, err := fetchAssetDetails(stub, palletData.AssetsSerialNumber[i])

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
		} //End Loop for Asset Nos
	}

	resp := BlockchainResponse{}
	resp.Err = "000"
	resp.Message = ""

	respString, _ := json.Marshal(resp)

	fmt.Println("Successfully saved Pallet Carton Asset Details")
	return []byte(respString), nil
}

/************** Update Pallet Carton Asset Details  Ends ************************/
