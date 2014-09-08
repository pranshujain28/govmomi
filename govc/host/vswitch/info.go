/*
Copyright (c) 2014 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vswitch

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/vmware/govmomi/govc/cli"
	"github.com/vmware/govmomi/govc/flags"
)

type info struct {
	*flags.ClientFlag
	*flags.OutputFlag
	*flags.SearchFlag
}

func init() {
	cli.Register("host.vswitch.info", &info{})
}

func (cmd *info) Register(f *flag.FlagSet) {
	cmd.SearchFlag = flags.NewSearchFlag(flags.SearchHosts)
}

func (cmd *info) Process() error { return nil }

func (cmd *info) Run(f *flag.FlagSet) error {
	ns, err := cmd.HostNetworkSystem()
	if err != nil {
		return err
	}

	if err = ns.Properties([]string{"networkInfo.vswitch"}); err != nil {
		return err
	}

	tw := tabwriter.NewWriter(os.Stdout, 2, 0, 2, ' ', 0)

	for i, s := range ns.NetworkInfo.Vswitch {
		if i > 0 {
			fmt.Fprintln(tw)
		}
		fmt.Fprintf(tw, "Name:\t%s\n", s.Name)
		fmt.Fprintf(tw, "Portgroup:\t%s\n", cmd.keys("key-vim.host.PortGroup-", s.Portgroup))
		fmt.Fprintf(tw, "Pnic:\t%s\n", cmd.keys("key-vim.host.PhysicalNic-", s.Pnic))
		fmt.Fprintf(tw, "MTU:\t%d\n", s.Mtu)
		fmt.Fprintf(tw, "Ports:\t%d\n", s.NumPorts)
		fmt.Fprintf(tw, "Ports Available:\t%d\n", s.NumPortsAvailable)
	}

	return tw.Flush()
}

func (cmd *info) keys(key string, vals []string) string {
	for i, val := range vals {
		vals[i] = strings.TrimPrefix(val, key)
	}
	return strings.Join(vals, ", ")
}