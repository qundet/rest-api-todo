package handlers

import (
	"encoding/json"
	"net/http"
	"rest-api-todo/internal/database"
	"rest-api-todo/internal/models"
	"strconv"
	"strings"
)

type Handlers struct {
	store *database.TaskStore
}

func NewHandlers(store *database.TaskStore) *Handlers {
	return &Handlers{store: store}
}

func responseWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(payload)
}

func responseWithError(w http.ResponseWriter, statusCode int, message string) {
	responseWithJSON(w, statusCode, map[string]string{"error": message})
}

func (h *Handlers) GetAllTask(w http.ResponseWriter, r *http.Request) {

	tasks, err := h.store.GetAll()

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Get data error")
		return
	}

	responseWithJSON(w, http.StatusOK, tasks)
}

func (h *Handlers) GetTask(w http.ResponseWriter, r *http.Request) {

	id, err := getIdFromPath(r.URL.Path)

	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Bad task id")
		return
	}

	task, err := h.store.GetById(id)

	if err != nil {
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, task)
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {

	var input models.CreateTaskInput

	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Bad request body")
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		responseWithError(w, http.StatusBadRequest, "Title can not be empty")
		return
	}

	createdTask, err := h.store.CreateTask(input)

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseWithJSON(w, http.StatusCreated, createdTask)
}

func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {

	id, err := getIdFromPath(r.URL.Path)

	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Bad task id")
		return
	}

	var input models.UpdateTaskInput

	err = json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Bad task body")
		return
	}

	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		responseWithError(w, http.StatusBadRequest, "Title can not be empty")
		return
	}

	updatedTask, err := h.store.UpdateTask(id, input)

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, updatedTask)
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromPath(r.URL.Path)

	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Bad task id")
		return
	}

	err = h.store.Delete(id)

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			responseWithError(w, http.StatusBadRequest, "Bad task id")
		} else {
			responseWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	responseWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
func getIdFromPath(path string) (int, error) {

	pathParts := strings.Split(strings.TrimPrefix(path, "/tasks/"), "/")

	id, err := strconv.Atoi(pathParts[0])

	return id, err
}
