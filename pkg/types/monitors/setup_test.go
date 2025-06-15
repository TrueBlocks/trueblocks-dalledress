package monitors

import (
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

// MockSDKMonitor provides controlled SDK responses for testing
type MockSDKMonitor struct {
	monitors []coreTypes.Monitor
	count    int
	err      error
	delay    time.Duration
}

func NewMockSDKMonitor() *MockSDKMonitor {
	return &MockSDKMonitor{
		monitors: []coreTypes.Monitor{
			{
				Address:     base.HexToAddress("0x1234567890123456789012345678901234567890"),
				Name:        "Test Monitor 1",
				NRecords:    100,
				FileSize:    1024,
				LastScanned: 12345,
				IsEmpty:     false,
				IsStaged:    true,
				Deleted:     false,
			},
			{
				Address:     base.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
				Name:        "Test Monitor 2",
				NRecords:    250,
				FileSize:    2048,
				LastScanned: 67890,
				IsEmpty:     true,
				IsStaged:    false,
				Deleted:     false,
			},
		},
		count: 2,
		err:   nil,
		delay: 0,
	}
}

func (m *MockSDKMonitor) SetError(err error) {
	m.err = err
}

func (m *MockSDKMonitor) SetDelay(delay time.Duration) {
	m.delay = delay
}

func (m *MockSDKMonitor) SetMonitors(monitors []coreTypes.Monitor) {
	m.monitors = monitors
	m.count = len(monitors)
}

// MessageCapture helps capture messages emitted during tests
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
