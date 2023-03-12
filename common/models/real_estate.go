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
	RealEstateId string `json:"id"`
	OwnerId      string `json:"owner_id"`
	Price        string `json:"price"` //0
	Bed          string `json:"bed"`
	Bath         string `json:"bath"`
	AcreLot      string `json:"acre_lot"`
	FullAddress  string `json:"full_address"`
	Street       string `json:"street"`
	City         string `json:"city"`
	State        string `json:"state"`
	ZipCode      string `json:"zip_code"`
	HouseSize    string `json:"house_size"`
	IsOpenToSell string `json:"is_open_to_sell"`
}
