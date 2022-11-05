package infrastructure

import (
	"ddd-hexa/eventbus/domain"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

type NATSEventBus struct {
	conn *nats.Conn
}

func NewNATSEventBus(host, port, user, password string) *NATSEventBus {
	url := fmt.Sprintf("nats://%s:%s", host, port)
	nc, err := nats.Connect(url, nats.UserInfo(user, password))
	if err != nil {
		panic(err)
	}
	return &NATSEventBus{
		conn: nc,
	}
}

func (eb *NATSEventBus) Publish(events ...domain.DomainEvent) error {
	for i := range events {
		event := events[i]
		data, err := json.Marshal(event)
		if err != nil {
			return err
		}

		err = eb.conn.Publish(event.GetEventTopic(), data)
		if err != nil {
			return err
		}
	}	
}

func (eb *NATSEventBus) RegisterAsyncSubscriber(topic string, subscriber domain.Subscriber) error {
	_, err := eb.conn.Subscribe(topic, func(m *nats.Msg) {
		subscriber.Consume(m.Data)
	})
	if err != nil {
		return err
	}
	return nil
}

func (eb *NATSEventBus) Close() {
	_ = eb.conn.Drain()	
}