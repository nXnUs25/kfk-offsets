package offset

import (
	"fmt"

	"github.com/nXnUs25/kfk-offsets/cmd"
	"github.com/nXnUs25/kfk-offsets/core"
)

func RunOffset(flag int, brokers []string, cgFrom string, cgsTo map[string]string) {
	cgs := make([]string, len(cgsTo))
	topics := make([]string, len(cgsTo))
	i := 0
	for k, v := range cgsTo {
		cgs[i] = v
		topics[i] = k
		i++
	}
	//fmt.Println(brokers)
	//fmt.Println(cgs)
	//fmt.Println(topics)
	core.SetKafkaCluster(brokers...)
	client := core.NewClient()
	groups := core.NewGroups(cgs...)
	cgroups := core.NewConsumerGroups(groups...)
	cgroups.SetTopicsToConsumerGroups(topics...)
	cgroups.SetGroupOffsets(client)

	switch flag {
	case cmd.Move:
		fmt.Println("moving")
		for _, cgroup := range cgroups.Groups {
			group := core.NewGroup(cgFrom)
			group.SetTopicToConsumerGroup(cgroup.Topic)
			group.SetGroupOffsets(client)
			core.MoveConsumersGroupOffset(client, group, &cgroup)
		}
		fmt.Println("done")
		//}
	case cmd.Export:
		fmt.Println("exporting")
	case cmd.Import:
		fmt.Println("importing")
	default:
		fmt.Println("Printing")
	}
}
