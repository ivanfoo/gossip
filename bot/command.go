package bot

import (
	_ "github.com/ivanfoo/rtop-bot/utils"
)

type CommandOptions struct {
	Action string
	Target string	
}

type Command struct {
	commandOptions CommandOptions
}

func NewCommand(opts CommandOptions) *Command {
	c := new(Command)
	c.commandOptions = opts

	return c
}