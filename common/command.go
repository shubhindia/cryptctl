package common

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var commands []*cli.Command
var versionString string

type Commander interface {
	Execute(c *cli.Context) error
}

func Version() string {
	return versionString
}
func RegisterCommand(command cli.Command) {
	logrus.Debugln("Registering", command.Name, "command...")
	commands = append(commands, &command)
}

func RegisterCommand2(name, usage string, data Commander, flags ...cli.Flag) {
	RegisterCommand(cli.Command{
		Name:   name,
		Usage:  usage,
		Action: data.Execute,
	})
}

func GetCommands() []*cli.Command {
	return commands
}
