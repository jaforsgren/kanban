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
	"github.com/gorilla/mux"
	"github.com/google/uuid"
    _ "github.com/jackc/pgx/v5/stdlib" // magic psql driver
)


// represents a task on the Kanban board
type Task struct {
    ID        uuid.UUID `json:"id"`
    BoardID   uuid.UUID `json:"board_id"`
    Title     string    `json:"title"`
    Status    string    `json:"status"`
    Description string `json:"description"`
}

type User struct {
    ID       uuid.UUID `json:"id"`
    Username string    `json:"username"`
    Email    string    `json:"email"`
}

type Board struct {
    ID      uuid.UUID `json:"id"`
    UserID  uuid.UUID `json:"user_id"`
    Title    string    `json:"title"`
}



func getAllBoardsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel() 
		
		// TODO: select board id , and 
		// TODO: get board details as well
		rows, err := db.QueryContext(ctx, "SELECT id, title FROM boards")
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
	
		boards := []Board{}
		for rows.Next() {
			var board Board
			if err := rows.Scan(&board.ID, &board.Title); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			boards = append(boards, board)
		}
	
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(boards)
	}
}

// Get specific board
func getBoardDetailsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		boardID := mux.Vars(r)["id"]

		enableCors(&w)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel() 
		
		// TODO: select board id , and 
		// TODO: get board details as well
		rows, err := db.QueryContext(ctx, "SELECT id, title, board_id, description, status FROM tasks WHERE board_id = $1", boardID)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Println("Request timeout exceeded.")
				http.Error(w, "Request timed out", http.StatusRequestTimeout)
				return
			}
			
			fmt.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		defer rows.Close()
		
		tasks := []Task{}
		for rows.Next() {
			var task Task
			if err := rows.Scan(&task.ID, &task.Title, &task.BoardID, &task.Description, &task.Status); err != nil {
				fmt.Println("Error:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tasks = append(tasks, task)
		}
	
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	}
}

func createUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
        
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user.ID = uuid.New()
		fmt.Println("Creating user:", user.ID)
		query := `INSERT INTO users (id, username, email) VALUES ($1, $2, $3)`
		_, err := db.ExecContext(ctx, query, user.ID, user.Username, user.Email)

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Println("user creation timeout:", err)
				http.Error(w, "Request timed out", http.StatusRequestTimeout)
				return
			}
			fmt.Println("user creation error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

func createBoardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
        
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

		var board Board
		if err := json.NewDecoder(r.Body).Decode(&board); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}	

		 // Creating a new board
		 board.ID = uuid.New()
		 fmt.Println("Creating board:", board.ID)
		 fmt.Println("Creating board:", board)
		 
		 query := `INSERT INTO boards (id, user_id, title) VALUES ($1, $2, $3)`
		 _, err := db.ExecContext(ctx, query, board.ID, board.UserID, board.Title, )
		 if err != nil {
			 if errors.Is(err, context.DeadlineExceeded) {
				 fmt.Println("board creation timeout:", err)
				 http.Error(w, "Request timed out", http.StatusRequestTimeout)
				 return
			 }
			 fmt.Println("board creation error:", err)
			 http.Error(w, err.Error(), http.StatusInternalServerError)
			 return
		 }

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(board)
	}
}


func deleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		vars := mux.Vars(r)
		id := vars["id"]

		fmt.Println("Deleting task:", id)

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
	}
}

// TODO: add this!
func deleteBoardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		vars := mux.Vars(r)
		id := vars["id"]

		fmt.Println("Deleting board:", id)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := db.ExecContext(ctx, "DELETE FROM boards WHERE id = $1", id)
		if err != nil {
			log.Printf("Failed to delete board: %v\n", err)
			http.Error(w, "Failed to delete board", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected() // Check if rows were affected ?
		if rowsAffected == 0 {
			http.Error(w, "Board not found", http.StatusNotFound)
			return
		}
	
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Task deleted successfully"}`))
	}
}

func createTaskHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        enableCors(&w)

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        var task Task
        if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

		task.ID = uuid.New()
		fmt.Println("Creating task:", task.ID , "for board" , task.BoardID )
		query := `INSERT INTO tasks (id, board_id, title, description, status) VALUES ($1, $2, $3, $4, $5)`
		_, err := db.ExecContext(ctx, query, task.ID, task.BoardID, task.Title, task.Description, task.Status)
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
}

func updateTaskHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        enableCors(&w)

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        var task Task
        if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

		if task.ID == uuid.Nil {
			http.Error(w, "Invalid task id", http.StatusBadRequest)
			log.Println("Invalid task id:", task.ID)
			return
		}
        
		fmt.Println("Updating task:", task.ID  )
		query := `UPDATE tasks SET title = $1, description = $2, status = $3 WHERE id = $4`
		_, err := db.ExecContext(ctx, query, task.Title, task.Description, task.Status, task.ID)
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
    
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(task)
    }
}

func updateBoardHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        enableCors(&w)

		boardID := mux.Vars(r)["board_id"]

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        var board Board
        if err := json.NewDecoder(r.Body).Decode(&board); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

		fmt.Println("Updating board:", board.ID )
		query := `UPDATE boards SET title = $1 WHERE id = $2`
		_, err := db.ExecContext(ctx, query, board.Title, boardID)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Println("board update timeout:", err)
				http.Error(w, "Request timed out", http.StatusRequestTimeout)
				return
			}
			fmt.Println("board update error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
    
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(board)
    }
}

func healtHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Go api is running!")
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, DELETE, PATCH")
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

    r.HandleFunc("/api/health", healtHandler)

	r.HandleFunc("/api/users", createUserHandler(db)).Methods("POST")

	r.HandleFunc("/api/boards", createBoardHandler(db)).Methods("POST")

	r.HandleFunc("/api/boards", getAllBoardsHandler(db)).Methods("GET")

	r.HandleFunc("/api/boards/{id}", getBoardDetailsHandler(db)).Methods("GET")

	r.HandleFunc("/api/boards/{id}", deleteBoardHandler(db)).Methods("DELETE")

	r.HandleFunc("/api/boards/{id}", updateBoardHandler(db)).Methods("PATCH")

	r.HandleFunc("/api/tasks", createTaskHandler(db)).Methods("POST")

    r.HandleFunc("/api/tasks", updateTaskHandler(db)).Methods("PATCH")

	r.HandleFunc("/api/tasks/{id}", deleteTaskHandler(db)).Methods("DELETE")

    // Start
    log.Println("Starting server on :9080")
    if err := http.ListenAndServe(":9080", r); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
