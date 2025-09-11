package tg

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"gopkg.in/telebot.v3"
)

func (bot *Bot) handleAddFood(c telebot.Context) error {
	var payload DefPayload
	if err := json.Unmarshal([]byte(c.Data()), &payload); err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Некорректные данные"})
	}
	if err := bot.svc.AddFood(context.Background(), payload.ID); err != nil {
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
		sb.WriteString("📖 Состав " + food.Name + ":\n")
	} else {
		sb.WriteString("📖 Состав блюда:\n")
	}
	for _, it := range items {
		sb.WriteString("- " + it.Ingredient.Name)
		if it.Quantity > 0 || it.Unit != "" {
			sb.WriteString(": ")
			if it.Quantity > 0 {
				sb.WriteString(strconv.Itoa(int(it.Quantity)))
			}
			if it.Unit != "" {
				sb.WriteString(" " + it.Unit)
			}
		}
		sb.WriteString("\n")
	}
	text := sb.String()

	if err := c.Respond(&telebot.CallbackResponse{Text: text, ShowAlert: true}); err == nil {
		return nil
	}
	return c.Send("📖 " + text)
}

func (bot *Bot) handleToggle(c telebot.Context) error {
	var payload TogglePayload
	if err := json.Unmarshal([]byte(c.Data()), &payload); err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Некорректные данные"})
	}

	if err := bot.svc.ToggleChecked(context.Background(), payload.IngredientID); err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка: " + err.Error()})
	}

	return bot.renderDepartment(
		c,
		payload.Index,
		true,
	)
}
