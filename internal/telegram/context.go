package telegram

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	ErrArgumentNotFound = errors.New("argument not found")
)

type Context struct {
	api     *tgbotapi.BotAPI
	message *tgbotapi.Message
	args    map[string]string
}

type Argument interface {
	Parse(value string) error
}

func (c Context) Argument(name string, value Argument) error {
	argValue, ok := c.args[name]
	if !ok {
		return ErrArgumentNotFound
	}
	return value.Parse(argValue)
}

func (c Context) Message() *tgbotapi.Message {
	return c.message
}

func (c Context) Reply(text string) error {
	message := tgbotapi.NewMessage(c.Message().Chat.ID, text)
	message.ReplyToMessageID = c.Message().MessageID

	_, err := c.api.Send(message)
	return err
}
