package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Neimess/food_tracker/internal/config"
	"github.com/Neimess/food_tracker/internal/repository"
	"github.com/Neimess/food_tracker/internal/service"
	"github.com/Neimess/food_tracker/internal/transport/tg"
	"github.com/Neimess/food_tracker/internal/transport/web"
	_ "modernc.org/sqlite"
)

var (
	Version   = "1.0.0"
	BuildTime = "unknown"
	Commit = "unknown"
)

func main() {
	cfg := config.MustLoad()
	cfg.App.Version = Version
	log.Printf(`
===========================================
  Application started
-------------------------------------------
  Version   : %s
  BuildTime : %s
  Commit    : %s
  Env       : %s
===========================================
`, Version, BuildTime, Commit, cfg.App.Env)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	dsn := fmt.Sprintf("file:%s?_fk=1", cfg.DB.Path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatalf("create db failed: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("db ping failed")
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}()
	

	cats := repository.NewFoodCategoriesRepo(db)
	deps := repository.NewDepartmentsRepo(db)
	ings := repository.NewIngredientsRepo(db)
	foods := repository.NewFoodsRepo(db)
	fi := repository.NewFoodIngredientsRepo(db)
	cart := repository.NewCartRepo("cart.json")

	websrv := web.NewServer(&cfg.HTTPServer, cats, deps, ings, foods, fi)
	srvErr := make(chan error, 1)

	go func() {
		if err := websrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srvErr <- err
		}
		close(srvErr)
	}()

	svcP := service.NewPlannerService(cfg.Cache.Path, foods, fi, ings, cart)
	if err := svcP.Preload(ctx); err != nil {
		log.Printf("preload failed: %v", err)
	}

	tgBot, err := tg.NewBot(ctx, &cfg.Telegram, svcP)
	if err != nil {
		panic(err)
	}

	go tgBot.Start()

	select {
	case <-ctx.Done():
	case err := <-srvErr:
		if err != nil {
			log.Printf("http server error: %v", err)
		}
	}

	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Println("shutting down...")
	if err := websrv.Shutdown(shutCtx); err != nil {
		log.Printf("http shutdown: %v", err)
	}
	if err := svcP.Close(shutCtx); err != nil {
		log.Printf("service shutting down failed: %v", err)
	}
	tgBot.Stop()
	log.Println("done")
}
