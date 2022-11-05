package domain

type IDomainEvent interface {
	GetProcessID() string
	GetEventID() string
	GetEventTopic() string
	GetAccount() string
}

type DomainEvent struct {
	ProcessID  string
	EventID    string
	EventTopic string
	Account    string
}

func (e *DomainEvent) GetProcessID() string {
	return e.ProcessID
}

func (e *DomainEvent) GetEventID() string {
	return e.EventID
}

func (e *DomainEvent) GetEventTopic() string {
	return e.EventTopic
}

func (e *DomainEvent) GetAccount() string {
	return e.Account
}
