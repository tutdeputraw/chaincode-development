package models

// type RealEstateModel struct {
// 	Id          string  `json:"id"`
// 	Price       string  `json:"price"` //0
// 	Bed         int     `json:"bed"`
// 	Bath        int     `json:"bath"`
// 	AcreLot     int     `json:"acre_lot"`
// 	FullAddress string  `json:"full_address"`
// 	Street      string  `json:"street"`
// 	City        string  `json:"city"`
// 	State       string  `json:"state"`
// 	ZipCode     int     `json:"zip_code"`
// 	HouseSize   float64 `json:"house_size"`
// }

type RealEstateModel struct {
	RealEstateId string `json:"realEstateModel_id"`
	OwnerId      string `json:"realEstateModel_owner_id"`
	Price        string `json:"realEstateModel_price"` //0
	Bed          string `json:"realEstateModel_bed"`
	Bath         string `json:"realEstateModel_bath"`
	AcreLot      string `json:"realEstateModel_acre_lot"`
	FullAddress  string `json:"realEstateModel_full_address"`
	Street       string `json:"realEstateModel_street"`
	City         string `json:"realEstateModel_city"`
	State        string `json:"realEstateModel_state"`
	ZipCode      string `json:"realEstateModel_zip_code"`
	HouseSize    string `json:"realEstateModel_house_size"`
	IsOpenToSell string `json:"realEstateModel_is_open_to_sell"`
}
