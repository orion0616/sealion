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
