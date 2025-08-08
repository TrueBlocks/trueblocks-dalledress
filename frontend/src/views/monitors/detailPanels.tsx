import { types } from '@models';
import { getDisplayAddress } from '@utils';

/**
 * Custom detail panel for Monitor rows
 * Displays monitor-specific information in a structured format
 */
export const renderMonitorDetailPanel = (
  rowData: Record<string, unknown> | null,
) => {
  if (!rowData) return null;
  const monitor = rowData as unknown as types.Monitor;
  return (
    <div>
      <h4>Monitor Details</h4>
      <div>
        <div>
          <strong>Name:</strong> {monitor.name}
        </div>
        <div>
          <strong>Address:</strong> {getDisplayAddress(monitor.address)}
        </div>
        <div>
          <strong>Records:</strong> {monitor.nRecords}
        </div>
        <div>
          <strong>File Size:</strong> {monitor.fileSize}
        </div>
        <div>
          <strong>Last Scanned:</strong> {monitor.lastScanned}
        </div>
      </div>
    </div>
  );
};

/**
 * Configuration object for monitors view detail panels
 * Maps data facets to their corresponding detail panel functions
 */
export const monitorsDetailPanels = {
  'monitors.monitors': renderMonitorDetailPanel,
};
