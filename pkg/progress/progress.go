package progress

import (
	"fmt"
	"strings"
	"time"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

const (
	firstPageCount   = 7
	initialIncrement = 10
	incrementGrowth  = 10
	MaxWaitTime      = 125 * time.Millisecond
)

// Progress manages the logic for sending progress updates.
type Progress struct {
	lastUpdate        time.Time
	nItemsSinceUpdate int
	nextThreshold     int
	currentIncrement  int
	listKind          types.ListKind
	onFirstDataFunc   func()
	firstDataSent     bool
}

// NewProgress creates and initializes a Progress.
func NewProgress(
	listKindCfg types.ListKind,
	onFirstData func(), // Can be nil
) *Progress {
	pr := &Progress{
		listKind:        listKindCfg,
		onFirstDataFunc: onFirstData,
	}
	// Initialize internal state
	pr.lastUpdate = time.Now()
	pr.nItemsSinceUpdate = 0
	pr.nextThreshold = firstPageCount + initialIncrement
	pr.currentIncrement = initialIncrement
	pr.firstDataSent = false
	return pr
}

func (pr *Progress) Tick(currentTotalCount, expectedTotal int) types.DataLoadedPayload {
	pr.nItemsSinceUpdate++
	shouldUpdate := false

	if !pr.firstDataSent && currentTotalCount == firstPageCount {
		shouldUpdate = true
		if pr.onFirstDataFunc != nil {
			go pr.onFirstDataFunc()
		}
		pr.firstDataSent = true
	} else if currentTotalCount >= pr.nextThreshold && currentTotalCount > firstPageCount {
		shouldUpdate = true
		pr.currentIncrement += incrementGrowth
		pr.nextThreshold = currentTotalCount + pr.currentIncrement
	}

	payload := types.DataLoadedPayload{
		CurrentCount:  currentTotalCount,
		ExpectedTotal: expectedTotal,
		ListKind:      pr.listKind,
	}
	if shouldUpdate {
		msgs.EmitLoaded("streaming", payload)
		msgs.EmitStatus(progress(currentTotalCount, pr.listKind, false))
		pr.nItemsSinceUpdate = 0
		pr.lastUpdate = time.Now()
	}

	return payload
}

func (pr *Progress) Heartbeat(currentTotalCount, expectedTotal int) types.DataLoadedPayload {
	payload := types.DataLoadedPayload{
		CurrentCount:  currentTotalCount,
		ExpectedTotal: expectedTotal,
		ListKind:      pr.listKind,
	}

	if time.Since(pr.lastUpdate) >= MaxWaitTime && pr.nItemsSinceUpdate > 0 {
		msgs.EmitLoaded("partial", payload)
		msgs.EmitStatus(progress(currentTotalCount, pr.listKind, true))

		pr.nItemsSinceUpdate = 0
		pr.lastUpdate = time.Now()
	}

	return payload
}

func progress(cnt int, kind types.ListKind, heartbeat bool) string {
	k := strings.Trim(strings.ToLower(string(kind)), " ")
	if heartbeat {
		return fmt.Sprintf("Loaded %d %s...", cnt, k)
	}
	return fmt.Sprintf("Loaded %d %s", cnt, k)
}
