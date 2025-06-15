package testing

import (
	"sync"
	"time"
)

type MessageCapture struct {
	messages []CapturedMessage
	mu       sync.Mutex
}

type CapturedMessage struct {
	EventType string
	MsgText   string
	Payload   []interface{}
	Timestamp time.Time
}

func NewMessageCapture() *MessageCapture {
	return &MessageCapture{
		messages: make([]CapturedMessage, 0),
	}
}

func (mc *MessageCapture) CaptureLoaded(args ...interface{}) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	msgText := ""
	payload := args
	if len(args) > 0 {
		if str, ok := args[0].(string); ok {
			msgText = str
			payload = args[1:]
		}
	}

	mc.messages = append(mc.messages, CapturedMessage{
		EventType: "DataLoaded",
		MsgText:   msgText,
		Payload:   payload,
		Timestamp: time.Now(),
	})
}

func (mc *MessageCapture) GetMessages() []CapturedMessage {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	result := make([]CapturedMessage, len(mc.messages))
	copy(result, mc.messages)
	return result
}

func (mc *MessageCapture) GetMessagesForText(msgText string) []CapturedMessage {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	var result []CapturedMessage
	for _, msg := range mc.messages {
		if msg.MsgText == msgText {
			result = append(result, msg)
		}
	}
	return result
}

func (mc *MessageCapture) Reset() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.messages = mc.messages[:0]
}
