/*
Copyright DASE@ECNU. 2016 All Rights Reserved.
*/

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"encoding/base64"

	"github.com/hyperledger/fabric/core/chaincode"
	"github.com/hyperledger/fabric/core/chaincode/platforms"
	"github.com/hyperledger/fabric/core/config"
	"github.com/hyperledger/fabric/core/container"
	"github.com/hyperledger/fabric/core/crypto"
	"github.com/hyperledger/fabric/core/peer"
	"github.com/hyperledger/fabric/core/util"
	pb "github.com/hyperledger/fabric/protos"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

var (
	confidentialityOn bool

	confidentialityLevel pb.ConfidentialityLevel
	chaincodeName string
)

func initNVP() (err error) {
	if err = initPeerClient(); err != nil {
		appLogger.Debugf("Failed deploying [%s]", err)
		return
	}

    if err = os.RemoveAll(viper.GetString("peer.fileSystemPath")+"/crypto"); err != nil {
        appLogger.Debugf("Failed removing [/crypto] [%s]\n", err)
        return
    }

	if err = initCryptoClients(); err != nil {
		appLogger.Debugf("Failed deploying [%s]", err)
		return
	}

	return
}

func initPeerClient() (err error) {
	//config.SetupTestConfig(".")
	config.SetupTestConfig("../")
	viper.Set("ledger.blockchain.deploy-system-chaincode", "false")
	viper.Set("peer.validator.validity-period.verification", "false")

	peerClientConn, err = peer.NewPeerClientConnection()
	if err != nil {
		fmt.Printf("error connection to server at host:port = %s\n", viper.GetString("peer.address"))
		return
	}
	serverClient = pb.NewPeerClient(peerClientConn)

	// Logging
	var formatter = logging.MustStringFormatter(
		`%{color}[%{module}] %{shortfunc} [%{shortfile}] -> %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	logging.SetFormatter(formatter)

	return
}

func initCryptoClients() error {
	crypto.Init()

	// deployer as jim
	if err := crypto.RegisterClient("jim", nil, "jim", "6avZQLwcUe9b"); err != nil {
		return err
	}
	var err error
	deployer, err = crypto.InitClient("jim", nil)
	if err != nil {
		return err
	}

	// warehouse as lukas
	if err := crypto.RegisterClient("lukas", nil, "lukas", "NPKYL39uKbkj"); err != nil {
		return err
	}
	warehouse, err = crypto.InitClient("lukas", nil)
	if err != nil {
		return err
	}

	// warehouse2 as assigner
	if err := crypto.RegisterClient("assigner", nil, "assigner", "Tc43PeqBl11"); err != nil {
		return err
	}
	warehouse2, err = crypto.InitClient("assigner", nil)
	if err != nil {
		return err
	}

	ioutil.WriteFile("../config/warehouse", []byte("1"), 0666)
	ioutil.WriteFile("../config/warehouse2", []byte("1"), 0666)
	ioutil.WriteFile("../config/customer", []byte(""), 0666)
	ioutil.WriteFile("../config/customer2", []byte(""), 0666)
	ioutil.WriteFile("../config/blockchain", []byte(""), 0666)
	ioutil.WriteFile("../config/member", []byte(""), 0666)

	return nil
}

func closeCryptoClient(client crypto.Client) {
	crypto.CloseClient(client)
}

func processTransaction(tx *pb.Transaction) (*pb.Response, error) {
	return serverClient.ProcessTransaction(context.Background(), tx)
}

func confidentiality(enabled bool) {
	confidentialityOn = enabled

	if confidentialityOn {
		confidentialityLevel = pb.ConfidentialityLevel_CONFIDENTIAL
	} else {
		confidentialityLevel = pb.ConfidentialityLevel_PUBLIC
	}
}

func deployInternal(deployer crypto.Client, warehouseCert crypto.CertificateHandler, warehouseCert2 crypto.CertificateHandler) (resp *pb.Response, err error) {
	// Prepare the spec. The metadata includes the identity of the warehouse
	spec := &pb.ChaincodeSpec{
		Type:        1,
		ChaincodeID: &pb.ChaincodeID{Path: "github.com/hyperledger/fabric/work/receipt/chaincode"},
		//ChaincodeID:          &pb.ChaincodeID{Name: chaincodeName},
		CtorMsg:              &pb.ChaincodeInput{Args: util.ToChaincodeArgs("init", base64.StdEncoding.EncodeToString(warehouseCert2.GetCertificate()))},
		Metadata:             warehouseCert.GetCertificate(),
		ConfidentialityLevel: confidentialityLevel,
	}

	// First build the deployment spec
	cds, err := getChaincodeBytes(spec)
	if err != nil {
		return nil, fmt.Errorf("Error getting deployment spec: %s ", err)
	}

	// Now create the Transactions message and send to Peer.
	transaction, err := deployer.NewChaincodeDeployTransaction(cds, cds.ChaincodeSpec.ChaincodeID.Name)
	if err != nil {
		return nil, fmt.Errorf("Error deploying chaincode: %s ", err)
	}

	resp, err = processTransaction(transaction)

	appLogger.Debugf("resp [%s]", resp.String())

	chaincodeName = cds.ChaincodeSpec.ChaincodeID.Name
	ioutil.WriteFile("../config/chaincodeName", []byte(chaincodeName), 0644)
	appLogger.Debugf("ChaincodeName [%s]", chaincodeName)

	return
}

func getChaincodeBytes(spec *pb.ChaincodeSpec) (*pb.ChaincodeDeploymentSpec, error) {
	mode := viper.GetString("chaincode.mode")
	var codePackageBytes []byte
	if mode != chaincode.DevModeUserRunsChaincode {
		appLogger.Debugf("Received build request for chaincode spec: %v", spec)
		var err error
		if err = checkSpec(spec); err != nil {
			return nil, err
		}

		codePackageBytes, err = container.GetChaincodePackageBytes(spec)
		if err != nil {
			err = fmt.Errorf("Error getting chaincode package bytes: %s", err)
			appLogger.Errorf("%s", err)
			return nil, err
		}
	}
	chaincodeDeploymentSpec := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: codePackageBytes}
	return chaincodeDeploymentSpec, nil
}

func checkSpec(spec *pb.ChaincodeSpec) error {
	// Don't allow nil value
	if spec == nil {
		return errors.New("Expected chaincode specification, nil received")
	}

	platform, err := platforms.Find(spec.Type)
	if err != nil {
		return fmt.Errorf("Failed to determine platform type: %s", err)
	}

	return platform.ValidateSpec(spec)
}
