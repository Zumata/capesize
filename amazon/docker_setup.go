package amazon

import (
	"fmt"
	"os"

	"github.com/Zumata/capesize/command"
	"github.com/Zumata/capesize/core"
)

func RunAmazonDockerSetup(m *core.Machine) {

	commands := []command.RemoteCommand{

		command.SSHExec{
			Log:     "Create the remote dir",
			Server:  m,
			Options: []string{"-n"},
			Command: fmt.Sprintf("mkdir %s", command.DefaultRemoteDirPath),
		},

		command.SSHExec{
			Log:     "Copy developer SSH keys",
			Server:  m,
			Options: []string{"-n"},
			Command: fmt.Sprintf("echo '%s' | cat >> ~/.ssh/authorized_keys", command.DefaultDeveloperKeys),
		},

		command.SSHExec{
			Log:     "Update packages",
			Server:  m,
			Options: []string{"-t", "-t"},
			Command: "sudo yum update -y",
		},

		command.SSHExec{
			Log:     "Install docker",
			Server:  m,
			Options: []string{"-t", "-t"},
			Command: "sudo yum install -y docker",
		},

		command.SSHExec{
			Log:     "Set options for docker daemon",
			Server:  m,
			Options: []string{"-t", "-t"},
			Command: `sudo tee /etc/sysconfig/docker <<<'OPTIONS="` + os.Getenv("DOCKER_OPTS") + `"' > /dev/null`,
		},

		command.SSHExec{
			Log:     "Start service",
			Server:  m,
			Options: []string{"-t", "-t"},
			Command: "sudo service docker start",
		},

		command.SSHExec{
			Log:     "Add ec2-user to docker group",
			Server:  m,
			Options: []string{"-t", "-t"},
			Command: "sudo usermod -a -G docker ec2-user",
		},
	}

	for _, command := range commands {
		command.Exec()
	}
}
