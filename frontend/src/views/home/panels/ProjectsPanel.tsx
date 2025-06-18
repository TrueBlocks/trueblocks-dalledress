import { useEffect, useState } from 'react';

import { GetOpenProjects } from '@app';
import { DashboardCard, StatusIndicator } from '@components';
import { useEvent, useIcons } from '@hooks';
import { Button, Group, List, Text, ThemeIcon } from '@mantine/core';
import { msgs } from '@models';
import { Log } from '@utils';

interface Project {
  id: string;
  name: string;
  path: string;
  isActive: boolean;
  isDirty: boolean;
  lastOpened: string;
  createdAt: string;
}

interface ProjectsPanelProps {
  onViewAll?: () => void;
  onNewProject?: () => void;
}

export const ProjectsPanel = ({
  onViewAll,
  onNewProject,
}: ProjectsPanelProps) => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { File } = useIcons();

  const fetchProjects = async (showLoading = true) => {
    try {
      if (showLoading) setLoading(true);
      const openProjects = await GetOpenProjects();
      const typedProjects = openProjects as unknown as Project[];
      const sorted = typedProjects
        .sort(
          (a, b) =>
            new Date(b.lastOpened).getTime() - new Date(a.lastOpened).getTime(),
        )
        .slice(0, 3);
      setProjects(sorted);
      setError(null);
    } catch (err) {
      Log(`Error fetching projects: ${err}`);
      setError('Failed to load projects');
    } finally {
      if (showLoading) setLoading(false);
    }
  };

  useEffect(() => {
    fetchProjects();
  }, []);

  // Listen for project changes and refresh silently (no loading state)
  useEvent(msgs.EventType.MANAGER, () => {
    fetchProjects(false);
  });

  const activeProject = projects.find((p) => p.isActive);
  const recentProjects = projects.filter((p) => !p.isActive);

  return (
    <DashboardCard
      title="Projects"
      subtitle={`${projects.length} open`}
      icon={<File size={20} />}
      loading={loading}
      error={error}
      onClick={onViewAll}
    >
      <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
        {activeProject && (
          <div>
            <StatusIndicator status="healthy" label="Active" size="xs" />
            <Text size="sm" fw={500} mt={4}>
              {activeProject.name}
              {activeProject.isDirty && ' *'}
            </Text>
          </div>
        )}

        {recentProjects.length > 0 && (
          <div>
            <Text size="xs" c="dimmed" mb={4}>
              Recent
            </Text>
            <List size="sm" spacing={2}>
              {recentProjects.map((project) => (
                <List.Item
                  key={project.id}
                  icon={
                    <ThemeIcon size={16} radius="xl" color="gray">
                      <File size={10} />
                    </ThemeIcon>
                  }
                >
                  <Text size="xs">
                    {project.name}
                    {project.isDirty && ' *'}
                  </Text>
                </List.Item>
              ))}
            </List>
          </div>
        )}

        <Group gap="xs" mt="auto">
          <Button size="xs" variant="light" onClick={onNewProject}>
            New Project
          </Button>
          {projects.length > 0 && (
            <Button size="xs" variant="outline" onClick={onViewAll}>
              View All
            </Button>
          )}
        </Group>
      </div>
    </DashboardCard>
  );
};
