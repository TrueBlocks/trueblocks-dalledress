import { IconCopy } from '@tabler/icons-react';
import styles from './StatusBar.module.css';

export type StatusLevel = 'progress' | 'success' | 'error';

type StatusBarProps = {
  visible: boolean;
  level: StatusLevel;
  message: string;
  meta?: string;
};

export function StatusBar({ visible, level, message, meta = '' }: StatusBarProps) {
  if (!message) return null;

  const copyError = () => {
    const text = meta ? `${message} — ${meta}` : message;
    navigator.clipboard.writeText(text);
  };

  return (
    <div className={`${styles.statusBar} ${visible ? styles.visible : ''}`}>
      <div className={styles.inner}>
        <span className={`${styles.level} ${styles[level]}`}>{level}</span>
        <span className={styles.message}>{message}</span>
        {meta && <span className={styles.meta}>{meta}</span>}
        {level === 'error' && (
          <button className={styles.copyButton} onClick={copyError} title="Copy error">
            <IconCopy size={14} />
          </button>
        )}
      </div>
    </div>
  );
}
