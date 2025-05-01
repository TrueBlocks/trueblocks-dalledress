import { useHotkeys } from 'react-hotkeys-hook';

interface UseFormHotkeysProps {
  mode?: 'display' | 'edit';
  setMode?: (mode: 'display' | 'edit') => void;
  onCancel?: () => void;
  submitButtonRef?: React.RefObject<HTMLButtonElement | null>;
}

export const useFormHotkeys = ({
  mode,
  setMode,
  onCancel,
}: UseFormHotkeysProps): void => {
  useHotkeys(
    'mod+a',
    (e) => {
      e.preventDefault();
      const activeElement = document.activeElement as HTMLInputElement;
      if (activeElement && activeElement.tagName.toLowerCase() === 'input') {
        const nativeInput = activeElement;
        nativeInput.setSelectionRange(0, nativeInput.value.length);
      }
    },
    { enableOnFormTags: true },
  );

  useHotkeys(
    'enter',
    (e) => {
      const activeElement = document.activeElement as HTMLInputElement;

      if (mode === 'display') {
        e.preventDefault();
        setMode?.('edit');
        return;
      }

      e.preventDefault();
      const form = activeElement.closest('form') as HTMLFormElement;
      if (form) {
        const submitButton = form.querySelector(
          'button[type="submit"]',
        ) as HTMLButtonElement;
        if (submitButton) {
          submitButton.click();
        }
      }
    },
    { enableOnFormTags: true },
  );

  useHotkeys(
    'esc',
    (e) => {
      e.preventDefault();
      setMode?.('display');
      onCancel?.();
    },
    { enableOnFormTags: true },
  );
};
