package abis

import (
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	testingPkg "github.com/TrueBlocks/trueblocks-dalledress/pkg/testing"
)

type MockSDKAbi struct {
	*testingPkg.MockSDKBase[coreTypes.Abi]
	functions []coreTypes.Function
}

func NewMockSDKAbi() *MockSDKAbi {
	return &MockSDKAbi{
		MockSDKBase: testingPkg.NewMockSDKBase(testingPkg.CreateTestAbis()),
		functions:   testingPkg.CreateTestFunctions(),
	}
}

func (m *MockSDKAbi) SetAbis(abis []coreTypes.Abi) {
	m.SetItems(abis)
}

func (m *MockSDKAbi) SetFunctions(functions []coreTypes.Function) {
	m.functions = functions
}

type MessageCapture = testingPkg.MessageCapture
type CapturedMessage = testingPkg.CapturedMessage

func NewMessageCapture() *MessageCapture {
	return testingPkg.NewMessageCapture()
}
