package main

import (
	"context"
	"net/http"

	"github.com/TristanHoladay/devops-deploy-practice/src/backend/api"

	"github.com/rs/zerolog/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var ctx = context.TODO()

func init() {
	log.Info().Msg("Creating database client connection")
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
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

func main() {
	log.Info().Msg("Starting FORM API")

	router := api.BuildRouter(collection, ctx)

	// envConfig := config.ProcessEnvConfig()
	log.Info().Msg("Running FORM API on port 8900")
	err := http.ListenAndServe(":8900", router)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start HTTP server")
	}
}
