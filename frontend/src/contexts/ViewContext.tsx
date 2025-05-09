import {
  ReactNode,
  createContext,
  useCallback,
  useContext,
  useState,
} from 'react';

// Pagination interfaces
export interface PaginationState {
  currentPage: number;
  pageSize: number;
  totalItems: number;
}

export const initialPaginationState: PaginationState = {
  currentPage: 0,
  pageSize: 10,
  totalItems: 0,
};

export interface ViewPaginationState {
  [viewName: string]: {
    [tabName: string]: PaginationState;
  };
}

interface ViewContextType {
  currentView: string;
  setCurrentView: (view: string) => void;
  viewPagination: ViewPaginationState;
  getViewPagination: (viewName: string, tabName: string) => PaginationState;
  updateViewPagination: (
    viewName: string,
    tabName: string,
    changes: Partial<PaginationState>,
  ) => void;
}

export const ViewContext = createContext<ViewContextType>({
  currentView: '',
  setCurrentView: () => {},
  viewPagination: {},
  getViewPagination: () => ({ ...initialPaginationState }),
  updateViewPagination: () => {},
});

export const ViewContextProvider = ({ children }: { children: ReactNode }) => {
  const [currentView, setCurrentView] = useState('');
  const [viewPagination, setViewPagination] = useState<ViewPaginationState>({});

  const getViewPagination = useCallback(
    (viewName: string, tabName: string) => {
      return (
        viewPagination[viewName]?.[tabName] || { ...initialPaginationState }
      );
    },
    [viewPagination],
  );

  const updateViewPagination = useCallback(
    (viewName: string, tabName: string, changes: Partial<PaginationState>) => {
      setViewPagination((prev) => {
        const currentPagination = prev[viewName]?.[tabName] || {
          ...initialPaginationState,
        };
        return {
          ...prev,
          [viewName]: {
            ...(prev[viewName] || {}),
            [tabName]: {
              ...currentPagination,
              ...changes,
            },
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
        getViewPagination,
        updateViewPagination,
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
