package api

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection
var ctx context.Context

type User struct {
	firstName string
	lastName  string
}

func BuildRouter(coll *mongo.Collection, c context.Context) chi.Router {
	log.Debug().Msg("Starting API")
	collection = coll
	ctx = c

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	setUpRoutes(router)

	return router
}

func setUpRoutes(router chi.Router) {
	log.Debug().Msg("Registering API routes")

	// Catch all routes
	router.Get("/*", handleCatchAllGet)
	router.Post("/form-submit", SubmitForm)
}

func createFileServer() http.Handler {
	log.Debug().Msg("Registering static file server, to serve the frontend")
	return http.FileServer(http.Dir("./ui"))
}

func handleCatchAllGet(w http.ResponseWriter, r *http.Request) {
	// If the request is not a real file, serve the index.html instead
	path := filepath.Join("./ui", strings.TrimPrefix(r.URL.Path, "/"))
	fs := createFileServer()
	if _, err := os.Stat(path); err != nil {
		r.URL.Path = "/"
	}
	fs.ServeHTTP(w, r)
}

func SubmitForm(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("submitting form data")
	r.ParseForm()
	user := User{firstName: r.FormValue("first-name"), lastName: r.FormValue("last-name")}
	log.Info().Msg("saving: " + r.FormValue("first-name") + " " + r.FormValue("first-name"))

	_, err := collection.InsertOne(ctx, bson.D{{"firstName", user.firstName}, {"lastName", user.lastName}})
	if err != nil {
		log.Err(err).Msg("Failed to save user to the database")
	}
}
