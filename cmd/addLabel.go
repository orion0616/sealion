/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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

	"github.com/spf13/cobra"
)

// addLabelCmd represents the addLabel command
var addLabelCmd = &cobra.Command{
	Use:   "label",
	Short: "Add a label to a task",

	Run: func(cmd *cobra.Command, args []string) {
		task, err := cmd.Flags().GetString("task")
		if err != nil {
			fmt.Println(err)
			return
		}
		client, err := todoist.NewClient()
		if err != nil {
			fmt.Println(err)
			return
		}
		err := client.AddLabel(args, task)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("added label successfully")
		}
	},
}

func init() {
	addCmd.AddCommand(addLabelCmd)
	getTasksCmd.Flags().StringP("task", "t", "", "select a task")
}
