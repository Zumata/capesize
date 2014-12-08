package core

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

type MachineStatusInfo struct {
	Name       string
	Status     string
	ProviderId string
	DNS        string
	Launched   string
}

func DisplayMachines(machines []MachineStatusInfo) {

	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	fmt.Fprintln(w, strings.Join([]string{
		"NAME",
		"STATUS",
		"PROVIDER_ID",
		"DNS",
		"LAUNCHED",
	}, "\t"))

	for _, m := range machines {
		fmt.Fprintln(w, strings.Join([]string{
			m.Name,
			m.Status,
			m.ProviderId,
			m.DNS,
			m.Launched,
		}, "\t"))
	}

	fmt.Fprintln(w)
	w.Flush()
}
