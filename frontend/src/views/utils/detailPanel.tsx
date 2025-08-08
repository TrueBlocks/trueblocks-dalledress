/* eslint-disable @typescript-eslint/no-explicit-any */

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
