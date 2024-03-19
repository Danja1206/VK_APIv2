package models

type CompletedQuest struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	QuestID int `json:"quest_id"`
	CompletedAt string `json:"completed_at"`
}