package domain

type Subscriber interface {
	Consume(event []byte)
}
