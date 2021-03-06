// Package skippedlinks contains the domain model for a skipped link
package skippedlinks

// SkippedLink is the domain model
type SkippedLink struct {
	ID          int64
	Link        string
	ErrorReason string `db:"error_reason"`
}

// Repository provides access to a links store
type Repository interface {
	Create(skippedLink SkippedLink) error
}
