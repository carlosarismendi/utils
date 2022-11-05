package domain

type EventBus interface {
	Publish(...DomainEvent) error
} 