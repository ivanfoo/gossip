package bot

import (
	_ "github.com/ivanfoo/gossip/utils"
)

type Command struct {
	Action string
	Target string
}

func NewCommand(action string, target string) *Command {
	c := new(Command)
	c.Action = action
	c.Target = target

	return c
}