/*
Copyright DASE@ECNU. 2016 All Rights Reserved.
*/

package main

import (
	"os"

	"github.com/hyperledger/fabric/core/crypto"
	pb "github.com/hyperledger/fabric/protos"
	"github.com/op/go-logging"
	"google.golang.org/grpc"
)

var (
	// Logging
	appLogger = logging.MustGetLogger("app")

	// NVP related objects
	peerClientConn *grpc.ClientConn
	serverClient   pb.PeerClient

	deployer crypto.Client
	warehouse crypto.Client
    warehouse2 crypto.Client

    warehouseCert crypto.CertificateHandler
	warehouseCert2 crypto.CertificateHandler
)

func deploy() (err error) {
	appLogger.Debug("------------- deploying chaincode...")

	//warehouseCert, err = warehouse.GetTCertificateHandlerNext()
	warehouseCert, err = warehouse.GetEnrollmentCertificateHandler()
	if err != nil {
		appLogger.Errorf("Failed getting warehouse TCert [%s]", err)
		return
	}

	warehouseCert2, err = warehouse2.GetEnrollmentCertificateHandler()
	if err != nil {
		appLogger.Errorf("Failed getting warehouse2 TCert [%s]", err)
		return
	}

	resp, err := deployInternal(deployer, warehouseCert, warehouseCert2)
	if err != nil {
		appLogger.Errorf("Failed deploying [%s]", err)
		return
	}
	appLogger.Debugf("Resp [%s]", resp.String())
	appLogger.Debugf("Chaincode NAME: [%s]-[%s]", chaincodeName, string(resp.Msg))

	appLogger.Debug("------------- Done!")
	return
}

func testReceiptManagementChaincode() (err error) {
	// Deploy
	err = deploy()
	if err != nil {
		appLogger.Errorf("Failed deploying [%s]", err)
		return
	}

	closeCryptoClient(deployer)
	closeCryptoClient(warehouse)
	closeCryptoClient(warehouse2)

	return
}

func main() {
	// Initialize a non-validating peer
	if err := initNVP(); err != nil {
		appLogger.Debugf("Failed initiliazing NVP [%s]", err)
		os.Exit(-1)
	}

	// Enable fabric 'confidentiality'
	confidentiality(false)

	// Exercise the 'receipt_management' chaincode
	if err := testReceiptManagementChaincode(); err != nil {
		appLogger.Debugf("Failed testing receipt management chaincode [%s]", err)
		os.Exit(-2)
	}
}