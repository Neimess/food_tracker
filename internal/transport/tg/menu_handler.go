package tg

import (
	"context"
	"fmt"
	"strings"

	"gopkg.in/telebot.v3"
)



func (bot *Bot) handleMenu(c telebot.Context) error {
	return bot.showMenu(c)
}

func (bot *Bot) buildMenu(ctx context.Context) (string, *telebot.ReplyMarkup, error) {
	foods, err := bot.svc.ListFoods(ctx)
	if err != nil {
		return "", nil, err
	}
	counts := bot.svc.SelectedFoodCounts()

	title := "Выбери блюда:"
	kb := &telebot.ReplyMarkup{}
	rows := [][]telebot.InlineButton{}

	for _, f := range foods {
		label := f.Name
		if c := counts[f.ID]; c > 0 {
			label = fmt.Sprintf("%s (%d)", f.Name, c)
		}
		add := telebot.InlineButton{Unique: "add_food", Text: "➕ " + label, Data: fmt.Sprintf("%d", f.ID)}
		rem := telebot.InlineButton{Unique: "remove_food", Text: "➖ " + f.Name, Data: fmt.Sprintf("%d", f.ID)}
		comp := telebot.InlineButton{Unique: "food_compose", Text: "ℹ️ состав", Data: fmt.Sprintf("%d", f.ID)}
		rows = append(rows, []telebot.InlineButton{add, rem, comp})
	}
	kb.InlineKeyboard = rows
	return title, kb, nil
}

func (bot *Bot) showMenu(c telebot.Context) error {
	text, kb, err := bot.buildMenu(context.Background())
	if err != nil {
		return c.Send("Ошибка при получении списка блюд")
	}
	return c.Send(text, kb)
}

func (bot *Bot) editMenu(c telebot.Context) error {
	text, kb, err := bot.buildMenu(context.Background())
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка обновления"})
	}
	if err := c.Edit(text, kb); err != nil {
		if strings.Contains(err.Error(), "message is not modified") {
			_, err = bot.B.EditReplyMarkup(c.Message(), kb)
			return err
		}
		_ = c.Edit(text+"\u200d", kb)
	}
	return nil
}
