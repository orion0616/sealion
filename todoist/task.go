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
	LegacyProjetID int         `json:"legacy_project_id"`
	IsDeleted      int         `json:"is_deleted"`
	AssignedByUID  int         `json:"assigned_by_uid"`
	Labels         []int64     `json:"labels"`
	SyncID         interface{} `json:"sync_id"`
	SectionID      interface{} `json:"section_id"`
	InHistory      int         `json:"in_history"`
	ChildOrder     int         `json:"child_order"`
	DateAdded      time.Time   `json:"date_added"`
	ID             int64       `json:"id"`
	Content        string      `json:"content"`
	Checked        int         `json:"checked"`
	AddedByUID     interface{} `json:"added_by_uid"`
	UserID         int         `json:"user_id"`
	Due            Due         `json:"due"`
	Priority       int         `json:"priority"`
	ParentID       interface{} `json:"parent_id"`
	ResponsibleUID interface{} `json:"responsible_uid"`
	ProjectID      int64       `json:"project_id"`
	DateCompleted  interface{} `json:"date_completed"`
	Collapsed      int         `json:"collapsed"`
}

type Due struct {
	Date        string      `json:"date"`
	Timezone    interface{} `json:"timezone"`
	IsRecurring bool        `json:"is_recurring"`
	String      string      `json:"string"`
	Lang        string      `json:"lang"`
}

func (c *Client) GetAllTasks() ([]Task, error) {
	values := url.Values{}
	values.Add("token", c.Token)
	values.Add("sync_token", "*")
	values.Add("resource_types", "[\"items\"]")

	resp, err := c.HTTPClient.PostForm("https://api.todoist.com/sync/v9/sync", values)
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
	values.Add("token", c.Token)
	values.Add("project_id", strconv.FormatInt(projectID, 10))

	resp, err := c.HTTPClient.PostForm("https://todoist.com/sync/v9/projects/get_data", values)
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
		fmt.Println("Failed to unmarshal in ExtractTasks. data = " + string(data))
		return nil, err
	}
	return projectData.Tasks, nil
}

// AddTasks adds tasks from a file
func (c *Client) AddTasks(fileName string) error {
	lines, err := util.ReadFile(fileName)
	commands := "["
	for _, line := range lines {
		words := strings.Split(line, " ")
		command, err := c.makeAddTaskCommand(words[0], words[1])
		if err != nil {
			return err
		}
		commands += command
		commands += ","
	}
	commands = strings.TrimRight(commands, ",")
	commands += "]"
	values := url.Values{}
	// values.Add("token", c.Token)
	values.Add("commands", commands)

	req, err := http.NewRequest("POST", "https://api.todoist.com/sync/v9/sync", nil)
	req.Header = map[string][]string{
		"Authorization": {"Bearer " + c.Token},
	}
	req.Header.Add("Authorization", "Bearer "+c.Token)
	// values := url.Values{}
	values.Add("sync_token", "*")
	// values.Add("resource_types", "[\"projects\"]")
	req.Form = values
	resp, err := c.HTTPClient.Do(req)

	// resp, err := c.HTTPClient.PostForm("https://api.todoist.com/sync/v9/sync", values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		return fmt.Errorf("failed to add task to a project. Status -> %s", resp.Status)
	}
	return nil
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
	fmt.Println(commands)
	values := url.Values{}
	values.Add("commands", commands)

	endpoint := "https://api.todoist.com/sync/v9/sync"
	u, err := url.Parse(endpoint)
	if err != nil {
		return err
	}
	values.Add("sync_token", "*")
	u.RawQuery = values.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to add task to a project.\nStatus -> %s.\nBody -> %s", resp.Status, string(b))
	} else {
		b, _ := io.ReadAll(resp.Body)
		fmt.Printf("Body -> %s", string(b))
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
