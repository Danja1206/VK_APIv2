package handlers

import (
	"net/http"
	"encoding/json"

	"github.com/Danaj1412/vk/database"
	"github.com/Danaj1412/vk/api/models"
)

func CreateBigQuest(w http.ResponseWriter, r *http.Request) {


	type BigQuestInput struct {
		models.BigQuest
		Steps []models.BigQuestStep `json:"steps"`
	}
	var params BigQuestInput
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	
	db, ok := database.FromContext(r.Context())
	if !ok {
		http.Error(w, "Не удалось получить подключение к базе данных из контекста", http.StatusInternalServerError)
		return
	}

	res, err := db.ExecContext(r.Context(),"INSERT INTO big_quests(name, cost, exp ) VALUES(?, ?, ?)", params.Name, params.Cost, params.Exp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bigQuestID, err := res.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, step := range params.Steps {
		_, err = db.ExecContext(r.Context(),"INSERT INTO big_quest_steps (big_quest_id, quest_id, step_order) VALUES (?, ?, ?)", bigQuestID, step.QuestID, step.Order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(params.BigQuest)
}