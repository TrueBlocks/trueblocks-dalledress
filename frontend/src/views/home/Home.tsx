import { Suspense, useEffect, useState } from 'react';

import { GetAppId } from '@app';
import { LazyPanel, PanelSkeleton } from '@components';
import { Container, Grid, Space, Text, Title } from '@mantine/core';
import { useLocation } from 'wouter';

import styles from './Home.module.css';
import {
  LazyAbisPanel,
  LazyChunksPanel,
  LazyExportsPanel,
  LazyMonitorsPanel,
  LazyNamesPanel,
  LazyProjectsPanel,
} from './LazyPanels';
import { QuickActions } from './QuickActions';
import { RecentActivity } from './RecentActivity';
import { SampleImageSection } from './SampleImageSection';
import { useUploadAddressesDialog } from './UploadAddressesDialog';

export const Home = () => {
  const [_, setAppName] = useState('Your App');
  const [, navigate] = useLocation();
  const uploadDialog = useUploadAddressesDialog();

  useEffect(() => {
    GetAppId().then((id) => {
      setAppName(id.appName);
    });
  }, []);

  // Mock recent activities for now
  const mockActivities = [
    {
      id: '1',
      type: 'image' as const,
      message: 'Generated new image from address 0x1234...5678',
      timestamp: new Date(Date.now() - 1000 * 60 * 30), // 30 minutes ago
      onClick: () => navigate('/dalledress'),
    },
    {
      id: '2',
      type: 'monitor' as const,
      message: 'Added monitor for 0xabcd...efgh',
      timestamp: new Date(Date.now() - 1000 * 60 * 60 * 2), // 2 hours ago
      onClick: () => navigate('/monitors'),
    },
    {
      id: '3',
      type: 'export' as const,
      message: 'Exported transaction data as CSV',
      timestamp: new Date(Date.now() - 1000 * 60 * 60 * 24), // 1 day ago
      onClick: () => navigate('/exports'),
    },
  ];

  const handleViewProjects = () => {
    // Navigate to home page (projects are shown there)
    navigate('/');
  };

  const handleViewNames = () => {
    navigate('/names');
  };

  const handleViewMonitors = () => {
    navigate('/monitors');
  };

  const handleViewExports = () => {
    navigate('/exports');
  };

  const handleViewChunks = () => {
    navigate('/chunks');
  };

  const handleViewAbis = () => {
    navigate('/abis');
  };

  const handleViewGallery = () => {
    // TODO: Implement gallery view or navigate to output folder
    console.log('Navigate to gallery');
  };

  const handleGenerateImage = () => {
    navigate('/dalledress');
  };

  const handleUploadAddresses = () => {
    uploadDialog.open();
  };

  const handleAddressesUpload = async (addresses: string[]) => {
    // TODO: Implement actual address upload to monitors
    console.log('Uploading addresses:', addresses);
    // This would typically call an API to add multiple monitors
    // For now, we'll just log and navigate to monitors page
    navigate('/monitors');
  };

  return (
    <Container size="xl" py="xl" className={styles.homeContainer}>
      <Space h="md" />

      {/* Header */}
      <header role="banner">
        <Title order={1} mb="xs">
          TrueBlocks DalleDress
        </Title>
        <Text size="lg" c="dimmed" mb="xl">
          Generate stunning images from Ethereum addresses, block hashes, or
          transaction hashes.
        </Text>
      </header>

      {/* Main Dashboard Grid */}
      <main role="main" aria-label="Dashboard Overview">
        <Grid className={styles.dashboardGrid}>
          {/* Left Column - Data Panels (2/3 width) */}
          <Grid.Col span={{ base: 12, md: 8 }}>
            <section aria-label="Data Overview Panels" aria-live="polite">
              <Grid>
                {/* Row 1: Projects and Names */}
                <Grid.Col span={{ base: 12, sm: 6 }}>
                  <div className={styles.panelCard}>
                    <LazyPanel priority="high">
                      <Suspense fallback={<PanelSkeleton title="Projects" />}>
                        <LazyProjectsPanel onViewAll={handleViewProjects} />
                      </Suspense>
                    </LazyPanel>
                  </div>
                </Grid.Col>
                <Grid.Col span={{ base: 12, sm: 6 }}>
                  <div className={styles.panelCard}>
                    <LazyPanel priority="high">
                      <Suspense fallback={<PanelSkeleton title="Names" />}>
                        <LazyNamesPanel onViewAll={handleViewNames} />
                      </Suspense>
                    </LazyPanel>
                  </div>
                </Grid.Col>
                {/* Row 2: Monitors and Exports */}
                <Grid.Col span={{ base: 12, sm: 6 }}>
                  <div className={styles.panelCard}>
                    <LazyPanel priority="normal">
                      <Suspense fallback={<PanelSkeleton title="Monitors" />}>
                        <LazyMonitorsPanel
                          onViewAll={handleViewMonitors}
                          onAddMonitor={() => console.log('Add monitor')}
                        />
                      </Suspense>
                    </LazyPanel>
                  </div>
                </Grid.Col>
                <Grid.Col span={{ base: 12, sm: 6 }}>
                  <div className={styles.panelCard}>
                    <LazyPanel priority="normal">
                      <Suspense fallback={<PanelSkeleton title="Exports" />}>
                        <LazyExportsPanel onViewAll={handleViewExports} />
                      </Suspense>
                    </LazyPanel>
                  </div>
                </Grid.Col>
                {/* Row 3: Chunks and ABIs */}
                <Grid.Col span={{ base: 12, sm: 6 }}>
                  <div className={styles.panelCard}>
                    <LazyPanel priority="low">
                      <Suspense fallback={<PanelSkeleton title="Chunks" />}>
                        <LazyChunksPanel onViewAll={handleViewChunks} />
                      </Suspense>
                    </LazyPanel>
                  </div>
                </Grid.Col>
                <Grid.Col span={{ base: 12, sm: 6 }}>
                  <div className={styles.panelCard}>
                    <LazyPanel priority="low">
                      <Suspense fallback={<PanelSkeleton title="ABIs" />}>
                        <LazyAbisPanel onViewAll={handleViewAbis} />
                      </Suspense>
                    </LazyPanel>
                  </div>
                </Grid.Col>{' '}
              </Grid>
            </section>
          </Grid.Col>

          {/* Right Column - Image and Actions (1/3 width) */}
          <Grid.Col span={{ base: 12, md: 4 }}>
            <aside
              role="complementary"
              aria-label="Quick Actions and Recent Activity"
            >
              <div className={styles.rightColumn}>
                {/* Sample Image Section */}
                <section aria-label="Sample Image" role="region">
                  <div className={styles.sampleImageContainer}>
                    <SampleImageSection onViewGallery={handleViewGallery} />
                  </div>
                </section>

                {/* Quick Actions */}
                <section aria-label="Quick Actions" role="region">
                  <div className={styles.quickActionsContainer}>
                    <QuickActions
                      onGenerateImage={handleGenerateImage}
                      onViewGallery={handleViewGallery}
                      onUploadAddresses={handleUploadAddresses}
                    />
                  </div>
                </section>

                {/* Recent Activity */}
                <section aria-label="Recent Activity" role="region">
                  <div className={styles.activityContainer}>
                    <LazyPanel priority="low">
                      <RecentActivity activities={mockActivities} />
                    </LazyPanel>
                  </div>
                </section>
              </div>
            </aside>
          </Grid.Col>
        </Grid>
      </main>

      {/* Upload Addresses Dialog */}
      <uploadDialog.Dialog
        opened={uploadDialog.opened}
        onClose={uploadDialog.close}
        onUpload={handleAddressesUpload}
      />
    </Container>
  );
};
