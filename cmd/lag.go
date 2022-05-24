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
	BrokersLag []string
	LagCG      []lagConsumerGroup
	Lag        bool
)

type (
	lagConsumerGroup map[string]interface{}
)

// lagCmd represents the lag command
var lagCmd = &cobra.Command{
	Use:     "lag",
	Aliases: []string{"lag", "l", "la"},
	Example: example,
	Short:   "Checks topic lag per consumer group.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		Lag = true
		confFlag := !cmd.Flags().Changed(config)
		if confFlag {
			if cmd.Flags().Changed(brokers) {
				brokers, err := cmd.Flags().GetStringArray(brokers)
				cobra.CheckErr(err)
				BrokersLag = brokers
			}
			if cmd.Flags().Changed(cgroup) && cmd.Flags().Changed(Topic) {
				cgcmd, err := cmd.Flags().GetString(cgroup)
				cobra.CheckErr(err)
				topcmd, err := cmd.Flags().GetString(Topic)
				cobra.CheckErr(err)
				LagCG = make([]lagConsumerGroup, 0)
				LagCG = append(LagCG, lagConsumerGroup{
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
	rootCmd.AddCommand(lagCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lagCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	lagCmd.Flags().StringArrayP(brokers, "b", []string{}, brokersUsage)
	lagCmd.Flags().StringP(cgroup, "g", "", consumerGroupNameUsage)
	lagCmd.Flags().StringP(Topic, "t", "", consumerGroupTopicUsage)
}

const (
	example string = `
kfk-offsets lag  - will use data from default configuration file. This way you can configure more Consumer Group instances with a topic.

If flags are specified, the default configuration file will skip equivalent config values/keys. It is allowing monitor only one topic per Consumer Group instance. 

kfk-offsets lag --brokers "kafka01.example.com:9092 kafka02.example.com:9092" --consumer-group cgroup-name --topic mytopic`

	cgroup        string = "consumer-group"
	Topic         string = "topic"
	Cg            string = "cg"
	brokers       string = "brokers"
	lagBrokersKey string = "lag.brokers"
	lagGroupsKey  string = "lag.groups"
)

func ExecuteLag() {
	err := lagCmd.Execute()
	cobra.CheckErr(err)
}
