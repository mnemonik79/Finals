package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	iterals "github.com/mnemonik79/Finals/internal/donetaskrepeat"
	"github.com/mnemonik79/Finals/internal/settings"
	"github.com/mnemonik79/Finals/internal/store"
	"github.com/mnemonik79/Finals/internal/tasks"
)

type ResponseJson struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func HandleNextDate(w http.ResponseWriter, r *http.Request) {

	strnow := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	strRepeat := r.URL.Query().Get("repeat")

	now, err := time.Parse(settings.Template, strnow)
	if err != nil {
		log.Fatal(err)
	}
	nextdate, err := iterals.NextDate(now, date, strRepeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = w.Write([]byte(nextdate))
	if err != nil {
		log.Fatal(err)
	}
}

func HandlePost(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t tasks.Task
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			http.Error(w, `{"error":"ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}
		id, err := store.CreateTask(t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp := ResponseJson{ID: id}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}

func HandleGet(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		task, err := store.GetTask(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(task); err != nil {
			http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}

	}
}

func HandlePut(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t tasks.Task
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			http.Error(w, `{"error":"ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}
		err = store.UpdateTask(t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
			http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}

func HandleRequests(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		err := store.DeleteTask(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
			http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}

func HandleTasksGet(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")
		tasksList, err := store.SearchTask(search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response := map[string][]tasks.Task{
			"tasks": tasksList,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}

func HandleTaskDone(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		err := store.DoneTask(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
			http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}
