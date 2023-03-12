package models

type UserModel struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	NPWPNumber  string `json:"npwp_number"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}
