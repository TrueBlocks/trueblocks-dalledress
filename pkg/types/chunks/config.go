package chunks

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// GetConfig returns the ViewConfig for the Chunks view
func (c *ChunksCollection) GetConfig() (*types.ViewConfig, error) {
	return &types.ViewConfig{
		ViewName: "chunks",
		Facets: map[string]types.FacetConfig{
			"stats": {
				Name:         "Stats",
				Store:        "chunks",
				IsForm:       false,
				Columns:      getStatsColumns(),
				DetailPanels: getStatsDetailPanels(),
				Actions:      []string{},
			},
			"index": {
				Name:         "Index",
				Store:        "chunks",
				IsForm:       false,
				Columns:      getIndexColumns(),
				DetailPanels: getIndexDetailPanels(),
				Actions:      []string{},
			},
			"blooms": {
				Name:         "Blooms",
				Store:        "chunks",
				IsForm:       false,
				Columns:      getBloomsColumns(),
				DetailPanels: getBloomsDetailPanels(),
				Actions:      []string{},
			},
			"manifest": {
				Name:         "Manifest",
				Store:        "chunks",
				IsForm:       true, // MANIFEST is a form view
				Columns:      nil,
				DetailPanels: getManifestDetailPanels(),
				Actions:      []string{},
			},
		},
		Actions: make(map[string]types.ActionConfig),
	}, nil
}

func getStatsColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{Key: "range", Header: "Range", Accessor: "range", Width: 150, Sortable: true, Filterable: true, Formatter: "blkrange"},
		{Key: "nAddrs", Header: "Addrs", Accessor: "nAddrs", Width: 120, Sortable: true, Filterable: true, Formatter: "number"},
		{Key: "nApps", Header: "Apps", Accessor: "nApps", Width: 100, Sortable: true, Filterable: true, Formatter: "number"},
		{Key: "nBlocks", Header: "Blocks", Accessor: "nBlocks", Width: 120, Sortable: true, Filterable: true, Formatter: "number"},
		{Key: "nBlooms", Header: "Blooms", Accessor: "nBlooms", Width: 120, Sortable: true, Filterable: true, Formatter: "number"},
		{Key: "recWid", Header: "Rec Wid", Accessor: "recWid", Width: 120, Sortable: true, Filterable: true, Formatter: "number"},
		{Key: "bloomSz", Header: "Bloom Sz", Accessor: "bloomSz", Width: 120, Sortable: true, Filterable: true, Formatter: "number"},
		{Key: "chunkSz", Header: "Chunk Sz", Accessor: "chunkSz", Width: 120, Sortable: true, Filterable: true, Formatter: "number"},
		{Key: "addrsPerBlock", Header: "Addrs Per Block", Accessor: "addrsPerBlock", Width: 100, Sortable: true, Filterable: true, Formatter: "float64"},
		{Key: "appsPerBlock", Header: "Apps Per Block", Accessor: "appsPerBlock", Width: 100, Sortable: true, Filterable: true, Formatter: "float64"},
		{Key: "appsPerAddr", Header: "Apps Per Addr", Accessor: "appsPerAddr", Width: 100, Sortable: true, Filterable: true, Formatter: "float64"},
		{Key: "ratio", Header: "Ratio", Accessor: "ratio", Width: 100, Sortable: true, Filterable: true, Formatter: "float64"},
	}
}

func getIndexColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{Key: "range", Header: "Range", Accessor: "range", Width: 150, Sortable: true, Filterable: true, Formatter: "blkrange"},
		{Key: "magic", Header: "Magic", Accessor: "magic", Width: 150, Sortable: true, Filterable: true, Formatter: "text"},
		{Key: "hash", Header: "Hash", Accessor: "hash", Width: 150, Sortable: true, Filterable: true, Formatter: "hash"},
		{Key: "nAddresses", Header: "Addresses", Accessor: "nAddresses", Width: 150, Sortable: true, Filterable: true, Formatter: "number"},
		{Key: "nAppearances", Header: "Appearances", Accessor: "nAppearances", Width: 150, Sortable: true, Filterable: true, Formatter: "number"},
		{Key: "size", Header: "Size", Accessor: "size", Width: 150, Sortable: true, Filterable: true, Formatter: "number"},
	}
}

func getBloomsColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{Key: "range", Header: "Range", Accessor: "range", Width: 150, Sortable: true, Filterable: true, Formatter: "blkrange"},
		{Key: "magic", Header: "Magic", Accessor: "magic", Width: 150, Sortable: true, Filterable: true, Formatter: "text"},
		{Key: "hash", Header: "Hash", Accessor: "hash", Width: 150, Sortable: true, Filterable: true, Formatter: "hash"},
		{Key: "nBlooms", Header: "Blooms", Accessor: "nBlooms", Width: 150, Sortable: true, Filterable: true, Formatter: "number"},
		{Key: "nInserted", Header: "Inserted", Accessor: "nInserted", Width: 150, Sortable: true, Filterable: true, Formatter: "number"},
		{Key: "size", Header: "Size", Accessor: "size", Width: 150, Sortable: true, Filterable: true, Formatter: "number"},
		{Key: "byteWidth", Header: "Byte Width", Accessor: "byteWidth", Width: 150, Sortable: true, Filterable: true, Formatter: "number"},
	}
}

func getStatsDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title: "Range",
			Fields: []types.DetailFieldConfig{
				{Key: "range", Label: "Range"},
			},
		},
		{
			Title: "Counts",
			Fields: []types.DetailFieldConfig{
				{Key: "nAddrs", Label: "Addresses"},
				{Key: "nApps", Label: "Apps"},
				{Key: "nBlocks", Label: "Blocks"},
				{Key: "nBlooms", Label: "Blooms"},
			},
		},
		{
			Title: "Sizes",
			Fields: []types.DetailFieldConfig{
				{Key: "recWid", Label: "Record Width"},
				{Key: "bloomSz", Label: "Bloom Size"},
				{Key: "chunkSz", Label: "Chunk Size"},
			},
		},
		{
			Title: "Efficiency",
			Fields: []types.DetailFieldConfig{
				{Key: "addrsPerBlock", Label: "Addrs/Block"},
				{Key: "appsPerBlock", Label: "Apps/Block"},
				{Key: "appsPerAddr", Label: "Apps/Addr"},
				{Key: "ratio", Label: "Ratio"},
			},
		},
	}
}

func getIndexDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title: "Range",
			Fields: []types.DetailFieldConfig{
				{Key: "range", Label: "Range"},
			},
		},
		{
			Title: "Identity",
			Fields: []types.DetailFieldConfig{
				{Key: "magic", Label: "Magic"},
				{Key: "hash", Label: "Hash", Formatter: "hash"},
			},
		},
		{
			Title: "Counts",
			Fields: []types.DetailFieldConfig{
				{Key: "nAddresses", Label: "Addresses"},
				{Key: "nAppearances", Label: "Appearances"},
			},
		},
		{
			Title: "Sizes",
			Fields: []types.DetailFieldConfig{
				{Key: "size", Label: "Size"},
			},
		},
	}
}

func getBloomsDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title: "Range",
			Fields: []types.DetailFieldConfig{
				{Key: "range", Label: "Range"},
			},
		},
		{
			Title: "Identity",
			Fields: []types.DetailFieldConfig{
				{Key: "magic", Label: "Magic"},
				{Key: "hash", Label: "Hash", Formatter: "hash"},
			},
		},
		{
			Title: "Counts",
			Fields: []types.DetailFieldConfig{
				{Key: "nBlooms", Label: "Blooms"},
				{Key: "nInserted", Label: "Inserted"},
			},
		},
		{
			Title: "Sizes",
			Fields: []types.DetailFieldConfig{
				{Key: "size", Label: "Size"},
				{Key: "byteWidth", Label: "Byte Width"},
			},
		},
	}
}

func getManifestDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title: "Manifest",
			Fields: []types.DetailFieldConfig{
				{Key: "version", Label: "Version"},
				{Key: "chain", Label: "Chain"},
				{Key: "specification", Label: "Specification", Formatter: "hash"},
			},
		},
	}
}
