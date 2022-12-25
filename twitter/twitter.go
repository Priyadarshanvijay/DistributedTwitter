package twitter

import (
	context "context"

	"github.com/twitter/auth"
	models "github.com/twitter/models"
	"github.com/twitter/posts"
	"github.com/twitter/storage"
	"github.com/twitter/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Server to be implemented for the defined twitter grpc server
type Server struct {
	UnimplementedTwitterServer
	AuthService    auth.Service
	StorageService storage.Storage
	PostService    posts.Service
	UserService    users.Service
}

func (s *Server) HealthCheck(_ context.Context, _ *models.Empty) (*models.Empty, error) {
	return &models.Empty{}, nil
}

func (s *Server) RegisterUser(ctx context.Context, userToCreate *models.User) (*models.User, error) {
	createdUser, err := s.UserService.RegisterUser(userToCreate)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return createdUser, nil
}

func (s *Server) LoginUser(ctx context.Context, userData *models.User) (*models.User, error) {
	storedUser, err := s.UserService.GetUser(userData)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	verifiedUser, err := s.AuthService.ValidateLogin(userData, storedUser)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	token, err := s.AuthService.GenerateToken(verifiedUser)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	mdWithToken := metadata.New(map[string]string{"token": token})
	if err := grpc.SendHeader(ctx, mdWithToken); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to send token")
	}
	return verifiedUser, nil
}

func (s *Server) FollowUser(ctx context.Context, userToFollow *models.User) (*models.Empty, error) {
	requestMadeBy, err := s.getUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	err = s.UserService.FollowUser(requestMadeBy, userToFollow)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &models.Empty{}, nil
}

func (s *Server) UnFollowUser(ctx context.Context, userToUnFollow *models.User) (*models.Empty, error) {
	requestMadeBy, err := s.getUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	err = s.UserService.UnFollowUser(requestMadeBy, userToUnFollow)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &models.Empty{}, nil
}

func (s *Server) CreatePost(ctx context.Context, postToCreate *models.Post) (*models.Post, error) {
	requestMadeBy, err := s.getUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	postToCreate.PostedBy = requestMadeBy.UserName
	postToCreate.PostedAt = timestamppb.Now()
	createdPost, err := s.PostService.CreatePost(postToCreate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to create post")
	}
	return createdPost, nil
}

func (s *Server) GetFeed(ctx context.Context, _ *models.Empty) (*models.MultiplePosts, error) {
	requestMadeBy, err := s.getUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	completeUserData, err := s.UserService.GetUser(requestMadeBy)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	followingList := make([]*models.User, 0)
	for _, curFollowingName := range completeUserData.Follows {
		followingList = append(followingList, &models.User{UserName: curFollowingName})
	}
	feed, err := s.PostService.GetFeed(followingList)
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot generate feed")
	}
	return &models.MultiplePosts{Posts: feed}, nil
}

func (s *Server) DeletePost(ctx context.Context, postToDelete *models.Post) (*models.Empty, error) {
	requestMadeBy, err := s.getUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	p, err := s.PostService.GetPost(postToDelete.PostID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if p == nil {
		return nil, status.Error(codes.Internal, "Post not found")
	}
	if p.PostedBy != requestMadeBy.UserName {
		return nil, status.Error(codes.PermissionDenied, "Only user can delete their posts")
	}
	err = s.PostService.DeletePost(postToDelete)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &models.Empty{}, nil
}

func (s *Server) GetUser(ctx context.Context, userToGet *models.User) (*models.User, error) {
	completeUserData, err := s.UserService.GetUser(userToGet)
	if err != nil {
		return nil, err
	}
	completeUserData.UserPassword = ""
	return completeUserData, nil
}

func (s *Server) GetUserProfile(ctx context.Context, userToGet *models.User) (*models.UserProfile, error) {
	completeUserData, err := s.UserService.GetUser(userToGet)
	if err != nil {
		return nil, err
	}
	completeUserData.UserPassword = ""
	postsToReturn, err := s.PostService.GetAllPosts(completeUserData)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &models.UserProfile{
		User:  completeUserData,
		Posts: postsToReturn,
	}, nil
}

func (s *Server) GetSelf(ctx context.Context, _ *models.Empty) (*models.User, error) {
	requestMadeBy, err := s.getUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	completeUserData, err := s.UserService.GetUser(requestMadeBy)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	completeUserData.UserPassword = ""
	return completeUserData, nil
}

func (s *Server) GetMyPosts(ctx context.Context, _ *models.Empty) (*models.MultiplePosts, error) {
	requestMadeBy, err := s.getUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	postsToReturn, err := s.PostService.GetAllPosts(requestMadeBy)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &models.MultiplePosts{Posts: postsToReturn}, nil
}

func (s *Server) GetPost(ctx context.Context, postToGet *models.Post) (*models.Post, error) {
	_, err := s.getUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	postToReturn, err := s.PostService.GetPost(postToGet.PostID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return postToReturn, nil
}

func (s *Server) getUserFromContext(ctx context.Context) (*models.User, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Missing token")
	}
	token, ok := md["token"]
	if !ok || len(token) == 0 || token[0] == "" {
		return nil, status.Error(codes.Unauthenticated, "Missing token")
	}
	requestMadeBy, err := s.AuthService.VerifyToken(token[0])
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	requestMadeBy.UserPassword = ""
	return requestMadeBy, nil
}
