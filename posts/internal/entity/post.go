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
	ID           int32
	AuthorID     int32
	Title        string
	Content      string
	CreationTime time.Time
	ThreadID     int32
	DB           DBExecutor
}

type ReturnPost struct {
	ID           int32
	AuthorID     int32
	Title        string
	Content      string
	CreationTime time.Time
	ThreadID     int32
}

type PostRepository interface {
	GetPosts() ([]ReturnPost, error)
	CreatePost(post Post) (Post, error)
	UpdatePost(post Post) (Post, error)
	DeletePost(post Post) (Post, error)
	GetPost(id int) (ReturnPost, error)
}

var (
	ErrInvalidPost      = errors.New("invalid post")
	ErrPostNotFound     = sql.ErrNoRows
	ErrPostExists       = errors.New("post already exists")
	ErrTitleEmpty       = errors.New("title is empty")
	ErrContentEmpty     = errors.New("content is empty")
	ErrAuthorIDZero     = errors.New("author id is zero")
	ErrCreationTimeZero = errors.New("creation time is zero")
	ErrThreadIDZero     = errors.New("thread id is zero")
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
	if p.AuthorID == 0 {
		return ErrAuthorIDZero
	}
	if p.ThreadID == 0 {
		return ErrThreadIDZero
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
		INSERT INTO posts (author_id, title, content, creation_time, thread_id)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := p.DB.Exec(query, p.AuthorID, p.Title, p.Content, p.CreationTime, p.ThreadID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	p.ID = int32(id)
	return nil
}

func (p *Post) Update() error {
	if p.ID == 0 {
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
	if p.ID == 0 {
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
	if p.ID == 0 {
		return ErrPostNotFound
	}

	query := `
		SELECT id, author_id, title, content, creation_time, thread_id
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
