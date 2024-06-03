package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"godb/static"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	router http.Handler
	db     *sql.DB
}

func New() *App {
	_, err := os.Stat("./database.db")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err := os.WriteFile("database.db", []byte(""), 0755)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}

	app := &App{
		db: db,
	}

	app.loadRoutes()

	return app
}

func (a *App) Start(ctx context.Context, port string) error {
	server := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%s", port),
		Handler: a.router,
	}

	err := a.db.Ping()
	if err != nil {
		return fmt.Errorf("failed to connect to sqlite: %w", err)
	}

	statement, err := a.db.Prepare(`CREATE TABLE IF NOT EXISTS` + "`user`" + `(
		user_id INTEGER PRIMARY KEY, 
		email TEXT, 
		password TEXT,
		first_name TEXT, 
		last_name TEXT)`)
	if err != nil {
		panic(err)
	}
	statement.Exec()

	static.ConvertStaticToDb(a.db)

	defer func() {
		if err := a.db.Close(); err != nil {
			fmt.Println("Failed to close SQLite", err)
		}
	}()

	fmt.Printf("Server is running... [http://127.0.0.1:%s]\n", port)

	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {

	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}
}
