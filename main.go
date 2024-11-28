package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"

	auth "github.com/mnemonik79/Finals/internal/authtentification"
	"github.com/mnemonik79/Finals/internal/database"
	"github.com/mnemonik79/Finals/internal/handlers"
	"github.com/mnemonik79/Finals/internal/settings"
	"github.com/mnemonik79/Finals/internal/store"
)

func main() {
	db := database.InitializeDatabase()
	defer db.Close()

	store := store.NewStorage(db)

	r := chi.NewRouter()
	r.Handle("/*", http.FileServer(http.Dir("./web")))
	r.Get("/api/nextdate", handlers.HandleNextDate)
	r.Post("/api/task", auth.Authentification(handlers.HandlePost(store)))
	r.Get("/api/tasks", auth.Authentification(handlers.HandleTasksGet(store)))
	r.Get("/api/task", handlers.HandleGet(store))
	r.Put("/api/task", handlers.HandlePost(store))
	r.Post("/api/task/done", auth.Authentification(handlers.HandleTaskDone(store)))
	r.Delete("/api/task", handlers.HandleRequests(store))
	r.Post("/api/signin", auth.HandleSiginingIn)

	environment := settings.GetEnv()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	err := http.ListenAndServe(":"+environment.Port, r)
	if err != nil {
		fmt.Println("Не удалось запустить сервер:\n", err)
	} else {
		infoLog.Print("Запуск сервера")
	}

}
