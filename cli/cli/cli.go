package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Run(args []string) {
	app := &cli.App{
		Name:                 I("upd"),
		HelpName:             I("upd"),
		Usage:                C("Command line application to update services"),
		Description:          Description(),
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			UpdateCommand(),
			ConfigCommand(),
		},
		UseShortOptionHandling: true,
	}

	if err := app.Run(args); err != nil {
		fmt.Fprintln(app.Writer, colored(err.Error(), colorError))
		panic(err)
	}
}

func Description() string {
	return D("Before using upd, you need to configure ssh connection via ") +
		I("./upd c") + "\n" +
		D("And specify directory for service if using local build") +
		"\n" + D("To get more precise help, type ") + "\n" + I("./upd c -h") +
		"\n" + D(" or ") + "\n" + I("./upd u -h")
}
