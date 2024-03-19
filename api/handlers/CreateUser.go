package handlers

import (
	"net/http"
	"encoding/json"
	"fmt"

	"github.com/Danaj1412/vk/database"
	"github.com/Danaj1412/vk/api/models"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)

	db, ok := database.FromContext(r.Context())
	if !ok {
		http.Error(w, "Не удалось получить подключение к базе данных из контекста", http.StatusInternalServerError)
		return
	}

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE name = ?)", user.Name).Scan(&exists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		_, err := db.ExecContext(r.Context(), "INSERT INTO users(name, balance) VALUES(?, ?)", user.Name, user.Balance)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(user)
	} else {
		
	fmt.Fprintf(w, "Пользователь существует!")
	}
}