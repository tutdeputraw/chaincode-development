package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"

	constant "tutdeputraw.com/common/constants"
	mock "tutdeputraw.com/common/mocks"
	"tutdeputraw.com/common/models"
)

//==========[CORE]==========//

func (s *RealEstateChaincode) RealEstate_Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	realEstates := mock.Mock_RealEstates

	i := 0
	realEstatesLen := len(realEstates)
	for i < realEstatesLen {
		realEstateAsBytes, _ := json.Marshal(realEstates[i])
		APIstub.PutState(constant.RealEstateKey+realEstates[i].RealEstateId, realEstateAsBytes)

		// write the user to the real estate ownership history
		compositeKey := constant.GetRealEstatesByOwnerKey
		ownerId := constant.UserKey + realEstates[i].OwnerId
		realEstateId := constant.RealEstateKey + realEstates[i].RealEstateId
		ownerRealEstateIdIndexKey, err := APIstub.CreateCompositeKey(
			compositeKey,
			[]string{ownerId, realEstateId},
		)
		if err != nil {
			return shim.Error(err.Error())
		}
		value := []byte{0x00}
		APIstub.PutState(ownerRealEstateIdIndexKey, value)

		i = i + 1
	}

	return shim.Success(nil)
}

func (s *RealEstateChaincode) RealEstate_CheckIfRealEstateHasAlreadyRegistered(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	bytes, _ := APIstub.GetState(constant.RealEstateKey + args[0])
	if bytes != nil {
		return shim.Success([]byte(strconv.FormatBool(true)))
	}

	return shim.Success([]byte(strconv.FormatBool(false)))
}

func (s *RealEstateChaincode) RealEstate_Create(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 12 {
		return shim.Error("Incorrect number of arguments. Expecting 12")
	}

	var realEstate = models.RealEstateModel{
		RealEstateId: args[0],
		OwnerId:      args[1],
		Price:        args[2],
		Bed:          args[3],
		Bath:         args[4],
		AcreLot:      args[5],
		FullAddress:  args[6],
		Street:       args[7],
		City:         args[8],
		State:        args[9],
		ZipCode:      args[10],
		HouseSize:    args[11],
	}

	realEstateAsBytes, _ := json.Marshal(realEstate)
	APIstub.PutState(constant.RealEstateKey+args[0], realEstateAsBytes)

	return shim.Success(realEstateAsBytes)
}

func (s *RealEstateChaincode) RealEstate_RegisterNewRealEstate(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 12 {
		return shim.Error("Incorrect number of arguments. Expecting 12")
	}

	var realEstate = models.RealEstateModel{
		RealEstateId: args[0],
		OwnerId:      args[1],
		Price:        args[2],
		Bed:          args[3],
		Bath:         args[4],
		AcreLot:      args[5],
		FullAddress:  args[6],
		Street:       args[7],
		City:         args[8],
		State:        args[9],
		ZipCode:      args[10],
		HouseSize:    args[11],
	}

	// return false if a user is not registered first
	responseIsUserRegistered := s.User_CheckIfUserExist(APIstub, []string{realEstate.OwnerId})
	userIsRegistered, _ := strconv.ParseBool(string(responseIsUserRegistered.Payload))
	userIsNotRegistered := !userIsRegistered
	if userIsNotRegistered {
		return shim.Error("user is not registered")
	}

	// return false if the real estate is already registered
	responseCheckIfRealEstateHasAlreadyRegistered := s.RealEstate_CheckIfRealEstateHasAlreadyRegistered(APIstub, args)
	realEstateIsRegistered, _ := strconv.ParseBool(string(responseCheckIfRealEstateHasAlreadyRegistered.Payload))
	if realEstateIsRegistered {
		return shim.Error("real estate is not registered")
	}

	// create real estate
	s.RealEstate_Create(APIstub, args)

	// create real estate history
	s.RealEstateHistory_Create(APIstub, []string{
		realEstate.RealEstateId, // id
		realEstate.OwnerId,
		realEstate.RealEstateId,
		time.Now().Local().String(),
	})

	// write the user to the real estate ownership history
	createCompositeRealEstatesByOwner := s.RealEstate_CreateCompositeRealEstatesByOwner(APIstub, realEstate)
	createCompositeRealEstatesByOwnerIsSuccess, _ := strconv.ParseBool(string(createCompositeRealEstatesByOwner.Payload))
	if !createCompositeRealEstatesByOwnerIsSuccess {
		return shim.Error("failed to create composite real estates by owner")
	}

	createCompositeOwnersByRealEstate := s.RealEstate_CreateCompositeOwnersByRealEstate(APIstub, realEstate)
	createCompositeOwnersByRealEstateIsSuccess, _ := strconv.ParseBool(string(createCompositeOwnersByRealEstate.Payload))
	if !createCompositeOwnersByRealEstateIsSuccess {
		return shim.Error("failed to create composite owners by real estate")
	}

	bytesRealEstate, _ := json.Marshal(realEstate)
	return shim.Success(bytesRealEstate)
}

func (s *RealEstateChaincode) RealEstate_QueryById(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	realEstateAsBytes, _ := APIstub.GetState(args[0])

	return shim.Success(realEstateAsBytes)
}

func (s *RealEstateChaincode) RealEstate_QueryAll(APIstub shim.ChaincodeStubInterface) sc.Response {
	startKey := constant.RealEstateKey + "0"
	endKey := constant.RealEstateKey + "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	results := []byte{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		results = append(results, queryResponse.Value...)
	}

	// fmt.Printf("- queryAllRealEstatese:\n%s\n", string(results))

	return shim.Success(results)

	// resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }
	// defer resultsIterator.Close()

	// // buffer is a JSON array containing QueryResults
	// var buffer bytes.Buffer
	// buffer.WriteString("[")

	// bArrayMemberAlreadyWritten := false
	// for resultsIterator.HasNext() {
	// 	queryResponse, err := resultsIterator.Next()
	// 	if err != nil {
	// 		return shim.Error(err.Error())
	// 	}
	// 	// Add a comma before array members, suppress it for the first array member
	// 	if bArrayMemberAlreadyWritten == true {
	// 		buffer.WriteString(",")
	// 	}
	// 	buffer.WriteString("{\"Key\":")
	// 	buffer.WriteString("\"")
	// 	buffer.WriteString(queryResponse.Key)
	// 	buffer.WriteString("\"")

	// 	buffer.WriteString(", \"Record\":")
	// 	// Record is a JSON object, so we write as-is
	// 	buffer.WriteString(string(queryResponse.Value))
	// 	buffer.WriteString("}")
	// 	bArrayMemberAlreadyWritten = true
	// }
	// buffer.WriteString("]")

	// fmt.Printf("- queryAllRealEstates:\n%s\n", buffer.String())

	// return shim.Success(buffer.Bytes())
}

func (s *RealEstateChaincode) RealEstate_QueryByOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments")
	}
	owner := args[0]

	ownerAndIdResultIterator, err := APIstub.GetStateByPartialCompositeKey(
		constant.GetRealEstatesByOwnerKey,
		[]string{constant.UserKey + owner},
	)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer ownerAndIdResultIterator.Close()

	var i int
	var id string

	var realestates []byte
	bArrayMemberAlreadyWritten := false

	realestates = append([]byte("["))

	for i = 0; ownerAndIdResultIterator.HasNext(); i++ {
		responseRange, err := ownerAndIdResultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		objectType, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}

		id = compositeKeyParts[1]
		assetAsBytes, err := APIstub.GetState(id)

		fmt.Println("INFOOOUSER")
		fmt.Println("compositeKeyParts\t: ", compositeKeyParts)
		fmt.Println("id\t\t\t: ", id)
		fmt.Println("assetAsBytes\t\t: \n", string(assetAsBytes))

		if bArrayMemberAlreadyWritten == true {
			newBytes := append([]byte(","), assetAsBytes...)
			realestates = append(realestates, newBytes...)

		} else {
			// newBytes := append([]byte(","), carsAsBytes...)
			realestates = append(realestates, assetAsBytes...)
		}

		fmt.Println("Found a asset\t: ", objectType)
		fmt.Println("for index\t: ", compositeKeyParts[0])
		fmt.Println("asset id\t: ", compositeKeyParts[1])
		bArrayMemberAlreadyWritten = true
	}

	realestates = append(realestates, []byte("]")...)

	return shim.Success(realestates)
}

func (s *RealEstateChaincode) RealEstate_CreateCompositeRealEstatesByOwner(APIstub shim.ChaincodeStubInterface, realEstate models.RealEstateModel) sc.Response {
	compositeKey := constant.GetRealEstatesByOwnerKey
	ownerId := constant.UserKey + realEstate.OwnerId
	realEstateId := constant.RealEstateKey + realEstate.RealEstateId

	ownerRealEstateIdIndexKey, err := APIstub.CreateCompositeKey(
		compositeKey,
		[]string{ownerId, realEstateId},
	)
	if err != nil {
		return shim.Success([]byte(strconv.FormatBool(false)))
	}

	value := []byte{0x00}
	fmt.Println("JJJA1: ", ownerRealEstateIdIndexKey)
	APIstub.PutState(ownerRealEstateIdIndexKey, value)

	return shim.Success([]byte(strconv.FormatBool(true)))
}

func (s *RealEstateChaincode) RealEstate_CreateCompositeOwnersByRealEstate(APIstub shim.ChaincodeStubInterface, realEstate models.RealEstateModel) sc.Response {
	compositeKey := constant.GetOwnersByRealEstateKey
	ownerId := constant.UserKey + realEstate.OwnerId
	realEstateId := constant.RealEstateKey + realEstate.RealEstateId
	realEstateHistoryKey := constant.RealEstateHistoryKey + realEstate.RealEstateId + realEstate.OwnerId

	ownerRealEstateIdIndexKey, err := APIstub.CreateCompositeKey(
		compositeKey,
		[]string{realEstateId, ownerId, realEstateHistoryKey},
	)
	if err != nil {
		return shim.Success([]byte(strconv.FormatBool(false)))
	}

	value := []byte{0x00}
	fmt.Println("JJJA2: ", ownerRealEstateIdIndexKey)
	APIstub.PutState(ownerRealEstateIdIndexKey, value)

	return shim.Success([]byte(strconv.FormatBool(true)))
}

func (s *RealEstateChaincode) RealEstate_ChangeRealEstateOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	realEstateKey := constant.RealEstateKey + args[0]
	newOwnerId := args[1]

	realEstateAsBytes, _ := APIstub.GetState(realEstateKey)
	realEstate := models.RealEstateModel{}
	json.Unmarshal(realEstateAsBytes, &realEstate)

	//==========[delete old key]==========//
	compositeKey := constant.GetRealEstatesByOwnerKey
	ownerId := constant.UserKey + realEstate.OwnerId
	realEstateId := constant.RealEstateKey + realEstate.RealEstateId

	ownerCaridIndexKey, err := APIstub.CreateCompositeKey(
		compositeKey,
		[]string{ownerId, realEstateId},
	)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = APIstub.DelState(ownerCaridIndexKey)
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}
	//----------[delete old key]----------//

	//==========[update realestate state]==========//
	realEstate.OwnerId = newOwnerId

	realEstateAsBytes, _ = json.Marshal(realEstate)

	APIstub.PutState(realEstateKey, realEstateAsBytes)
	//----------[update realestate state]----------//

	//==========[create real estate history]==========//
	s.RealEstateHistory_Create(APIstub, []string{
		realEstate.RealEstateId, // id
		realEstate.OwnerId,
		realEstate.RealEstateId,
		time.Now().Local().String(),
	})
	//----------[create real estate history]----------//

	//==========[add new user into the composite of real estate history]==========//
	createCompositeRealEstatesByOwner := s.RealEstate_CreateCompositeRealEstatesByOwner(APIstub, realEstate)
	createCompositeRealEstatesByOwnerIsSuccess, _ := strconv.ParseBool(string(createCompositeRealEstatesByOwner.Payload))
	if !createCompositeRealEstatesByOwnerIsSuccess {
		return shim.Error("failed to create composite real estates by owner")
	}

	createCompositeOwnersByRealEstate := s.RealEstate_CreateCompositeOwnersByRealEstate(APIstub, realEstate)
	createCompositeOwnersByRealEstateIsSuccess, _ := strconv.ParseBool(string(createCompositeOwnersByRealEstate.Payload))
	if !createCompositeOwnersByRealEstateIsSuccess {
		return shim.Error("failed to create composite owners by real estate")
	}
	//----------[add new user into the composite of real estate history]----------//

	return shim.Success(nil)
}