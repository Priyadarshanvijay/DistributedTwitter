package posts

import (
	"sort"
	"sync"

	"github.com/twitter/models"
	"github.com/twitter/storage"
)

type Service interface {
	GetPost(string) (*models.Post, error)
	CreatePost(*models.Post) (*models.Post, error)
	DeletePost(*models.Post) error
	GetAllPosts(*models.User) ([]*models.Post, error)
	GetFeed([]*models.User) ([]*models.Post, error)
}

type PostService struct {
	db storage.Storage
}

func (ps *PostService) GetPost(postId string) (*models.Post, error) {
	return ps.db.PostStore().GetPost(postId)
}

func (ps *PostService) CreatePost(newPost *models.Post) (*models.Post, error) {
	return ps.db.PostStore().CreatePost(newPost)
}

func (ps *PostService) DeletePost(postToDelete *models.Post) error {
	return ps.db.PostStore().DeletePost(postToDelete)
}

func (ps *PostService) GetAllPosts(userData *models.User) ([]*models.Post, error) {
	return ps.db.PostStore().GetPosts(userData)
}

func (ps *PostService) GetFeed(followingList []*models.User) ([]*models.Post, error) {
	feed := make([]*models.Post, 0)
	wg := sync.WaitGroup{}
	wg.Add(len(followingList))

	queue := make(chan []*models.Post, 1)

	for _, curUser := range followingList {
		go func(cu *models.User) {
			latestPosts, err := ps.db.PostStore().GetPosts(cu)
			if err != nil {
				wg.Done()
				return
			}
			queue <- latestPosts
		}(curUser)
	}

	go func() {
		for curPosts := range queue {
			feed = append(feed, curPosts...)
			wg.Done()
		}
	}()

	wg.Wait()
	sort.Slice(feed, func(i int, j int) bool {
		return feed[i].PostedAt.Seconds < feed[j].PostedAt.Seconds
	})
	return feed, nil
}

func New(db storage.Storage) Service {
	return &PostService{
		db: db,
	}
}
