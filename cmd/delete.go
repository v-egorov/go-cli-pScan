/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"vegorov.ru/go-cli/pScan/scan"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete <host1>...<hostN>",
	Aliases: []string{"d"},
	Short:   "Удалить хост[ы] из списка",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostFile := viper.GetString("hosts-file")
		return deleteAction(os.Stdout, hostFile, args)
	},
}

func init() {
	hostsCmd.AddCommand(deleteCmd)
}

func deleteAction(out io.Writer, hostsFile string, args []string) error {
	hl := &scan.HostsList{}
	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	for _, h := range args {
		if err := hl.Remove(h); err != nil {
			return err
		}
		fmt.Fprintln(out, "Удалён хост:", h)
	}

	return hl.Save(hostsFile)
}
