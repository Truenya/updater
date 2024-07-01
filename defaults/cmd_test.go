package defaults

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var b, r, m, c = CmdCreator{"back"}, CmdCreator{"resowatch"}, //nolint:gochecknoglobals
	CmdCreator{"monitor"}, CmdCreator{"checker"}

func TestInheritSprintf(t *testing.T) {
	t.Parallel()

	t.Run("Test NewSession", func(t *testing.T) {
		t.Parallel()

		require.Equal(t,
			fmt.Sprintf("docker cp %s %s%s/bin/", "%s:"+"service", DefaultRemoteDir, "service"),
			"docker cp %s:service /tmp/update_dir/service/bin/")
	})
}

func TestRmDir(t *testing.T) {
	t.Parallel()

	t.Run("Test rmdir", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, b.RmDir(), fmt.Sprintf("rm -rf %s%s/*", DefaultRemoteDir, b.s))
	})
}

func TestCreateBinDir(t *testing.T) {
	t.Parallel()

	t.Run("Test mkdir bin", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, b.MkBinDir(), fmt.Sprintf("mkdir -p %s%s/bin", DefaultRemoteDir, b.s))
	})
	t.Run("Test mkdir checker", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, c.MkBinDir(), fmt.Sprintf("mkdir -p %s%s/bin", DefaultRemoteDir, c.s))
	})
	t.Run("Test mkdir resowatch", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, r.MkBinDir(), fmt.Sprintf("mkdir -p %s%s/bin", DefaultRemoteDir, r.s))
	})
}

func TestCreateDir(t *testing.T) {
	t.Parallel()

	t.Run("Test mkdir", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, b.MkDir(""), fmt.Sprintf("mkdir -p %s%s/", DefaultRemoteDir, b.s))
		require.Equal(t, b.MkBinDir(), fmt.Sprintf("mkdir -p %s%s/bin", DefaultRemoteDir, b.s))
	})
}

func TestPullImg(t *testing.T) {
	t.Parallel()

	t.Run("Test pullimg", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, b.PullImg(), "docker pull registry-dev.ispsystem.net/team/vm/back:%s")
	})
}

func TestCreateContainer(t *testing.T) {
	t.Parallel()

	t.Run("Test create container", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, b.CreateContainer(), "docker create --name %s registry-dev.ispsystem.net/team/vm/back:%s sh")
	})
}

func TestCopyPackages(t *testing.T) {
	t.Parallel()

	t.Run("Test copy packages", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, b.DCpFrom("/python_packages"), "docker cp %s:/python_packages /tmp/update_dir/back")
		require.Equal(t, DockerCpFromCmd("/monitor/"), "docker cp %s:/monitor/ /tmp/update_dir/")
	})
}

func TestCopyVM(t *testing.T) {
	t.Parallel()

	t.Run("Test copy VM", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, b.DCpAllFrom("vm/"), "docker cp %s:/opt/ispsystem/vm/ /tmp/update_dir/back")
	})
}

func TestMvEtc(t *testing.T) {
	t.Parallel()

	t.Run("Test mv etc", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, b.MvInService("vm/etc/", "."), "mv /tmp/update_dir/back/vm/etc/ /tmp/update_dir/back/.")
	})
}

func TestMvBins(t *testing.T) {
	t.Parallel()

	t.Run("Test mv bins", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, b.MvInService("vm/bin/", "."), "mv /tmp/update_dir/back/vm/bin/ /tmp/update_dir/back/.")
		require.Equal(t, m.MvInService(m.s, "bin"), "mv /tmp/update_dir/monitor/monitor /tmp/update_dir/monitor/bin")
		require.Equal(t, c.MvInService(c.s, "bin/"), "mv /tmp/update_dir/checker/checker /tmp/update_dir/checker/bin/")
	})
}

func TestMvScripts(t *testing.T) {
	t.Parallel()

	t.Run("Test mv scripts", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, b.MvInService("vm/scripts/", "."), "mv /tmp/update_dir/back/vm/scripts/ /tmp/update_dir/back/.")
	})
}

func TestCpLibs(t *testing.T) {
	t.Parallel()

	t.Run("Test cp libs", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, b.DCpToDefaultBoxDest("vm/lib/", "vm"),
			"docker cp /tmp/update_dir/back/vm/lib/ vm_box:/opt/ispsystem/vm")
	})
}

func TestCpBins(t *testing.T) {
	t.Parallel()

	t.Run("Test cp bin", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, b.DCpToDefaultBoxDest("bin/", "vm"), "docker cp /tmp/update_dir/back/bin/ vm_box:/opt/ispsystem/vm")
		require.Equal(t,
			r.DCpToDefaultBoxDest(r.s, r.s+"/bin/"),
			"docker cp /tmp/update_dir/resowatch/resowatch vm_box:/opt/ispsystem/resowatch/bin/")
		require.Equal(t, c.DCpBinToBox(), "docker cp /tmp/update_dir/checker/bin/ vm_box:/opt/ispsystem/checker")
		require.Equal(t, m.DCpBinToBox(), "docker cp /tmp/update_dir/monitor/bin/ vm_box:/opt/ispsystem/monitor")
	})
}

func TestDCpSupervisordToBox(t *testing.T) {
	t.Parallel()

	t.Run("Test cp Supervisord to box", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, m.DCpToBoxDest("supervisord.conf.d/", "/etc/"),
			"docker cp /tmp/update_dir/monitor/supervisord.conf.d/ vm_box:/etc/")
	})
}

func TestDCpBinToBox(t *testing.T) {
	t.Parallel()

	t.Run("Test cp bin to box", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, m.DCpBinToBox(), "docker cp /tmp/update_dir/monitor/bin/ vm_box:/opt/ispsystem/monitor")
		// "docker cp ~/update_dir/monitor/supervisord.conf.d/ vm_box:/etc/"
		require.Equal(t, m.DCpToBoxDest("supervisord.conf.d/", "/etc/"),
			"docker cp /tmp/update_dir/monitor/supervisord.conf.d/ vm_box:/etc/")
	})
}

func TestDCpBinTo(t *testing.T) {
	t.Parallel()

	t.Run("Test cp bin to box", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, m.DCpBinTo("vm_box"), "docker cp /tmp/update_dir/monitor/bin/ vm_box:/opt/ispsystem/monitor")
		require.Equal(t, m.DCpBinTo("hell"), "docker cp /tmp/update_dir/monitor/bin/ hell:/opt/ispsystem/monitor")
	})
}

func TestDCpEtcToBox(t *testing.T) {
	t.Parallel()

	t.Run("Test cp etc to box", func(t *testing.T) {
		t.Parallel()

		// "docker cp ~/update_dir/monitor/etc/ vm_box:/opt/ispsystem/monitor/"
		require.Equal(t, m.DCpToDefaultBox("etc/"), "docker cp /tmp/update_dir/monitor/etc/ vm_box:/opt/ispsystem/monitor")
	})
}

func TestDCpToBoxAll(t *testing.T) {
	t.Parallel()

	t.Run("Test cp to box all", func(t *testing.T) {
		t.Parallel()

		// TODO
		// require.Equal(t, r.DCpTo("", "vm_box", ""), "docker cp /tmp/update_dir/resowatch/ vm_box:/opt/ispsystem/")
		require.Equal(t, r.DCpToDefaultBoxDest("", ""), "docker cp /tmp/update_dir/resowatch/ vm_box:/opt/ispsystem/")
		require.Equal(t, r.DCpToBoxAll(), "docker cp /tmp/update_dir/resowatch/ vm_box:/opt/ispsystem/")
	})
}

func TestDcpToSameDirAsService(t *testing.T) {
	t.Parallel()

	t.Run("Test cp to same dir as service", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, m.DCpToSameDirAsService("bin/", "vm_box"),
			"docker cp /tmp/update_dir/monitor/bin/ vm_box:/opt/ispsystem/monitor")
	})
}

func TestCpScripts(t *testing.T) {
	t.Parallel()

	t.Run("Test cp scripts", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, b.DCpToDefaultBoxDest("scripts/", "vm"),
			"docker cp /tmp/update_dir/back/scripts/ vm_box:/opt/ispsystem/vm")
	})
}

func TestCpEtc(t *testing.T) {
	t.Parallel()

	t.Run("Test cp etc", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, b.DCpToDefaultBoxDest("etc/", "vm"), "docker cp /tmp/update_dir/back/etc/ vm_box:/opt/ispsystem/vm")
	})
}

func TestRestart(t *testing.T) {
	t.Parallel()

	t.Run("Test restart", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, m.DRestartBoxService(), "docker exec vm_box supervisorctl restart monitor")
		require.Equal(t, r.DRestartBoxService(), "docker exec vm_box supervisorctl restart resowatch")
		require.Equal(t, r.DRestartService("vm_box"), "docker exec vm_box supervisorctl restart resowatch")
	})
}

func TestCpFromBox(t *testing.T) {
	t.Parallel()

	t.Run("Test cp from ", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, m.DCpFrom("/etc/supervisord.conf.d/"),
			"docker cp %s:/etc/supervisord.conf.d/ /tmp/update_dir/monitor")
	})
}

func TestRsyncArgs(t *testing.T) {
	t.Parallel()

	t.Run("Test rsync args", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, RsyncArgs("/etc/supervisord.conf.d/", "/etc/supervisord.conf.d/"),
			[]string{"/etc/supervisord.conf.d/", "/etc/supervisord.conf.d/"})
	})
}
