package main

import (
	"log/slog"
	"os"

	"github.com/arsenydubrovin/ad-submission/src/internal/controller"
	"github.com/arsenydubrovin/ad-submission/src/internal/models"
	echo "github.com/labstack/echo/v4"
	slogecho "github.com/samber/slog-echo"
)

func main() {
	cfg := loadConfig()

	log := setupLogger(cfg.app.env)
	slog.SetDefault(log)

	db, err := openDB(cfg.postgres.host, cfg.postgres.port, cfg.postgres.user, cfg.postgres.db)
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

	err = ctrl.Serve(cfg.app.httpPort)
	if err != nil {
		log.Error("failed to start server", wrapErr(err))
		os.Exit(1)
	}
}
