package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"github.com/twitter/models"
	"github.com/twitter/twitter"
	"github.com/twitter/web"
	"google.golang.org/grpc"
)

type Config struct {
	Version         int    `map_structure:"version"`
	ServiceGRPCPort string `map_structure:"serviceGrpcPort"`
	ServiceHostname string `map_structure:"serviceHostname"`
	Hostname        string `map_structure:"hostName"`
	HTTPPort        string `map_structure:"httpPort"`
	ServiceName     string `map_structure:"serviceName"`
}

func GetConfig(config *Config) error {
	viper.AddConfigPath("cmd/web/")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&config)
	return err
}

func Cleanup(serviceConnection *grpc.ClientConn) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		if s == syscall.SIGTERM || s == syscall.SIGINT {
			log.Printf("Gracefully shutting down the server, closing connections...\n")
			err := serviceConnection.Close()
			if err != nil {
				log.Printf("Error while closing the client connection: %v\n", err)
			}
			os.Exit(1)
		}
	}

}

func main() {
	config := &Config{}
	err := GetConfig(config)

	if err != nil {
		log.Fatalf("Unable to load config file, exiting x_x \n%v\n", err)
	} else {
		log.Printf("Config loaded successfully, config file: %+v\n", viper.ConfigFileUsed())
	}

	log.Printf("%s Client version %d\n", config.ServiceName, config.Version)

	twitterConn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", config.ServiceHostname, config.ServiceGRPCPort),
		grpc.WithInsecure(),
	)

	if err != nil {
		log.Fatalf("Could not connect: %v\n", err)
	}

	defer func() {
		log.Printf("Gracefully shutting down the server, closing connections...\n")
		err := twitterConn.Close()
		if err != nil {
			log.Printf("Error while closing the client connection: %v\n", err)
		}
	}()
	go Cleanup(twitterConn)

	twitterClient := twitter.NewTwitterClient(twitterConn)

	if _, err := twitterClient.HealthCheck(context.Background(), &models.Empty{}); err != nil {
		log.Fatalf("%s service not running", config.ServiceName)
	}

	webService := &web.WebService{
		TwitterService: twitterClient,
	}

	http.HandleFunc("/", webService.Index)
	http.HandleFunc("/login", webService.Login)
	http.HandleFunc("/register", webService.Register)
	http.HandleFunc("/home", webService.Home)
	http.HandleFunc("/createPost", webService.CreatePost)
	http.HandleFunc("/followUser", webService.FollowUser)
	http.HandleFunc("/unFollowUser", webService.DeleteFollowing)
	http.HandleFunc("/profile", webService.Profile)
	http.HandleFunc("/otherUser", webService.OtherUser)
	http.HandleFunc("/deletePost", webService.DeletePost)
	http.HandleFunc("/logout", webService.Logout)
	err = http.ListenAndServe(
		fmt.Sprintf("%s:%s", config.Hostname, config.HTTPPort),
		nil,
	)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	log.Printf("Starting server at %s:%s\n", config.Hostname, config.HTTPPort)
}
