/*
Copyright © 2019 orion0616 earth.nobu.light@gmail.com

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

	"github.com/orion0616/sealion/todoist"
	"github.com/spf13/cobra"
)

// getLabelsCmd represents the getLabels command
var getLabelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "list of todoist labels",

	Run: func(cmd *cobra.Command, args []string) {
		client, err := todoist.NewClient()
		if err != nil {
			fmt.Println(err)
		}
		labels, err := client.GetLabels()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("ID         NAME\n")
		for _, label := range labels {
			fmt.Printf("%-10d %s\n", label.ID, label.Name)
		}
	},
}

func init() {
	getCmd.AddCommand(getLabelsCmd)
}
