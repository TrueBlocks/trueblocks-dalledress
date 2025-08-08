import { types } from '@models';
import { getDisplayAddress } from '@utils';

/**
 * Custom bank statement-style detail panel for Statement rows
 * Displays transaction data in a professional bank statement format with ether values
 */
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
      style={{
        fontFamily: 'monospace',
        fontSize: '14px',
        lineHeight: '1.6',
      }}
    >
      <div
        style={{
          textAlign: 'center',
          marginBottom: '20px',
          borderBottom: '2px solid #333',
          paddingBottom: '10px',
        }}
      >
        <h3>Jay&apos;s Bank Statement</h3>
        <div style={{ fontSize: '12px', color: '#666' }}>
          Block: {statement.blockNumber} | Tx: {statement.transactionIndex} |
          Date: {new Date(statement.timestamp * 1000).toLocaleDateString()}
        </div>
      </div>

      <div style={{ marginBottom: '15px' }}>
        <strong>Transaction Identification:</strong>
        <br />
        <div style={{ marginLeft: '10px' }}>
          Account: {getDisplayAddress(statement.accountedFor)}
          <br />
          Asset: {statement.symbol} ({getDisplayAddress(statement.asset)})
          <br />
          From: {getDisplayAddress(statement.sender)}
          <br />
          To: {getDisplayAddress(statement.recipient)}
        </div>
      </div>

      <table
        style={{
          width: '100%',
          borderCollapse: 'collapse',
          marginBottom: '15px',
        }}
      >
        <thead>
          <tr style={{ borderBottom: '1px solid #333' }}>
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
          <tr style={{ borderTop: '1px solid #333', fontWeight: 'bold' }}>
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
          marginTop: '10px',
        }}
      >
        <span style={{ marginRight: '10px' }}>Reconciled:</span>
        {isReconciled ? (
          <span style={{ color: 'green', fontSize: '18px' }}>✓</span>
        ) : (
          <span style={{ color: 'red', fontSize: '18px' }}>✗</span>
        )}
        {!isReconciled && statement.correctingReasons && (
          <span style={{ marginLeft: '10px', fontSize: '12px', color: '#666' }}>
            ({statement.correctingReasons})
          </span>
        )}
      </div>
    </div>
  );
};

/**
 * Configuration object for exports view detail panels
 * Maps data facets to their corresponding detail panel functions
 */
export const exportsDetailPanels = {
  'exports.statements': renderStatementDetailPanel,
};
