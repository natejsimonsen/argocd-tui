package model

import (
	"fmt"
)

type Context string

const (
	App      = "App"
	Global   = "Global"
	AppList  = "AppList"
	MainPage = "MainPage"
	Help     = "Help"
)

type Command struct {
	Description string
	Handler     func()
	Context     Context
}

func (c *Command) String() string {
	return fmt.Sprintf("%-10s - %s", c.Context, c.Description)
}

type CommandModel struct {
	Commands map[Context]map[rune]*Command
}

func NewCommandModel() *CommandModel {
	commands := map[Context]map[rune]*Command{}

	commands[App] = map[rune]*Command{}
	commands[Global] = map[rune]*Command{}
	commands[AppList] = map[rune]*Command{}
	commands[MainPage] = map[rune]*Command{}
	commands[Help] = map[rune]*Command{}

	return &CommandModel{
		Commands: commands,
	}
}

func (m *CommandModel) Add(r rune, context Context, desc string, handler func()) error {
	if cmd, ok := m.Commands[context][r]; ok {
		return fmt.Errorf("Error: command already exists, %s", cmd)
	}

	cmd := Command{
		Context:     context,
		Description: desc,
		Handler:     handler,
	}

	m.Commands[context][r] = &cmd
	return nil
}
