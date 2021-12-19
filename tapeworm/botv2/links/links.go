// Package links contains the domain model for a link
package links

// Link is the domain model
type Link struct {
	ID        int64             `json:"id"`
	CreatedTS int64             `json:"created_ts"`
	CreatedBy int64             `json:"created_by"`
	DeletedAt int64             `json:"deleted_at"`
	Link      string            `json:"link"`
	Title     string            `json:"title"`
	ExtraData map[string]string `json:"extra_data"`
}

// Repository provides access to a links store
type Repository interface {
	CreateMany(link []Link) error
	ListMatchingIDs(IDs []int64) ([]Link, error)
	List() ([]Link, error)
}
