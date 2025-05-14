package repo

import (
	"post_service/internal/entity"
	"post_service/internal/interfaces"
	"time"
)

type ThreadRepository struct {
	db interfaces.DatabaseDriver
}

func NewThreadRepository(db interfaces.DatabaseDriver) *ThreadRepository {
	return &ThreadRepository{
		db: db,
	}
}

func (r *ThreadRepository) GetThreads() ([]entity.ReturnThread, error) {
	query := `
		SELECT id, author_id, title, description
		FROM threads
		ORDER BY id DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []entity.ReturnThread
	for rows.Next() {
		var thread entity.ReturnThread
		err := rows.Scan(
			&thread.ID,
			&thread.AuthorID,
			&thread.Title,
			&thread.Description,
		)
		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return threads, nil
}

func (r *ThreadRepository) CreateThread(thread entity.Thread) (entity.Thread, error) {
	thread.SetDB(r.db)
	if err := thread.Create(); err != nil {
		return entity.Thread{}, err
	}
	return thread, nil
}

func (r *ThreadRepository) UpdateThread(thread entity.Thread) (entity.Thread, error) {
	thread.SetDB(r.db)
	if err := thread.Update(); err != nil {
		return entity.Thread{}, err
	}
	return thread, nil
}

func (r *ThreadRepository) DeleteThread(thread entity.Thread) (entity.Thread, error) {
	thread.SetDB(r.db)
	if err := thread.Delete(); err != nil {
		return entity.Thread{}, err
	}
	return thread, nil
}

func (r *ThreadRepository) GetThread(id int) (entity.ReturnThread, error) {
	thread := entity.Thread{
		ID: int32(id),
	}
	thread.SetDB(r.db)
	if err := thread.Get(); err != nil {
		return entity.ReturnThread{}, err
	}
	return entity.ReturnThread{
		ID:          thread.ID,
		AuthorID:    thread.AuthorID,
		Title:       thread.Title,
		Description: thread.Description,
	}, nil
}

func (r *ThreadRepository) GetThreadPosts(threadID int) ([]entity.ReturnPost, error) {
	query := `
		SELECT id, author_id, title, content, creation_time, thread_id
		FROM posts
		WHERE thread_id = ?
		ORDER BY creation_time DESC
	`
	rows, err := r.db.Query(query, threadID)
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
			&post.ThreadID,
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
