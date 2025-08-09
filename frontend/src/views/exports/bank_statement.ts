// Legacy Bank Statement Renderer (preserved for reference)
//
// This code was previously used for custom rendering of Exports statements.
// It is now commented out and not used in production.
//
// To revive, uncomment and integrate as needed.

/*
import { types } from '@models';
import { getDisplayAddress } from '@utils';

/ * *
 * Custom bank statement-style detail panel for Statement rows
 * Displays transaction data in a professional bank statement format with ether values
 * /
export const renderStatementDetailPanel = (
  rowData: Record<string, unknown> | null,
) => {
  if (!rowData) return null;
  const statement = rowData as unknown as types.Statement;

  // Helper function to format values as ether
  const formatValue = (value: unknown) => {
    if (!value || value === '0') return '-';

    // Convert wei to ether using the statement's decimals
    const valueStr = typeof value === 'string' ? value : String(value);
    const decimals = statement.decimals || 18; // Default to 18 for ETH

    try {
      // Simple conversion: divide by 10^decimals
      const divisor = Math.pow(10, decimals);
      const etherValue = parseFloat(valueStr) / divisor;

      // Format to 3 decimal places
      return etherValue.toFixed(3);
    } catch {
      return valueStr; // Fallback to original value if conversion fails
    }
  };

  // Check if reconciled (no correcting reasons means it's reconciled)
  const isReconciled =
    !statement.correctingReasons || statement.correctingReasons === '';

  return (
    <div
      className="detail-table fixed-prompt-width"
      style={{
        width: '100%',
        fontFamily: 'monospace',
        fontSize: 14,
        lineHeight: 1.6,
        background: 'var(--mantine-color-gray-0)',
        borderRadius: 6,
        border: '1px solid var(--mantine-color-gray-4)',
        boxShadow: 'var(--mantine-shadow-xs)',
        margin: 0,
        padding: 0,
      }}
    >
      <div
        style={{
          textAlign: 'center',
          marginBottom: 16,
          borderBottom: '1px solid var(--mantine-color-gray-4)',
          padding: '12px 0 8px 0',
        }}
      >
        <span style={{ fontWeight: 600, fontSize: 18 }}>Bank Statement</span>
        <div
          style={{
            fontSize: 12,
            color: 'var(--mantine-color-dimmed)',
            marginTop: 2,
          }}
        >
          Block: {statement.blockNumber} | Tx: {statement.transactionIndex} |
          Date: {new Date(statement.timestamp * 1000).toLocaleDateString()}
        </div>
      </div>
      <div
        style={{
          display: 'flex',
          gap: 32,
          margin: '0 0 12px 0',
          padding: '0 16px',
        }}
      >
        <div>
          <b>Account:</b> {getDisplayAddress(statement.accountedFor)}
        </div>
        <div>
          <b>Asset:</b> {statement.symbol} ({getDisplayAddress(statement.asset)}
          )
        </div>

        <div>
          <b>From:</b> {getDisplayAddress(statement.sender)}
        </div>
        <div>
          <b>To:</b> {getDisplayAddress(statement.recipient)}
        </div>
      </div>
      <table
        className="detail-table"
        style={{ width: '100%', borderCollapse: 'collapse', marginBottom: 12 }}
      >
        <thead>
          <tr style={{ borderBottom: '1px solid var(--mantine-color-gray-4)' }}>
            <th style={{ textAlign: 'left', padding: '5px' }}>Type</th>
            <th style={{ textAlign: 'right', padding: '5px' }}>In</th>
            <th style={{ textAlign: 'right', padding: '5px' }}>Out</th>
            <th style={{ textAlign: 'right', padding: '5px' }}>Gas</th>
            <th style={{ textAlign: 'right', padding: '5px' }}>Balance</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td style={{ padding: '5px' }}>Beginning Balance</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>
              {formatValue(statement.begBal)}
            </td>
          </tr>
          <tr>
            <td style={{ padding: '5px' }}>Regular Amount</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>
              {formatValue(statement.amountIn)}
            </td>
            <td style={{ textAlign: 'right', padding: '5px' }}>
              {formatValue(statement.amountOut)}
            </td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
          </tr>
          <tr>
            <td style={{ padding: '5px' }}>Internal Transfers</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>
              {formatValue(statement.internalIn)}
            </td>
            <td style={{ textAlign: 'right', padding: '5px' }}>
              {formatValue(statement.internalOut)}
            </td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
          </tr>
          <tr>
            <td style={{ padding: '5px' }}>Self Destruct</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>
              {formatValue(statement.selfDestructIn)}
            </td>
            <td style={{ textAlign: 'right', padding: '5px' }}>
              {formatValue(statement.selfDestructOut)}
            </td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
          </tr>
          <tr>
            <td style={{ padding: '5px' }}>Gas Costs</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>
              {formatValue(statement.gasOut)}
            </td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
          </tr>
          <tr>
            <td style={{ padding: '5px' }}>Miner Rewards</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>
              {formatValue(statement.minerBaseRewardIn)}
            </td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
          </tr>
          <tr
            style={{
              borderTop: '1px solid var(--mantine-color-gray-4)',
              fontWeight: 'bold',
            }}
          >
            <td style={{ padding: '5px' }}>Ending Balance</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>
              {formatValue(statement.endBal)}
            </td>
          </tr>
        </tbody>
      </table>
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          margin: '8px 0 12px 0',
          fontSize: 14,
        }}
      >
        <span style={{ marginRight: 10 }}>Reconciled:</span>
        {isReconciled ? (
          <span style={{ color: 'var(--mantine-color-green-6)', fontSize: 18 }}>
            ✓
          </span>
        ) : (
          <span style={{ color: 'var(--mantine-color-red-6)', fontSize: 18 }}>
            ✗
          </span>
        )}
        {!isReconciled && statement.correctingReasons && (
          <span
            style={{
              marginLeft: 10,
              fontSize: 12,
              color: 'var(--mantine-color-dimmed)',
            }}
          >
            ({statement.correctingReasons})
          </span>
        )}
      </div>
    </div>
  );
};

/ * *
 * Configuration object for exports view detail panels
 * Maps data facets to their corresponding detail panel functions
 * /
export const exportsDetailPanels = {
  'exports.statements': renderStatementDetailPanel,
};
*/
