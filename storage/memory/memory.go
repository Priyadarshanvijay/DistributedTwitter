package memory

import (
	"errors"
	"fmt"
	"sync"

	"github.com/twitter/models"
	"github.com/twitter/storage"
)

type threadSafeUser struct {
	user     *models.User
	userLock sync.RWMutex
}

type userStore struct {
	mtx      sync.RWMutex
	usersMap map[string]*threadSafeUser
}

type userPostMap struct {
	userPostMtx sync.RWMutex
	lastPostId  int64
	posts       map[string]*models.Post
}

type postStore struct {
	mtx         sync.RWMutex
	postTillNow int64
	userPost    map[string]*userPostMap
	postUser    map[string]string
}

type memory struct {
	users *userStore
	posts *postStore
}

func (m *memory) UserStore() storage.UserStore {
	return m.users
}

func (m *memory) PostStore() storage.PostStore {
	return m.posts
}

func (m *memory) Close() {
}

func (u *userStore) AddUser(newUser *models.User) (*models.User, error) {
	if _, userExistsError := u.GetUser(newUser.UserName); userExistsError == nil {
		return nil, errors.New("User already exists")
	}
	u.mtx.Lock()
	userWithLock := &threadSafeUser{}

	userWithLock.userLock.Lock()
	defer userWithLock.userLock.Unlock()
	u.mtx.Unlock()

	userWithLock.user = newUser
	u.usersMap[newUser.UserName] = userWithLock

	return u.usersMap[newUser.UserName].user, nil
}

func (u *userStore) GetUser(userName string) (*models.User, error) {
	u.mtx.RLock()
	userWithLock, exists := u.usersMap[userName]
	u.mtx.RUnlock()
	if !exists {
		return nil, errors.New("User doesn't exists")
	}
	return userWithLock.user, nil
}

func (u *userStore) UpdateUser(updatedUser *models.User) (*models.User, error) {
	if _, userExistsError := u.GetUser(updatedUser.UserName); userExistsError != nil {
		return nil, userExistsError
	}
	u.mtx.Lock()

	u.usersMap[updatedUser.UserName].userLock.Lock()
	defer u.usersMap[updatedUser.UserName].userLock.Unlock()
	u.mtx.Unlock()

	if updatedUser.UserEmail != "" {
		u.usersMap[updatedUser.UserName].user.UserEmail = updatedUser.UserEmail
	}

	if updatedUser.UserPassword != "" {
		u.usersMap[updatedUser.UserName].user.UserPassword = updatedUser.UserPassword
	}

	return u.usersMap[updatedUser.UserName].user, nil
}

func (u *userStore) FollowUser(curUser *models.User, userToFollow *models.User) error {
	if _, userExistsError := u.GetUser(curUser.UserName); userExistsError != nil {
		return userExistsError
	}
	if _, userExistsError := u.GetUser(userToFollow.UserName); userExistsError != nil {
		return userExistsError
	}
	u.mtx.Lock()

	u.usersMap[curUser.UserName].userLock.Lock()
	u.usersMap[userToFollow.UserName].userLock.Lock()
	defer func() {
		u.usersMap[curUser.UserName].userLock.Unlock()
		u.usersMap[userToFollow.UserName].userLock.Unlock()
	}()
	u.mtx.Unlock()

	u.usersMap[curUser.UserName].user.Follows = append(u.usersMap[curUser.UserName].user.Follows, userToFollow.UserName)
	u.usersMap[userToFollow.UserName].user.Followers = append(u.usersMap[userToFollow.UserName].user.Followers, curUser.UserName)
	return nil
}

func getIndexOfValue(collection []string, value string) int {
	for idx, curValue := range collection {
		if curValue == value {
			return idx
		}
	}
	return -1
}

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (u *userStore) UnFollowUser(curUser *models.User, userToUnFollow *models.User) error {
	if _, userExistsError := u.GetUser(curUser.UserName); userExistsError != nil {
		return userExistsError
	}
	if _, userExistsError := u.GetUser(userToUnFollow.UserName); userExistsError != nil {
		return userExistsError
	}
	u.mtx.Lock()

	u.usersMap[curUser.UserName].userLock.Lock()
	u.usersMap[userToUnFollow.UserName].userLock.Lock()
	defer func() {
		u.usersMap[curUser.UserName].userLock.Unlock()
		u.usersMap[userToUnFollow.UserName].userLock.Unlock()
	}()
	u.mtx.Unlock()
	idxToRemove := getIndexOfValue(u.usersMap[curUser.UserName].user.Follows, userToUnFollow.UserName)
	if idxToRemove >= 0 {
		u.usersMap[curUser.UserName].user.Follows = remove(u.usersMap[curUser.UserName].user.Follows, idxToRemove)
	}
	idxToRemove = getIndexOfValue(u.usersMap[userToUnFollow.UserName].user.Followers, curUser.UserName)
	if idxToRemove >= 0 {
		u.usersMap[userToUnFollow.UserName].user.Followers = remove(u.usersMap[userToUnFollow.UserName].user.Followers, idxToRemove)
	}
	return nil
}

func (p *postStore) CreatePost(newPost *models.Post) (*models.Post, error) {
	p.mtx.Lock()
	if p.userPost[newPost.PostedBy] == nil {
		curUserPostMap := &userPostMap{}
		curUserPostMap.lastPostId = -1
		curUserPostMap.posts = make(map[string]*models.Post)
		p.userPost[newPost.PostedBy] = curUserPostMap
	}

	p.postTillNow += 1

	postId := p.postTillNow

	p.postUser[fmt.Sprint(postId)] = newPost.PostedBy

	// acquire users personal post lock
	p.userPost[newPost.PostedBy].userPostMtx.Lock()
	defer p.userPost[newPost.PostedBy].userPostMtx.Unlock()

	// release the main map lock
	p.mtx.Unlock()

	p.userPost[newPost.PostedBy].lastPostId = postId
	newPost.PostID = fmt.Sprint(postId)

	p.userPost[newPost.PostedBy].posts[fmt.Sprint(postId)] = newPost
	return newPost, nil
}

func (p *postStore) DeletePost(postToDelete *models.Post) error {
	p.mtx.Lock()
	if p.userPost[postToDelete.PostedBy] == nil {
		p.mtx.Unlock()
		return errors.New("Post does not exists")
	}
	delete(p.postUser, postToDelete.PostID)
	p.userPost[postToDelete.PostedBy].userPostMtx.Lock()
	defer p.userPost[postToDelete.PostedBy].userPostMtx.Unlock()
	p.mtx.Unlock()
	if _, postExists := p.userPost[postToDelete.PostedBy].posts[postToDelete.PostID]; !postExists {
		p.userPost[postToDelete.PostedBy].userPostMtx.Unlock()
		return errors.New("Post does not exists")
	}
	delete(p.userPost[postToDelete.PostedBy].posts, postToDelete.PostID)
	return nil
}

func (p *postStore) GetPosts(postedBy *models.User) ([]*models.Post, error) {
	p.mtx.RLock()
	if p.userPost[postedBy.UserName] == nil {
		p.mtx.RUnlock()
		return []*models.Post{}, nil
	}
	p.userPost[postedBy.UserName].userPostMtx.RLock()
	defer p.userPost[postedBy.UserName].userPostMtx.RUnlock()
	p.mtx.RUnlock()

	postsToReturn := make([]*models.Post, 0)

	// TODO: fix the pagination logic
	// postsReadTillNow := 0
	// curPostId := p.userPost[postedBy.UserName].lastPostId

	// if lastReadPost > 0 {
	// 	if int64(lastReadPost) < curPostId {
	// 		curPostId = int64(lastReadPost)
	// 	}
	// }

	// for curPostId > 0 && postsReadTillNow < maxNumOfPosts {
	// 	if validPost, postExists := p.userPost[postedBy.UserName].posts[fmt.Sprint(curPostId)]; postExists {
	// 		postsToReturn = append(postsToReturn, validPost)
	// 		postsReadTillNow -= 1
	// 	}
	// 	curPostId -= 1
	// }
	for _, val := range p.userPost[postedBy.UserName].posts {
		postsToReturn = append(postsToReturn, val)
	}

	return postsToReturn, nil
}

func (p *postStore) GetPost(postId string) (*models.Post, error) {
	p.mtx.RLock()
	createdBy, postExists := p.postUser[postId]
	if !postExists {
		return nil, errors.New("Post doesn't exists")
	}
	p.userPost[createdBy].userPostMtx.RLock()
	p.mtx.RUnlock()
	defer p.userPost[createdBy].userPostMtx.RUnlock()
	postToReturn, postExists := p.userPost[createdBy].posts[postId]
	if !postExists {
		return nil, errors.New("Post doesn't exists")
	}
	return postToReturn, nil
}

func New() storage.Storage {
	m := &memory{}
	m.users = &userStore{}
	m.users.usersMap = make(map[string]*threadSafeUser)
	m.posts = &postStore{
		postTillNow: 0,
	}
	m.posts.userPost = make(map[string]*userPostMap)
	m.posts.postUser = make(map[string]string)
	return m
}
