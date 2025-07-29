import { useCallback, useEffect, useRef, useState } from 'react';

import {
  GetOpenProjects,
  NewProject,
  OpenProjectFile,
  SaveProject,
} from '@app';
import { Action } from '@components';
import { useIconSets } from '@hooks';
import {
  Button,
  Group,
  Modal,
  Paper,
  Stack,
  Text,
  TextInput,
  Title,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { Log } from '@utils';

interface ProjectSelectionModalProps {
  opened: boolean;
  onProjectSelected: () => void;
}

interface RecentProject {
  name: string;
  path: string;
  last_opened: string;
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
  const [error, setError] = useState<string | null>(null);
  const { File, Add } = useIconSets();
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

  const loadRecentProjects = useCallback(async () => {
    try {
      const projects = await GetOpenProjects();
      // Convert to our interface and sort by last_opened
      const getTime = (dateStr: string) => {
        const d = new Date(dateStr);
        return isNaN(d.getTime()) ? 0 : d.getTime();
      };
      const recentProjects: RecentProject[] = projects
        .map((p: Record<string, string>) => ({
          name: p.name || '',
          path: p.path || '',
          last_opened: p.last_opened || '',
        }))
        .sort((a, b) => getTime(b.last_opened) - getTime(a.last_opened))
        .slice(0, 10); // Show only 10 most recent
      if (isMounted.current) {
        setRecentProjects(recentProjects);
      }
    } catch (err) {
      Log(`Failed to load recent projects: ${err}`);
      setRecentProjects([]);
    }
  }, []);

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
                  <Stack gap="xs">
                    {recentProjects.map((project, index) => (
                      <Group
                        key={index}
                        gap="xs"
                        style={{
                          padding: '8px',
                          borderRadius: '4px',
                        }}
                      >
                        <File size={16} />
                        <div style={{ flex: 1 }}>
                          <Text size="sm" fw={500}>
                            {project.name}
                          </Text>
                          <Text size="xs" c="dimmed">
                            {new Date(project.last_opened).toLocaleDateString()}
                          </Text>
                        </div>
                        <Action
                          icon="Switch"
                          size="sm"
                          variant="light"
                          onClick={() => handleOpenProject(project.path)}
                          disabled={isLoading}
                        />
                      </Group>
                    ))}
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
