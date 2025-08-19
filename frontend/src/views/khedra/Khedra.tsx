import { useCallback, useEffect, useState } from 'react';

import { GetKhedraControlURL } from '@app';
import { useHotkeys } from '@mantine/hooks';
import { Log } from '@utils';

const KhedraControl = () => {
  const [dashboardURL, setDashboardURL] = useState<string>('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [refreshKey, setRefreshKey] = useState(0);

  const getDashboardURL = useCallback(async () => {
    try {
      setLoading(true);
      setError('');

      Log('Getting Khedra dashboard URL...');

      // Get the dashboard URL from the backend (reads control.json file)
      const url = await GetKhedraControlURL();
      Log(`Got dashboard URL: ${url}`);

      setDashboardURL(url);
    } catch (err) {
      const errorMessage =
        err instanceof Error
          ? err.message
          : 'Failed to get Khedra dashboard URL';
      Log(`Error getting dashboard URL: ${errorMessage}`);
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  }, []);

  const refreshDashboard = useCallback(() => {
    Log('Refreshing Khedra dashboard...');
    setRefreshKey((prev) => prev + 1);
  }, []);

  const handleReload = useCallback(() => {
    Log('Reloading Khedra dashboard...');
    refreshDashboard();
  }, [refreshDashboard]);

  useEffect(() => {
    getDashboardURL();
  }, [getDashboardURL]);

  useHotkeys([['mod+r', handleReload]]);

  if (loading) {
    return <div>Loading Khedra dashboard...</div>;
  }

  if (error) {
    return (
      <div>
        <div>Error: {error}</div>
        <button onClick={getDashboardURL}>Retry</button>
      </div>
    );
  }

  if (!dashboardURL) {
    return <div>No dashboard URL available</div>;
  }

  return (
    <iframe
      key={refreshKey}
      src={dashboardURL}
      style={{
        width: '100%',
        height: '100vh',
        border: 'none',
      }}
      title="Khedra Dashboard"
    />
  );
};

export const Khedra = () => {
  return <KhedraControl />;
};
