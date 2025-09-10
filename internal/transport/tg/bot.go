package tg

import (
	"log"

	"time"

	"github.com/Neimess/food_tracker/internal/service"
	"gopkg.in/telebot.v3"
)

type Bot struct {
	b   *telebot.Bot
	svc *service.PlannerService
}

func NewBot(token string, svc *service.PlannerService) (*Bot, error) {
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}
	return &Bot{b: b, svc: svc}, nil
}

func (bot *Bot) Start() {
	log.Println("Telegram bot started")
	bot.registerHandlers()
	bot.b.Start()
}
