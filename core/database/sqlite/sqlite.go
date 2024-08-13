package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"main/core/config"
	"main/core/logger"
	"main/internal/models/tasks"
	"os"
)

const dbDriver = "sqlite3"
const limit = 10

var ErrNoSuchTask = errors.New("no such task")

type Storage struct {
	db *sql.DB
}

var db Storage

func Init() {
	if err := db.initDB(); err != nil {
		logger.Get().Fatal("failed to init database", zap.Error(err))
	}
	return
}

func Get() *Storage {
	if db.db == nil {
		Init()
	}
	return &db
}

func (s *Storage) initDB() error {
	if err := createDB(config.Get().DB.Path); err != nil {
		return err
	}

	db, err := sql.Open(dbDriver, config.Get().DB.Path)
	if err != nil {
		logger.Get().Error("failed to open database")
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	if err = createNewTable(db); err != nil {
		return err
	}

	s.db = db

	return nil
}

func createDB(dbPath string) error {
	if _, err := os.Create(dbPath); err != nil {
		return fmt.Errorf("failed to create db file: %w", err)
	}

	logger.Get().Info("database created")

	return nil
}

func createNewTable(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS scheduler (
		   id INTEGER PRIMARY KEY AUTOINCREMENT,
		   date VARCHAR(8) NOT NULL,
		   title TEXT NOT NULL,
		   comment TEXT DEFAULT "",
		   repeat VARCHAR(128) NOT NULL
   		);
	
   		CREATE INDEX scheduler_date ON scheduler (date);
   `); err != nil {
		return fmt.Errorf("failed to create new table: %w", err)
	}

	return nil
}

func (s *Storage) AddTaskDB(task tasks.Task) (int64, error) {
	result, err := s.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Storage) FindTask(id string) (*tasks.Task, error) {
	task := tasks.Task{}
	if id == "" {
		return nil, errors.New("empty task id")
	}

	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	err := s.db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoSuchTask
		}
		logger.Get().Error("failed to query task", zap.Error(err))
		return nil, err
	}

	return &task, nil
}

func (s *Storage) UpdateTask(task tasks.Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	_, err := s.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoSuchTask
		}
		logger.Get().Error("failed to update task", zap.Error(err))
		return err
	}

	return nil
}

func (s *Storage) Tasks(offset int) ([]tasks.Task, error) {
	query := fmt.Sprintf("SELECT * FROM scheduler ORDER BY date LIMIT %d OFFSET %d", limit, offset)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []tasks.Task
	for rows.Next() {
		var task tasks.Task
		if err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		result = append(result, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Storage) SearchTasks(search string, offset int) ([]tasks.Task, error) {
	query := "SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? LIMIT ? OFFSET ?"
	rows, err := s.db.Query(query, "%"+search+"%", "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []tasks.Task
	for rows.Next() {
		var task tasks.Task
		if err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		result = append(result, task)
	}
	return result, nil
}

func (s *Storage) TasksByDate(date string) ([]tasks.Task, error) {
	query := "SELECT * FROM scheduler WHERE date = ?"
	rows, err := s.db.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []tasks.Task
	for rows.Next() {
		var task tasks.Task
		if err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		result = append(result, task)
	}
	return result, nil
}

func (s *Storage) DoneTask(id string) error {
	var task tasks.Task
	if err := s.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(
		&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		return ErrNoSuchTask
	}

	//if task.Repeat == "" {
	_, err := s.db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return errors.New("задача не найдена")
	}
	//} else {
	//	date, err := task.NextDate(time.Now(), task.Date, task.Repeat)
	//	if err != nil {
	//		return err
	//	}
	//	task.Date = date
	//	if err = s.UpdateTask(task); err != nil {
	//		return err
	//	}
	//}

	return nil
}

func (s *Storage) DeleteTask(id string) error {
	var exists bool
	err := s.db.QueryRow("SELECT exists(SELECT 1 FROM scheduler WHERE id=?)", id).Scan(&exists)
	if err != nil || !exists {
		return ErrNoSuchTask
	}

	if _, err = s.db.Exec("DELETE FROM scheduler WHERE id = ?", id); err != nil {
		return err
	}

	return nil
}
