package core

type Delivery interface {
	Payload() string
	Ack() bool    // consume success
	Reject() bool // consume failed
}
