package ollama

import "sync"

type Context struct {
	ollamas map[string]*Ollama
	mu      *sync.Mutex
}

func NewContext() Context {
	return Context{
		ollamas: make(map[string]*Ollama),
		mu:      &sync.Mutex{},
	}
}

func (c *Context) AddOllama(name string, prompt string, model string) bool {
	if _, ok := c.ollamas[name]; ok {
		return false
	}

	c.mu.Lock()
	c.ollamas[name] = NewOllama(name, prompt, model)
	c.mu.Unlock()

	return true
}

func (c *Context) RemoveOllama(name string) bool {
	if _, ok := c.ollamas[name]; !ok {
		return false
	}

	c.mu.Lock()
	delete(c.ollamas, name)
	c.mu.Unlock()

	return true
}

func (c *Context) Ollama(name string) (*Ollama, bool) {
	ollama, ok := c.ollamas[name]
	return ollama, ok
}

func (c *Context) GetCompletion(prompt string, model string) (string, error) {
	resp, err := GetCompletion(prompt, model)
	if err != nil {
		return "", err
	}
	return resp.Response, nil
}
