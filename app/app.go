package app

import (
	"analytics/handler"
	"analytics/storage"
	"context"
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var dbUserName string
var dbPassword string

type App struct {
	router *mux.Router
	logger *log.Logger
}

// New function creates and set up instance of App
func New() (App, error) {
	flag.StringVar(&dbUserName, "dbuser", "default", "Specify DB username")
	flag.StringVar(&dbPassword, "dbpass", "gfFrjbkVSPiB", "Specify DB password")
	flag.Parse()

	l := log.New(os.Stdout, "analytics ", log.LstdFlags)

	l.Println("Connceting to DB...")
	s, err := storage.NewClickHouse(dbUserName, dbPassword)

	ev := handler.NewEventsManager(l, s)
	r := mux.NewRouter()

	getRouter := r.Methods("GET").Subrouter()
	getRouter.HandleFunc("/", ev.GetAvgForLastMin)

	postRouter := r.Methods("POST").Subrouter()
	postRouter.HandleFunc("/", ev.ProcessEvent)

	app := App{
		router: r,
		logger: l,
	}

	return app, err
}

func (a App) Start() {
	srv := &http.Server{
		Addr:         ":9090",
		Handler:      a.router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		a.logger.Println("Server is starting")
		log.Fatal(srv.ListenAndServe())
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	a.logger.Println("Received terminate, graceful shutdown...", sig)
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(tc)
}
