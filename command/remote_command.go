package command

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/Zumata/capesize/config"
)

var DefaultDeveloperKeys string
var DefaultRemoteDirPath string
var DefaultIdentityFile string

func init() {

	DefaultDeveloperKeys = config.SetConfig("DEVELOPER_KEYS", "")
	DefaultRemoteDirPath = config.SetConfig("REMOTE_DIR_PATH", "opt")
	DefaultIdentityFile = config.SetConfig("IDENTITY_FILE", "")

}

type Server interface {
	User() string
	Hostname() string
	DisplayName() string
}

type RemoteCommand interface {
	Exec()
}

type CommandSCP struct {
	Log         string
	Server      Server
	Options     []string
	Source      string
	Sources     []string
	Destination string
}

func (c CommandSCP) Exec() {

	params := []string{"-oStrictHostKeyChecking=no", "-oUserKnownHostsFile=/dev/null", "-oIdentityFile=" + DefaultIdentityFile}

	if len(c.Options) > 0 {
		params = append(params, c.Options...)
	}

	if c.Destination == "" {
		c.Destination = fmt.Sprintf("~/%s", DefaultRemoteDirPath)
	}

	switch {
	case c.Source != "" && len(c.Sources) == 0:
		params = append(params, c.Source)
	case c.Source == "" && len(c.Sources) != 0:
		params = append(params, c.Sources...)
	default:
		panic("Executing SCP with invalid Source(s)")
	}

	params = append(params, fmt.Sprintf("%s@%s:%s", c.Server.User(), c.Server.Hostname(), c.Destination))

	cmd := exec.Command("scp", params...)
	fmt.Println(c.Log)
	fmt.Printf("%s::%s\n", c.Server.DisplayName(), cmd)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Println(cmd)
		panic(err)
	}

}

type SSHExec struct {
	Log     string
	Server  Server
	Options []string
	Command string
}

func (c SSHExec) Exec() {

	params := []string{"-oStrictHostKeyChecking=no", "-oUserKnownHostsFile=/dev/null", "-oIdentityFile=" + DefaultIdentityFile}
	if len(c.Options) > 0 {
		params = append(params, c.Options...)
	}
	params = append(params, fmt.Sprintf("%s@%s", c.Server.User(), c.Server.Hostname()), c.Command)

	cmd := exec.Command("ssh", params...)
	fmt.Println(c.Log)
	fmt.Printf("%s::%s\n", c.Server.DisplayName(), cmd)

	var out bytes.Buffer
	cmd.Stdout = &out

	var outErr bytes.Buffer
	cmd.Stderr = &outErr

	err := cmd.Run()
	if err != nil {
		fmt.Println(cmd)
		fmt.Println("===> stderr")
		fmt.Println(outErr.String())
		fmt.Println("===> panic")
		panic(err)
	}

}

func SSHPing(server Server) error {
	cmd := exec.Command("ssh", "-oStrictHostKeyChecking=no", "-oUserKnownHostsFile=/dev/null", "-oConnectTimeout=1", "-oIdentityFile="+DefaultIdentityFile, server.User()+"@"+server.Hostname(), "echo 'test' 2>/dev/null || true")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return err
}
