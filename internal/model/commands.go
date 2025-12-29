package model

import (
	"fmt"
)

type Command struct {
	Description string
	Handler     func()
}

func (c *Command) String() string {
	return fmt.Sprintf("%s", c.Description)
}

type CommandModel struct {
	Commands map[rune]*Command
}

func NewCommandModel() *CommandModel {
	return &CommandModel{
		Commands: map[rune]*Command{},
	}
}

func (m *CommandModel) Add(r rune, command Command) error {
	if cmd, ok := m.Commands[r]; ok {
		return fmt.Errorf("Error: command already exists, %s", cmd)
	}

	m.Commands[r] = &command
	return nil
}
