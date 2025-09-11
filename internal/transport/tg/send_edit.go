package tg

import (
	"strings"

	"gopkg.in/telebot.v3"
)

const zeroWidthJoiner = "\u200d"
const errMsgNotModified = "message is not modified"

type sendEditCfg struct {
	Edit                bool              // true => редактируем текущее сообщение
	EditMarkupIfSame    bool              // при "not modified" попробовать обновить только клавиатуру
	ForceChangeOnSame   bool              // при "not modified" добавить невидимый символ к тексту и повторить
	ReplyTo             *telebot.Message  // ответ на сообщение
	ParseMode           telebot.ParseMode // telebot.ModeMarkdown / HTML и т.п.
	DisableNotification bool              // отправка без звука
	ProtectContent      bool              // запрет пересылки/сохранения
	AllowWithoutReply   bool              // разрешить отправку без ReplyTo
}

type SendEditOption func(*sendEditCfg)

// Включить режим редактирования
func WithEdit() SendEditOption { return func(c *sendEditCfg) { c.Edit = true } }

// Если контент тот же — обновить только клавиатуру
func WithEditMarkupIfSame() SendEditOption {
	return func(c *sendEditCfg) { c.EditMarkupIfSame = true }
}

// Если контент тот же — форсировать изменение текста (добавить \u200d)
func WithForceChangeOnSame() SendEditOption {
	return func(c *sendEditCfg) { c.ForceChangeOnSame = true }
}

// Установить ParseMode (Markdown/HTML/Plain)
func WithParseMode(m telebot.ParseMode) SendEditOption {
	return func(c *sendEditCfg) { c.ParseMode = m }
}

// Ответ на конкретное сообщение
func WithReplyTo(msg *telebot.Message) SendEditOption {
	return func(c *sendEditCfg) { c.ReplyTo = msg }
}

// Отправить без звука
func WithSilent() SendEditOption { return func(c *sendEditCfg) { c.DisableNotification = true } }

// Защитить контент (нельзя переслать/сохранить)
func WithProtectContent() SendEditOption { return func(c *sendEditCfg) { c.ProtectContent = true } }

// Разрешить отправку без ReplyTo
func WithAllowWithoutReply() SendEditOption {
	return func(c *sendEditCfg) { c.AllowWithoutReply = true }
}

// SendOrEdit отправляет или редактирует сообщение.
// text: текст сообщения (если пустой и edit=true, а markup != nil — можно обновить только клавиатуру)
func (bot *Bot) SendOrEdit(
    c telebot.Context,
    text string,
    markup *telebot.ReplyMarkup,
    opts ...SendEditOption,
) (*telebot.Message, error) {
    cfg := sendEditCfg{ParseMode: telebot.ModeMarkdown}
    for _, opt := range opts {
        opt(&cfg)
    }

    sendOpts := &telebot.SendOptions{
        ParseMode:           cfg.ParseMode,
        ReplyMarkup:         markup,
        ReplyTo:             cfg.ReplyTo,
        DisableNotification: cfg.DisableNotification,
        Protected:           cfg.ProtectContent,
        AllowWithoutReply:   cfg.AllowWithoutReply,
    }

    if cfg.Edit && text == "" && markup != nil {
        return bot.B.EditReplyMarkup(c.Message(), markup)
    }

    if cfg.Edit {
        m, err := bot.B.Edit(c.Message(), text, sendOpts)
        if err != nil {
            if strings.Contains(err.Error(), errMsgNotModified) {
                if cfg.EditMarkupIfSame && markup != nil {
                    return bot.B.EditReplyMarkup(c.Message(), markup)
                }
                if cfg.ForceChangeOnSame {
                    return bot.B.Edit(c.Message(), text+zeroWidthJoiner, sendOpts)
                }
                return nil, nil
            }
            return nil, err
        }
        return m, nil
    }

    return bot.B.Send(c.Recipient(), text, sendOpts)
}
