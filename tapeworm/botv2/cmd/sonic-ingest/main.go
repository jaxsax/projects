package main

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/expectedsh/go-sonic/sonic"
	"github.com/jaxsax/projects/tapeworm/botv2/links"
	"github.com/jaxsax/projects/tapeworm/botv2/models"

	_ "github.com/mattn/go-sqlite3"
)

var (
	host     = flag.String("sonic-host", "", "host address to sonic server")
	port     = flag.Int("sonic-port", 0, "port number to sonic server")
	password = flag.String("sonic-password", "", "password of sonic server")

	sqlitePath         = flag.String("sqlite-path", "", "path to sqlite db")
	ingest             = flag.Bool("ingest", false, "whether we should ingest into sonic")
	linksCollection    = "links"
	linksByTitleBucket = "title"
)

func main() {
	// This binary demonstrates ingesting data into sonic
	// https://github.com/valeriansaliou/sonic
	flag.Parse()

	if *host == "" {
		log.Println("sonic-host cannot be empty")
		os.Exit(1)
	}

	if *port == 0 {
		*port = 1491
	}

	ingester, err := sonic.NewIngester(*host, *port, *password)
	if err != nil {
		panic(err)
	}

	if *sqlitePath == "" {
		log.Println("sqlite-path cannot be empty")
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", *sqlitePath)
	if err != nil {
		panic(err)
	}

	if *ingest {
		err = ingestFromSqlite(db, ingester)
		if err != nil {
			panic(err)
		}
	}

	search, err := sonic.NewSearch(*host, *port, *password)
	if err != nil {
		panic(err)
	}

	linksRepository := links.NewSqliteRepository(db)

	log.Println("Ready for search!")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		txt := scanner.Text()
		if txt == "quit" {
			break
		}

		txt = strings.Trim(txt, "\n")
		fmt.Printf("Search query: '%v'\n", txt)

		suggestions, err := search.Suggest("links", "title", txt, 10)
		if err == nil {
			fmt.Printf("suggestions: %v\n", suggestions)
		}

		results, err := search.Query(linksCollection, linksByTitleBucket, txt, 25, 0, sonic.LangEng)
		if err != nil {
			fmt.Printf("error searching: %v\n", err)
			continue
		}

		// We can't simply use len(results) because Query internally
		// returns results from a subslice of string[:3], instead for emptyness check
		// We'll see if its len(1) and the first entry is an empty string
		if len(results) == 0 {
			fmt.Printf("No results found\n")
			continue
		}

		if len(results) == 1 {
			if results[0] == "" {
				fmt.Printf("No results found\n")
				continue
			}
		}

		fmt.Printf("Results (%v):\n", len(results))

		matchingLinks, err := convertToLinks(linksRepository, results)
		if err != nil {
			fmt.Printf("error fetching links: %v\n", err)
			continue
		}

		whitespaceCleaner := regexp.MustCompile(`\s+`)
		for _, r := range matchingLinks {
			// See if we can enhance by username in extra_data
			creator := strconv.FormatInt(r.CreatedBy, 10)
			if username, ok := r.ExtraData["created_username"]; ok {
				creator = username
			}

			t := time.Unix(r.CreatedTS, 0)
			ts := humanize.RelTime(t, time.Now(), "ago", "later")
			title := strings.Replace(r.Title, "\n", " ", -1)
			title = whitespaceCleaner.ReplaceAllString(title, " ")
			fmt.Printf("[%v] %v\n", ts, title)
			fmt.Printf("\t%v\n\t%v\n", r.Link, creator)
		}
	}

	if scanner.Err() != nil {
		panic(scanner.Err())
	}
}

func convertToLinks(repository links.Repository, objects []string) ([]links.Link, error) {
	// Convert from id:<int64> -> list of int64 ids
	linkIDs := make([]int64, 0, len(objects))

	for _, obj := range objects {
		p := strings.Split(obj, ":")
		if len(p) != 2 {
			continue
		}

		if p[0] != "id" {
			continue
		}

		id, err := strconv.ParseInt(p[1], 10, 64)
		if err != nil {
			continue
		}

		linkIDs = append(linkIDs, id)
	}

	matchingLinks, err := repository.ListMatchingIDs(linkIDs)
	if err != nil {
		return []links.Link{}, fmt.Errorf("list: %w", err)
	}

	return matchingLinks, nil
}

func ingestFromSqlite(db *sql.DB, ingestable sonic.Ingestable) error {
	links, err := models.Links().All(context.Background(), db)
	if err != nil {
		return fmt.Errorf("find all: %w", err)
	}

	for _, link := range links {
		object := fmt.Sprintf("id:%v", link.ID)
		text := link.Title
		err = ingestable.Push(linksCollection, linksByTitleBucket, object, text, sonic.LangNone)
		if err != nil {
			log.Printf("ingest error: [object: %v, text: %v] err: %v", object, text, err)
			continue
		}

		log.Printf("indexed [%v:%v]", object, text)
	}

	return nil
}
