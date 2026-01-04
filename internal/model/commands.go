package model

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type Context string

const (
	App        = "App"
	Global     = "Global"
	CommandBar = "CommandBar"
	AppTable   = "AppTable"
	MainPage   = "MainPage"
	Help       = "Help"
)

type Command struct {
	Description string
	Handler     func()
	Context     Context
}

func (c *Command) String() string {
	return fmt.Sprintf("%-10s - %s", c.Context, c.Description)
}

type KeyStroke struct {
	Rune rune
	Key  tcell.Key
}

type CommandModel struct {
	Commands map[Context]map[KeyStroke]*Command
	Context  Context
}

func NewCommandModel() *CommandModel {
	commands := map[Context]map[KeyStroke]*Command{}

	commands[App] = map[KeyStroke]*Command{}
	commands[Global] = map[KeyStroke]*Command{}
	commands[CommandBar] = map[KeyStroke]*Command{}
	commands[AppTable] = map[KeyStroke]*Command{}
	commands[MainPage] = map[KeyStroke]*Command{}
	commands[Help] = map[KeyStroke]*Command{}

	return &CommandModel{
		Commands: commands,
		Context:  Global,
	}
}

func (m *CommandModel) Add(ks KeyStroke, context Context, desc string, handler func(ctx Context)) error {
	if cmd, ok := m.Commands[context][ks]; ok {
		return fmt.Errorf("error: command already exists, %s", cmd)
	}

	cmd := Command{
		Context:     context,
		Description: desc,
		Handler: func() {
			handler(context)
			if context != Global {
				m.Context = context
			}
		},
	}

	m.Commands[context][ks] = &cmd
	return nil
}
