package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Parsa-Sedigh/rebottle/internal/driver"
	"github.com/Parsa-Sedigh/rebottle/internal/models"
	"github.com/Parsa-Sedigh/rebottle/pkg/validation"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type config struct {
	port int
	db   struct{ dsn string }
}

type application struct {
	config     config
	logger     *zap.Logger
	version    string
	DB         models.Models
	DBPool     *sql.DB
	Session    *scs.SessionManager
	Validate   *validator.Validate
	Translator ut.Translator
	server     *http.Server
}

func newApp() *application {
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

	var cfg config
	flag.IntVar(&cfg.port, "port", 5001, "Server port to listen on")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("DSN"), "DSN")

	conn, err := driver.OpenDB(cfg.db.dsn)
	if err != nil {
		sugar.Fatal(err)
	}

	session := scs.New()
	session.Lifetime = 2 * time.Minute
	session.Store = memstore.New()

	validate, trans, err := validation.Register()
	if err != nil {
		sugar.Fatal(err)
	}

	app := application{
		config:     cfg,
		logger:     logger,
		DB:         models.NewModels(conn),
		DBPool:     conn,
		Session:    session,
		Validate:   validate,
		Translator: trans,
	}

	app.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	logger.Info(fmt.Sprintf("Starting Back end server on port %d", cfg.port))

	return &app
}

func (app *application) start() {
	app.logger.Info("Starting server...")
	defer app.logger.Sync()

	go func() {
		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.Fatal("Could not listen on", zap.String("addr", app.server.Addr), zap.Error(err))
		}
	}()

	app.logger.Info("Server is ready to handle requests", zap.String("addr", app.server.Addr))
	app.gracefulShutdown()
}

func (app *application) gracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	sig := <-quit

	app.logger.Info("Server is shutting down", zap.String("reason", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer func(DBPool *sql.DB) {
		err := DBPool.Close()
		if err != nil {
			app.logger.Info("err closing DB: ", zap.Error(err))
		}
	}(app.DBPool)

	app.server.SetKeepAlivesEnabled(false)

	if err := app.server.Shutdown(ctx); err != nil {
		app.logger.Fatal("Could not gracefully shutdown the server", zap.Error(err))
	}

	app.logger.Info("Server stopped")
}

func main() {
	app := newApp()
	app.start()
}
