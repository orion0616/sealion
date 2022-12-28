package todoist

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/orion0616/sealion/util"
)

// GetLabelsResult express a result of getting all label
type GetLabelsResult struct {
	SyncToken     string `json:"sync_token"`
	TempIDMapping struct {
	} `json:"temp_id_mapping"`
	Labels   []Label
	FullSync bool `json:"full_sync"`
}

// Label express a todoist label
type Label struct {
	ItemOrder  int    `json:"item_order"`
	IsDeleted  int    `json:"is_deleted"`
	Name       string `json:"name"`
	Color      int    `json:"color"`
	IsFavorite int    `json:"is_favorite"`
	ID         int64  `json:"id"`
}

// AddLabels adds labels to tasks in a project
func (c *Client) AddLabels(labelNames []string, project string) error {
	tasks, err := c.GetTasks(project)
	if err != nil {
		return err
	}
	commands := "["
	for _, task := range tasks {
		id, err := strconv.ParseInt(task.ID, 10, 64)
		command, err := makeAddLabelsCommand(labelNames, id)
		if err != nil {
			return err
		}
		commands += command
		commands += ","
	}
	commands = strings.TrimRight(commands, ",")
	commands += "]"

	resp, err := c.do("", commands, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to add labels to projects.\nStatus -> %s\nBody -> %s\n", resp.Status, string(b))
	}
	return nil
}

func makeAddLabelsCommand(labelIDs []string, taskID int64) (string, error) {
	uuid, err := util.CreateUUID()
	if err != nil {
		return "", err
	}
	command := fmt.Sprintf("{\"type\": \"item_update\", \"uuid\": \"%s\", \"args\": {\"id\": %d, \"labels\": %s}}", uuid, taskID, createLabelIDsString(labelIDs))
	return command, nil
}

func createLabelIDsString(ids []string) string {
	var strs []string
	for _, id := range ids {
		strs = append(strs, "\""+id+"\"")
	}
	return "[" + strings.Join(strs, ",") + "]"
}
