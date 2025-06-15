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
	MinTickTime      = 75 * time.Millisecond
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
	pr.currentIncrement = initialIncrement
	pr.nextThreshold = firstPageCount + pr.currentIncrement
	pr.firstDataSent = false
	return pr
}

func (pr *Progress) Tick(currentTotalCount, expectedTotal int) types.DataLoadedPayload {
	pr.nItemsSinceUpdate++
	shouldUpdate := false
	isFirstPageEvent := false

	now := time.Now()

	if !pr.firstDataSent && currentTotalCount >= firstPageCount {
		shouldUpdate = true
		isFirstPageEvent = true
	} else if pr.firstDataSent && currentTotalCount >= pr.nextThreshold {
		shouldUpdate = true
	}

	sendUpdateThisTick := false
	if shouldUpdate {
		if isFirstPageEvent {
			sendUpdateThisTick = true
			if pr.onFirstDataFunc != nil {
				go pr.onFirstDataFunc()
			}
			pr.firstDataSent = true
		} else {
			if now.Sub(pr.lastUpdate) >= MinTickTime {
				sendUpdateThisTick = true
			}
		}
	}

	if pr.firstDataSent && !isFirstPageEvent && currentTotalCount >= pr.nextThreshold {
		pr.currentIncrement += incrementGrowth
		pr.nextThreshold = currentTotalCount + pr.currentIncrement
	}

	payload := types.DataLoadedPayload{
		CurrentCount:  currentTotalCount,
		ExpectedTotal: expectedTotal,
		ListKind:      pr.listKind,
	}

	if sendUpdateThisTick {
		msgs.EmitLoaded("streaming", payload)
		msgs.EmitStatus(progress(currentTotalCount, pr.listKind, false))
		pr.nItemsSinceUpdate = 0
		pr.lastUpdate = now
		if isFirstPageEvent && currentTotalCount >= pr.nextThreshold {
			pr.currentIncrement += incrementGrowth
			pr.nextThreshold = currentTotalCount + pr.currentIncrement
		}
	}

	return payload
}

func (pr *Progress) Heartbeat(currentTotalCount, expectedTotal int) types.DataLoadedPayload {
	payload := types.DataLoadedPayload{
		CurrentCount:  currentTotalCount,
		ExpectedTotal: expectedTotal,
		ListKind:      pr.listKind,
	}

	now := time.Now()
	if now.Sub(pr.lastUpdate) >= MaxWaitTime && pr.nItemsSinceUpdate > 0 {
		msgs.EmitLoaded("partial", payload)
		msgs.EmitStatus(progress(currentTotalCount, pr.listKind, true))

		pr.nItemsSinceUpdate = 0
		pr.lastUpdate = now
	}

	return payload
}

func progress(cnt int, kind types.ListKind, heartbeat bool) string {
	k := strings.Trim(strings.ToLower(string(kind)), " ")
	if heartbeat {
		return fmt.Sprintf("Loaded %d %s...", cnt, k)
	}
	return fmt.Sprintf("Loaded %d %s.", cnt, k)
}
