package api

import (
	"encoding/json"
	"net/http"
)

type Message struct {
    Message string            `json:"message"`
    Routes  map[string]string `json:"routes"`
}

// ИЗМЕНИТЬ ЭТО: было Handler, стало IndexHandler
func IndexHandler(w http.ResponseWriter, r *http.Request) {
    msg := Message{
        Message: "Go API на Vercel работает!",
        Routes: map[string]string{
            "/api":         "Главная",
            "/api/users":   "Получить список пользователей",
        },
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    json.NewEncoder(w).Encode(msg)
}