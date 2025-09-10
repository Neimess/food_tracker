package tg

import (
	"context"
	"log"

	"time"

	"github.com/Neimess/food_tracker/internal/service"
	"gopkg.in/telebot.v3"
)

type Bot struct {
	B   *telebot.Bot
	svc *service.PlannerService
}

func NewBot(ctx context.Context, token string, svc *service.PlannerService) (*Bot, error) {
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}
	return &Bot{B: b, svc: svc}, nil
}

func (bot *Bot) Start() {
	log.Println("Telegram bot started")
	bot.registerHandlers()
	bot.B.Start()
}

func (bot *Bot) Stop() {
	bot.B.Stop()
}