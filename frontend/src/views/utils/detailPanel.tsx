/* eslint-disable @typescript-eslint/no-explicit-any */
import { types } from '@models';

import { DetailTable } from '../../components/detail/DetailTable';

/**
 * Type definition for MetaOverlay configuration
 */
export interface MetaOverlay {
  section: string;
  detailLabel?: string;
  detailFormat?: 'address' | 'hash' | string;
  detailOrder?: number;
  detailOnly?: boolean;
  formatters?: Record<string, (value: any) => string>;
  extras?: Record<string, any>;
}

/**
 * Creates a detail panel from ViewConfig or falls back to default
 * This function encapsulates all ViewConfig integration logic and returns
 * a simple detail panel function for use in BaseTab.
 */
export const createDetailPanelFromViewConfig = <
  T extends Record<string, unknown>,
>(
  viewConfig: types.ViewConfig | null | undefined,
  getCurrentDataFacet: () => string,
  fallbackName = 'Details',
) => {
  // Get the current facet configuration
  const currentFacetConfig = viewConfig?.facets?.[getCurrentDataFacet()];

  // If we have ViewConfig detail panels, use them
  if (currentFacetConfig?.detailPanels?.length) {
    return buildDetailPanelFromConfigs(currentFacetConfig.detailPanels);
  }

  // Otherwise, return a fallback detail panel
  const FallbackDetailPanel = (rowData: T | null) => {
    if (!rowData) {
      return <div className="no-selection">Select a row to view details</div>;
    }
    return (
      <div>
        <h4>{fallbackName}</h4>
        <div className="detail-panel-default">
          {Object.entries(rowData).map(([key, value]) => (
            <div key={key} className="detail-row">
              <strong>{key}:</strong> {String(value)}
            </div>
          ))}
        </div>
      </div>
    );
  };
  FallbackDetailPanel.displayName = `${fallbackName}DetailPanel`;
  return FallbackDetailPanel;
};

/**
 * Creates a default detail panel function that renders key-value pairs
 * for any row data object. This is the fallback used when no custom
 * detail panel is provided.
 */
export const createDefaultDetailPanel = <
  T extends Record<string, unknown>,
>() => {
  const DefaultDetailPanel = (rowData: T | null) => {
    if (!rowData) {
      return <div className="no-selection">Select a row to view details</div>;
    }

    return (
      <div>
        <h4>
          {(rowData as any).name
            ? `${(rowData as any).name} Details`
            : 'Row Details'}
        </h4>
        <div className="detail-panel-default">
          {Object.entries(rowData).map(([key, value]) => (
            <div key={key} className="detail-row">
              <strong>{key}:</strong> {String(value)}
            </div>
          ))}
        </div>
      </div>
    );
  };

  DefaultDetailPanel.displayName = 'DefaultDetailPanel';
  return DefaultDetailPanel;
};

/**
 * Bridge function to build detail panels from column configurations.
 * This is a temporary bridge until all views migrate to backend ViewConfig.
 */
export const buildDetailPanelFromColumns = <T extends Record<string, unknown>>(
  columns: any[],
  meta: Record<string, MetaOverlay>,
  options?: {
    title?: string;
    collapsedSections?: string[];
    formatters?: Record<string, (value: any) => string>;
    extras?: Record<string, any>;
  },
) => {
  const DetailPanel = (rowData: T | null) => {
    if (!rowData) {
      return <div className="no-selection">Select a row to view details</div>;
    }

    // Group columns by section if meta is provided
    const sections: Record<
      string,
      Array<{
        key: string;
        column: any;
        meta: MetaOverlay | undefined;
      }>
    > = {};

    columns.forEach((column) => {
      const key = column.accessor || column.key;
      const columnMeta = meta[key];
      const sectionName = columnMeta?.section || 'General';

      if (!sections[sectionName]) {
        sections[sectionName] = [];
      }

      sections[sectionName].push({ key, column, meta: columnMeta });
    });

    return (
      <div>
        {options?.title && <h4>{options.title}</h4>}
        <div className="detail-panel-columns">
          {Object.entries(sections).map(([sectionName, sectionColumns]) => (
            <div key={sectionName} className="detail-section">
              <h5>{sectionName}</h5>
              {sectionColumns
                .sort(
                  (a, b) =>
                    (a.meta?.detailOrder || 99) - (b.meta?.detailOrder || 99),
                )
                .map(({ key, column, meta }) => {
                  const value = (rowData as any)[key];
                  if (value === undefined || value === null) return null;

                  return (
                    <div key={key} className="detail-row">
                      <strong>
                        {meta?.detailLabel ||
                          column.header ||
                          column.Header ||
                          key}
                        :
                      </strong>{' '}
                      {formatFieldValue(value, meta?.detailFormat)}
                    </div>
                  );
                })}
            </div>
          ))}
        </div>
      </div>
    );
  };

  DetailPanel.displayName = `DetailPanel_${options?.title || 'Columns'}`;
  return DetailPanel;
};

/**
 * Builds detail panels from backend ViewConfig DetailPanelConfig.
 * This is the new approach that uses backend-generated configuration.
 */
export const buildDetailPanelFromConfig = <T extends Record<string, unknown>>(
  panelConfig: types.DetailPanelConfig,
) => {
  const DetailPanel = (rowData: T | null) => {
    if (!rowData) {
      return <div className="no-selection">Select a row to view details</div>;
    }

    return (
      <div>
        <h4>{panelConfig.title}</h4>
        <div className="detail-panel-config">
          {panelConfig.fields.map((field) => {
            const value = (rowData as any)[field.key];
            if (value === undefined || value === null) return null;

            return (
              <div key={field.key} className="detail-row">
                <strong>{field.label}:</strong>{' '}
                {formatFieldValue(value, field.formatter)}
              </div>
            );
          })}
        </div>
      </div>
    );
  };

  DetailPanel.displayName = `DetailPanel_${panelConfig.title}`;
  return DetailPanel;
};

/**
 * Builds detail panels from multiple DetailPanelConfig objects.
 * This creates a combined panel with multiple sections using proper DetailTable styling.
 */
export const buildDetailPanelFromConfigs = <T extends Record<string, unknown>>(
  panelConfigs: types.DetailPanelConfig[],
) => {
  const DetailPanel = (rowData: T | null) => {
    if (!rowData) {
      return <div className="no-selection">Select a row to view details</div>;
    }

    // Convert DetailPanelConfig to DetailTable format
    const sections = panelConfigs.map((panelConfig) => ({
      name: panelConfig.title,
      rows: panelConfig.fields
        .map((field) => {
          const value = (rowData as any)[field.key];
          if (value === undefined || value === null) return null;

          return {
            label: field.label,
            value: formatFieldValue(value, field.formatter),
          };
        })
        .filter(Boolean) as { label: string; value: React.ReactNode }[],
    }));

    return (
      <DetailTable
        sections={sections}
        defaultCollapsedSections={['Statistics']} // Start with some sections collapsed
      />
    );
  };

  DetailPanel.displayName = `DetailPanel_Combined`;
  return DetailPanel;
};

/**
 * Format field values based on formatter type
 */
function formatFieldValue(value: any, formatter?: string): string {
  if (!formatter) return String(value);

  switch (formatter) {
    case 'address':
      return `${String(value).slice(0, 8)}...${String(value).slice(-6)}`;
    case 'hash':
      return `${String(value).slice(0, 10)}...${String(value).slice(-8)}`;
    case 'number':
      return Number(value).toLocaleString();
    case 'fileSize':
      // Simple file size formatting
      const bytes = Number(value);
      if (bytes === 0) return '0 B';
      const k = 1024;
      const sizes = ['B', 'KB', 'MB', 'GB'];
      const i = Math.floor(Math.log(bytes) / Math.log(k));
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    case 'boolean':
      return value ? 'Yes' : 'No';
    case 'timestamp':
      return new Date(Number(value) * 1000).toLocaleString();
    case 'computed':
      // For computed fields, we'd need additional logic here
      return String(value);
    default:
      return String(value);
  }
}

/**
 * Generic function to get the appropriate detail panel for a view and facet.
 * Looks up custom detail panels from the provided configuration, falls back
 * to the default renderer if no custom panel is found.
 *
 * @param viewName - The name of the view (e.g., 'monitors', 'exports')
 * @param facet - The current data facet
 * @param customPanels - Configuration object with custom detail panel functions
 * @returns The appropriate detail panel function
 */
export const getDetailPanel = <T extends Record<string, unknown>>(
  viewName: string,
  facet: string,
  customPanels: Record<string, (rowData: T | null) => React.ReactNode> = {},
) => {
  // Check for a custom panel for this specific facet
  const customPanelKey = `${viewName}.${facet}`;
  if (customPanels[customPanelKey]) {
    return customPanels[customPanelKey];
  }

  // Check for a view-wide custom panel
  if (customPanels[viewName]) {
    return customPanels[viewName];
  }

  // Fall back to the default panel
  return createDefaultDetailPanel<T>();
};
