package app

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/contracts"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/exports"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/names"
)

// ExportData handles export requests with full context logging and CSV generation
func (a *App) ExportData(payload *types.Payload) error {
	activeProject := a.Projects.GetActiveProject()
	if activeProject == nil {
		err := fmt.Errorf("no active project")
		msgs.EmitError("export failed: no active project", err)
		return err
	}
	payload.ProjectPath = activeProject.Path

	collection := getCollection(payload)
	if collection == nil {
		err := fmt.Errorf("unsupported collection type: %s", payload.Collection)
		msgs.EmitError("unsupported collection type", err)
		return err
	}

	exportFilename, err := collection.ExportData(payload)
	if err != nil {
		msgs.EmitError("failed to export data", err)
		return fmt.Errorf("failed to export data: %w", err)
	}

	cmd := "open \"" + exportFilename + "\""
	exitCode := utils.System(cmd)
	if exitCode != 0 {
		logging.LogBackend(fmt.Sprintf("Failed to open export file, exit code: %d", exitCode))
	}

	statusMsg := fmt.Sprintf("Export completed: %s %s data", payload.Collection, payload.DataFacet)
	if payload.Address != "" && payload.Address != "0x0" {
		statusMsg += fmt.Sprintf(" for %s", payload.Address[:10]+"...")
	}
	if payload.Chain != "" {
		statusMsg += fmt.Sprintf(" on %s", payload.Chain)
	}
	msgs.EmitStatus(statusMsg)

	return nil
}

func getCollection(payload *types.Payload) types.Collection {
	// TODO: BOGUS - THIS NEEDS TO HAVE ALL TYPES
	switch payload.Collection {
	case "exports":
		return exports.GetExportsCollection(payload)
	case "names":
		return names.GetNamesCollection(payload)
	case "contracts":
		return contracts.GetContractsCollection(payload)
	default:
		logging.LogBackend(fmt.Sprintf("Warning: Unknown collection type: %s", payload.Collection))
		return nil
	}
}
