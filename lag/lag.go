package lag

import (
	"fmt"
	"os"

	"github.com/nXnUs25/kfk-offsets/core"
)

func RunLag(brokers []string, cg map[string]string) {
	cgs := make([]string, len(cg))
	topics := make([]string, len(cg))
	i := 0
	for k, v := range cg {
		cgs[i] = v
		topics[i] = k
		i++
	}
	core.SetKafkaCluster(brokers...)
	client := core.NewClient()
	groups := core.NewGroups(cgs...)
	cgroups := core.NewConsumerGroups(groups...)
	cgroups.SetTopicsToConsumerGroups(topics...)
	cgroups.SetGroupOffsets(client)
	fmt.Fprintln(os.Stdout, cgroups)
}
