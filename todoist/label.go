package todoist

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
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

	resp, err := c.HTTPClient.PostForm("https://todoist.com/api/v8/sync", values)
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
		return nil, err
	}
	return getLabelsResult.Labels, nil
}
