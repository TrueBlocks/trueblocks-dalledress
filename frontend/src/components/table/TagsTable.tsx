import {
  ForwardedRef,
  forwardRef,
  useEffect,
  useImperativeHandle,
  useRef,
  useState,
} from 'react';

import './Table.css';
import './TagsTable.css';

export interface TagsTableProps {
  tags: string[];
  onTagSelect: (tag: string | null, focusMainTable?: boolean) => void;
  selectedTag: string | null;
  visible: boolean;
}

export interface TagsTableHandle {
  focus: () => void;
}

export const TagsTable = forwardRef(function TagsTable(
  { tags, onTagSelect, selectedTag, visible }: TagsTableProps,
  ref: ForwardedRef<TagsTableHandle>,
) {
  const [focusedIndex, setFocusedIndex] = useState<number>(-1);
  const tagTableRef = useRef<HTMLDivElement | null>(null);

  // Expose focus method via ref
  useImperativeHandle(ref, () => ({
    focus: () => {
      tagTableRef.current?.focus();
    },
  }));

  // Handle tag selection with toggle functionality
  const handleTagSelect = (tag: string, shouldFocus: boolean = false) => {
    // Toggle tag selection if already selected
    if (selectedTag === tag) {
      onTagSelect(null, shouldFocus);
    } else {
      onTagSelect(tag, shouldFocus);
    }
  };

  // Handle keyboard navigation within the tags table
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!visible) return;

    const tagsLength = tags.length;

    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        setFocusedIndex((prev) => (prev < tagsLength - 1 ? prev + 1 : prev));
        break;
      case 'ArrowUp':
        e.preventDefault();
        setFocusedIndex((prev) => (prev > 0 ? prev - 1 : prev));
        break;
      case 'Enter':
        e.preventDefault();
        if (focusedIndex >= 0 && focusedIndex < tagsLength) {
          const tag = tags[focusedIndex];
          if (tag) {
            // Pass true to indicate focus should move to main table
            handleTagSelect(tag, true);
          }
        }
        break;
      case 'Escape':
        e.preventDefault();
        onTagSelect(null); // Clear tag selection
        break;
      default:
        break;
    }
  };

  // Clear selection when tags change
  useEffect(() => {
    setFocusedIndex(-1);
  }, [tags]);

  // Reset focused index when selected tag changes
  useEffect(() => {
    if (selectedTag) {
      const index = tags.findIndex((tag) => tag === selectedTag);
      setFocusedIndex(index !== -1 ? index : -1);
    } else {
      setFocusedIndex(-1);
    }
  }, [selectedTag, tags]);

  if (!visible) return null;

  return (
    <div className="tags-table-container">
      <div className="tags-table-header">
        <h3>Tags</h3>
      </div>
      <div
        className="tags-table"
        tabIndex={0}
        onKeyDown={handleKeyDown}
        ref={tagTableRef}
        data-testid="tags-table"
      >
        {tags.length === 0 ? (
          <div className="no-tags">No tags available</div>
        ) : (
          <ul>
            {tags.map((tag, index) => (
              <li
                key={tag}
                className={`tag-item ${selectedTag === tag ? 'selected' : ''} ${
                  focusedIndex === index ? 'focused' : ''
                }`}
                onClick={() => handleTagSelect(tag, false)}
              >
                {tag}
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
});
