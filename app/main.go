package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

var db *sql.DB

func main() {
	// DSN: host=postgres — это имя сервиса в docker-compose
	dsn := fmt.Sprintf("postgres://%s:%s@postgres:5432/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	var err error
	db, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}

	// Retry: ждём, пока Postgres поднимется
	for i := 1; i <= 10; i++ {
		if err = db.Ping(); err == nil {
			log.Println("db connected")
			break
		}
		log.Printf("db ping attempt %d/10 failed: %v", i, err)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Fatalf("db not reachable: %v", err)
	}

	// Создаём таблицу, если её нет
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatalf("create table: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /metrics", metricsHandler)
	mux.HandleFunc("GET /tasks", getTasksHandler)
	mux.HandleFunc("POST /tasks", createTasksHandler)

	log.Println("server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("listen: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("# placeholder\n"))
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.QueryContext(r.Context(), `SELECT id, title FROM tasks ORDER BY id`)
	if err != nil {
		http.Error(w, `{"error":"db query failed"}`, 500)
		return
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title); err != nil {
			http.Error(w, `{"error":"scan failed"}`, 500)
			return
		}
		tasks = append(tasks, t)
	}
	json.NewEncoder(w).Encode(tasks)
}

func createTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	var input struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(400)
		w.Write([]byte(`{"error":"invalid JSON"}`))
		return
	}
	if strings.TrimSpace(input.Title) == "" {
		w.WriteHeader(400)
		w.Write([]byte(`{"error":"title required"}`))
		return
	}

	var newID int
	err := db.QueryRowContext(r.Context(),
		`INSERT INTO tasks (title) VALUES ($1) RETURNING id`,
		input.Title,
	).Scan(&newID)
	if err != nil {
		http.Error(w, `{"error":"insert failed"}`, 500)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(Task{ID: newID, Title: input.Title})
}
