package sql

import (
	"github.com/jaxsax/projects/tapeworm/botv2/skippedlinks"
	"github.com/jmoiron/sqlx"
)

type skippedLinksRepository struct {
	db *sqlx.DB
}

func NewSkippedLinksRepository(db *sqlx.DB) *skippedLinksRepository {
	return &skippedLinksRepository{
		db: db,
	}
}

func (repo *skippedLinksRepository) Create(skippedLink skippedlinks.SkippedLink) error {
	query := `INSERT INTO skipped_links(link, error_reason)
				VALUES(:link, :error_reason)`

	_, err := repo.db.NamedExec(query, skippedLink)
	if err != nil {
		return err
	}

	return nil
}
