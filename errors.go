package begger

type Error struct {
	StatusCode int
	Status     string
	Message    string
	Metadata   map[string]interface{}
}
