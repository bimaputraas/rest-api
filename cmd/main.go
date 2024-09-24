package main

import (
	"github.com/bimaputraas/rest-api/internal/api/controller"
	"github.com/bimaputraas/rest-api/internal/api/middleware"
	"github.com/bimaputraas/rest-api/internal/api/routes"
	"github.com/bimaputraas/rest-api/internal/config"
	"github.com/bimaputraas/rest-api/internal/repository/mysql"
	"github.com/bimaputraas/rest-api/internal/usecase"
	_ "github.com/bimaputraas/rest-api/pkg/database"
	pkgdatabase "github.com/bimaputraas/rest-api/pkg/database"
	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	db, err := pkgdatabase.NewGorm(cfg.MySQLURI)
	if err != nil {
		log.Fatal(err)
	}
	repo := mysql.New(db)
	uc := usecase.New(repo, cfg)
	mw := middleware.New(uc)
	ctr := controller.New(uc)

	r := routes.New(mw, ctr)

	log.Fatal(r.Run(":8080"))
}
