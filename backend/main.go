package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"k8s-tests/backend/database"
)

func getPosts() ([]database.Post, error) {
	// ! I Know...
	// ! This ia NOT right. DB connection should not be here.

	ctx := context.Background()

	db, err := sql.Open("mysql", fmt.Sprintf("%s?parseTime=true", os.Getenv("DB_STRING")))
	if err != nil {
		return nil, err
	}

	queries := database.New(db)

	// list all posts
	posts, err := queries.ListPosts(ctx)
	if err != nil {
		return nil, err
	}
	// log.Println(posts)

	return posts, nil
}

type DataResultMeta struct {
	Timestamp string `json:"timestamp"`
	Hostname  string `json:"hostname"`
}

type DataResult[T any] struct {
	Data T              `json:"data"`
	Meta DataResultMeta `json:"meta"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file.")
	}

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	fmt.Println("Lets GO...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from POD: %s", hostname)
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := getPosts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}

		fmt.Fprint(w, "Ok")
	})
	http.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {

		posts, err := getPosts()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}

		result, err := json.Marshal(DataResult[[]database.Post]{
			Data: posts,
			Meta: DataResultMeta{
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Hostname:  hostname,
			},
		})

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, string(result))
	})
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Auth from POD: %s", hostname)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
