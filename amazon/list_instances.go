package amazon

import (
	"github.com/Zumata/capesize/core"
	"github.com/crowdmob/goamz/ec2"
)

func (a *Amazon) ListInstances() {

	filter := ec2.NewFilter()
	filter.Add("tag:capesize", "launching", "ready")
	filter.Add("instance-state-name", "running")

	resp, err := a.Client.DescribeInstances([]string{}, filter)
	if err != nil {
		panic(err)
	}

	machineInfo := []core.MachineStatusInfo{}

	for _, r := range resp.Reservations {
		for _, i := range r.Instances {

			tags := make(map[string]string)

			for _, t := range i.Tags {
				tags[t.Key] = t.Value
			}

			machineInfo = append(machineInfo, core.MachineStatusInfo{
				Name:       tags["Name"],
				Status:     tags["capesize"],
				ProviderId: i.InstanceId,
				DNS:        i.DNSName,
				Launched:   i.LaunchTime,
			})

		}
	}

	core.DisplayMachines(machineInfo)

}
