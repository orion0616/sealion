package cmd

import (
	"testing"

	"github.com/orion0616/sealion/todoist"
)

func TestCreateResult(t *testing.T) {
	p1 := todoist.Project{ID: 111, Name: "test1"}
	p2 := todoist.Project{ID: 222, Name: "test2"}
	p3 := todoist.Project{ID: 3333333333, Name: "test3"}
	projects := []todoist.Project{p1, p2, p3}

	actual := createResult(projects)
	expected := `ID         NAME
111        test1
222        test2
3333333333 test3
`
	if actual != expected {
		t.Errorf("Actual:\n%v\n Expected:\n%v", actual, expected)
	}
}
