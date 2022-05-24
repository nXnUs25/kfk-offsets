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
package main

import (
	"github.com/nXnUs25/kfk-offsets/cg"
	"github.com/nXnUs25/kfk-offsets/cmd"
	"github.com/nXnUs25/kfk-offsets/lag"
	"github.com/nXnUs25/kfk-offsets/offset"
)

func main() {
	cmd.Execute()
	if cmd.Root {
		return
	}

	if cmd.Lag {
		cgroup := make(map[string]string, 0)
		for _, v := range cmd.LagCG {
			cgroup[v[cmd.Topic].(string)] = v[cmd.Cg].(string)

		}
		lag.RunLag(cmd.BrokersLag, cgroup)
		return
	}

	if cmd.Offset.Call {
		cgroup := make(map[string]string, 0)
		for _, v := range cmd.OffsetGroups {
			cgroup[v[cmd.Topic].(string)] = v[cmd.Cg].(string)

		}
		offset.RunOffset(cmd.Offset.Action, cmd.BrokersOff, cmd.OffsetGroup, cgroup)
		return
	}
	if cmd.CG.Call {
		cgroup := make(map[string]string, 0)
		for _, v := range cmd.ListCG {
			cgroup[v[cmd.Topic].(string)] = v[cmd.Cg].(string)
		}
		cg.RunCG(cmd.CG.Create, cmd.BrokersCG, cgroup)
		return
	}

}
