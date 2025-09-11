package tg

import "gopkg.in/telebot.v3"

var (
	BtnAddFood     = telebot.Btn{Unique: "add_food"}
	BtnRemoveFood  = telebot.Btn{Unique: "remove_food"}
	BtnFoodCompose = telebot.Btn{Unique: "food_compose"}

	BtnCartShow  = telebot.Btn{Unique: "cart_show"}
	BtnCartClear = telebot.Btn{Unique: "clear_cart"}
	BtnMenuBack  = telebot.Btn{Unique: "menu_back"}

	BtnToggle = telebot.Btn{Unique: "toggle"}
	BtnNav    = telebot.Btn{Unique: "nav"}
)
