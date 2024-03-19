package internal

import (
	"net/http"
	"github.com/Danaj1412/vk/database"
	"math"

	_ "github.com/go-sql-driver/mysql"

)

func GetUserLvl(user_id int, r *http.Request, w http.ResponseWriter) int {

	db, ok := database.FromContext(r.Context())
	if !ok {
		http.Error(w, "Не удалось получить подключение к базе данных из контекста", http.StatusInternalServerError)
		return -1
	}

	var result int
	err := db.QueryRow("SELECT lvl FROM users WHERE id = ?", user_id).Scan(&result)
	if err != nil {
		return -1
	}

	return result
}

func GetUserExp(user_id int, r *http.Request, w http.ResponseWriter) int {

	db, ok := database.FromContext(r.Context())
	if !ok {
		http.Error(w, "Не удалось получить подключение к базе данных из контекста", http.StatusInternalServerError)
		return -1
	}

	var result int
	err := db.QueryRow("SELECT exp FROM users WHERE id = ?", user_id).Scan(&result)
	if err != nil {
		return -1
	}

	return result
}

func GetExpToNextLvl(user_id int, r *http.Request, w http.ResponseWriter) int {

	db, ok := database.FromContext(r.Context())
	if !ok {
		http.Error(w, "Не удалось получить подключение к базе данных из контекста", http.StatusInternalServerError)
		return -1
	}

	var result int
	err := db.QueryRow("SELECT lvl.exp_to_lvl FROM levels lvl JOIN users us ON us.lvl = lvl.lvl WHERE us.id = ?", user_id).Scan(&result)
	if err != nil {
		return -1
	}
	return result;
}

func GetExpToLvlUp(user_exp int , user_lvl int, r *http.Request, w http.ResponseWriter) (int) {
	db, ok := database.FromContext(r.Context())
	if !ok {
		http.Error(w, "Не удалось получить подключение к базе данных из контекста", http.StatusInternalServerError)
		return -1
	}
    var data int
    err := db.QueryRow("SELECT (? - exp_total) AS data FROM levels WHERE lvl = ?", user_exp, user_lvl).Scan(&data)
    if err != nil {
        return 0
    }
    return data

}

func GetLvlPercent(user_id int, r *http.Request, w http.ResponseWriter) int {

	percent := float64(GetExpToLvlUp(GetUserExp(user_id,r,w), GetUserLvl(user_id,r,w),r,w)) / float64(GetExpToNextLvl(user_id,r,w)) * 100;
	return int(math.Round(percent));

}

func UpdateUserLvl(user_id int, r *http.Request, w http.ResponseWriter) {

	db, ok := database.FromContext(r.Context())
	if !ok {
		http.Error(w, "Не удалось получить подключение к базе данных из контекста", http.StatusInternalServerError)
		return
	}

    var data int
    err := db.QueryRow("SELECT MAX(lvl) AS max_lvl FROM levels WHERE exp_total <= ?", GetUserExp(user_id,r,w)).Scan(&data)
    if err != nil {
        return
    }

	_, err = db.ExecContext(r.Context(),"UPDATE users SET lvl = ? WHERE id = ?", data, user_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func GetExpToUser(user_id int, exp_count int , r *http.Request, w http.ResponseWriter) {
	db, ok := database.FromContext(r.Context())
	if !ok {
		http.Error(w, "Не удалось получить подключение к базе данных из контекста", http.StatusInternalServerError)
		return
	}

	_, err := db.ExecContext(r.Context(),"UPDATE users SET exp = exp + ?  WHERE id = ?", exp_count, user_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if GetLvlPercent(user_id, r, w) >= 100 {
		UpdateUserLvl(user_id, r, w);
	}
}