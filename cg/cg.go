package cg

import (
	"fmt"
	"os"

	"github.com/nXnUs25/kfk-offsets/core"
)

func RunCG(flag bool, brokers []string, cg map[string]string) {
	core.SetKafkaCluster(brokers...)
	client := core.NewClient()
	if flag && false {
		cgs := make([]string, len(cg))
		topics := make([]string, len(cg))

		i := 0
		for k, v := range cg {
			cgs[i] = v
			topics[i] = k
			i++
		}

		groups := core.NewGroups(cgs...)
		cgroups := core.NewConsumerGroups(groups...)
		cgroups.SetTopicsToConsumerGroups(topics...)

		for _, cg := range cgroups.Groups {
			cg.CreateConsumerGroup(brokers, []string{cg.Topic}, client)
		}
	} else {
		if flag {
			fmt.Fprintln(os.Stdout, "Usage --create work in progress")
		}
		core.ListAllGroups(client)
	}

}
