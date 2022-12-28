package todoist

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

// GetLabels returns a list of todoist labels
func (c *Client) GetLabels() ([]Label, error) {
	values := url.Values{}
	values.Add("token", c.Token)
	values.Add("sync_token", "*")
	values.Add("resource_types", "[\"labels\"]")

	resp, err := c.HTTPClient.PostForm("https://todoist.com/api/v9/sync", values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	labels, err := ExtractLabels(resp)
	if err != nil {
		return nil, err
	}
	return labels, err
}

// ExtractLabels extracts labels from http.Response
func ExtractLabels(resp *http.Response) ([]Label, error) {
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var getLabelsResult GetLabelsResult
	if err := json.Unmarshal(data, &getLabelsResult); err != nil {
		fmt.Println("Failed to unmarshal in ExtractLabels. data = " + string(data))
		return nil, err
	}
	return getLabelsResult.Labels, nil
}

// AddLabels adds labels to tasks in a project
func (c *Client) AddLabels(labelNames []string, project string) error {
	tasks, err := c.GetTasks(project)
	if err != nil {
		return err
	}
	commands := "["
	for _, task := range tasks {
		labelIDs, err := c.CreateLabelIDs(labelNames)
		if err != nil {
			return err
		}

		// TODO : using sync API
		command, err := makeAddLabelsCommand(labelIDs, task.ID)
		if err != nil {
			return err
		}
		commands += command
		commands += ","
	}
	commands = strings.TrimRight(commands, ",")
	commands += "]"
	values := url.Values{}
	values.Add("token", c.Token)
	values.Add("commands", commands)

	resp, err := c.HTTPClient.PostForm("https://api.todoist.com/sync/v9/sync", values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		return fmt.Errorf("failed to add labels to projects. Status -> %s", resp.Status)
	}
	return nil
}

func makeAddLabelsCommand(labelIDs []int64, taskID int64) (string, error) {
	uuid, err := util.CreateUUID()
	if err != nil {
		return "", err
	}
	command := fmt.Sprintf("{\"type\": \"item_update\", \"uuid\": \"%s\", \"args\": {\"id\": %d, \"labels\": %s}}", uuid, taskID, createLabelIDsString(labelIDs))
	return command, nil
}

func (c *Client) CreateLabelIDs(labelNames []string) ([]int64, error) {
	var ret []int64
	var added bool
	labels, err := c.GetLabels()
	if err != nil {
		return nil, err
	}
	for _, name := range labelNames {
		added = false
		for _, label := range labels {
			if name == label.Name {
				ret = append(ret, label.ID)
				added = true
				break
			}
		}
		if !added {
			return nil, errors.New("failed to find a label")
		}
	}
	return ret, nil
}

func createLabelIDsString(ids []int64) string {
	var strs []string
	for _, id := range ids {
		strs = append(strs, strconv.FormatInt(id, 10))
	}
	return "[" + strings.Join(strs, ",") + "]"
}
