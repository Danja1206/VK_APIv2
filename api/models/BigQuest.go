package models

type BigQuest struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Cost int `json:"cost"`
	Exp int `json:"exp"`
}

type BigQuestStep struct {
	QuestID int `json:"quest_id"`
	Order   int `json:"order"`
}