package repo

import (
	"post_service/config"
	"post_service/internal/entity"
	"post_service/internal/interfaces"
	"post_service/pkg/logger"
	"time"
)

var l *logger.Logger
var cfg *config.Config

func init() {
	cfg = config.NewConfig()
	l = logger.New(cfg.App.LogLevel)
}

type PostRepository struct {
	db interfaces.DatabaseDriver
}

func NewPostRepository(db interfaces.DatabaseDriver) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r *PostRepository) GetPosts() ([]entity.ReturnPost, error) {
	query := `
		SELECT id, author_id, title, content, creation_time
		FROM posts
		ORDER BY creation_time DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []entity.ReturnPost
	for rows.Next() {
		var post entity.ReturnPost
		var creationTimeStr string
		err := rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Title,
			&post.Content,
			&creationTimeStr,
		)
		if err != nil {
			return nil, err
		}

		// Parse the MySQL timestamp
		post.CreationTime, err = time.Parse("2006-01-02 15:04:05", creationTimeStr)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) CreatePost(post entity.Post) (entity.Post, error) {
	post.SetDB(r.db)
	if err := post.Create(); err != nil {
		return entity.Post{}, err
	}
	return post, nil
}

func (r *PostRepository) UpdatePost(post entity.Post) (entity.Post, error) {
	post.SetDB(r.db)
	if err := post.Update(); err != nil {
		return entity.Post{}, err
	}
	return post, nil
}

func (r *PostRepository) DeletePost(post entity.Post) (entity.Post, error) {
	post.SetDB(r.db)
	if err := post.Delete(); err != nil {
		return entity.Post{}, err
	}
	return post, nil
}

func (r *PostRepository) GetPost(id int) (entity.ReturnPost, error) {
	post := entity.Post{
		ID: int32(id),
	}
	post.SetDB(r.db)
	if err := post.Get(); err != nil {
		return entity.ReturnPost{}, err
	}
	return entity.ReturnPost{ID: post.ID, AuthorID: post.AuthorID, Title: post.Title, Content: post.Content, CreationTime: post.CreationTime}, nil
}
