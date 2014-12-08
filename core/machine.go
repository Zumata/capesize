package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/Zumata/capesize/config"
)

var BuildIdentifier string
var DefaultSecurityGroup string

func init() {

	BuildIdentifier = config.SetConfig("BUILD_IDENTIFIER", "capesize")
	DefaultSecurityGroup = config.SetConfig("SECURITY_GROUP", "")

}

type Machine struct {
	Name          string
	SecurityGroup string
	Tag           string
	AssignedDNS   string
	IPAddress     string
	Login         string
}

type ProviderMachineConfig struct {
	Login string
}

func NewMachine(machineConfig ProviderMachineConfig) *Machine {
	formatted_date := fmt.Sprintf("%d_%s_%d", time.Now().Year(), time.Now().Month(), time.Now().Day())
	return &Machine{
		Name:          config.GenerateName(),
		SecurityGroup: DefaultSecurityGroup,
		Tag:           fmt.Sprintf("_docker_host_%s", formatted_date),
		Login:         machineConfig.Login,
	}
}

// Satisfy the Server interface - to allow remote command execution
func (m *Machine) User() string {
	return m.Login
}

func (m *Machine) Hostname() string {
	return m.AssignedDNS
}

func (m *Machine) DisplayName() string {
	return m.Name
}

func (m *Machine) SuccessTag() string {
	return strings.Join([]string{
		BuildIdentifier,
		m.Name,
		fmt.Sprintf("%d_%s_%d", time.Now().Year(), time.Now().Month(), time.Now().Day()),
	}, "_")
}
