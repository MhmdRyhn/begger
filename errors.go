package begger

type Error struct {
	HTTPStatusCode int
	StatusName     string
	Message        string
	Metadata       map[string]interface{}
}
