package concurrency

// Message - Struct that represents an message in the context of the concurrency package. Contains the ID of the
// message, the Message, the time that was produced and the type of the message
type Message struct {
	ID         string
	Message    interface{}
	TimeInNano int64
	MsgType    string
}

// Event - Struct that represents an event in the context of the concurrency package. Contains an InitMessage,
// An InMessageSequence to give traceability of the Event, an OutMessage, and a sequence number that is usefull to
// define the processing order of the event
type Event struct {
	InitMessage       *Message
	InMessageSequence Slice //Slice of Messages - we can keep track of the entire flow
	OutMessage        *Message
	Sequence          interface{}
}

// EventMap - A concurrency.Map of events with a common InitMessage
type EventMap struct {
	InitMessage *Message
	Events      *Map
}

// EventSortedMap - A concurrency.SortedMap of events with a common InitMessage
type EventSortedMap struct {
	InitMessage *Message
	Events      *SortedMap
}
