package models

type Quest struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Cost int `json:"cost"`
	IsUnique bool `json:"is_unique"`
	MaxLimit int `json:"max_limit"`
	Exp int `json:"exp"`
}