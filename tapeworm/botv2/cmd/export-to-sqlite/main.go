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

	sqlitePath               = flag.String("sqlite-path", "bot.db", "Export to this sqlite database")
	reallyInsertLinks        = flag.Bool("insert-links", false, "really insert links")
	reallyInsertUpdates      = flag.Bool("insert-updates", false, "really insert updates")
	reallyInsertSkippedLinks = flag.Bool("insert-skipped", false, "really insert skipped links")
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

	if syncLinks(pgDB, sqliteDB); err != nil {
		panic(err)
	}

	if syncMessages(pgDB, sqliteDB); err != nil {
		panic(err)
	}

	if syncSkippedLinks(pgDB, sqliteDB); err != nil {
		panic(err)
	}
}

func doInTx(f func(tx *sql.Tx) error, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("open tx: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	err = f(tx)
	return err
}

func syncLinks(pgDB *sqlx.DB, sqliteDB *sql.DB) error {
	// Find which IDs already exist
	existingLinks, err := models.Links().All(context.TODO(), sqliteDB)
	if err != nil {
		return fmt.Errorf("find all in sqlite: %w", err)
	}

	existingLinkIDs := map[int64]bool{}
	for _, link := range existingLinks {
		existingLinkIDs[link.ID] = true
	}

	postgresLinks, err := models.Links().All(context.TODO(), pgDB)
	if err != nil {
		return fmt.Errorf("find all in postgres: %w", err)
	}
	fmt.Printf(
		"links|existing=%v,total=%v,new=%v\n",
		len(existingLinkIDs),
		len(postgresLinks),
		len(postgresLinks)-len(existingLinkIDs),
	)

	tx, err := sqliteDB.Begin()
	if err != nil {
		return fmt.Errorf("create tx: %w", err)
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

		if *reallyInsertLinks {
			err = m.Insert(context.TODO(), tx, boil.Infer())
			if err != nil {
				spew.Printf("%v insert_error: err: %v\n", m, err)
				return nil
			}
		} else {
			spew.Printf("Would have inserted %v\n", m)
		}
	}

	return nil
}

func syncMessages(pgDB *sqlx.DB, sqliteDB *sql.DB) error {
	allUpdates, err := models.Updates().All(context.TODO(), sqliteDB)
	if err != nil {
		return fmt.Errorf("find existing updates: %w", err)
	}

	existingUpdateIDs := map[int64]bool{}
	for _, update := range allUpdates {
		existingUpdateIDs[update.ID] = true
	}

	postgresUpdates, err := models.Updates().All(context.TODO(), pgDB)
	if err != nil {
		return fmt.Errorf("find all pg updates: %w", err)
	}
	fmt.Printf(
		"updates|existing=%v,total=%v,new=%v\n",
		len(existingUpdateIDs),
		len(postgresUpdates),
		len(postgresUpdates)-len(existingUpdateIDs),
	)

	err = doInTx(func(tx *sql.Tx) error {
		for _, update := range postgresUpdates {
			_, exists := existingUpdateIDs[update.ID]
			if exists {
				continue
			}

			var m models.Update
			m.ID = update.ID
			m.Data = update.Data

			if *reallyInsertUpdates {
				err = m.Insert(context.TODO(), tx, boil.Infer())
				if err != nil {
					spew.Printf("%v insert_update_error, err: %v\n", m, err)
					return err
				}
			} else {
				spew.Printf("Would have inserted update: %v\n", m)
			}
		}

		return nil
	}, sqliteDB)
	if err != nil {
		return fmt.Errorf("sync updates: %w", err)
	}

	return nil
}

func syncSkippedLinks(pgDB *sqlx.DB, sqliteDB *sql.DB) error {
	allSkipped, err := models.SkippedLinks().All(context.TODO(), sqliteDB)
	if err != nil {
		return fmt.Errorf("find existing skipped links: %w", err)
	}

	existingIDs := map[int64]bool{}
	for _, skipped := range allSkipped {
		existingIDs[skipped.ID] = true
	}

	postgresSkipped, err := models.SkippedLinks().All(context.TODO(), pgDB)
	if err != nil {
		return fmt.Errorf("find all pg updates: %w", err)
	}

	fmt.Printf(
		"skipped_link|existing=%v,total=%v,new=%v\n",
		len(existingIDs),
		len(postgresSkipped),
		len(postgresSkipped)-len(existingIDs),
	)

	err = doInTx(func(tx *sql.Tx) error {
		for _, skippedLink := range postgresSkipped {
			_, exists := existingIDs[skippedLink.ID]
			if exists {
				continue
			}

			var m models.SkippedLink
			m.ID = skippedLink.ID
			m.Link = skippedLink.Link
			m.ErrorReason = skippedLink.ErrorReason

			if *reallyInsertSkippedLinks {
				err = m.Insert(context.TODO(), tx, boil.Infer())
				if err != nil {
					spew.Printf("%v insert_skipped_link_error, err: %v\n", m, err)
					return err
				}
			} else {
				spew.Printf("Would have inserted skipped link: %v\n", m)
			}
		}
		return nil
	}, sqliteDB)
	if err != nil {
		return fmt.Errorf("sync skipped links: %w", err)
	}

	return nil
}
