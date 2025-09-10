package services

import (
	"errors"
	"time"

	"github.com/ARQAP/ARQAP-Backend/src/middleware"
	"github.com/ARQAP/ARQAP-Backend/src/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new instance of UserService
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// GetAllUsers retrieves all User records from the database
func (s *UserService) GetAllUsers() ([]models.UserModel, error) {
	var users []models.UserModel
	result := s.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// CreateUser creates a new User record in the database
func (s *UserService) CreateUser(user *models.UserModel) (*models.UserModel, error) {
	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	result := s.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// DeleteUser deletes a User record by ID
func (s *UserService) DeleteUser(id int) error {
	result := s.db.Delete(&models.UserModel{}, id)
	return result.Error
}

// AuthenticateUser checks user credentials and returns a JWT token if valid
func (s *UserService) AuthenticateUser(username, password string) (string, error) {
	var user models.UserModel
	result := s.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid username or password")
		}
		return "", result.Error
	}

	// Compare the provided password with the hashed password in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid username or password")
	}

	claims := jwt.MapClaims{
		"id":  user.Id,
		"exp": time.Now().Add(time.Hour * 12).Unix(), // Token expires in 12 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(middleware.GetSecretKey()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
