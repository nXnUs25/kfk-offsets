/*
Copyright © 2022 Augustyn Chmiel <chmielua@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	BrokersCG []string
	ListCG    []listConsumerGroup
	CG        cmdCG
)

type (
	listConsumerGroup map[string]interface{}
	cmdCG             struct {
		Call   bool
		List   bool
		Create bool
	}
)

// consumerGroupCmd represents the consumerGroup command
var consumerGroupCmd = &cobra.Command{
	Use:     "consumer-group",
	Aliases: []string{"consumergroup", "consumerGroup", "cg", "consumer-group"},
	Short:   "Create / List consumer groups.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Example: exampleCG,
	Run: func(cmd *cobra.Command, args []string) {
		CG = cmdCG{Call: true, List: true, Create: false}
		if cmd.Flags().Changed(create) {
			CG.List = false
			CG.Create = true
		}
		confFlag := !cmd.Flags().Changed(config) || CG.Create
		if confFlag {
			if cmd.Flags().Changed(brokers) {
				var err error
				BrokersCG, err = cmd.Flags().GetStringArray(brokers)
				cobra.CheckErr(err)
			}
			if cmd.Flags().Changed(cgroup) && cmd.Flags().Changed(Topic) {
				cgcmd, err := cmd.Flags().GetString(cgroup)
				cobra.CheckErr(err)
				topcmd, err := cmd.Flags().GetString(Topic)
				cobra.CheckErr(err)
				ListCG = make([]listConsumerGroup, 0)
				ListCG = append(ListCG, listConsumerGroup{
					Cg:    cgcmd,
					Topic: topcmd,
				})
			}
			if !(cmd.Flags().Changed(cgroup) && cmd.Flags().Changed(Topic)) && !viper.IsSet(lagGroupsKey) {
				fmt.Fprintln(os.Stderr, "WARNING missiong a flag")
				cmd.Usage()
				os.Exit(1)
				//fmt.Println(cgroup, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(consumerGroupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// consumerGroupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	consumerGroupCmd.Flags().StringArrayP(brokers, "b", []string{}, brokersUsage)
	consumerGroupCmd.Flags().StringP(cgroup, "g", "", consumerGroupNameUsage)
	consumerGroupCmd.Flags().StringP(Topic, "t", "", consumerGroupTopicUsage)
	consumerGroupCmd.Flags().BoolP(create, "c", false, descCreate)
	consumerGroupCmd.Flags().BoolP(list, "l", true, descList)
}

const (
	exampleCG string = `
kfk-offsets consumer-group  - will use data from default configuration file. This way you can configure more Consumer Group instances with a topic.
--list option is the default behavior.

If flags are specified, the default configuration file will skip equivalent config values/keys. It is allowing to create only one topic per Consumer Group instance. 

kfk-offsets consumer-group --brokers "kafka01.example.com:9092 kafka02.example.com:9092" --consumer-group cgroup-name --topic mytopic [--create|--list]
`

	create         string = "create"
	list           string = "list"
	listBrokersKey string = "cg.brokers"
	listGroupsKey  string = "cg.groups"
	descList       string = "List Consumer Groups"
	descCreate     string = "Create Consumer Groups"
)
