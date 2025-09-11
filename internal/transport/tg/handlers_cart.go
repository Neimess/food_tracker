package tg

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/Neimess/food_tracker/internal/domain"
	"gopkg.in/telebot.v3"
)

func (bot *Bot) handleStart(c telebot.Context) error {
	return c.Send("Привет 👋 Я помогу тебе спланировать блюда и покупки.\n\n" +
		"Команды:\n/menu – выбрать блюда\n/cart – показать корзину\n/clear – очистить корзину")
}

func (bot *Bot) handleCart(c telebot.Context) error     { return bot.renderDepartment(c, 0, false) }
func (bot *Bot) handleCartShow(c telebot.Context) error { return bot.renderDepartment(c, 0, true) }

func (bot *Bot) renderDepartment(c telebot.Context, reqIndex int, edit bool) error {
	items, _ := bot.svc.BuildCart(context.Background())

	if len(items) == 0 {
		kb := &telebot.ReplyMarkup{}
		btn := BtnMenuBack
		btn.Text, btn.Data = "⬅️ В меню", "back"
		kb.Inline(kb.Row(btn))
		if edit {
			_, err := bot.SendOrEdit(c, "Корзина пуста 🛒", kb, WithEdit())
			return err
		}
		_, err := bot.SendOrEdit(c, "Корзина пуста 🛒", kb)
		return err
	}

	groups := make(map[string][]domain.CartItem)
	for _, it := range items {
		groups[it.Department] = append(groups[it.Department], it)
	}

	depts := make([]string, 0, len(groups))
	for d := range groups {
		depts = append(depts, d)
	}
	sort.Strings(depts)

	n := len(depts)
	if n == 0 {
		kb := &telebot.ReplyMarkup{}
		btn := BtnMenuBack
		btn.Text, btn.Data = "⬅️ В меню", "back"
		kb.Inline(kb.Row(btn))
		if edit {
			_, err := bot.SendOrEdit(c, "Корзина пуста 🛒", kb, WithEdit())
			return err
		}
		_, err := bot.SendOrEdit(c, "Корзина пуста 🛒", kb)
		return err
	}

	idx := ((reqIndex % n) + n) % n
	dept := depts[idx]

	title := fmt.Sprintf("🗂 *%s* — %d поз.", dept, len(groups[dept]))

	kb := &telebot.ReplyMarkup{}
	var rows []telebot.Row

	for _, it := range groups[dept] {
		data, _ := json.Marshal(TogglePayload{Index: idx, IngredientID: it.IngredientID})
		check := "⬜"
		if it.Checked {
			check = "✅"
		}

		btn := BtnToggle
		btn.Text = fmt.Sprintf("%s %s (%d %s)", check, it.Name, it.Qty, it.Unit)
		btn.Data = string(data)
		rows = append(rows, kb.Row(btn))
	}

	if n > 1 {
		prevIdx := (idx - 1 + n) % n
		nextIdx := (idx + 1) % n

		btnPrev := BtnNav
		btnPrev.Text = "⬅️ " + depts[prevIdx]
		btnPrev.Data = strconv.Itoa(prevIdx)

		btnNext := BtnNav
		btnNext.Text = depts[nextIdx] + " ➡️"
		btnNext.Data = strconv.Itoa(nextIdx)

		if n == 2 {
			btnNext.Text = "Перейти к: " + depts[nextIdx]
			rows = append(rows, kb.Row(btnNext))
		} else {
			rows = append(rows, kb.Row(btnPrev, btnNext))
		}
	}

	btnMenu := BtnMenuBack
	btnMenu.Text, btnMenu.Data = "⬅️ В меню", "back"
	rows = append(rows, kb.Row(btnMenu))

	kb.Inline(rows...)
	if edit {
		_, err := bot.SendOrEdit(c, title, kb, WithEdit(), WithEditMarkupIfSame(), WithForceChangeOnSame())
		return err
	}
	_, err := bot.SendOrEdit(c, title, kb)
	return err
}
