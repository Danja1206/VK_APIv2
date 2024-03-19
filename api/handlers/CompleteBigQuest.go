package handlers

import (
	"net/http"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Danaj1412/vk/database"
	"github.com/Danaj1412/vk/internal"
)

func CompleteBigQuest(w http.ResponseWriter, r *http.Request) {
	var params struct {
		UserID  int `json:"user_id"`
		BigQuestID  int `json:"big_quest_id"`
	}

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

	var currentStep int
	err = db.QueryRow("SELECT current_step FROM user_big_quests_steps WHERE user_id = ? AND big_quest_id = ?", params.UserID, params.BigQuestID).Scan(&currentStep)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = db.ExecContext(r.Context(),"INSERT INTO user_big_quests_steps (user_id, big_quest_id, current_step) VALUES (?, ?, 1)", params.UserID, params.BigQuestID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "Первый шаг выполнен!")
			return
		}
	}

	var nextStep int
	var maxStep int

	err = db.QueryRow("SELECT MAX(step_order) FROM big_quest_steps WHERE big_quest_id = ? AND step_order > ?", params.BigQuestID, currentStep).Scan(&maxStep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.QueryRow("SELECT MIN(step_order) FROM big_quest_steps WHERE big_quest_id = ? AND step_order > ?", params.BigQuestID, currentStep).Scan(&nextStep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
		_, err = db.ExecContext(r.Context(),"UPDATE user_big_quests_steps SET current_step = current_step + 1 WHERE user_id = ? AND big_quest_id = ?", params.UserID, params.BigQuestID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if maxStep == nextStep {
			_, err = db.ExecContext(r.Context(),"INSERT INTO user_big_quests(user_id, big_quest_id, completed_at) VALUES(?, ?, NOW())", params.UserID, params.BigQuestID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = db.Exec("UPDATE users SET balance = balance + (SELECT cost FROM big_quests WHERE id = ?) WHERE id = ?", params.BigQuestID, params.UserID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var exp int
			err = db.QueryRow("SELECT exp FROM big_quests WHERE id = ? ", params.BigQuestID).Scan(&exp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		
			internal.GetExpToUser(params.UserID, exp, r,w);
			fmt.Fprintf(w, "Квест выполнен")
		}else {
			fmt.Fprintf(w, "Шаг квеста обновлен до %d", nextStep)
		}



}