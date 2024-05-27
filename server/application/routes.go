package application

import (
	"godb/handler"
	"godb/repository/level"
	"godb/repository/user"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *App) loadRoutes() {
	fsd := http.FileServer(http.Dir("dist"))
	fsa := http.FileServer(http.Dir("assets"))
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Handle("/*", fsd)
	router.Handle("/assets/*", http.StripPrefix("/assets/", fsa))

	router.Route("/user", a.loadUserRoutes)
	router.Route("/session", a.loadSessionRoutes)
	router.Route("/level", a.loadLevelRoutes)

	a.router = router
}

func (a *App) loadUserRoutes(router chi.Router) {
	userHandler := &handler.User{
		Repo: &user.SQLRepo{
			DB: a.db,
		},
	}

	router.Post("/", userHandler.Create)
	router.Get("/{id}", userHandler.GetByID)
	router.Delete("/{id}", userHandler.DeleteByID)
}

func (a *App) loadSessionRoutes(router chi.Router) {
	sessionHandler := &handler.Session{
		URepo: &user.SQLRepo{
			DB: a.db,
		},
		Key: []byte("Your key here"),
	}

	router.Get("/", sessionHandler.Create)
}

func (a *App) loadLevelRoutes(router chi.Router) {
	levelHandler := &handler.Level{
		LRepo: &level.SQLRepo{
			DB: a.db,
		},
		Key: []byte("Your key here"),
	}

	router.Post("/", levelHandler.Post)
	router.Post("/completed", levelHandler.Completed)
}
