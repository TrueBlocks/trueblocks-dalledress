package testing

import (
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

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
