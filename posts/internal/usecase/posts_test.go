package usecase

import (
	"post_service/internal/entity"
	"testing"
	"time"
)

// MockPostRepository implements PostRepository interface for testing
type MockPostRepository struct {
	posts []entity.Post
}

func (m *MockPostRepository) GetPosts() ([]entity.Post, error) {
	return m.posts, nil
}

func (m *MockPostRepository) CreatePost(post entity.Post) (entity.Post, error) {
	post.ID = int32(len(m.posts) + 1)
	m.posts = append(m.posts, post)
	return post, nil
}

func (m *MockPostRepository) UpdatePost(post entity.Post) (entity.Post, error) {
	for i, p := range m.posts {
		if p.ID == post.ID {
			m.posts[i] = post
			return post, nil
		}
	}
	return entity.Post{}, entity.ErrPostNotFound
}

func (m *MockPostRepository) DeletePost(post entity.Post) (entity.Post, error) {
	for i, p := range m.posts {
		if p.ID == post.ID {
			m.posts = append(m.posts[:i], m.posts[i+1:]...)
			return p, nil
		}
	}
	return entity.Post{}, entity.ErrPostNotFound
}

func (m *MockPostRepository) GetPost() (entity.Post, error) {
	if len(m.posts) > 0 {
		return m.posts[0], nil
	}
	return entity.Post{}, entity.ErrPostNotFound
}

func TestPostUsecase_GetPosts(t *testing.T) {
	mockRepo := &MockPostRepository{
		posts: []entity.Post{
			{
				ID:           1,
				AuthorID:     123,
				Title:        "Test Post 1",
				Content:      "Content 1",
				CreationTime: time.Now(),
			},
			{
				ID:           2,
				AuthorID:     456,
				Title:        "Test Post 2",
				Content:      "Content 2",
				CreationTime: time.Now(),
			},
		},
	}

	useCase := NewPostUsecase(mockRepo)
	posts, err := useCase.GetPostsUsecase()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(posts) != 2 {
		t.Errorf("expected 2 posts, got %d", len(posts))
	}
}

func TestPostUsecase_CreatePost(t *testing.T) {
	mockRepo := &MockPostRepository{}
	useCase := NewPostUsecase(mockRepo)

	post := entity.Post{
		AuthorID:     123,
		Title:        "New Post",
		Content:      "New Content",
		CreationTime: time.Now(),
	}

	createdPost, err := useCase.CreatePostUsecase(post)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if createdPost.ID != 1 {
		t.Errorf("expected post ID 1, got %d", createdPost.ID)
	}
}

func TestPostUsecase_UpdatePost(t *testing.T) {
	mockRepo := &MockPostRepository{
		posts: []entity.Post{
			{
				ID:           1,
				AuthorID:     123,
				Title:        "Original Title",
				Content:      "Original Content",
				CreationTime: time.Now(),
			},
		},
	}

	useCase := NewPostUsecase(mockRepo)

	updatedPost := entity.Post{
		ID:           1,
		AuthorID:     123,
		Title:        "Updated Title",
		Content:      "Updated Content",
		CreationTime: time.Now(),
	}

	post, err := useCase.UpdatePostUsecase(updatedPost)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if post.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got %s", post.Title)
	}
}

func TestPostUsecase_DeletePost(t *testing.T) {
	mockRepo := &MockPostRepository{
		posts: []entity.Post{
			{
				ID:           1,
				AuthorID:     123,
				Title:        "Test Post",
				Content:      "Test Content",
				CreationTime: time.Now(),
			},
		},
	}

	useCase := NewPostUsecase(mockRepo)

	post := entity.Post{ID: 1}
	deletedPost, err := useCase.DeletePostUsecase(post)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if deletedPost.ID != 1 {
		t.Errorf("expected deleted post ID 1, got %d", deletedPost.ID)
	}
}

func TestPostUsecase_GetPost(t *testing.T) {
	mockRepo := &MockPostRepository{
		posts: []entity.Post{
			{
				ID:           1,
				AuthorID:     123,
				Title:        "Test Post",
				Content:      "Test Content",
				CreationTime: time.Now(),
			},
		},
	}

	useCase := NewPostUsecase(mockRepo)

	post, err := useCase.GetPostUsecase(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if post.ID != 1 {
		t.Errorf("expected post ID 1, got %d", post.ID)
	}
}
