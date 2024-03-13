package ollama

import (
	"errors"
	"sync"
)

type Ollama struct {
	name    string
	model   string
	history []Message
	mu      *sync.Mutex
}

func NewOllama(name string, prompt string, model string) Ollama {
	ollama := Ollama{
		name:    name,
		model:   model,
		history: make([]Message, 0),
		mu:      &sync.Mutex{},
	}

	ollama.history = append(ollama.history, Message{
		Role:    "system",
		Content: prompt,
	})

	return ollama
}

func (o *Ollama) GetChatCompletion(prompt string) (string, error) {
	o.historyMessage("user", prompt)

	completion, err := GetChatCompletion(o.history, o.model)
	if err != nil {
		return "", err
	}

	if len(completion.Choices) == 0 {
		return "", errors.New("no completions")
	}

	message := completion.Choices[0].Message
	o.historyMessage("assistant", message.Content)

	return message.Content, nil
}

func (o *Ollama) Clear() {
	o.mu.Lock()
	history := make([]Message, 0)
	o.history = append(history, o.history[0])
	o.mu.Unlock()
}

func (o *Ollama) historyMessage(role string, prompt string) {
	o.mu.Lock()
	o.history = append(o.history, Message{
		Role:    role,
		Content: prompt,
	})
	o.mu.Unlock()
}
