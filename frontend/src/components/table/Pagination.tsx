import './Pagination.css';

// PaginationProps defines the props for the Pagination component.
interface PaginationProps {
  totalPages: number;
  currentPage: number;
  handlePageChange: (page: number) => void;
  focusControls: () => void;
}

// Pagination renders pagination controls for navigating between pages in the table.
export const Pagination = ({
  totalPages,
  currentPage,
  handlePageChange,
  focusControls,
}: PaginationProps) => {
  const maxButtons = 5;
  let pageButtons = [];
  let start = 0;
  let end = maxButtons - 1;

  // Special case: if totalPages <= 1, all buttons should be disabled
  const allDisabled = totalPages <= 1;

  if (totalPages > maxButtons) {
    if (currentPage <= 2) {
      start = 0;
      end = maxButtons - 1;
    } else if (currentPage >= totalPages - 3) {
      start = totalPages - maxButtons;
      end = totalPages - 1;
    } else {
      start = currentPage - 2;
      end = currentPage + 2;
    }
  } else {
    start = 0;
    end = maxButtons - 1;
  }

  for (let i = start; i <= end; i++) {
    const pageNum = i + 1;
    const exists = i < totalPages && totalPages > 0;
    pageButtons.push(
      <button
        key={i}
        onClick={exists && !allDisabled ? () => handlePageChange(i) : undefined}
        className={i === currentPage ? 'active' : ''}
        aria-current={i === currentPage ? 'page' : undefined}
        onFocus={focusControls}
        disabled={!exists || allDisabled}
      >
        {pageNum}
      </button>,
    );
  }

  return (
    <div className="pagination">
      <button
        onClick={!allDisabled ? () => handlePageChange(0) : undefined}
        disabled={currentPage === 0 || allDisabled}
        title="First Page"
        onFocus={focusControls}
      >
        &laquo;
      </button>
      <button
        onClick={
          !allDisabled
            ? () => handlePageChange(Math.max(0, currentPage - 10))
            : undefined
        }
        disabled={currentPage === 0 || allDisabled}
        title="Previous Page"
        onFocus={focusControls}
      >
        &lsaquo;
      </button>
      {pageButtons}
      <button
        onClick={
          !allDisabled
            ? () => handlePageChange(Math.min(totalPages - 1, currentPage + 10))
            : undefined
        }
        disabled={currentPage >= totalPages - 1 || allDisabled}
        title="Next Page"
        onFocus={focusControls}
      >
        &rsaquo;
      </button>
      <button
        onClick={
          !allDisabled ? () => handlePageChange(totalPages - 1) : undefined
        }
        disabled={currentPage >= totalPages - 1 || allDisabled}
        title="Last Page"
        onFocus={focusControls}
      >
        &raquo;
      </button>
    </div>
  );
};
