package domain

import "time"

// User represents the user entity in our domain
type User struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"-" bson:"password"` // "-" means this field won't be included in JSON
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
}

// UserUseCase defines the interface for user business logic
type UserUseCase interface {
	Register(email, password string) (*User, error)
	Login(email, password string) (string, error) // Returns JWT token
	ValidateToken(token string) (*User, error)
}
