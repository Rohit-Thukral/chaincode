package main

/*****Chaincode to perform Compliance tasks*****
Methods Involved
uploadComplianceDocument
getAllComplianceDocument


Author: Santosh Kumar
Dated: 30/7/2017
/*****************************************************/
import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//////////////////////////@@@@@@@@@@@@@@@@@  santosh compliance document   @@@@@@@@@@@@@@@\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\

//method for storing complaince document metadata and hash
func uploadComplianceDocument(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	resp := BlockchainResponse{}
	fmt.Println("uploading compliance document", args[0])
	compDoc, _ := parseComplianceDocument(args[0])
	fmt.Println("uploading compliance document", compDoc)
	complianceId := compDoc.Compliance_id
	saveErr := saveComplianceDocument(stub, complianceId, compDoc)
	if saveErr != nil {
		resp.Err = "000"
		resp.ErrMsg = complianceId
		resp.Message = "Document Not saved"
		respString, _ := json.Marshal(resp)
		return []byte(respString), saveErr
	}
	entityCompMapRequest := EntityComplianceDocMapping{}
	entityCompMap, err := fetchEntityComplianceDocumentMapping(stub, compDoc.Manufacturer)
	fmt.Println("entitycompliancemapping previous entry..", entityCompMap)
	if err != nil {
		entityCompMapRequest.ComplianceIds = append(entityCompMapRequest.ComplianceIds, complianceId)
		fmt.Println("entity comp mapping entity..", entityCompMapRequest)
		saveEntityComplianceDocumentMapping(stub, entityCompMapRequest, compDoc.Manufacturer)
	} else {
		entityCompMapRequest.ComplianceIds = append(entityCompMap.ComplianceIds, complianceId)
		fmt.Println("Updated entity compliance document mapping", entityCompMapRequest)
		saveEntityComplianceDocumentMapping(stub, entityCompMapRequest, compDoc.Manufacturer)
	}
	complianceidsRequest := ComplianceIds{}
	complianceids, err := fetchComplianceDocumentIds(stub, "CompDocIDs")
	fmt.Println("compliance ids.....", complianceids)
	if err != nil {
		complianceidsRequest.ComplianceIds = append(complianceidsRequest.ComplianceIds, complianceId)
		saveComplianceDocumentIds(stub, complianceidsRequest)
	} else {
		complianceidsRequest.ComplianceIds = append(complianceids.ComplianceIds, complianceId)
		fmt.Println("Updated entity compliance document mapping", entityCompMapRequest)
		saveComplianceDocumentIds(stub, complianceidsRequest)
	}
	resp.Err = "000"
	resp.ErrMsg = "Data Saved"
	resp.Message = "Successfully uploaded compliance document to ledger"
	respString, _ := json.Marshal(resp)

	fmt.Println("Successfully uploaded compliance document to ledger")
	return []byte(respString), nil
}

//save entity compliance document mapping in blockchain
func saveEntityComplianceDocumentMapping(stub shim.ChaincodeStubInterface, entityCompMapRequest EntityComplianceDocMapping, entityname string) ([]byte, error) {
	dataToStore, _ := json.Marshal(entityCompMapRequest)
	fmt.Println("entity compliance document mappping data..", dataToStore)
	entitykey := entityname + "ComDoc"
	fmt.Println("key for entity compliance document mapping", entitykey)
	err := stub.PutState(entitykey, []byte(dataToStore))
	if err != nil {
		fmt.Println("Could not save Entity compliance Mapping to ledger", err)
		return nil, err
	}

	resp := BlockchainResponse{}
	resp.Err = "000"
	resp.Message = entityname

	respString, _ := json.Marshal(resp)
	fmt.Println("Successfully saved Entity WayBill Mapping")
	return []byte(respString), nil

}

//save compliance document ids in blockchain
func saveComplianceDocumentIds(stub shim.ChaincodeStubInterface, comids ComplianceIds) ([]byte, error) {
	dataToStore, _ := json.Marshal(comids)
	fmt.Println("save compliance document ids..", dataToStore)
	err := stub.PutState("CompDocIDs", []byte(dataToStore))
	if err != nil {
		fmt.Println("Could not save complianceIds to ledger", err)
		return nil, err
	}

	resp := BlockchainResponse{}
	resp.Err = "000"
	resp.Message = "CompDocIDs"
	respString, _ := json.Marshal(resp)

	fmt.Println("Successfully saved compliance IDs")
	return []byte(respString), nil

}

//get entity name from compliance Document json
func parseComplianceDocument(jsonComDoc string) (ComplianceDocument, error) {
	complianceDoc := ComplianceDocument{}
	fmt.Println("compliance document unmarshal", jsonComDoc)
	json.Unmarshal([]byte(jsonComDoc), &complianceDoc)
	fmt.Println("Unmarshal compliance Document", complianceDoc)
	return complianceDoc, nil
}

//save compliance document to blockchain
func saveComplianceDocument(stub shim.ChaincodeStubInterface, complianceId string, compDoc ComplianceDocument) error {
	dataToStore, _ := json.Marshal(compDoc)
	fmt.Println("compliance id ....", complianceId)
	fmt.Println("compliance documnt to store.....", dataToStore)
	err := stub.PutState(complianceId, []byte(dataToStore))
	if err != nil {
		fmt.Println("compliance document not uploaded to ledger", err)
		return err
	}
	return err
}

//fetch entity compliance document mapping
func fetchEntityComplianceDocumentMapping(stub shim.ChaincodeStubInterface, entityname string) (EntityComplianceDocMapping, error) {
	entityComplianceDocMapping := EntityComplianceDocMapping{}
	entitykey := entityname + "ComDoc"
	indexByte, err := stub.GetState(entitykey)
	if err != nil {
		fmt.Println("Could not retrive entity compliance mapping ", err)
		return entityComplianceDocMapping, err
	}

	if marshErr := json.Unmarshal(indexByte, &entityComplianceDocMapping); marshErr != nil {
		fmt.Println("Could not retrive entity compliance mapping from ledger", marshErr)
		return entityComplianceDocMapping, marshErr
	}

	return entityComplianceDocMapping, nil

}

//fetch compliance ids collection
func fetchComplianceDocumentIds(stub shim.ChaincodeStubInterface, compkey string) (ComplianceIds, error) {
	complianceids := ComplianceIds{}
	indexByte, err := stub.GetState(compkey)
	fmt.Println("already stored compliance document keys", indexByte)
	if err != nil {
		fmt.Println("Could not retrive complianceids", err)
		return complianceids, err
	}

	json.Unmarshal(indexByte, &complianceids)
	fmt.Println("retrieved compliance keys after unmarshal..", complianceids)
	return complianceids, nil

}

func getComplianceDocumentByEntityName(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	complianceDocumentList := ComplianceDocumentList{}
	entityComplianceMapping, err := fetchEntityComplianceDocumentMapping(stub, args[0])
	if err != nil {
		return nil, nil
	} else {
		iterator := len(entityComplianceMapping.ComplianceIds)
		for i := 0; i < iterator; i++ {
			complianceDocuments, _ := fetchComplianceDocumentByComplianceId(stub, entityComplianceMapping.ComplianceIds[i])
			complianceDocumentList.ComplianceDocumentList = append(complianceDocumentList.ComplianceDocumentList, complianceDocuments)
		}
		dataToReturn, _ := json.Marshal(complianceDocumentList)
		return []byte(dataToReturn), nil
	}
	return nil, nil
}
func fetchComplianceDocumentByComplianceId(stub shim.ChaincodeStubInterface, complianceid string) (ComplianceDocument, error) {
	complianceDocument := ComplianceDocument{}
	indexByte, err := stub.GetState(complianceid)
	if err != nil {
		fmt.Println("Could not retrive compliance document", err)
		return complianceDocument, err
	}

	if marshErr := json.Unmarshal(indexByte, &complianceDocument); marshErr != nil {
		fmt.Println("Could not retrive complianceids from ledger", marshErr)
		return complianceDocument, marshErr
	}

	return complianceDocument, nil

}

func getAllComplianceDocument(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	complianceDocumentList := ComplianceDocumentList{}
	complianceIds, err := fetchComplianceDocumentIds(stub, args[0])
	if err != nil {
		return nil, err
	} else {
		iterator := len(complianceIds.ComplianceIds)
		for i := 0; i < iterator; i++ {
			complianceDocuments, _ := fetchComplianceDocumentByComplianceId(stub, complianceIds.ComplianceIds[i])
			complianceDocumentList.ComplianceDocumentList = append(complianceDocumentList.ComplianceDocumentList, complianceDocuments)
		}
		dataToReturn, _ := json.Marshal(complianceDocumentList)
		return []byte(dataToReturn), nil
	}

}

///////////////////////////////////////////////////////end compliance docuent \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
