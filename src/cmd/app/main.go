package main

import (
	"log/slog"
	"os"

	"github.com/arsenydubrovin/ad-submission/src/internal/config"
	"github.com/arsenydubrovin/ad-submission/src/internal/controller"
	"github.com/arsenydubrovin/ad-submission/src/internal/models"
	echo "github.com/labstack/echo/v4"
	slogecho "github.com/samber/slog-echo"
)

func main() {
	err := config.Load(".env")
	if err != nil {
		slog.Error("failed to load config", wrapErr(err))
		os.Exit(1)
	}

	appCfg, err := config.NewApplicationConfig()
	if err != nil {
		slog.Error("failed to load application config", wrapErr(err))
		os.Exit(1)
	}

	log := setupLogger(appCfg.Env())
	slog.SetDefault(log)

	dbCfg, err := config.NewPostgresConfig()
	if err != nil {
		log.Error("failed to load database config", wrapErr(err))
		os.Exit(1)
	}
	db, err := openDB(dbCfg.DSN())
	if err != nil {
		log.Error("failed to open database", wrapErr(err))
		os.Exit(1)
	}
	log.Info("connection to the database is set")

	md := models.New(db)

	e := echo.New()
	e.Use(slogecho.NewWithConfig(log,
		slogecho.Config{
			WithRequestBody: true,
		}))

	ctrl := controller.New(e, &md)
	ctrl.RegisterRoutes()
	log.Info("controller is initialized")

	httpCfg, err := config.NewHTTPConfig()
	if err != nil {
		log.Error("failed to load http server config", wrapErr(err))
		os.Exit(1)
	}

	err = ctrl.Serve(httpCfg.Port())
	if err != nil {
		log.Error("failed to start server", wrapErr(err))
		os.Exit(1)
	}
}
