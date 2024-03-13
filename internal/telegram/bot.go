package telegram

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api    *tgbotapi.BotAPI
	router Router
}

func NewBot(token string, router Router) (Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return Bot{}, err
	}

	return Bot{
		api:    api,
		router: router,
	}, nil
}

func (b Bot) Listen() error {
	updateCfg := tgbotapi.NewUpdate(0)
	updateChan := b.api.GetUpdatesChan(updateCfg)

	for update := range updateChan {
		if update.Message == nil {
			continue
		}

		if !strings.HasPrefix(update.Message.Text, "/") {
			continue
		}

		go b.handleCommand(update.Message)
	}

	return nil
}

func (b Bot) handleCommand(message *tgbotapi.Message) {
	ctx := &Context{
		api:     b.api,
		message: message,
		args:    make(map[string]string),
	}

	if err := b.router.Execute(ctx); err != nil {
		message := tgbotapi.NewMessage(message.Chat.ID, err.Error())
		b.api.Send(message)
	}
}
