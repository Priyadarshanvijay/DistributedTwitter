package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/spf13/viper"
	"github.com/twitter/auth"
	"github.com/twitter/posts"
	"github.com/twitter/storage/etcd"
	"github.com/twitter/storage/memory"
	"github.com/twitter/twitter"
	"github.com/twitter/users"
	"google.golang.org/grpc"
)

type Config struct {
	Version            int      `map_structure:"version"`
	GRPCPort           string   `map_structure:"grpcPort"`
	EtcdEndpoints      []string `map_structure:"etcdEndpoints"`
	SigningSecret      string   `map_structure:"signingSecret"`
	MemoryType         string   `map_structure:"memoryType"`
	TokenValidityHours int      `map_structure:"tokenValidityHours"`
	Hostname           string   `map_structure:"hostName"`
}

func GetConfig(config *Config) error {
	viper.AddConfigPath("cmd/server/")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&config)
	return err
}

func main() {
	config := &Config{}
	err := GetConfig(config)

	if err != nil {
		log.Fatalf("Unable to load config file, exiting x_x \n%v\n", err)
	} else {
		log.Printf("Config loaded successfully, config file: %+v\n", viper.ConfigFileUsed())
	}

	log.Printf("Twitter version %d\n", config.Version)

	listner, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.Hostname, config.GRPCPort))
	log.Printf("Starting server at %s:%s\n", config.Hostname, config.GRPCPort)
	if err != nil {
		log.Fatalln(err)
	}
	s := grpc.NewServer()
	twtServer := &twitter.Server{}

	if config.MemoryType == "memory" {
		twtServer.StorageService = memory.New()
	} else if config.MemoryType == "raft" {
		storageService, err := etcd.New(config.EtcdEndpoints)
		if err != nil {
			log.Println("Unable to connect to etcd")
			log.Fatal(err)
		}
		twtServer.StorageService = storageService
	} else {
		log.Fatalf("Unrecognized type of memory supplied: %s\n", config.MemoryType)
	}

	twtServer.AuthService = auth.New(time.Duration(config.TokenValidityHours)*time.Hour, config.SigningSecret)
	twtServer.PostService = posts.New(twtServer.StorageService)
	twtServer.UserService = users.New(twtServer.AuthService, twtServer.StorageService)
	defer twtServer.StorageService.Close()
	twitter.RegisterTwitterServer(s, twtServer)

	if err := s.Serve(listner); err != nil {
		log.Fatalln(err)
	}
}
