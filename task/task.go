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

func (t *Task) Equal(other *Task) bool {
	return t.ID == other.ID &&
		t.Name == other.Name &&
		t.Status == other.Status &&
		t.CreationTime.Format("2006-01-02") == other.CreationTime.Format("2006-01-02")
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
