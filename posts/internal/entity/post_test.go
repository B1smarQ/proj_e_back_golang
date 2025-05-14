package entity

import (
	"database/sql"
	"testing"
	"time"
)

// MockDB is a mock implementation of sql.DB
type MockDB struct {
	execFunc func(query string, args ...interface{}) (sql.Result, error)
}

func (m *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return m.execFunc(query, args...)
}

func TestPost_Validate(t *testing.T) {
	tests := []struct {
		name    string
		post    Post
		wantErr error
	}{
		{
			name: "valid post",
			post: Post{
				ID:           1,
				AuthorID:     123,
				Title:        "Test Post",
				Content:      "This is a test post content",
				CreationTime: time.Now(),
			},
			wantErr: nil,
		},
		{
			name: "empty title",
			post: Post{
				ID:           2,
				AuthorID:     123,
				Title:        "",
				Content:      "This is a test post content",
				CreationTime: time.Now(),
			},
			wantErr: ErrTitleEmpty,
		},
		{
			name: "empty content",
			post: Post{
				ID:           3,
				AuthorID:     123,
				Title:        "Test Post",
				Content:      "",
				CreationTime: time.Now(),
			},
			wantErr: ErrContentEmpty,
		},
		{
			name: "zero author id",
			post: Post{
				ID:           4,
				AuthorID:     0,
				Title:        "Test Post",
				Content:      "This is a test post content",
				CreationTime: time.Now(),
			},
			wantErr: ErrAuthorIDZero,
		},
		{
			name: "zero creation time",
			post: Post{
				ID:           5,
				AuthorID:     123,
				Title:        "Test Post",
				Content:      "This is a test post content",
				CreationTime: time.Time{},
			},
			wantErr: ErrCreationTimeZero,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.post.Validate()
			if err != tt.wantErr {
				t.Errorf("Post.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPost_Create(t *testing.T) {
	mockDB := &MockDB{
		execFunc: func(query string, args ...interface{}) (sql.Result, error) {
			// Verify the query and args
			if query != "INSERT INTO posts (author_id, title, content, creation_time) VALUES (?, ?, ?, ?)" {
				t.Errorf("unexpected query: %s", query)
			}
			return &mockResult{lastID: 1}, nil
		},
	}

	validPost := Post{
		AuthorID:     123,
		Title:        "Test Post",
		Content:      "This is a test post content",
		CreationTime: time.Now(),
		db:           mockDB,
	}

	tests := []struct {
		name    string
		post    Post
		wantErr bool
	}{
		{
			name:    "valid post creation",
			post:    validPost,
			wantErr: false,
		},
		{
			name: "invalid post creation - empty title",
			post: Post{
				AuthorID:     123,
				Title:        "",
				Content:      "This is a test post content",
				CreationTime: time.Now(),
				DB:           mockDB,
			},
			wantErr: true,
		},
		{
			name: "invalid post creation - nil db",
			post: Post{
				AuthorID:     123,
				Title:        "Test Post",
				Content:      "This is a test post content",
				CreationTime: time.Now(),
				DB:           nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.post.Create()
			if (err != nil) != tt.wantErr {
				t.Errorf("Post.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// mockResult implements sql.Result for testing
type mockResult struct {
	lastID int64
}

func (m *mockResult) LastInsertId() (int64, error) {
	return m.lastID, nil
}

func (m *mockResult) RowsAffected() (int64, error) {
	return 1, nil
}
