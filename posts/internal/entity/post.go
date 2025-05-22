package entity

import (
	"database/sql"
	"errors"
	"post_service/pkg/logger"
	"time"
)

type DBExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) (*sql.Row, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type Post struct {
	ID           string
	AuthorID     string
	AuthorName   string
	Title        string
	Content      string
	CreationTime time.Time
	ThreadID     string
	DB           DBExecutor
}

type ReturnPost struct {
	ID           string
	AuthorID     string
	AuthorName   string
	Title        string
	Content      string
	CreationTime time.Time
	ThreadID     string
}

type PostRepository interface {
	GetPosts() ([]ReturnPost, error)
	CreatePost(post Post) (Post, error)
	UpdatePost(post Post) (Post, error)
	DeletePost(post Post) (Post, error)
	GetPost(id string) (ReturnPost, error)
}

var (
	ErrInvalidPost      = errors.New("invalid post")
	ErrPostNotFound     = sql.ErrNoRows
	ErrPostExists       = errors.New("post already exists")
	ErrTitleEmpty       = errors.New("title is empty")
	ErrContentEmpty     = errors.New("content is empty")
	ErrAuthorIDEmpty    = errors.New("author id is empty")
	ErrAuthorNameEmpty  = errors.New("author name is empty")
	ErrCreationTimeZero = errors.New("creation time is zero")
	ErrThreadIDEmpty    = errors.New("thread id is empty")
	log                 *logger.Logger
)

func init() {
	log = logger.New("info")
}

func NewPost(db DBExecutor) *Post {
	return &Post{DB: db}
}

func (p *Post) Validate() error {
	if p.Title == "" {
		return ErrTitleEmpty
	}
	if p.Content == "" {
		return ErrContentEmpty
	}
	if p.AuthorID == "" {
		return ErrAuthorIDEmpty
	}
	if p.AuthorName == "" {
		return ErrAuthorNameEmpty
	}
	if p.ThreadID == "" {
		return ErrThreadIDEmpty
	}
	if p.CreationTime.IsZero() {
		return ErrCreationTimeZero
	}
	if err := p.Moderate(); err != nil {
		return err
	}
	return nil
}

func (p *Post) SetDB(db DBExecutor) {
	p.DB = db
}

func (p *Post) Create() error {
	if err := p.Validate(); err != nil {
		return err
	}

	query := `
		INSERT INTO posts (id, author_id, author_name, title, content, creation_time, thread_id)
		VALUES (UUID(), ?, ?, ?, ?, ?, ?)
	`
	_, err := p.DB.Exec(query, p.AuthorID, p.AuthorName, p.Title, p.Content, p.CreationTime, p.ThreadID)
	if err != nil {
		return err
	}

	// Get the generated UUID
	row, err := p.DB.QueryRow("SELECT id FROM posts WHERE id = LAST_INSERT_ID()")
	if err != nil {
		return err
	}
	if err := row.Scan(&p.ID); err != nil {
		return err
	}

	return nil
}

func (p *Post) Update() error {
	if p.ID == "" {
		return ErrPostNotFound
	}

	if err := p.Validate(); err != nil {
		return err
	}

	query := `
		UPDATE posts 
		SET title = ?, content = ?
		WHERE id = ?
	`
	_, err := p.DB.Exec(query, p.Title, p.Content, p.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Post) Delete() error {
	if p.ID == "" {
		return ErrPostNotFound
	}

	query := `DELETE FROM posts WHERE id = ?`
	_, err := p.DB.Exec(query, p.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Post) Get() error {
	if p.ID == "" {
		return ErrPostNotFound
	}

	query := `
		SELECT id, author_id, author_name, title, content, creation_time, thread_id
		FROM posts
		WHERE id = ?
	`
	row, err := p.DB.QueryRow(query, p.ID)
	if err != nil {
		return err
	}

	var creationTimeStr string
	err = row.Scan(
		&p.ID,
		&p.AuthorID,
		&p.AuthorName,
		&p.Title,
		&p.Content,
		&creationTimeStr,
		&p.ThreadID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrPostNotFound
		}
		return err
	}

	// Parse the MySQL timestamp
	p.CreationTime, err = time.Parse("2006-01-02 15:04:05", creationTimeStr)
	if err != nil {
		return err
	}

	return nil
}

// Moderation placeholder
func (p *Post) Moderate() error {
	return nil
}
