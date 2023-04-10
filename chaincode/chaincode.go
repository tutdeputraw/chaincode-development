package chaincode

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	"tutdeputraw.com/common/helpers"
)

type RealEstateChaincode struct{}

func (s *RealEstateChaincode) Init(APIStub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *RealEstateChaincode) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()

	helpers.Logger.Infof("Function name\t: %v", function)
	helpers.Logger.Infof("Args length\t: %v", len(args))
	helpers.Logger.Infof("Args\t\t: %v", args)

	switch function {
	// Query
	case "Query":
		return s.Query(APIstub, args)
	case "QueryWithPagination":
		return s.QueryWithPagination(APIstub, args)

	// Real Estate
	case "RealEstate_Init":
		return s.RealEstate_Init(APIstub)
	case "RealEstate_Create":
		return s.RealEstate_Create(APIstub, args)
	case "RealEstate_GetById":
		return s.RealEstate_GetById(APIstub, args)
	case "RealEstate_GetAll":
		return s.RealEstate_GetAll(APIstub)
	case "RealEstate_GetByOwner":
		return s.RealEstate_GetByOwner(APIstub, args)
	case "RealEstate_CheckIfRealEstateHasAlreadyRegistered":
		return s.RealEstate_CheckIfRealEstateHasAlreadyRegistered(APIstub, args)
	case "RealEstate_RegisterNewRealEstate":
		return s.RealEstate_RegisterNewRealEstate(APIstub, args)
	case "RealEstate_ChangeRealEstateOwner":
		return s.RealEstate_ChangeRealEstateOwner(APIstub, args)
	case "RealEstate_ChangeRealEstateSellStatus":
		return s.RealEstate_ChangeRealEstateSellStatus(APIstub, args)

	// User
	case "User_Init":
		return s.User_Init(APIstub)
	case "User_CheckIfUserExist":
		return s.User_CheckIfUserExist(APIstub, args)
	case "User_GetById":
		return s.User_GetById(APIstub, args)
	case "User_GetAll":
		return s.User_GetAll(APIstub)
	case "User_Create":
		return s.User_Create(APIstub, args)
	case "NYOBAK":
		return s.NYOBAK(APIstub)

	// History
	case "RealEstateHistory_Create":
		return s.RealEstateHistory_Create(APIstub, args)
	case "RealEstateHistory_GetByRealEstateId":
		return s.RealEstateHistory_GetByRealEstateId(APIstub, args)

		// Real Estate Sales Record
	case "RealEstateSalesRecord_Create":
		return s.RealEstateSalesRecord_Create(APIstub, args)
	case "RealEstateSalesRecord_Delete":
		return s.RealEstateSalesRecord_Delete(APIstub, args)
	case "RealEstateSalesRecord_IncrementInterestUsers":
		return s.RealEstateSalesRecord_IncrementInterestUsers(APIstub, args)
	case "RealEstateSalesRecord_UpdateSalesPhase":
		return s.RealEstateSalesRecord_UpdateSalesPhase(APIstub, args)
	case "RealEstateSalesRecord_UpdateRealEstateAssessment":
		return s.RealEstateSalesRecord_UpdateRealEstateAssessment(APIstub, args)
	case "RealEstateSalesRecord_GetByRealEstateIdComposite":
		return s.RealEstateSalesRecord_GetByRealEstateIdComposite(APIstub, args)

	default:
		return shim.Error("Invalid Smart Contract function name.")
	}
}
