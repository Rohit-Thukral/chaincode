package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type TransactionService struct {
}

type TransactionRequest struct {
	CallingEntityName string `json:"callingEntityName"`
}

type TransactionResponse struct {
	TransactionDetailsArr []TransactionDetails `json:"transactionDetailsArr"`
}

func (t *TransactionService) GetTransactionRecords(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Transaction Details " + args[0])
	var resp TransactionResponse
	var thisClass TransactionService
	request := thisClass.parseTransactionRequest(args[0])

	tmpEntity, err := thisClass.fetchTransactionEntity(stub, "TransactionDetailsKey")
	if err != nil {
		fmt.Println("Error while retrieveing the Transaction Details", err)
		return nil, err
	}
	resp.TransactionDetailsArr = thisClass.checkTransactionCondition(request.CallingEntityName, tmpEntity.TransactionDetailsArr)
	datatoreturn, _ := json.Marshal(resp)
	return []byte(datatoreturn), nil
}

func (t *TransactionService) fetchTransactionEntity(stub shim.ChaincodeStubInterface, transactionKey string) (TransactionDetailsList, error) {
	fmt.Println("Entering fetchTransactionEntity" + transactionKey)
	var transactionDetailsData TransactionDetailsList

	indexByte, err := stub.GetState(transactionKey)
	if err != nil {
		fmt.Println("Could not retrive Transaction Details", err)
		return transactionDetailsData, err
	}
	fmt.Println("entities Bytes :  " + string(indexByte))

	marshErr := json.Unmarshal(indexByte, &transactionDetailsData)
	if marshErr != nil {
		fmt.Println("Could not retrieve transaction details", marshErr)
		return transactionDetailsData, marshErr
	}

	fmt.Println("transactionDetailsData : ======================")
	fmt.Println(transactionDetailsData)
	fmt.Println("transactionDetailsData : ======================")
	fmt.Println("Exiting transactionDetailsData ")
	return transactionDetailsData, nil

}

func (t *TransactionService) checkTransactionCondition(entityId string, txArr []TransactionDetails) []TransactionDetails {

	fmt.Println("In Check transaction condition method entityID-->", entityId, "TxArr", txArr)
	var txDetailsArr []TransactionDetails
	lenOfArray := len(txArr)
	fmt.Println("===lenOfArray all Transaction details===", lenOfArray)

	for i := 0; i < lenOfArray; i++ {
		txDetails := txArr[i]
		toUserArr := txArr[i].ToUserId
		if txDetails.FromUserId == entityId {
			txDetailsArr = append(txDetailsArr, txDetails)
		} else {
			var toUserlen = len(toUserArr)
			for j := 0; j < toUserlen; j++ {
				if toUserArr[j] == entityId {
					txDetailsArr = append(txDetailsArr, txDetails)
				}
			}
		}

	}

	return txDetailsArr
}

func (t *TransactionService) parseTransactionRequest(jsondata string) TransactionRequest {
	fmt.Println("Entering parseTransactionRequest ", jsondata)
	res := TransactionRequest{}
	json.Unmarshal([]byte(jsondata), &res)

	fmt.Println("======================")
	fmt.Println(res)
	fmt.Println("======================")

	return res
}
