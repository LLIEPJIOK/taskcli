package database

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/LLIEPJIOK/taskcli/task"
)

var (
	daysQuantityMap = make(map[task.Day]int)
)

func checkDatabaseExistence() (bool, error) {
	var exists bool
	err := dataBase.QueryRow(`
		SELECT EXISTS (
			SELECT 1 
			FROM pg_database 
			WHERE datname = $1
	)`, databaseName).Scan(&exists)
	return exists, err
}

func deleteDatabase() {
	Close()

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		host, port, user, password)
	db, err := sql.Open("postgres", connStr)
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error while closing database: %v\n", err)
		}
	}()
	if err != nil {
		log.Fatal("error open postgres:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("error connecting to postgres:", err)
	}

	_ = db.QueryRow(`
		DELETE 
			FROM pg_database 
			WHERE datname = $1
		`, databaseName)
}

func checkTableExistence(tableName string) (bool, error) {
	var exists bool
	err := dataBase.QueryRow(`
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.tables 
			WHERE table_name = $1
		)`, tableName).Scan(&exists)
	return exists, err
}

func containsInSlice(tasks []task.Task, target task.Task) bool {
	for _, task := range tasks {
		if task.Equal(&target) {
			return true
		}
	}
	return false
}

func compareMaps(first, second map[task.Day]int) bool {
	if len(first) != len(second) {
		return false
	}
	for key, firstVal := range first {
		secondVal, ok := second[key]
		if !ok || firstVal != secondVal {
			return false
		}
	}
	return true
}

func TestConfiguration(t *testing.T) {
	databaseName = "configuration_test_db"
	defer deleteDatabase()
	for i := 0; i < 2; i++ {
		configure()

		exists, err := checkDatabaseExistence()
		if err != nil {
			t.Fatalf("%v) error while checking database existence: %v", i, err)
		} else if !exists {
			t.Fatalf("%v) database hasn't been created", i)
		}

		exists, err = checkTableExistence("tasks")
		if err != nil {
			t.Fatalf("%v) error while checking User table existence: %v", i, err)
		} else if !exists {
			t.Fatalf("%v) User table hasn't been created", i)
		}
	}
}

func checkGettingAllTasks(t *testing.T, tasks []task.Task) {
	allTasks, err := GetAllTasks()
	if err != nil {
		t.Fatalf("error while getting all tasks: %v", err)
	}
	for _, task := range tasks {
		if !containsInSlice(allTasks, task) {
			t.Fatalf("task %v isn't contained in %v", task, allTasks)
		}
	}
}

func checkGettingTaskByIds(t *testing.T, tasks []task.Task) {
	for _, task := range tasks {
		gettingTask, err := GetTaskById(task.ID)
		if err != nil {
			t.Fatalf("error while getting task by id: %v", err)
		}
		if !task.Equal(&gettingTask) {
			t.Fatalf("find by id: expected: %v, but got: %v", task, gettingTask)
		}
	}
}

func checkGettingTaskByStatus(t *testing.T, tasks []task.Task, status string) {
	gettingTasks, err := GetTasksByStatus(status)
	if err != nil {
		t.Fatalf("error while getting tasks by status = %v: %v", status, err)
	}
	for _, task := range tasks {
		if task.Status == status && !containsInSlice(gettingTasks, task) {
			t.Fatalf("task %v isn't contained in %v", task, gettingTasks)
		}
	}
}

func checkGettingTaskByStatuses(t *testing.T, tasks []task.Task) {
	t.Run("tasks by status = To do", func(t *testing.T) {
		checkGettingTaskByStatus(t, tasks, "To do")
	})
	t.Run("tasks by status = In progress", func(t *testing.T) {
		checkGettingTaskByStatus(t, tasks, "In progress")
	})
	t.Run("tasks by status = Done", func(t *testing.T) {
		checkGettingTaskByStatus(t, tasks, "Done")
	})
}

func checkGettingDaysWithQuantity(t *testing.T) {
	gettingMap, err := GetDaysWithQuantity()
	if err != nil {
		t.Fatalf("error while getting days with quantity: %v", err)
	}
	if !compareMaps(gettingMap, daysQuantityMap) {
		t.Fatalf("expected: %v, but got: %v", daysQuantityMap, gettingMap)
	}
}

func checkAll(t *testing.T, tasks []task.Task) {
	t.Run("all tasks", func(t *testing.T) {
		checkGettingAllTasks(t, tasks)
	})
	t.Run("tasks by id", func(t *testing.T) {
		checkGettingTaskByIds(t, tasks)
	})
	t.Run("tasks by status", func(t *testing.T) {
		checkGettingTaskByStatuses(t, tasks)
	})
	t.Run("days with quantity", func(t *testing.T) {
		checkGettingDaysWithQuantity(t)
	})
}

func TestTasksTable(t *testing.T) {
	tasks := []task.Task{
		{
			ID:           1,
			Name:         "1",
			Status:       task.ToDo.String(),
			CreationTime: time.Now().Add(-time.Minute * 100),
		},
		{
			ID:           2,
			Name:         "3",
			Status:       task.InProgress.String(),
			CreationTime: time.Now().Add(-time.Hour * 30),
		},
		{
			ID:           3,
			Name:         "3",
			Status:       task.Done.String(),
			CreationTime: time.Now().Add(-time.Second * 100),
		},
		{
			ID:           4,
			Name:         "4",
			Status:       task.ToDo.String(),
			CreationTime: time.Now().Add(-time.Minute * 12340),
		},
		{
			ID:           5,
			Name:         "5",
			Status:       task.InProgress.String(),
			CreationTime: time.Now(),
		},
	}

	databaseName = "tasks_test_db"
	configure()
	defer deleteDatabase()

	for _, curTask := range tasks {
		Insert(&curTask)
		daysQuantityMap[task.Day{
			Year:  curTask.CreationTime.Year(),
			Month: int(curTask.CreationTime.Month()),
			Day:   curTask.CreationTime.Day(),
		}]++
	}

	t.Run("before update", func(t *testing.T) {
		checkAll(t, tasks)
	})

	Delete(5)
	daysQuantityMap[task.Day{
		Year:  tasks[4].CreationTime.Year(),
		Month: int(tasks[4].CreationTime.Month()),
		Day:   tasks[4].CreationTime.Day(),
	}]--
	if daysQuantityMap[task.Day{
		Year:  tasks[4].CreationTime.Year(),
		Month: int(tasks[4].CreationTime.Month()),
		Day:   tasks[4].CreationTime.Day(),
	}] == 0 {
		delete(daysQuantityMap, task.Day{
			Year:  tasks[4].CreationTime.Year(),
			Month: int(tasks[4].CreationTime.Month()),
			Day:   tasks[4].CreationTime.Day(),
		})
	}
	tasks = tasks[:len(tasks)-1]

	tasks[1].Name = "new"
	Update(&tasks[1])
	tasks[3].Status = task.Done.String()
	Update(&tasks[3])

	t.Run("after update", func(t *testing.T) {
		checkAll(t, tasks)
	})
}
