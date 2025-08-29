package dalledress

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	dalle "github.com/TrueBlocks/trueblocks-dalle/v2"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// SeriesCrud handles create, update, remove for series facet
func (c *DalleDressCollection) SeriesCrud(
	payload *types.Payload,
	op crud.Operation,
	item *Series,
) error {
	if payload.DataFacet != DalleDressSeries {
		return fmt.Errorf("SeriesCrud invalid facet: %s", payload.DataFacet)
	}
	if item == nil {
		return fmt.Errorf("SeriesCrud missing item")
	}
	seriesDir := filepath.Join(dalle.DataDir(), "series")
	// Support a pseudo duplicate operation encoded by passing an item whose Suffix ends with "-copy" pattern AND op == crud.Create with a source indicated in item.Last (temporary convention) is overkill.
	// Instead, frontend will call Create with full cloned object; so duplicate maps to plain Create here.
	switch op {
	case crud.Create, crud.Update:
		if op == crud.Create {
			// uniqueness check
			for _, existing := range c.seriesFacet.GetStore().GetItems() {
				if existing != nil && existing.Suffix == item.Suffix {
					return fmt.Errorf("series suffix already exists: %s", item.Suffix)
				}
			}
		}
		item.SaveSeries(item.Suffix, item.Last)
		// update ModifiedAt from file system
		if fi, err := os.Stat(seriesDir + "/" + item.Suffix + ".json"); err == nil {
			item.ModifiedAt = fi.ModTime().UTC().Format(time.RFC3339)
		}
	case crud.Remove:
		_ = dalle.DeleteSeries(seriesDir, item.Suffix)
	case crud.Autoname:
		// not applicable
	case crud.Delete, crud.Undelete:
		// unsupported logical deletion
	default:
		return fmt.Errorf("unsupported op %v", op)
	}
	store := c.seriesFacet.GetStore()
	store.UpdateData(func(data []*Series) []*Series {
		switch op {
		case crud.Remove:
			out := make([]*Series, 0, len(data))
			for _, s := range data {
				if s.Suffix != item.Suffix {
					out = append(out, s)
				}
			}
			return out
		case crud.Create:
			for _, s := range data {
				if s.Suffix == item.Suffix {
					*s = *item
					return data
				}
			}
			return append(data, item)
		case crud.Update:
			for _, s := range data {
				if s.Suffix == item.Suffix {
					*s = *item
					break
				}
			}
			return data
		}
		return data
	})
	c.seriesFacet.SyncWithStore()
	// emit enriched event so frontend refreshes with context
	currentItems := c.seriesFacet.GetStore().GetItems()
	currentCount := len(currentItems)
	payloadSummary := types.Summary{TotalCount: currentCount, FacetCounts: map[types.DataFacet]int{DalleDressSeries: currentCount}, LastUpdated: time.Now().Unix()}
	operation := "update"
	switch op {
	case crud.Create:
		operation = "create"
	case crud.Update:
		operation = "update"
	case crud.Remove:
		operation = "remove"
	}
	msgs.EmitLoaded(types.DataLoadedPayload{
		Payload:       types.Payload{Collection: "dalledress", DataFacet: DalleDressSeries},
		CurrentCount:  currentCount,
		ExpectedTotal: currentCount,
		State:         types.StateLoaded,
		Summary:       payloadSummary,
		Timestamp:     time.Now().Unix(),
		EventPhase:    "complete",
		Operation:     operation,
	})
	return nil
}
