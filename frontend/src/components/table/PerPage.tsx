import './PerPage.css';

// PerPageProps defines the props for the PerPage component.
interface PerPageProps {
  pageSize: number;
  onPageSizeChange: (size: number) => void;
  focusTable: () => void;
  focusControls: () => void;
}

// PerPage renders a page size selector for the table, allowing the user to change the number of rows displayed per page.
export const PerPage = ({
  pageSize,
  onPageSizeChange,
  focusTable,
  focusControls,
}: PerPageProps) => (
  <div className="page-size-selector">
    <select
      value={pageSize}
      onChange={(e) => {
        onPageSizeChange(Number(e.target.value));
        setTimeout(focusTable, 100);
      }}
      aria-label="Items per page"
      onFocus={focusControls}
    >
      {[10, 25, 50, 100].map((size, index) => (
        <option key={size} value={size}>
          {size + (index === 0 ? ' per page' : '')}
        </option>
      ))}
    </select>
  </div>
);
