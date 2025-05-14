package entity

import (
	"database/sql"
	"errors"
	"time"
)

type Thread struct {
	ID          int32
	AuthorID    int32
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DB          DBExecutor
}

type ReturnThread struct {
	ID          int32
	AuthorID    int32
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ThreadRepository interface {
	GetThreads() ([]ReturnThread, error)
	CreateThread(thread Thread) (Thread, error)
	UpdateThread(thread Thread) (Thread, error)
	DeleteThread(thread Thread) (Thread, error)
	GetThread(id int) (ReturnThread, error)
	GetThreadPosts(threadID int) ([]ReturnPost, error)
}

var (
	ErrInvalidThread      = errors.New("invalid thread")
	ErrThreadNotFound     = sql.ErrNoRows
	ErrThreadExists       = errors.New("thread already exists")
	ErrThreadTitleEmpty   = errors.New("title is empty")
	ErrThreadDescEmpty    = errors.New("description is empty")
	ErrThreadAuthorIDZero = errors.New("author ID is zero")
)

func NewThread(db DBExecutor) *Thread {
	return &Thread{DB: db}
}

func (t *Thread) Validate() error {
	if t.Title == "" {
		return ErrThreadTitleEmpty
	}
	if t.Description == "" {
		return ErrThreadDescEmpty
	}
	if t.AuthorID == 0 {
		return ErrThreadAuthorIDZero
	}
	return nil
}

func (t *Thread) SetDB(db DBExecutor) {
	t.DB = db
}

func (t *Thread) Create() error {
	if err := t.Validate(); err != nil {
		return err
	}

	query := `
		INSERT INTO threads (author_id, title, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := t.DB.Exec(query, t.AuthorID, t.Title, t.Description, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	t.ID = int32(id)
	return nil
}

func (t *Thread) Update() error {
	if t.ID == 0 {
		return ErrThreadNotFound
	}

	if err := t.Validate(); err != nil {
		return err
	}

	query := `
		UPDATE threads 
		SET title = ?, description = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := t.DB.Exec(query, t.Title, t.Description, t.UpdatedAt, t.ID)
	if err != nil {
		return err
	}

	return nil
}

func (t *Thread) Delete() error {
	if t.ID == 0 {
		return ErrThreadNotFound
	}

	query := `DELETE FROM threads WHERE id = ?`
	_, err := t.DB.Exec(query, t.ID)
	if err != nil {
		return err
	}

	return nil
}

func (t *Thread) Get() error {
	if t.ID == 0 {
		return ErrThreadNotFound
	}

	query := `
		SELECT id, author_id, title, description, created_at, updated_at
		FROM threads
		WHERE id = ?
	`
	row, err := t.DB.QueryRow(query, t.ID)
	if err != nil {
		return err
	}

	var createdAtStr, updatedAtStr string
	err = row.Scan(
		&t.ID,
		&t.AuthorID,
		&t.Title,
		&t.Description,
		&createdAtStr,
		&updatedAtStr,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrThreadNotFound
		}
		return err
	}

	// Parse the MySQL timestamps
	t.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return err
	}
	t.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
	if err != nil {
		return err
	}

	return nil
}
