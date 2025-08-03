import { useCallback, useEffect, useRef, useState } from 'react';

import {
  GetOpenProjects,
  NewProject,
  OpenProjectFile,
  SaveProject,
} from '@app';
import { Action, StatusIndicator } from '@components';
import { useViewContext } from '@contexts';
import { useActiveProject2, useIconSets } from '@hooks';
import {
  Badge,
  Button,
  Card,
  Group,
  Modal,
  Paper,
  Stack,
  Text,
  TextInput,
  Title,
  Tooltip,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { types } from '@models';
import { Log } from '@utils';
import { useLocation } from 'wouter';

interface ProjectSelectionModalProps {
  opened: boolean;
  onProjectSelected: () => void;
}

interface RecentProject {
  name: string;
  path: string;
  last_opened: string;
  id?: string;
  isDirty?: boolean;
  isActive?: boolean;
}

interface ProjectContext {
  lastView: string;
  activeAddress: string;
  activeChain: string;
  activeContract: string;
  lastFacetMap: Record<string, types.DataFacet>;
}

interface NewProjectForm {
  name: string;
  address: string;
}

export const ProjectSelectionModal = ({
  opened,
  onProjectSelected,
}: ProjectSelectionModalProps) => {
  const [recentProjects, setRecentProjects] = useState<RecentProject[]>([]);
  const [loadingCreate, setLoadingCreate] = useState(false);
  const [loadingOpen, setLoadingOpen] = useState(false);
  const [_loadingContext, setLoadingContext] = useState<string | null>(null);
  const [projectContexts, setProjectContexts] = useState<
    Record<string, ProjectContext>
  >({});
  const [error, setError] = useState<string | null>(null);
  const { File, Add } = useIconSets();
  const { restoreProjectFilterStates } = useViewContext();
  const { lastView } = useActiveProject2();
  const [, navigate] = useLocation();
  const isMounted = useRef(true);

  const form = useForm<NewProjectForm>({
    initialValues: {
      name: '',
      address: '',
    },
    validate: {
      name: (value) => (!value ? 'Project name is required' : null),
      address: (value) => {
        if (!value) return 'Primary address is required';
        // Allow .eth names or standard Ethereum addresses
        if (value.endsWith('.eth')) {
          // ENS names ending with .eth are valid - backend will resolve them
          if (value.length < 5) {
            // minimum: "a.eth" = 5 chars
            return 'ENS name too short';
          }
          if (!/^[a-z0-9-]+\.eth$/i.test(value)) {
            return 'Invalid ENS name format';
          }
          return null;
        }
        // Standard Ethereum address validation
        if (!/^0x[a-fA-F0-9]{40}$/.test(value)) {
          return 'Invalid Ethereum address format (use 0x... or .eth name)';
        }
        return null;
      },
    },
  });

  const loadProjectContext = useCallback(
    async (projectId: string, _projectPath: string) => {
      if (projectContexts[projectId]) return; // Already loaded

      setLoadingContext(projectId);
      try {
        // For now, we'll use placeholder data since we need the project to be active to get its context
        // In a real implementation, we might need a backend API to get project context without activating it
        setProjectContexts((prev) => ({
          ...prev,
          [projectId]: {
            lastView: '/', // Default placeholder
            activeAddress: '', // Would need backend API to get this
            activeChain: 'mainnet', // Default placeholder
            activeContract: '',
            lastFacetMap: {},
          },
        }));
      } catch (err) {
        Log(`Failed to load context for project ${projectId}: ${err}`);
      } finally {
        setLoadingContext(null);
      }
    },
    [projectContexts],
  );

  const loadRecentProjects = useCallback(async () => {
    try {
      const projects = await GetOpenProjects();
      // Convert to our interface and sort by last_opened
      const getTime = (dateStr: string) => {
        const d = new Date(dateStr);
        return isNaN(d.getTime()) ? 0 : d.getTime();
      };
      const recentProjects: RecentProject[] = projects
        .map((p: Record<string, unknown>) => ({
          id: String(p.id || ''),
          name: String(p.name || ''),
          path: String(p.path || ''),
          last_opened: String(p.last_opened || p.lastOpened || ''),
          isDirty: Boolean(p.isDirty),
          isActive: Boolean(p.isActive),
        }))
        .sort((a, b) => getTime(b.last_opened) - getTime(a.last_opened))
        .slice(0, 10); // Show only 10 most recent

      if (isMounted.current) {
        setRecentProjects(recentProjects);

        // Load context for recent projects
        recentProjects.forEach((project) => {
          if (project.id) {
            loadProjectContext(project.id, project.path);
          }
        });
      }
    } catch (err) {
      Log(`Failed to load recent projects: ${err}`);
      setRecentProjects([]);
    }
  }, [loadProjectContext]);

  useEffect(() => {
    isMounted.current = true;
    return () => {
      isMounted.current = false;
    };
  }, []);

  useEffect(() => {
    if (opened) {
      loadRecentProjects();
    }
  }, [opened, loadRecentProjects]);

  const handleCreateProject = async (values: NewProjectForm) => {
    if (loadingCreate || loadingOpen) return;
    setLoadingCreate(true);
    setError(null);

    try {
      await NewProject(values.name, values.address);

      // Immediately save the project after creation
      await SaveProject();

      await restoreProjectFilterStates();

      const targetView = lastView || '/';
      navigate(targetView);

      onProjectSelected();
    } catch (err) {
      const errorMsg = `Failed to create project: ${err}`;
      setError(errorMsg);
    } finally {
      setLoadingCreate(false);
    }
  };

  const handleOpenProject = async (projectPath: string) => {
    if (loadingCreate || loadingOpen) return;
    setLoadingOpen(true);
    setError(null);

    try {
      await OpenProjectFile(projectPath);
      await restoreProjectFilterStates();

      const targetView = lastView || '/';
      navigate(targetView);

      onProjectSelected();
    } catch (err) {
      const errorMsg = `Failed to open project: ${err}`;
      setError(errorMsg);
    } finally {
      setLoadingOpen(false);
    }
  };

  const handleOpenFile = async () => {
    if (loadingCreate || loadingOpen) return;
    setLoadingOpen(true);
    setError(null);

    try {
      // This will trigger the native file picker
      await OpenProjectFile('');
      await restoreProjectFilterStates();

      const targetView = lastView || '/';
      navigate(targetView);

      onProjectSelected();
    } catch (err) {
      const errorMsg = `Failed to open project: ${err}`;
      setError(errorMsg);
    } finally {
      setLoadingOpen(false);
    }
  };

  const isLoading = loadingCreate || loadingOpen;

  return (
    <Modal
      opened={opened}
      onClose={() => {}} // Cannot be closed without selecting a project
      centered
      size="lg"
      withCloseButton={false}
      closeOnClickOutside={false}
      closeOnEscape={false}
      overlayProps={{
        backgroundOpacity: 0.8,
        blur: 3,
      }}
    >
      <Stack gap="xl">
        <div style={{ textAlign: 'center' }}>
          <Title order={2} mb="xs">
            Select or Create a Project
          </Title>
          <Text c="dimmed">
            TrueBlocks requires an active project with at least one address to
            continue
          </Text>
        </div>

        {error && (
          <Text c="red" size="sm" style={{ textAlign: 'center' }}>
            {error}
          </Text>
        )}

        <Group gap="xl" align="flex-start" grow>
          {/* Create New Project */}
          <Paper p="md" withBorder>
            <Stack gap="md">
              <Group gap="xs">
                <Add size={20} />
                <Title order={4}>Create New Project</Title>
              </Group>

              <form onSubmit={form.onSubmit(handleCreateProject)}>
                <Stack gap="md">
                  <TextInput
                    label="Project Name"
                    placeholder="My Analysis Project"
                    required
                    {...form.getInputProps('name')}
                  />
                  <TextInput
                    label="Primary Address"
                    placeholder="0x..."
                    required
                    {...form.getInputProps('address')}
                    description="The main Ethereum address for this project"
                  />
                  <Button
                    type="submit"
                    loading={loadingCreate}
                    disabled={isLoading}
                    leftSection={<Add size={16} />}
                    fullWidth
                  >
                    Create Project
                  </Button>
                </Stack>
              </form>
            </Stack>
          </Paper>

          {/* Open Existing Project */}
          <Paper p="md" withBorder>
            <Stack gap="md">
              <Group gap="xs">
                <File size={20} />
                <Title order={4}>Open Project</Title>
              </Group>

              <Button
                variant="outline"
                leftSection={<File size={16} />}
                onClick={handleOpenFile}
                loading={loadingOpen}
                disabled={isLoading}
                fullWidth
              >
                Browse for Project File
              </Button>

              {recentProjects.length > 0 && (
                <>
                  <Text size="sm" fw={500} mt="md">
                    Recent Projects
                  </Text>
                  <Stack gap="sm">
                    {recentProjects.map((project, index) => {
                      const context = projectContexts[project.id || ''];
                      const isRecent =
                        new Date().getTime() -
                          new Date(project.last_opened).getTime() <
                        24 * 60 * 60 * 1000; // Within 24 hours

                      return (
                        <Card key={index} withBorder padding="sm" radius="md">
                          <Stack gap="xs">
                            {/* Project Header */}
                            <Group justify="space-between" wrap="nowrap">
                              <Group gap="xs" style={{ flex: 1, minWidth: 0 }}>
                                <File size={16} />
                                <div style={{ flex: 1, minWidth: 0 }}>
                                  <Group gap="xs" wrap="nowrap">
                                    <Text
                                      size="sm"
                                      fw={project.isActive ? 600 : 500}
                                      style={{
                                        textOverflow: 'ellipsis',
                                        overflow: 'hidden',
                                        whiteSpace: 'nowrap',
                                      }}
                                    >
                                      {project.name}
                                    </Text>
                                    {project.isDirty && (
                                      <Tooltip label="Has unsaved changes">
                                        <Badge
                                          size="xs"
                                          color="orange"
                                          variant="dot"
                                        >
                                          Modified
                                        </Badge>
                                      </Tooltip>
                                    )}
                                    {project.isActive && (
                                      <Badge
                                        size="xs"
                                        color="blue"
                                        variant="light"
                                      >
                                        Active
                                      </Badge>
                                    )}
                                    {isRecent && (
                                      <Badge
                                        size="xs"
                                        color="green"
                                        variant="light"
                                      >
                                        Recent
                                      </Badge>
                                    )}
                                  </Group>
                                </div>
                              </Group>
                              <Action
                                icon="Switch"
                                size="sm"
                                variant="light"
                                onClick={() => handleOpenProject(project.path)}
                                disabled={isLoading}
                              />
                            </Group>

                            {/* Project Context Preview */}
                            {context && (
                              <Group
                                gap="md"
                                style={{
                                  fontSize: '12px',
                                  color: 'var(--mantine-color-dimmed)',
                                }}
                              >
                                {context.activeAddress && (
                                  <Group gap="xs">
                                    <Text size="xs" c="dimmed">
                                      Address:
                                    </Text>
                                    <Text
                                      size="xs"
                                      fw={500}
                                      style={{ fontFamily: 'monospace' }}
                                    >
                                      {context.activeAddress.slice(0, 8)}...
                                      {context.activeAddress.slice(-6)}
                                    </Text>
                                  </Group>
                                )}
                                {context.activeChain && (
                                  <Group gap="xs">
                                    <Text size="xs" c="dimmed">
                                      Chain:
                                    </Text>
                                    <Badge
                                      size="xs"
                                      variant="light"
                                      color="blue"
                                    >
                                      {context.activeChain}
                                    </Badge>
                                  </Group>
                                )}
                                {context.lastView &&
                                  context.lastView !== '/' && (
                                    <Group gap="xs">
                                      <Text size="xs" c="dimmed">
                                        Last View:
                                      </Text>
                                      <Text size="xs" fw={500}>
                                        {context.lastView
                                          .replace('/', '')
                                          .replace(/^./, (c) =>
                                            c.toUpperCase(),
                                          )}
                                      </Text>
                                    </Group>
                                  )}
                              </Group>
                            )}

                            {/* Project Metadata */}
                            <Group
                              justify="space-between"
                              style={{ fontSize: '11px' }}
                            >
                              <Text size="xs" c="dimmed">
                                {new Date(
                                  project.last_opened,
                                ).toLocaleDateString()}{' '}
                                at{' '}
                                {new Date(
                                  project.last_opened,
                                ).toLocaleTimeString([], {
                                  hour: '2-digit',
                                  minute: '2-digit',
                                })}
                              </Text>
                              <StatusIndicator
                                status={
                                  project.isActive
                                    ? 'healthy'
                                    : project.isDirty
                                      ? 'warning'
                                      : 'inactive'
                                }
                                label=""
                                size="xs"
                              />
                            </Group>
                          </Stack>
                        </Card>
                      );
                    })}
                  </Stack>
                </>
              )}
            </Stack>
          </Paper>
        </Group>
      </Stack>
    </Modal>
  );
};
