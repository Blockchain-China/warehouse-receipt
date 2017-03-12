/*
Copyright DASE@ECNU. 2016 All Rights Reserved.
*/

package query

import (
	"fmt"
	"reflect"
	"encoding/json"
	"strings"
	"strconv"
	"io/ioutil"

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
	regulator     crypto.Client
	ouyeel     crypto.Client
	customer     crypto.Client
)

func queryOwner(ownerCert crypto.CertificateHandler, receiptId string) (err error) {
	appLogger.Debug("------------- Query owner...")
    
    //getOwner
	queryTx, theOwnerIs, err := whoIsTheOwner(warehouse, receiptId, "getOwner")
	if err != nil {
		return
	}
	appLogger.Debugf("Resp [%s]", theOwnerIs.String())
	appLogger.Debug("Query....done")

	var res []byte
	if confidentialityOn {
		// Decrypt result
		res, err = warehouse.DecryptQueryResult(queryTx, theOwnerIs.Msg)
		if err != nil {
			appLogger.Errorf("Failed decrypting result [%s]", err)
			return
		}
	} else {
		res = theOwnerIs.Msg
	}

	if !reflect.DeepEqual(res, ownerCert.GetCertificate()) {
		appLogger.Error("the parameter is not the owner.")

		appLogger.Debugf("Query result  : [% x]", res)
		appLogger.Debugf("parameter's cert: [% x]", ownerCert.GetCertificate())

		return fmt.Errorf("the parameter is not the owner.")
	}
	appLogger.Debug("the parameter is the owner!")
	appLogger.Debug("------------- Done!")
	return
}

func queryReceipt(receiptId string) (err error) {
	appLogger.Debug("------------- Query receipt...")
    
	queryTx, theOwnerIs, err := whoIsTheOwner(warehouse, receiptId, "getReceipt")
	if err != nil {
		return
	}

	appLogger.Debugf("Resp=[%s]", theOwnerIs)
	// if theOwnerIs.Status==200 {
	//     appLogger.Debugf("Resp.Msg=[%x]", theOwnerIs.Msg)
	// }else{
	//     appLogger.Errorf("Resp.Status=[%x]", theOwnerIs.Status)
	// 	   return
	// }
	appLogger.Debug("Query....done")

	var res []byte
	if confidentialityOn {
		// Decrypt result
		res, err = warehouse.DecryptQueryResult(queryTx, theOwnerIs.Msg)
		if err != nil {
			appLogger.Errorf("Failed decrypting result [%s]", err)
			return
		}
	} else {
		res = theOwnerIs.Msg
	}
    appLogger.Debugf("res＝%s", res)

	appLogger.Debug("------------- Done!")
	return
}

func queryReceipts(receiptIds []string) (receipts []entity.WarehouseReceipt, err error) {
	appLogger.Debug("------------- Query receipts...")
    
	// receipts = make([]entity.WarehouseReceipt, 0, len(receiptIds))
	for _, receiptId := range receiptIds {
		
		queryTx, theOwnerIs, err := whoIsTheOwner(warehouse, receiptId, "getReceipt")
		if err != nil {
		    appLogger.Errorf("Failed query [%s]", err)
			return nil, err
		}

		appLogger.Debug("Query receiptId=["+receiptId+"]...")

		var res []byte
		if confidentialityOn {
			// Decrypt result
			res, err = warehouse.DecryptQueryResult(queryTx, theOwnerIs.Msg)
			if err != nil {
				appLogger.Errorf("Failed decrypting result [%s]", err)
				return nil, err
			}
		} else {
			res = theOwnerIs.Msg
		}
		appLogger.Debugf("res＝%s", res)
		
		var receipt_ entity.WarehouseReceipt
        err = json.Unmarshal(res, &receipt_)
        if err != nil {
		    appLogger.Errorf("Failed Unmarshal [%s]", err)
			return nil, err
		}
        receipts = append(receipts, receipt_)
    }

	appLogger.Debug("------------- Done!")
	return receipts, nil
}

func testReceiptManagementChaincode(username string) (msg string, err error) {
	// var ownerCert crypto.CertificateHandler
	// ownerCert, err = regulator.GetEnrollmentCertificateHandler()
	// ownerCert, err = ouyeel.GetEnrollmentCertificateHandler()
	// ownerCert, err = customer.GetEnrollmentCertificateHandler()
	// if err != nil {
	// 	appLogger.Errorf("Failed getting TCert [%s]", err)
	// 	return
	// }
	// err = queryOwner(ownerCert, receiptId)

	// err = queryReceipt(receiptId)
	// if err != nil {
	// 	appLogger.Errorf("Failed query [%s]", err)
	// 	return
	// }

	var receiptIds []string

	if username == "warehouse" {
		input, err := ioutil.ReadFile("../config/"+username)
		line := string(input)
		if err!= nil{
        	appLogger.Errorf("Error readFile [%s]", err)
        	return "", err
    	}

		n,_:=strconv.Atoi(line)
		for i:=1;i<n;i++ {
            receiptIds = append(receiptIds, "warehouse-"+strconv.Itoa(i))    
        }
	} else if username == "warehouse2" {
		input, err := ioutil.ReadFile("../config/"+username)
		line := string(input)
		if err!= nil{
        	appLogger.Errorf("Error readFile [%s]", err)
        	return "", err
    	}

		n,_:=strconv.Atoi(line)
		for i:=1;i<n;i++ {
            receiptIds = append(receiptIds, "warehouse2-"+strconv.Itoa(i))     
        }
	} else if username == "customer" {
		input, err := ioutil.ReadFile("../config/"+username)
		line := string(input)
		if err!= nil{
        	appLogger.Errorf("Error readFile [%s]", err)
        	return "", err
    	}

	    arr := strings.Split(string(line), ",")
	    receiptIds = arr[:len(arr)-1]
	} else if username == "customer2" {
		input, err := ioutil.ReadFile("../config/"+username)
		line := string(input)
		if err!= nil{
        	appLogger.Errorf("Error readFile [%s]", err)
        	return "", err
    	}
		
		arr := strings.Split(string(line), ",")
	    receiptIds = arr[:len(arr)-1]
	} else {
		input, err := ioutil.ReadFile("../config/"+"warehouse")
		line := string(input)
	    if err!= nil{
            appLogger.Errorf("Error readFile [%s]", err)
            return "", err
        }
        n,err := strconv.Atoi(line)
		for i:=1;i<n;i++ {
            receiptIds = append(receiptIds, "warehouse-"+strconv.Itoa(i))    
        }

        input, err = ioutil.ReadFile("../config/"+"warehouse2")
        line = string(input)
	    if err!= nil{
            appLogger.Errorf("Error readFile [%s]", err)
            return "", err
        }
        n,err = strconv.Atoi(line)
		for i:=1;i<n;i++ {
            receiptIds = append(receiptIds, "warehouse2-"+strconv.Itoa(i))     
        }
	}

	receipts, err:= queryReceipts(receiptIds)
	if err != nil {
		appLogger.Errorf("Failed query [%s]", err)
		return "", err
	}

	closeCryptoClient(warehouse)
	closeCryptoClient(regulator)
	closeCryptoClient(ouyeel)
	closeCryptoClient(customer)

	receiptsBytes, err := json.Marshal(&receipts)
    if err != nil {
		appLogger.Errorf("Failed query [%s]", err)
		return "", err
	}

	return string(receiptsBytes), nil
}

func Query(username string)(msg string) {
	// Initialize a non-validating peer whose role is to submit
	// transactions to the fabric network.
	if err := initNVP(); err != nil {
		appLogger.Debugf("Failed initiliazing NVP [%s]", err)
		return "{\"msg\":\""+err.Error()+"\"}"
	}

	// Enable fabric 'confidentiality'
	confidentiality(false)

	msg, err := testReceiptManagementChaincode(username)
	if err != nil {
		appLogger.Debugf("Failed testing receipt management chaincode [%s]", err)
		return "{\"msg\":\""+err.Error()+"\"}"
	}
	return msg
}

// func main() {
// 	//  ./query warehouse-1
// 	//  ./query warehouse2-1
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
