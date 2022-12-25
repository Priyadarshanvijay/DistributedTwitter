package storage

import (
	"github.com/twitter/models"
)

type UserStore interface {
	AddUser(*models.User) (*models.User, error)
	GetUser(string) (*models.User, error)
	UpdateUser(*models.User) (*models.User, error)
	FollowUser(*models.User, *models.User) error
	UnFollowUser(curUser *models.User, userToUnFollow *models.User) error
}

type PostStore interface {
	CreatePost(*models.Post) (*models.Post, error)
	DeletePost(*models.Post) error
	GetPosts(*models.User) ([]*models.Post, error)
	GetPost(string) (*models.Post, error)
}

type Storage interface {
	UserStore() UserStore
	PostStore() PostStore
	Close()
}
