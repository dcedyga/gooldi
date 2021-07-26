package concurrency

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// MessageOption - option to initialize the message
type MessageOption func(*Message)

// MessagePairOption - option to initialize the message pair
type MessagePairOption func(*MessagePair)

/*
Message - Struct that represents an message in the context of the concurrency package.
Contains:
 - ID of the message
 - the Message,
 - the time that was produced
 - the type of the message,
 - a correlation key
 - the index related to the order that can be used to produce deterministic outputs
*/
type Message struct {
	ID             string
	Message        interface{}
	TimeInNano     int64
	MsgType        string
	CorrelationKey int64
	Index          int64
}

func IsMessagePtr(t interface{}) bool {
	switch t.(type) {
	case *Message:
		return true
	default:
		return false
	}
}

// MessagePair - Struct that represents an input/output message Pair in the context of the concurrency package. Contains an In Message,
// an Out Message, and a index number that is usefull to define the processing order of the MessagePair and a correlationKey
type MessagePair struct {
	In             *Message
	Out            *Message
	Index          int64
	CorrelationKey int64
}

//Message and MessagePair constructors
// NewMessage - Constructor
func NewMessage(msg interface{}, mType string, opts ...MessageOption) *Message {
	m := &Message{
		ID:             uuid.NewV4().String(),
		TimeInNano:     time.Now().UnixNano(),
		Message:        msg,
		MsgType:        mType,
		CorrelationKey: time.Now().UnixNano(),
		Index:          0,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

//MessageWithIndex - initialize index
func MessageWithIndex(idx int64) MessageOption {
	return func(m *Message) {
		m.Index = idx
	}
}

//MessageWithCorrelationKey - initialize correlationKey
func MessageWithCorrelationKey(cKey int64) MessageOption {
	return func(m *Message) {
		m.CorrelationKey = cKey
	}
}

// NewMessagePair - Constructor
func NewMessagePair(in *Message, out *Message, opts ...MessagePairOption) *MessagePair {
	m := &MessagePair{
		In:             in,
		Out:            out,
		CorrelationKey: in.CorrelationKey,
		Index:          out.Index,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

//MessagePairWithIndex - initialize index
func MessagePairWithIndex(idx int64) MessagePairOption {
	return func(m *MessagePair) {
		m.Index = idx
	}
}

//MessagePairWithCorrelationKey - initialize correlationKey
func MessagePairWithCorrelationKey(cKey int64) MessagePairOption {
	return func(m *MessagePair) {
		m.CorrelationKey = cKey
	}
}
