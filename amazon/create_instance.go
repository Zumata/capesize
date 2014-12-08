package amazon

import (
	"log"
	"os"
	"regexp"
	"time"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/ec2"

	"github.com/Zumata/capesize/command"
	"github.com/Zumata/capesize/config"
	"github.com/Zumata/capesize/core"
)

var DefaultEC2Config AmazonConfig

func init() {

	DefaultEC2Config = AmazonConfig{
		ImageID:          config.SetConfig("EC2_IMAGE_ID", "ami-56b7eb04"),
		InstanceType:     config.SetConfig("EC2_INSTANCE_TYPE", "m3.medium"),
		KeyPairName:      config.SetConfig("EC2_KEY_PAIR_NAME", ""),
		AvailabilityZone: config.SetConfig("EC2_AVAILABILITY_ZONE", "ap-southeast-1b"),
		HostUser:         config.SetConfig("EC2_HOST_USER", "ec2-user"),
	}

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" && os.Getenv("AWS_SECRET_ACCESS_KEY") != "" {
		config.RegisterProvider("amazon", true)
		return
	}
	config.RegisterProvider("amazon", false)
}

type Amazon struct {
	Client *ec2.EC2
	Config *AmazonConfig
}

type AmazonConfig struct {
	ImageID          string
	InstanceType     string
	KeyPairName      string
	AvailabilityZone string
	HostUser         string
}

func NewAmazon(config AmazonConfig) *Amazon {
	auth, err := aws.EnvAuth()
	if err != nil {
		panic("capesize.amazon: Amazon auth environment settings have not been configured")
	}
	client := ec2.New(auth, aws.APSoutheast)
	return &Amazon{Client: client, Config: &config}
}

type SpawnedEC2Instance struct {
	machine       *core.Machine
	client        *ec2.EC2
	config        *AmazonConfig
	InstanceId    string
	DNSName       string
	InstanceState string
	IPAddress     string
}

func NewSpawnedEC2Instance(m *core.Machine, c *ec2.EC2, config *AmazonConfig) *SpawnedEC2Instance {
	return &SpawnedEC2Instance{
		machine: m,
		client:  c,
		config:  config,
	}
}

// Satisfy the Server interface - to allow remote command execution

func (i *SpawnedEC2Instance) User() string {
	return i.config.HostUser
}

func (i *SpawnedEC2Instance) Hostname() string {
	return i.DNSName
}

func (i *SpawnedEC2Instance) DisplayName() string {
	return i.DNSName
}

// helpers

func (i *SpawnedEC2Instance) Run() {

	options := ec2.RunInstancesOptions{
		ImageId:      i.config.ImageID,
		InstanceType: i.config.InstanceType,
		KeyName:      i.config.KeyPairName,
		SecurityGroups: []ec2.SecurityGroup{
			ec2.SecurityGroup{
				Name: i.machine.SecurityGroup,
			},
		},
		MinCount:         1,
		MaxCount:         1,
		AvailabilityZone: i.config.AvailabilityZone,
	}

	// Spawn instance
	resp, err := i.client.RunInstances(&options)
	if err != nil {
		log.Fatalln("Failure to run instance", err)
	}

	i.InstanceId = resp.Instances[0].InstanceId

}

func (i *SpawnedEC2Instance) UpdateDNSInfo() {

	returnedInstance, err := i.client.DescribeInstances([]string{i.InstanceId}, nil)
	if err != nil {
		match, _ := regexp.MatchString("does not exist", err.Error())
		if match {
			return
		}
		log.Fatalln("Failing to call DescribeInstances", err)
	}

	instanceDescription := returnedInstance.Reservations[0].Instances[0]
	i.DNSName = instanceDescription.DNSName
	i.InstanceState = instanceDescription.State.Name
	i.IPAddress = instanceDescription.IPAddress

}

func (i *SpawnedEC2Instance) RunningWithDNSName() bool {
	return i.InstanceState == "running" && i.DNSName != ""
}

func (i *SpawnedEC2Instance) LogStatus() {
	log.Println("Current State: ", i.InstanceState)
}

func (i *SpawnedEC2Instance) AddOrUpdateTag(key, val string) {
	_, err := i.client.CreateTags([]string{i.InstanceId}, []ec2.Tag{{key, val}})
	if err != nil {
		panic(err)
	}
}

func (a *Amazon) CreateInstance(machine *core.Machine, done chan bool) {

	ec2Instance := NewSpawnedEC2Instance(machine, a.Client, a.Config)
	ec2Instance.Run()
	ec2Instance.AddOrUpdateTag("capesize", "launching")

	// wait for DNS
	for {
		ec2Instance.UpdateDNSInfo()
		ec2Instance.LogStatus()
		if ec2Instance.RunningWithDNSName() {
			break
		}
		time.Sleep(5 * time.Second)
	}

	// wait until SSH-able
	for maxRetries := 0; ; maxRetries++ {

		err := command.SSHPing(ec2Instance)

		if err == nil {
			break
		}

		// Wait x seconds and re-ping
		if maxRetries < 30 {
			log.Println("Failed ping", ec2Instance)
			time.Sleep(20 * time.Second)
		}

		// Eventually give up
		if maxRetries >= 30 {
			log.Fatalln("Failed ping (final attempt)", ec2Instance)
		}

	}

	// update machine with info if successful
	machine.AssignedDNS = ec2Instance.DNSName
	machine.IPAddress = ec2Instance.IPAddress

	RunAmazonDockerSetup(machine)
	ec2Instance.AddOrUpdateTag("Name", machine.SuccessTag())
	ec2Instance.AddOrUpdateTag("capesize", "ready")

	done <- true

}

func (a *Amazon) InstanceConfig() core.ProviderMachineConfig {
	return core.ProviderMachineConfig{Login: a.Config.HostUser}
}
