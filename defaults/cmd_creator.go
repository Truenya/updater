package defaults

type CmdCreator struct {
	s string // service
}

func (c CmdCreator) MkBinDir() string {
	return CreateBinDirCmd(c.s)
}

func (c CmdCreator) MkDir(dir string) string {
	return CreateDirCmd(c.s + "/" + dir)
}

func (c CmdCreator) RmDir() string {
	return RemoveExistingDirCmd(c.s)
}

func (c CmdCreator) PullImg() string {
	return PullImageCmd(c.s)
}

func (c CmdCreator) CreateContainer() string {
	return CreateContainerCmd(c.s)
}

func (c CmdCreator) RsyncBin() []string {
	return RsyncArg(c.s)
}

func (c CmdCreator) DCpBinToBox() string {
	return DockerCpBinToBoxDefaultCmd(c.s)
}

func (c CmdCreator) DCpBinTo(container string) string {
	return DockerCpBinToDefaultCmd(c.s, container)
}

func (c CmdCreator) DCpToBoxAll() string {
	return c.DCpToDefault("", "vm_box", "")
}

func (c CmdCreator) DCpToDefaultBox(what string) string {
	return DockerCpToDefaultBoxCmd(c.s, what)
}

func (c CmdCreator) DCpToDefaultBoxDest(what, to string) string {
	return DockerCpToDefaultBoxDestCmd(c.s, what, to)
}

func (c CmdCreator) DCpToSameDirAsService(what, container string) string {
	return DockerCpToDefaultCmd(c.s, what, container)
}

func (c CmdCreator) DCpToDefault(what, container, to string) string {
	return DockerCpToDefaultDestCmd(c.s, what, container, to)
}

func (c CmdCreator) DRestartBoxService() string {
	return DockerRestartBoxServiceCmd(c.s)
}

func (c CmdCreator) DRestartService(container string) string {
	return DockerRestartServiceCmd(c.s, container)
}

func (c CmdCreator) DCpBinFrom(from string) string {
	return DockerCpBinFromCmd(c.s, from)
}

func (c CmdCreator) DCpFrom(what string) string {
	return DockerCpFromToCmd(what, c.s)
}

func (c CmdCreator) DCpAllFrom(name string) string {
	return DockerCpServiceFromCmd(c.s, name)
}

func (c CmdCreator) DCpToBoxDest(from, to string) string {
	return DockerCpToBoxDestCmd(c.s, from, to)
}

func (c CmdCreator) MvInService(from, to string) string {
	return MoveServiceDefDir(c.s, from, to)
}
