package main

import (
	"fmt"
	"net"
	"testing"
	"time"

	//"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric/core/chaincode"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/hyperledger/fabric/core/chaincode/shim/crypto/attr"
	//"github.com/hyperledger/fabric/core/container"
	"github.com/hyperledger/fabric/core/crypto"
	//"github.com/hyperledger/fabric/core/ledger"
	"github.com/hyperledger/fabric/membersrvc/ca"
	pb "github.com/hyperledger/fabric/protos"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	//"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	//"reflect"
)

const (
	chaincodeStartupTimeoutDefault int = 5000
)

var (
	testLogger = logging.MustGetLogger("test")

	lis net.Listener

	citi    crypto.Client
	auditor crypto.Client
	bermuda crypto.Client
	art     crypto.Client
	allianz crypto.Client
	nigeria crypto.Client


	server *grpc.Server
	aca    *ca.ACA
	eca    *ca.ECA
	tca    *ca.TCA
	tlsca  *ca.TLSCA
)

func TestMain(m *testing.M) {
	setup()

	go initMembershipSrvc()
	fmt.Println("Wait for some secs for OBCCA")
	time.Sleep(2 * time.Second)

	go initVP()
	fmt.Println("Wait for some secs for VP")
	time.Sleep(2 * time.Second)

	go initInsuranceFrontingChaincode()
	fmt.Println("Wait for some secs for Chaincode")
	time.Sleep(2 * time.Second)

	if err := initClients(); err != nil {
		panic(err)
	}

	fmt.Println("Wait for 10 secs for chaincode to be started")
	time.Sleep(10 * time.Second)

	ret := m.Run()

	lis.Close()
	fmt.Println("Wait for some secs for Listener to close")
	time.Sleep(2 * time.Second)

	defer removeFolders()
	os.Exit(ret)
}

func TestInsuranceFronting(t *testing.T) {

}

func initClients() error {
	// Citi
	if err := crypto.RegisterClient("citi", nil, "citi", "4nXSrfoYGFCP"); err != nil {
		return err
	}
	var err error
	citi, err = crypto.InitClient("citi", nil)
	if err != nil {
		return err
	}

	// Auditor
	if err := crypto.RegisterClient("auditor", nil, "auditor", "yg5DVhm0er1z"); err != nil {
		return err
	}
	auditor, err = crypto.InitClient("auditor", nil)
	if err != nil {
		return err
	}

	// Bermuda
	if err := crypto.RegisterClient("bermuda", nil, "bermuda", "b7pmSxzKNFiw"); err != nil {
		return err
	}
	bermuda, err = crypto.InitClient("bermuda", nil)
	if err != nil {
		return err
	}

	// Art
	if err := crypto.RegisterClient("art", nil, "art", "YsWZD4qQmYxo"); err != nil {
		return err
	}
	art, err = crypto.InitClient("art", nil)
	if err != nil {
		return err
	}

	// Allianz
	if err := crypto.RegisterClient("allianz", nil, "allianz", "W8G0usrU7jRk"); err != nil {
		return err
	}
	allianz, err = crypto.InitClient("allianz", nil)
	if err != nil {
		return err
	}

	// Nigeria
	if err := crypto.RegisterClient("nigeria", nil, "nigeria", "H80SiB5ODKKQ"); err != nil {
		return err
	}
	allianz, err = crypto.InitClient("nigeria", nil)
	if err != nil {
		return err
	}

	return nil
}

func initInsuranceFrontingChaincode() {
	err := shim.Start(new(InsuranceFrontingChaincode))
	if err != nil {
		panic(err)
	}
}

func initVP() {
	var opts []grpc.ServerOption
	if viper.GetBool("peer.tls.enabled") {
		creds, err := credentials.NewServerTLSFromFile(viper.GetString("peer.tls.cert.file"), viper.GetString("peer.tls.key.file"))
		if err != nil {
			grpclog.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)

	//lis, err := net.Listen("tcp", viper.GetString("peer.address"))

	//use a different address than what we usually use for "peer"
	//we override the peerAddress set in chaincode_support.go
	peerAddress := "0.0.0.0:40404"
	var err error
	lis, err = net.Listen("tcp", peerAddress)
	if err != nil {
		return
	}

	getPeerEndpoint := func() (*pb.PeerEndpoint, error) {
		return &pb.PeerEndpoint{ID: &pb.PeerID{Name: "testpeer"}, Address: peerAddress}, nil
	}

	ccStartupTimeout := time.Duration(chaincodeStartupTimeoutDefault) * time.Millisecond
	userRunsCC := true

	// Install security object for peer
	var secHelper crypto.Peer
	if viper.GetBool("security.enabled") {
		enrollID := viper.GetString("security.enrollID")
		enrollSecret := viper.GetString("security.enrollSecret")
		var err error

		if viper.GetBool("peer.validator.enabled") {
			testLogger.Debugf("Registering validator with enroll ID: %s", enrollID)
			if err = crypto.RegisterValidator(enrollID, nil, enrollID, enrollSecret); nil != err {
				panic(err)
			}
			testLogger.Debugf("Initializing validator with enroll ID: %s", enrollID)
			secHelper, err = crypto.InitValidator(enrollID, nil)
			if nil != err {
				panic(err)
			}
		} else {
			testLogger.Debugf("Registering non-validator with enroll ID: %s", enrollID)
			if err = crypto.RegisterPeer(enrollID, nil, enrollID, enrollSecret); nil != err {
				panic(err)
			}
			testLogger.Debugf("Initializing non-validator with enroll ID: %s", enrollID)
			secHelper, err = crypto.InitPeer(enrollID, nil)
			if nil != err {
				panic(err)
			}
		}
	}

	pb.RegisterChaincodeSupportServer(grpcServer,
		chaincode.NewChaincodeSupport(chaincode.DefaultChain, getPeerEndpoint, userRunsCC,
			ccStartupTimeout, secHelper))

	grpcServer.Serve(lis)
}

func initMembershipSrvc() {
	ca.LogInit(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr, os.Stdout)

	aca = ca.NewACA()
	eca = ca.NewECA()
	tca = ca.NewTCA(eca)
	tlsca = ca.NewTLSCA(eca)

	var opts []grpc.ServerOption
	if viper.GetBool("peer.pki.tls.enabled") {
		// TLS configuration
		creds, err := credentials.NewServerTLSFromFile(
			filepath.Join(viper.GetString("server.rootpath"), "tlsca.cert"),
			filepath.Join(viper.GetString("server.rootpath"), "tlsca.priv"),
		)
		if err != nil {
			panic("Failed creating credentials for Membersrvc: " + err.Error())
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	fmt.Printf("open socket...\n")
	sockp, err := net.Listen("tcp", viper.GetString("server.port"))
	if err != nil {
		panic("Cannot open port: " + err.Error())
	}
	fmt.Printf("open socket...done\n")

	server = grpc.NewServer(opts...)

	aca.Start(server)
	eca.Start(server)
	tca.Start(server)
	tlsca.Start(server)

	fmt.Printf("start serving...\n")
	server.Serve(sockp)
}

func setup() {
	// Conf
	viper.SetConfigName("test") // name of config file (without extension)
	viper.AddConfigPath(".")     // path to look for the config file in
	err := viper.ReadInConfig()  // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file [%s] \n", err))
	}

	// Logging
	var formatter = logging.MustStringFormatter(
		`%{color}[%{module}] %{shortfunc} [%{shortfile}] -> %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	logging.SetFormatter(formatter)

	logging.SetLevel(logging.DEBUG, "peer")
	logging.SetLevel(logging.DEBUG, "chaincode")
	logging.SetLevel(logging.DEBUG, "cryptochain")

	// Init the crypto layer
	if err := crypto.Init(); err != nil {
		panic(fmt.Errorf("Failed initializing the crypto layer [%s]", err))
	}

	removeFolders()
}

func removeFolders() {
	if err := os.RemoveAll(viper.GetString("peer.fileSystemPath")); err != nil {
		fmt.Printf("Failed removing [%s] [%s]\n", "hyperledger", err)
	}
}
