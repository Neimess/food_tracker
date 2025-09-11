package tg

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"gopkg.in/telebot.v3"
)

func (bot *Bot) handleMenu(c telebot.Context) error         { return bot.showMenu(c, false) }
func (bot *Bot) handleMenuCallback(c telebot.Context) error { return bot.showMenu(c, true) }

func (bot *Bot) buildMenu(ctx context.Context) (string, *telebot.ReplyMarkup, error) {
	foods, err := bot.svc.ListFoods(ctx)
	if err != nil {
		return "", nil, err
	}
	counts := bot.svc.SelectedFoodCounts()

	title := "Выбери блюда:"
	kb := &telebot.ReplyMarkup{}
	var rows []telebot.Row

	for _, f := range foods {
		label := f.Name
		if c := counts[f.ID]; c > 0 {
			label = fmt.Sprintf("%s (%d)", f.Name, c)
		}

		dataAdd, _ := json.Marshal(map[string]int64{"id": f.ID})
		btnAdd := BtnAddFood
		btnAdd.Text = "➕ " + label
		btnAdd.Data = string(dataAdd)

		btnRem := BtnRemoveFood
		btnRem.Text = "➖ " + f.Name
		btnRem.Data = strconv.FormatInt(f.ID, 10)

		btnComp := BtnFoodCompose
		btnComp.Text = "ℹ️ Состав"
		btnComp.Data = strconv.FormatInt(f.ID, 10)

		rows = append(rows, kb.Row(btnAdd, btnRem, btnComp))
	}

	btnCart := BtnCartShow
	btnCart.Text, btnCart.Data = "🛒 Корзина", "show"
	btnClear := BtnCartClear
	btnClear.Text, btnClear.Data = "🧹 Очистить", "clear"
	rows = append(rows, kb.Row(btnCart, btnClear))

	btnSettings := kb.WebApp("⚙️ Настройки", &telebot.WebApp{
		URL: bot.webAppURL,
	})
	rows = append(rows, kb.Row(btnSettings))

	kb.Inline(rows...)
	return title, kb, nil
}

func (bot *Bot) showMenu(c telebot.Context, edit bool) error {
	text, kb, err := bot.buildMenu(context.Background())
	if err != nil {
		log.Println(err)
		return c.Send("Ошибка при получении списка блюд")
	}

	var msg *telebot.Message
	if edit {
		msg, err = bot.SendOrEdit(
			c, text, kb,
			WithEdit(),
			WithEditMarkupIfSame(),
			WithForceChangeOnSame(),
		)
	} else {
		msg, err = bot.SendOrEdit(c, text, kb)
	}
	if err != nil {
		return err
	}

	if msg != nil {
		bot.mu.Lock()
		bot.lastMenus[msg.Chat.ID] = msg.ID
		bot.mu.Unlock()
	}
	return nil
}

func (bot *Bot) editMenu(c telebot.Context) error {
	text, kb, err := bot.buildMenu(context.Background())
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка обновления"})
	}

	_, err = bot.SendOrEdit(
		c, text, kb,
		WithEdit(),
		WithEditMarkupIfSame(),
		WithForceChangeOnSame(),
	)
	return err
}
