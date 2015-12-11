package monitors

import (
	"github.com/ivanfoo/gossip/utils"
	"golang.org/x/crypto/ssh"
)

type Containers struct {
    Report string
    Client *ssh.Client
}

func NewContainersMonitor(client *ssh.Client) *Containers {
    c := new(Containers)
    c.Client = client
    return c
}

func (c *Containers) getContainers() (string, error) {
	report, err := utils.RunCommand(c.Client, "docker ps")
	return report, err
}

func (c *Containers) Run() string {
	output, _ := c.getContainers()
	return output
}

func RunContainersMonitor(client *ssh.Client) string {
	report, _ := utils.RunCommand(client, "docker ps | awk '{print $2}' | tail -n +2")
	return report
}
