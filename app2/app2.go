/*
Copyright DASE@ECNU. 2016 All Rights Reserved.
*/

package app2

import (
	"encoding/json"
	"errors"

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

	regulator     crypto.Client
	ouyeel     crypto.Client
	customer     crypto.Client
	customer2     crypto.Client
)

func transferOwnership(owner crypto.Client, newOwner crypto.Client, receipt entity.WarehouseReceipt) (err error) {
	appLogger.Debug("------------- owner transfer the warehouse receipt to newOwner...")

	//ownerCert, err = owner.GetTCertificateHandlerNext()
	ownerCert, err := owner.GetEnrollmentCertificateHandler()
	if err != nil {
		appLogger.Errorf("Failed getting owner TCert [%s]", err)
		return
	}

	newOwnerCert, err := newOwner.GetEnrollmentCertificateHandler()
	if err != nil {
		appLogger.Errorf("Failed getting newOwner TCert [%s]", err)
		return
	}

	// receiptJson, err := queryReceipt(owner, receiptId)
	// if err != nil {
	// 	appLogger.Errorf("Failed query receiptId=[%s], err:[%s]", receiptId, err)
	// 	return
	// }

	receiptBytes,err := json.Marshal(&receipt)
	if err!= nil{
		appLogger.Errorf("Error marshalling data")
	}
	var receiptJson = string(receiptBytes)
	appLogger.Debug("receiptJson=", receiptJson)

	resp, err := transferOwnershipInternal(owner, ownerCert, receipt.ID, receiptJson, newOwnerCert)
	if err != nil {
		appLogger.Errorf("Failed transfering ownership [%s]", err)
		return
	}
	appLogger.Debugf("Resp [%s]", resp.String())

	appLogger.Debug("------------- Done!")
	return
}

func queryReceipt(owner crypto.Client, receiptId string) (receiptJson string, err error) {
	appLogger.Debug("------------- Query by receiptId=", receiptId)
    
	queryTx, theOwnerIs, err := whoIsTheOwner(owner, receiptId, "getReceipt")
	if err != nil {
		return "", err
	}

	appLogger.Debugf("Resp=[%s]", theOwnerIs)
	appLogger.Debug("Query....done")

	var res []byte
	if confidentialityOn {
		// Decrypt result
		res, err = owner.DecryptQueryResult(queryTx, theOwnerIs.Msg)
		if err != nil {
			appLogger.Errorf("Failed decrypting result [%s]", err)
			return "", err
		}
	} else {
		res = theOwnerIs.Msg
	}
    appLogger.Debugf("res＝%s", res)
	appLogger.Debug("------------- Query Done!")
	return string(res), nil 
}

func testReceiptManagementChaincode(receipt entity.WarehouseReceipt) (err error) {
    if receipt.CREGULEOPTION=="监管公司未核" {
    	receipt.CREGULEOPTION="监管公司已核"
        err = transferOwnership(regulator, ouyeel, receipt) 
        appLogger.Debug("------------- 监管公司已核!")	
    } else if receipt.CPLATFORMOPTION=="平台未核" {
    	receipt.CPLATFORMOPTION="平台已核"
    	if receipt.CUSTOMER_CODE=="customer"{
    	    err = transferOwnership(ouyeel, customer, receipt)
        } else {
        	err = transferOwnership(ouyeel, customer2, receipt)
        }
        appLogger.Debug("------------- 平台已核!")
    } else if receipt.COWNEROPTION=="存货人未验收"{
    	receipt.COWNEROPTION="存货人已验收"
    	if receipt.CUSTOMER_CODE=="customer"{
            err = transferOwnership(customer, customer, receipt)
         } else {
        	err = transferOwnership(customer2, customer2, receipt)
        }
        appLogger.Debug("------------- 存货人已验收!")  
    } else {
    	return errors.New("非法状态的审核数据！")
    }
	
	closeCryptoClient(regulator)
	closeCryptoClient(ouyeel)
	closeCryptoClient(customer)
	closeCryptoClient(customer2)

	if err != nil {
		appLogger.Errorf("Failed tarnsfering ownership [%s]", err)
		return
	}

	return
}

func Transfer(receipt entity.WarehouseReceipt)(msg string){
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
// 	if len(os.Args) != 3 {
// 		appLogger.Debugf("Incorrect number of arguments. Expecting 3")
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

// 	if err := testReceiptManagementChaincode(os.Args[1], os.Args[2]); err != nil {
// 		appLogger.Debugf("Failed testing receipt management chaincode [%s]", err)
// 		os.Exit(-2)
// 	}
// }
