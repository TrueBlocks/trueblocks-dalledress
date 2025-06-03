// ADD_ROUTE
package abis

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/repository"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

func (ac *AbisCollection) LoadData(listKind types.ListKind) {
	if !ac.NeedsUpdate(listKind) {
		return
	}

	switch listKind {
	case AbisDownloaded:
		go ac.loadDownloadedAbis()
	case AbisKnown:
		go ac.loadKnownAbis()
	case AbisFunctions:
		go ac.loadFunctions()
	case AbisEvents:
		go ac.loadEvents()
	default:
		logger.Error(fmt.Sprintf("AbisCollection.LoadData: unexpected list kind: %v", listKind))
	}
}

func (ac *AbisCollection) loadDownloadedAbis() {
	result, err := ac.downloadedRepo.Load(repository.LoadOptions{})
	if err != nil {
		logging.LogError("AbisCollection.loadDownloadedAbis: %v", err, repository.ErrorAlreadyLoading)
		return
	}
	msgs.EmitStatus(result.Status)
	msgs.EmitPayload(msgs.EventDataLoaded, "", result.Payload)
}

func (ac *AbisCollection) loadKnownAbis() {
	result, err := ac.knownRepo.Load(repository.LoadOptions{})
	if err != nil {
		logging.LogError("AbisCollection.loadKnownAbis: %v", err, repository.ErrorAlreadyLoading)
		return
	}
	msgs.EmitStatus(result.Status)
	msgs.EmitPayload(msgs.EventDataLoaded, result.Payload)
}

func (ac *AbisCollection) loadFunctions() {
	result, err := ac.functionsRepo.Load(repository.LoadOptions{})
	if err != nil {
		logging.LogError("AbisCollection.loadFunctions: %v", err, repository.ErrorAlreadyLoading)
		return
	}
	msgs.EmitStatus(result.Status)
	msgs.EmitPayload(msgs.EventDataLoaded, result.Payload)
}

func (ac *AbisCollection) loadEvents() {
	result, err := ac.eventsRepo.Load(repository.LoadOptions{})
	if err != nil {
		logging.LogError("AbisCollection.loadEvents: %v", err, repository.ErrorAlreadyLoading)
		return
	}
	msgs.EmitStatus(result.Status)
	msgs.EmitPayload(msgs.EventDataLoaded, result.Payload)
}

// ADD_ROUTE
