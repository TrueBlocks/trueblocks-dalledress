import { ReactNode, createContext, useState } from 'react';

interface ViewContextType {
  currentView: string;
  setCurrentView: (view: string) => void;
}

export const ViewContext = createContext<ViewContextType>({
  currentView: '',
  setCurrentView: () => {},
});

export const ViewContextProvider = ({ children }: { children: ReactNode }) => {
  const [currentView, setCurrentView] = useState('');
  return (
    <ViewContext.Provider value={{ currentView, setCurrentView }}>
      {children}
    </ViewContext.Provider>
  );
};
