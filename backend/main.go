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
	"strconv"
	"github.com/gorilla/mux"
    _ "github.com/jackc/pgx/v5/stdlib" // magic psql driver
)


// represents a task on the Kanban board
type Task struct {
    ID     *int    `json:"id"`
    Title  string `json:"title"`
    Status string `json:"status"`
	Description string `json:"description"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
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
	r := mux.NewRouter() // gorilla Mux Router

	// ROUTES
	r.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		w.WriteHeader(http.StatusNoContent)
	})

    r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, "Go api is running!")
    })

	  // Fetch all tasks
	r.HandleFunc("/api/tasks/list", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel() 
	
		rows, err := db.QueryContext(ctx, "SELECT id, title, status FROM tasks")
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Println("Request timeout exceeded.")
				http.Error(w, "Request timed out", http.StatusRequestTimeout)
				return
			}
			fmt.Println(err)
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
	}).Methods("GET", "OPTIONS")

    // Create / Update task
    r.HandleFunc("/api/tasks/create", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		// createTaskHandler(§db, &w, &r)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel() 

		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if task.ID != nil {
			// Updating task
			fmt.Println("Updating task:", *task.ID)
			query := `UPDATE tasks SET title = $1, description = $2, status = $3 WHERE id = $4`
			_, err := db.ExecContext(ctx, query, task.Title, task.Description, task.Status, *task.ID)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					fmt.Println("task update timeout:", err)
					http.Error(w, "Request timed out", http.StatusRequestTimeout)
					return
				}
				fmt.Println("task update error:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			// Creating a new task
			fmt.Println("Creating task:", *task.ID)
			query := `INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3)`
			_, err := db.ExecContext(ctx, query, task.Title, task.Description, task.Status)
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
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)

	}).Methods("POST")

	r.HandleFunc("/api/tasks/{id}", func(w http.ResponseWriter, r *http.Request){
		enableCors(&w)

		vars := mux.Vars(r)
		id := vars["id"]

		fmt.Println("Deleting task:", id)
		_ ,err := strconv.Atoi(id)
		if err != nil {
			log.Printf("Failed to delete task: %v\n", err)
			http.Error(w, "Failed to delete task", http.StatusInternalServerError)
			return
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := db.ExecContext(ctx, "DELETE FROM tasks WHERE id = $1", id)
		if err != nil {
			log.Printf("Failed to delete task: %v\n", err)
			http.Error(w, "Failed to delete task", http.StatusInternalServerError)
			return
		}
		
		rowsAffected, _ := result.RowsAffected() // Check if rows were affected ?
		if rowsAffected == 0 {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
	
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Task deleted successfully"}`))
	}).Methods("DELETE")

    // Start
    log.Println("Starting server on :9080")
    if err := http.ListenAndServe(":9080", r); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
