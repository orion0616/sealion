package todoist

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	name := "TODOIST_TOKEN"
	current := os.Getenv(name)
	os.Setenv(name, "xxx")
	m.Run()
	os.Setenv(name, current)
}

func TestGetToken(t *testing.T) {
	actual, _ := getToken()
	expected := "xxx"
	if actual != expected {
		t.Errorf("Actual:\n%v\n Expected:\n%v", actual, expected)
	}
}
