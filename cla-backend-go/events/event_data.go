package events

type EventData interface {
	GetEventString(args *LogEventArgs) (eventData string, containsPII bool)
}

type CreateUserEventData struct{}

func (ed *CreateUserEventData) GetEventString(args *LogEventArgs) (eventData string, containsPII bool) {
	return CreateUser, true
}
