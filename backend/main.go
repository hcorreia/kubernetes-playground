package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	// _ "github.com/golang-migrate/migrate/v4/source/github"

	"k8s-tests/backend/database"
)

func connectDB() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s?parseTime=true", os.Getenv("DB_STRING")))
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)

	return db
}

func mainMigration(db *sql.DB, direction string) error {

	fmt.Println("Running migrations",
		strings.ToUpper(direction),
		"...")

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./database/migrations",
		"mysql", driver)
	if err != nil {
		return err
	}

	if direction == "status" {
		ver, _, err := m.Version()

		if err == nil {
			fmt.Println("Version: ", ver)
		} else {
			fmt.Println("Version: ", err.Error())
		}

	} else if direction == "up" {
		err = m.Up()
	} else if direction == "down" {
		// err = m.Down()
		err = m.Steps(-1)
	} else {
		panic("Wrong direction!")
	}

	if err != nil {
		fmt.Println("Done:", err.Error())
		// return err
	} else {
		fmt.Println("Done")
	}

	return nil
}

func getPosts(db *sql.DB) ([]database.Post, error) {
	ctx := context.Background()

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

func panicHelpText() {
	panic(
		"Wrong command line args.\n" +
			"E.g.:\n" +
			"  - backend migrate up\n" +
			"  - backend migrate down\n")
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

	if len(os.Args) > 1 {
		if len(os.Args) == 3 &&
			os.Args[1] == "migrate" &&
			slices.Contains([]string{"up", "down", "status"}, os.Args[2]) {

			err = mainMigration(connectDB(), os.Args[2])
			if err != nil {
				panic(err)
			}
			return
		} else {
			panicHelpText()
		}
	}

	db := connectDB()

	fmt.Println("Lets GO...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from POD: %s", hostname)
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := getPosts(db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}

		fmt.Fprint(w, "Ok")
	})
	http.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {

		posts, err := getPosts(db)
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

	port := ":8080"

	fmt.Println("Running on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
