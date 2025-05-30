// ADD_ROUTE
package types

import (
	"sync/atomic"
)

var refreshRate = 31

// LoadData kicks off a go routine that streams the requested data making it available
// to GetPage as soon as it becomes available (even partially).
func (ac *AbisCollection) LoadData(listKind ListKind) {
	if !atomic.CompareAndSwapInt32(&ac.isLoading, 0, 1) {
		return
	}

	needsUpdate := ac.NeedsUpdate(listKind)
	if needsUpdate {
		go ac.loadInternal(listKind)
	} else {
		atomic.StoreInt32(&ac.isLoading, 0)
	}
}

// NeedsUpdate checks if the AbisCollection needs to be updated.
func (ac *AbisCollection) NeedsUpdate(listKind ListKind) bool {
	ac.mutex.RLock()
	defer ac.mutex.RUnlock()

	switch listKind {
	case AbisDownloaded:
		return !ac.isDownloadedLoaded
	case AbisKnown:
		return !ac.isKnownLoaded
	case AbisFunctions:
		return !ac.isFuncsLoaded
	case AbisEvents:
		return !ac.isEventsLoaded
	default:
		return true
	}
}

// ADD_ROUTE
