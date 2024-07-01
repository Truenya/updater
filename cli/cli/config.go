package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/cli/tui"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/config"
)

func ConfigCommand() *cli.Command {
	return &cli.Command{
		Name:                   "config",
		HelpName:               I("c[onfig]"),
		Description:            configDescription(),
		Aliases:                []string{"c"},
		Usage:                  C("Allow to change config.json"),
		UsageText:              configUsage(),
		Action:                 configAction,
		UseShortOptionHandling: true,
		Subcommands: []*cli.Command{
			{
				Name:      "set",
				HelpName:  I("set"),
				Aliases:   []string{"s"},
				Usage:     C("Set key to value for service"),
				UsageText: setUsage(),
				Action:    setAction,
			},
			{
				Name:      "unset",
				HelpName:  I("unset"),
				Aliases:   []string{"u"},
				Usage:     C("Uhset key for service"),
				UsageText: unsetUsage(),
				Action:    unsetAction,
			},
		},
	}
}

func configDescription() string {
	return D("To use local build, directory for service should be set") + "\n" +
		D("Build dir also can be specified. Default is 'build'")
}

func configAction(cCtx *cli.Context) error {
	if !cCtx.Args().Present() {
		// Show all settings, and allow to change them
		return tui.ConfigSSH()
	}

	return nil
}

func setAction(cCtx *cli.Context) error {
	if cCtx.Args().Len() != 3 {
		// TODO: tui config all services
		// upd c s -- show list of services
		// upd c s <service> -- show list of keys
		return fmt.Errorf("wrong number of arguments")
	}

	config.Set(cCtx.Args().First(), cCtx.Args().Get(1), cCtx.Args().Get(2))

	return config.Write()
}

func unsetAction(cCtx *cli.Context) error {
	if cCtx.Args().Len() != 2 {
		// TODO: tui config all services
		// upd c s -- show list of services
		// upd c s <service> -- show list of keys
		return fmt.Errorf("wrong number of arguments")
	}

	config.Unset(cCtx.Args().First(), cCtx.Args().Get(1))

	return config.Write()
}

func unsetUsage() string {
	return I("./upd c u[nset] <service> <key>") + "\t— " +
		C("Unsetting <key> for <service>")
}

func setUsage() string {
	return I("./upd c s[et] <service> <key> <value>") + "\t— " +
		C("Setting <key> to <value> for <service>")
}

func configUsage() string {
	return "" +
		I("./upd c") + "\t— " + C("Configuring ssh connections interactively") + "\n" +
		I("./upd c s back dir /home/user/vm_back") + "\t— " + C("Setting directory, from which back will be taken") + "\n" +
		I("./upd c u back build cmake_build_debug") + "\t— " + C("Setting build dir") + "\n" +
		I("./upd c u back build") + "\t— " + C("Unsetting build dir") + "\n" +
		I("./upd c s back branch rc-import-fixes") + "\t— " + C("Setting branch")
}
