package handlers

import (
	"net/http"
	"encoding/json"

	"github.com/Danaj1412/vk/database"
	"github.com/Danaj1412/vk/api/models"
)


func CreateQuest(w http.ResponseWriter, r *http.Request) {
	var quest models.Quest

	_ = json.NewDecoder(r.Body).Decode(&quest)

	db, ok := database.FromContext(r.Context())
	if !ok {
		http.Error(w, "Не удалось получить подключение к базе данных из контекста", http.StatusInternalServerError)
		return
	}

	_, err := db.ExecContext(r.Context(),"INSERT INTO quests(name, cost, is_unique, max_limit,exp) VALUES(?, ?, ?, ?, ?)", quest.Name, quest.Cost, quest.IsUnique, quest.MaxLimit, quest.Exp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(quest)
}