import {
  ReactNode,
  createContext,
  useCallback,
  useContext,
  useState,
} from 'react';

import { TableKey, tableKeyToString } from '.';

// Pagination interfaces
export interface PaginationState {
  currentPage: number;
  pageSize: number;
  totalItems: number;
}

// Create stable reference for initial state to prevent new object creation
export const initialPaginationState: PaginationState = Object.freeze({
  currentPage: 0,
  pageSize: 10,
  totalItems: 0,
});

export interface ViewPaginationState {
  [key: string]: PaginationState;
}

interface ViewContextType {
  currentView: string;
  setCurrentView: (view: string) => void;
  viewPagination: ViewPaginationState;
  getPagination: (tableKey: TableKey) => PaginationState;
  updatePagination: (
    tableKey: TableKey,
    changes: Partial<PaginationState>,
  ) => void;
}

export const ViewContext = createContext<ViewContextType>({
  currentView: '',
  setCurrentView: () => {},
  viewPagination: {},
  getPagination: () => initialPaginationState,
  updatePagination: () => {},
});

export const ViewContextProvider = ({ children }: { children: ReactNode }) => {
  const [currentView, setCurrentView] = useState('');
  const [viewPagination, setViewPagination] = useState<ViewPaginationState>({});

  const getPagination = useCallback(
    (tableKey: TableKey) => {
      const key = tableKeyToString(tableKey);
      return viewPagination[key] || initialPaginationState;
    },
    [viewPagination],
  );

  const updatePagination = useCallback(
    (tableKey: TableKey, changes: Partial<PaginationState>) => {
      setViewPagination((prev) => {
        const key = tableKeyToString(tableKey);
        const currentPagination = prev[key] || { ...initialPaginationState };
        return {
          ...prev,
          [key]: {
            ...currentPagination,
            ...changes,
          },
        };
      });
    },
    [],
  );

  return (
    <ViewContext.Provider
      value={{
        currentView,
        setCurrentView,
        viewPagination,
        getPagination,
        updatePagination,
      }}
    >
      {children}
    </ViewContext.Provider>
  );
};

export const useViewContext = () => {
  const context = useContext(ViewContext);
  if (!context) {
    throw new Error('useViewContext must be used within a ViewContextProvider');
  }
  return context;
};
