/**
 * ViewStateKey identifies a specific data view/facet for UI state management.
 *
 * ## Architecture Overview
 * The key represents a hierarchical path through the UI:
 * - viewName: Top-level view ("/exports", "/names", etc.)
 * - tabName: Facet/tab within that view ("transactions", "receipts", etc.)
 *
 * ## Scope Usage Patterns
 * Different UI state types use different parts of the key:
 *
 * ### Full Key Scope (viewName + tabName)
 * Used for state that should be isolated per tab:
 * - **Pagination**: Each tab has different data and pagination needs
 * - **Sorting**: Each tab has different columns and sort requirements
 * - **Keyboard Navigation**: Tab-specific row selection and focus
 *
 * ### View-Only Scope (viewName only)
 * Used for state that should be shared across all tabs in a view:
 * - **Filtering**: Filter text applies to all tabs within a view
 * - **View-level Preferences**: Settings that affect the entire view
 *
 * ## Examples
 * ```typescript
 * // For exports view, transactions tab
 * { viewName: "/exports", tabName: "transactions" }
 *
 * // For names view, entities tab
 * { viewName: "/names", tabName: "entities" }
 * ```
 *
 * ## String Serialization
 * Keys are converted to strings for use in state maps:
 * `viewStateKeyToString({ viewName: "/exports", tabName: "transactions" })`
 * → `"/exports/transactions/"`
 */
export interface ViewStateKey {
  /** The top-level view identifier (e.g., "/exports", "/names") */
  viewName: string;
  /** The facet/tab identifier within the view (e.g., "transactions", "receipts") */
  tabName: string;
}

/**
 * Converts a ViewStateKey to a string for use in state storage.
 *
 * @param key - The ViewStateKey to convert
 * @returns Formatted string: "viewName/tabName/"
 *
 * @example
 * ```typescript
 * const key = { viewName: "/exports", tabName: "transactions" };
 * const keyString = viewStateKeyToString(key);
 * // Result: "/exports/transactions/"
 * ```
 */
export const viewStateKeyToString = (key: ViewStateKey): string => {
  return `${key.viewName}/${key.tabName}/`;
};

/**
 * Type guard to check if an object is a valid ViewStateKey.
 *
 * @param obj - Object to check
 * @returns True if object has valid viewName and tabName strings
 */
export const isViewStateKey = (obj: unknown): obj is ViewStateKey => {
  return (
    typeof obj === 'object' &&
    obj !== null &&
    typeof (obj as ViewStateKey).viewName === 'string' &&
    typeof (obj as ViewStateKey).tabName === 'string'
  );
};

/**
 * Creates a ViewStateKey with validation.
 *
 * @param viewName - The view identifier
 * @param tabName - The tab identifier
 * @returns Validated ViewStateKey
 * @throws Error if either parameter is empty or invalid
 */
export const createViewStateKey = (
  viewName: string,
  tabName: string,
): ViewStateKey => {
  if (!viewName?.trim()) {
    throw new Error('ViewStateKey viewName cannot be empty');
  }
  if (!tabName?.trim()) {
    throw new Error('ViewStateKey tabName cannot be empty');
  }

  return { viewName: viewName.trim(), tabName: tabName.trim() };
};
