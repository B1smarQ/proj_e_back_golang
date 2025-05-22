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

func (m *MockPostRepository) GetPosts() ([]entity.ReturnPost, error) {
	var returnPosts []entity.ReturnPost
	for _, post := range m.posts {
		returnPosts = append(returnPosts, entity.ReturnPost{
			ID:           post.ID,
			AuthorID:     post.AuthorID,
			AuthorName:   post.AuthorName,
			Title:        post.Title,
			Content:      post.Content,
			CreationTime: post.CreationTime,
			ThreadID:     post.ThreadID,
		})
	}
	return returnPosts, nil
}

func (m *MockPostRepository) CreatePost(post entity.Post) (entity.Post, error) {
	post.ID = "post_" + string(len(m.posts)+1)
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

func (m *MockPostRepository) GetPost(id string) (entity.ReturnPost, error) {
	for _, post := range m.posts {
		if post.ID == id {
			return entity.ReturnPost{
				ID:           post.ID,
				AuthorID:     post.AuthorID,
				AuthorName:   post.AuthorName,
				Title:        post.Title,
				Content:      post.Content,
				CreationTime: post.CreationTime,
				ThreadID:     post.ThreadID,
			}, nil
		}
	}
	return entity.ReturnPost{}, entity.ErrPostNotFound
}

func TestPostUsecase_GetPosts(t *testing.T) {
	mockRepo := &MockPostRepository{
		posts: []entity.Post{
			{
				ID:           "post_1",
				AuthorID:     "user1",
				AuthorName:   "Test User 1",
				Title:        "Test Post 1",
				Content:      "Content 1",
				CreationTime: time.Now(),
				ThreadID:     "thread1",
			},
			{
				ID:           "post_2",
				AuthorID:     "user2",
				AuthorName:   "Test User 2",
				Title:        "Test Post 2",
				Content:      "Content 2",
				CreationTime: time.Now(),
				ThreadID:     "thread1",
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
		AuthorID:     "user1",
		AuthorName:   "Test User",
		Title:        "New Post",
		Content:      "New Content",
		CreationTime: time.Now(),
		ThreadID:     "thread1",
	}

	createdPost, err := useCase.CreatePostUsecase(post)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if createdPost.ID != "post_1" {
		t.Errorf("expected post ID post_1, got %s", createdPost.ID)
	}
}

func TestPostUsecase_UpdatePost(t *testing.T) {
	mockRepo := &MockPostRepository{
		posts: []entity.Post{
			{
				ID:           "post_1",
				AuthorID:     "user1",
				AuthorName:   "Test User",
				Title:        "Original Title",
				Content:      "Original Content",
				CreationTime: time.Now(),
				ThreadID:     "thread1",
			},
		},
	}

	useCase := NewPostUsecase(mockRepo)

	updatedPost := entity.Post{
		ID:           "post_1",
		AuthorID:     "user1",
		AuthorName:   "Test User",
		Title:        "Updated Title",
		Content:      "Updated Content",
		CreationTime: time.Now(),
		ThreadID:     "thread1",
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
				ID:           "post_1",
				AuthorID:     "user1",
				AuthorName:   "Test User",
				Title:        "Test Post",
				Content:      "Test Content",
				CreationTime: time.Now(),
				ThreadID:     "thread1",
			},
		},
	}

	useCase := NewPostUsecase(mockRepo)

	post := entity.Post{ID: "post_1"}
	deletedPost, err := useCase.DeletePostUsecase(post)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if deletedPost.ID != "post_1" {
		t.Errorf("expected deleted post ID post_1, got %s", deletedPost.ID)
	}
}

func TestPostUsecase_GetPost(t *testing.T) {
	mockRepo := &MockPostRepository{
		posts: []entity.Post{
			{
				ID:           "post_1",
				AuthorID:     "user1",
				AuthorName:   "Test User",
				Title:        "Test Post",
				Content:      "Test Content",
				CreationTime: time.Now(),
				ThreadID:     "thread1",
			},
		},
	}

	useCase := NewPostUsecase(mockRepo)

	post, err := useCase.GetPostUsecase("post_1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if post.ID != "post_1" {
		t.Errorf("expected post ID post_1, got %s", post.ID)
	}
}
