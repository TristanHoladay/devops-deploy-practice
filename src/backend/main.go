package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/TristanHoladay/devops-deploy-practice/src/backend/api"

	"github.com/rs/zerolog/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/joho/godotenv"
)

var collection *mongo.Collection
var ctx = context.TODO()

func loadEnvConfig() {
	log.Info().Msg("Loading env config")
	err := godotenv.Load("env/.env")
	if err != nil {
		log.Err(err).Msg("failed to get env file")
	}

}

func createDBConn() {
	log.Info().Msg("Creating database client connection")
	mongoURL := fmt.Sprintf("mongodb://%s:%s/", os.Getenv("MONGO_SERVER"), os.Getenv("MONGO_PORT"))
	clientOptions := options.Client().ApplyURI(mongoURL)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to make mongodb connection")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot ping mongodb server")
	}

	log.Info().Msg("Creating collection users in database: devops_deploy")
	collection = client.Database("devops_deploy").Collection("users")
}

func init() {
	loadEnvConfig()
	createDBConn()
}

func main() {
	log.Info().Msg("Starting FORM API")
	router := api.BuildRouter(collection, ctx)

	log.Info().Msg("Running FORM API on port 8900")
	err := http.ListenAndServe(":8900", router)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start HTTP server")
	}
}
