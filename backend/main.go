package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "time"

    _ "github.com/jackc/pgx/v5/stdlib" // magic psql driver
)

func main() {
    // Connect
    db, err := sql.Open("pgx", "postgres://kanban_user:kanban_password@db:5432/kanban_db")
    if err != nil {
        log.Fatalf("Unable to connect to the database: %v", err)
    }
    defer db.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        log.Fatalf("Database ping failed: %v", err)
    }
    log.Println("Connected to the database!")

    mux := http.NewServeMux()

	// Routes
    mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, "Go api is running")
    })

    // Start
    log.Println("Starting server on :9080")
    if err := http.ListenAndServe(":9080", mux); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
