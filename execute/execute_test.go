package execute

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/defaults"
)

func TestLocal(t *testing.T) {
	t.Parallel()
	t.Run("local_err", func(t *testing.T) {
		t.Parallel()
		_, err := Local("Not existing command")
		require.Error(t, err)
	})
	t.Run("local_write_to_file", func(t *testing.T) {
		t.Parallel()
		_, err := Local("touch test.txt")
		require.NoError(t, err)
		_, err = Local("echo test > test.txt")
		require.NoError(t, err)
		out, err := Local("cat test.txt")
		require.NoError(t, err)
		require.Equal(t, "test\n", string(out))
	})
}

func TestRemoterAddKey(t *testing.T) {
	t.Parallel()

	r := Remoter{}
	l := r.AddKey([]string{})

	require.Equal(t, []string{}, l)

	r.Key = "key"

	l = r.AddKey([]string{})
	require.Equal(t, []string{"-i", "key"}, l)
}

func TestRemoterAddPassword(t *testing.T) {
	t.Parallel()

	r := Remoter{}
	args, needPass := r.AddPassword([]string{})

	require.False(t, needPass)
	require.Equal(t, []string{}, args)

	r.Pass = "pass"

	args, needPass = r.AddPassword([]string{})
	require.True(t, needPass)
	require.Equal(t, []string{"-p", "pass", "ssh"}, args)
}

func TestRemoterGetUserAndAddr(t *testing.T) {
	t.Parallel()

	r := Remoter{}
	addr, err := r.GetUserAndAddr()

	require.Error(t, err)
	require.Equal(t, "", addr)

	r.Addr = "addr"

	addr, err = r.GetUserAndAddr()
	require.NoError(t, err)
	require.Equal(t, "addr", addr)

	r.User = "user"
	addr, err = r.GetUserAndAddr()
	require.NoError(t, err)
	require.Equal(t, "user@addr", addr)
}

func TestGetRsyncSsh(t *testing.T) {
	t.Parallel()

	r := Remoter{}
	require.Equal(t, "", r.GetRsyncSSH())
	r.Key = "key"
	require.Equal(t, `-e "ssh -i key"`, r.GetRsyncSSH())
	r.Key = ""
	r.Pass = "pass"
	require.Equal(t, `-e "sshpass -p pass ssh"`, r.GetRsyncSSH())
}

func TestRsyncCommand(t *testing.T) {
	t.Parallel()

	r := Remoter{Addr: "addr"}
	localDir := "localDir"
	file := "file"
	serviceName := "serviceName"
	rsyncPart := fmt.Sprintf("rsync %s -r ", r.GetRsyncSSH())
	localFilePath := fmt.Sprintf(" %s/%s ", localDir, file)
	addr, err := r.GetUserAndAddr()
	require.NoError(t, err)

	remoteFilePath := fmt.Sprintf(" %s:%s%s/%s", addr, defaults.DefaultRemoteDir, serviceName, file)
	v, err := r.RsyncCommand(localDir, file, serviceName)
	require.NoError(t, err)
	require.Equal(t, rsyncPart+localFilePath+remoteFilePath, v)
}

func TestSSHCommand(t *testing.T) {
	t.Parallel()

	r := Remoter{Addr: "addr"}
	l, needSshpass, err := r.SSHCommand()
	gold := AddStrictHostOpts([]string{})
	gold = append(gold, "addr")

	require.NoError(t, err)
	require.False(t, needSshpass)
	require.Equal(t, gold, l)
}

func TestWithSshPass(t *testing.T) {
	t.Parallel()

	a := WithSSHPass([]string{})
	require.Contains(t, a.Path, "sshpass")
	require.Equal(t, a.Args, []string{"sshpass"})
}

func TestJustSsh(t *testing.T) {
	t.Parallel()

	a := JustSSH([]string{})
	require.Contains(t, a.Path, "ssh")
	require.Equal(t, a.Args, []string{"ssh", "-o", "BatchMode=yes"})
}

func TestGetExecCmd(t *testing.T) {
	t.Parallel()

	a := GetExecCmd([]string{}, false)
	require.Contains(t, a.Path, "ssh")
	require.Equal(t, a.Args, []string{"ssh", "-o", "BatchMode=yes"})
	a = GetExecCmd([]string{}, true)
	require.Contains(t, a.Path, "sshpass")
	require.Equal(t, a.Args, []string{"sshpass"})
}

func TestCheckConnection(t *testing.T) {
	t.Parallel()

	a := CheckConnection("addr")
	require.Equal(t, a, true)
}

// idk how to replicate not ExitError type here.
// func TestBySSH(t *testing.T) {
// 	t.Parallel()
// 	// TODO need to configure in pipelines first
// 	// a, err := BySSH(Remoter{Addr: "localhost"}, "echo test")
// 	// require.NoError(t, err)
// 	// require.Equal(t, "test\n", a)
// 	a, err := BySSH(Remoter{Addr: "localhost"}, "notexistingcommand")
// 	require.Error(t, err)
// 	require.Equal(t, "", a)
// }

// func TestRsync(_ *testing.T) {
// TODO need to configure in pipelines first
// r := Remoter{Addr: "localhost"}
// localDir := "./"
// file := "file"
// serviceName := "serviceName"
// err := Rsync(r, localDir, file, serviceName)
// require.Error(t, err)
// _, err = BySSH(Remoter{Addr: "localhost"}, "echo test > file")
// require.NoError(t, err)
// err = Rsync(r, localDir, file, serviceName)
// require.Error(t, err)
// _, err = BySSH(Remoter{Addr: "localhost"}, "rm file")
// require.NoError(t, err)
// }

// func TestRemoterFromConfig(t *testing.T) {
// 	t.Parallel()
//
// 	r := RemoterFromConfig()
// 	goldkey, _ := config.Get("ssh", "key")
// 	goldpass, _ := config.Get("ssh", "pass")
// 	goldaddr, _ := config.Get("ssh", "addr")
// 	golduser, _ := config.Get("ssh", "user")
//
// 	require.Equal(t, goldkey, r.Key)
// 	require.Equal(t, goldpass, r.Pass)
// 	require.Equal(t, goldaddr, r.Addr)
// 	require.Equal(t, golduser, r.User)
// }
