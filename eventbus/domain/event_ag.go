package domain

type EventAG struct {
	domainEvents []*DomainEvent
}

func (eag *EventAG) RecordDomainEvent(event *DomainEvent) {
	if eag.domainEvents == nil {
		eag.domainEvents = make([]*DomainEvent, 0, 1)
	}

	eag.domainEvents = append(eag.domainEvents, event)
}

func (eag *EventAG) PullDomainEvents() []*DomainEvent {
	return eag.domainEvents
}

func (eag *EventAG) ClearDomainEvents() {
	eag.domainEvents = make([]*DomainEvent, 0)
}
