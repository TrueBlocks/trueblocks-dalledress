import { useEffect, useState } from 'react';
import { IconDatabase, IconHome, IconPhoto, IconSettings, IconStack2 } from '@tabler/icons-react';
import { AppLayout, useViewHotkeys, useWindowGeometry } from '@trueblocks/ui';
import {
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

const VIEWS = [
  { num: 1, id: 'dashboard', label: 'Dashboard', icon: IconHome },
  { num: 2, id: 'images', label: 'Images', icon: IconPhoto },
  { num: 3, id: 'series', label: 'Series', icon: IconStack2 },
  { num: 4, id: 'databases', label: 'Databases', icon: IconDatabase },
  { num: 5, id: 'settings', label: 'Settings', icon: IconSettings },
];

export function App() {
  const [route, setRoute] = useState('dashboard');
  const [selectedImageId, setSelectedImageId] = useState('');
  const [initialSidebarWidth, setInitialSidebarWidth] = useState(220);
  const [ready, setReady] = useState(false);

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
    navigate('images');
  };

  useViewHotkeys({ views: VIEWS, activeView: route, onNavigate: navigate });

  if (!ready) return null;

  return (
    <AppLayout
      navItems={VIEWS.map(({ id, label, icon }) => ({ id, label, icon }))}
      activeNav={route}
      onNavigate={navigate}
      initialSidebarWidth={initialSidebarWidth}
      saveSidebarWidth={SetSidebarWidth}
    >
      {route === 'dashboard' && <Dashboard onGeneratedImage={showGeneratedImage} />}
      {route === 'images' && <Images selectedImageId={selectedImageId} />}
      {route === 'series' && <Series />}
      {route === 'databases' && <Databases />}
      {route === 'settings' && <Settings />}
    </AppLayout>
  );
}
