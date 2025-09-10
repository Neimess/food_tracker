package tg

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/telebot.v3"
)

func (bot *Bot) handleStart(c telebot.Context) error {
	return c.Send("Привет 👋 Я помогу тебе спланировать блюда и покупки.\n\n" +
		"Команды:\n/menu – выбрать блюда\n/cart – показать корзину\n/clear – очистить корзину")
}

func (bot *Bot) handleMenu(c telebot.Context) error {
	foods, err := bot.svc.ListFoods(context.Background())
	if err != nil {
		return c.Send("Ошибка при получении списка блюд")
	}

	if len(foods) == 0 {
		return c.Send("Блюда не найдены 🙈")
	}

	var rows [][]telebot.InlineButton
	for _, f := range foods {
		addBtn := telebot.InlineButton{
			Unique: "add_food",
			Text:   fmt.Sprintf("➕ %s", f.Name),
			Data:   fmt.Sprintf("%d", f.ID),
		}
		remBtn := telebot.InlineButton{
			Unique: "remove_food",
			Text:   fmt.Sprintf("➖ %s", f.Name),
			Data:   fmt.Sprintf("%d", f.ID),
		}
		compBtn := telebot.InlineButton{
			Unique: "food_compose",
			Text:   "ℹ️ состав",
			Data:   fmt.Sprintf("%d", f.ID),
    	}
		rows = append(rows, []telebot.InlineButton{addBtn, remBtn, compBtn})
	}

	return c.Send("Выбери блюда:", &telebot.ReplyMarkup{InlineKeyboard: rows})
}

func (bot *Bot) handleCart(c telebot.Context) error {
	items, _ := bot.svc.BuildCart(context.Background())
	if len(items) == 0 {
		return c.Send("Корзина пуста 🛒")
	}

	var msg strings.Builder
	msg.WriteString("🛒 Твоя корзина:\n\n")
	for _, it := range items {
		check := "⬜"
		if it.Checked {
			check = "✅"
		}
		msg.WriteString(fmt.Sprintf("%s %s — %.2f %s (%s)\n",
			check, it.Name, it.Qty, it.Unit, it.Department))
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
	foodID, _ := strconv.ParseInt(c.Data(), 10, 64)
	if err := bot.svc.AddFood(context.Background(), foodID); err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка при добавлении"})
	}
	return c.Respond(&telebot.CallbackResponse{Text: "Блюдо добавлено ✅"})
}

func (bot *Bot) handleRemoveFood(c telebot.Context) error {
	foodID, err := strconv.ParseInt(c.Data(), 10, 64)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "❌ неверный ID"})
	}

	if err := bot.svc.RemoveFood(context.Background(), foodID); err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка при удалении"})
	}

	return c.Respond(&telebot.CallbackResponse{Text: "Блюдо удалено 🗑"})
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
        return c.Send("Готовый продукт")
    }

    var sb strings.Builder
    if food != nil {
        sb.WriteString(fmt.Sprintf("📖 Состав *%s*:\n", food.Name))
    } else {
        sb.WriteString("📖 Состав блюда:\n")
    }
    // items: []domain.CompositionItem
    for _, it := range items {
        sb.WriteString(fmt.Sprintf("- %s: %.2f %s", it.Ingredient.Name, it.Quantity, it.Unit))
        if dep := it.Ingredient.Department; dep != "" {
            sb.WriteString(fmt.Sprintf("  _(отдел: %s)_", dep))
        }
        sb.WriteString("\n")
    }

    return c.Send(sb.String(), &telebot.SendOptions{ParseMode: "Markdown"})
}

