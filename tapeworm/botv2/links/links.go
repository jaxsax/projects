// Package links contains the domain model for a link
package links

// Link is the domain model
type Link struct {
	ID        int64
	CreatedTS int64
	CreatedBy int64
	Link      string
	Title     string
	ExtraData map[string]string
}

// Repository provides access to a links store
type Repository interface {
	CreateMany(link []Link) error
	ListMatchingIDs(IDs []int64) ([]Link, error)
	List() ([]Link, error)
}
