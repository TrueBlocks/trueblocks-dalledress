package names

import (
	testingPkg "github.com/TrueBlocks/trueblocks-dalledress/pkg/testing"
)

type MockSDKName struct {
	*testingPkg.MockSDKBase[Name]
}

func NewMockSDKName() *MockSDKName {
	return &MockSDKName{
		MockSDKBase: testingPkg.NewMockSDKBase(testingPkg.CreateTestNames()),
	}
}

func (m *MockSDKName) SetNames(names []Name) {
	m.SetItems(names)
}

type MessageCapture = testingPkg.MessageCapture
type CapturedMessage = testingPkg.CapturedMessage

func NewMessageCapture() *MessageCapture {
	return testingPkg.NewMessageCapture()
}
