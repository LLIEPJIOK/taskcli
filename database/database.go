package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/LLIEPJIOK/taskcli/task"

	_ "github.com/lib/pq"
)

const (
	host     = "db"
	port     = 5432
	user     = "postgres"
	password = "secret"
)

var (
	databaseName = "taskcli"
	dataBase     *sql.DB
)

func Insert(task *task.Task) error {
	_, err := dataBase.Exec(`
		INSERT INTO tasks(name, status, creation_time) 
		VALUES($1, $2, $3)`,
		task.Name, task.Status, task.CreationTime.Format("2006-01-02"))
	if err != nil {
		return fmt.Errorf("error insert %#v in database: %v", *task, err)
	}
	return nil
}

func Delete(id uint) error {
	_, err := dataBase.Exec(`
		DELETE FROM tasks
		WHERE id = $1`,
		id)
	if err != nil {
		return fmt.Errorf("error delete task by id = %v in database: %v", id, err)
	}
	return nil
}

func Update(task *task.Task) error {
	origTask, err := GetTaskById(task.ID)
	if err != nil {
		return err
	}
	origTask.Merge(task)
	_, err = dataBase.Exec(`
		UPDATE tasks
		SET 
			name = $1,
			status = $2
		WHERE id = $3`,
		origTask.Name, origTask.Status, origTask.ID)
	if err != nil {
		return fmt.Errorf("error update task by id = %v in database: %v", task.ID, err)
	}
	return nil
}

func GetAllTasks() ([]task.Task, error) {
	var tasks []task.Task
	rows, err := dataBase.Query(`
		SELECT *
		FROM tasks`,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to get tasks from database: %w", err)
	}
	for rows.Next() {
		var task task.Task
		err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.Status,
			&task.CreationTime,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to read row in database: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func GetTaskById(id uint) (task.Task, error) {
	var task task.Task
	err := dataBase.QueryRow(`
		SELECT *
		FROM tasks
		WHERE id = $1`,
		id).Scan(
		&task.ID,
		&task.Name,
		&task.Status,
		&task.CreationTime,
	)
	if err != nil {
		return task, fmt.Errorf("error get task by id = %v in database: %v", id, err)
	}
	return task, nil
}

func GetTasksByStatus(status string) ([]task.Task, error) {
	var tasks []task.Task
	rows, err := dataBase.Query(`
		SELECT *
		FROM tasks
		WHERE status = $1`,
		status)
	if err != nil {
		return nil, fmt.Errorf("unable to get tasks from database: %w", err)
	}
	for rows.Next() {
		var task task.Task
		err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.Status,
			&task.CreationTime,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to read row in database: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func GetDaysWithQuantity() (map[task.Day]int, error) {
	daysQuantity := make(map[task.Day]int)
	rows, err := dataBase.Query(`
		SELECT 
			creation_time, COUNT(creation_time)
		FROM tasks
		WHERE DATE_PART('month', now() - creation_time) <= 6
		GROUP BY creation_time
		ORDER BY creation_time
		`)
	if err != nil {
		return nil, fmt.Errorf("unable to get DayQuantity from database: %w", err)
	}
	for rows.Next() {
		var day time.Time
		var quantity int
		err := rows.Scan(
			&day,
			&quantity,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to read row DayQuantity in database: %w", err)
		}
		daysQuantity[task.Day{
			Year:  day.Year(),
			Month: int(day.Month()),
			Day:   day.Day(),
		}] = quantity
	}
	return daysQuantity, nil
}

func createDBIfNotExist() error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		host, port, user, password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error open database: %w", err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}
	rows, err := db.Query(`
		SELECT 1 
		FROM pg_database 
		WHERE datname = $1`,
		databaseName)
	if err != nil {
		return fmt.Errorf("error checking database existence: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		_, err = db.Exec(`CREATE DATABASE ` + databaseName)
		if err != nil {
			return fmt.Errorf("error creating database: %w", err)
		}
	}
	return nil
}

func createTableIfNotExist() error {
	_, err := dataBase.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			name CHARACTER VARYING,
			status CHARACTER VARYING,
			creation_time TIMESTAMP WITHOUT TIME ZONE
		)`)
	if err != nil {
		return fmt.Errorf("error create table in database: %v", err)
	}
	return nil
}

func configure() error {
	err := createDBIfNotExist()
	if err != nil {
		return err
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, databaseName)
	dataBase, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error open database: %w", err)
	}
	if err = dataBase.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	return createTableIfNotExist()
}

func Close() {
	if err := dataBase.Close(); err != nil {
		log.Fatal("cannot close db:", err)
	}
}

func init() {
	err := configure()
	if err != nil {
		log.Fatal(err)
	}
}
