package api

import (
	"encoding/json"
	"my-go-app/db"
	"net/http"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

// ИЗМЕНИТЬ ЭТО: было Handler, стало UsersHandler
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
        rows, err := db.GetDB().Query("SELECT id, name FROM users")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var users []User
        for rows.Next() {
            var user User
            if err := rows.Scan(&user.ID, &user.Name); err != nil {
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

        err := db.GetDB().QueryRow(
            "INSERT INTO users(name) VALUES($1) RETURNING id",
            user.Name,
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