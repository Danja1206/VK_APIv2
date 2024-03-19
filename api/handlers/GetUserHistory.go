package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"

	"github.com/Danaj1412/vk/database"
	"github.com/Danaj1412/vk/api/models"
)


func GetUserHistory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]

	db, ok := database.FromContext(r.Context())
	if !ok {
		http.Error(w, "Не удалось получить подключение к базе данных из контекста", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT q.id, q.name, q.cost, uq.completed_at FROM quests q JOIN user_quests uq ON q.id = uq.quest_id WHERE uq.user_id = ?", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	rowsBig, errBig := db.Query("SELECT bq.id, bq.name, bq.cost, ubq.completed_at FROM big_quests bq JOIN user_big_quests ubq ON ubq.user_id = bq.id WHERE ubq.user_id = ?", userID)
	if errBig != nil {
		http.Error(w, errBig.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type QuestHistory struct {
		models.Quest
		CompletedAt string `json:"completed_at"`
	}
	type BigQuestHistory struct {
		models.BigQuest
		CompletedAt string `json:"completed_at"`
	}

	var questsHistory []QuestHistory
	for rows.Next() {
		var quest QuestHistory
		err := rows.Scan(&quest.Quest.ID, &quest.Quest.Name, &quest.Quest.Cost, &quest.CompletedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		questsHistory = append(questsHistory, quest)
	}


	var bigQuestsHistory []BigQuestHistory
	for rowsBig.Next() {
		var bigQuest BigQuestHistory
		err := rowsBig.Scan(&bigQuest.BigQuest.ID, &bigQuest.BigQuest.Name, &bigQuest.BigQuest.Cost, &bigQuest.CompletedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bigQuestsHistory = append(bigQuestsHistory, bigQuest)
	}

	var user models.User
	err = db.QueryRow("SELECT id, name, balance, exp, lvl FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Name, &user.Balance, &user.Exp, &user.Lvl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		User models.User `json:"user"`
		Quests []QuestHistory `json:"quests"`
		BigQuests []BigQuestHistory `json:"big_quests"`
	}{
		User:   user,
		Quests: questsHistory,
		BigQuests: bigQuestsHistory,
	}

	json.NewEncoder(w).Encode(response)
}