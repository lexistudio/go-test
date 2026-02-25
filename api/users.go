package api

import (
	"encoding/json"
	"my-go-app/db"
	"net/http"
)

// Структура с обычными типами (без sql.Null*)
type User struct {
    ID        int    `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Age       int    `json:"age"`
    Avatar    string `json:"avatar"`
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }

    db.InitDB()
    defer db.GetDB().Close()

    switch r.Method {
    case http.MethodGet:
        rows, err := db.GetDB().Query("SELECT id, first_name, last_name, age, avatar FROM users ORDER BY id")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var users []User
        for rows.Next() {
            var user User
            // Сканируем напрямую в string/int (если в БД есть NULL, будет ошибка)
            if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Age, &user.Avatar); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            users = append(users, user)
        }

        json.NewEncoder(w).Encode(users)

    case http.MethodPost:
        var user User
        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // Валидация
        if user.FirstName == "" || user.LastName == "" {
            http.Error(w, "First name and last name are required", http.StatusBadRequest)
            return
        }
        if user.Age < 1 || user.Age > 120 {
            http.Error(w, "Age must be between 1 and 120", http.StatusBadRequest)
            return
        }

        err := db.GetDB().QueryRow(
            "INSERT INTO users (first_name, last_name, age, avatar) VALUES ($1, $2, $3, $4) RETURNING id",
            user.FirstName, user.LastName, user.Age, user.Avatar,
        ).Scan(&user.ID)

        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(user)

    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}