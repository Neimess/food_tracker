package tg

import (
	"log"

	"gopkg.in/telebot.v3"
)

func ShortLogger(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		upd := c.Update()
		switch {
		case upd.Message != nil:
			m := upd.Message
			log.Printf("[MSG] from=%s (id=%d) text=%q",
				m.Sender.Username, m.Sender.ID, m.Text)

		case upd.Callback != nil:
			cb := upd.Callback
			log.Printf("[CALLBACK] from=%s (id=%d) data=%q",
				cb.Sender.Username, cb.Sender.ID, cb.Data)
		}
		return next(c)
	}
}
