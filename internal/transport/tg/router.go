package tg

import "gopkg.in/telebot.v3"

func (bot *Bot) registerHandlers() {
	b := bot.b

	b.Handle("/start", bot.handleStart)
	b.Handle("/menu", bot.handleMenu)
	b.Handle("/cart", bot.handleCart)
	b.Handle("/clear", bot.handleClear)

	// inline кнопки на добавление блюда
	b.Handle(&telebot.Btn{Unique: "add_food"}, bot.handleAddFood)
	b.Handle(&telebot.Btn{Unique: "remove_food"}, bot.handleRemoveFood)
	b.Handle(&telebot.Btn{Unique: "food_compose"}, bot.handleFoodCompose)
}
