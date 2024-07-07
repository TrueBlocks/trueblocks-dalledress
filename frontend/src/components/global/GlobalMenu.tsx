import React, { ReactNode, useEffect, useState } from 'react';
import { NavLink } from '@mantine/core';
import { IconHome, IconSettings, IconTag, IconList, IconSpider } from '@tabler/icons-react';
import { Link, useRoute } from 'wouter';
import { GetLastRoute, SetLastRoute } from '@gocode/app/App';

// StyledNavLink is a helper component that renders navigation link
// and sets its `active` prop (makes active link more visible)
type StyledNavLinkParams = {
  label: string;
  href: string;
  icon?: ReactNode;
  children?: ReactNode;
  onClick?: () => void;
  activeRoute: string;
};

function StyledNavLink(params: StyledNavLinkParams) {
  const [isActive] = useRoute(params.href);
  const isActiveRoute = params.activeRoute === params.href;
  return (
    <Link href={params.href}>
      <NavLink
        label={params.label}
        active={isActive || isActiveRoute}
        leftSection={params.icon}
        onClick={params.onClick}
      >
        {params.children}
      </NavLink>
    </Link>
  );
}

function GlobalMenu() {
  const [activeRoute, setActiveRoute] = useState('/');

  useEffect(() => {
    const lastRoute = (GetLastRoute() || '/').then((route) => {
      setActiveRoute(route);
    });
  }, []);

  const handleRouteChange = (route: string) => {
    SetLastRoute(route);
    setActiveRoute(route);
  };

  return (
    <>
      <StyledNavLink
        label="Home"
        icon={<IconHome />}
        href="/"
        onClick={() => handleRouteChange('/')}
        activeRoute={activeRoute}
      />
      <StyledNavLink
        label="Dalle"
        icon={<IconSpider />}
        href="/dalle"
        onClick={() => handleRouteChange('/dalle')}
        activeRoute={activeRoute}
      />
      <StyledNavLink
        label="Names"
        icon={<IconTag />}
        href="/names"
        onClick={() => handleRouteChange('/names')}
        activeRoute={activeRoute}
      />
      <StyledNavLink
        label="Settings"
        icon={<IconSettings />}
        href="/settings"
        onClick={() => handleRouteChange('/settings')}
        activeRoute={activeRoute}
      />
    </>
  );
}

export default GlobalMenu;
