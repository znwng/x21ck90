package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"todo-go/internal/repository"
	"todo-go/internal/utils"
)

func TodosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		todos, err := repository.GetAll()
		if err != nil {
			utils.WriteJSON(w, 500, map[string]string{"error": "db error"})
			return
		}
		utils.WriteJSON(w, 200, todos)

	case http.MethodPost:
		var payload struct {
			Item string `json:"item"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.WriteJSON(w, 400, map[string]string{"error": "bad body"})
			return
		}

		id, err := repository.Create(payload.Item)
		if err != nil {
			utils.WriteJSON(w, 500, map[string]string{"error": "insert failed"})
			return
		}

		utils.WriteJSON(w, 201, map[string]int64{"id": id})
	}
}

func TodoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteJSON(w, 400, map[string]string{"error": "invalid id"})
		return
	}

	switch r.Method {

	case http.MethodGet:
		t, err := repository.GetByID(id)
		if err != nil {
			utils.WriteJSON(w, 404, map[string]string{"error": "not found"})
			return
		}
		utils.WriteJSON(w, 200, t)

	case http.MethodPatch:
		if err := repository.Toggle(id); err != nil {
			utils.WriteJSON(w, 500, map[string]string{"error": "update failed"})
			return
		}
		utils.WriteJSON(w, 200, map[string]string{"message": "toggled"})
	}
}
