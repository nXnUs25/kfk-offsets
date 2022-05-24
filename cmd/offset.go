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
	"github.com/spf13/cobra"
)

var (
	Offset       cmdOffset
	ExportPath   string
	ImportFile   string
	BrokersOff   []string
	OffsetGroup  string
	OffsetGroups []offsetCGs
)

type (
	cmdOffset struct {
		Call   bool
		Action int
	}
	offsetCGs map[string]interface{}
)

// offsetCmd represents the offset command
var offsetCmd = &cobra.Command{
	Use:     "offset",
	Aliases: []string{"offset", "offs", "o", "off"},
	Short:   "A brief description of your command",
	Long: `You can use the [offset] option to reset the offsets of a consumer group to a particular value. 
The tool can be used to reset all offsets on all topics. However, this is something you probably won't ever want to do. 
Therefore, it is highly recommended that you exercise caution when resetting offsets.`,
	Run: func(cmd *cobra.Command, args []string) {
		Offset = cmdOffset{
			Call: true,
		}
		var err error
		if cmd.Flags().Changed(export) {
			Offset.Action = Export
			ExportPath, err = cmd.Flags().GetString(export)
			cobra.CheckErr(err)
		} else if cmd.Flags().Changed(importFile) {
			Offset.Action = Import
			ImportFile, err = cmd.Flags().GetString(importFile)
			cobra.CheckErr(err)
		} else if cmd.Flags().Changed(move) {
			Offset.Action = Move
		} else {
			Offset.Action = Print
		}
	},
}

func init() {
	rootCmd.AddCommand(offsetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// offsetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// offsetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	offsetCmd.Flags().StringP("export-path", "e", "", "Path to exports offset to a file. File name will be autogenerate base on consumer name and topic.")
	offsetCmd.Flags().StringP("import", "i", "", "File to import from the offset, needs CSV format")
	offsetCmd.Flags().BoolP("move", "m", false, "Will move the offset from the consumer group A to the consumer group B for topic Z")
}

const (
	Move          int    = 1
	Export        int    = 2
	Import        int    = 3
	Print         int    = 0
	importFile    string = "import"
	export        string = "export-path"
	move          string = "move"
	offBrokersKey string = "offsetReset.brokers"
	offGroupKey   string = "offsetReset.consumers.from"
	offGroupsKey  string = "offsetReset.consumers.to"
)
