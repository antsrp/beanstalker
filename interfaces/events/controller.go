package events

type Controller interface {
	Parse(string) (*string, error)
}
