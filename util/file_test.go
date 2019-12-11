package util

import "testing"

func TestReadFile(t *testing.T) {
	fileName := "../testdata/sample.txt"
	actual, _ := ReadFile(fileName)

	expected := []string{
		"This is a sample.",
		"Task1 xxx",
		"Task2 yyy",
	}

	if len(actual) != len(expected) {
		t.Errorf("Length is wrong. Actual:\n%v\n Expected:\n%v", actual, expected)
	}
	for i := range actual {
		if actual[i] != expected[i] {
			t.Errorf("Contents is wrong. Actual:\n%v\n Expected:\n%v", actual, expected)
		}
	}
}
