package db

type LinkDimension struct {
	ID     uint64 `db:"id"`
	LinkID uint64 `db:"link_id"`
	Kind   string `db:"kind"`
	Data   string `db:"data"`
}
