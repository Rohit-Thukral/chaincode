package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type ShipmentPageLoadService struct {
}

type ShipmentPageLoadRequest struct {
	CallingEntityName string `json:"callingEntityName"`
}

type ConsignerShipmentPageLoadResponse struct {
	ConsignerId        string `json:"consignerId"`
	ConsignerName      string `json:"consignerName"`
	ConsignerType      string `json:"consignerType"`
	ConsignerAddress   string `json:"consignerAddress"`
	ConsignerRegNumber string `json:"consignerRegNumber"`
	ConsignerCountry   string `json:"consignerCountry"`
}

type ConsigneeShipmentPageLoadResponse struct {
	ConsigneeId        string `json:"consigneeId"`
	ConsigneeName      string `json:"consigneeName"`
	ConsigneeAddress   string `json:"consigneeAddress"`
	ConsigneeCountry   string `json:"consigneeCountry"`
	ConsigneeRegNumber string `json:"consigneeRegNumber"`
}
type RetailerShipmentPageLoadResponse struct {
	RetailerId        string `json:"RetailerId"`
	RetailerName      string `json:"RetailerName"`
	RetailerAddress   string `json:"RetailerAddress"`
	RetailerCountry   string `json:"RetailerCountry"`
	RetailerRegNumber string `json:"RetailerRegNumber"`
}

type CarrierResponse struct {
	CarrierID string `json:"CarrierID"`
	Name      string `json:"Name"`
}

type ShipmentPageLoadResponse struct {
	CallingEntityName  string                              `json:"callingEntityName"`
	ConsignerId        string                              `json:"consignerId"`
	ConsignerName      string                              `json:"consignerName"`
	ConsignerType      string                              `json:"consignerType"`
	ConsignerAddress   string                              `json:"consignerAddress"`
	ConsignerRegNumber string                              `json:"consignerRegNumber"`
	ConsignerCountry   string                              `json:"consignerCountry"`
	Consignee          []ConsigneeShipmentPageLoadResponse `json:"consignee"`
	Carrier            []CarrierResponse                   `json:"carrier"`
	ModelNames         []string                            `json:"modelNames"`
	WaybillIds         []EntityWayBillMappingDetail        `json:"waybillIds"`
	Retailer           []RetailerShipmentPageLoadResponse  `json:"retailer"`
}

type CountryEntityMappingRequest struct {
	CountryFrom string `json:"countryFrom"`
}
type CountryEntityMappingResponse struct {
	WareHouseList []ConsigneeShipmentPageLoadResponse `json:"wareHouseList"`
}

func (t *ShipmentPageLoadService) GetCountryWarehouse(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering GetCountryWarehouse New " + args[0])
	var thisClass ShipmentPageLoadService

	var err error
	var consigneeArr CountryEntityMappingResponse
	var allEntities AllEntities
	request := CountryEntityMappingRequest{}

	json.Unmarshal([]byte(args[0]), &request)

	allEntities, err = thisClass.fetchAllEntities(stub)
	if err == nil {
		lenOfArray := len(allEntities.EntityArr)
		fmt.Println("===lenOfArray all entities===", lenOfArray)

		for i := 0; i < lenOfArray; i++ {
			var tmpConsigneeResponse ConsigneeShipmentPageLoadResponse
			var tmpEntity Entity

			tmpEntity, err = thisClass.fetchEntities(stub, allEntities.EntityArr[i])

			if err == nil {
				if tmpEntity.EntityCountry == request.CountryFrom && tmpEntity.EntityType == "Warehouse" {
					tmpConsigneeResponse.ConsigneeId = tmpEntity.EntityId
					tmpConsigneeResponse.ConsigneeName = tmpEntity.EntityName
					tmpConsigneeResponse.ConsigneeAddress = tmpEntity.EntityAddress
					tmpConsigneeResponse.ConsigneeCountry = tmpEntity.EntityCountry
					tmpConsigneeResponse.ConsigneeRegNumber = tmpEntity.EntityRegNumber
					consigneeArr.WareHouseList = append(consigneeArr.WareHouseList, tmpConsigneeResponse)
				}
			}
		}
	} else {
		fmt.Println("Error while fetching workflow data", err)
		return json.Marshal(consigneeArr)
	}

	fmt.Println("consigneeArr : ======================")
	fmt.Println(consigneeArr)
	fmt.Println("consigneeArr : ======================")

	fmt.Println("Exiting GetCountryWarehouse ")
	datatoreturn, _ := json.Marshal(consigneeArr)
	return []byte(datatoreturn), nil

}

func (t *ShipmentPageLoadService) ShipmentPageLoad(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering ShipmentPageLoad New " + args[0])
	var thisClass ShipmentPageLoadService

	var err error

	var tmpEntity Entity

	var consignerDetails ConsignerShipmentPageLoadResponse
	var consigneeArr []ConsigneeShipmentPageLoadResponse
	var carrier []CarrierResponse
	var response ShipmentPageLoadResponse
	var assetModelDetails AssetModelDetails
	var waybillIds EntityWayBillMapping
	var entityWayBillMappingDetail []EntityWayBillMappingDetail
	request := thisClass.parseShipmentPageLoadRequest(args[0])

	tmpEntity, err = thisClass.fetchEntities(stub, request.CallingEntityName)
	if err == nil {
		consignerDetails.ConsignerId = tmpEntity.EntityId
		consignerDetails.ConsignerName = tmpEntity.EntityName
		consignerDetails.ConsignerType = tmpEntity.EntityType
		consignerDetails.ConsignerAddress = tmpEntity.EntityAddress
		consignerDetails.ConsignerRegNumber = tmpEntity.EntityRegNumber
		consignerDetails.ConsignerCountry = tmpEntity.EntityCountry

	}

	assetModelDetails, err = thisClass.fetchAssetModelName(stub)
	fmt.Println("consignerDetails.Id == " + consignerDetails.ConsignerId)
	if consignerDetails.ConsignerId != "" {
		consigneeArr, carrier, err = thisClass.fetchCorrespondingConsignees(stub, consignerDetails)

		var werr error
		if tmpEntity.EntityType == "Warehouse" {
			waybillIds, werr = fetchEntityWayBillMappingData(stub, tmpEntity.EntityId)
			fmt.Println("After getting entitywaybill mapping-----", waybillIds)
			if werr != nil {
				fmt.Println("Error while fetching waybill id based on entity ", werr)
			}
			entityWayBillMappingDetail = waybillIds.WayBillsNumber
		}
		response.Retailer, _ = thisClass.fetchCorrespondingRetailer(stub, consignerDetails)
		response.CallingEntityName = request.CallingEntityName

		response.ConsignerId = consignerDetails.ConsignerId
		response.ConsignerName = consignerDetails.ConsignerName
		response.ConsignerType = consignerDetails.ConsignerType
		response.ConsignerAddress = consignerDetails.ConsignerAddress
		response.ConsignerRegNumber = consignerDetails.ConsignerRegNumber
		response.ConsignerCountry = consignerDetails.ConsignerCountry

		response.Consignee = consigneeArr
		response.Carrier = carrier
		response.WaybillIds = entityWayBillMappingDetail
		response.ModelNames = assetModelDetails.ModelNames
	}

	fmt.Println("response : ======================")
	fmt.Println(response)
	fmt.Println("response : ======================")
	fmt.Println("Exiting ShipmentPageLoad ")

	return json.Marshal(response)

	//return nil,nil

}

func (t *ShipmentPageLoadService) fetchCorrespondingConsignees(stub shim.ChaincodeStubInterface, consignerDetails ConsignerShipmentPageLoadResponse) ([]ConsigneeShipmentPageLoadResponse, []CarrierResponse, error) {
	fmt.Println("Entering fetchCorrespondingConsignees consignerDetails : ")
	fmt.Println("===consignerDetails===", consignerDetails)

	var err error
	var thisClass ShipmentPageLoadService

	var consigneeArr []ConsigneeShipmentPageLoadResponse

	var carrier []CarrierResponse
	var allEntities AllEntities

	allEntities, err = thisClass.fetchAllEntities(stub)
	if err == nil {
		lenOfArray := len(allEntities.EntityArr)
		fmt.Println("===lenOfArray all entities===", lenOfArray)

		for i := 0; i < lenOfArray; i++ {
			var tmpConsigneeResponse ConsigneeShipmentPageLoadResponse
			var tmpEntity Entity

			tmpEntity, err = thisClass.fetchEntities(stub, allEntities.EntityArr[i])
			fmt.Println("===tmpEntity===", tmpEntity)
			fmt.Println("===consignerDetails.ConsignerType===", consignerDetails.ConsignerType)
			fmt.Println("===tmpEntity.EntityType===", tmpEntity.EntityType)
			fmt.Println("===consignerDetails.ConsignerCountry===", consignerDetails.ConsignerCountry)
			fmt.Println("===tmpEntity.EntityCountry===", tmpEntity.EntityCountry)
			fmt.Println("===consignerDetails.ConsignerId===", consignerDetails.ConsignerId)
			fmt.Println("===tmpEntity.EntityId ===", tmpEntity.EntityId)

			if err == nil {
				if tmpEntity.EntityId != consignerDetails.ConsignerId && ((consignerDetails.ConsignerType == "Manufacturer" && tmpEntity.EntityType == "DC" && consignerDetails.ConsignerCountry == tmpEntity.EntityCountry) || (consignerDetails.ConsignerType == "DC" && tmpEntity.EntityType == "DC" && consignerDetails.ConsignerCountry != tmpEntity.EntityCountry) || (consignerDetails.ConsignerType == "Warehouse" && tmpEntity.EntityType == "Warehouse" && consignerDetails.ConsignerCountry != tmpEntity.EntityCountry) || (consignerDetails.ConsignerType == "DC" && tmpEntity.EntityType == "Retailer" && consignerDetails.ConsignerCountry == tmpEntity.EntityCountry)) {
					tmpConsigneeResponse.ConsigneeId = tmpEntity.EntityId
					tmpConsigneeResponse.ConsigneeName = tmpEntity.EntityName
					tmpConsigneeResponse.ConsigneeAddress = tmpEntity.EntityAddress
					tmpConsigneeResponse.ConsigneeCountry = tmpEntity.EntityCountry
					tmpConsigneeResponse.ConsigneeRegNumber = tmpEntity.EntityRegNumber
					consigneeArr = append(consigneeArr, tmpConsigneeResponse)

				}

				if consignerDetails.ConsignerCountry == tmpEntity.EntityCountry && tmpEntity.EntityType == "ThirdPartyLogistic" {
					var tmpCarrier CarrierResponse
					tmpCarrier.CarrierID = tmpEntity.EntityId
					tmpCarrier.Name = tmpEntity.EntityName
					carrier = append(carrier, tmpCarrier)
				}
			}

		}
	} else {
		fmt.Println("Error while fetching workflow data", err)
		return consigneeArr, carrier, err
	}

	fmt.Println("consigneeArr : ======================")
	fmt.Println(consigneeArr)
	fmt.Println("consigneeArr : ======================")
	fmt.Println("carrier : ======================")
	fmt.Println(carrier)
	fmt.Println("carrier : ======================")
	fmt.Println("Exiting fetchCorrespondingConsignees ")

	return consigneeArr, carrier, nil

}
func (t *ShipmentPageLoadService) fetchCorrespondingRetailer(stub shim.ChaincodeStubInterface, consignerDetails ConsignerShipmentPageLoadResponse) ([]RetailerShipmentPageLoadResponse, error) {
	fmt.Println("Entering fetchCorrespondingRetailer consignerDetails : ")
	fmt.Println("===consignerDetails===", consignerDetails)

	var err error
	var thisClass ShipmentPageLoadService

	var retailerArr []RetailerShipmentPageLoadResponse
	var allEntities AllEntities

	allEntities, err = thisClass.fetchAllEntities(stub)
	if err == nil {
		lenOfArray := len(allEntities.EntityArr)
		fmt.Println("===lenOfArray all entities===", lenOfArray)

		for i := 0; i < lenOfArray; i++ {
			var tmpRetailerResponse RetailerShipmentPageLoadResponse
			var tmpEntity Entity

			tmpEntity, err = thisClass.fetchEntities(stub, allEntities.EntityArr[i])

			if err == nil {
				if tmpEntity.EntityId != consignerDetails.ConsignerId && (consignerDetails.ConsignerType == "DC" && tmpEntity.EntityType == "Retailer" && consignerDetails.ConsignerCountry == tmpEntity.EntityCountry) {
					tmpRetailerResponse.RetailerId = tmpEntity.EntityId
					tmpRetailerResponse.RetailerName = tmpEntity.EntityName
					tmpRetailerResponse.RetailerAddress = tmpEntity.EntityAddress
					tmpRetailerResponse.RetailerCountry = tmpEntity.EntityCountry
					tmpRetailerResponse.RetailerRegNumber = tmpEntity.EntityRegNumber
					retailerArr = append(retailerArr, tmpRetailerResponse)

				}
			}

		}
	} else {
		fmt.Println("Error while fetching workflow data", err)
		return retailerArr, err
	}

	return retailerArr, nil

}

func (t *ShipmentPageLoadService) fetchAssetModelName(stub shim.ChaincodeStubInterface) (AssetModelDetails, error) {
	fmt.Println("Entering fetchAssetModelName ")
	var modelnames AssetModelDetails

	indexByte, err := stub.GetState("ASSET_MODEL_NAMES")
	if err != nil {
		fmt.Println("Could not retrive Shipment Index", err)
		return modelnames, err
	}

	if marshErr := json.Unmarshal(indexByte, &modelnames); marshErr != nil {
		fmt.Println("Could not save Shipment to ledger", marshErr)
		return modelnames, marshErr
	}
	fmt.Println(modelnames)
	fmt.Println("Exiting fetchAssetModelName ")
	return modelnames, nil

}

func (t *ShipmentPageLoadService) fetchWorkflows(stub shim.ChaincodeStubInterface) (AllWorkflows, error) {
	fmt.Println("Entering fetchWorkflows ")
	var workflows AllWorkflows

	indexByte, err := stub.GetState("ALL_WORKFLOWS")
	if err != nil {
		fmt.Println("Could not retrive Shipment Index", err)
		return workflows, err
	}

	if marshErr := json.Unmarshal(indexByte, &workflows); marshErr != nil {
		fmt.Println("Could not save Shipment to ledger", marshErr)
		return workflows, marshErr
	}
	fmt.Println(workflows)
	fmt.Println("Exiting fetchWorkflows ")
	return workflows, nil

}

func (t *ShipmentPageLoadService) fetchEntities(stub shim.ChaincodeStubInterface, entityID string) (Entity, error) {
	fmt.Println("Entering fetchEntities " + entityID)
	var entities Entity

	indexByte, err := stub.GetState(entityID)
	if err != nil {
		fmt.Println("Could not retrive Shipment Index", err)
		return entities, err
	}
	fmt.Println("entities Bytes :  " + string(indexByte))

	if marshErr := json.Unmarshal(indexByte, &entities); marshErr != nil {
		fmt.Println("Could not save Shipment to ledger", marshErr)
		return entities, marshErr
	}

	fmt.Println("entities : ======================")
	fmt.Println(entities)
	fmt.Println("entities : ======================")
	fmt.Println("Exiting fetchEntities ")
	return entities, nil

}

func (t *ShipmentPageLoadService) fetchAllEntities(stub shim.ChaincodeStubInterface) (AllEntities, error) {
	fmt.Println("Entering fetchAllEntities ")
	var allEntities AllEntities

	indexByte, err := stub.GetState("ALL_ENTITIES")
	if err != nil {
		fmt.Println("Could not retrive Shipment Index", err)
		return allEntities, err
	}

	if marshErr := json.Unmarshal(indexByte, &allEntities); marshErr != nil {
		fmt.Println("Could not save Shipment to ledger", marshErr)
		return allEntities, marshErr
	}

	fmt.Println("allEntities : ======================")
	fmt.Println(allEntities)
	fmt.Println("allEntities : ======================")
	fmt.Println("Exiting fetchAllEntities ")
	return allEntities, nil

}

func (t *ShipmentPageLoadService) parseShipmentPageLoadRequest(jsondata string) ShipmentPageLoadRequest {
	fmt.Println("Entering parseShipmentPageLoadRequest ")
	res := ShipmentPageLoadRequest{}
	json.Unmarshal([]byte(jsondata), &res)
	fmt.Println(res)
	fmt.Println("Exiting parseShipmentPageLoadRequest ")
	return res
}

/************** ShipmentPageLoad Ends ***************************/
