/*
Copyright DASE@ECNU. 2016 All Rights Reserved.
*/

package query

import (
	"fmt"
	"io/ioutil"

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
