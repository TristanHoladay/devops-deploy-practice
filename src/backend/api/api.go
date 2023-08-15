package api

import (
	"context"
	"fmt"
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

	log.Debug().Msg("Registering API routes")

	// router.Route("/", func(r chi.Router) {

	// })
	router.Post("/form-submit", SubmitForm)

	log.Debug().Msg("Registering static file server, to serve the frontend")

	// Setup a file server for the static UI files
	fs := http.FileServer(http.Dir("./build/ui"))

	// Catch all routes
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// If the request is not a real file, serve the index.html instead
		path := filepath.Join("./build/ui", strings.TrimPrefix(r.URL.Path, "/"))
		if _, err := os.Stat(path); err != nil {
			r.URL.Path = "/"
		}
		fs.ServeHTTP(w, r)
	})

	return router
}

func SubmitForm(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("submitting form data")
	r.ParseForm()
	fmt.Println(r.FormValue("first-name"))
	fmt.Println(r.FormValue("last-name"))
	user := User{firstName: r.FormValue("first-name"), lastName: r.FormValue("last-name")}

	// fmt.Println(user)

	_, err := collection.InsertOne(ctx, bson.D{{"firstName", user.firstName}, {"lastName", user.lastName}})
	if err != nil {
		log.Err(err).Msg("Failed to save user to the database")
	}
}
