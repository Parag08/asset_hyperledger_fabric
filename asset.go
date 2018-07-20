/*
 * Copyright Wavenet All Rights Reserved
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"
	"crypto/sha256"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"strconv"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Wallet implements a structure to hold wallet information
type wallet struct {
	Name string `json:"name"`
	Owner string `json:"owner"`
	Balance float64 `json:"balance"`
	Password string `json:"password"`
}

type stakeholders struct {
        Fraction float64 `json:"fraction"`
        WalletAdress string `json:"walletAddress"`
}

type asset struct {
        Name string `json:"name"`
        Password string `json:"password"`
        Owner string `json:"Owner"`
        OwnerInfo string `json:"OwnerInfo"`
        Stakeholders []stakeholders `json:"stakeholder"`
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the args from the transaction proposal

	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	if function == "createAsset" {
		return t.createAsset(stub,args)
	} else if function == "getAsset" {
		return t.getAsset(stub,args)
	}
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}


func convertStringToSHA256 (inputString string) string {
	sha_256 := sha256.New()
	sha_256.Write([]byte(inputString))
	return fmt.Sprintf("%X", sha_256.Sum(nil))
}

func (t *SimpleChaincode) Authenticate (walletInfo wallet,password string) bool {
	if (walletInfo.Password == convertStringToSHA256(password)) {
		return true
	}
	return false
}

func (t *SimpleChaincode) createAsset (stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// len of args should be 5
	// 0              1          2          4             5
	// assetName    password   ownerName  ownerInfo   stakeholderList
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3st argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4nd argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return shim.Error("5th argument must be a non-empty string")
	}
	assetAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get asset: " + err.Error())
	} else if assetAsBytes != nil {
		fmt.Println("asset already exists")
		return shim.Error("asset already exists")
	}
	assetName := args[0]
	password := convertStringToSHA256(args[1])
	ownerName := args[2]
	ownerInfo := args[3]
	stakeholdersList := []stakeholders{}
	json.Unmarshal([]byte(args[4]), &stakeholdersList)
	sum := 0.0
	for _ , elem := range stakeholdersList {
		sum = sum + elem.Fraction
	}
	fmt.Println("Sum of fraction is:",sum)
	if sum != 1.0 {
		return shim.Error("sum of fractional distribution among stakeholder should be 1")
	}
	assetData := &asset{assetName, password,ownerName, ownerInfo, stakeholdersList}
	assetJSONasBytes, err := json.Marshal(assetData)
	if err != nil {
		return shim.Error(err.Error())
	}
	// === Save wallet to state ===
	err = stub.PutState(assetName, assetJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}


func (t *SimpleChaincode) getAsset (stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// len of args should be 5
	// 0              1
	// assetName    password
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	assetAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get asset: " + err.Error())
	} else if assetAsBytes == nil {
		fmt.Println("asset not found")
		return shim.Error("asset doesn't exists")
	}
	assetInfo := asset{}
	err = json.Unmarshal(assetAsBytes, &assetInfo) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	if (assetInfo.Password == convertStringToSHA256(args[1])){
		return shim.Success(assetAsBytes)
	} else {
		return shim.Error("password is wrong")
	}
}


func (t *SimpleChaincode) buyAsset (stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// len of args should be 4
	// 0           1       2             3
	// assetName price customerWallet  password
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	assetAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get asset: " + err.Error())
	} else if assetAsBytes == nil {
		fmt.Println("asset not found")
		return shim.Error("asset doesn't exists")
	}
	price, err := strconv.ParseFloat(args[1],64)
	if err != nil {
		return shim.Error("3rd argument must be a numeric string")
	}
	password := args[3]
	assetInfo := asset{}
	err = json.Unmarshal(assetAsBytes, &assetInfo) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	walletAsBytes, err := stub.GetState(args[2])
	if err != nil {
		return shim.Error("Failed to get asset: " + err.Error())
	} else if walletAsBytes == nil {
		fmt.Println("asset not found")
		return shim.Error("asset doesn't exists")
	}
	walletInfo := wallet{}
	err = json.Unmarshal(walletAsBytes, &walletInfo) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	if t.Authenticate(walletInfo, password) {
		// subtract money from customer wallet
		// divide and add money to all wallets of stakeholders based on there ratio owned
		if (walletInfo.Balance >= price  && price > 0) {
			walletInfo.Balance = walletInfo.Balance - price
			stakeholdersList := assetInfo.Stakeholders
			for _ , elem := range stakeholdersList {
				walletAsBytes, err := stub.GetState(elem.WalletAdress)
				if err != nil {
					return shim.Error("Failed to get asset: " + err.Error())
				} else if walletAsBytes == nil {
					fmt.Println("asset not found")
					return shim.Error("asset doesn't exists")
				}
				walletInfo := wallet{}
				err = json.Unmarshal(walletAsBytes, &walletInfo) //unmarshal it aka JSON.parse()
				if err != nil {
					return shim.Error(err.Error())
				}
				walletInfo.Balance = walletInfo.Balance + elem.Fraction*price
				stakeholdersWalletJSONasBytes, err := json.Marshal(walletInfo)
				if err != nil {
					return shim.Error(err.Error())
				}
				// === Save wallet to state ===
				err = stub.PutState(walletInfo.Name, stakeholdersWalletJSONasBytes)
				if err != nil {
					return shim.Error(err.Error())
				}
			}
			customerWalletJSONasBytes, err := json.Marshal(walletInfo)
			if err != nil {
				return shim.Error(err.Error())
			}
			// === Save wallet to state ===
			err = stub.PutState(walletInfo.Name, customerWalletJSONasBytes)
			if err != nil {
				return shim.Error(err.Error())
			}
			return shim.Success(nil)
		} else {
			return shim.Error("insufficient balance to make this transaction")
		}
	} else {
		return shim.Error("password is wrong")
	}
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SimpleChaincode)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}

