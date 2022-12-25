package web

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/twitter/models"
	"github.com/twitter/twitter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Service interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	Home(w http.ResponseWriter, r *http.Request)
	Profile(w http.ResponseWriter, r *http.Request)
	OtherUser(w http.ResponseWriter, r *http.Request)
	DeletePost(w http.ResponseWriter, r *http.Request)
	CreatePost(w http.ResponseWriter, r *http.Request)
	FollowUser(w http.ResponseWriter, r *http.Request)
	DeleteFollowing(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Index(w http.ResponseWriter, r *http.Request)
}

type WebService struct {
	TwitterService twitter.TwitterClient
}

type HomeContext struct {
	Username  string
	Posts     []map[string]string
	Following int
	Followers int
}

type ProfileContext struct {
	Username     string
	Posts        []map[string]string
	FollowingNum int
	FollowersNum int
	Following    []string
	Followers    []string
}

func (ws *WebService) getContextWithToken(r *http.Request) (context.Context, error) {
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}
	token := tokenCookie.Value
	newContext := metadata.AppendToOutgoingContext(r.Context(), "token", token)
	return newContext, nil
}

func (ws *WebService) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("web/login.gtpl")
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/", http.StatusPermanentRedirect)
			return
		}
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		var header metadata.MD
		_, err := ws.TwitterService.LoginUser(
			r.Context(), &models.User{
				UserName:     r.Form.Get("username"),
				UserPassword: r.Form.Get("password"),
			},
			grpc.Header(&header),
		)
		if err != nil {
			fmt.Fprintf(w, "Login Failed : %s", err)
		} else {
			tokens := header.Get("token")
			if len(tokens) == 0 {
				fmt.Fprintf(w, "Failed to generate token : %s", err)
			} else {
				http.SetCookie(w, &http.Cookie{
					Name: "token", Value: tokens[0],
					Expires: time.Now().Add(24 * time.Hour),
				})
				http.Redirect(w, r, "/Home", http.StatusFound)
			}
		}
	}
}

func (ws *WebService) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("web/register.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		_, err := ws.TwitterService.RegisterUser(r.Context(), &models.User{
			UserName:     r.Form.Get("username"),
			UserPassword: r.Form.Get("password"),
		})
		if err != nil {
			fmt.Fprintf(w, "Register Failed : %s", err)
			fmt.Println("error:", err)
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	}
}

func (ws *WebService) Home(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("web/home.gtpl")
		newContext, err := ws.getContextWithToken(r)
		if err != nil {
			fmt.Fprintf(w, "Status Unauthorized")
			return
		}
		self, err := ws.TwitterService.GetSelf(newContext, &models.Empty{})
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		AllPosts := []map[string]string{}

		posts, err := ws.TwitterService.GetFeed(newContext, &models.Empty{})
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		} else {
			for _, post := range posts.Posts {
				AllPosts = append(AllPosts, map[string]string{
					"author":    post.PostedBy,
					"content":   post.Content,
					"createdAt": post.PostedAt.AsTime().Format("15:04, Jan 2, 2006"),
				})
			}
		}
		context := HomeContext{
			Username:  self.UserName,
			Posts:     AllPosts,
			Following: len(self.Follows),
			Followers: len(self.Followers),
		}
		err = t.Execute(w, context)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
	}
}

func (ws *WebService) Profile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("web/profile.gtpl")
		newContext, err := ws.getContextWithToken(r)
		if err != nil {
			fmt.Fprintf(w, "Status Unauthorized")
			return
		}
		self, err := ws.TwitterService.GetSelf(newContext, &models.Empty{})
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		selfProfile, err := ws.TwitterService.GetMyPosts(newContext, &models.Empty{})
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		AllPosts := []map[string]string{}

		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		} else {
			for _, post := range selfProfile.Posts {
				AllPosts = append(AllPosts, map[string]string{
					"author":    post.PostedBy,
					"content":   post.Content,
					"createdAt": post.PostedAt.AsTime().Format("15:04, Jan 2, 2006"),
					"postId":    post.PostID,
				})
			}
		}
		context := ProfileContext{
			Username:     self.UserName,
			Posts:        AllPosts,
			Following:    self.Follows,
			Followers:    self.Followers,
			FollowingNum: len(self.Follows),
			FollowersNum: len(self.Followers),
		}
		err = t.Execute(w, context)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
	}
}

func (ws *WebService) OtherUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("web/otherProfile.gtpl")
		userId := r.URL.Query().Get("id")
		if userId == "" {
			http.Redirect(w, r, "/home", http.StatusNotFound)
		}
		newContext, err := ws.getContextWithToken(r)
		if err != nil {
			fmt.Fprintf(w, "Status Unauthorized")
			return
		}
		userProfile, err := ws.TwitterService.GetUserProfile(newContext, &models.User{
			UserName: userId,
		})
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		self := userProfile.User

		AllPosts := []map[string]string{}

		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		} else {
			for _, post := range userProfile.Posts {
				AllPosts = append(AllPosts, map[string]string{
					"author":    post.PostedBy,
					"content":   post.Content,
					"createdAt": post.PostedAt.AsTime().Format("15:04, Jan 2, 2006"),
					"postId":    post.PostID,
				})
			}
		}
		context := ProfileContext{
			Username:     self.UserName,
			Posts:        AllPosts,
			Following:    self.Follows,
			Followers:    self.Followers,
			FollowingNum: len(self.Follows),
			FollowersNum: len(self.Followers),
		}
		err = t.Execute(w, context)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
	}
}

func (ws *WebService) DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		newContext, err := ws.getContextWithToken(r)
		if err != nil {
			fmt.Fprintf(w, "Status Unauthorized")
			return
		}
		r.ParseForm()
		post := models.Post{
			PostID: r.Form.Get("postId"),
		}
		_, err = ws.TwitterService.DeletePost(newContext, &post)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}
}

func (ws *WebService) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		newContext, err := ws.getContextWithToken(r)
		if err != nil {
			fmt.Fprintf(w, "Status Unauthorized")
			return
		}
		r.ParseForm()
		post := models.Post{
			Content: r.Form.Get("content"),
		}
		_, err = ws.TwitterService.CreatePost(newContext, &post)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func (ws *WebService) FollowUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		newContext, err := ws.getContextWithToken(r)
		if err != nil {
			fmt.Fprintf(w, "Status Unauthorized")
			return
		}
		following := models.User{
			UserName: r.Form.Get("username"),
		}
		_, err = ws.TwitterService.FollowUser(newContext, &following)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func (ws *WebService) DeleteFollowing(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		newContext, err := ws.getContextWithToken(r)
		if err != nil {
			fmt.Fprintf(w, "Status Unauthorized")
			return
		}
		following := models.User{
			UserName: r.Form.Get("username"),
		}
		_, err = ws.TwitterService.UnFollowUser(newContext, &following)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		http.Redirect(w, r, "/profile", http.StatusFound)
	}
}

func (ws *WebService) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Expires: time.Now(),
		})
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func (ws *WebService) Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tokenCookie, err := r.Cookie("token")
		if err != nil || tokenCookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
			return
		}
		http.Redirect(w, r, "/home", http.StatusFound)
	}
}
