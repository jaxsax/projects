package botv2

import "github.com/jmoiron/sqlx"

type Link struct {
	ID        int64
	CreatedTS int64
	CreatedBy string
	Link      string
	Title     string
	ExtraData string
}

type LinksDB struct {
	db *sqlx.DB
}

func NewLinksDB(db *sqlx.DB) *LinksDB {
	return &LinksDB{
		db: db,
	}
}
