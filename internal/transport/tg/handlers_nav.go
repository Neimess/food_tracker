package tg

import (
	"context"
	"strconv"
	"strings"

	"gopkg.in/telebot.v3"
)

func (bot *Bot) handleNav(c telebot.Context) error {
	idx, err := strconv.Atoi(strings.TrimSpace(c.Data()))
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Некорректная навигация"})
	}
	return bot.renderDepartment(c, idx, true)
}

func (bot *Bot) handleClear(c telebot.Context) error {
	if err := bot.svc.Clear(context.Background()); err != nil {
		return c.Send("Ошибка при очистке корзины")
	}
	return c.Send("Корзина очищена 🧹")
}

func (bot *Bot) handleClearShow(c telebot.Context) error {
	if err := bot.svc.Clear(context.Background()); err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка удаления", ShowAlert: true})
	}
	_ = c.Respond(&telebot.CallbackResponse{Text: "Корзина очищена", ShowAlert: true})
	return bot.showMenu(c, true)
}
