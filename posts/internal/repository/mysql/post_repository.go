package mysql

import (
	"database/sql"
	"time"

	"post_service/internal/entity"
)

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) entity.PostRepository {
	return &postRepository{
		db: db,
	}
}

func (r *postRepository) GetPosts() ([]entity.ReturnPost, error) {
	query := `
		SELECT id, author_id, author_name, title, content, creation_time, thread_id
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
			&post.AuthorName,
			&post.Title,
			&post.Content,
			&creationTimeStr,
			&post.ThreadID,
		)
		if err != nil {
			return nil, err
		}

		post.CreationTime, err = time.Parse("2006-01-02 15:04:05", creationTimeStr)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *postRepository) GetPost(id string) (entity.ReturnPost, error) {
	query := `
		SELECT id, author_id, author_name, title, content, creation_time, thread_id
		FROM posts
		WHERE id = ?
	`
	var post entity.ReturnPost
	var creationTimeStr string
	err := r.db.QueryRow(query, id).Scan(
		&post.ID,
		&post.AuthorID,
		&post.AuthorName,
		&post.Title,
		&post.Content,
		&creationTimeStr,
		&post.ThreadID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.ReturnPost{}, entity.ErrPostNotFound
		}
		return entity.ReturnPost{}, err
	}

	post.CreationTime, err = time.Parse("2006-01-02 15:04:05", creationTimeStr)
	if err != nil {
		return entity.ReturnPost{}, err
	}

	return post, nil
}

func (r *postRepository) CreatePost(post entity.Post) (entity.Post, error) {
	if err := post.Create(); err != nil {
		return entity.Post{}, err
	}
	return post, nil
}

func (r *postRepository) UpdatePost(post entity.Post) (entity.Post, error) {
	if err := post.Update(); err != nil {
		return entity.Post{}, err
	}
	return post, nil
}

func (r *postRepository) DeletePost(post entity.Post) (entity.Post, error) {
	if err := post.Delete(); err != nil {
		return entity.Post{}, err
	}
	return post, nil
}
