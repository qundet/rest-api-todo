package main

import (
	"log"
	"net/http"
	"os"
	"rest-api-todo/internal/database"
	"rest-api-todo/internal/handlers"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://taskuser:taskpass@localhost:5432/taskdb"
	}

	serverPort := os.Getenv("SERVICE_PORT")
	if serverPort == "" {
		serverPort = "8888"
	}

	log.Printf("Server is starting on %s", serverPort)

	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("Successful connect to db")

	handler := handlers.NewHandlers(database.NewTaskStore(db))

	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", methodHandler(handler.GetAllTask, http.MethodGet))
	mux.HandleFunc("/tasks/create", methodHandler(handler.CreateTask, http.MethodPost))
	mux.HandleFunc("/tasks/", taskIdHandler(handler))

	loggedMux := loggingMiddleware(mux)

	serverAddress := ":" + serverPort

	err = http.ListenAndServe(serverAddress, loggedMux)

	if err != nil {
		log.Fatal(err)
	}
}

func methodHandler(handlerFunc http.HandlerFunc, allowedMethod string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
		handlerFunc(w, r)
	}
}

func taskIdHandler(handler *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetTask(w, r)
		case http.MethodPut:
			handler.UpdateTask(w, r)
		case http.MethodDelete:
			handler.DeleteTask(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}