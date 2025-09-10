package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Neimess/food_tracker/internal/config"
	"github.com/Neimess/food_tracker/internal/repository"
	"github.com/Neimess/food_tracker/internal/service"
	"github.com/Neimess/food_tracker/internal/transport/tg"
	"github.com/Neimess/food_tracker/internal/transport/web"
	_ "modernc.org/sqlite"
)

func main() {
	cfg := config.MustLoad()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	db, _ := sql.Open("sqlite", "file:foods.db?_fk=1")
	defer db.Close()

	cats := repository.NewFoodCategoriesRepo(db)
	deps := repository.NewDepartmentsRepo(db)
	ings := repository.NewIngredientsRepo(db)
	foods := repository.NewFoodsRepo(db)
	fi := repository.NewFoodIngredientsRepo(db)
	cart := repository.NewCartRepo()

	websrv := web.NewServer(cats, deps, ings, foods, fi)
	srvErr := make(chan error, 1)

	go func() {
		if err := websrv.ListenAndServe(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srvErr <- err
		}
		close(srvErr)
	}()

	svcP := service.NewPlannerService(foods, fi, ings, cart)

	tgBot, err := tg.NewBot(ctx, cfg.TGToken, svcP)
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
	tgBot.Stop()
	log.Println("done")
}
