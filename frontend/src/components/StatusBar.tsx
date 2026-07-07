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

  return (
    <div className={`${styles.statusBar} ${visible ? styles.visible : ''}`}>
      <div className={styles.inner}>
        <span className={`${styles.level} ${styles[level]}`}>{level}</span>
        <span className={styles.message}>{message}</span>
        {meta && <span className={styles.meta}>{meta}</span>}
      </div>
    </div>
  );
}
