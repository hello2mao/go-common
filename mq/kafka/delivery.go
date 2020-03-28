package kafka

type Delivery struct {
}

func (d *Delivery) Payload() string {
	return ""
}
func (d *Delivery) Ack() bool {
	return true
}

func (d *Delivery) Reject() bool {
	return true
}
