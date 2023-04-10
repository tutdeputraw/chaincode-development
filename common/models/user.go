package models

type UserModel struct {
	Id          string `json:"userModel_id"`
	Name        string `json:"userModel_name"`
	NPWPNumber  string `json:"userModel_npwp_number"`
	PhoneNumber string `json:"userModel_phone_number"`
	Email       string `json:"userModel_email"`
}
