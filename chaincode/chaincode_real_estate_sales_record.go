package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	constant "tutdeputraw.com/common/constants"
	"tutdeputraw.com/common/models"
)

func (s *RealEstateChaincode) RealEstateSalesRecord_Create(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	var realEstateSalesRecord = models.RealEstateSalesRecordModel{
		RealEstateId: args[0],
		SellerId:     args[1],
	}

	realEstateSalesRecordAsBytes, _ := json.Marshal(realEstateSalesRecord)
	realEstateSalesRecordKey := constant.State_RealEstateSalesRecord + realEstateSalesRecord.RealEstateId + realEstateSalesRecord.SellerId
	APIstub.PutState(
		realEstateSalesRecordKey,
		realEstateSalesRecordAsBytes,
	)

	return shim.Success(realEstateSalesRecordAsBytes)
}

func (s *RealEstateChaincode) RealEstateSalesRecord_Delete(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	realEstateSalesRecordId := args[0]

	err := APIstub.DelState(realEstateSalesRecordId)
	if err != nil {
		return shim.Error("failed to delete state: " + err.Error())
	}

	return shim.Success(nil)
}

func (s *RealEstateChaincode) RealEstateSalesRecord_IncrementInterestUsers(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	realEstateSalesRecordId := constant.State_RealEstateSalesRecord + args[0]

	// get the realEstateSalesRecord
	realEstateSalesRecordAsBytes, _ := APIstub.GetState(realEstateSalesRecordId)
	realEstateSalesRecord := models.RealEstateSalesRecordModel{}
	json.Unmarshal(realEstateSalesRecordAsBytes, &realEstateSalesRecord)

	// convert the interestUsers value to int
	interestUser, error := strconv.Atoi(realEstateSalesRecord.InterestUsers)
	if error != nil {
		fmt.Println("failed to convert string to int")
	}

	// increment the value
	incrementedValue := interestUser + 1
	realEstateSalesRecord.InterestUsers = strconv.Itoa(incrementedValue)

	realEstateSalesRecordAsBytes, _ = json.Marshal(realEstateSalesRecord)

	APIstub.PutState(realEstateSalesRecordId, realEstateSalesRecordAsBytes)

	return shim.Success(nil)
}

func (s *RealEstateChaincode) RealEstateSalesRecord_UpdateSalesPhase(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	realEstateSalesRecordId := constant.State_RealEstateSalesRecord + args[0]
	phase := args[1]

	// get the realEstateSalesRecord
	realEstateSalesRecordAsBytes, _ := APIstub.GetState(realEstateSalesRecordId)
	realEstateSalesRecord := models.RealEstateSalesRecordModel{}
	json.Unmarshal(realEstateSalesRecordAsBytes, &realEstateSalesRecord)

	realEstateSalesRecord.Phase = phase

	realEstateSalesRecordAsBytes, _ = json.Marshal(realEstateSalesRecord)

	APIstub.PutState(realEstateSalesRecordId, realEstateSalesRecordAsBytes)

	return shim.Success(nil)
}

func (s *RealEstateChaincode) RealEstateSalesRecord_UpdateRealEstateAssessment(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	realEstateSalesRecordId := constant.State_RealEstateSalesRecord + args[0]
	realEstateAssessment := args[1]

	// get the realEstateSalesRecord
	realEstateSalesRecordAsBytes, _ := APIstub.GetState(realEstateSalesRecordId)
	realEstateSalesRecord := models.RealEstateSalesRecordModel{}
	json.Unmarshal(realEstateSalesRecordAsBytes, &realEstateSalesRecord)

	realEstateSalesRecord.Phase = realEstateAssessment

	realEstateSalesRecordAsBytes, _ = json.Marshal(realEstateSalesRecord)

	APIstub.PutState(realEstateSalesRecordId, realEstateSalesRecordAsBytes)

	return shim.Success(nil)
}

func (s *RealEstateChaincode) RealEstateSalesRecord_GetByRealEstateIdComposite(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments")
	}
	realEstateId := args[0]

	resultsIterator, err := APIstub.GetStateByPartialCompositeKey(
		constant.Composite_GetRealEstateSalesRecordByRealEstateKey,
		[]string{constant.State_RealEstate + realEstateId},
	)
	if err != nil {
		return shim.Error("failed to get resultsIterator: " + err.Error())
	}
	defer resultsIterator.Close()

	var i int
	var id string

	var realestatesRecord []byte
	realestatesRecord = append([]byte("["))

	bArrayMemberAlreadyWritten := false

	for i = 0; resultsIterator.HasNext(); i++ {
		responseRange, err := resultsIterator.Next()
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
			realestatesRecord = append(realestatesRecord, newBytes...)
		} else {
			// newBytes := append([]byte(","), carsAsBytes...)
			realestatesRecord = append(realestatesRecord, assetAsBytes...)
		}

		fmt.Println("Found a asset\t: ", objectType)
		fmt.Println("for index\t: ", compositeKeyParts[0])
		fmt.Println("asset id\t: ", compositeKeyParts[1])
		bArrayMemberAlreadyWritten = true
	}

	realestatesRecord = append(realestatesRecord, []byte("]")...)

	return shim.Success(realestatesRecord)
}

func (s *RealEstateChaincode) RealEstateSalesRecord_DeleteComposite(APIstub shim.ChaincodeStubInterface, realEstate models.RealEstateModel) sc.Response {
	compositeKey := constant.Composite_GetRealEstateSalesRecordByRealEstateKey
	realEstateId := constant.State_RealEstate + realEstate.RealEstateId
	realEstateSalesRecordId := constant.State_RealEstateSalesRecord + realEstate.RealEstateId + realEstate.OwnerId

	realEstateSalesRecordIndexKey, err := APIstub.CreateCompositeKey(
		compositeKey,
		[]string{realEstateId, realEstateSalesRecordId},
	)
	if err != nil {
		return shim.Success([]byte(strconv.FormatBool(false)))
	}

	// check the availbility of state
	// state, _ := APIstub.GetState(combinedCompositeKey)
	// fmt.Println("SAKJANE1: " + combinedCompositeKey)
	// if state == nil {
	// 	fmt.Println("GAONO")
	// 	return shim.Success([]byte(strconv.FormatBool(true)))
	// }

	err = APIstub.DelState(realEstateSalesRecordIndexKey)
	if err != nil {
		return shim.Success([]byte(strconv.FormatBool(false)))
	}

	return shim.Success([]byte(strconv.FormatBool(true)))
}

func (s *RealEstateChaincode) RealEstateSalesRecord_CreateComposite(APIstub shim.ChaincodeStubInterface, realEstate models.RealEstateModel) sc.Response {
	compositeKey := constant.Composite_GetRealEstateSalesRecordByRealEstateKey
	realEstateId := constant.State_RealEstate + realEstate.RealEstateId
	realEstateSalesRecordId := constant.State_RealEstateSalesRecord + realEstate.RealEstateId + realEstate.OwnerId

	realEstateSalesRecordIndexKey, err := APIstub.CreateCompositeKey(
		compositeKey,
		[]string{realEstateId, realEstateSalesRecordId},
	)
	if err != nil {
		return shim.Success([]byte(strconv.FormatBool(false)))
	}

	value := []byte{0x00}
	err = APIstub.PutState(realEstateSalesRecordIndexKey, value)
	if err != nil {
		return shim.Success([]byte(strconv.FormatBool(false)))
	}

	kk, _ := APIstub.GetState(realEstateSalesRecordIndexKey)
	fmt.Println("SAKJANE0: " + realEstateSalesRecordIndexKey)
	if kk == nil {
		fmt.Println("GAONO JILID 2")
	}

	return shim.Success([]byte(strconv.FormatBool(true)))
}
