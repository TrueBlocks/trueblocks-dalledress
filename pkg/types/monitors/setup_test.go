package monitors

import (
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	testingPkg "github.com/TrueBlocks/trueblocks-dalledress/pkg/testing"
)

type MockSDKMonitor struct {
	*testingPkg.MockSDKBase[coreTypes.Monitor]
}

func NewMockSDKMonitor() *MockSDKMonitor {
	return &MockSDKMonitor{
		MockSDKBase: testingPkg.NewMockSDKBase(testingPkg.CreateTestMonitors()),
	}
}

func (m *MockSDKMonitor) SetMonitors(monitors []coreTypes.Monitor) {
	m.SetItems(monitors)
}
