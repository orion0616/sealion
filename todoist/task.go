package todoist

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/orion0616/sealion/util"
)

// ProjectData is data of project and its items
type ProjectData struct {
	Tasks []Task `json:"items"`
}

// Task is correspond to a todoist item
type Task struct {
	IsDeleted      bool        `json:"is_deleted"`
	AssignedByUID  string      `json:"assigned_by_uid"`
	Labels         []string    `json:"labels"`
	SyncID         interface{} `json:"sync_id"`
	SectionID      string      `json:"section_id"`
	InHistory      bool        `json:"in_history"`
	ChildOrder     int         `json:"child_order"`
	DateAdded      time.Time   `json:"added_at"`
	ID             string      `json:"id"`
	Content        string      `json:"content"`
	Checked        bool        `json:"checked"`
	UserID         string      `json:"user_id"`
	Due            interface{} `json:"due"`
	Priority       int         `json:"priority"`
	ParentID       string      `json:"parent_id"`
	ResponsibleUID interface{} `json:"responsible_uid"`
	ProjectID      string      `json:"project_id"`
	Collapsed      bool        `json:"collapsed"`
}

type Due struct {
	Date        string      `json:"date"`
	Timezone    interface{} `json:"timezone"`
	IsRecurring bool        `json:"is_recurring"`
	String      string      `json:"string"`
	Lang        string      `json:"lang"`
}

// GetTasks returns a list of todoist task in a project
func (c *Client) GetTasks(projectName string) ([]Task, error) {
	projects, err := c.GetProjects()
	if err != nil {
		return nil, err
	}
	var projectID int64
	for _, project := range projects {
		if project.Name == projectName {
			projectID, _ = strconv.ParseInt(project.ID, 10, 64)
		}
	}

	values := url.Values{}
	values.Add("project_id", strconv.FormatInt(projectID, 10))
	endpoint := "https://api.todoist.com/sync/v9/projects/get_data"
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	u.RawQuery = values.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	tasks, err := ExtractTasks(resp)
	if err != nil {
		return nil, err
	}
	return tasks, err
}

// ExtractTasks extracts tasks from http.Response
func ExtractTasks(resp *http.Response) ([]Task, error) {
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var projectData ProjectData
	if err := json.Unmarshal(data, &projectData); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal in ExtractTasks. data = " + string(data))
	}
	return projectData.Tasks, nil
}

func (c *Client) AddSeqTasks(projectName string, number int) error {
	commands := "["
	for i := 1; i <= number; i++ {
		taskName := strconv.Itoa(i)
		command, err := c.makeAddTaskCommand(taskName, projectName)
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
		return fmt.Errorf("failed to add task to a project.\nStatus -> %s.\nBody -> %s\n", resp.Status, string(b))
	}
	return nil
}

func (c *Client) makeAddTaskCommand(taskName, projectName string) (string, error) {
	projects, err := c.GetProjects()
	if err != nil {
		return "", err
	}
	var projectID int64
	for _, project := range projects {
		if project.Name == projectName {
			projectID, _ = strconv.ParseInt(project.ID, 10, 64)
		}
	}

	uuid1, err := util.CreateUUID()
	if err != nil {
		return "", err
	}
	uuid2, err := util.CreateUUID()
	if err != nil {
		return "", err
	}
	command := fmt.Sprintf("{\"type\": \"item_add\", \"temp_id\": \"%s\", \"uuid\": \"%s\",\"args\": {\"content\": \"%s\", \"project_id\": %d}}",
		uuid1, uuid2, taskName, projectID)
	return command, nil
}
