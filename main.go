package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/shubhindia/hcictl/commands"
	"github.com/shubhindia/hcictl/common"
)

func main() {

	app := cli.NewApp()
	app.Name = "hcictl"
	app.Usage = "hcictl is a command line tool"
	app.Version = common.Version()
	app.Authors = []*cli.Author{
		{
			Name:  "Shubham Gopale",
			Email: "shubhindia123@gmail.com",
		},
	}
	app.Commands = common.GetCommands()
	app.CommandNotFound = func(context *cli.Context, command string) {
		logrus.Fatalln("Command", command, "not found.")
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}