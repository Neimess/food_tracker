package tg

func (bot *Bot) registerHandlers() {
	b := bot.B

	b.Handle("/start", bot.handleStart)
	b.Handle("/menu", bot.handleMenu)
	b.Handle("/cart", bot.handleCart)
	b.Handle("/clear", bot.handleClear)

	b.Handle(&BtnAddFood, bot.handleAddFood)
	b.Handle(&BtnRemoveFood, bot.handleRemoveFood)
	b.Handle(&BtnFoodCompose, bot.handleFoodCompose)

	b.Handle(&BtnCartShow, bot.handleCartShow)
	b.Handle(&BtnCartClear, bot.handleClearShow)
	b.Handle(&BtnMenuBack, bot.handleMenuCallback)

	b.Handle(&BtnToggle, bot.handleToggle)
	b.Handle(&BtnNav, bot.handleNav)

}
