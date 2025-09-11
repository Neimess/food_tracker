package tg

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Neimess/food_tracker/internal/config"
	"github.com/Neimess/food_tracker/internal/service"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type Bot struct {
	B   *telebot.Bot
	svc *service.PlannerService

	lastMenus map[int64]int
	mu        sync.Mutex

	cfg *config.TelegramConfig
}

func NewBot(ctx context.Context, cfg *config.TelegramConfig, svc *service.PlannerService) (*Bot, error) {
	pref := telebot.Settings{
		Token: cfg.Token,
		Poller: &telebot.Webhook{
			Endpoint: &telebot.WebhookEndpoint{
				PublicURL: cfg.URL + cfg.Token,
			},
			Listen: cfg.Address,
		},
		Synchronous: false,
		OnError: func(err error, c telebot.Context) {
			chat := "<nil>"
			if c != nil && c.Sender() != nil {
				chat = fmt.Sprintf("%s (%d)", c.Sender().Username, c.Sender().ID)
			}
			log.Printf("[ERROR] chat=%s err=%v", chat, err)
		},
	}
	b, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}
	return &Bot{B: b, svc: svc, lastMenus: make(map[int64]int), cfg: cfg}, nil
}

func (bot *Bot) Start() {
	log.Printf("Telegram bot started at: %s", bot.cfg.Address)
	bot.B.Use(ShortLogger)
	bot.B.Use(middleware.Whitelist(bot.cfg.AllowedUsers...))
	bot.B.Use(middleware.AutoRespond())
	bot.B.Use(middleware.Recover())
	bot.registerHandlers()
	bot.B.Start()
}

func (bot *Bot) Stop() {
	bot.B.Stop()
}
