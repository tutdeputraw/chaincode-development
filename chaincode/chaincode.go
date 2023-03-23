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
	case "QueryAssets":
		return s.QueryAssets(APIstub, args)
	// Real Estate
	case "RealEstate_Init":
		return s.RealEstate_Init(APIstub)
	case "RealEstate_Create":
		return s.RealEstate_Create(APIstub, args)
	case "RealEstate_QueryById":
		return s.RealEstate_QueryById(APIstub, args)
	case "RealEstate_QueryAll":
		return s.RealEstate_QueryAll(APIstub)
	case "RealEstate_CheckIfRealEstateHasAlreadyRegistered":
		return s.RealEstate_CheckIfRealEstateHasAlreadyRegistered(APIstub, args)
	case "RealEstate_RegisterNewRealEstate":
		return s.RealEstate_RegisterNewRealEstate(APIstub, args)
	case "RealEstate_QueryByOwner":
		return s.RealEstate_QueryByOwner(APIstub, args)
	case "RealEstate_ChangeRealEstateOwner":
		return s.RealEstate_ChangeRealEstateOwner(APIstub, args)
	case "RealEstate_ChangeRealEstateSellStatus":
		return s.RealEstate_ChangeRealEstateSellStatus(APIstub, args)

	// User
	case "User_Init":
		return s.User_Init(APIstub)
	case "User_CheckIfUserExist":
		return s.User_CheckIfUserExist(APIstub, args)
	case "User_QueryById":
		return s.User_QueryById(APIstub, args)
	case "User_QueryAll":
		return s.User_QueryAll(APIstub)
	case "User_Create":
		return s.User_Create(APIstub, args)
	case "NYOBAK":
		return s.NYOBAK(APIstub)

	// History
	case "RealEstateHistory_Create":
		return s.RealEstateHistory_Create(APIstub, args)
	case "RealEstateHistory_QueryByRealEstateId":
		return s.RealEstateHistory_QueryByRealEstateId(APIstub, args)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}
}
