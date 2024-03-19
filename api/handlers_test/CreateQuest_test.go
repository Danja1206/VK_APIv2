package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Danaj1412/vk/api/models"
	"github.com/Danaj1412/vk/internal"
	"github.com/Danaj1412/vk/database"
	"github.com/Danaj1412/vk/api/handlers"
)

func TestCreateQuest(t *testing.T) {

	db, err := database.InitDB("root", "root", "127.0.0.1", "3306", "vk")
	if err != nil {
		log.Fatalf("Ошибка при инициализации базы данных: %v", err)
	}
	defer db.Close()

	quest := models.Quest{
		Name:     "Test Quest",
		Cost:     100,
		IsUnique: true,
		MaxLimit: 1,
		Exp:      50,
	}
	jsonData, err := json.Marshal(quest)
	assert.NoError(t, err, "Ошибка при маршалинге JSON")

	req, err := http.NewRequest("POST", "/quest", bytes.NewBuffer(jsonData))
	assert.NoError(t, err, "Ошибка при создании запроса")

	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), "db", db)
	req = req.WithContext(ctx)

	handler := http.HandlerFunc(internal.WithDBContextMiddleware(handlers.CreateQuest, db))
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Ожидался код ответа 200")

	var createdQuest models.Quest
	err = json.Unmarshal(rr.Body.Bytes(), &createdQuest)
	assert.NoError(t, err, "Ошибка при анмаршалинге JSON")

}