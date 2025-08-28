package store

type Documents struct {
	Uuid       int
	Title      string
	Body       string
	Occurrence int
	Score      float64 // pertinence score calculated by bm25() // + score low + pertinent
}

type Store interface {
	Insert(d Documents) error
	Search(query string) ([]Documents, error)
	Close() error
}
