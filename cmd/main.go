package main

import (
	"log"

	"github.com/bimaputraas/rest-api/internal/api/controller"
	"github.com/bimaputraas/rest-api/internal/api/middleware"
	"github.com/bimaputraas/rest-api/internal/api/routes"
	"github.com/bimaputraas/rest-api/internal/config"
	"github.com/bimaputraas/rest-api/internal/repository"
	gormRepository "github.com/bimaputraas/rest-api/internal/repository/gorm"
	redisRepository "github.com/bimaputraas/rest-api/internal/repository/redis"
	"github.com/bimaputraas/rest-api/internal/usecase"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	gormDB, err := gorm.Open(mysql.Open(cfg.MySQLURI), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	redisClient := &redis.Client{}

	uc := usecase.New(repository.New(gormRepository.NewStorage(gormDB), redisRepository.NewCache(redisClient)), cfg)

	r := routes.New(middleware.New(uc), controller.New(uc))

	log.Fatal(r.Run(":8080"))
}
