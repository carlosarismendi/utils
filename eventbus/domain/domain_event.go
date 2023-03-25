package domain

type IDomainEvent interface {
	GetTopic() string
	GetProcessID() string
	GetEventID() string
	GetAccount() string
}

// nolint:revive // I want this to be called DomainEvent because it's a Domain Driven Design concept.
type DomainEvent struct {
	Topic     string `json:"topic"`
	ProcessID string `json:"processID"`
	EventID   string `json:"eventID"`
	Account   string `json:"account"`
}

func NewDomainEvent(topic, processID, eventID, account string) *DomainEvent {
	return &DomainEvent{
		Topic:     topic,
		ProcessID: processID,
		EventID:   eventID,
		Account:   account,
	}
}

func (e *DomainEvent) GetTopic() string {
	return e.Topic
}

func (e *DomainEvent) GetProcessID() string {
	return e.ProcessID
}

func (e *DomainEvent) GetEventID() string {
	return e.EventID
}

func (e *DomainEvent) GetAccount() string {
	return e.Account
}
