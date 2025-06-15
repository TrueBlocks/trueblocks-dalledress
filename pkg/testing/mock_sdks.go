package testing

import (
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

type MockSDKAbi struct {
	*MockSDKBase[coreTypes.Abi]
	functions []coreTypes.Function
}

func NewMockSDKAbi() *MockSDKAbi {
	return &MockSDKAbi{
		MockSDKBase: NewMockSDKBase(CreateTestAbis()),
		functions:   CreateTestFunctions(),
	}
}

func (m *MockSDKAbi) SetFunctions(functions []coreTypes.Function) {
	m.functions = functions
}

func (m *MockSDKAbi) GetFunctions() []coreTypes.Function {
	if m.GetError() != nil {
		return nil
	}
	return m.functions
}

type MockSDKMonitor struct {
	*MockSDKBase[coreTypes.Monitor]
}

func NewMockSDKMonitor() *MockSDKMonitor {
	return &MockSDKMonitor{
		MockSDKBase: NewMockSDKBase(CreateTestMonitors()),
	}
}

type MockSDKName struct {
	*MockSDKBase[coreTypes.Name]
}

func NewMockSDKName() *MockSDKName {
	return &MockSDKName{
		MockSDKBase: NewMockSDKBase(CreateTestNames()),
	}
}
