package main

import (
	"log"
	"net/http"

	"github.com/Rumm1/eduhub-backend/internal/app"
	"github.com/Rumm1/eduhub-backend/internal/config"
	"github.com/Rumm1/eduhub-backend/internal/platform/db"
	platformjwt "github.com/Rumm1/eduhub-backend/internal/platform/jwt"
)

func main() {
	cfg := config.Load()

	postgresPool, err := db.NewPostgresPool(cfg.Database.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer postgresPool.Close()

	jwtManager := platformjwt.NewManager(
		cfg.JWT.AccessSecret,
		cfg.JWT.AccessTTLMinutes,
	)

	router := app.NewRouter(postgresPool, jwtManager)

	log.Println("EduHub backend started on port", cfg.App.Port)

	if err := http.ListenAndServe(":"+cfg.App.Port, router); err != nil {
		log.Fatal(err)
	}
}
