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
  const tagItemsRef = useRef<(HTMLLIElement | null)[]>([]);

  // Initialize or resize the refs array when tags change
  useEffect(() => {
    tagItemsRef.current = tagItemsRef.current.slice(0, tags.length);
    while (tagItemsRef.current.length < tags.length) {
      tagItemsRef.current.push(null);
    }
  }, [tags]);

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

  // Scroll the focused item into view
  useEffect(() => {
    if (focusedIndex >= 0 && focusedIndex < tags.length) {
      const focusedElement = tagItemsRef.current[focusedIndex];
      if (
        focusedElement &&
        typeof focusedElement.scrollIntoView === 'function'
      ) {
        // Use scrollIntoView with options to ensure smooth scrolling and proper positioning
        try {
          focusedElement.scrollIntoView({
            behavior: 'smooth',
            block: 'nearest',
          });
        } catch (e) {
          // Handle potential issues in test environments
          console.error(`Error scrolling: ${e}`);
        }
      }
    }
  }, [focusedIndex, tags.length]);

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
      if (index !== -1) {
        setFocusedIndex(index);
      } else {
        setFocusedIndex(-1);
      }
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
                ref={(el) => {
                  tagItemsRef.current[index] = el;
                  return undefined;
                }}
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
