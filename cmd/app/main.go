package main

import (
	"database/sql"

	"github.com/Neimess/food_tracker/internal/repository"
	"github.com/Neimess/food_tracker/internal/service"
	"github.com/Neimess/food_tracker/internal/transport/tg"
	"github.com/Neimess/food_tracker/internal/transport/web"
	_ "modernc.org/sqlite"
)

func main() {
	db, _ := sql.Open("sqlite", "file:foods.db?_fk=1")
	defer db.Close()

	cats := repository.NewFoodCategoriesRepo(db)
	deps := repository.NewDepartmentsRepo(db)
	ings := repository.NewIngredientsRepo(db)
	foods := repository.NewFoodsRepo(db)
	fi := repository.NewFoodIngredientsRepo(db)
	cart := repository.NewCartRepo()

	websrv := web.NewServer(cats, deps, ings, foods, fi)
	go websrv.ListenAndServe(":8080")

	svcP := service.NewPlannerService(foods, fi, ings, cart)
	tgBot, err := tg.NewBot("7647335277:AAGzOJEg2ujkExYuukZOnkwtGWVsAkfAnTo", svcP)
	if err != nil {
		panic("stop")
	}
	tgBot.Start()
}
