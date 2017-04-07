package main


import (
	"encoding/json"
	"fmt"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type InboxService struct {

}

type InboxRequest struct {
	CallingEntityName string `json:"callingEntityName"`
	Status            string `json:"status"`
}

type InboxResponse struct {
	Data []ShipmentWayBill `json:"data"`
}



func (t *InboxService) Inbox(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Inbox " + args[0])
	var util Utility
	var shipmentArray []ShipmentWayBill
	var resp InboxResponse
	var err error

	request := parseInboxRequest(args[0])
	shipmentArray, err = util.fetchShipmentIndex(stub, request.CallingEntityName, request.Status)

	if(err == nil) {
		resp.Data =  shipmentArray
	}else {
		return nil,err
	}

	fmt.Println("======================")
	fmt.Println(resp)
	fmt.Println("======================")
	fmt.Println("Exiting Inbox")
	return json.Marshal(resp)

}

func parseInboxRequest(jsondata string) InboxRequest {
	fmt.Println("Entering parseInboxRequest " + jsondata)
	res := InboxRequest{}
	json.Unmarshal([]byte(jsondata), &res)
	
	fmt.Println("======================")
	fmt.Println(res)
	fmt.Println("======================")

	return res
}

