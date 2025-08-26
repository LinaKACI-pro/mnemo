package store

type Entry struct {
	ID    int
	Value string
}

type Store interface {
	Insert(e Entry) error
	Search(query string) ([]Entry, error)
	Close() error
}
