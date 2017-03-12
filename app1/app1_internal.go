/*
Copyright DASE@ECNU. 2016 All Rights Reserved.
*/

package app1

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/config"
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
	if err = initCryptoClients(); err != nil {
		appLogger.Debugf("Failed deploying [%s]", err)
		return
	}

	return
}

func initPeerClient() (err error) {
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

	// warehouse
	if err := crypto.RegisterClient("lukas", nil, "lukas", "NPKYL39uKbkj"); err != nil {
		return err
	}
	var err error
	warehouse, err = crypto.InitClient("lukas", nil)
	if err != nil {
		return err
	}

	// regulator
	if err := crypto.RegisterClient("diego", nil, "diego", "DRJ23pEQl16a"); err != nil {
		return err
	}
	regulator, err = crypto.InitClient("diego", nil)
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

func assignRegulatorInternal(invoker crypto.Client, invokerCert crypto.CertificateHandler, invokerStr string, receiptId string, receiptJson string, regulatorCert crypto.CertificateHandler) (resp *pb.Response, err error) {
	// Get a transaction handler to be used to submit the execute transaction
	// and bind the chaincode access control logic using the binding
	submittingCertHandler, err := invoker.GetTCertificateHandlerNext()
	if err != nil {
		return nil, err
	}
	txHandler, err := submittingCertHandler.GetTransactionHandler()
	if err != nil {
		return nil, err
	}
	binding, err := txHandler.GetBinding()
	if err != nil {
		return nil, err
	}
    
    input, _:= ioutil.ReadFile("../config/chaincodeName")
    chaincodeName = string(input)
	
	chaincodeInput := &pb.ChaincodeInput{
		Args: util.ToChaincodeArgs("assign", receiptId, receiptJson, invokerStr, base64.StdEncoding.EncodeToString(regulatorCert.GetCertificate())),
	}
	chaincodeInputRaw, err := proto.Marshal(chaincodeInput)
	if err != nil {
		return nil, err
	}

	// Access control. warehouse signs chaincodeInputRaw || binding to confirm his identity
	sigma, err := invokerCert.Sign(append(chaincodeInputRaw, binding...))
	if err != nil {
		return nil, err
	}

	// Prepare spec and submit
	spec := &pb.ChaincodeSpec{
		Type:                 1,
		//ChaincodeID: &pb.ChaincodeID{Path: "github.com/hyperledger/fabric/work/receipt/chaincode"},
		ChaincodeID:          &pb.ChaincodeID{Name: chaincodeName},
		CtorMsg:              chaincodeInput,
		Metadata:             sigma, // Proof of identity
		ConfidentialityLevel: confidentialityLevel,
	}

	chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

	// Now create the Transactions message and send to Peer.
	var uuid = util.GenerateUUID()
	transaction, err := txHandler.NewChaincodeExecute(chaincodeInvocationSpec, uuid)
	if err != nil {
		return nil, fmt.Errorf("Error deploying chaincode: %s ", err)
	}

	line, _ := ioutil.ReadFile("../config/blockchain")
    ioutil.WriteFile("../config/blockchain", []byte(string(line)+receiptId+","+invokerStr+","+uuid+"\n"), 0644)

	return processTransaction(transaction)
}