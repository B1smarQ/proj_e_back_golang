package entity

import "time"

type MockPostRepository struct{}

func NewMockPostRepository() *MockPostRepository {
	return &MockPostRepository{}
}

func (m *MockPostRepository) GetPosts() ([]ReturnPost, error) {
	return []ReturnPost{
		{
			ID:           "1",
			AuthorID:     "user1",
			AuthorName:   "Test User 1",
			Title:        "Test Post 1",
			Content:      "This is a test post",
			CreationTime: time.Now(),
			ThreadID:     "thread1",
		},
		{
			ID:           "2",
			AuthorID:     "user2",
			AuthorName:   "Test User 2",
			Title:        "Test Post 2",
			Content:      "This is another test post",
			CreationTime: time.Now(),
			ThreadID:     "thread1",
		},
	}, nil
}

func (m *MockPostRepository) CreatePost(post Post) (Post, error) {
	post.ID = "1" // Mock ID assignment
	return post, nil
}

func (m *MockPostRepository) UpdatePost(post Post) (Post, error) {
	return post, nil
}

func (m *MockPostRepository) DeletePost(post Post) (Post, error) {
	return post, nil
}

func (m *MockPostRepository) GetPost(id string) (ReturnPost, error) {
	return ReturnPost{
		ID:           id,
		AuthorID:     "user1",
		AuthorName:   "Test User",
		Title:        "Test Post",
		Content:      "This is a test post",
		CreationTime: time.Now(),
		ThreadID:     "thread1",
	}, nil
}
