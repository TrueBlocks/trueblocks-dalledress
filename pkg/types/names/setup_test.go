package names

import (
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	testingPkg "github.com/TrueBlocks/trueblocks-dalledress/pkg/testing"
)

type MockSDKName struct {
	*testingPkg.MockSDKBase[coreTypes.Name]
}

func NewMockSDKName() *MockSDKName {
	return &MockSDKName{
		MockSDKBase: testingPkg.NewMockSDKBase(testingPkg.CreateTestNames()),
	}
}

func (m *MockSDKName) SetNames(names []coreTypes.Name) {
	m.SetItems(names)
}

type MessageCapture = testingPkg.MessageCapture
type CapturedMessage = testingPkg.CapturedMessage

func NewMessageCapture() *MessageCapture {
	return testingPkg.NewMessageCapture()
}
