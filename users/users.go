package users

import (
	"github.com/twitter/auth"
	"github.com/twitter/models"
	"github.com/twitter/storage"
)

type Service interface {
	GetUser(*models.User) (*models.User, error)
	RegisterUser(*models.User) (*models.User, error)
	FollowUser(*models.User, *models.User) error
	UnFollowUser(*models.User, *models.User) error
}

type UserService struct {
	authService auth.Service
	db          storage.Storage
}

func (us *UserService) GetUser(userToGet *models.User) (*models.User, error) {
	return us.db.UserStore().GetUser(userToGet.UserName)
}

func (us *UserService) RegisterUser(newUser *models.User) (*models.User, error) {
	hashedPassword := us.authService.SecureValue(newUser.UserPassword)
	newUser.UserPassword = hashedPassword
	return us.db.UserStore().AddUser(newUser)
}

func (us *UserService) FollowUser(curUser *models.User, userToFollow *models.User) error {
	return us.db.UserStore().FollowUser(curUser, userToFollow)
}

func (us *UserService) UnFollowUser(curUser *models.User, userToUnFollow *models.User) error {
	return us.db.UserStore().UnFollowUser(curUser, userToUnFollow)
}

func New(as auth.Service, db storage.Storage) Service {
	return &UserService{
		authService: as,
		db:          db,
	}
}
