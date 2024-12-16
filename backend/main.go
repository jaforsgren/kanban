package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "time"
	"encoding/json"
	"errors"

    _ "github.com/jackc/pgx/v5/stdlib" // magic psql driver
)


// represents a task on the Kanban board
type Task struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Status string `json:"status"`
}

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

	// ROUTES

    mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, "Go api is running!")
    })

	  // Fetch all tasks
	mux.HandleFunc("/api/tasks/list", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Listing tasks...")
	
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()  // Ensure context is canceled after the function completes
	
		rows, err := db.QueryContext(ctx, "SELECT id, title, status FROM tasks")
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Println("Request timeout exceeded.")
				http.Error(w, "Request timed out", http.StatusRequestTimeout)
				return
			}
			fmt.Println("read error...", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
	
		tasks := []Task{}
		for rows.Next() {
			var task Task
			if err := rows.Scan(&task.ID, &task.Title, &task.Status); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tasks = append(tasks, task)
		}
	
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	})

    // Create a new task
    mux.HandleFunc("/api/tasks/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()  // Ensure context is canceled after the function completes
	
			var task Task
			if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
	
			fmt.Println("Creating task:", task)
	
			_, err := db.ExecContext(ctx, "INSERT INTO tasks (title, status) VALUES ($1, $2)", task.Title, task.Status)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					fmt.Println("task creation timeout:", err)
					http.Error(w, "Request timed out", http.StatusRequestTimeout)
					return
				}
				fmt.Println("task creation error:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
	
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(task)
		}
	})
	

    // Start
    log.Println("Starting server on :9080")
    if err := http.ListenAndServe(":9080", mux); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
