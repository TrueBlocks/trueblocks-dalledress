import {
  ReactNode,
  createContext,
  useCallback,
  useContext,
  useState,
} from 'react';

import { sorting } from '@models';

import { TableKey, tableKeyToString } from '.';

// Pagination interfaces
export interface PaginationState {
  currentPage: number;
  pageSize: number;
  totalItems: number;
}

// Sorting and filtering interfaces
export interface SortingState {
  [key: string]: sorting.SortDef | null; // keyed by tableKeyToString(tableKey)
}

export interface FilteringState {
  [key: string]: string; // keyed by tableKey.viewName
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
  viewSorting: SortingState;
  getSorting: (tableKey: TableKey) => sorting.SortDef | null;
  updateSorting: (tableKey: TableKey, sort: sorting.SortDef | null) => void;
  viewFiltering: FilteringState;
  getFiltering: (tableKey: TableKey) => string;
  updateFiltering: (tableKey: TableKey, filter: string) => void;
}

export const ViewContext = createContext<ViewContextType>({
  currentView: '',
  setCurrentView: () => {},
  viewPagination: {},
  getPagination: () => initialPaginationState,
  updatePagination: () => {},
  viewSorting: {},
  getSorting: () => null,
  updateSorting: () => {},
  viewFiltering: {},
  getFiltering: () => '',
  updateFiltering: () => {},
});

export const ViewContextProvider = ({ children }: { children: ReactNode }) => {
  const [currentView, setCurrentView] = useState('');
  const [viewPagination, setViewPagination] = useState<ViewPaginationState>({});
  const [viewSorting, setViewSorting] = useState<SortingState>({});
  const [viewFiltering, setViewFiltering] = useState<FilteringState>({});

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

  const getSorting = useCallback(
    (tableKey: TableKey) => {
      const key = tableKeyToString(tableKey);
      return viewSorting[key] || null;
    },
    [viewSorting],
  );

  const updateSorting = useCallback(
    (tableKey: TableKey, sort: sorting.SortDef | null) => {
      setViewSorting((prev) => {
        const key = tableKeyToString(tableKey);
        return {
          ...prev,
          [key]: sort,
        };
      });
    },
    [],
  );

  const getFiltering = useCallback(
    (tableKey: TableKey) => {
      return viewFiltering[tableKey.viewName] || '';
    },
    [viewFiltering],
  );

  const updateFiltering = useCallback((tableKey: TableKey, filter: string) => {
    setViewFiltering((prev) => ({
      ...prev,
      [tableKey.viewName]: filter,
    }));
  }, []);

  return (
    <ViewContext.Provider
      value={{
        currentView,
        setCurrentView,
        viewPagination,
        getPagination,
        updatePagination,
        viewSorting,
        getSorting,
        updateSorting,
        viewFiltering,
        getFiltering,
        updateFiltering,
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

// Hook for sorting state (per-tab)
export const useSorting = (tableKey: TableKey) => {
  const { getSorting, updateSorting } = useViewContext();

  const sort = getSorting(tableKey);
  const setSorting = useCallback(
    (sort: sorting.SortDef | null) => {
      updateSorting(tableKey, sort);
    },
    [tableKey, updateSorting],
  );

  return { sort, setSorting };
};

// Hook for filtering state (per-view)
export const useFiltering = (tableKey: TableKey) => {
  const { getFiltering, updateFiltering } = useViewContext();

  const filter = getFiltering(tableKey);
  const setFiltering = useCallback(
    (filter: string) => {
      updateFiltering(tableKey, filter);
    },
    [tableKey, updateFiltering],
  );

  return { filter, setFiltering };
};
