import { FormField, Header } from '@components';
import { render, screen } from '@testing-library/react';

describe('Header', () => {
  it('renders all column headers', () => {
    const columns: FormField<any>[] = [
      { key: 'name', header: 'Name' },
      { key: 'address', header: 'Address' },
      { key: 'tags', header: 'Tags' },
      { key: 'source', header: 'Source' },
      { key: 'status', header: 'Status' },
    ];
    render(
      <table>
        <Header columns={columns} />
      </table>,
    );
    expect(screen.getByText('Name')).toBeInTheDocument();
    expect(screen.getByText('Address')).toBeInTheDocument();
    expect(screen.getByText('Tags')).toBeInTheDocument();
    expect(screen.getByText('Source')).toBeInTheDocument();
    expect(screen.getByText('Status')).toBeInTheDocument();
  });
});
