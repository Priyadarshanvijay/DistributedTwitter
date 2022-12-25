package etcd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/twitter/models"
	"github.com/twitter/storage"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/proto"
)

type userStore struct {
	client     *clientv3.Client
	userPrefix string
	// users who the cur user is followed by
	followersPrefix string
	// users who the cur user follows
	followsPrefix string
}

type postStore struct {
	client      *clientv3.Client
	postsPrefix string
}

type etcd struct {
	client *clientv3.Client
	users  *userStore
	posts  *postStore
}

func (e *etcd) UserStore() storage.UserStore {
	return e.users
}

func (e *etcd) PostStore() storage.PostStore {
	return e.posts
}

func (e *etcd) Close() {
	log.Println("closing etcd connection")
	err := e.client.Close()
	if err != nil {
		log.Fatalln("Unable to close etcd connection")
	}
}

func (u *userStore) AddUser(newUser *models.User) (*models.User, error) {
	if _, userExistsError := u.GetUser(newUser.UserName); userExistsError == nil {
		return nil, errors.New("User already exists")
	}
	userInBytes, err := proto.Marshal(newUser)
	if err != nil {
		return nil, err
	}
	stringifiedUser := string(userInBytes)
	key := fmt.Sprintf("%s/%s/", u.userPrefix, newUser.UserName)
	_, err = u.client.Put(context.Background(), key, stringifiedUser)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (u *userStore) GetUser(userName string) (*models.User, error) {

	key := fmt.Sprintf("%s/%s/", u.userPrefix, userName)
	resp, err := u.client.Get(context.Background(), key)

	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, errors.New("User doesn't exists")
	}

	userToReturn := &models.User{}

	err = proto.Unmarshal([]byte(resp.Kvs[0].Value), userToReturn)

	if err != nil {
		return nil, err
	}

	keyPrefixForFollowing := fmt.Sprintf("%s/%s/", u.followsPrefix, userToReturn.UserName)

	resp, err = u.client.Get(context.Background(), keyPrefixForFollowing, clientv3.WithPrefix(), clientv3.WithKeysOnly())

	if err != nil {
		return nil, err
	}

	for _, curFollowing := range resp.Kvs {
		s := strings.TrimPrefix(string(curFollowing.Key), keyPrefixForFollowing)
		userToReturn.Follows = append(userToReturn.Follows, s)
	}

	keyPrefixForFollowers := fmt.Sprintf("%s/%s/", u.followersPrefix, userToReturn.UserName)

	resp, err = u.client.Get(context.Background(), keyPrefixForFollowers, clientv3.WithPrefix(), clientv3.WithKeysOnly())

	if err != nil {
		return nil, err
	}

	for _, curFollowers := range resp.Kvs {
		s := strings.TrimPrefix(string(curFollowers.Key), keyPrefixForFollowers)
		userToReturn.Followers = append(userToReturn.Followers, s)
	}

	return userToReturn, nil
}

func (u *userStore) UpdateUser(updatedUser *models.User) (*models.User, error) {
	userInDB, userExistsError := u.GetUser(updatedUser.UserName)

	if userExistsError != nil {
		return nil, userExistsError
	}

	if updatedUser.UserEmail != "" {
		userInDB.UserEmail = updatedUser.UserEmail
	}

	if updatedUser.UserPassword != "" {
		userInDB.UserPassword = updatedUser.UserPassword
	}

	userInBytes, err := proto.Marshal(userInDB)
	if err != nil {
		return nil, err
	}
	stringifiedUser := string(userInBytes)
	key := fmt.Sprintf("%s/%s/", u.userPrefix, updatedUser.UserName)
	_, err = u.client.Put(context.Background(), key, stringifiedUser)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (u *userStore) FollowUser(curUser *models.User, userToFollow *models.User) error {
	_, userExistsError := u.GetUser(curUser.UserName)
	if userExistsError != nil {
		return userExistsError
	}
	_, userExistsError = u.GetUser(userToFollow.UserName)
	if userExistsError != nil {
		return userExistsError
	}

	// TODO: make the next 2 statements concurrent

	key := fmt.Sprintf("%s/%s/%s", u.followsPrefix, curUser.UserName, userToFollow.UserName)

	_, err := u.client.Put(context.Background(), key, "")

	if err != nil {
		return err
	}

	key = fmt.Sprintf("%s/%s/%s", u.followersPrefix, userToFollow.UserName, curUser.UserName)

	_, err = u.client.Put(context.Background(), key, "")

	if err != nil {
		return err
	}

	return nil
}

func (u *userStore) UnFollowUser(curUser *models.User, userToUnFollow *models.User) error {
	_, userExistsError := u.GetUser(curUser.UserName)
	if userExistsError != nil {
		return userExistsError
	}
	_, userExistsError = u.GetUser(userToUnFollow.UserName)
	if userExistsError != nil {
		return userExistsError
	}

	// TODO: make the next 2 statements concurrent

	key := fmt.Sprintf("%s/%s/%s", u.followsPrefix, curUser.UserName, userToUnFollow.UserName)

	_, err := u.client.Delete(context.Background(), key)

	if err != nil {
		return err
	}

	key = fmt.Sprintf("%s/%s/%s", u.followersPrefix, userToUnFollow.UserName, curUser.UserName)

	_, err = u.client.Delete(context.Background(), key)

	if err != nil {
		return err
	}

	return nil
}

func (p *postStore) CreatePost(newPost *models.Post) (*models.Post, error) {
	postId := uuid.New()
	newPost.PostID = fmt.Sprintf("%s/%s", newPost.PostedBy, postId.String())

	key := fmt.Sprintf("%s/%s", p.postsPrefix, newPost.PostID)
	postInBytes, err := proto.Marshal(newPost)
	if err != nil {
		return nil, err
	}
	stringifiedPost := string(postInBytes)
	_, err = p.client.Put(context.Background(), key, stringifiedPost)
	if err != nil {
		return nil, err
	}
	return newPost, nil
}

func (p *postStore) DeletePost(postToDelete *models.Post) error {
	key := fmt.Sprintf("%s/%s", p.postsPrefix, postToDelete.PostID)
	_, err := p.client.Delete(context.Background(), key)
	return err
}

func (p *postStore) GetPosts(postedBy *models.User) ([]*models.Post, error) {
	prefixKey := fmt.Sprintf("%s/%s/", p.postsPrefix, postedBy.UserName)
	resp, err := p.client.Get(context.Background(), prefixKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	postsToReturn := make([]*models.Post, 0)
	wg := sync.WaitGroup{}
	wg.Add(len(resp.Kvs))

	queue := make(chan *models.Post, 1)

	for _, curPostResp := range resp.Kvs {
		go func(cPR *mvccpb.KeyValue) {
			curPost := &models.Post{}
			errInUnmarshal := proto.Unmarshal(cPR.Value, curPost)
			if errInUnmarshal != nil {
				wg.Done()
				return
			}
			queue <- curPost
		}(curPostResp)
	}

	go func() {
		for curPost := range queue {
			postsToReturn = append(postsToReturn, curPost)
			wg.Done()
		}
	}()

	wg.Wait()

	return postsToReturn, nil
}

func (p *postStore) GetPost(postId string) (*models.Post, error) {
	key := fmt.Sprintf("%s/%s", p.postsPrefix, postId)
	resp, err := p.client.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) < 1 {
		return nil, errors.New("No matching posts found")
	}
	postToReturn := &models.Post{}
	err = proto.Unmarshal(resp.Kvs[0].Value, postToReturn)
	if err != nil {
		return nil, err
	}
	return postToReturn, nil
}

func New(endpoints []string) (storage.Storage, error) {
	newEtcd := &etcd{}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	newEtcd.client = cli
	newEtcd.posts = &postStore{
		client:      cli,
		postsPrefix: "twitter-key-posts",
	}
	newEtcd.users = &userStore{
		client:          cli,
		userPrefix:      "twitter-key-users",
		followersPrefix: "twitter-key-followers",
		followsPrefix:   "twitter-key-follows",
	}
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
		for {
			s := <-c
			if s == syscall.SIGTERM || s == syscall.SIGINT {
				newEtcd.Close()
				os.Exit(1)
			}
		}

	}()
	return newEtcd, nil
}
