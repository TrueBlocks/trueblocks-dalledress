import { Fragment, useState } from 'react';

import './DetailTable.css';

type Row = { label: string; value: React.ReactNode };

type Section = { name: string; rows: Row[] };

interface DetailTableProps {
  sections: Section[];
  promptWidthPx?: number;
  className?: string;
  defaultCollapsedSections?: string[];
}

export const DetailTable = ({
  sections,
  promptWidthPx = 220,
  className,
  defaultCollapsedSections = [],
}: DetailTableProps) => {
  const [collapsed, setCollapsed] = useState<Set<string>>(
    new Set(defaultCollapsedSections),
  );
  const toggle = (name: string) => {
    setCollapsed((prev) => {
      const next = new Set(prev);
      if (next.has(name)) next.delete(name);
      else next.add(name);
      return next;
    });
  };
  const styleVar = {
    ['--detail-prompt-width']: `${promptWidthPx}px`,
  } as Record<string, string>;
  return (
    <table
      className={`detail-table fixed-prompt-width${className ? ` ${className}` : ''}`}
      style={styleVar}
    >
      <tbody>
        {sections.map((s) => {
          const isCollapsed = collapsed.has(s.name);
          return (
            <Fragment key={s.name}>
              <tr key={`${s.name}-header`}>
                <th
                  colSpan={2}
                  onClick={() => toggle(s.name)}
                  style={{ cursor: 'pointer' }}
                >
                  <span style={{ fontSize: '0.75em' }}>
                    {isCollapsed ? '▶' : '▼'}
                  </span>{' '}
                  {s.name}
                </th>
              </tr>
              {!isCollapsed &&
                s.rows.map((r, i) => (
                  <tr key={`${s.name}-${i}`}>
                    <td
                      className="prompt"
                      style={{ width: `${promptWidthPx}px` }}
                    >
                      {r.label}
                    </td>
                    <td>{r.value}</td>
                  </tr>
                ))}
              <tr className="separator" key={`${s.name}-separator`}>
                <td colSpan={2} />
              </tr>
            </Fragment>
          );
        })}
      </tbody>
    </table>
  );
};
