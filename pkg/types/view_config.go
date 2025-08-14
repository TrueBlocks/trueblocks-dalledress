package types

// ViewConfig represents the complete configuration for a view
type ViewConfig struct {
	ViewName   string                  `json:"viewName"`
	Facets     map[string]FacetConfig  `json:"facets"`
	Actions    map[string]ActionConfig `json:"actions"`
	FacetOrder []string                `json:"facetOrder"`
}

// FacetConfig represents configuration for a single facet within a view
type FacetConfig struct {
	Name          string              `json:"name"`
	Store         string              `json:"store"`
	IsForm        bool                `json:"isForm"`
	DividerBefore bool                `json:"dividerBefore"`
	Fields        []FieldConfig       `json:"fields"`
	Columns       []ColumnConfig      `json:"columns"`
	DetailPanels  []DetailPanelConfig `json:"detailPanels"`
	Actions       []string            `json:"actions"`
	HeaderActions []string            `json:"headerActions"`
}

// FieldConfig is the single source-of-truth for facet fields
// Header is derived from Label; do not duplicate.
type FieldConfig struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Formatter   string `json:"formatter"`
	Section     string `json:"section"`
	Width       int    `json:"width"`
	Sortable    bool   `json:"sortable"`
	Filterable  bool   `json:"filterable"`
	Order       int    `json:"order"`
	DetailOrder int    `json:"detailOrder"`
	NoTable     bool   `json:"-"`
	NoDetail    bool   `json:"-"`
	ColumnLabel string `json:"-"`
	DetailLabel string `json:"-"`
}

// ColumnConfig represents a table column configuration
type ColumnConfig struct {
	Key        string `json:"key"`
	Header     string `json:"header"`
	Width      int    `json:"width"`
	Sortable   bool   `json:"sortable"`
	Filterable bool   `json:"filterable"`
	Formatter  string `json:"formatter"`
	Order      int    `json:"order"`
}

// DetailPanelConfig represents a detail panel section
type DetailPanelConfig struct {
	Title     string              `json:"title"`
	Collapsed bool                `json:"collapsed"`
	Fields    []DetailFieldConfig `json:"fields"`
}

// DetailFieldConfig represents a field in a detail panel
type DetailFieldConfig struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Formatter   string `json:"formatter"`
	DetailOrder int    `json:"detailOrder"`
}

// ActionConfig represents an available action (delete, remove, etc.)
type ActionConfig struct {
	Name         string `json:"name"`
	Label        string `json:"label"`
	Icon         string `json:"icon"`
	Confirmation bool   `json:"confirmation"`
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
