package models


type User struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Balance int `json:"balance"`
	Exp int `json:"exp"`
	Lvl int `json:"lvl"`
}