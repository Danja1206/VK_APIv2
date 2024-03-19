package main

import (
	//"database/sql"
	"fmt"
	"os"
	"log"
	"net/http"
	"github.com/Danaj1412/vk/database"
	"github.com/Danaj1412/vk/api/handlers"
	"github.com/Danaj1412/vk/internal"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

)

func main() {
 
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbIP := os.Getenv("DB_IP")
	dbPort := os.Getenv("DB_PORT")
	db, err := database.InitDB(dbUser,dbPassword,dbIP,dbPort,dbName)
	if err != nil {
		log.Fatalf("Ошибка при инициализации базы данных: %v", err)
	}
	defer db.Close()


	fmt.Println("Сервер запущен")

	router := mux.NewRouter()
	router.HandleFunc("/user", internal.WithDBContextMiddleware(handlers.CreateUser, db)).Methods("POST")
	router.HandleFunc("/quest", internal.WithDBContextMiddleware(handlers.CreateQuest, db)).Methods("POST")
	router.HandleFunc("/big-quest", internal.WithDBContextMiddleware(handlers.CreateBigQuest, db)).Methods("POST")
	router.HandleFunc("/complete", internal.WithDBContextMiddleware(handlers.CompleteQuest, db)).Methods("POST")
	router.HandleFunc("/daily-reward", internal.WithDBContextMiddleware(handlers.DailyQuest, db)).Methods("POST")
	router.HandleFunc("/complete-big", internal.WithDBContextMiddleware(handlers.CompleteBigQuest, db)).Methods("POST")
	router.HandleFunc("/history/{user_id}", internal.WithDBContextMiddleware(handlers.GetUserHistory, db)).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
