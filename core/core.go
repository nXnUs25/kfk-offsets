package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	kfk "gopkg.in/Shopify/sarama.v1"
)

const (
	clientID string = "cg-migration"
)

var (
	KfkC *KafkaCluster
)

type KafkaCluster struct {
	Brokers []string
}

type Group struct {
	Name     string
	Topic    string
	Offsets  []GroupOffset
	TotalLag int64
}

type GroupOffset struct {
	Partition     int32
	CurrentOffset int64
	OldOffset     int64
	Lag           int64
}

type ConsumerGroups struct {
	Groups []Group
}

func (c *ConsumerGroups) String() string {
	var b strings.Builder
	if c.Groups != nil {
		for _, g := range c.Groups {
			fmt.Fprint(&b, g.String())
		}
	}
	return b.String()
}
func (g *Group) String() string {
	var b strings.Builder
	if g != nil {
		fmt.Fprintf(&b, Header, "GROUP", "TOPIC", "PARTITION", "CURRENT-OFFSET", "LOG-END-OFFSET", "LAG")
		fmt.Fprint(&b, g.OffsetsToString(g.Name, g.Topic))
		fmt.Fprintf(&b, TOTAL_LAG, "TOTAL LAG:", g.TotalLag)
	}
	return b.String()
}

func (g *Group) OffsetsToString(name, topic string) string {
	var b strings.Builder
	if g != nil {
		for i, _ := range g.Offsets {
			fmt.Fprintf(&b, CGF, name, topic, g.Offsets[i].String())
		}
	}
	return b.String()
}

func (o *GroupOffset) String() string {
	var b strings.Builder
	if o != nil {
		fmt.Fprintf(&b, OFFSET, o.Partition, o.CurrentOffset, o.OldOffset, o.Lag)
	}
	return b.String()
}

func NewGroup(name string) *Group {
	return &Group{Name: name}
}

func NewGroups(names ...string) []Group {
	groups := make([]Group, 0, len(names))
	for _, name := range names {
		groups = append(groups, *NewGroup(name))
	}
	return groups
}

func (g *Group) SetTopicToConsumerGroup(topicName string) {
	if g != nil {
		g.Topic = topicName
	}
}

func (g *ConsumerGroups) SetTopicsToConsumerGroups(topicNames ...string) {
	if g != nil {
		for i, _ := range g.Groups {
			topic := topicNames[i]
			g.Groups[i].SetTopicToConsumerGroup(topic)
		}
	}
}

func NewConsumerGroups(groups ...Group) *ConsumerGroups {
	if len(groups) > 0 {
		return &ConsumerGroups{Groups: groups}
	} else {
		return &ConsumerGroups{Groups: make([]Group, 0, 5)}
	}
}

const (
	TOTAL_LAG string = "%-130v%-15v\n\n"
	CGF       string = `%-40v %-40v %v`
	OFFSET    string = "%-15v %-15v %-15v %-15v\n"
	Header    string = "%-40v %-40v %-15v %-15v %-15v %-15v\n"
)

// this method is called to rest consumer group offsets on [tpCGs] from [migrationCG]
func MoveConsumerGroupsOffset(client kfk.Client, migrationCG *Group, toCGs *ConsumerGroups) {
	consumer, _ := kfk.NewConsumerFromClient(client)
	defer consumer.Close()
	partitions, _ := consumer.Partitions(migrationCG.Topic)
	for i, _ := range toCGs.Groups {
		offmg, err := kfk.NewOffsetManagerFromClient(toCGs.Groups[i].Name, client)
		if err != nil {
			fmt.Println(err)
		}
		for j, _ := range partitions {
			pom, err := offmg.ManagePartition(toCGs.Groups[i].Topic, toCGs.Groups[i].Offsets[j].Partition)
			if err != nil {
				fmt.Println(err)
			}
			defer pom.Close()
			if migrationCG.Offsets[j].CurrentOffset < toCGs.Groups[i].Offsets[j].CurrentOffset {
				pom.ResetOffset(migrationCG.Offsets[j].CurrentOffset, "")
			} else {
				pom.MarkOffset(migrationCG.Offsets[j].CurrentOffset, "")
			}
		}

	}
}

func MoveConsumersGroupOffset(client kfk.Client, migrationCG *Group, toCG *Group) {
	consumer, _ := kfk.NewConsumerFromClient(client)
	defer consumer.Close()
	partitions, _ := consumer.Partitions(migrationCG.Topic)
	offmg, err := kfk.NewOffsetManagerFromClient(toCG.Name, client)
	if err != nil {
		fmt.Println(err)
	}
	for j, _ := range partitions {
		pom, err := offmg.ManagePartition(toCG.Topic, toCG.Offsets[j].Partition)
		if err != nil {
			fmt.Println(err)
		}
		defer pom.Close()
		if migrationCG.Offsets[j].CurrentOffset < toCG.Offsets[j].CurrentOffset {
			pom.ResetOffset(migrationCG.Offsets[j].CurrentOffset, "")
		} else {
			pom.MarkOffset(migrationCG.Offsets[j].CurrentOffset, "")
		}
	}

}

func (c ConsumerGroups) SetGroupOffsets(client kfk.Client) {
	consumer, _ := kfk.NewConsumerFromClient(client)
	defer consumer.Close()

	for i, _ := range c.Groups {
		offsetMg, _ := kfk.NewOffsetManagerFromClient(c.Groups[i].Name, client)
		defer offsetMg.Close()

		partitions, _ := consumer.Partitions(c.Groups[i].Topic)
		cgOffset := GroupOffset{}
		for _, p := range partitions {
			rofs, err := client.GetOffset(c.Groups[i].Topic, p, kfk.OffsetNewest)
			if err != nil {
				fmt.Println(err)
			}
			pom, err := offsetMg.ManagePartition(c.Groups[i].Topic, p)
			if err != nil {
				fmt.Println(err)
			}
			ofs, _ := pom.NextOffset()
			//jif ofs == kfk.OffsetNewest || ofs == kfk.OffsetOldest {
			//	ofs = 0
			//}
			switch ofs {
			case kfk.OffsetNewest:
				ofs = rofs
			case kfk.OffsetOldest:
				ofs = 0
			}
			cgOffset.Partition = p
			cgOffset.CurrentOffset = ofs
			cgOffset.OldOffset = rofs
			cgOffset.Lag = rofs - ofs
			c.Groups[i].TotalLag += cgOffset.Lag
			c.Groups[i].Offsets = append(c.Groups[i].Offsets, cgOffset)
		}
	}

}

func (g *Group) SetGroupOffsets(client kfk.Client) {
	consumer, _ := kfk.NewConsumerFromClient(client)
	defer consumer.Close()
	offsetMg, _ := kfk.NewOffsetManagerFromClient(g.Name, client)
	defer offsetMg.Close()

	partitions, _ := consumer.Partitions(g.Topic)
	cgOffset := GroupOffset{}
	for _, p := range partitions {
		rofs, err := client.GetOffset(g.Topic, p, kfk.OffsetNewest)
		if err != nil {
			fmt.Println(err)
		}
		pom, err := offsetMg.ManagePartition(g.Topic, p)
		if err != nil {
			fmt.Println(err)
		}

		ofs, _ := pom.NextOffset()
		switch ofs {
		case kfk.OffsetNewest:
			ofs = rofs
		case kfk.OffsetOldest:
			ofs = 0
		}
		cgOffset.Partition = p
		cgOffset.CurrentOffset = ofs
		cgOffset.OldOffset = rofs
		cgOffset.Lag = rofs - ofs
		g.TotalLag += cgOffset.Lag
		g.Offsets = append(g.Offsets, cgOffset)
	}
}

func (k *KafkaCluster) append(brokers ...string) {
	if k != nil {
		k.Brokers = append(k.Brokers, brokers...)
	}
}

func SetKafkaCluster(brokers ...string) {
	KfkC = &KafkaCluster{make([]string, 0, 3)}
	KfkC.append(brokers...)
}

func NewClient() kfk.Client {
	cfg := NewConfig()
	client, err := kfk.NewClient(KfkC.Brokers, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v", err)
	}
	return client
}

func NewConfig() *kfk.Config {
	cfg := kfk.NewConfig()
	cfg.ClientID = clientID
	cfg.Version = kfk.V2_1_0_0
	return cfg
}

func ListAllGroups(client kfk.Client) {
	for _, broker := range client.Brokers() {
		err := broker.Open(client.Config())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Problem opening connection to broker %v - %v\n", broker, err)
		}
		defer broker.Close()
		broker_groups, err := broker.ListGroups(&kfk.ListGroupsRequest{})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Connot retrieving groups from broker")
		}
		for group, _ := range broker_groups.Groups {
			fmt.Fprintln(os.Stdout, group)
		}
	}
}

func (c *ConsumerGroups) append(groups map[string]string) {
	if c != nil {
		for g, _ := range groups {
			c.Groups = append(c.Groups, *NewGroup(g))
		}
	}
}

func (cgroups *ConsumerGroups) CreateConsumerGroup(client kfk.Client) {
	wg := sync.WaitGroup{}
	defer wg.Done()
	wg.Add(len(cgroups.Groups))
	for _, cgroup := range cgroups.Groups {
		cg, err := kfk.NewConsumerGroupFromClient(cgroup.Name, client)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Cannot create new consumer group.", err)
		}
		defer cg.Close()
		om, err := kfk.NewOffsetManagerFromClient(cgroup.Name, client)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Cannot create new offset manager.", err)
		}
		defer om.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Minute*5))
		defer cancel()

		go func(g Group) {
			fmt.Fprintf(os.Stdout, "Creating Consumer group [%s] to consume from topic [%s]\n", g.Name, g.Topic)
			cg.Consume(ctx, []string{g.Name}, &CGHandler{})
			fmt.Fprintf(os.Stdout, "Consumer group [%s] created to consume from topic [%s]\n", g.Name, g.Topic)
			wg.Wait()
		}(cgroup)
	}
}

type CGHandler struct {
	ready chan bool
}

func (CGHandler) Setup(_ kfk.ConsumerGroupSession) error   { return nil }
func (CGHandler) Cleanup(_ kfk.ConsumerGroupSession) error { return nil }
func (h CGHandler) ConsumeClaim(sess kfk.ConsumerGroupSession, claim kfk.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		sess.MarkMessage(message, "")
	}
	return nil
}

func (g *Group) CreateConsumerGroup(brokers []string, topics []string, c kfk.Client) {
	config := c.Config()
	consumer := CGHandler{
		ready: make(chan bool),
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := kfk.NewConsumerGroup(brokers, g.Name, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, topics, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			log.Println("djfad")
			consumer.ready = make(chan bool)
		}

	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}
