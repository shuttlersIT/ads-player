// backend/models/user.go

package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User model
type User struct {
	gorm.Model
	UserID         uint       `gorm:"primaryKey" json:"user_id"`
	Username       string     `gorm:"unique;not null" json:"username"`
	Role           string     `gorm:"not null" json:"role"`
	Email          string     `json:"email" gorm:"unique"`
	Password       string     `gorm:"not null" json:"-"` //Password not exposed in JSON
	ProfilePicture string     `json:"profilePicture"`
	FirstName      string     `json:"firstName"`
	LastName       string     `json:"lastName"`
	Bio            string     `json:"bio"`
	Subscriptions  []Channel  `gorm:"many2many:user_subscriptions;"`
	Playlists      []Playlist `json:"playlists"`
	WatchHistory   []Video    `json:"watchHistory" gorm:"many2many:user_watch_history;"`
	LikedVideos    []Video    `json:"likedVideos" gorm:"many2many:user_liked_videos;"`
	DislikedVideos []Video    `json:"dislikedVideos" gorm:"many2many:user_disliked_videos;"`
}

// TableName sets the table name for the Playlist model.
func (User) TableName() string {
	return "playlists"
}

// UserDBModel provides methods for database operations related to User.
type UserDBModel struct {
	DB *gorm.DB
}

// NewUserDBModel creates a new instance of UserDBModel.
func NewUserDBModel(db *gorm.DB) *UserDBModel {
	return &UserDBModel{DB: db}
}

// GetUserByUsername retrieves a user by username from the database.
func (m *UserDBModel) GetUserByUsername(username string) (*User, error) {
	var user User
	err := m.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

// UpdateUser updates an existing user in the database.
func (m *UserDBModel) UpdateUser(user *User) error {
	return m.DB.Save(user).Error
}

// DeleteUser deletes a user from the database.
func (m *UserDBModel) DeleteUser(userID uint) error {
	return m.DB.Delete(&User{UserID: userID}).Error
}

// CreateUser creates a new user in the database.
func (m *UserDBModel) CreateUser(user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return m.DB.Create(user).Error
}

// AuthenticateUser authenticates a user with the provided username and password.
func (m *UserDBModel) AuthenticateUser(username, password string) (*User, error) {
	var user User
	if err := m.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return &user, nil
}
