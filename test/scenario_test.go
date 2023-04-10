package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	cc "tutdeputraw.com/chaincode"
	helper "tutdeputraw.com/common/helpers"
	mock "tutdeputraw.com/common/mocks"
	"tutdeputraw.com/common/models"
)

func Test_registerNewUser(t *testing.T) {
	cc := new(cc.RealEstateChaincode)
	stub := shimtest.NewMockStub("real_estate", cc)

	expect := models.UserModel{
		Id:          "0",
		Name:        "username0",
		NPWPNumber:  "usernpwpnumber0",
		PhoneNumber: "phonenumber0",
		Email:       "email0",
	}

	result := helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("User_Create"),
		[]byte(expect.Id),
		[]byte(expect.Name),
		[]byte(expect.NPWPNumber),
		[]byte(expect.PhoneNumber),
		[]byte(expect.Email),
	})
	resultInModel := models.UserModel{}
	err1 := json.Unmarshal(result, &resultInModel)
	if err1 != nil {
		panic(err1)
	}
	if resultInModel != expect {
		t.FailNow()
		return
	}

	getUserById := helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("User_GetById"),
		[]byte(expect.Id),
	})
	getUserByIdInModel := models.UserModel{}
	err2 := json.Unmarshal(getUserById, &getUserByIdInModel)
	if err2 != nil {
		panic(err2)
	}
	if getUserByIdInModel != expect {
		t.FailNow()
		return
	}
}

func Test_registerNewRealEstate(t *testing.T) {
	cc := new(cc.RealEstateChaincode)
	stub := shimtest.NewMockStub("real_estate", cc)

	expect := mock.Mock_RealEstates_Owner1

	helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("User_Init"),
	})

	for _, v := range expect {
		helper.Test_CheckInvoke(t, stub, [][]byte{
			[]byte("RealEstate_RegisterNewRealEstate"),
			[]byte(v.RealEstateId), // realestate id
			[]byte(v.OwnerId),      // user id
			[]byte(v.Price),
			[]byte(v.Bed),
			[]byte(v.Bath),
			[]byte(v.AcreLot),
			[]byte(v.FullAddress),
			[]byte(v.Street),
			[]byte(v.City),
			[]byte(v.State),
			[]byte(v.ZipCode),
			[]byte(v.HouseSize),
			[]byte(v.IsOpenToSell),
		})
	}

	result := helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstate_GetByOwner"),
		[]byte("1"),
	})
	resultInModel := []models.RealEstateModel{}
	err := json.Unmarshal(result, &resultInModel)
	if err != nil {
		panic(err)
	}

	for i, v := range expect {
		if v != resultInModel[i] {
			t.FailNow()
			return
		}
	}
}

func Test_queryRealEstateByOwner(t *testing.T) {
	cc := new(cc.RealEstateChaincode)
	stub := shimtest.NewMockStub("real_estate", cc)
	expect := mock.Mock_RealEstates_Owner1

	helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("User_Init"),
	})
	helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstate_Init"),
	})

	queryResult := helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstate_GetByOwner"),
		[]byte("1"),
	})
	queryResultInModel := []models.RealEstateModel{}
	err := json.Unmarshal(queryResult, &queryResultInModel)
	if err != nil {
		panic(err)
	}

	for i, v := range expect {
		if v != queryResultInModel[i] {
			t.FailNow()
			return
		}
	}
}

func Test_BuyRealEstate(t *testing.T) {
	cc := new(cc.RealEstateChaincode)
	stub := shimtest.NewMockStub("real_estate", cc)

	//==========[init user]==========//
	helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("User_Init"),
	})

	queryResult := helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("User_GetById"),
		[]byte("3"),
	})
	queryResultInModelA := models.UserModel{}
	json.Unmarshal(queryResult, &queryResultInModelA)

	if queryResultInModelA != mock.Mock_Users[3] {
		t.Error("user_init not equal")
	}

	//----------[init user]----------//

	//==========[real estate 3 should have one owner history]==========//
	// mock := mock.Mock_RealEstates_TransactionHistory

	helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstate_Init"),
	})

	// for _, v := range mock {
	// 	helper.Test_CheckInvoke(t, stub, [][]byte{
	// 		[]byte("RealEstate_RegisterNewRealEstate"),
	// 		[]byte(v.RealEstateId), // realestate id
	// 		[]byte(v.OwnerId),      // user id
	// 		[]byte(v.Price),
	// 		[]byte(v.Bed),
	// 		[]byte(v.Bath),
	// 		[]byte(v.AcreLot),
	// 		[]byte(v.FullAddress),
	// 		[]byte(v.Street),
	// 		[]byte(v.City),
	// 		[]byte(v.State),
	// 		[]byte(v.ZipCode),
	// 		[]byte(v.HouseSize),
	// 		[]byte(v.IsOpenToSell),
	// 	})
	// }

	expect := []models.RealEstateHistoryModel{
		{
			OwnerID:      "2",
			RealEstateId: "3",
		},
	}

	queryResult = helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstateHistory_GetByRealEstateId"),
		[]byte("3"),
	})
	queryResultInModel := []models.RealEstateHistoryModel{}
	json.Unmarshal(queryResult, &queryResultInModel)

	for i, v := range expect {
		if v != queryResultInModel[i] {
			t.Error("= expect & queryResultInModel dont have the same value")
		}
	}
	//----------[real estate 3 should have one owner history]----------//

	//==========[user 2 should have a real estate]==========//
	queryResult = helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstate_GetByOwner"),
		[]byte("2"), // owner id
	})
	queryResultInModelb := []models.RealEstateModel{}
	json.Unmarshal(queryResult, &queryResultInModelb)

	expectb := []models.RealEstateModel{
		{
			RealEstateId: "3",
			OwnerId:      "2",
			Price:        "13000",
			Bed:          "1",
			Bath:         "1",
			AcreLot:      "150",
			FullAddress:  "cibinong",
			Street:       "mbongso",
			City:         "ndarjo",
			State:        "indo",
			ZipCode:      "61271",
			HouseSize:    "5",
			IsOpenToSell: "true",
		},
	}

	for i, v := range expectb {
		if v != queryResultInModelb[i] {
			t.Error("expectb &  queryResultInModelb values are not equal")
		}
	}
	//----------[user 2 should have a real estate]----------//

	//==========[user 3 should have a real estate]==========//
	queryResult = helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstate_GetByOwner"),
		[]byte("3"), // owner id
	})
	queryResultInModelb = []models.RealEstateModel{}
	json.Unmarshal(queryResult, &queryResultInModelb)

	expectb = []models.RealEstateModel{
		{
			RealEstateId: "4",
			OwnerId:      "3",
			Price:        "16000",
			Bed:          "12",
			Bath:         "11",
			AcreLot:      "1500",
			FullAddress:  "bangkalan",
			Street:       "meduro",
			City:         "madura",
			State:        "jerman",
			ZipCode:      "121414",
			HouseSize:    "53",
			IsOpenToSell: "false",
		},
	}

	for i, v := range expectb {
		if v != queryResultInModelb[i] {
			t.Error("expectb &  queryResultInModelb values are not equal")
		}
	}

	//----------[user 3 should have a real estate]----------//

	//==========[change real estate 3 ownership]==========//
	helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstate_ChangeRealEstateOwner"),
		[]byte("3"), // real estate id
		[]byte("3"), // new owner id
	})

	expect = []models.RealEstateHistoryModel{
		{
			OwnerID:      "2",
			RealEstateId: "3",
		},
		{
			OwnerID:      "3",
			RealEstateId: "3",
		},
	}

	queryResult = helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstateHistory_GetByRealEstateId"),
		[]byte("3"),
	})
	queryResultInModel = []models.RealEstateHistoryModel{}
	json.Unmarshal(queryResult, &queryResultInModel)

	if len(expect) != len(queryResultInModel) {
		t.Error("expect and queryResultInModel don't have the exac same array length")
	}

	for i, v := range expect {
		if v != queryResultInModel[i] {
			t.Error("expect and queryResultInModel don't have the same value")
		}
	}
	//----------[change real estate ownership]----------//

	//==========[user 3 should have 2 real estate]==========//
	queryResult = helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstate_GetByOwner"),
		[]byte("3"), // owner id
	})
	queryResultInModelb = []models.RealEstateModel{}
	json.Unmarshal(queryResult, &queryResultInModelb)

	expectb = []models.RealEstateModel{
		{
			RealEstateId: "3",
			OwnerId:      "3",
			Price:        "13000",
			Bed:          "1",
			Bath:         "1",
			AcreLot:      "150",
			FullAddress:  "cibinong",
			Street:       "mbongso",
			City:         "ndarjo",
			State:        "indo",
			ZipCode:      "61271",
			HouseSize:    "5",
			IsOpenToSell: "false",
		},
		{
			RealEstateId: "4",
			OwnerId:      "3",
			Price:        "16000",
			Bed:          "12",
			Bath:         "11",
			AcreLot:      "1500",
			FullAddress:  "bangkalan",
			Street:       "meduro",
			City:         "madura",
			State:        "jerman",
			ZipCode:      "121414",
			HouseSize:    "53",
			IsOpenToSell: "false",
		},
	}

	for i, v := range expectb {
		if v != queryResultInModelb[i] {
			t.Error("expectb &  queryResultInModelb values are not equal")
		}
	}
	//----------[user 3 should have 2 real estate]----------//

	//==========[user 2 should have no real estate]==========//
	queryResult = helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstate_GetByOwner"),
		[]byte("2"), // owner id
	})
	queryResultInModelb = []models.RealEstateModel{}
	json.Unmarshal(queryResult, &queryResultInModelb)

	expectb = []models.RealEstateModel{}

	if len(expectb) != len(queryResultInModelb) {
		t.Error("dont have the same length")
	}
	//----------[user 2 should have no real estate]----------//

	//==========[user 3 should have 2 real estates]==========//
	queryResult = helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstate_GetByOwner"),
		[]byte("3"), // owner id
	})
	queryResultInModelb = []models.RealEstateModel{}
	json.Unmarshal(queryResult, &queryResultInModelb)

	expectb = []models.RealEstateModel{
		{
			RealEstateId: "3",
			OwnerId:      "3",
			Price:        "13000",
			Bed:          "1",
			Bath:         "1",
			AcreLot:      "150",
			FullAddress:  "cibinong",
			Street:       "mbongso",
			City:         "ndarjo",
			State:        "indo",
			ZipCode:      "61271",
			HouseSize:    "5",
			IsOpenToSell: "false",
		},
		{
			RealEstateId: "4",
			OwnerId:      "3",
			Price:        "16000",
			Bed:          "12",
			Bath:         "11",
			AcreLot:      "1500",
			FullAddress:  "bangkalan",
			Street:       "meduro",
			City:         "madura",
			State:        "jerman",
			ZipCode:      "121414",
			HouseSize:    "53",
			IsOpenToSell: "false",
		},
	}

	if len(expectb) != len(queryResultInModelb) {
		t.Error("dont have the same length")
	}

	for i, v := range expectb {
		if v != queryResultInModelb[i] {
			t.Error("= expectb & queryResultInModelb dont't have exac same value")
		}
	}
	//----------[user 3 should have 2 real estates]----------//
}

func Test_TestRealEstate_GetByOwner(t *testing.T) {
	cc := new(cc.RealEstateChaincode)
	stub := shimtest.NewMockStub("real_estate", cc)

	//==========[example]==========//
	//----------[example]----------//

	//==========[init user]==========//
	helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("User_Init"),
	})
	//----------[init user]----------//

	//==========[init real estates]==========//
	helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstate_Init"),
	})
	//----------[init real estates]----------//

	//==========[user 2 should have 1 real estate]==========//
	expect := []models.RealEstateModel{
		{
			RealEstateId: "3",
			OwnerId:      "2",
			Price:        "13000",
			Bed:          "1",
			Bath:         "1",
			AcreLot:      "150",
			FullAddress:  "cibinong",
			Street:       "mbongso",
			City:         "ndarjo",
			State:        "indo",
			ZipCode:      "61271",
			HouseSize:    "5",
			IsOpenToSell: "true",
		},
	}

	queryResultAsBytes := helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstate_GetByOwner"),
		[]byte("2"), // real estate id
	})
	queryResult := []models.RealEstateModel{}
	json.Unmarshal(queryResultAsBytes, &queryResult)

	fmt.Println("expect: ", expect)
	fmt.Println("queryResult: ", queryResult)
	fmt.Println("len(expect): ", len(expect))
	fmt.Println("len(queryResult): ", len(queryResult))

	if len(expect) != len(queryResult) {
		t.Error("panjang array tidak sama")
	}

	for i, v := range expect {
		if v != queryResult[i] {
			t.FailNow()
			return
		}
	}
	//----------[user 2 should have 1 real estate]----------//
}

func Test_OwnerSetRealEstateToSell(t *testing.T) {
	cc := new(cc.RealEstateChaincode)
	stub := shimtest.NewMockStub("real_estate", cc)

	//==========[template]==========//
	//----------[template]----------//

	//==========[init real estates]==========//
	queryResultAsBytes := helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstate_Init"),
	})
	//----------[init real estates]----------//

	//==========[real estate with id 3 should have the true value of the IsOpenToSell field]==========//
	queryResultAsBytes = helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstate_GetById"),
		[]byte("0"), // real estate id
	})
	queryResult := models.RealEstateModel{}
	json.Unmarshal(queryResultAsBytes, &queryResult)

	expect := "false"

	if expect != queryResult.IsOpenToSell {
		t.Errorf("expect!=queryResult")
	}
	//----------[real estate with id 3 should have the true value of the IsOpenToSell field]----------//

	//==========[real estate record should not be there]==========//
	realEstateSaleRecordAsBytes := helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstateSalesRecord_GetByRealEstateIdComposite"),
		[]byte("0"), // real estate id
	})

	realEstateSalesRecord := []models.RealEstateSalesRecordModel{}
	json.Unmarshal(realEstateSaleRecordAsBytes, &realEstateSalesRecord)

	if len(realEstateSalesRecord) != 0 {
		t.Error("realEstateSalesRecord length is not 0")
	}

	fmt.Println("DADAD: " + string(realEstateSaleRecordAsBytes))
	//----------[real estate record should not be there]----------//

	//==========[change the real estate isopentosell value to be true]==========//
	queryResultAsBytes = helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstate_ChangeRealEstateSellStatus"),
		[]byte("0"),    // real estate id
		[]byte("true"), // status
	})
	//----------[change the real estate isopentosell value to be true]----------//

	//==========[real estate record should be there]==========//
	realEstateSaleRecordAsBytes = helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstateSalesRecord_GetByRealEstateIdComposite"),
		[]byte("0"), // real estate id
	})

	realEstateSalesRecord = []models.RealEstateSalesRecordModel{}
	json.Unmarshal(realEstateSaleRecordAsBytes, &realEstateSalesRecord)

	if len(realEstateSalesRecord) == 0 {
		t.Error("realEstateSalesRecord length shouldn't be 0")
	}
	fmt.Println("DADAD: " + string(realEstateSaleRecordAsBytes))
	//----------[real estate record should be there]----------//

	//==========[real estate with id 3 should have the true value of the IsOpenToSell field]==========//
	queryResultAsBytes = helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstate_GetById"),
		[]byte("0"), // real estate id
	})
	queryResult = models.RealEstateModel{}
	json.Unmarshal(queryResultAsBytes, &queryResult)

	expect = "true"

	if expect != queryResult.IsOpenToSell {
		t.Errorf("expect!=queryResult")
	}
	//----------[real estate with id 3 should have the true value of the IsOpenToSell field]----------//

	//==========[change the real estate isopentosell value to be false]==========//
	queryResultAsBytes = helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstate_ChangeRealEstateSellStatus"),
		[]byte("0"),     // real estate id
		[]byte("false"), // status
	})
	//----------[change the real estate isopentosell value to be false]----------//

	//==========[real estate record should not be there]==========//
	realEstateSaleRecordAsBytes = helper.Test_CheckInvoke(t, stub, [][]byte{
		[]byte("RealEstateSalesRecord_GetByRealEstateIdComposite"),
		[]byte("0"), // real estate id
	})

	realEstateSalesRecord = []models.RealEstateSalesRecordModel{}
	json.Unmarshal(realEstateSaleRecordAsBytes, &realEstateSalesRecord)

	if len(realEstateSalesRecord) != 0 {
		t.Error("realEstateSalesRecord length should be 0")
	}
	//----------[real estate record should not be there]----------//
}

func Test_ExternalAdvisorAssessTheRealEstate(t *testing.T) {
	// soon
}
