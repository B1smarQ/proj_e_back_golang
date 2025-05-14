package entity

import "time"

type MockPostRepository struct{}

func NewMockPostRepository() *MockPostRepository {
	return &MockPostRepository{}
}

func (m *MockPostRepository) GetPosts() ([]Post, error) {
	return []Post{
		{
			ID:           1,
			Title:        "Test Post 1",
			Content:      "This is a test post",
			CreationTime: time.Now(),
		},
		{
			ID:           2,
			Title:        "Test Post 2",
			Content:      "This is another test post",
			CreationTime: time.Now(),
		},
	}, nil
}

func (m *MockPostRepository) CreatePost(post Post) (Post, error) {
	post.ID = 1 // Mock ID assignment
	return post, nil
}

func (m *MockPostRepository) UpdatePost(post Post) (Post, error) {
	return post, nil
}

func (m *MockPostRepository) DeletePost(post Post) (Post, error) {
	return post, nil
}

func (m *MockPostRepository) GetPost() (Post, error) {
	return Post{
		ID:           1,
		Title:        "Test Post",
		Content:      "This is a test post",
		CreationTime: time.Now(),
	}, nil
}
