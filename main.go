package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rachit2501/goserver/internal/database"
);

type apiConfig struct {
	DB *database.Queries
}

func main(){
	godotenv.Load();

	portString := os.Getenv("PORT")
	
	if portString == "" {
		log.Fatal("PORT is not found in the environmaent")
	}

	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("Can't connect to the server",err)
	}	

	db := database.New(conn)
	apiCFG := apiConfig{
		DB: db,
	}

	go startScraping(db,10, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*","https://*"},
		AllowedMethods: []string{"GET","POST","PUT","DELETE","OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))

	v1Router := chi.NewRouter();

	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerError)
	
	v1Router.Post("/users", apiCFG.handlerCreateUser)
	v1Router.Get("/users", apiCFG.middlewareAuth(apiCFG.handlerGetUser))
	v1Router.Post("/feeds", apiCFG.middlewareAuth(apiCFG.handlerFeedCreate))
	
	v1Router.Get("/feed_follows", apiCFG.middlewareAuth(apiCFG.handlerFeedFollowsGet))
	v1Router.Post("/feed_follows", apiCFG.middlewareAuth(apiCFG.handlerFeedFollowCreate))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCFG.middlewareAuth(apiCFG.handlerFeedFollowDelete))

	v1Router.Get("/posts", apiCFG.middlewareAuth(apiCFG.handlerPostsGet))


	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:   ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	err = srv.ListenAndServe()

	if err!=nil {
		log.Fatal(err)
	}

}