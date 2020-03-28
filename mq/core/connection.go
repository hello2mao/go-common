package core

// Connection is an interface that can be used to publishing
type Connection interface {
	OpenQueue(name string) Queue
}
