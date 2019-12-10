package cmd

import (
	"testing"

	"github.com/orion0616/sealion/todoist"
)

func TestCreateGetTasksResult(t *testing.T) {
	t1 := todoist.Task{ID: 111, Content: "test1"}
	t2 := todoist.Task{ID: 222, Content: "test2"}
	t3 := todoist.Task{ID: 3333333333, Content: "test3"}
	tasks := []todoist.Task{t1, t2, t3}

	actual := createGetTasksResult(tasks)
	expected := `ID         NAME
111        test1
222        test2
3333333333 test3
`
	if actual != expected {
		t.Errorf("Actual:\n%v\n Expected:\n%v", actual, expected)
	}
}
