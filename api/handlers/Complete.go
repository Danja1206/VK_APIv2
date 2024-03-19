package handlers

import (
	"net/http"
	"encoding/json"
	"fmt"

	"github.com/Danaj1412/vk/database"
	"github.com/Danaj1412/vk/internal"
)


func CompleteQuest(w http.ResponseWriter, r *http.Request) {
	var params struct {
		UserID  int `json:"user_id"`
		QuestID int `json:"quest_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&params)

	db, ok := database.FromContext(r.Context())
	if !ok {
		http.Error(w, "Не удалось получить подключение к базе данных из контекста", http.StatusInternalServerError)
		return
	}

	var isUnique bool
	err := db.QueryRow("SELECT is_unique FROM quests WHERE id = ?", params.QuestID).Scan(&isUnique)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch isUnique {
	case true:
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM user_quests WHERE quest_id = ?", params.QuestID).Scan(&count)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if count > 0 {
			http.Error(w, "Уникальный квест выполняется единожды!", http.StatusBadRequest)
			return
		}
		break
	case false:
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM user_quests WHERE quest_id = ? AND user_id = ?", params.QuestID, params.UserID).Scan(&count)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var maxLimit int
		err = db.QueryRow("SELECT max_limit FROM quests WHERE id = ?", params.QuestID).Scan(&maxLimit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if count >= maxLimit {
			http.Error(w, "Вы достигли лимита выполнений", http.StatusBadRequest)
			return
		}

		break
		
	}

	// Начисляем награду и помечаем квест как выполненный
	_, err = db.Exec("UPDATE users SET balance = balance + (SELECT cost FROM quests WHERE id = ?) WHERE id = ?", params.QuestID, params.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var exp int
	err = db.QueryRow("SELECT exp FROM quests WHERE id = ? ", params.QuestID).Scan(&exp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	internal.GetExpToUser(params.UserID, exp, r,w);
	

	_, err = db.Exec("INSERT INTO user_quests(user_id, quest_id, completed_at) VALUES(?, ?, NOW())", params.UserID, params.QuestID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Квест выполнен!")
}