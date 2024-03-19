package handlers

import (
	"net/http"
	"encoding/json"
	"database/sql"
	"fmt"
	"time"

	"github.com/Danaj1412/vk/database"
	"github.com/Danaj1412/vk/api/models"
	"github.com/Danaj1412/vk/internal"
)

func DailyQuest(w http.ResponseWriter, r *http.Request) {

	var quest models.DailyRewards


	_ = json.NewDecoder(r.Body).Decode(&quest)

	db, ok := database.FromContext(r.Context())
	if !ok {
		http.Error(w, "Не удалось получить подключение к базе данных из контекста", http.StatusInternalServerError)
		return
	}

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM user_daily_quests WHERE user_id = ?)", quest.UserID).Scan(&exists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		result, err := db.ExecContext(r.Context(),"INSERT INTO daily_quests (name, description, cost, exp) VALUES (?, ?, ?, ?)", fmt.Sprintf("user_%v_daily", quest.UserID), "Get you daily reward!", 100,125)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {

				dailyQuestID, err := result.LastInsertId()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				_, err = db.ExecContext(r.Context(),"INSERT INTO user_daily_quests (user_id, daily_quest_id, completed_at) VALUES (?, ?, NOW())", quest.UserID, dailyQuestID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				} else {
					_, err = db.Exec("UPDATE users SET balance = (balance + (SELECT cost FROM daily_quests WHERE id = ?) + lvl * 0.1) WHERE id = ?", dailyQuestID, quest.UserID)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					internal.GetExpToUser(quest.UserID, 125, r,w);
				}
				fmt.Fprintf(w, "Вы получили ежедневную награду!")
			}
		return
	} else {
		var lastCompletedAt sql.NullTime
		err := db.QueryRow(`
			SELECT completed_at
			FROM user_daily_quests
			WHERE user_id = ?
			ORDER BY completed_at DESC
			LIMIT 1
		`, quest.UserID).Scan(&lastCompletedAt)
		if err != nil {
			if err == sql.ErrNoRows {
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		
		if !lastCompletedAt.Valid || time.Since(lastCompletedAt.Time) > 24*time.Hour {
			var dailyQuestID int
			err := db.QueryRow("SELECT daily_quests_id FROM user_daily_quests WHERE user_id = ?", quest.UserID).Scan(&dailyQuestID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			_, err = db.ExecContext(r.Context(),"UPDATE user_daily_quests SET completed_at = NOW() WHERE user_id = ?", quest.UserID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			_, err = db.ExecContext(r.Context(),"UPDATE users SET balance = (balance + (SELECT cost FROM daily_quests WHERE id = ?) + lvl * 0.1) WHERE id = ?", dailyQuestID, quest.UserID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			internal.GetExpToUser(quest.UserID, 125, r,w);
			
			internal.GetExpToUser(quest.UserID, 10,r,w)
			fmt.Fprintf(w, "Вы получили ежедневную награду!")
			return
		}
		
		remainingTime := 24*time.Hour - time.Since(lastCompletedAt.Time)
		
		remainingTimeFormatted := fmt.Sprintf("%02d:%02d:%02d",
			int(remainingTime.Hours()),
			int(remainingTime.Minutes())%60,
			int(remainingTime.Seconds())%60,
		)
		
		fmt.Fprintf(w, "Оставшееся время до получения награды квеста: %s", remainingTimeFormatted)
	}

}