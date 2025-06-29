import { useMemo } from 'react';

import { Action } from '@components';
import { FormField } from '@components';
import { types } from '@models';

import { ActionData, useActionConfig } from '.';

interface ColumnConfig<T extends Record<string, unknown>> {
  showActions?: boolean | ((pageData: PageDataUnion) => boolean);
  actions: ActionType[];
  getCanRemove?: (row: T, pageData?: PageDataUnion) => boolean;
}

type ActionType = 'delete' | 'undelete' | 'remove' | 'autoname';

type PageDataUnion = {
  facet: types.DataFacet;
  [key: string]: unknown;
} | null;

export function toPageDataProp<T>(pageData: T | null): PageDataUnion {
  return pageData as unknown as PageDataUnion;
}

interface ActionHandlers {
  handleToggle?: (addressStr: string) => void;
  handleRemove?: (addressStr: string) => void;
  handleAutoname?: (addressStr: string) => void;
}

export const useColumns = (
  baseColumns: FormField<Record<string, unknown>>[],
  config: ColumnConfig<Record<string, unknown>>,
  handlers: ActionHandlers,
  pageData: PageDataUnion,
  actionConfig: ReturnType<typeof useActionConfig>,
  perRowCrud: boolean = true,
) => {
  return useMemo(() => {
    // Determine if actions should be shown
    const shouldShow =
      typeof config.showActions === 'function'
        ? config.showActions(pageData)
        : (config.showActions ?? true);

    if (!shouldShow || config.actions.length === 0) {
      return actionConfig.injectActionColumn(baseColumns, () => null);
    }

    const renderActions = (actionData: ActionData) => {
      const isDeleted = actionData.isDeleted;
      const effectiveDeletedState = actionData.isProcessing
        ? !isDeleted
        : isDeleted;

      return (
        <div className="action-buttons-container">
          {config.actions.includes('delete') &&
            handlers.handleToggle && ( // undelete is permitted but ignored
              <Action
                icon="Delete"
                iconOff="Undelete"
                isOn={!effectiveDeletedState}
                onClick={() => handlers.handleToggle?.(actionData.addressStr)}
                disabled={actionData.isProcessing}
                title={effectiveDeletedState ? 'Undelete' : 'Delete'}
                size="sm"
              />
            )}
          {config.actions.includes('remove') && handlers.handleRemove && (
            <Action
              icon="Remove"
              onClick={() => handlers.handleRemove?.(actionData.addressStr)}
              disabled={
                actionData.isProcessing ||
                (perRowCrud ? !effectiveDeletedState : false)
              }
              title="Remove"
              size="sm"
            />
          )}
          {config.actions.includes('autoname') && handlers.handleAutoname && (
            <Action
              icon="Autoname"
              onClick={() => handlers.handleAutoname?.(actionData.addressStr)}
              disabled={actionData.isProcessing}
              title="Auto-generate name"
              size="sm"
            />
          )}
        </div>
      );
    };

    return actionConfig.injectActionColumn(
      baseColumns,
      renderActions,
      config.getCanRemove,
    );
  }, [baseColumns, config, handlers, pageData, actionConfig, perRowCrud]);
};

export type { ColumnConfig, ActionType, ActionHandlers };
