package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/alexthvest/ollama-bot/internal/ollama"
	"github.com/alexthvest/ollama-bot/internal/telegram"
	"github.com/joho/godotenv"
)

var supportedModels = []string{"openhermes", "mistral", "wizard-vicuna-uncensored:13b-q4_0", "codellama:7b", "llama2:13b", "llama2-uncensored:7b"}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	token, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Fatal("BOT_TOKEN is missing")
	}

	router := setupRouter()
	bot, err := telegram.NewBot(token, router)
	if err != nil {
		log.Fatal(err)
	}

	if err := bot.Listen(); err != nil {
		log.Fatal(err)
	}
}

func setupRouter() telegram.Router {
	router := telegram.NewRouter()
	ollamaCtx := ollama.NewContext()

	router.Command("ollama:{name} {...prompt}", func(ctx *telegram.Context) error {
		return ollamaNamed(ctx, ollamaCtx)
	})

	router.Command("ollama create {model} {name} {...prompt}", func(ctx *telegram.Context) error {
		return ollamaCreate(ctx, ollamaCtx)
	})

	router.Command("ollama clear {name}", func(ctx *telegram.Context) error {
		return ollamaClear(ctx, ollamaCtx)
	})

	router.Command("ollama rm {name}", func(ctx *telegram.Context) error {
		return ollamaRemove(ctx, ollamaCtx)
	})

	router.Command("ollama {model} {...prompt}", func(ctx *telegram.Context) error {
		return ollamaDefault(ctx, ollamaCtx)
	})

	return router
}

func ollamaDefault(ctx *telegram.Context, ollamaCtx ollama.Context) error {
	var modelArg telegram.String
	if err := ctx.Argument("model", &modelArg); err != nil {
		return err
	}

	var promptArg telegram.String
	if err := ctx.Argument("prompt", &promptArg); err != nil {
		return err
	}

	model := string(modelArg)
	prompt := string(promptArg)

	if !slices.Contains(supportedModels, model) {
		return fmt.Errorf("unsupported model. supported models: %s", strings.Join(supportedModels, ", "))
	}

	completion, err := ollamaCtx.GetCompletion(prompt, model)
	if err != nil {
		return err
	}

	return ctx.Reply(completion)
}

func ollamaNamed(ctx *telegram.Context, ollamaCtx ollama.Context) error {
	var nameArg telegram.String
	if err := ctx.Argument("name", &nameArg); err != nil {
		return err
	}

	var promptArg telegram.String
	if err := ctx.Argument("prompt", &promptArg); err != nil {
		return err
	}

	name := string(nameArg)
	prompt := string(promptArg)

	ollama, ok := ollamaCtx.Ollama(name)
	if !ok {
		return fmt.Errorf("ollama '%s' not exists", name)
	}

	completion, err := ollama.GetChatCompletion(prompt)
	if err != nil {
		return err
	}

	return ctx.Reply(completion)
}

func ollamaCreate(ctx *telegram.Context, ollamaCtx ollama.Context) error {
	var modelArg telegram.String
	if err := ctx.Argument("model", &modelArg); err != nil {
		return err
	}

	var nameArg telegram.String
	if err := ctx.Argument("name", &nameArg); err != nil {
		return err
	}

	var promptArg telegram.String
	if err := ctx.Argument("prompt", &promptArg); err != nil {
		return err
	}

	model := string(modelArg)
	name := string(nameArg)
	prompt := string(promptArg)

	if !slices.Contains(supportedModels, model) {
		return fmt.Errorf("unsupported model. supported models: %s", strings.Join(supportedModels, ", "))
	}

	if ok := ollamaCtx.AddOllama(name, prompt, model); !ok {
		return errors.New("ollama with this name already exists")
	}

	return ctx.Reply(fmt.Sprintf("ollama '%s' created", name))
}

func ollamaClear(ctx *telegram.Context, ollamaCtx ollama.Context) error {
	var nameArg telegram.String
	if err := ctx.Argument("name", &nameArg); err != nil {
		return err
	}

	name := string(nameArg)

	ollama, ok := ollamaCtx.Ollama(name)
	if !ok {
		return fmt.Errorf("ollama '%s' not exists", name)
	}

	ollama.Clear()
	return ctx.Reply(fmt.Sprintf("ollama '%s' restored to initial state", name))
}

func ollamaRemove(ctx *telegram.Context, ollamaCtx ollama.Context) error {
	var nameArg telegram.String
	if err := ctx.Argument("name", &nameArg); err != nil {
		return err
	}

	name := string(nameArg)

	if ok := ollamaCtx.RemoveOllama(name); !ok {
		return fmt.Errorf("ollama '%s' not exists", name)
	}

	return ctx.Reply(fmt.Sprintf("ollama '%s' removed", name))
}
