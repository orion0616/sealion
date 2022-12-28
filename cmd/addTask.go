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
	"strconv"

	"github.com/orion0616/sealion/todoist"
	"github.com/spf13/cobra"
)

// addTaskCmd represents the addTask command
var addTaskCmd = &cobra.Command{
	Use:   "task",
	Short: "add tasks written in a file to a project",
	Long: `A formart of file is like below
<taskname1> <projectname1>
<taskname2> <projectname2>
...`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := todoist.NewClient()
		if err != nil {
			fmt.Println(err)
			return
		}
		projectName, _ := cmd.Flags().GetString("project")
		number, _ := cmd.Flags().GetString("number")
		if projectName != "" && number != "" {
			err = addSeqTasks(client, projectName, number)
		} else {
			err = fmt.Errorf("You cannot use -p/--project and -n/--number alone.")
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("succeded to add tasks.")
	},
}

func addSeqTasks(client *todoist.Client, projectName string, number string) error {
	num, err := strconv.Atoi(number)
	if err != nil {
		return err
	}
	return client.AddSeqTasks(projectName, num)
}

func init() {
	addCmd.AddCommand(addTaskCmd)
	addTaskCmd.Flags().StringP("project", "p", "", "specify a project")
	addTaskCmd.Flags().StringP("number", "n", "", "specify the number of tasks")
}
