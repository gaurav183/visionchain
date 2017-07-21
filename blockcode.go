/*
Copyright 2016 IBM

Licensed under the Apache License, Version 2.0 (the "License")
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Licensed Materials - Property of IBM
© Copyright IBM Corp. 2016
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var cpPrefix = "cp:"
var accountPrefix = "acct:"

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}



const (
	millisPerSecond = int64(time.Second / time.Millisecond)
	nanosPerMillisecond = int64(time.Millisecond / time.Nanosecond)
)

func msToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(msInt / millisPerSecond,
		(msInt % millisPerSecond) * nanosPerMillisecond), nil
}

type video_frame struct {
	Hash  string    `json:"Hash"`
	Timecode string      `json:"Timecode"`
}

type footage struct {
	vID     string  `json:"vID"`
	Owner    Account  `json:"owner"`
	Frames    []video_frame `json:"frames"`
}

type Account struct {
	ID          string  `json:"id"`
	Name      string  `json:"name"`
	AssetsIds   []string `json:"assetIds"`
}



func (t *SimpleChaincode) createAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Creating account")

	// Obtain the username to associate with the account
	if len(args) != 2 {
		fmt.Println("Error obtaining username")
		return nil, errors.New("createAccount accepts a single username argument")
	}
	username := args[0]
	fullname := args[1]

	// Build an account object for the user
	var assetIds []string
	var account = Account{ID: username, Name: fullname, AssetsIds: assetIds}
	accountBytes, err := json.Marshal(&account)
	if err != nil {
		fmt.Println("error creating account" + account.ID)
		return nil, errors.New("Error creating account " + account.ID)
	}

	fmt.Println("Attempting to get state of any existing account for " + account.ID)
	existingBytes, err := stub.GetState(account.ID)
	if err == nil {

		var company Account
		err = json.Unmarshal(existingBytes, &company)
		if err != nil {
			fmt.Println("Error unmarshalling account " + account.ID + "\n--->: " + err.Error())

			if strings.Contains(err.Error(), "unexpected end") {
				fmt.Println("No data means existing account found for " + account.ID + ", initializing account.")
				err = stub.PutState(accountPrefix + account.ID, accountBytes)

				if err == nil {
					fmt.Println("created account" + account.ID)
					return nil, nil
				} else {
					fmt.Println("failed to create initialize account for " + account.ID)
					return nil, errors.New("failed to initialize an account for " + account.ID + " => " + err.Error())
				}
			} else {
				return nil, errors.New("Error unmarshalling existing account " + account.ID)
			}
		} else {
			fmt.Println("Account already exists for " + account.ID + " " + company.ID)
			return nil, errors.New("Can't reinitialize existing user " + account.ID)
		}
	} else {

		fmt.Println("No existing account found for " + account.ID + ", initializing account.")
		err = stub.PutState( account.ID, accountBytes)

		if err == nil {
			fmt.Println("created account" + account.ID)
			return nil, nil
		} else {
			fmt.Println("failed to create initialize account for " + account.ID)
			return nil, errors.New("failed to initialize an account for " + account.ID + " => " + err.Error())
		}

	}

}


func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Init firing. Function will be ignored: " + function)

	// Initialize the collection of commercial paper keys
	fmt.Println("Initializing paper keys collection")
	var blank []string
	blankBytes, _ := json.Marshal(&blank)
	err := stub.PutState("FootageKeys", blankBytes)
	if err != nil {
		fmt.Println("Failed to initialize paper key collection")
	}

	fmt.Println("Initialization complete")
	return nil, nil
}

func (t *SimpleChaincode) createNewFootage(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Creating commercial paper")

	/*		0
		json
	  	{
			"ticker":  "string",
			"par": 0.00,
			"qty": 10,
			"discount": 7.5,
			"maturity": 30,
			"owners": [ // This one is not required
				{
					"company": "company1",
					"quantity": 5
				},
				{
					"company": "company3",
					"quantity": 3
				},
				{
					"company": "company4",
					"quantity": 2
				}
			],				
			"issuer":"company2",
			"issueDate":"1456161763790"  (current time in milliseconds as a string)

		}
	*/
	//need one arg
	if len(args) != 1 {
		fmt.Println("error invalid arguments")
		return nil, errors.New("Incorrect number of arguments. Expecting commercial paper record")
	}

	var newfootage footage
	var err error
	var account Account

	fmt.Println("Unmarshalling Footage")
	err = json.Unmarshal([]byte(args[0]), &account)
	if err != nil {
		fmt.Println("error invalid footage issue")
		return nil, errors.New("Invalid footage issue")
	}
	newfootage.Owner = account.ID
	account.AssetsIds = append(account.AssetsIds, footage.vID)

	// Set the issuer to be the owner of all quantity
	var videoframe video_frame

	newfootage.Frames = append(newfootage.Frames, videoframe)

	fmt.Println("Getting State on CP " + newfootage.vID)
	cpRxBytes, err := stub.GetState(newfootage.vID)
	if cpRxBytes == nil {
		fmt.Println("vID does not exist, creating it")
		cpBytes, err := json.Marshal(&cp)
		if err != nil {
			fmt.Println("Error marshalling foortage")
			return nil, errors.New("Error issuing footage")
		}
		err = stub.PutState(newfootage.vID, cpBytes)
		if err != nil {
			fmt.Println("Error issuing footage")
			return nil, errors.New("Error issuing footage")
		}

		fmt.Println("Marshalling account bytes to write")
		accountBytesToWrite, err := json.Marshal(&account)
		if err != nil {
			fmt.Println("Error marshalling account")
			return nil, errors.New("Error issuing footage")
		}
		err = stub.PutState( newfootage.Owner, accountBytesToWrite)
		if err != nil {
			fmt.Println("Error putting state on accountBytesToWrite")
			return nil, errors.New("Error issuing commercial paper")
		}


		// Update the paper keys by adding the new key
		fmt.Println("Getting Paper Keys")
		keysBytes, err := stub.GetState("PaperKeys")
		if err != nil {
			fmt.Println("Error retrieving paper keys")
			return nil, errors.New("Error retrieving paper keys")
		}
		var keys []string
		err = json.Unmarshal(keysBytes, &keys)
		if err != nil {
			fmt.Println("Error unmarshel keys")
			return nil, errors.New("Error unmarshalling paper keys ")
		}

		fmt.Println("Appending the new key to Paper Keys")
		foundKey := false
		for _, key := range keys {
			if key == cpPrefix + cp.CUSIP {
				foundKey = true
			}
		}
		if foundKey == false {
			keys = append(keys, cpPrefix + cp.CUSIP)
			keysBytesToWrite, err := json.Marshal(&keys)
			if err != nil {
				fmt.Println("Error marshalling keys")
				return nil, errors.New("Error marshalling the keys")
			}
			fmt.Println("Put state on PaperKeys")
			err = stub.PutState("PaperKeys", keysBytesToWrite)
			if err != nil {
				fmt.Println("Error writting keys back")
				return nil, errors.New("Error writing the keys back")
			}
		}

		fmt.Println("Issue commercial paper %+v\n", cp)
		return nil, nil
	} else {
		fmt.Println("CUSIP exists")

		var cprx CP
		fmt.Println("Unmarshalling CP " + cp.CUSIP)
		err = json.Unmarshal(cpRxBytes, &cprx)
		if err != nil {
			fmt.Println("Error unmarshalling cp " + cp.CUSIP)
			return nil, errors.New("Error unmarshalling cp " + cp.CUSIP)
		}

		cprx.Qty = cprx.Qty + cp.Qty

		for key, val := range cprx.Owners {
			if val.Company == cp.Issuer {
				cprx.Owners[key].Quantity += cp.Qty
				break
			}
		}

		cpWriteBytes, err := json.Marshal(&cprx)
		if err != nil {
			fmt.Println("Error marshalling cp")
			return nil, errors.New("Error issuing commercial paper")
		}
		err = stub.PutState(cpPrefix + cp.CUSIP, cpWriteBytes)
		if err != nil {
			fmt.Println("Error issuing paper")
			return nil, errors.New("Error issuing commercial paper")
		}

		fmt.Println("Updated commercial paper %+v\n", cprx)
		return nil, nil
	}
}

func getAllFootage(stub shim.ChaincodeStubInterface) ([]CP, error) {

	var allFootage []footage

	// Get list of all the keys
	keysBytes, err := stub.GetState("PaperKeys") //get keys of this account's footages
	if err != nil {
		fmt.Println("Error retrieving paper keys")
		return nil, errors.New("Error retrieving paper keys")
	}
	var keys []string
	err = json.Unmarshal(keysBytes, &keys)
	if err != nil {
		fmt.Println("Error unmarshalling paper keys")
		return nil, errors.New("Error unmarshalling paper keys")
	}

	// Get all the cps
	for _, value := range keys {
		cpBytes, err := stub.GetState(value)

		var feet footage
		err = json.Unmarshal(cpBytes, &feet) //gross
		if err != nil {
			fmt.Println("Error retrieving footage " + value)
			return nil, errors.New("Error retrieving cp " + value)
		}

		fmt.Println("Appending CP" + value)
		allFootage = append(allFootage, feet)
	}

	return allFootage, nil
}

func getFootage(cpid string, stub shim.ChaincodeStubInterface) (CP, error) {
	var feet footage //here feet is a single footage object

	cpBytes, err := stub.GetState(cpid) //here 'cpid' refers to vID field of footage
	if err != nil {
		fmt.Println("Error retrieving footage " + cpid)
		return cp, errors.New("Error retrieving footage " + cpid)
	}

	err = json.Unmarshal(cpBytes, &feet)
	if err != nil {
		fmt.Println("Error unmarshalling footage " + cpid)
		return cp, errors.New("Error unmarshalling cp " + cpid)
	}

	return cp, nil
}

func getAccount(companyID string, stub shim.ChaincodeStubInterface) (Account, error) {
	var shooter Account  //shooter of footgage / account holder
	companyBytes, err := stub.GetState(companyID)
	if err != nil {
		fmt.Println("Account not found " + companyID)
		return company, errors.New("Account not found " + companyID)
	}

	err = json.Unmarshal(companyBytes, &company)
	if err != nil {
		fmt.Println("Error unmarshalling account " + companyID + "\n err:" + err.Error())
		return company, errors.New("Error unmarshalling account " + companyID)
	}

	return company, nil
}



func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Query running. Function: " + function)

	if function == "GetAllFootage" {
		fmt.Println("Getting all Footages")
		allCPs, err := getAllFootage(stub)
		if err != nil {
			fmt.Println("Error from getAllFootage")
			return nil, err
		} else {
			allCPsBytes, err1 := json.Marshal(&allCPs)
			if err1 != nil {
				fmt.Println("Error marshalling allcps")
				return nil, err1
			}
			fmt.Println("All success, returning allcps")
			return allCPsBytes, nil
		}
	} else if function == "GetFootage" {
		fmt.Println("Getting particular footage")
		cp, err := getFootage(args[0], stub)
		if err != nil {
			fmt.Println("Error Getting particular footage")
			return nil, err
		} else {
			cpBytes, err1 := json.Marshal(&cp)
			if err1 != nil {
				fmt.Println("Error marshalling the footage")
				return nil, err1
			}
			fmt.Println("All success, returning the footage")
			return cpBytes, nil
		}
	} else if function == "GetAccount" {
		fmt.Println("Getting the account")
		company, err := getAccount(args[0], stub)
		if err != nil {
			fmt.Println("Error from getAccount")
			return nil, err
		} else {
			companyBytes, err1 := json.Marshal(&company)
			if err1 != nil {
				fmt.Println("Error marshalling the account")
				return nil, err1
			}
			fmt.Println("All success, returning the account")
			return companyBytes, nil
		}
	} else {
		fmt.Println("Generic Query call")
		bytes, err := stub.GetState(args[0])

		if err != nil {
			fmt.Println("Some error happenend: " + err.Error())
			return nil, err
		}

		fmt.Println("All success, returning from generic")
		return bytes, nil
	}
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Invoke running. Function: " + function)

	if function == "craeteNewFootage" {
		return t.craeteNewFootage(stub, args)
	} else if function == "createAccount" {
		return t.createAccount(stub, args)
	}

	return nil, errors.New("Received unknown function invocation: " + function)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Println("Error starting Simple chaincode: %s", err)
	}
}

