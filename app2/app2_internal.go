/*
Copyright DASE@ECNU. 2016 All Rights Reserved.
*/

package app2

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
		appLogger.Debugf("Failed initPeerClient [%s]", err)
		return

	}
	if err = initCryptoClients(); err != nil {
		appLogger.Debugf("Failed initCryptoClients [%s]", err)
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

	// regulator
	if err := crypto.RegisterClient("diego", nil, "diego", "DRJ23pEQl16a"); err != nil {
		return err
	}
	var err error
	regulator, err = crypto.InitClient("diego", nil)
	if err != nil {
		return err
	}

	// ouyeel
	if err := crypto.RegisterClient("binhn", nil, "binhn", "7avZQLwcUe9q"); err != nil {
		return err
	}
	ouyeel, err = crypto.InitClient("binhn", nil)
	if err != nil {
		return err
	}

	// customer
	if err := crypto.RegisterClient("alice", nil, "alice", "CMS10pEQlB16"); err != nil {
		return err
	}
	customer, err = crypto.InitClient("alice", nil)
	if err != nil {
		return err
	}

	// customer2
	if err := crypto.RegisterClient("bob", nil, "bob", "NOE63pEQbL25"); err != nil {
		return err
	}
	customer2, err = crypto.InitClient("bob", nil)
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

func transferOwnershipInternal(invoker crypto.Client, invokerCert crypto.CertificateHandler, receiptId string, receiptJson string, newOwnerCert crypto.CertificateHandler) (resp *pb.Response, err error) {
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
		Args: util.ToChaincodeArgs("transfer", receiptId, receiptJson, base64.StdEncoding.EncodeToString(newOwnerCert.GetCertificate())),
	}
	chaincodeInputRaw, err := proto.Marshal(chaincodeInput)
	if err != nil {
		return nil, err
	}

	// Access control. signs chaincodeInputRaw || binding to confirm his identity
	sigma, err := invokerCert.Sign(append(chaincodeInputRaw, binding...))
	if err != nil {
		return nil, err
	}

	// Prepare spec and submit
	spec := &pb.ChaincodeSpec{
		Type:                 1,
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
		return nil, fmt.Errorf("Error invoking chaincode: %s ", err)
	}

    var invokerStr string
    if invoker==regulator {
        invokerStr="regulator" 
    } else if invoker==ouyeel {
        invokerStr="ouyeel"
    } else if invoker==customer{
    	invokerStr="customer"
    } else {
    	invokerStr="customer2"
    }
	line, _ := ioutil.ReadFile("../config/blockchain")
    ioutil.WriteFile("../config/blockchain", []byte(string(line)+receiptId+","+invokerStr+","+uuid+"\n"), 0644)

	return processTransaction(transaction)
}

func whoIsTheOwner(invoker crypto.Client, asset string, getCol string) (transaction *pb.Transaction, resp *pb.Response, err error) {
	input, _:= ioutil.ReadFile("../config/chaincodeName")
    chaincodeName = string(input)
	
	chaincodeInput := &pb.ChaincodeInput{Args: util.ToChaincodeArgs("query", getCol, asset)}

	// Prepare spec and submit
	spec := &pb.ChaincodeSpec{
		Type:                 1,
		//ChaincodeID: &pb.ChaincodeID{Path: "github.com/hyperledger/fabric/work/receipt/chaincode"},
		ChaincodeID:          &pb.ChaincodeID{Name: chaincodeName},
		CtorMsg:              chaincodeInput,
		ConfidentialityLevel: confidentialityLevel,
	}

	chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

	// Now create the Transactions message and send to Peer.
	transaction, err = invoker.NewChaincodeQuery(chaincodeInvocationSpec, util.GenerateUUID())
	if err != nil {
		return nil, nil, fmt.Errorf("Error deploying chaincode: %s ", err)
	}

	resp, err = processTransaction(transaction)
	return
}
