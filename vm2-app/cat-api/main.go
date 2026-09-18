package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "os"

    _ "github.com/jackc/pgx/v5/stdlib"
)

type Cat struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
}

func main() {
    dsn := os.Getenv("DATABASE_URL")
    addr := os.Getenv("API_ADDR")

    if dsn == "" || addr == "" {
        log.Fatal("DATABASE_URL and API_ADDR must be set")
    }

    db, err := sql.Open("pgx", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok\n"))
    })

    http.HandleFunc("/cat", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var cat Cat
        err := db.QueryRow(`
            SELECT id, name, description
            FROM cats
            ORDER BY random()
            LIMIT 1
        `).Scan(&cat.ID, &cat.Name, &cat.Description)

        if err != nil {
            http.Error(w, "database error", http.StatusInternalServerError)
            log.Println(err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(cat)
    })

    log.Printf("API listens on %s", addr)
    log.Fatal(http.ListenAndServe(addr, nil))
}
