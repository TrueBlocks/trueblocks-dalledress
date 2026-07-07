import { useCallback, useEffect, useState } from 'react';
import { IconDatabase, IconHome, IconPhoto, IconSettings, IconStack2 } from '@tabler/icons-react';
import { AppLayout, useViewHotkeys, useWindowGeometry } from '@trueblocks/ui';
import { StatusBar, StatusLevel } from './components/StatusBar';
import {
  GetGenerationProgress,
  GetLastRoute,
  GetSidebarWidth,
  SaveWindowGeometry,
  SetLastRoute,
  SetSidebarWidth,
} from '../wailsjs/go/app/App';
import { WindowGetPosition, WindowGetSize } from '../wailsjs/runtime/runtime';
import { Databases } from './views/Databases';
import { Dashboard } from './views/Dashboard';
import { Images } from './views/Images';
import { Series } from './views/Series';
import { Settings } from './views/Settings';
import { app } from '../wailsjs/go/models';
import { dalle } from '../wailsjs/go/models';

const VIEWS = [
  { num: 1, id: 'dashboard', label: 'Dashboard', icon: IconHome },
  { num: 2, id: 'images', label: 'Images', icon: IconPhoto },
  { num: 3, id: 'series', label: 'Series', icon: IconStack2 },
  { num: 4, id: 'databases', label: 'Databases', icon: IconDatabase },
  { num: 5, id: 'settings', label: 'Settings', icon: IconSettings },
];

type GlobalStatus = {
  visible: boolean;
  level: StatusLevel;
  message: string;
  meta?: string;
  percent?: number;
};

type ProgressTarget = {
  series: string;
  seed: string;
};

const PHASE_LABELS: Record<string, string> = {
  setup: 'Preparing generation run',
  base_prompts: 'Selecting records and building prompts',
  enhance_prompt: 'Enhancing prompt',
  image_prep: 'Preparing image request',
  image_wait: 'Waiting for image provider',
  image_download: 'Receiving image artifact',
  annotate: 'Annotating generated image',
  failed: 'Generation failed',
  completed: 'Generation complete',
};

function statusForProgress(progress: app.GenerationProgress): GlobalStatus {
  if (progress.error) {
    return { visible: true, level: 'error', message: progress.error };
  }
  const message = progress.cacheHit
    ? 'Using cached image artifacts'
    : PHASE_LABELS[progress.phase] || 'Working';
  return {
    visible: true,
    level: progress.done ? 'success' : 'progress',
    message,
    percent: progress.percent > 0 ? progress.percent : undefined,
  };
}

export function App() {
  const [route, setRoute] = useState('dashboard');
  const [selectedImageId, setSelectedImageId] = useState('');
  const [currentImage, setCurrentImage] = useState<dalle.ImageMetadataRecord | null>(null);
  const [initialSidebarWidth, setInitialSidebarWidth] = useState(220);
  const [ready, setReady] = useState(false);
  const [globalStatus, setGlobalStatus] = useState<GlobalStatus>({
    visible: false,
    level: 'progress',
    message: '',
  });
  const [progressTarget, setProgressTarget] = useState<ProgressTarget | null>(null);

  useEffect(() => {
    Promise.all([GetLastRoute().catch(() => 'dashboard'), GetSidebarWidth().catch(() => 220)]).then(
      ([savedRoute, savedWidth]) => {
        if (VIEWS.some((v) => v.id === savedRoute)) setRoute(savedRoute);
        if (savedWidth > 0) setInitialSidebarWidth(savedWidth);
        setReady(true);
      },
    );
  }, []);

  useWindowGeometry(SaveWindowGeometry, WindowGetPosition, WindowGetSize);

  const navigate = (id: string) => {
    setRoute(id);
    SetLastRoute(id);
  };

  const showGeneratedImage = (imageId: string) => {
    setSelectedImageId(imageId);
    setProgressTarget(null);
    navigate('images');
  };

  const handleProgressStart = useCallback((series: string, seed: string) => {
    setProgressTarget({ series, seed });
  }, []);

  useEffect(() => {
    if (!progressTarget) return;
    const poll = () => {
      GetGenerationProgress(progressTarget.series, progressTarget.seed)
        .then((progress) => {
          if (progress.active) setGlobalStatus(statusForProgress(progress));
        })
        .catch((err: unknown) => {
          const message = err instanceof Error ? err.message : String(err);
          setGlobalStatus({ visible: true, level: 'error', message });
        });
    };
    poll();
    const interval = window.setInterval(poll, 750);
    return () => window.clearInterval(interval);
  }, [progressTarget]);

  useViewHotkeys({ views: VIEWS, activeView: route, onNavigate: navigate });

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (!(event.metaKey || event.ctrlKey) || event.shiftKey || event.altKey) return;
      if (event.key.toLowerCase() !== 'r') return;
      if (event.target instanceof HTMLElement) {
        const editableTags = ['INPUT', 'TEXTAREA', 'SELECT'];
        if (editableTags.includes(event.target.tagName) || event.target.isContentEditable) return;
      }
      event.preventDefault();
      window.dispatchEvent(new CustomEvent('view:refresh', { detail: route }));
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [route]);

  if (!ready) return null;

  return (
    <>
      <AppLayout
        navItems={VIEWS.map(({ id, label, icon }) => ({ id, label, icon }))}
        activeNav={route}
        onNavigate={navigate}
        initialSidebarWidth={initialSidebarWidth}
        saveSidebarWidth={SetSidebarWidth}
      >
        {route === 'dashboard' && (
          <Dashboard
            onGeneratedImage={showGeneratedImage}
            currentImage={currentImage}
            onCurrentImageChange={setCurrentImage}
            onStatusChange={setGlobalStatus}
            onProgressStart={handleProgressStart}
          />
        )}
        {route === 'images' && (
          <Images
            selectedImageId={selectedImageId}
            onCurrentImageChange={setCurrentImage}
            onStatusChange={setGlobalStatus}
            onProgressStart={handleProgressStart}
          />
        )}
        {route === 'series' && <Series />}
        {route === 'databases' && <Databases />}
        {route === 'settings' && <Settings />}
      </AppLayout>
      <StatusBar
        visible={globalStatus.visible}
        level={globalStatus.level}
        message={globalStatus.message}
        meta={globalStatus.meta}
        percent={globalStatus.percent}
      />
    </>
  );
}
