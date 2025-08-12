package types

// ViewConfig represents the complete configuration for a view
type ViewConfig struct {
	ViewName   string                  `json:"viewName" wails:"viewName"`
	Facets     map[string]FacetConfig  `json:"facets" wails:"facets"`
	Actions    map[string]ActionConfig `json:"actions" wails:"actions"`
	FacetOrder []string                `json:"facetOrder" wails:"facetOrder"`
}

// FacetConfig represents configuration for a single facet within a view
type FacetConfig struct {
	Name          string              `json:"name" wails:"name"`
	Store         string              `json:"store" wails:"store"`
	IsForm        bool                `json:"isForm" wails:"isForm"`
	DividerBefore bool                `json:"dividerBefore" wails:"dividerBefore"`
	Columns       []ColumnConfig      `json:"columns" wails:"columns"`
	DetailPanels  []DetailPanelConfig `json:"detailPanels" wails:"detailPanels"`
	Actions       []string            `json:"actions" wails:"actions"`
	HeaderActions []string            `json:"headerActions" wails:"headerActions"`
}

// ColumnConfig represents a table column configuration
type ColumnConfig struct {
	Key        string `json:"key" wails:"key"`
	Header     string `json:"header" wails:"header"`
	Accessor   string `json:"accessor" wails:"accessor"`
	Width      int    `json:"width" wails:"width"`
	Sortable   bool   `json:"sortable" wails:"sortable"`
	Filterable bool   `json:"filterable" wails:"filterable"`
	Formatter  string `json:"formatter" wails:"formatter"`
}

// DetailPanelConfig represents a detail panel section
type DetailPanelConfig struct {
	Title     string              `json:"title" wails:"title"`
	Collapsed bool                `json:"collapsed" wails:"collapsed"`
	Fields    []DetailFieldConfig `json:"fields" wails:"fields"`
}

// DetailFieldConfig represents a field in a detail panel
type DetailFieldConfig struct {
	Key       string `json:"key" wails:"key"`
	Label     string `json:"label" wails:"label"`
	Formatter string `json:"formatter" wails:"formatter"`
}

// ActionConfig represents an available action (delete, remove, etc.)
type ActionConfig struct {
	Name         string   `json:"name" wails:"name"`
	Label        string   `json:"label" wails:"label"`
	Icon         string   `json:"icon" wails:"icon"`
	Confirmation bool     `json:"confirmation" wails:"confirmation"`
	Facets       []string `json:"facets" wails:"facets"` // Which facets allow this action
}

/*
FUTURE FEATURES (commented for reference):

// Conditional Detail Panels
type PanelSelectorConfig struct {
	Condition ConditionConfig     `json:"condition"`
	Panels    []DetailPanelConfig `json:"panels"`
}

type ConditionConfig struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"` // "equals", "contains", "exists"
	Value    interface{} `json:"value"`
}

// User Customization
type UserViewConfig struct {
	BaseConfig ViewConfig          `json:"baseConfig"`
	Overrides  ViewConfigOverrides `json:"overrides"`
}

type ViewConfigOverrides struct {
	HiddenColumns []string                   `json:"hiddenColumns"`
	ColumnOrder   []string                   `json:"columnOrder"`
	CustomPanels  []DetailPanelConfig        `json:"customPanels"`
	Facets        map[string]FacetOverrides  `json:"facets"`
}

// Computed Fields
type ComputedFieldConfig struct {
	Key        string `json:"key"`
	Expression string `json:"expression"` // JS expression or function name
	DependsOn  []string `json:"dependsOn"` // Field dependencies
}
*/
