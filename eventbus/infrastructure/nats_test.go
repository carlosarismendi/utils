package infrastructure

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/carlosarismendi/dddhelper/eventbus/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type event struct {
	*domain.DomainEvent
	Msg string `json:"msg"`
}

type subscriber struct {
	msgs chan string
}

func newSubscriber(ch chan string) *subscriber {
	return &subscriber{
		msgs: ch,
	}
}

func (s *subscriber) Consume(e []byte) {
	var event event
	err := json.Unmarshal(e, &event)
	if err != nil {
		panic(err)
	}

	s.msgs <- event.Msg
}

func TestNatsEventBus(t *testing.T) {
	t.Run("PublishingAnEventThatIsConsumedByASubscriberWithoutError_returnsNoError", func(t *testing.T) {
		const TOPIC = "TEST_TOPIC"

		// Create event bus
		eb := NewNATSEventBus("localhost", "4222")
		defer eb.Close()

		// Subscribe subscriber to topic in event bus
		ch := make(chan string, 3)
		sub := newSubscriber(ch)
		err := eb.RegisterAsyncSubscriber(TOPIC, sub)
		require.NoError(t, err)

		// Publish message to topic
		processID := "ad5898c0-d901-4de2-8123-ba73c8adc190"
		eventID := uuid.New().String()
		event := &event{
			DomainEvent: domain.NewDomainEvent(TOPIC, processID, eventID, "account"),
			Msg:         "MESSAGE 1 FOR TESTING",
		}

		// Check that subscriber has not received any messages before publishing
		select {
		case msg := <-ch:
			t.Error("Unexpected message: ", msg)
		case <-time.After(200 * time.Millisecond):
			// Expected behaviour, since there are no messages published nor received by any subscriber
			break
		}

		err = eb.Publish(event)
		require.NoError(t, err)

		// Read message received by subscriber
		select {
		case msg := <-ch:
			require.Equal(t, "MESSAGE 1 FOR TESTING", msg)
		case <-time.After(250 * time.Millisecond):
			t.Error("Unexpected error: subscriber did not receive the message")
		}
	})
}
