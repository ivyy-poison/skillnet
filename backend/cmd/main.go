package main

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

type serverConfig struct {
	db          *gorm.DB
	store       redis.Store
	router      *gin.Engine
	googleCloud *storage.Client
}

func main() {
	serverConfig := initialiseProdServer()
	serverConfig.setupRoutes()
	serverConfig.runRouter()
	log.Println("Setup complete!")
}

func initialiseProdServer() *serverConfig {
	router := gin.Default()
	db := database.ConnectProdDatabase()
	server := serverConfig{
		router: router,
		db:     db,
	}
	server.setupRedis()
	server.setupGoogleCloud()
	return &server
}

func (server *serverConfig) setupRedis() {
	env := helpers.RetrieveRedisEnv()
	redisAddress := env.Address()
	const redisNetwork = "tcp"
	store, err := redis.NewStore(env.MaxConn, redisNetwork, redisAddress, "", []byte(env.Secret))
	if err != nil {
		log.Fatal(err.Error())
	}
	server.store = store
}

func (s *serverConfig) setupGoogleCloud() {
	ctx := context.Background()
	env := helpers.RetrieveGoogleCloudEnv()

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(env.Filepath))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	s.googleCloud = client
}

func (server *serverConfig) runRouter() {
	env := helpers.RetrieveWebAppEnv()
	routerAddress := env.Address()
	err := server.router.Run(routerAddress)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
}
