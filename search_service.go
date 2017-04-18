package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)


type SearchService struct {
}


type Shipment struct {
	ShipmentNumber        string           `json:"shipmentNumber"`
	WayBillNo             string           `json:"wayBillNo"`
	WayBillType           string           `json:"wayBillType"`
	PersonConsigningGoods string           `json:"personConsigningGoods"`
	Consigner             string           `json:"consigner"`
	ConsignerAddress      string           `json:"consignerAddress"`
	Consignee             string           `json:"consignee"`
	ConsigneeAddress      string           `json:"consigneeAddress"`
	ConsigneeRegNo        string           `json:"consigneeRegNo"`
	Quantity              string           `json:"quantity"`
	Pallets               []string         `json:"pallets"`
	Cartons               []string         `json:"cartons"`
	Status                string           `json:"status"`
	ModelNo               string           `json:"modelNo"`
	VehicleNumber         string           `json:"vehicleNumber"`
	VehicleType           string           `json:"vehicleType"`
	PickUpTime            string           `json:"pickUpTime"`
	ValueOfGoods          string           `json:"valueOfGoods"`
	ContainerId           string           `json:"containerId"`
	MasterWayBillRef      []string         `json:"masterWayBillRef"`
	WayBillHistorys       []WayBillHistory `json:"wayBillHistorys"`
	Carrier               string           `json:"carrier"`
	Acl                   []string         `json:"acl"`
	CreatedBy             string           `json:"createdBy"`
	Custodian             string           `json:"custodian"`
	CreatedTimeStamp      string           `json:"createdTimeStamp"`
	UpdatedTimeStamp      string           `json:"updatedTimeStamp"`
}
/************** Asset Search Service Starts ************************/

type SearchAssetRequest struct {
	CallingEntityName string `json:"callingEntityName"`
	AssetId           string `json:"assetId"`
}

type SearchAssetData struct {
	AssetSerialNo  				string     	`json:"assetSerialNo"`
	AssetModel      			string     	`json:"assetModel"`
	AssetType  					string     	`json:"assetType"`
	CartonSerialNumber       	string     	`json:"cartonSerialNumber"`
	PalletSerialNumber       	string     	`json:"palletSerialNumber"`
	Custodian             		string   	`json:"custodian"`
	custodianTime   			string   	`json:"custodianTime"`
	Status             	  		string   	`json:"status"`
}

type SearchAssetResponse struct {
	Data 						[]SearchAssetData		`json:"data"`
	ErrCode        				string     	`json:"errCode"`
	ErrMessage     				string     	`json:"errMessage"`
}

func retrieveShipment(stub shim.ChaincodeStubInterface, shipmentId string) (Shipment, error) {
	var shipment Shipment

	shipmentBytes, err := stub.GetState(shipmentId)
	if err != nil {
		return shipment, err
	} else {
		if marshErr := json.Unmarshal(shipmentBytes, &shipment); marshErr != nil {
			fmt.Println("Could not Unmarshal Asset", marshErr)
			return shipment, marshErr
		}
		return shipment, nil
	}
}
func PrepareSearchAssetResponse(stub shim.ChaincodeStubInterface, asset AssetDetails) ([]byte, error) {
	var resp SearchAssetResponse
	var tmpShipment ShipmentWayBill
	var err error
	var respDataArr []SearchAssetData
	

	
	var respData SearchAssetData
	respData.AssetSerialNo = asset.AssetSerialNumber
	respData.AssetModel = asset.AssetModel
	respData.AssetType = asset.AssetType
	respData.CartonSerialNumber = asset.CartonSerialNumber
	respData.PalletSerialNumber = asset.PalletSerialNumber
	
	tmpShipment, err = fetchShipmentWayBillData(stub, asset.MshipmentNumber)
	
	if err != nil {
		resp.ErrCode = "ERR_DATA"
		resp.ErrMessage = "Unable to get Shipment MshipmentNumber"
		fmt.Println("Error while retrieveing the Shipment Details", err)
		return nil, err
	} else {
		respData.Custodian = tmpShipment.Custodian;
		respData.custodianTime = tmpShipment.WayBillModifiedDate;
		respData.Status = tmpShipment.Status;
		respDataArr = append(respDataArr, respData)
	}
//-----------------------------------------------------------------------//
	var respData2 SearchAssetData
	respData2.AssetSerialNo = asset.AssetSerialNumber
	respData2.AssetModel = asset.AssetModel
	respData2.AssetType = asset.AssetType
	respData2.CartonSerialNumber = asset.CartonSerialNumber
	respData2.PalletSerialNumber = asset.PalletSerialNumber
	
	tmpShipment, err = fetchShipmentWayBillData(stub, asset.DcWayBillNumber)
	
	if err != nil {
		resp.ErrCode = "ERR_DATA"
		resp.ErrMessage = "Unable to get Shipment DcWayBillNumber"
		fmt.Println("Error while retrieveing the Shipment Details", err)
		return nil, err
	} else {
		respData2.Custodian = tmpShipment.Custodian;
		respData2.custodianTime = tmpShipment.WayBillModifiedDate;
		respData2.Status = tmpShipment.Status;
		respDataArr = append(respDataArr, respData2)
	}
	
	resp.Data = respDataArr
	return json.Marshal(resp)

}

func parseSearchAssetRequest(requestParam string) (SearchAssetRequest, error) {
	var request SearchAssetRequest

	if marshErr := json.Unmarshal([]byte(requestParam), &request); marshErr != nil {
		fmt.Println("Could not Unmarshal Asset", marshErr)
		return request, marshErr
	}
	return request, nil

}

func SearchAsset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering SearchAsset " + args[0])
	var asset AssetDetails
	var err error
	var request SearchAssetRequest
	var resp SearchAssetResponse

	request, err = parseSearchAssetRequest(args[0])
	if err != nil {
		resp.ErrCode = INVALID_INPUT_ERROR_CODE
		resp.ErrMessage = INVALID_INPUT_ERROR_MSG
		return json.Marshal(resp)
	}

	asset, err = fetchAssetDetails(stub, request.AssetId)

	if err != nil {
		resp.ErrCode = NODATA_ERROR_CODE
		resp.ErrMessage = NODATA_ERROR_MSG
		return json.Marshal(resp)
	}

	return PrepareSearchAssetResponse(stub, asset)

}

/************** Asset Search Service Ends ************************/

/************** Carton Search Service Starts ************************/

type SearchCartonRequest struct {
	CallingEntityName string `json:"callingEntityName"`
	CartonId          string `json:"cartonId"`
}

type SearchCartonResponse struct {
	CartonId       string     `json:"cartonId"`
	PalletId       string     `json:"palletId"`
	ShipmentDetail []Shipment `json:"shipmentDetail"`
	ErrCode        string     `json:"errCode"`
	ErrMessage     string     `json:"errMessage"`
}

func parseSearchCartonRequest(requestParam string) (SearchCartonRequest, error) {
	var request SearchCartonRequest

	if marshErr := json.Unmarshal([]byte(requestParam), &request); marshErr != nil {
		fmt.Println("Could not Unmarshal Asset", marshErr)
		return request, marshErr
	}
	return request, nil

}

func parseCarton(stub shim.ChaincodeStubInterface, cartonId string) (Carton, error) {
	var carton Carton

	cartonBytes, err := stub.GetState(cartonId)
	if err != nil {
		return carton, err
	} else {
		if marshErr := json.Unmarshal(cartonBytes, &carton); marshErr != nil {
			fmt.Println("Could not Unmarshal Asset", marshErr)
			return carton, marshErr
		}
		return carton, nil
	}

}

func PrepareSearchCartontResponse(stub shim.ChaincodeStubInterface, carton Carton) ([]byte, error) {
	var resp SearchCartonResponse
	var shipmentArr []Shipment
	var tmpShipment Shipment
	var err error

	resp.CartonId = carton.CartonId
	resp.PalletId = carton.PalletId

	lenOfArray := len(carton.ShipmentIds)

	for i := 0; i < lenOfArray; i++ {
		tmpShipment, err = retrieveShipment(stub, carton.ShipmentIds[i])
		if err != nil {
			fmt.Println("Error while retrieveing the Shipment Details", err)
			return nil, err
		}
		shipmentArr = append(shipmentArr, tmpShipment)
	}

	resp.ShipmentDetail = shipmentArr
	return json.Marshal(resp)

}

func SearchCarton(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering SearchCarton " + args[0])
	var carton Carton
	var err error
	var request SearchCartonRequest
	var resp SearchCartonResponse

	request, err = parseSearchCartonRequest(args[0])
	if err != nil {
		resp.ErrCode = INVALID_INPUT_ERROR_CODE
		resp.ErrMessage = INVALID_INPUT_ERROR_MSG
		return json.Marshal(resp)
	}

	carton, err = parseCarton(stub, request.CartonId)

	if err != nil {
		resp.ErrCode = NODATA_ERROR_CODE
		resp.ErrMessage = NODATA_ERROR_MSG
		return json.Marshal(resp)
	}

	return PrepareSearchCartontResponse(stub, carton)

}

/************** Carton Search Service Ends ************************/

/************** Pallet Search Service Starts ************************/

type SearchPalletRequest struct {
	CallingEntityName string `json:"callingEntityName"`
	PalletId          string `json:"palletId"`
}

type SearchPalletResponse struct {
	PalletId       string     `json:"palletId"`
	ShipmentDetail []Shipment `json:"shipmentDetail"`
	ErrCode        string     `json:"errCode"`
	ErrMessage     string     `json:"errMessage"`
}

func parseSearchPalletRequest(requestParam string) (SearchPalletRequest, error) {
	var request SearchPalletRequest

	if marshErr := json.Unmarshal([]byte(requestParam), &request); marshErr != nil {
		fmt.Println("Could not Unmarshal Asset", marshErr)
		return request, marshErr
	}
	return request, nil

}

func parsePallet(stub shim.ChaincodeStubInterface, palletId string) (Pallet, error) {

	var pallet Pallet

	palletBytes, err := stub.GetState(palletId)
	if err != nil {
		return pallet, err
	} else {
		if marshErr := json.Unmarshal(palletBytes, &pallet); marshErr != nil {
			fmt.Println("Could not Unmarshal Asset", marshErr)
			return pallet, marshErr
		}
		return pallet, nil
	}

}

func PrepareSearchPalletResponse(stub shim.ChaincodeStubInterface, pallet Pallet) ([]byte, error) {
	var resp SearchPalletResponse
	var shipmentArr []Shipment
	var tmpShipment Shipment
	var err error

	resp.PalletId = pallet.PalletId

	lenOfArray := len(pallet.ShipmentIds)

	for i := 0; i < lenOfArray; i++ {
		tmpShipment, err = retrieveShipment(stub, pallet.ShipmentIds[i])
		if err != nil {
			fmt.Println("Error while retrieveing the Shipment Details", err)
			return nil, err
		}
		shipmentArr = append(shipmentArr, tmpShipment)
	}

	resp.ShipmentDetail = shipmentArr
	return json.Marshal(resp)

}

func SearchPallet(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering SearchPallet " + args[0])
	var pallet Pallet
	var err error
	var request SearchPalletRequest
	var resp SearchPalletResponse

	request, err = parseSearchPalletRequest(args[0])
	if err != nil {
		resp.ErrCode = INVALID_INPUT_ERROR_CODE
		resp.ErrMessage = INVALID_INPUT_ERROR_MSG
		return json.Marshal(resp)
	}

	pallet, err = parsePallet(stub, request.PalletId)

	if err != nil {
		resp.ErrCode = NODATA_ERROR_CODE
		resp.ErrMessage = NODATA_ERROR_MSG
		return json.Marshal(resp)
	}

	return PrepareSearchPalletResponse(stub, pallet)

}

/************** Pallet Search Service Ends ************************/
//END COMMNENTED BY ARSHAD AS THESE ARE NO MORE USED - KARTHIK PLEASE REVIEW AND PERMANNETLY DELETE IT IF NOT USED BY ANY OF YOU MODULE

/************** Date Search Service Starts ************************/

type SearchDateRequest struct {
	CallingEntityName string `json:"callingEntityName"`
	StartDate         string `json:"startDate"`
	EndDate           string `json:"endDate"`
}

type SearchDateResponse struct {
	ShipmentDetail []Shipment `json:"shipmentDetail"`
}

func parseAllShipmentDump() (AllShipmentDump, error) {
	var dump AllShipmentDump

	if marshErr := json.Unmarshal([]byte("ALL_SHIPMENT_DUMP"), &dump); marshErr != nil {
		fmt.Println("Could not Unmarshal Asset", marshErr)
		return dump, marshErr
	}
	return dump, nil

}

func SearchDateRange(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	//var shipmentDump AllShipmentDump
	var err error
	var shipmentArr []Shipment
	var tmpShipment Shipment
	var resp SearchDateResponse

	/*shipmentDump, err = parseAllShipmentDump()
	if err != nil {
		return nil, err
	}*/

	lenOfArray := len(args)

	for i := 0; i < lenOfArray; i++ {
		//	tmpShipment, err = retrieveShipment(stub, args[i])
		if err != nil {
			fmt.Println("Error while retrieveing the Shipment Details", err)
			return nil, err
		}
		shipmentArr = append(shipmentArr, tmpShipment)
	}
	resp.ShipmentDetail = shipmentArr

	return json.Marshal(resp)

}

/************** Date Search Service Ends ************************/