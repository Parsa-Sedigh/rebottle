package main

import (
	"flag"
	"fmt"
	"github.com/Parsa-Sedigh/rebottle/internal/driver"
	"github.com/Parsa-Sedigh/rebottle/internal/models"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	port int
	db   struct{ dsn string }
}

type application struct {
	config  config
	logger  *zap.Logger
	version string
	DB      models.Models
	Session *scs.SessionManager
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.logger.Info(fmt.Sprintf("Starting Back end server on port %d", app.config.port))

	return srv.ListenAndServe()
}

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	//infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	//errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	var cfg config
	flag.IntVar(&cfg.port, "port", 5001, "Server port to listen on")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("DSN"), "DSN")

	conn, err := driver.OpenDB(cfg.db.dsn)
	if err != nil {
		sugar.Fatal(err)
	}

	defer conn.Close()

	session := scs.New()
	session.Lifetime = 2 * time.Minute
	session.Store = memstore.New()

	app := application{
		config:  cfg,
		logger:  logger,
		DB:      models.NewModels(conn),
		Session: session,
	}

	err = app.serve()
	if err != nil {
		app.logger.Fatal(err.Error())
	}
}
