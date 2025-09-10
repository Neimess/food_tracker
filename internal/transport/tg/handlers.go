package tg

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Neimess/food_tracker/internal/domain"
	"gopkg.in/telebot.v3"
)

func (bot *Bot) handleStart(c telebot.Context) error {
	return c.Send("Привет 👋 Я помогу тебе спланировать блюда и покупки.\n\n" +
		"Команды:\n/menu – выбрать блюда\n/cart – показать корзину\n/clear – очистить корзину")
}

func (bot *Bot) handleCart(c telebot.Context) error {
	items, _ := bot.svc.BuildCart(context.Background())
	if len(items) == 0 {
		return c.Send("Корзина пуста 🛒")
	}

	var msg strings.Builder
	msg.WriteString("🛒 Твоя корзина:\n")
	groups := make(map[string][]domain.CartItem)
	for _, it := range items {
		groups[it.Department] = append(groups[it.Department], it)
	}
	deps := make([]string, 0, len(groups))
	for dep := range groups {
		deps = append(deps, dep)
	}

	sort.Strings(deps)

	for _, dep := range deps {
		msg.WriteString(fmt.Sprintf("📂 %s\n", dep))
		sort.Slice(groups[dep], func(i, j int) bool {
			return groups[dep][i].Name < groups[dep][j].Name
		})
		for _, it := range groups[dep] {
			check := "⬜"
			if it.Checked {
				check = "✅"
			}
			msg.WriteString(fmt.Sprintf(" %s %s — %d %s\n",
				check, it.Name, it.Qty, it.Unit))
		}
		msg.WriteString("\n")
	}
	return c.Send(msg.String())
}

func (bot *Bot) handleClear(c telebot.Context) error {
	if err := bot.svc.Clear(context.Background()); err != nil {
		return c.Send("Ошибка при очистке корзины")
	}
	return c.Send("Корзина очищена 🧹")
}

func (bot *Bot) handleAddFood(c telebot.Context) error {
	id, _ := strconv.ParseInt(c.Data(), 10, 64)
	if err := bot.svc.AddFood(context.Background(), id); err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка при добавлении"})
	}
	_ = c.Respond(&telebot.CallbackResponse{Text: "Добавлено ✅"})
	return bot.editMenu(c)
}

func (bot *Bot) handleRemoveFood(c telebot.Context) error {
	id, err := strconv.ParseInt(c.Data(), 10, 64)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Неверный ID"})
	}
	if err := bot.svc.RemoveFood(context.Background(), id); err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Отсутствует в корзине"})
	}
	_ = c.Respond(&telebot.CallbackResponse{Text: "Убрано 🗑"})
	return bot.editMenu(c)
}

func (bot *Bot) handleFoodCompose(c telebot.Context) error {
	foodID, err := strconv.ParseInt(c.Data(), 10, 64)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "❌ неверный ID"})
	}

	food, items, err := bot.svc.GetComposition(context.Background(), foodID)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Не удалось получить состав"})
	}
	if len(items) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "Готовый продукт", ShowAlert: true})
	}

	var sb strings.Builder
	if food != nil {
		sb.WriteString(fmt.Sprintf("📖 Состав %s:\n", food.Name))
	} else {
		sb.WriteString("📖 Состав блюда:\n")
	}
	for _, it := range items {
		sb.WriteString(fmt.Sprintf("- %s: %d %s", it.Ingredient.Name, it.Quantity, it.Unit))
		sb.WriteString("\n")
	}
	text := sb.String()
	if err := c.Respond(&telebot.CallbackResponse{
		Text:      text,
		ShowAlert: true,
	}); err == nil {
		return nil
	}
	return c.Send("📖 " + text)
}
