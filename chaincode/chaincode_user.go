package chaincode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	constant "tutdeputraw.com/common/constants"
	mock "tutdeputraw.com/common/mocks"
	"tutdeputraw.com/common/models"
)

func (s *RealEstateChaincode) User_Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	users := mock.Mock_Users
	// id, err := cid.GetID(APIstub)
	// if err != nil {
	// 	fmt.Println("failed to get cid, err= ", err)
	// 	return shim.Error("")
	// }
	// fmt.Println("ikiloo: ", id)

	i := 0
	usersLen := len(users)
	for i < usersLen {
		realEstateAsBytes, _ := json.Marshal(users[i])
		APIstub.PutState(constant.State_User+users[i].Id, realEstateAsBytes)
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *RealEstateChaincode) User_CheckIfUserExist(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	fmt.Println("IDID: ", constant.State_User+args[0])

	bytes, _ := APIstub.GetState(constant.State_User + args[0])
	if bytes != nil {
		return shim.Success([]byte(strconv.FormatBool(true)))
	}

	return shim.Success([]byte(strconv.FormatBool(false)))
}

func (s *RealEstateChaincode) User_Create(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	// check if user has not registered yet
	responseIsUserRegistered := s.User_CheckIfUserExist(APIstub, []string{args[0]})
	userIsRegistered, _ := strconv.ParseBool(string(responseIsUserRegistered.Payload))
	if userIsRegistered {
		return shim.Error("user has already registered")
	}

	bytes, _ := APIstub.GetState(args[0])
	if bytes != nil {
		return shim.Error("user is already registered")
	}

	var user = models.UserModel{
		Id:          args[0],
		Name:        args[1],
		NPWPNumber:  args[2],
		PhoneNumber: args[3],
		Email:       args[4],
	}

	userAsBytes, _ := json.Marshal(user)
	APIstub.PutState(constant.State_User+args[0], userAsBytes)

	return shim.Success(userAsBytes)
}

func (s *RealEstateChaincode) User_QueryById(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	realEstateAsBytes, _ := APIstub.GetState(constant.State_User + args[0])
	// fmt.Printf("- queryRealEstateById:\n%s\n", string(realEstateAsBytes))
	return shim.Success(realEstateAsBytes)
}

func (s *RealEstateChaincode) User_QueryAll(APIstub shim.ChaincodeStubInterface) sc.Response {
	startKey := constant.State_User + "0"
	endKey := constant.State_User + "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	// fmt.Printf("- queryAllUsers:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
