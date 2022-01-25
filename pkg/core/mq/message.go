package mq

import (
	"context"
	"sync"
)

var closedChan = make(chan struct{})

func init() {
	close(closedChan)
}

type ackType int

const (
	noAckSent ackType = iota
	ack
	nack
)

// Message 消息
type Message struct {
	Header map[string]string
	Body   []byte

	ctx  context.Context
	ack  chan struct{}
	nack chan struct{}

	ackMutex    sync.Mutex
	ackSentType ackType
}

// NewMessage
// 	@Description 创建 消息
//	@Param body
// 	@Return *Message
func NewMessage(body []byte) *Message {
	return &Message{
		Body:   body,
		Header: make(map[string]string),
		ack:    make(chan struct{}),
		nack:   make(chan struct{}),
	}
}

// Ack
// 	@Description 确认消息，幂等，
// 	@Receiver m Message
// 	@Return bool 如果已经nack 了，为false
func (m *Message) Ack() bool {
	m.ackMutex.Lock()
	defer m.ackMutex.Unlock()

	if m.ackSentType == nack {
		return false
	}
	if m.ackSentType != noAckSent {
		return true
	}

	m.ackSentType = ack
	if m.ack == nil {
		m.ack = closedChan
	} else {
		close(m.ack)
	}

	return true
}

// NAck
// 	@Description 不确认消息，幂等
// 	@Receiver m
// 	@Return bool 如果已经ack 了为false
func (m *Message) NAck() bool {
	m.ackMutex.Lock()
	defer m.ackMutex.Unlock()

	if m.ackSentType == ack {
		return false
	}
	if m.ackSentType != noAckSent {
		return true
	}

	m.ackSentType = nack

	if m.nack == nil {
		m.nack = closedChan
	} else {
		close(m.nack)
	}

	return true
}

// Acked
// 	@Description
// 	@Receiver m
// 	@Return <-chan
func (m *Message) Acked() <-chan struct{} {
	return m.ack
}

// NAcked
// 	@Description
// 	@Receiver m
// 	@Return <-chan
func (m *Message) NAcked() <-chan struct{} {
	return m.nack
}

// Context
// 	@Description
// 	@Receiver m
// 	@Return context.Context
func (m *Message) Context() context.Context {
	if m.ctx != nil {
		return m.ctx
	}
	return context.Background()
}

// SetContext
// 	@Description
// 	@Receiver m
//	@Param ctx
func (m *Message) SetContext(ctx context.Context) {
	m.ctx = ctx
}

// Copy
// 	@Description
// 	@Receiver m
// 	@Return *Message
func (m *Message) Copy() *Message {
	msg := NewMessage(m.Body)
	for k, v := range m.Header {
		msg.Header[k] = v
	}
	return msg
}
