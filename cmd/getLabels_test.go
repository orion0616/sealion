package cmd

import (
	"testing"

	"github.com/orion0616/sealion/todoist"
)

func TestCreateGetLabelsResult(t *testing.T) {
	l1 := todoist.Label{ID: 111, Name: "test1"}
	l2 := todoist.Label{ID: 222, Name: "test2"}
	l3 := todoist.Label{ID: 3333333333, Name: "test3"}
	labels := []todoist.Label{l1, l2, l3}

	actual := createGetLabelsResult(labels)
	expected := `ID         NAME
111        test1
222        test2
3333333333 test3
`
	if actual != expected {
		t.Errorf("Actual:\n%v\n Expected:\n%v", actual, expected)
	}
}
