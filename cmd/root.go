/*
Copyright © 2022 --license

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
	cfgFile       string
	defConfigFile string = fmt.Sprintf("config file (default is $HOME/%v.yaml)", confFileName)
	Root          bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kfk-offsets",
	Short: "Tool to monitor modify kafka consumer groups offsets",
	Long: `This tool is primarily used for describing consumer groups 
and debugging any consumer offset issues, like consumer lag. 
The output from the tool shows the log and consumer offsets for each 
partition connected to the consumer group that is being described.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		Root = true
		cmd.Usage()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, config, "", defConfigFile)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//rootCmd.Flags().StringArrayP("brokers", "b", []string{}, brokersUsage)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".kfk-offsets" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(confFileName)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())

		if viper.IsSet(lagBrokersKey) {
			BrokersLag = make([]string, 0)
			err := viper.UnmarshalKey(lagBrokersKey, &BrokersLag)
			cobra.CheckErr(err)
		}

		if viper.IsSet(lagGroupsKey) {
			LagCG = make([]lagConsumerGroup, 0)
			err := viper.UnmarshalKey(lagGroupsKey, &LagCG)
			cobra.CheckErr(err)
			//fmt.Println("LagCG: ", LagCG)
		}

		if viper.IsSet(listBrokersKey) {
			BrokersCG = make([]string, 0)
			err := viper.UnmarshalKey(listBrokersKey, &BrokersCG)
			cobra.CheckErr(err)
		}

		if viper.IsSet(listGroupsKey) {
			ListCG = make([]listConsumerGroup, 0)
			err := viper.UnmarshalKey(listGroupsKey, &ListCG)
			cobra.CheckErr(err)
		}

		if viper.IsSet(offBrokersKey) {
			BrokersOff = make([]string, 0)
			err := viper.UnmarshalKey(offBrokersKey, &BrokersOff)
			cobra.CheckErr(err)
		}

		if viper.IsSet(offGroupsKey) {
			OffsetGroups = make([]offsetCGs, 0)
			err := viper.UnmarshalKey(offGroupsKey, &OffsetGroups)
			cobra.CheckErr(err)
		}

		if viper.IsSet(offGroupKey) {
			OffsetGroup = ""
			err := viper.UnmarshalKey(offGroupKey, &OffsetGroup)
			cobra.CheckErr(err)
		}
	}
}

const (
	brokersUsage string = `List of brokers to connect to.`

	consumerGroupNameUsage  string = `consumer group name`
	consumerGroupTopicUsage string = `topic name from which consumer group consume messages`
	confFileName            string = `kfk-offset`
	config                  string = "config"
)
