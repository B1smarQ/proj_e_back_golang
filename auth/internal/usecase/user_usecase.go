package usecase

import (
	"errors"
	"time"

	"auth/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepo domain.UserRepository
	jwtKey   []byte
}

// NewUserUseCase creates a new instance of UserUseCase
func NewUserUseCase(userRepo domain.UserRepository, jwtKey string) domain.UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
		jwtKey:   []byte(jwtKey),
	}
}

func (uc *userUseCase) Register(email, password string) (*domain.User, error) {
	// Check if user already exists
	existingUser, _ := uc.userRepo.FindByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *userUseCase) Login(email, password string) (string, error) {
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	tokenString, err := token.SignedString(uc.jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (uc *userUseCase) ValidateToken(tokenString string) (*domain.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return uc.jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"].(string)
		return uc.userRepo.FindByID(userID)
	}

	return nil, errors.New("invalid token")
}
