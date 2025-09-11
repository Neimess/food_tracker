package tg

import (
	"context"
	"log"
	"strings"

	"gopkg.in/telebot.v3"
)

//nolint:unused
func (bot *Bot) updateKeyboardOnInvalidate() {
	for range bot.svc.InvalidateChan {
		bot.mu.Lock()
		menus := make(map[int64]int, len(bot.lastMenus))
		for chatID, msgID := range bot.lastMenus {
			menus[chatID] = msgID
		}
		bot.mu.Unlock()

		for chatID, msgID := range menus {
			text, kb, err := bot.buildMenu(context.Background())
			if err != nil {
				log.Printf("updateKeyboardOnInvalidate buildMenu failed: %v", err)
				continue
			}

			_, err = bot.B.Edit(
				&telebot.Message{Chat: &telebot.Chat{ID: chatID}, ID: msgID},
				text,
				&telebot.SendOptions{ReplyMarkup: kb},
			)
			if err != nil {
				if strings.Contains(err.Error(), "message is not modified") {
					continue
				}
				log.Printf("updateKeyboardOnInvalidate edit failed: chat=%d msg=%d err=%v", chatID, msgID, err)
			}
		}
	}
}
