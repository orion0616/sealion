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
package todoist

import (
	"errors"
	"net/http"
	"os"
)

// Client is todoist client
type Client struct {
	Client http.Client
	Token  string
}

func getToken() (string, error) {
	token := os.Getenv("TODOIST_TOKEN")
	if token == "" {
		return "", errors.New("token is not set")
	}
	return token, nil
}

// NewClient provides new todoist client
func NewClient() (*Client, error) {
	token, err := getToken()
	if err != nil {
		return nil, err
	}
	return &Client{http.Client{}, token}, nil
}
