/*****All structs used in the chaincode*****
Structs Involved
BlockchainResponse
AssetDetails
CartonDetails
PalletDetails
ShipmentWayBill
EWWayBill
EntityWayBillMapping
CreateEntityWayBillMappingRequest
WayBillShipmentMapping

Author: Mohd Arshad
Dated: 30/7/2017
/*****************************************************/
package main

const (
	SHIPMENT         = "SHIPMENT"
	WAYBILL          = "WAYBILL"
	DCSHIPMENT       = "DCSHIPMENT"
	DCWAYBILL        = "DCWAYBILL"
	EWWAYBILL        = "EWWAYBILL"
	RETAILERSHIPMENT = "RETAILERSHIPMENT"
	RETAILERWAYBILL  = "RETAILERWAYBILL"
)

type BlockchainResponse struct {
	Err     string `json:"err"`
	ErrMsg  string `json:"errMsg"`
	Message string `json:"message"`
	TxID    string `json:"txid"`
}

type AssetDetails struct {
	AssetSerialNumber      string `json:"assetSerialNumber"`
	AssetModel             string `json:"assetModel"`
	AssetType              string `json:"assetType"`
	AssetMake              string `json:"assetMake"`
	AssetCOO               string `json:"assetCOO"`
	AssetMaufacture        string `json:"assetMaufacture"`
	AssetStatus            string `json:"assetStatus"`
	CreatedBy              string `json:"createdBy"`
	CreatedDate            string `json:"createdDate"`
	ModifiedBy             string `json:"modifiedBy"`
	ModifiedDate           string `json:"modifiedDate"`
	PalletSerialNumber     string `json:"palletSerialNumber"`
	CartonSerialNumber     string `json:"cartonSerialNumber"`
	MshipmentNumber        string `json:"mshipmentNumber"`
	DcShipmentNumber       string `json:"dcShipmentNumber"`
	MwayBillNumber         string `json:"mwayBillNumber"`
	DcWayBillNumber        string `json:"dcWayBillNumber"`
	EwWayBillNumber        string `json:"ewWayBillNumber"`
	MShipmentDate          string `json:"mShipmentDate"`
	DcShipmentDate         string `json:"dcShipmentDate"`
	MWayBillDate           string `json:"mWayBillDate"`
	DcWayBillDate          string `json:"dcWayBillDate"`
	EwWayBillDate          string `json:"ewWayBillDate"`
	RetailerShipmentNumber string `json:"retailerShipmentNumber"`
	RetailerWaybillNumber  string `json:"retailerWaybillNumber"`
}

type CartonDetails struct {
	CartonSerialNumber     string   `json:"cartonSerialNumber"`
	CartonModel            string   `json:"cartonModel"`
	CartonStatus           string   `json:"cartonStatus"`
	CartonCreatedBy        string   `json:"cartonCreatedBy"`
	CartonCreationDate     string   `json:"cartonCreationDate"`
	PalletSerialNumber     string   `json:"palletSerialNumber"`
	AssetSerialNumber      []string `json:"assetsSerialNumber"`
	MshipmentNumber        string   `json:"mshipmentNumber"`
	DcShipmentNumber       string   `json:"dcShipmentNumber"`
	MwayBillNumber         string   `json:"mwayBillNumber"`
	DcWayBillNumber        string   `json:"dcWayBillNumber"`
	EwWayBillNumber        string   `json:"ewWayBillNumber"`
	Dimensions             string   `json:"dimensions"`
	Weight                 string   `json:"weight"`
	MShipmentDate          string   `json:"mShipmentDate"`
	DcShipmentDate         string   `json:"dcShipmentDate"`
	MWayBillDate           string   `json:"mWayBillDate"`
	DcWayBillDate          string   `json:"dcWayBillDate"`
	EwWayBillDate          string   `json:"ewWayBillDate"`
	RetailerShipmentNumber string   `json:"retailerShipmentNumber"`
	RetailerWaybillNumber  string   `json:"retailerWaybillNumber"`
}

type PalletDetails struct {
	PalletSerialNumber     string   `json:"palletSerialNumber"`
	PalletModel            string   `json:"palletModel"`
	PalletStatus           string   `json:"palletStatus"`
	CartonSerialNumber     []string `json:"cartonSerialNumber"`
	PalletCreationDate     string   `json:"palletCreationDate"`
	AssetSerialNumber      []string `json:"assetsSerialNumber"`
	MshipmentNumber        string   `json:"mshipmentNumber"`
	DcShipmentNumber       string   `json:"dcShipmentNumber"`
	MwayBillNumber         string   `json:"mwayBillNumber"`
	DcWayBillNumber        string   `json:"dcWayBillNumber"`
	EwWayBillNumber        string   `json:"ewWayBillNumber"`
	Dimensions             string   `json:"dimensions"`
	Weight                 string   `json:"weight"`
	MShipmentDate          string   `json:"mShipmentDate"`
	DcShipmentDate         string   `json:"dcShipmentDate"`
	MWayBillDate           string   `json:"mWayBillDate"`
	DcWayBillDate          string   `json:"dcWayBillDate"`
	EwWayBillDate          string   `json:"ewWayBillDate"`
	RetailerShipmentNumber string   `json:"retailerShipmentNumber"`
	RetailerWaybillNumber  string   `json:"retailerWaybillNumber"`
}

type AllShipmentWayBills struct {
	AllWayBillNumber []string `json:"allWayBillNumber"`
}

/*This is common struct across Shipment and Waybill*/
type ShipmentWayBill struct {
	WayBillNumber         string                    `json:"waybillNumber"`
	ShipmentNumber        string                    `json:"shipmentNumber"`
	CountryFrom           string                    `json:"countryFrom"`
	CountryTo             string                    `json:"countryTo"`
	Consigner             string                    `json:"consigner"`
	Consignee             string                    `json:"consignee"`
	Status                string                    `json:"status"`
	Custodian             string                    `json:"custodian"`
	CustodianHistory      []CustodianHistoryDetail  `json:"custodianHistory"`
	PersonConsigningGoods string                    `json:"personConsigningGoods"`
	Comments              string                    `json:"comments"`
	TpComments            string                    `json:"tpComments"`
	VehicleNumber         string                    `json:"vehicleNumber"`
	VehicleType           string                    `json:"vehicleType"`
	PickupDate            string                    `json:"pickupDate"`
	PalletsSerialNumber   []string                  `json:"palletsSerialNumber"`
	CartonsSerialNumber   []string                  `json:"cartonsSerialNumber"`
	AssetsSerialNumber    []string                  `json:"assetsSerialNumber"`
	AddressOfConsigner    string                    `json:"addressOfConsigner"`
	AddressOfConsignee    string                    `json:"addressOfConsignee"`
	ConsignerRegNumber    string                    `json:"consignerRegNumber"`
	Carrier               string                    `json:"carrier"`
	VesselType            string                    `json:"vesselType"`
	VesselNumber          string                    `json:"vesselNumber"`
	ContainerNumber       string                    `json:"containerNumber"`
	ServiceType           string                    `json:"serviceType"`
	ShipmentModel         string                    `json:"shipmentModel"`
	PalletsQuantity       string                    `json:"palletsQuantity"`
	CartonsQuantity       string                    `json:"cartonsQuantity"`
	AssetsQuantity        string                    `json:"assetsQuantity"`
	ShipmentValue         string                    `json:"shipmentValue"`
	EntityName            string                    `json:"entityName"`
	ShipmentCreationDate  string                    `json:"shipmentCreationDate"`
	EWWayBillNumber       string                    `json:"ewWaybillNumber"`
	SupportiveDocuments   []SupportiveDocumentsList `json:"supportiveDocumentsList"`
	ShipmentCreatedBy     string                    `json:"shipmentCreatedBy"`
	ShipmentModifiedDate  string                    `json:"shipmentModifiedDate"`
	ShipmentModifiedBy    string                    `json:"shipmentModifiedBy"`
	WayBillCreationDate   string                    `json:"waybillCreationDate"`
	WayBillCreatedBy      string                    `json:"waybillCreatedBy"`
	WayBillModifiedDate   string                    `json:"waybillModifiedDate"`
	WayBillModifiedBy     string                    `json:"waybillModifiedBy"`
	ShipmentImage         []byte                    `json:"shipmentImage"`
	WaybillImage          []byte                    `json:"waybillImage"`
	DCShipmentImage       []byte                    `json:"dcshipmentImage"`
	DCWaybillImage        []byte                    `json:"dcwaybillImage"`
}
type SupportiveDocumentsList struct {
	DocumentType  string `json:"documentType"`
	DocumentHash  string `json:"documentHash"`
	DocumentTitle string `json:"documentTitle"`
}
type AllEWWayBill struct {
	AllWayBillNumber []string `json:"allWayBillNumber"`
}

type EWWayBill struct {
	EwWayBillNumber       string                    `json:"ewWaybillNumber"`
	WayBillsNumber        []string                  `json:"waybillsNumber"`
	ShipmentsNumber       []string                  `json:"shipmentsNumber"`
	CountryFrom           string                    `json:"countryFrom"`
	CountryTo             string                    `json:"countryTo"`
	Consigner             string                    `json:"consigner"`
	Consignee             string                    `json:"consignee"`
	Custodian             string                    `json:"custodian"`
	Status                string                    `json:"status"`
	CustodianHistory      []CustodianHistoryDetail  `json:"custodianHistory"`
	CustodianTime         string                    `json:"custodianTime"`
	PersonConsigningGoods string                    `json:"personConsigningGoods"`
	Comments              string                    `json:"comments"`
	PalletsSerialNumber   []string                  `json:"palletsSerialNumber"`
	CartonsSerialNumber   []string                  `json:"cartonsSerialNumber"`
	AssetsSerialNumber    []string                  `json:"assetsSerialNumber"`
	AddressOfConsigner    string                    `json:"addressOfConsigner"`
	AddressOfConsignee    string                    `json:"addressOfConsignee"`
	ConsignerRegNumber    string                    `json:"consignerRegNumber"`
	VesselType            string                    `json:"vesselType"`
	VesselNumber          string                    `json:"vesselNumber"`
	ContainerNumber       string                    `json:"containerNumber"`
	ServiceType           string                    `json:"serviceType"`
	SupportiveDocuments   []SupportiveDocumentsList `json:"supportiveDocumentsList"`
	EwWayBillCreationDate string                    `json:"ewWaybillCreationDate"`
	EwWayBillCreatedBy    string                    `json:"ewWaybillCreatedBy"`
	EwWayBillModifiedDate string                    `json:"ewWaybillModifiedDate"`
	EwWayBillModifiedBy   string                    `json:"ewWaybillModifiedBy"`
	EWWayBillImage        []byte                    `json:"ewWaybillImage"`
}

type EntityWayBillMapping struct {
	WayBillsNumber []EntityWayBillMappingDetail
}
type EntityWayBillMappingDetail struct {
	WayBillNumber string
	Country       string
}

type CreateEntityWayBillMappingRequest struct {
	EntityName     string
	WayBillsNumber []EntityWayBillMappingDetail
}
type WayBillShipmentMapping struct {
	DCWayBillsNumber string
	DCShipmentNumber string
}

//storing compliance document mdetadata and hash
type ComplianceDocument struct {
	Compliance_id      string `json:"compliance_id"`
	Manufacturer       string `json:"manufacturer"`
	Regulator          string `json:"regulator"`
	DocumentTitle      string `json:"documentTitle"`
	Document_mime_type string `json:"document_mime_type"`
	DocumentHash       string `json:"documentHash"`
	DocumentType       string `json:"documentType"`
	CreatedBy          string `json:"createdBy"`
	CreatedDateStr     string `json:"createdDateStr"`
}

//mapping for entity and corresponding document
type EntityComplianceDocMapping struct {
	ComplianceIds []string `json:"complianceIds"`
}

//collection of all the compliance document ids
type ComplianceIds struct {
	ComplianceIds []string `json:"complianceIds"`
}

//list of compliance document
type ComplianceDocumentList struct {
	ComplianceDocumentList []ComplianceDocument `json:"complianceList"`
}
type CustodianHistoryDetail struct {
	CustodianName string `json:"CustodianName"`
	Comments      string `json:"comments"`
	CustodianTime string `json:"custodianTime"`
}
type CustodianHistory struct {
	CustodianHistoryList []CustodianHistoryDetail `json:"custodianHistory"`
}

//storing transaction details
type TransactionDetails struct {
	FromUserId      string   `json:"fromUserId"`
	BlockCount      string   `json:"blockCount"`
	CurrentBlock    string   `json:"currentBlock"`
	ToUserId        []string `json:"toUserId"`
	TransactionId   string   `json:"transactionId"`
	Status          string   `json:"status"`
	TransactionTime string   `json:"transactionTime"`
}

//list of transaction document
type TransactionDetailsList struct {
	TransactionDetailsArr []TransactionDetails `json:"transactionDetailsArr"`
}
