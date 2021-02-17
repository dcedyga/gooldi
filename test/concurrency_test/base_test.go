package concurrency_test

import (
	"fmt"
	concurrency "go-channels/concurrency"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
}

//
func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

//
func (suite *Suite) SetupSuite() {}

//
func (suite *Suite) SetupTest() {}

//
func (suite *Suite) TearDownSuite() {

}

func process1(input interface{}, params ...interface{}) interface{} {
	event := input.(*concurrency.Event)
	id := params[0]
	e := &concurrency.Message{
		ID:         uuid.NewV4().String(),
		Message:    fmt.Sprintf("Client %d got message: %v", id, event.OutMessage.Message),
		TimeInNano: event.OutMessage.TimeInNano,
		MsgType:    event.OutMessage.MsgType,
	}
	return e
}
func process2(input interface{}, params ...interface{}) interface{} {
	event := input.(*concurrency.Event)
	id := 2
	i := event.OutMessage.Message.(int) * 2
	e := &concurrency.Message{
		ID:         uuid.NewV4().String(),
		Message:    fmt.Sprintf("This is a different client that multiplies by 2 -> ClientId %v,  message: %v, value: %v", id, event.OutMessage.Message, i),
		TimeInNano: event.OutMessage.TimeInNano,
		MsgType:    event.OutMessage.MsgType,
	}
	return e

}
func process3(input interface{}, params ...interface{}) interface{} {
	event := input.(*concurrency.Event)
	id := params[0]
	i := event.OutMessage.Message.(int) * event.OutMessage.Message.(int)
	e := &concurrency.Message{
		ID:         uuid.NewV4().String(),
		Message:    fmt.Sprintf("This is the final client that squares the message -> ClientId %v,  message: %v, value: %v", id, event.OutMessage.Message, i),
		TimeInNano: event.OutMessage.TimeInNano,
		MsgType:    event.OutMessage.MsgType,
	}
	return e
}
func process4(input interface{}, params ...interface{}) interface{} {
	event := input.(*concurrency.Event)
	id := params[0]
	i := event.OutMessage.Message.(int) * 3
	e := &concurrency.Message{
		ID:         uuid.NewV4().String(),
		Message:    fmt.Sprintf("This is a different client that multiplies by 3 -> ClientId %v,  message: %v, value: %v", id, event.OutMessage.Message, i),
		TimeInNano: event.OutMessage.TimeInNano,
		MsgType:    event.OutMessage.MsgType,
	}
	return e
}

func filterFn(input interface{}, params ...interface{}) bool {
	event := input.(*concurrency.Event)
	i := event.InitMessage.Message.(int) % 2
	return i == 0

}
