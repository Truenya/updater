package cli

import (
	"errors"
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/cli/tui"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/config"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/statistic"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/updater"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/util"
)

func UpdateAction(cCtx *cli.Context) error {
	serviceName, uType, err := ParseArgs(cCtx)
	if err != nil {
		// It is not an error, just user choice
		if errors.Is(err, NoServiceError{}) {
			return nil
		}

		return err
	}

	u, err := GetUpdater(serviceName, uType)
	if err != nil {
		return err
	}

	u.SetBusy(true)

	p := Update(u, &err)

	WaitForUpdateAction(cCtx, p, GetMaxTime(serviceName, uType))

	return err // got in the goroutine
}

func Update(u *updater.Updater, err *error) chan util.Progress {
	p := make(chan util.Progress)
	go func(p chan util.Progress) {
		*err = u.Update(p)
		close(p)
	}(p)

	return p
}

func GetMaxTime(serviceName string, uType util.UpdaterType) float64 { // progressbar is working with floats
	maxTime := statistic.GetSumDuration(serviceName, uType)
	if maxTime == 0 {
		maxTime = bigNumber
	}

	return float64(maxTime.Milliseconds())
}

func WaitForUpdateAction(cCtx *cli.Context, p chan util.Progress, maxTime float64) {
	if cCtx.Bool("quiet") {
		// just wait
		for {
			_, ok := <-p
			if !ok {
				return
			}
		}
	}

	// tui is receiver of progress
	tui.RunSingleProgress(p, maxTime)
}

type NoServiceError struct{}

func (e NoServiceError) Error() string {
	return "" // it is not exactly an error, just user choice
}

type NotSupportedError struct {
	serviceName string
}

func (e NotSupportedError) Error() string {
	return fmt.Sprintf("service %s is not supported. Please add it in config.json and write scripts to it.", e.serviceName)
}

func ParseArgs(cCtx *cli.Context) (string, util.UpdaterType, error) {
	serviceName := GetServiceName(cCtx.Args())
	if serviceName == "" {
		return "", util.UpdateLocal, NoServiceError{}
	}

	if _, ok := config.GetSupported()[serviceName]; !ok {
		return "", util.UpdateLocal, NotSupportedError{serviceName: serviceName}
	}

	phrase := "Updating service "
	uType := util.UpdateLocal

	if cCtx.Bool("container") {
		uType = util.UpdateRemote
		phrase += "using container"
	} else if cCtx.Bool("inplace") {
		uType = util.UpdateInPlace
		phrase += "in place"
	}

	if !cCtx.Bool("quiet") {
		fmt.Println(phrase+":", ColorService(serviceName))
	}

	return serviceName, uType, SaveArgs(cCtx, serviceName, uType)
}

func GetServiceName(args cli.Args) string {
	if args.Present() {
		return args.First()
	}

	return tui.ChooseList(config.GetServicesFiltered(""))
}

func SaveArgs(cCtx *cli.Context, serviceName string, uType util.UpdaterType) error {
	switch {
	case cCtx.NArg() < 2:
		if _, ok := config.Get("ssh", "addr"); !ok {
			return fmt.Errorf("no ip specified or configured")
		}

		return nil
	case cCtx.NArg() >= 2 && UpdateIP(cCtx): // nothing
	case cCtx.NArg() >= 3 && uType == util.UpdateRemote:
		config.Set(serviceName, "branch", cCtx.Args().Get(2))
	}

	return config.Write()
}

func UpdateIP(cCtx *cli.Context) bool {
	ip := cCtx.Args().Get(1)
	config.Set("ssh", "addr", ip)
	config.AddSelect("addr", ip)

	return false
}

func GetUpdater(serviceName string, uType util.UpdaterType) (*updater.Updater, error) {
	dataForService := config.Args(serviceName)

	return updater.Get(serviceName, uType, dataForService, config.Args("ssh")), nil
}

const bigNumber = 10000 * time.Second
const shortMsg = 5

func UpdateCommand() *cli.Command {
	return &cli.Command{
		Name:      "update",
		HelpName:  I("u[pdate]"),
		Aliases:   []string{"u"},
		Usage:     C("Update chosen service with latest used or specified ip"),
		UsageText: updateUsage(),
		Action:    UpdateAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "container", Aliases: []string{"c"}},
			&cli.BoolFlag{Name: "inplace", Aliases: []string{"i"}},
			&cli.BoolFlag{Name: "quiet", Aliases: []string{"q"}},
		},
		UseShortOptionHandling: true,
		BashComplete:           updBashComplete,
	}
}

func updBashComplete(cCtx *cli.Context) {
	// switch is sequential by its nature, so we can use it to compute all and stop, when needed
	switch {
	// Len 0 when showing all and 1 when it is not completed, also 1 when it is completed.
	// So will check completness and return if it is not completed
	case cCtx.NArg() < 2 && serviceNameBashComplete(cCtx):
		return
	case cCtx.NArg() < 3 && addrBashComplete(cCtx):
		return
	case cCtx.Bool("container") && cCtx.NArg() < 4 && branchBashComplete(cCtx):
	}
}

// When we already have service name
// We dont want to show all services
// Instead we are going ahead to addresses.
func serviceNameBashComplete(cCtx *cli.Context) bool {
	// TODO: filter by supported
	// Possible problem: if we have full service name as part of another - completes will show shit until first letter of
	// next arg will be written.
	// But we haven't, as for now
	cur := cCtx.Args().First()

	services := config.GetServicesFiltered(cur)
	if len(services) == 0 {
		return false
	}

	for _, service := range services {
		fmt.Fprintln(cCtx.App.Writer, service)
	}

	return true
}

func addrBashComplete(cCtx *cli.Context) bool {
	const addrArgPos = 1

	return selectsBashComplete(cCtx, "addr", addrArgPos)
}

func branchBashComplete(cCtx *cli.Context) bool {
	const branchArgPos = 2

	serviceName := cCtx.Args().First()

	return selectsBashComplete(cCtx, serviceName+"_branch", branchArgPos)
}

func selectsBashComplete(cCtx *cli.Context, option string, pos int) bool {
	cur := cCtx.Args().Get(pos)

	selects, ok := config.GetSelectsFiltered(option, cur)
	if !ok || len(selects) == 0 || len(selects) == 1 && cur == selects[0] {
		// no data in selects or only one option and it is the same as current one
		// so we dont need to show anything and go ahead
		return false
	}

	for _, comp := range selects {
		fmt.Fprintln(cCtx.App.Writer, comp)
	}

	// when some options are showed we need to return in upstream switch
	return true
}

func updateUsage() string {
	return "" +
		I("./upd u back") + "\t— " +
		C("Updating back from local build, using last used ip") + "\n" +
		I("./upd u -q back") + "\t— " +
		C("Updating back from local build, using last used ip in quiet mode") + "\n" +
		I("./upd u back 172.31.48.84") + "\t— " +
		C("Updating back from local build, using specified ip") + "\n" +
		I("./upd u back voronezh") + "\t— " +
		C("Updating back from local build, using specified ssh alias, known or resolvable hostname") + "\n" +
		I("./upd u -c back") + "\t— " +
		C("Updating back from docker container usyng last used ip and last used branch") + "\n" +
		I("./upd u -cq back 172.31.48.84") + "\t— " +
		C("Updating using specified ip and last used branch in quiet mode") + "\n" +
		I("./upd u -c back 172.31.48.84 rc-import-fixes") + "\t— " +
		C("Updating back from docker container on specified ip and branch") + "\n" +
		I("./upd u -i back") + "\t— " +
		C("Updating back in place")
}
