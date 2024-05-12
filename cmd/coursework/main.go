package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/HeadGardener/coursework/internal/config"
	"github.com/HeadGardener/coursework/internal/handlers"
	"github.com/HeadGardener/coursework/internal/lib/auth"
	"github.com/HeadGardener/coursework/internal/server"
	"github.com/HeadGardener/coursework/internal/service"
	"github.com/HeadGardener/coursework/internal/storage"
)

const shutdownTimeout = 5 * time.Second

var confPath = flag.String("conf-path", "./config/.env", "path to config .env file")

func main() {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	conf, err := config.Init(*confPath)
	if err != nil {
		stop()
		log.Fatalf("[FATAL] error while initializing config: %s", err.Error())
	}

	db, err := storage.NewDB(ctx, conf.DBConfig)
	if err != nil {
		stop()
		log.Fatalf("[FATAL] error while establishing db connection: %s", err.Error())
	}

	rdb := storage.NewRedisDB(conf.RedisConfig)
	if err != nil {
		stop()
		log.Fatalf("[FATAL] error while establishing db connection: %s", err.Error())
	}

	var (
		userStorage  = storage.NewUserStorage(db)
		drinkStorage = storage.NewDrinkStorage(db)
		tokenStorage = storage.NewTokenStorage(rdb)
	)

	var (
		tokenManager = auth.NewTokenManager(&conf.TokensConfig)
	)

	var (
		authService  = service.NewAuthService(tokenManager, tokenStorage, userStorage)
		drinkService = service.NewDrinkService(drinkStorage)
	)

	handler := handlers.NewHandler(authService, drinkService)

	srv := &server.Server{}
	go func() {
		if err = srv.Run(conf.ServerConfig, handler.InitRoutes()); err != nil {
			log.Printf("[ERROR] failed to run server: %e", err)
		}
	}()
	log.Println("[INFO] server start working")

	<-ctx.Done()
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Printf("[INFO] server forced to shutdown: %s", err.Error())
	}

	if err = db.Close(); err != nil {
		log.Printf("[INFO] db connection forced to shutdown: %s", err.Error())
	}

	if err = rdb.Close(); err != nil {
		log.Printf("[INFO] redis db connection forced to shutdown: %s", err.Error())
	}

	log.Println("[INFO] server exiting")
}
