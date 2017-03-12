/*
Copyright DASE@ECNU. 2016 All Rights Reserved.
*/

package app1

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/hyperledger/fabric/core/crypto"
	pb "github.com/hyperledger/fabric/protos"
	"github.com/op/go-logging"
	"google.golang.org/grpc"

	"github.com/hyperledger/fabric/work/receipt/entity"
)

var (
	// Logging
	appLogger = logging.MustGetLogger("app")

	// NVP related objects
	peerClientConn *grpc.ClientConn
	serverClient   pb.PeerClient

	warehouse     crypto.Client
	warehouse2     crypto.Client
	regulator     crypto.Client
)

func assignRegulator(receipt entity.WarehouseReceipt) (err error) {
	appLogger.Debug("------------- warehouse assigns receipt to regulator...")
    
    warehouseStr := receipt.WAREHOUSE_CODE
    var invoker crypto.Client
    if warehouseStr=="warehouse" {
        invoker = warehouse
    }else{
    	invoker = warehouse2
    }

    var invokerCert crypto.CertificateHandler
    //warehouseCert, err = warehouse.GetTCertificateHandlerNext()
    invokerCert, err = invoker.GetEnrollmentCertificateHandler()
	if err != nil {
		appLogger.Errorf("Failed getting warehouse TCert [%s]", err)
		return
	}

	var regulatorCert crypto.CertificateHandler
	//regulatorCert, err = regulator.GetTCertificateHandlerNext()
	regulatorCert, err = regulator.GetEnrollmentCertificateHandler()
	if err != nil {
		appLogger.Errorf("Failed getting regulator TCert [%s]", err)
		return
	}

	line, _ := ioutil.ReadFile("../config/"+warehouseStr)
    numStr := string(line)
    num, _ := strconv.Atoi(numStr)
    ioutil.WriteFile("../config/"+warehouseStr, []byte(strconv.Itoa(num+1)), 0644)
    receipt.ID = warehouseStr+"-"+numStr
    
    receiptId := receipt.ID
    customerFile := receipt.CUSTOMER_CODE
    line, _ = ioutil.ReadFile("../config/"+customerFile)
    ioutil.WriteFile("../config/"+customerFile, []byte(string(line)+receiptId+","), 0644)
 
    var receiptBytes []byte
    receiptBytes,err = json.Marshal(&receipt)
	if err!= nil{
		appLogger.Errorf("Error marshalling data")
	}
	var receiptJson string = string(receiptBytes)
	appLogger.Debug("receiptJson=", receiptJson)

	resp, err := assignRegulatorInternal(invoker, invokerCert, warehouseStr, receiptId, receiptJson, regulatorCert)
	if err != nil {
		appLogger.Errorf("Failed assigning regulator [%s]", err)
		return
	}
	appLogger.Debugf("Resp [%s]", resp.String())
	appLogger.Debug("------------- Done!")
	return
}

func testReceiptManagementChaincode(receipt entity.WarehouseReceipt) (err error) {
	err = assignRegulator(receipt)

	closeCryptoClient(warehouse)
	closeCryptoClient(warehouse2)
	closeCryptoClient(regulator)

	if err != nil {
		appLogger.Errorf("Failed assigning regulator [%s]", err)
		return
	}
	
	return
}

func Assign(receipt entity.WarehouseReceipt)(msg string){
	// Initialize a non-validating peer whose role is to submit
	// transactions to the fabric network.
	if err := initNVP(); err != nil {
		appLogger.Debugf("Failed initiliazing NVP [%s]", err)
		return "{\"msg\":\""+err.Error()+"\"}"
	}

	// Enable fabric 'confidentiality'
	confidentiality(false)

	if err := testReceiptManagementChaincode(receipt); err != nil {
		appLogger.Debugf("Failed testing receipt management chaincode [%s]", err)
		return "{\"msg\":\""+err.Error()+"\"}"
	}
	return "{\"msg\":\"ok\"}"
}

// func main() {
// 	//  ./app1 warehouse
// 	//  ./app1 warehouse2
// 	if len(os.Args) != 2 {
// 		appLogger.Debugf("Incorrect number of arguments. Expecting 2")
// 		os.Exit(-3)
// 	}
// 	// Initialize a non-validating peer whose role is to submit
// 	// transactions to the fabric network.
// 	if err := initNVP(); err != nil {
// 		appLogger.Debugf("Failed initiliazing NVP [%s]", err)
// 		os.Exit(-1)
// 	}

// 	// Enable fabric 'confidentiality'
// 	confidentiality(true)

// 	if err := testReceiptManagementChaincode(os.Args[1]); err != nil {
// 		appLogger.Debugf("Failed testing receipt management chaincode [%s]", err)
// 		os.Exit(-2)
// 	}
// }

