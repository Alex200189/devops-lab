package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// Task представляет структуру задачи
type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// Хранилище в оперативной памяти (In-Memory) с поддержкой конкурентности
var (
	tasks       = []Task{}
	nextID      = 1
	tasksMu     sync.Mutex
	tasksPathPrefix = "/tasks"
)

func main() {
	mux := http.NewServeMux()

	// Регистрация эндпоинтов (синтаксис Go 1.22+)
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /metrics", metricsHandler)
	mux.HandleFunc("GET /tasks", getTasksHandler)
	mux.HandleFunc("POST /tasks", createTasksHandler)

	// Обертка для обработки неизвестных путей (404)
	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, совпадает ли запрос с зарегистрированными маршрутами
		_, pattern := mux.Handler(r)
		if pattern == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"Not Found"}`))
			return
		}
		mux.ServeHTTP(w, r)
	})

	fmt.Println("Сервер запущен на порту :8080...")
	if err := http.ListenAndServe(":8080", rootHandler); err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}
}

// GET /health
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// GET /metrics
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("# HELP http_requests_total Total number of HTTP requests.\n# TYPE http_requests_total counter\nhttp_requests_total 0\n"))
}

// GET /tasks
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	tasksMu.Lock()
	defer tasksMu.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

// POST /tasks
func createTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input struct {
		Title string `json:"title"`
	}

	// Декодируем входящий JSON
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Invalid JSON structure"}`))
		return
	}

	// Валидация: пустой или отсутствующий title
	if strings.TrimSpace(input.Title) == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Title is required and cannot be empty"}`))
		return
	}

	// Безопасное добавление в in-memory слайс
	tasksMu.Lock()
	newTask := Task{
		ID:    nextID,
		Title: input.Title,
	}
	nextID++
	tasks = append(tasks, newTask)
	tasksMu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

