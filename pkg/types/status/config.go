package status

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// GetConfig returns the ViewConfig for the Status view
func (c *StatusCollection) GetConfig() (*types.ViewConfig, error) {
	return &types.ViewConfig{
		ViewName: "status",
		Facets: map[string]types.FacetConfig{
			"status": {
				Name:          "Status",
				Store:         "status",
				IsForm:        true,
				Columns:       getStatusColumns(),
				DetailPanels:  getStatusDetailPanels(),
				Actions:       []string{}, // Status view has no actions
				HeaderActions: []string{},
			},
			"caches": {
				Name:          "Caches",
				Store:         "caches",
				Columns:       getCachesColumns(),
				DetailPanels:  getCachesDetailPanels(),
				Actions:       []string{}, // Caches view has no actions
				HeaderActions: []string{"export"},
			},
			"chains": {
				Name:          "Chains",
				Store:         "chains",
				Columns:       getChainsColumns(),
				DetailPanels:  getChainsDetailPanels(),
				Actions:       []string{}, // Chains view has no actions
				HeaderActions: []string{"export"},
			},
		},
		Actions: map[string]types.ActionConfig{
			"export": {Name: "export", Label: "Export Data", Icon: "Export"},
		},
		FacetOrder: []string{"status", "caches", "chains"},
	}, nil
}

func getStatusColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{Key: "cachePath", Header: "Cache Path", Accessor: "cachePath", Width: 200},
		{Key: "indexPath", Header: "Index Path", Accessor: "indexPath", Width: 200},
		{Key: "chain", Header: "Chain", Accessor: "chain", Width: 100},
		{Key: "chainId", Header: "Chain ID", Accessor: "chainId", Width: 100},
		{Key: "networkId", Header: "Network ID", Accessor: "networkId", Width: 100},
		{Key: "chainConfig", Header: "Chain Config", Accessor: "chainConfig", Width: 150},
		{Key: "rootConfig", Header: "Root Config", Accessor: "rootConfig", Width: 150},
		{Key: "clientVersion", Header: "Client Version", Accessor: "clientVersion", Width: 150},
		{Key: "version", Header: "Version", Accessor: "version", Width: 100},
		{Key: "progress", Header: "Progress", Accessor: "progress", Width: 100},
		{Key: "rpcProvider", Header: "RPC Provider", Accessor: "rpcProvider", Width: 200},
		{Key: "hasEsKey", Header: "Has ES Key", Accessor: "hasEsKey", Width: 100},
		{Key: "hasPinKey", Header: "Has Pin Key", Accessor: "hasPinKey", Width: 100},
		{Key: "isApi", Header: "Is API", Accessor: "isApi", Width: 100},
		{Key: "isArchive", Header: "Is Archive", Accessor: "isArchive", Width: 100},
		{Key: "isScraping", Header: "Is Scraping", Accessor: "isScraping", Width: 100},
		{Key: "isTesting", Header: "Is Testing", Accessor: "isTesting", Width: 100},
		{Key: "isTracing", Header: "Is Tracing", Accessor: "isTracing", Width: 100},
	}
}

func getStatusDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title: "Paths",
			Fields: []types.DetailFieldConfig{
				{Key: "cachePath", Label: "Cache Path"},
				{Key: "indexPath", Label: "Index Path"},
			},
		},
		{
			Title: "Chain",
			Fields: []types.DetailFieldConfig{
				{Key: "chain", Label: "Chain"},
				{Key: "chainId", Label: "Chain ID"},
				{Key: "networkId", Label: "Network ID"},
				{Key: "chainConfig", Label: "Chain Config"},
			},
		},
		{
			Title: "Config",
			Fields: []types.DetailFieldConfig{
				{Key: "rootConfig", Label: "Root Config"},
				{Key: "clientVersion", Label: "Client Version"},
				{Key: "version", Label: "Version"},
			},
		},
		{
			Title: "Progress",
			Fields: []types.DetailFieldConfig{
				{Key: "progress", Label: "Progress"},
			},
		},
		{
			Title: "Providers",
			Fields: []types.DetailFieldConfig{
				{Key: "rpcProvider", Label: "RPC Provider"},
			},
		},
		{
			Title: "Flags",
			Fields: []types.DetailFieldConfig{
				{Key: "hasEsKey", Label: "Has ES Key"},
				{Key: "hasPinKey", Label: "Has Pin Key"},
				{Key: "isApi", Label: "Is API"},
				{Key: "isArchive", Label: "Is Archive"},
				{Key: "isScraping", Label: "Is Scraping"},
				{Key: "isTesting", Label: "Is Testing"},
				{Key: "isTracing", Label: "Is Tracing"},
			},
		},
	}
}

func getCachesColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{Key: "type", Header: "Type", Accessor: "type", Width: 100},
		{Key: "path", Header: "Path", Accessor: "path", Width: 300},
		{Key: "nFiles", Header: "Files", Accessor: "nFiles", Width: 100},
		{Key: "nFolders", Header: "Folders", Accessor: "nFolders", Width: 100},
		{Key: "sizeInBytes", Header: "Size (Bytes)", Accessor: "sizeInBytes", Width: 150},
		{Key: "lastCached", Header: "Last Cached", Accessor: "lastCached", Width: 150},
	}
}

func getCachesDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title: "General",
			Fields: []types.DetailFieldConfig{
				{Key: "type", Label: "Type"},
				{Key: "path", Label: "Path"},
			},
		},
		{
			Title: "Statistics",
			Fields: []types.DetailFieldConfig{
				{Key: "nFiles", Label: "Files"},
				{Key: "nFolders", Label: "Folders"},
				{Key: "sizeInBytes", Label: "Size (Bytes)"},
			},
		},
		{
			Title: "Timestamps",
			Fields: []types.DetailFieldConfig{
				{Key: "lastCached", Label: "Last Cached"},
			},
		},
	}
}

func getChainsColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{Key: "chain", Header: "Chain", Accessor: "chain", Width: 100},
		{Key: "chainId", Header: "Chain ID", Accessor: "chainId", Width: 100},
		{Key: "symbol", Header: "Symbol", Accessor: "symbol", Width: 100},
		{Key: "rpcProvider", Header: "RPC Provider", Accessor: "rpcProvider", Width: 200},
		{Key: "ipfsGateway", Header: "IPFS Gateway", Accessor: "ipfsGateway", Width: 200},
		{Key: "localExplorer", Header: "Local Explorer", Accessor: "localExplorer", Width: 200},
		{Key: "remoteExplorer", Header: "Remote Explorer", Accessor: "remoteExplorer", Width: 200},
	}
}

func getChainsDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title: "General",
			Fields: []types.DetailFieldConfig{
				{Key: "chain", Label: "Chain"},
				{Key: "chainId", Label: "Chain ID"},
				{Key: "symbol", Label: "Symbol"},
			},
		},
		{
			Title: "Providers",
			Fields: []types.DetailFieldConfig{
				{Key: "rpcProvider", Label: "RPC Provider"},
				{Key: "ipfsGateway", Label: "IPFS Gateway"},
			},
		},
		{
			Title: "Explorers",
			Fields: []types.DetailFieldConfig{
				{Key: "localExplorer", Label: "Local Explorer"},
				{Key: "remoteExplorer", Label: "Remote Explorer"},
			},
		},
	}
}
