package models


type DailyRewards struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	QuestID int `json:"quest_id"`
	CompletedAt string `json:"completed_at"`
}