package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/jaxsax/projects/tapeworm/botv2/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var (
	hostname = flag.String("hostname", "", "postgres host")
	database = flag.String("database", "", "postgres database name")
	username = flag.String("user", "postgres", "postgres username")
	password = flag.String("password", "", "postgres database")

	sqlitePath = flag.String("sqlite-path", "bot.db", "Export to this sqlite database")
	doInsert   = flag.Bool("insert", false, "really perform insertion")
)

type link struct {
	ID        int64 `db:"id"`
	CreatedTS int64 `db:"created_ts"`
	CreatedBy int64 `db:"created_by"`
	Link      string
	Title     string
	ExtraData string `db:"extra_data"`
}

func main() {
	flag.Parse()

	// Retrieve postgres connection
	pgDB, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable",
			*username, *password, *hostname, *database,
		),
	)

	if err != nil {
		panic(err)
	}
	defer pgDB.Close()

	// Retrieve sqlite connection
	sqliteDB, err := sql.Open("sqlite3", *sqlitePath)
	if err != nil {
		panic(err)
	}
	defer sqliteDB.Close()

	// Read links from postgres
	var links []link
	err = pgDB.Select(&links, "SELECT * FROM links ORDER BY created_ts desc")
	if err != nil {
		panic(err)
	}

	// Write to sqlite

	schema := `
	CREATE TABLE IF NOT EXISTS links  (
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		created_ts BIGINT NOT NULL,
		created_by BIGINT NOT NULL,
		link VARCHAR(1024) NOT NULL,
		title VARCHAR(1024) NOT NULL,
		extra_data TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS skipped_links  (
		id BIGINT PRIMARY KEY NOT NULL,
		error_reason VARCHAR(512) NOT NULL,
		link VARCHAR(1024) NOT NULL
	);

	CREATE TABLE  IF NOT EXISTS updates (
		id BIGINT PRIMARY KEY NOT NULL,
	   "data" TEXT NOT NULL
	);
	`

	_, err = sqliteDB.Exec(schema)
	if err != nil {
		panic(err)
	}

	// Find which IDs already exist
	existingLinks, err := models.Links().All(context.TODO(), sqliteDB)
	if err != nil {
		panic(err)
	}

	existingLinkIDs := map[int64]bool{}
	for _, link := range existingLinks {
		existingLinkIDs[link.ID] = true
	}

	fmt.Printf("%#v\n", existingLinkIDs)

	postgresLinks, err := models.Links().All(context.TODO(), pgDB)
	if err != nil {
		panic(err)
	}

	tx, err := sqliteDB.Begin()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()

			return
		}

		_ = tx.Commit()
	}()

	for _, link := range postgresLinks {
		_, exists := existingLinkIDs[link.ID]
		if exists {
			continue
		}

		var m models.Link
		m.ID = link.ID
		m.CreatedBy = link.CreatedBy
		m.CreatedTS = link.CreatedTS
		m.Link = link.Link
		m.Title = link.Title
		m.ExtraData = link.ExtraData

		if *doInsert {
			err = m.Insert(context.TODO(), tx, boil.Infer())
			if err != nil {
				return
			}
		} else {
			spew.Printf("Would have inserted %v\n", m)
		}
	}
}
