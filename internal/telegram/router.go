package telegram

import (
	"regexp"
	"strings"
)

var argumentRegex = regexp.MustCompile(`\{(\w+)\}`)
var argumentRestRegex = regexp.MustCompile(`\{\.\.\.(\w+)\}`)

type Router struct {
	commands []Command
}

type Command struct {
	pattern *regexp.Regexp
	handler CommandHandler
}

type CommandHandler = func(ctx *Context) error

func NewRouter() Router {
	return Router{
		commands: make([]Command, 0),
	}
}

func (r *Router) Command(pattern string, handler CommandHandler) {
	pattern = strings.ToLower(pattern)
	r.addCommand(pattern, handler)
}

func (r *Router) Execute(ctx *Context) error {
	userCommand := strings.TrimPrefix(ctx.Message().Text, "/")

	var command Command
	for _, command = range r.commands {
		if !command.pattern.MatchString(userCommand) {
			continue
		}

		groupNames := command.pattern.SubexpNames()
		for _, match := range command.pattern.FindAllStringSubmatch(userCommand, -1) {
			for groupIdx, group := range match {
				name := groupNames[groupIdx]
				if name == "" {
					continue
				}
				ctx.args[name] = group
			}
		}

		break
	}

	return command.handler(ctx)
}

func (r *Router) addCommand(pattern string, handler CommandHandler) {
	r.commands = append(r.commands, Command{
		pattern: r.compilePattern(pattern),
		handler: handler,
	})
}

func (r *Router) compilePattern(pattern string) *regexp.Regexp {
	pattern = argumentRegex.ReplaceAllString(pattern, `(?P<$1>\S+)`)
	pattern = argumentRestRegex.ReplaceAllString(pattern, `(?P<$1>.*)`)
	return regexp.MustCompile(pattern)
}
