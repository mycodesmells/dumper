package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Post represents a structure of a blog post.
type Post struct {
	ID       int
	Title    string
	Slug     string
	Abstract string
	Content  string
	Dates    struct {
		Create      time.Time
		Publication time.Time
	}
}

var (
	username = flag.String("username", "", "username to access database")
	password = flag.String("password", "", "password to access database")
	host     = flag.String("host", "", "host (address) of the location of database")
	dbName   = flag.String("db", "", "name of database")
	filename = flag.String("out", "/tmp/dumper.json", "location of output JSON file")
)

func main() {
	flag.Parse()

	// Set up connection to MySQL database
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", *username, *password, *host, *dbName)
	db, err := sql.Open("mysql", connectionString)
	defer db.Close()
	if err != nil {
		panic(fmt.Errorf("failed to set up connection with MySQL database: %+v", err))
	}

	res, err := db.Query("SELECT id, title, slug, abstract, content, create_date, publication_date FROM posts_post")
	if err != nil {
		panic(fmt.Errorf("failed to fetch data: %+v", err))
	}
	defer res.Close()

	f, err := os.Create(*filename)
	defer f.Close()

	fmt.Println("Parsing and saving posts...")
	var p Post
	for res.Next() {
		res.Scan(&p.ID, &p.Title, &p.Slug, &p.Abstract, &p.Content, &p.Dates.Create, &p.Dates.Publication)

		err = savePostToFile(*f, p)
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("\nDone - Posts saved as JSON to file: %s\n", *filename)
}

func savePostToFile(f os.File, p Post) error {
	fmt.Printf(".%d.", p.ID)

	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	f.Write(b)
	f.WriteString("\n")
	return nil
}
