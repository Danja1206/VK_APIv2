package models


type DailyQuest struct {
    ID int `json:"id"`
    Name string `json:"name"`
    Description string `json:"description"`
    Cost int `json:"reward_points"`
}