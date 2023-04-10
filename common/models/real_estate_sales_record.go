package models

type RealEstateSalesRecordModel struct {
	RealEstateId      string `json:"RealEstateSalesRecordModel_id"`
	SellerId          string `json:"RealEstateSalesRecordModel_seller_id"`
	PropertyAdvisorId string `json:"RealEstateSalesRecordModel_property_advisor_id"`
	PropertyAgentId   string `json:"RealEstateSalesRecordModel_property_agent_id"`

	// this assessment is done by property advisor
	// only real estate advisor can change this state
	RealEstateAssessment string `json:"RealEstateSalesRecordModel_real_estate_assessment"`

	// this data contains integer value and always increments everytime users clicking the detail page
	InterestUsers string `json:"RealEstateSalesRecordModel_interest_users"`

	// there are 4 phase in the real estate transaction [preparation, due diligence, completion]
	// only broker can change this state
	Phase string `json:"RealEstateSalesRecordModel_phase"`
}
