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


type TrailItems struct {
	TrailName			string           `json:"trailName"`
	TimeStamp			string           `json:"timeStamp"`
	TrailAddress		string           `json:"trailAddress"`
	TrailLat			float64           `json:"trailLat"`
	TrailLng			float64           `json:"trailLng"`
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
	ShipmentNumber       		string     	`json:"shipmentNumber"`
	Custodian             		string   	`json:"custodian"`
	CustodianTime   			string   	`json:"custodianTime"`
	Status             	  		string   	`json:"status"`

	SenderName					string   	`json:"senderName"`
	SenderAddress				string   	`json:"senderAddress"`
	ReceiverName				string   	`json:"receiverName"`
	ReceiverAddress				string   	`json:"receiverAddress"`

	ServiceType					string   	`json:"serviceType"`
	NoOfItems					string   	`json:"noOfItems"`
	TransportMode				string   	`json:"transportMode"`

	PalletsArr					[]string   	`json:"palletsArr"`

	History						[]TrailItems	`json:"history"`
}

type SearchAssetResponse struct {
	Data 						[]SearchAssetData		`json:"data"`
	ErrCode        				string     	`json:"errCode"`
	ErrMessage     				string     	`json:"errMessage"`
}

func prepareHistroyTrail(stub shim.ChaincodeStubInterface, histroy []CustodianHistoryDetail) ([]TrailItems) {
	var trailHistroyArr []TrailItems
	var thisClass ShipmentPageLoadService
	var err error

	lenOfArray := len(histroy)

	fmt.Println("===lenOfArray all histroy===", lenOfArray)

	for i := 0; i < lenOfArray; i++ {
		var tmpEntity Entity
		var trailHistroy TrailItems
		tmpEntity, err = thisClass.fetchEntities(stub, histroy[i].CustodianName)
		if err != nil {
			fmt.Println("Could not fetch Entities", err)
		}
		trailHistroy.TrailName = tmpEntity.EntityName
		trailHistroy.TimeStamp = histroy[i].CustodianTime
		trailHistroy.TrailAddress = tmpEntity.EntityAddress
		trailHistroy.TrailLat = tmpEntity.Latitude
		trailHistroy.TrailLng = tmpEntity.Longitude

		trailHistroyArr = append(trailHistroyArr, trailHistroy)
	}

	return trailHistroyArr
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
	
	if(asset.MshipmentNumber != "") {
		tmpShipment, err = fetchShipmentWayBillData(stub, asset.MshipmentNumber)
		
		if err != nil {
			resp.ErrCode = "ERR_DATA"
			resp.ErrMessage = "Unable to get Shipment MshipmentNumber"
			fmt.Println("Error while retrieveing the Shipment Details", err)
			return nil, err
		} else {
			respData.ShipmentNumber = asset.MshipmentNumber;
			respData.Custodian = tmpShipment.Custodian;
			respData.CustodianTime = tmpShipment.WayBillModifiedDate;
			respData.Status = tmpShipment.Status;

			respData.SenderName = tmpShipment.Consigner;
			respData.SenderAddress = tmpShipment.AddressOfConsigner;
			respData.ReceiverName = tmpShipment.Consignee;
			respData.ReceiverAddress = tmpShipment.AddressOfConsignee;

			respData.ServiceType = tmpShipment.ServiceType;
			respData.NoOfItems = tmpShipment.AssetsQuantity;
			respData.TransportMode = tmpShipment.VesselType;
			respData.PalletsArr = tmpShipment.PalletsSerialNumber;


			respData.History = prepareHistroyTrail(stub, tmpShipment.CustodianHistory)
			respDataArr = append(respDataArr, respData)
		}
	}
//-----------------------------------------------------------------------//
	/*var respData2 SearchAssetData
	respData2.AssetSerialNo = asset.AssetSerialNumber
	respData2.AssetModel = asset.AssetModel
	respData2.AssetType = asset.AssetType
	respData2.CartonSerialNumber = asset.CartonSerialNumber
	respData2.PalletSerialNumber = asset.PalletSerialNumber
	var tmpEWWayBill EWWayBill
	if(asset.EwWayBillNumber != "") {
		tmpEWWayBill, err = fetchEWWayBillData(stub, asset.EwWayBillNumber)
		
		if err != nil {
			resp.ErrCode = "ERR_DATA"
			resp.ErrMessage = "Unable to get Shipment DcWayBillNumber"
			fmt.Println("Error while retrieveing the Shipment Details", err)
			return nil, err
		} else {
			respData2.ShipmentNumber = asset.EwWayBillNumber;
			respData2.Custodian = tmpEWWayBill.Custodian;
			respData2.CustodianTime = tmpEWWayBill.EwWayBillModifiedDate;
			respData2.Status = tmpEWWayBill.Status;

			respData2.SenderName = tmpEWWayBill.Consigner;
			respData2.SenderAddress = tmpEWWayBill.AddressOfConsigner;
			respData2.ReceiverName = tmpEWWayBill.Consignee;
			respData2.ReceiverAddress = tmpEWWayBill.AddressOfConsignee;

			respData2.ServiceType = tmpEWWayBill.ServiceType;
			respData2.NoOfItems = "-";
			respData2.TransportMode = tmpEWWayBill.VesselType;
			respData2.PalletsArr = nil;


			respData2.History = prepareHistroyTrail(stub, tmpEWWayBill.CustodianHistory)
			respDataArr = append(respDataArr, respData2)
		}
	}*/
//-----------------------------------------------------------------------//
	var respData3 SearchAssetData
	respData3.AssetSerialNo = asset.AssetSerialNumber
	respData3.AssetModel = asset.AssetModel
	respData3.AssetType = asset.AssetType
	respData3.CartonSerialNumber = asset.CartonSerialNumber
	respData3.PalletSerialNumber = asset.PalletSerialNumber
	
	if(asset.DcShipmentNumber != "") {
		tmpShipment, err = fetchShipmentWayBillData(stub, asset.DcShipmentNumber)
		
		if err != nil {
			resp.ErrCode = "ERR_DATA"
			resp.ErrMessage = "Unable to get Shipment DcWayBillNumber"
			fmt.Println("Error while retrieveing the Shipment Details", err)
			return nil, err
		} else {
			respData3.ShipmentNumber = asset.DcShipmentNumber;
			respData3.Custodian = tmpShipment.Custodian;
			respData3.CustodianTime = tmpShipment.WayBillModifiedDate;
			respData3.Status = tmpShipment.Status;

			respData3.SenderName = tmpShipment.Consigner;
			respData3.SenderAddress = tmpShipment.AddressOfConsigner;
			respData3.ReceiverName = tmpShipment.Consignee;
			respData3.ReceiverAddress = tmpShipment.AddressOfConsignee;

			respData3.ServiceType = tmpShipment.ServiceType;
			respData3.NoOfItems = tmpShipment.AssetsQuantity;
			respData3.TransportMode = tmpShipment.VesselType;
			respData3.PalletsArr = tmpShipment.PalletsSerialNumber;


			respData3.History = prepareHistroyTrail(stub, tmpShipment.CustodianHistory)
			respDataArr = append(respDataArr, respData3)
		}
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
	CartonId           string `json:"cartonId"`
}

type SearchCartonData struct {
	CartonSerialNumber       	string     	`json:"cartonSerialNumber"`
	PalletSerialNumber       	string     	`json:"palletSerialNumber"`
	ShipmentNumber       		string     	`json:"shipmentNumber"`
	Custodian             		string   	`json:"custodian"`
	CustodianTime   			string   	`json:"custodianTime"`
	Status             	  		string   	`json:"status"`

	SenderName					string   	`json:"senderName"`
	SenderAddress				string   	`json:"senderAddress"`
	ReceiverName				string   	`json:"receiverName"`
	ReceiverAddress				string   	`json:"receiverAddress"`

	ServiceType					string   	`json:"serviceType"`
	NoOfItems					string   	`json:"noOfItems"`
	TransportMode				string   	`json:"transportMode"`

	PalletsArr					[]string   	`json:"palletsArr"`

	History						[]TrailItems	`json:"history"`
}

type SearchCartonResponse struct {
	Data 						[]SearchCartonData		`json:"data"`
	ErrCode        				string     	`json:"errCode"`
	ErrMessage     				string     	`json:"errMessage"`
}

func PrepareSearchCartonResponse(stub shim.ChaincodeStubInterface, carton CartonDetails) ([]byte, error) {
	fmt.Println("PrepareSearchCartonResponse ", carton)
	var resp SearchCartonResponse
	var tmpShipment ShipmentWayBill
	var err error
	var respDataArr []SearchCartonData
	

	
	var respData SearchCartonData
	respData.CartonSerialNumber = carton.CartonSerialNumber
	respData.PalletSerialNumber = carton.PalletSerialNumber
	
	if(carton.MshipmentNumber != "") {
		tmpShipment, err = fetchShipmentWayBillData(stub, carton.MshipmentNumber)
		
		if err != nil {
			resp.ErrCode = "ERR_DATA"
			resp.ErrMessage = "Unable to get Shipment MshipmentNumber"
			fmt.Println("Error while retrieveing the Shipment Details", err)
			return nil, err
		} else {
			respData.ShipmentNumber = carton.MshipmentNumber;
			respData.Custodian = tmpShipment.Custodian;
			respData.CustodianTime = tmpShipment.WayBillModifiedDate;
			respData.Status = tmpShipment.Status;

			respData.SenderName = tmpShipment.Consigner;
			respData.SenderAddress = tmpShipment.AddressOfConsigner;
			respData.ReceiverName = tmpShipment.Consignee;
			respData.ReceiverAddress = tmpShipment.AddressOfConsignee;

			respData.ServiceType = tmpShipment.ServiceType;
			respData.NoOfItems = tmpShipment.AssetsQuantity;
			respData.TransportMode = tmpShipment.VesselType;
			respData.PalletsArr = tmpShipment.PalletsSerialNumber;


			respData.History = prepareHistroyTrail(stub, tmpShipment.CustodianHistory)
			respDataArr = append(respDataArr, respData)
		}
	}
//-----------------------------------------------------------------------//
	/*var respData2 SearchCartonData
	respData2.CartonSerialNumber = carton.CartonSerialNumber
	respData2.PalletSerialNumber = carton.PalletSerialNumber
	var tmpEWWayBill EWWayBill
	if(carton.EwWayBillNumber != "") {
		tmpEWWayBill, err = fetchShipmentWayBillData(stub, carton.EwWayBillNumber)
		
		if err != nil {
			resp.ErrCode = "ERR_DATA"
			resp.ErrMessage = "Unable to get Shipment DcWayBillNumber"
			fmt.Println("Error while retrieveing the Shipment Details", err)
			return nil, err
		} else {
			respData2.ShipmentNumber = carton.DcWayBillNumber;
			respData2.Custodian = tmpEWWayBill.Custodian;
			respData2.CustodianTime = tmpEWWayBill.EwWayBillModifiedDate;
			respData2.Status = tmpEWWayBill.Status;

			respData2.SenderName = tmpEWWayBill.Consigner;
			respData2.SenderAddress = tmpEWWayBill.AddressOfConsigner;
			respData2.ReceiverName = tmpEWWayBill.Consignee;
			respData2.ReceiverAddress = tmpEWWayBill.AddressOfConsignee;

			respData2.ServiceType = tmpEWWayBill.ServiceType;
			respData2.NoOfItems = '-';
			respData2.TransportMode = tmpEWWayBill.VesselType;
			respData2.PalletsArr = nil;


			respData2.History = prepareHistroyTrail(stub, tmpEWWayBill.CustodianHistory)
			respDataArr = append(respDataArr, respData2)
		}
	}*/
//-----------------------------------------------------------------------//
	var respData3 SearchCartonData
	respData3.CartonSerialNumber = carton.CartonSerialNumber
	respData3.PalletSerialNumber = carton.PalletSerialNumber
	
	if(carton.DcShipmentNumber != "") {
		tmpShipment, err = fetchShipmentWayBillData(stub, carton.DcShipmentNumber)
		
		if err != nil {
			resp.ErrCode = "ERR_DATA"
			resp.ErrMessage = "Unable to get Shipment DcWayBillNumber"
			fmt.Println("Error while retrieveing the Shipment Details", err)
			return nil, err
		} else {
			respData3.ShipmentNumber = carton.DcShipmentNumber;
			respData3.Custodian = tmpShipment.Custodian;
			respData3.CustodianTime = tmpShipment.WayBillModifiedDate;
			respData3.Status = tmpShipment.Status;

			respData3.SenderName = tmpShipment.Consigner;
			respData3.SenderAddress = tmpShipment.AddressOfConsigner;
			respData3.ReceiverName = tmpShipment.Consignee;
			respData3.ReceiverAddress = tmpShipment.AddressOfConsignee;

			respData3.ServiceType = tmpShipment.ServiceType;
			respData3.NoOfItems = tmpShipment.AssetsQuantity;
			respData3.TransportMode = tmpShipment.VesselType;
			respData3.PalletsArr = tmpShipment.PalletsSerialNumber;


			respData3.History = prepareHistroyTrail(stub, tmpShipment.CustodianHistory)
			respDataArr = append(respDataArr, respData3)
		}
	}
	
	resp.Data = respDataArr
	return json.Marshal(resp)

}

func parseSearchCartonRequest(requestParam string) (SearchCartonRequest, error) {
	var request SearchCartonRequest

	if marshErr := json.Unmarshal([]byte(requestParam), &request); marshErr != nil {
		fmt.Println("Could not Unmarshal Carton", marshErr)
		return request, marshErr
	}
	return request, nil

}

func SearchCarton(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering SearchCarton " + args[0])
	var carton CartonDetails
	var err error
	var request SearchCartonRequest
	var resp SearchCartonResponse

	request, err = parseSearchCartonRequest(args[0])
	if err != nil {
		resp.ErrCode = INVALID_INPUT_ERROR_CODE
		resp.ErrMessage = INVALID_INPUT_ERROR_MSG
		return json.Marshal(resp)
	}

	carton, err = fetchCartonDetails(stub, request.CartonId)

	if err != nil {
		resp.ErrCode = NODATA_ERROR_CODE
		resp.ErrMessage = NODATA_ERROR_MSG
		return json.Marshal(resp)
	}

	return PrepareSearchCartonResponse(stub, carton)

}

/************** Carton Search Service Ends ************************/

/************** Pallet Search Service Starts ************************/


type SearchPalletRequest struct {
	CallingEntityName string `json:"callingEntityName"`
	PalletId           string `json:"palletId"`
}

type SearchPalletData struct {
	PalletSerialNumber       	string     	`json:"palletSerialNumber"`
	ShipmentNumber       		string     	`json:"shipmentNumber"`
	Custodian             		string   	`json:"custodian"`
	CustodianTime   			string   	`json:"custodianTime"`
	Status             	  		string   	`json:"status"`

	SenderName					string   	`json:"senderName"`
	SenderAddress				string   	`json:"senderAddress"`
	ReceiverName				string   	`json:"receiverName"`
	ReceiverAddress				string   	`json:"receiverAddress"`

	ServiceType					string   	`json:"serviceType"`
	NoOfItems					string   	`json:"noOfItems"`
	TransportMode				string   	`json:"transportMode"`

	PalletsArr					[]string   	`json:"palletsArr"`

	History						[]TrailItems	`json:"history"`
}

type SearchPalletResponse struct {
	Data 						[]SearchPalletData		`json:"data"`
	ErrCode        				string     	`json:"errCode"`
	ErrMessage     				string     	`json:"errMessage"`
}

func PrepareSearchPalletResponse(stub shim.ChaincodeStubInterface, pallet PalletDetails) ([]byte, error) {
	fmt.Println("PrepareSearchPalletResponse ", pallet)
	var resp SearchPalletResponse
	var tmpShipment ShipmentWayBill
	var err error
	var respDataArr []SearchPalletData
	

	
	var respData SearchPalletData
	respData.PalletSerialNumber = pallet.PalletSerialNumber

	if(pallet.MshipmentNumber != "") {
	tmpShipment, err = fetchShipmentWayBillData(stub, pallet.MshipmentNumber)
	
	
		if err != nil {
			resp.ErrCode = "ERR_DATA"
			resp.ErrMessage = "Unable to get Shipment MshipmentNumber"
			fmt.Println("Error while retrieveing the Shipment Details", err)
			return nil, err
		} else {
			respData.ShipmentNumber = pallet.MshipmentNumber;
			respData.Custodian = tmpShipment.Custodian;
			respData.CustodianTime = tmpShipment.WayBillModifiedDate;
			respData.Status = tmpShipment.Status;

			respData.SenderName = tmpShipment.Consigner;
			respData.SenderAddress = tmpShipment.AddressOfConsigner;
			respData.ReceiverName = tmpShipment.Consignee;
			respData.ReceiverAddress = tmpShipment.AddressOfConsignee;

			respData.ServiceType = tmpShipment.ServiceType;
			respData.NoOfItems = tmpShipment.AssetsQuantity;
			respData.TransportMode = tmpShipment.VesselType;
			respData.PalletsArr = tmpShipment.PalletsSerialNumber;


			respData.History = prepareHistroyTrail(stub, tmpShipment.CustodianHistory)
			respDataArr = append(respDataArr, respData)
		}
	}
//-----------------------------------------------------------------------//
	/*var respData2 SearchPalletData
	respData2.PalletSerialNumber = pallet.PalletSerialNumber
	var tmpEWWayBill EWWayBill

	if(pallet.EwWayBillNumber != "") {
		tmpEWWayBill, err = fetchShipmentWayBillData(stub, pallet.EwWayBillNumber)
		
		if err != nil {
			resp.ErrCode = "ERR_DATA"
			resp.ErrMessage = "Unable to get Shipment DcWayBillNumber"
			fmt.Println("Error while retrieveing the Shipment Details", err)
			return nil, err
		} else {
			respData2.ShipmentNumber = pallet.EwWayBillNumber;
			respData2.Custodian = tmpEWWayBill.Custodian;
			respData2.CustodianTime = tmpEWWayBill.EwWayBillModifiedDate;
			respData2.Status = tmpEWWayBill.Status;

			respData2.SenderName = tmpEWWayBill.Consigner;
			respData2.SenderAddress = tmpEWWayBill.AddressOfConsigner;
			respData2.ReceiverName = tmpEWWayBill.Consignee;
			respData2.ReceiverAddress = tmpEWWayBill.AddressOfConsignee;

			respData2.ServiceType = tmpEWWayBill.ServiceType;
			respData2.NoOfItems = '-';
			respData2.TransportMode = tmpEWWayBill.VesselType;
			respData2.PalletsArr = nil;


			respData2.History = prepareHistroyTrail(stub, tmpEWWayBill.CustodianHistory)
			respDataArr = append(respDataArr, respData2)
		}
	}*/
//-----------------------------------------------------------------------//
	var respData3 SearchPalletData
	respData3.PalletSerialNumber = pallet.PalletSerialNumber
	
	if(pallet.DcShipmentNumber != "") {
		tmpShipment, err = fetchShipmentWayBillData(stub, pallet.DcShipmentNumber)
		
		if err != nil {
			resp.ErrCode = "ERR_DATA"
			resp.ErrMessage = "Unable to get Shipment DcWayBillNumber"
			fmt.Println("Error while retrieveing the Shipment Details", err)
			return nil, err
		} else {
			respData3.ShipmentNumber = pallet.DcShipmentNumber;
			respData3.Custodian = tmpShipment.Custodian;
			respData3.CustodianTime = tmpShipment.WayBillModifiedDate;
			respData3.Status = tmpShipment.Status;

			respData3.SenderName = tmpShipment.Consigner;
			respData3.SenderAddress = tmpShipment.AddressOfConsigner;
			respData3.ReceiverName = tmpShipment.Consignee;
			respData3.ReceiverAddress = tmpShipment.AddressOfConsignee;

			respData3.ServiceType = tmpShipment.ServiceType;
			respData3.NoOfItems = tmpShipment.AssetsQuantity;
			respData3.TransportMode = tmpShipment.VesselType;
			respData3.PalletsArr = tmpShipment.PalletsSerialNumber;


			respData3.History = prepareHistroyTrail(stub, tmpShipment.CustodianHistory)
			respDataArr = append(respDataArr, respData3)
		}
	}
	
	resp.Data = respDataArr
	return json.Marshal(resp)

}

func parseSearchPalletRequest(requestParam string) (SearchPalletRequest, error) {
	var request SearchPalletRequest

	if marshErr := json.Unmarshal([]byte(requestParam), &request); marshErr != nil {
		fmt.Println("Could not Unmarshal Asset", marshErr)
		return request, marshErr
	}
	return request, nil

}

func SearchPallet(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering SearchPallet " + args[0])
	var pallet PalletDetails
	var err error
	var request SearchPalletRequest
	var resp SearchPalletResponse

	request, err = parseSearchPalletRequest(args[0])
	if err != nil {
		resp.ErrCode = INVALID_INPUT_ERROR_CODE
		resp.ErrMessage = INVALID_INPUT_ERROR_MSG
		return json.Marshal(resp)
	}

	pallet, err = fetchPalletDetails(stub, request.PalletId)

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

type SearchDateResponseData struct {
	AssetSerialNo  				string     	`json:"assetSerialNo"`
	AssetModel      			string     	`json:"assetModel"`
	AssetType  					string     	`json:"assetType"`
	CartonSerialNumber       	string     	`json:"cartonSerialNumber"`
	PalletSerialNumber       	string     	`json:"palletSerialNumber"`
	ShipmentNumber       		string     	`json:"shipmentNumber"`
	Custodian             		string   	`json:"custodian"`
	CustodianTime   			string   	`json:"custodianTime"`
	Status             	  		string   	`json:"status"`

	SenderName					string   	`json:"senderName"`
	SenderAddress				string   	`json:"senderAddress"`
	ReceiverName				string   	`json:"receiverName"`
	ReceiverAddress				string   	`json:"receiverAddress"`

	ServiceType					string   	`json:"serviceType"`
	NoOfItems					string   	`json:"noOfItems"`
	TransportMode				string   	`json:"transportMode"`

	PalletsArr					[]string   	`json:"palletsArr"`

	History						[]TrailItems	`json:"history"`
}

type SearchDateResponse struct {
	Data 						[]SearchDateResponseData		`json:"data"`
	ErrCode        				string     	`json:"errCode"`
	ErrMessage     				string     	`json:"errMessage"`
}

type AllAssetData struct {
	Index 		[]string		`json:"index"`
}

func fetchAllAssetData(stub shim.ChaincodeStubInterface) ([]string, error) {
	var allAssetData AllAssetData
	var dummydata []string

	indexByte, err := stub.GetState("ALL_ASSET_INDEX")
	if err != nil {
		fmt.Println("Could not retrive  Shipment WayBill ", err)
		return dummydata, err
	}

	json.Unmarshal(indexByte, &allAssetData)

	fmt.Println("======================allAssetData data-->")
	fmt.Println(string(indexByte))
	fmt.Println(allAssetData)
	fmt.Println("======================")

	return allAssetData.Index, nil

}

func PrepareSearchDateRangeResponse(stub shim.ChaincodeStubInterface, assetIds []string, request SearchDateRequest) ([]byte, error) {
	var resp SearchDateResponse
	var err error
	var respDataArr []SearchDateResponseData
	var asset AssetDetails
	var tmpShipment ShipmentWayBill

	lenOfArray := len(assetIds)
	for i := 0; i < lenOfArray; i++ {
		fmt.Println("======================assetIds[i] data-->")
		fmt.Println(assetIds[i])
		fmt.Println("======================")
		asset, err = fetchAssetDetails(stub, assetIds[i])

		var respData SearchDateResponseData
		respData.AssetSerialNo = asset.AssetSerialNumber
		respData.AssetModel = asset.AssetModel
		respData.AssetType = asset.AssetType
		respData.CartonSerialNumber = asset.CartonSerialNumber
		respData.PalletSerialNumber = asset.PalletSerialNumber
		
		fmt.Println("======================asset data-->")
		fmt.Println(asset)
		fmt.Println("======================")

		if(asset.MshipmentNumber != "") {
			tmpShipment, err = fetchShipmentWayBillData(stub, asset.MshipmentNumber)
			
			if err != nil {
				resp.ErrCode = "ERR_DATA"
				resp.ErrMessage = "Unable to get Shipment MshipmentNumber"
				fmt.Println("Error while retrieveing the Shipment Details", err)
				return nil, err
			} else {
				respData.ShipmentNumber = asset.MshipmentNumber;
				respData.Custodian = tmpShipment.Custodian;
				respData.CustodianTime = tmpShipment.WayBillModifiedDate;
				respData.Status = tmpShipment.Status;

				respData.SenderName = tmpShipment.Consigner;
				respData.SenderAddress = tmpShipment.AddressOfConsigner;
				respData.ReceiverName = tmpShipment.Consignee;
				respData.ReceiverAddress = tmpShipment.AddressOfConsignee;

				respData.ServiceType = tmpShipment.ServiceType;
				respData.NoOfItems = tmpShipment.AssetsQuantity;
				respData.TransportMode = tmpShipment.VesselType;
				respData.PalletsArr = tmpShipment.PalletsSerialNumber;


				respData.History = prepareHistroyTrail(stub, tmpShipment.CustodianHistory)
				if(compareDates(request.StartDate, request.EndDate, tmpShipment.ShipmentCreationDate)) {
					respDataArr = append(respDataArr, respData)
				}else {
					fmt.Println("tmpShipment.ShipmentCreationDate is not in range", tmpShipment.ShipmentCreationDate)
				}

				
			}
		}

		var respData3 SearchDateResponseData
		respData3.AssetSerialNo = asset.AssetSerialNumber
		respData3.AssetModel = asset.AssetModel
		respData3.AssetType = asset.AssetType
		respData3.CartonSerialNumber = asset.CartonSerialNumber
		respData3.PalletSerialNumber = asset.PalletSerialNumber
		
		if(asset.DcShipmentNumber != "") {
			tmpShipment, err = fetchShipmentWayBillData(stub, asset.DcShipmentNumber)
			
			if err != nil {
				resp.ErrCode = "ERR_DATA"
				resp.ErrMessage = "Unable to get Shipment DcWayBillNumber"
				fmt.Println("Error while retrieveing the Shipment Details", err)
				return nil, err
			} else {
				respData3.ShipmentNumber = asset.DcShipmentNumber;
				respData3.Custodian = tmpShipment.Custodian;
				respData3.CustodianTime = tmpShipment.WayBillModifiedDate;
				respData3.Status = tmpShipment.Status;

				respData3.SenderName = tmpShipment.Consigner;
				respData3.SenderAddress = tmpShipment.AddressOfConsigner;
				respData3.ReceiverName = tmpShipment.Consignee;
				respData3.ReceiverAddress = tmpShipment.AddressOfConsignee;

				respData3.ServiceType = tmpShipment.ServiceType;
				respData3.NoOfItems = tmpShipment.AssetsQuantity;
				respData3.TransportMode = tmpShipment.VesselType;
				respData3.PalletsArr = tmpShipment.PalletsSerialNumber;


				respData3.History = prepareHistroyTrail(stub, tmpShipment.CustodianHistory)
				if(compareDates(request.StartDate, request.EndDate, tmpShipment.ShipmentCreationDate)) {
					respDataArr = append(respDataArr, respData3)
				}else {
					fmt.Println("tmpShipment.ShipmentCreationDate is not in range", tmpShipment.ShipmentCreationDate)
				}
				
			}
		}
	}
	resp.Data = respDataArr
	return json.Marshal(resp)

}

func parseSearchDateRequest(requestParam string) (SearchDateRequest) {
	var request SearchDateRequest

	if marshErr := json.Unmarshal([]byte(requestParam), &request); marshErr != nil {
		fmt.Println("Could not Unmarshal Asset", marshErr)
	}
	return request

}


func SearchDateRange(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	//var shipmentDump AllShipmentDump
	var err error
	var resp SearchDateResponse
	var assetIDArr []string
	var request SearchDateRequest

	request = parseSearchDateRequest(args[0])
	
	assetIDArr,err =fetchAllAssetData(stub)
	fmt.Println("======================allAssetData data-->")
	fmt.Println(assetIDArr)
	fmt.Println(err)
	fmt.Println("======================")
	if(err == nil) {
		
		return PrepareSearchDateRangeResponse(stub, assetIDArr,request )

	}

	return json.Marshal(resp)

}

/************** Date Search Service Ends ************************/