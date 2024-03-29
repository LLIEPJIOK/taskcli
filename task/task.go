package task

import (
	"reflect"
	"time"
)

type Task struct {
	ID           uint
	Name         string
	Status       string
	CreationTime time.Time
}

func New(name string) *Task {
	task := &Task{
		Name:         name,
		Status:       ToDo.String(),
		CreationTime: time.Now(),
	}
	return task
}

func (orig *Task) Merge(t *Task) {
	updateValues := reflect.ValueOf(t).Elem()
	oldValues := reflect.ValueOf(orig).Elem()
	if oldValues.CanSet() {
		for i := 0; i < updateValues.NumField(); i++ {
			updateField := updateValues.Field(i).Interface()
			if v, ok := updateField.(string); ok && v != "" {
				oldValues.Field(i).SetString(v)
			}
		}
	}
}
