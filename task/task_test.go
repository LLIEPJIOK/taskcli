package task

import (
	"testing"
	"time"
)

type taskTestCase struct {
	first  *Task
	second *Task
}

func TestTask(t *testing.T) {
	testCases := []taskTestCase{
		{
			first: &Task{
				ID:           1,
				Name:         "task",
				Status:       ToDo.String(),
				CreationTime: time.Now(),
			},
			second: &Task{
				ID:           1,
				Name:         "updatedTask",
				Status:       Done.String(),
				CreationTime: time.Now(),
			},
		},
		{
			first: &Task{
				ID:           0,
				Name:         "task",
				Status:       InProgress.String(),
				CreationTime: time.Now(),
			},
			second: New("task"),
		},
		{
			first: New("task"),
			second: &Task{
				ID:           0,
				Name:         "updatedTask",
				Status:       ToDo.String(),
				CreationTime: time.Now(),
			},
		},
	}
	for _, testCase := range testCases {
		testCase.first.Merge(testCase.second)
		if !testCase.first.Equal(testCase.second) {
			t.Fatalf("incorrect merge: expected: %v, but got: %v", testCase.first, testCase.second)
		}
	}
}
