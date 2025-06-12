import { useEmitters } from '@utils';

type EntityType = 'names' | 'monitors' | 'abis';
type ActionType =
  | 'create'
  | 'update'
  | 'delete'
  | 'undelete'
  | 'remove'
  | 'autoname'
  | 'clean'
  | 'reload';

interface EntityConfig {
  displayName: string;
  singularName: string;
  pluralName: string;
  useAddress: boolean; // whether to show "for address" in messages
}

const ENTITY_CONFIGS: Record<EntityType, EntityConfig> = {
  names: {
    displayName: 'Name',
    singularName: 'name',
    pluralName: 'names',
    useAddress: true,
  },
  monitors: {
    displayName: 'Monitor',
    singularName: 'monitor',
    pluralName: 'monitors',
    useAddress: false,
  },
  abis: {
    displayName: 'ABI',
    singularName: 'address',
    pluralName: 'ABI',
    useAddress: false,
  },
};

export const useActionMsgs = (entityType: EntityType) => {
  const { emitStatus, emitError } = useEmitters();
  const config = ENTITY_CONFIGS[entityType];

  const generateSuccessMessage = (
    action: ActionType,
    identifier?: string | number,
  ): string => {
    switch (action) {
      case 'create':
        return `${config.displayName} '${identifier}' was created successfully`;

      case 'update':
        return `${config.displayName} '${identifier}' was updated successfully`;

      case 'delete':
        if (config.useAddress && entityType === 'names') {
          return `${config.displayName} for address ${identifier} was deleted successfully`;
        }
        return `${config.displayName} ${identifier} was deleted successfully`;

      case 'undelete':
        if (config.useAddress && entityType === 'names') {
          return `${config.displayName} for address ${identifier} was undeleted successfully`;
        }
        return `${config.displayName} ${identifier} was undeleted successfully`;

      case 'remove':
        if (config.useAddress && entityType === 'names') {
          return `${config.displayName} for address ${identifier} was removed successfully`;
        }
        return `${config.displayName} ${identifier} was removed successfully`;

      case 'autoname':
        return `Address ${identifier} was auto-named successfully`;

      case 'clean':
        if (typeof identifier === 'number') {
          return `Cleaned ${identifier} ${config.singularName}(s) successfully`;
        }
        return `${config.displayName}s data cleaned successfully`;

      case 'reload':
        return `Reloaded ${config.singularName} data. Fetching fresh data...`;

      default:
        return `${config.displayName} ${action} completed successfully`;
    }
  };

  const generateFailureMessage = (
    action: ActionType,
    identifier?: string,
    error?: string,
  ): string => {
    switch (action) {
      case 'create':
        return `Failed to create ${config.singularName}: ${error}`;

      case 'update':
        return `Failed to update ${config.singularName}: ${error}`;

      case 'delete':
        if (config.useAddress && entityType === 'names') {
          return `Failed to delete ${config.singularName} for address ${identifier}: ${error}`;
        }
        return `Failed to delete ${config.singularName} ${identifier}: ${error}`;

      case 'undelete':
        if (config.useAddress && entityType === 'names') {
          return `Failed to undelete ${config.singularName} for address ${identifier}: ${error}`;
        }
        return `Failed to undelete ${config.singularName} ${identifier}: ${error}`;

      case 'remove':
        if (config.useAddress && entityType === 'names') {
          return `Failed to remove ${config.singularName} for address ${identifier}: ${error}`;
        }
        return `Failed to remove ${config.singularName} ${identifier}: ${error}`;

      case 'autoname':
        return `Failed to auto-name address ${identifier}: ${error}`;

      case 'clean':
        return `Failed to clean ${config.pluralName}: ${error}`;

      default:
        return `Failed to ${action} ${config.singularName}: ${error}`;
    }
  };

  return {
    // Direct emit methods
    emitSuccess: (action: ActionType, identifier?: string | number) => {
      emitStatus(generateSuccessMessage(action, identifier));
    },

    emitFailure: (action: ActionType, identifier?: string, error?: string) => {
      emitError(generateFailureMessage(action, identifier, error));
    },

    // Message generators (for cases where you need the string without emitting)
    success: generateSuccessMessage,
    failure: generateFailureMessage,

    // Convenience methods for common patterns
    emitReloadStatus: () => {
      emitStatus(
        `Reloaded ${config.singularName} data. Fetching fresh data...`,
      );
    },

    emitCleaningStatus: (count?: number) => {
      if (count !== undefined) {
        emitStatus(`Cleaning ${count} ${config.singularName}(s)...`);
      } else {
        emitStatus(`Cleaning all ${config.pluralName}...`);
      }
    },
  };
};
